package gallery

import (
	"context"
	"fuu/v/pkg/domain"
	"fuu/v/pkg/utils"
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

	files = slices.Filter(files, func(file fs.DirEntry) bool {
		mimeType := mime.TypeByExtension(filepath.Ext(file.Name()))
		return utils.ValidType.MatchString(mimeType) && utils.ValidFile(file.Name())
	})

	res := make([]string, len(files))

	for i, file := range files {
		if !file.IsDir() {
			res[i] = file.Name()
		}
	}

	sort.SliceStable(res, func(i, j int) bool {
		idx1, err := utils.GetImageIndex(res[i])
		if err != nil {
			return false
		}
		idx2, err := utils.GetImageIndex(res[j])
		if err != nil {
			return false
		}
		return idx1 < idx2
	})

	content := domain.Content{
		List: res,
	}

	encoded, err := json.Marshal(content)

	if err != nil {
		return domain.Content{}, err
	}

	// Write-through caching
	r.rdb.SetNX(ctx, path, encoded, time.Second*30)

	return content, nil
}
