package gallery

import (
	"fuu/v/pkg/domain"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	handler        *Handler
	repository     *Repository
	handlerOnce    sync.Once
	repositoryOnce sync.Once
)

func ProvideHandler(repository domain.DirectoryRepository) *Handler {
	handlerOnce.Do(func() {
		handler = &Handler{
			repo: repository,
		}
	})
	return handler
}

func ProvideRepository(rdb *redis.Client, logger *zap.SugaredLogger, root string) *Repository {
	repositoryOnce.Do(func() {
		repository = &Repository{
			rdb:        rdb,
			logger:     logger,
			workingDir: root,
		}
	})
	return repository
}
