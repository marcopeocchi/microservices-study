package user

import (
	"fuu/v/pkg/domain"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	handler    *Handler
	repository *Repository
	service    *Service

	handlerOnce    sync.Once
	repositoryOnce sync.Once
	serviceOnce    sync.Once
)

func provideHandler(service domain.UserService, logger *zap.SugaredLogger) *Handler {
	handlerOnce.Do(func() {
		handler = &Handler{
			service: service,
			logger:  logger,
		}
	})
	return handler
}

func provideRepository(db *gorm.DB, logger *zap.SugaredLogger) *Repository {
	repositoryOnce.Do(func() {
		repository = &Repository{
			db:     db,
			logger: logger,
		}
	})
	return repository
}

func provideService(repository domain.UserRepository, logger *zap.SugaredLogger) *Service {
	serviceOnce.Do(func() {
		service = &Service{
			repo:   repository,
			logger: logger,
		}
	})
	return service
}
