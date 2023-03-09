package workers

import (
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"fuu/v/internal/domain"
	config "fuu/v/pkg/config"
	utils "fuu/v/pkg/utils"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/google/uuid"
	"github.com/marcopeocchi/fazzoletti/slices"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Thumbnailer struct {
	BaseDir    string
	ImgHeight  int
	ImgQuality int
	CacheDir   string
	Database   *gorm.DB
	Logger     *zap.SugaredLogger
}

type job struct {
	Id             string
	InputFile      string
	OutputFile     string
	WorkingDirPath string
	WorkingDirName string
	IsImage        bool
}

func (t *Thumbnailer) Generate() {
	t.prune()

	files, err := os.ReadDir(t.BaseDir)

	if err != nil {
		t.Logger.Fatalln(err)
	}

	workQueue := make([]job, len(files))

	for i, file := range files {
		if file.IsDir() {
			current := file.Name()
			workingDir := filepath.Join(t.BaseDir, current)
			content, err := os.ReadDir(workingDir)
			if err != nil {
				t.Logger.Fatalln(err)
			}

			test := &domain.Directory{}
			t.Database.Where("path = ?", workingDir).First(&test)

			if test.Thumbnail != "" {
				continue
			}

			for _, f := range content {
				mimeType := mime.TypeByExtension(filepath.Ext(f.Name()))
				if utils.ValidType.MatchString(mimeType) && utils.ValidFile(f.Name()) {

					uuid, err := uuid.NewRandom()
					if err != nil {
						t.Logger.Fatalln(err)
					}

					workQueue[i] = job{
						Id:             uuid.String(),
						WorkingDirName: current,
						WorkingDirPath: workingDir,
						InputFile:      filepath.Join(t.BaseDir, current, f.Name()),
						OutputFile:     filepath.Join(t.CacheDir, uuid.String()),
						IsImage:        utils.IsImage(mimeType),
					}
					break
				}
			}
		}
	}
	t.mainThread(slices.Filter(workQueue, func(entry job) bool {
		return entry.InputFile != ""
	}))
	// GC
	workQueue = nil
}

func (t *Thumbnailer) Remove(dirpath string) {
	target := &domain.Directory{}
	t.Database.Where("path = ?", dirpath).First(&target)

	os.Remove(filepath.Join(t.CacheDir, target.Thumbnail))
	t.Database.Delete(&target, target.ID)
}

func (t *Thumbnailer) mainThread(queue []job) {
	// generate n thumbnails at time where n is core number
	maxConcurrency := runtime.NumCPU()
	t.Logger.Infow(
		"starting thumbnailer",
		"cores", maxConcurrency,
		"directories", len(queue),
	)

	// block if guard channel is filled with n jobs
	pipeline := make(chan int, maxConcurrency)
	messages := make(chan job)

	format := config.Instance().ImageOptimizationFormat

	go t.thumbnailRefSaver(messages)

	for _, work := range queue {
		// take
		pipeline <- 1
		// job closure
		go func(w job) {
			var cmd *exec.Cmd

			if w.IsImage {
				cmd = exec.Command(
					"convert", w.InputFile,
					"-geometry", fmt.Sprintf("x%d", t.ImgHeight),
					"-format", format,
					"-quality", strconv.Itoa(t.ImgQuality),
					w.OutputFile,
				)
			} else {
				cmd = exec.Command(
					"ffmpeg",
					"-i", w.InputFile,
					"-ss", "00:00:01.000",
					"-vframes", "1",
					"-filter:v", fmt.Sprintf("scale=-1:%d", t.ImgHeight),
					"-f", format,
					w.OutputFile,
				)
			}
			err := cmd.Start()
			if err != nil {
				t.Logger.Fatalln(err)
			}
			// join
			err = cmd.Wait()
			if err == nil {
				t.Logger.Infow("generated thumbnail", "file", w.InputFile)
			}
			// Save to db
			messages <- w
			<-pipeline
		}(work)
	}
}

// Execute a db query for-each message received from the channel.
// The operations should be serialized and so doable for sqlite.
// A transaction setup to ensure the lock of the db.
func (t *Thumbnailer) thumbnailRefSaver(messages chan job) {
	for w := range messages {
		if w.Id != "" {
			t.Database.Create(&domain.Directory{
				Name:      w.WorkingDirName,
				Path:      w.WorkingDirPath,
				Thumbnail: w.Id,
			})
		}
	}
}

func (t *Thumbnailer) prune() {
	all := new([]domain.Directory)
	t.Database.Find(all)

	filter := bloom.NewWithEstimates(uint(len(*all)), 0.01)

	t.Logger.Infoln("started database prune")
	count := 0

	for _, entry := range *all {
		_, err := os.Stat(entry.Path)
		if os.IsNotExist(err) {
			t.Database.Where("path = ?", entry.Path).Delete(&domain.Directory{})
			count++
		}
		if err == nil {
			filter.AddString(entry.Thumbnail)
		}
	}

	files, _ := os.ReadDir(t.CacheDir)
	for _, file := range files {
		if !filter.TestString(file.Name()) && filepath.Ext(file.Name()) != ".db" {
			toRemove := filepath.Join(t.CacheDir, file.Name())
			t.Logger.Infow("deleting dead enrty", "file", toRemove)
			os.Remove(toRemove)
		}
	}

	t.Logger.Infow("finished database prune", "count", count)
}
