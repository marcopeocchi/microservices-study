package gallery

import (
	"context"
	"fmt"
	config "fuu/v/pkg/config"
	"fuu/v/pkg/domain"
	"fuu/v/pkg/utils"
	"fuu/v/pkg/workers"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/goccy/go-json"
	"github.com/marcopeocchi/fazzoletti/slices"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	rdb        *redis.Client
	workingDir string
}

var (
	imageFormat = config.Instance().ImageOptimizationFormat
)

func (r *Repository) FindByPath(ctx context.Context, path string) (domain.Content, error) {
	cached, _ := r.rdb.Get(ctx, path).Bytes()

	if len(cached) > 0 {
		res := domain.Content{}
		err := json.Unmarshal(cached, &res)

		return res, err
	}

	wd := filepath.Join(r.workingDir, path)

	files, _ := os.ReadDir(wd)
	filesAvif, _ := os.ReadDir(filepath.Join(wd, "avif"))
	filesWebp, _ := os.ReadDir(filepath.Join(wd, "webp"))

	filterFunc := func(file fs.DirEntry) bool {
		mimeType := mime.TypeByExtension(filepath.Ext(file.Name()))
		return utils.ValidType.MatchString(mimeType) && utils.ValidFile(file.Name())
	}

	files = slices.Filter(files, func(file fs.DirEntry) bool {
		return filterFunc(file)
	})

	filesAvif = slices.Filter(filesAvif, func(file fs.DirEntry) bool {
		return filterFunc(file)
	})

	resOrig := make([]string, len(files))
	resAvif := make([]string, len(filesAvif))
	resWebp := make([]string, len(filesWebp))

	for i, file := range files {
		if !file.IsDir() {
			resOrig[i] = file.Name()
		}
	}

	// Lazy convert all pictures
	go workers.Converter(wd, resOrig, imageFormat)

	for i, file := range filesAvif {
		if !file.IsDir() {
			resAvif[i] = fmt.Sprintf("/avif/%s", file.Name())
		}
	}

	for i, file := range filesWebp {
		if !file.IsDir() {
			resWebp[i] = fmt.Sprintf("/webp/%s", file.Name())
		}
	}

	sort.SliceStable(resOrig, func(i, j int) bool {
		return utils.FilesSortFunc(i, j, resOrig)
	})

	sort.SliceStable(resAvif, func(i, j int) bool {
		return utils.FilesSortFunc(i, j, resAvif)
	})

	sort.SliceStable(resWebp, func(i, j int) bool {
		return utils.FilesSortFunc(i, j, resWebp)
	})

	content := domain.Content{
		Source:        resOrig,
		Avif:          resAvif,
		WebP:          resWebp,
		AvifAvailable: len(resOrig) == len(resAvif),
		WebPAvailable: len(resOrig) == len(resWebp),
	}

	encoded, err := json.Marshal(content)

	if err != nil {
		return domain.Content{}, err
	}

	// Write-through caching
	r.rdb.SetNX(ctx, path, encoded, time.Second*30)

	return content, nil
}
