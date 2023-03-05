package workers

import (
	"fmt"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	config "fuu/v/pkg/config"
	"fuu/v/pkg/domain"
	"fuu/v/pkg/utils"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/google/uuid"
	"github.com/marcopeocchi/fazzoletti/slices"
	"gorm.io/gorm"
)

type Thumbnailer struct {
	BaseDir           string
	ImgHeight         int
	ImgQuality        int
	ForceRegeneration bool
	CacheDir          string
	Database          *gorm.DB
}

type job struct {
	Id             string
	InputFile      string
	OutputFile     string
	WorkingDirPath string
	WorkingDirName string
	IsImage        bool
}

func (t *Thumbnailer) Start() {
	t.prune()

	files, err := os.ReadDir(t.BaseDir)
	log.Printf("Creating thumbnails for %d entries incrementally\n", len(files))

	if err != nil {
		log.Fatal(err)
	}

	workQueue := make([]job, len(files))

	for i, file := range files {
		if file.IsDir() {
			current := file.Name()
			workingDir := filepath.Join(t.BaseDir, current)
			content, err := os.ReadDir(filepath.Join(t.BaseDir, current))
			if err != nil {
				log.Fatal(err)
			}

			if !t.ForceRegeneration {
				var row domain.Directory
				t.Database.First(&row, "path = ?", workingDir)

				if row.Thumbnail != "" {
					continue
				}
			}

			for _, f := range content {
				mimeType := mime.TypeByExtension(filepath.Ext(f.Name()))
				if utils.ValidType.MatchString(mimeType) && utils.ValidFile(f.Name()) {

					uuid, err := uuid.NewRandom()
					if err != nil {
						log.Fatalln("cannot generate id for thumbnail")
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
	target := new(domain.Directory)
	t.Database.Where("path = ?", fmt.Sprintf("`%s`", dirpath)).First(&target)

	os.Remove(filepath.Join(t.CacheDir, target.Thumbnail))
	t.Database.Delete(&target, target.ID)
}

func (t *Thumbnailer) mainThread(queue []job) {
	// generate n thumbnails at time where n is core number
	maxConcurrency := runtime.NumCPU()
	log.Printf("Starting thumbnailer on %d cores\n", maxConcurrency)
	log.Println(len(queue), "directories needs a thumbnail")

	// block if guard channel is filled with n jobs
	pipeline := make(chan int, maxConcurrency)
	messages := make(chan job)

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
					"-format", config.Instance().ImageOptimizationFormat,
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
					"-f", config.Instance().ImageOptimizationFormat,
					w.OutputFile,
				)
			}
			err := cmd.Start()
			if err != nil {
				log.Panicln(err)
			}
			// join
			err = cmd.Wait()
			if err == nil {
				log.Println("Generated thumbnail for", w.InputFile)
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
		t.Database.Create(&domain.Directory{
			Name:      w.WorkingDirName,
			Path:      w.WorkingDirPath,
			Thumbnail: w.Id,
		})
	}
}

func (t *Thumbnailer) prune() {
	all := new([]domain.Directory)
	t.Database.Find(all)

	filter := bloom.NewWithEstimates(uint(len(*all)), 0.01)

	log.Println("Start database prune")
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
			log.Println("Deleting", toRemove)
			os.Remove(toRemove)
		}
	}

	log.Println("Database pruned removed", count, "rows")
}
