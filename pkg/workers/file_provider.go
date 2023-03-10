package workers

import (
	"context"
	"mime"
	"os"
	"path/filepath"
	"runtime"

	thumbnailspb "fuu/v/gen/go/grpc/thumbnails/v1"
	"fuu/v/internal/domain"
	"fuu/v/pkg/config"
	"fuu/v/pkg/utils"

	"github.com/google/uuid"
	"github.com/marcopeocchi/fazzoletti/slices"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Thumbnailer struct {
	BaseDir    string
	ImgHeight  int
	ImgQuality int
	CacheDir   string
	Database   *gorm.DB
	Logger     *zap.SugaredLogger
	Conn       *grpc.ClientConn
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
	t.send(slices.Filter(workQueue, func(entry job) bool {
		return entry.InputFile != ""
	}))
	// GC
	workQueue = nil
}

func (t *Thumbnailer) Remove(dirpath string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := thumbnailspb.NewThumbnailServiceClient(t.Conn)
	client.Delete(ctx, &thumbnailspb.DeleteRequest{
		Path: dirpath,
	})
}

func (t *Thumbnailer) send(queue []job) {
	t.Logger.Infow(
		"starting thumbnailer",
		"cores", runtime.NumCPU(),
		"directories", len(queue),
	)

	format := config.Instance().ImageOptimizationFormat
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := thumbnailspb.NewThumbnailServiceClient(t.Conn)

	for _, work := range queue {
		client.Generate(ctx, &thumbnailspb.GenerateRequest{
			Path:   work.InputFile,
			Folder: work.WorkingDirPath,
			Format: format,
		})
		t.Database.FirstOrCreate(&domain.Directory{
			Path: work.WorkingDirPath,
			Name: work.WorkingDirName,
		})
	}
}
