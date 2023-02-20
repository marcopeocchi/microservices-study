package pkg

import (
	"embed"
	"fmt"
	"fuu/v/pkg/cli"
	config "fuu/v/pkg/config"
	"fuu/v/pkg/listing"
	"fuu/v/pkg/static"
	"fuu/v/pkg/user"
	"fuu/v/pkg/workers"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var cfg = config.Instance()

func RunBlocking(db *gorm.DB, frontend *embed.FS) {

	thumbnailer := workers.Thumbnailer{
		BaseDir:           cfg.WorkingDir,
		ImgHeight:         cfg.ThumbnailHeight,
		ImgQuality:        cfg.ThumbnailQuality,
		ForceRegeneration: cfg.ForceRegeneration,
		CacheDir:          cfg.CacheDir,
		Database:          db,
	}

	fileWatcher := workers.FileWatcher{
		WorkingDir: cfg.WorkingDir,
	}
	fileWatcher.New()

	log.Println("Starting server")

	wg := new(sync.WaitGroup)
	wg.Add(3)

	// Discourage the execution of this program as SuperUser.
	// Unless in executed docker because of obvious reasons.
	uid := os.Getuid()
	if uid == 0 {
		log.Println(cli.Yellow, "You're running this program as root (UID 0)", cli.Reset)
		log.Println(cli.Yellow, "This isn't reccomended unless you're using Docker", cli.Reset)
	}

	// ********** MAIN COMPONENTS GOROUTINES **********

	// HTTP Server
	go func() {
		server := createServer(cfg.Port, frontend, db)
		server.ListenAndServe()
		wg.Done()
	}()

	// Thumbnailer worker
	go func() {
		log.Println("Starting thumbnailer")

		start := time.Now()
		thumbnailer.Start()
		stop := time.Since(start)

		log.Println("Thumbnailer took", cli.Format(stop, cli.BgBlue))
		wg.Done()
	}()

	// Ionotify filewatcher worker
	go func() {
		defer func() {
			fileWatcher.Close()
			wg.Done()
		}()

		fileWatcher.Start(
			func(event string) { thumbnailer.Start() },
			func(event string) { thumbnailer.Remove(event) },
		)
	}()

	log.Println("Server started")

	// wait for the waitgroup to finish, which it will not.
	// effectively blocks.
	wg.Wait()
}

func createServer(port int, app *embed.FS, db *gorm.DB) *http.Server {
	r := mux.NewRouter()

	// User group router
	ur := r.PathPrefix("/user").Subrouter()
	ur.HandleFunc("/login", user.Container(db).Login())
	ur.HandleFunc("/logout", user.Container(db).Logout())
	ur.Use(CORS)

	// Overlay functionalites router
	or := r.PathPrefix("/overlay").Subrouter()
	or.HandleFunc("/list", listing.Container(db).ListAllDirectories())
	or.HandleFunc("/stream", http.HandlerFunc(static.StreamVideoFile))
	or.HandleFunc("/gallery", http.HandlerFunc(static.ListDirectoryContentHandler))
	or.Use(CORS)
	or.Use(authenticated)

	// Static resources related router
	sr := r.PathPrefix("/static").Subrouter()
	sr.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(cfg.WorkingDir))))
	sr.Use(neuter)

	// Thumbnails related router
	tr := r.PathPrefix("/thumbs").Subrouter()
	tr.PathPrefix("/").Handler(http.StripPrefix("/thumbs", http.FileServer(http.Dir(cfg.CacheDir))))
	tr.Use(neuter)
	tr.Use(serveThumbnail)

	// Frontend
	build, _ := fs.Sub(*app, "frontend/dist")

	sh := SpaHandler{
		entrypoint: "index.html",
		filesystem: &build,
	}

	sh.AddRoute("/login")
	sh.AddRoute("/gallery")
	sh.AddRoute("/help")

	r.PathPrefix("/").Handler(sh.Handler())
	r.Use(CORS)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
