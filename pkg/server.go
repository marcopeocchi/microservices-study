package pkg

import (
	"context"
	"embed"
	"fmt"
	"fuu/v/pkg/cli"
	"fuu/v/pkg/models"
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
	appContext = context.Background()
	config     Config
	db         *gorm.DB
)

func createAppServer(name string, port int, app *embed.FS) *http.Server {
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
	mux.Handle("/list", CORS(authenticated(http.HandlerFunc(listDirectoryHandler))))
	mux.Handle("/stream", CORS(authenticated(http.HandlerFunc(streamVideoFile))))
	mux.Handle("/gallery", CORS(authenticated(http.HandlerFunc(listDirectoryContentHandler))))

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}

	return &server
}

func RunBlocking(ctx context.Context) {
	appContext = ctx
	config = ctx.Value("config").(Config)
	db = ctx.Value("db").(*gorm.DB)

	db.AutoMigrate(&models.Directory{})

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

	go func() {
		log.Println("Starting thumbnailer")
		start := time.Now()
		thumbnailer.Start()
		log.Println("Thumbnailer took", cli.Format(time.Since(start), cli.BgBlue))
		wg.Done()
	}()
	go func() {
		server := createAppServer("app", config.Port, ctx.Value("react").(*embed.FS))
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
