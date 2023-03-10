package listing

import (
	"fuu/v/internal/domain"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	handler        *Handler
	repository     *Repository
	service        *Service
	handlerSingle  sync.Once
	repositoryOnce sync.Once
	serviceOnce    sync.Once
)

func ProvideHandler(service domain.ListingService) *Handler {
	handlerSingle.Do(func() {
		handler = &Handler{
			service: service,
		}
	})

	return handler
}

func ProvideService(repository domain.ListingRepository) *Service {
	serviceOnce.Do(func() {
		service = &Service{
			repo: repository,
		}
	})

	return service
}

func ProvideRepository(db *gorm.DB, rdb *redis.Client, conn *grpc.ClientConn, logger *zap.SugaredLogger) *Repository {
	repositoryOnce.Do(func() {
		repository = &Repository{
			db:     db,
			rdb:    rdb,
			logger: logger,
			conn:   conn,
		}
	})

	return repository
}
