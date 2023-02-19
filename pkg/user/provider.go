package user

import (
	"fuu/v/pkg/domain"
	"sync"

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

func provideHandler(service domain.UserService) *Handler {
	handlerOnce.Do(func() {
		handler = &Handler{
			service: service,
		}
	})
	return handler
}

func provideRepository(db *gorm.DB) *Repository {
	repositoryOnce.Do(func() {
		repository = &Repository{
			db: db,
		}
	})
	return repository
}

func provideService(repository domain.UserRepository) *Service {
	serviceOnce.Do(func() {
		service = &Service{
			repo: repository,
		}
	})
	return service
}
