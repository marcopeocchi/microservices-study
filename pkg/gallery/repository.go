package gallery

import (
	"context"
	"fmt"
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

func (r *Repository) FindByPath(ctx context.Context, path string) (domain.Content, error) {
	cached, _ := r.rdb.Get(ctx, path).Bytes()

	if len(cached) > 0 {
		res := domain.Content{}
		err := json.Unmarshal(cached, &res)

		return res, err
	}

	files, _ := os.ReadDir(filepath.Join(r.workingDir, path))
	filesAvif, _ := os.ReadDir(filepath.Join(r.workingDir, path, "avif"))

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

	for i, file := range files {
		if !file.IsDir() {
			resOrig[i] = file.Name()
		}
	}

	// Lazy convert all pictures
	go workers.Avifier(filepath.Join(r.workingDir, path), resOrig)

	for i, file := range filesAvif {
		if !file.IsDir() {
			if filepath.Ext(file.Name()) == ".avif" {
				resAvif[i] = fmt.Sprintf("/avif/%s", file.Name())
				continue
			}
		}
	}

	sortFunc := func(i, j int, v []string) bool {
		idx1, err := utils.GetImageIndex(v[i])
		if err != nil {
			return false
		}
		idx2, err := utils.GetImageIndex(v[j])
		if err != nil {
			return false
		}
		return idx1 < idx2
	}

	sort.SliceStable(resOrig, func(i, j int) bool {
		return sortFunc(i, j, resOrig)
	})

	sort.SliceStable(resAvif, func(i, j int) bool {
		return sortFunc(i, j, resAvif)
	})

	content := domain.Content{
		List:          resOrig,
		Avif:          resAvif,
		AvifAvailable: len(resOrig) == len(resAvif),
	}

	encoded, err := json.Marshal(content)

	if err != nil {
		return domain.Content{}, err
	}

	// Write-through caching
	r.rdb.SetNX(ctx, path, encoded, time.Second*30)

	return content, nil
}
