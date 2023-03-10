package workers

import (
	"context"
	"mime"
	"os"
	"path/filepath"
	"runtime"
	"time"

	thumbnailspb "fuu/v/gen/go/grpc/thumbnails/v1"
	"fuu/v/internal/domain"
	"fuu/v/pkg/config"
	utils "fuu/v/pkg/utils"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/google/uuid"
	"github.com/marcopeocchi/fazzoletti/slices"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type Thumbnailer struct {
	BaseDir    string
	ImgHeight  int
	ImgQuality int
	CacheDir   string
	Database   *gorm.DB
	Logger     *zap.SugaredLogger
	conn       *grpc.ClientConn
	client     thumbnailspb.ThumbnailServiceClient
}

type job struct {
	Id             string
	InputFile      string
	OutputFile     string
	WorkingDirPath string
	WorkingDirName string
	IsImage        bool
}

func getGrpcClient(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}

	return grpc.DialContext(ctx, addr, opts...)
}

func (t *Thumbnailer) Generate() {
	conn, err := getGrpcClient("localhost:10099")
	t.client = thumbnailspb.NewThumbnailServiceClient(conn)

	if err != nil {
		panic(err)
	}
	t.conn = conn

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
	t.Logger.Infow(
		"starting thumbnailer",
		"cores", runtime.NumCPU(),
		"directories", len(queue),
	)

	format := config.Instance().ImageOptimizationFormat
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, work := range queue {
		t.client.Generate(ctx, &thumbnailspb.GenerateRequest{
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

func (t *Thumbnailer) prune() {
	all := &[]domain.Directory{}
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
