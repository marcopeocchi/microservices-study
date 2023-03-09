package workers

import (
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type FileWatcher struct {
	watcher       *fsnotify.Watcher
	WorkingDir    string
	OnFileCreated func(event string)
	OnFileDeleted func(event string)
	Logger        *zap.SugaredLogger
}

func (f *FileWatcher) Setup() {
	var err error

	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		f.Logger.Fatalln(err)
	}
	err = f.watcher.Add(f.WorkingDir)
	if err != nil {
		f.Logger.Fatalln(err)
	}
}

func (f *FileWatcher) Start() {
	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				f.Logger.Infow("added directory", "event", event.Name)
				f.OnFileCreated(event.Name)
			}
			if event.Has(fsnotify.Remove) {
				f.Logger.Infow("removing directory", "event", event.Name)
				f.OnFileDeleted(event.Name)
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}
			f.Logger.Errorln(err)
		}
	}
}

func (f *FileWatcher) Close() {
	defer f.watcher.Close()
}
