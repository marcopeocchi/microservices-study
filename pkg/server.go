package pkg

import (
	"context"
	"embed"
	"fmt"
	"fuu/v/pkg/cli"
	"fuu/v/pkg/domain"
	"fuu/v/pkg/listing"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gorm.io/gorm"
)

var (
	config Config
	db     *gorm.DB
)

func createAppServer(ctx context.Context, port int, app *embed.FS) *http.Server {
	reactBuild, _ := fs.Sub(*app, "frontend/dist")

	mux := http.NewServeMux()
	mux.Handle("/", reactHandler(&reactBuild))

	mux.Handle(
		"/static/",
		http.StripPrefix("/static",
			neuter(authenticated(http.FileServer(http.Dir(config.WorkingDir)))),
		),
	)

	mux.Handle(
		"/thumbnails/",
		http.StripPrefix("/thumbnails",
			neuter(authenticated(serveThumbnail(http.FileServer(http.Dir(config.CacheDir))))),
		),
	)
	mux.Handle("/user", CORS(http.HandlerFunc(loginHandler)))

	mux.Handle("/list", CORS(authenticated(listing.Wire(db).ListAllDirectories())))
	mux.Handle("/stream", CORS(authenticated(http.HandlerFunc(streamVideoFile))))
	mux.Handle("/gallery", CORS(authenticated(http.HandlerFunc(listDirectoryContentHandler))))

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}

func RunBlocking(cfg Config, localdb *gorm.DB, frontend *embed.FS) {
	config = cfg
	db = localdb
	db.AutoMigrate(&domain.Directory{})

	thumbnailer := Thumbnailer{
		BaseDir:           config.WorkingDir,
		ImgHeight:         config.ThumbnailHeight,
		ImgQuality:        config.ThumbnailQuality,
		ForceRegeneration: config.ForceRegeneration,
		CacheDir:          config.CacheDir,
		Database:          db,
	}

	log.Println("Starting server")

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// Discourage the execution of this program as SuperUser.
	// Unless in executed docker because of obvious reasons.
	uid := os.Getuid()
	if uid == 0 {
		log.Println(cli.Yellow, "You're running this program as root (UID 0)", cli.Reset)
		log.Println(cli.Yellow, "This isn't reccomended unless you're using Docker", cli.Reset)
	}

	type ctxKey string
	serverContext := context.Background()
	serverContext = context.WithValue(serverContext, ctxKey("config"), config)

	go func() {
		log.Println("Starting thumbnailer")
		start := time.Now()
		thumbnailer.Start()
		log.Println("Thumbnailer took", cli.Format(time.Since(start), cli.BgBlue))
		wg.Done()
	}()
	go func() {
		server := createAppServer(serverContext, config.Port, frontend)
		server.ListenAndServe()
		wg.Done()
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// Start a light NON-Recursive Filesystem watcher as a background routine.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) {
					log.Println("Added directory:", event.Name)
					thumbnailer.Start()
				}
				if event.Has(fsnotify.Remove) {
					log.Println("Removing directory:", event.Name)
					thumbnailer.Remove(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(config.WorkingDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server started")

	// wait for the waitgroup to finish, which it will not.
	// effectively blocks.
	wg.Wait()
}
