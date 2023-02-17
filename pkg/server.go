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
		server := createAppServer(cfg.Port, frontend, db)
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

func createAppServer(port int, app *embed.FS, db *gorm.DB) *http.Server {
	reactBuild, _ := fs.Sub(*app, "frontend/dist")

	mux := http.NewServeMux()
	mux.Handle("/", reactHandler(&reactBuild))

	mux.Handle("/static/", http.StripPrefix("/static", neuter(authenticated(http.FileServer(http.Dir(cfg.WorkingDir))))))
	mux.Handle("/thumbs/", http.StripPrefix("/thumbs", neuter(authenticated(serveThumbnail(http.FileServer(http.Dir(cfg.CacheDir)))))))

	mux.Handle("/list", CORS(authenticated(listing.Wire(db).ListAllDirectories())))
	mux.Handle("/stream", CORS(authenticated(http.HandlerFunc(static.StreamVideoFile))))
	mux.Handle("/gallery", CORS(authenticated(http.HandlerFunc(static.ListDirectoryContentHandler))))

	mux.Handle("/user/login", CORS(user.Wire(db).Login()))
	mux.Handle("/user/logout", CORS(user.Wire(db).Logout()))

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
