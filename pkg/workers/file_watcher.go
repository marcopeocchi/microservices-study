package workers

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher    *fsnotify.Watcher
	WorkingDir string
}

func (f *FileWatcher) New() {
	var err error

	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	err = f.watcher.Add(f.WorkingDir)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FileWatcher) Start(onFileCreated, onFileDeleted func(event string)) {
	// Start a light NON-Recursive Filesystem watcher as a background routine.
	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				log.Println("Added directory:", event.Name)
				onFileCreated(event.Name)
			}
			if event.Has(fsnotify.Remove) {
				log.Println("Removing directory:", event.Name)
				onFileDeleted(event.Name)
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (f *FileWatcher) Close() {
	defer f.watcher.Close()
}
