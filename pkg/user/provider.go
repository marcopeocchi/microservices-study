package user

import (
	"fuu/v/pkg/domain"
	"sync"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	handler    *Handler
	repository *Repository
	service    *Service

	handlerOnce    sync.Once
	repositoryOnce sync.Once
	serviceOnce    sync.Once

	ProviderSet wire.ProviderSet = wire.NewSet(
		provideHandler,
		provideRepository,
		provideService,
		wire.Bind(new(domain.UserHandler), new(*Handler)),
		wire.Bind(new(domain.UserRepository), new(*Repository)),
		wire.Bind(new(domain.UserService), new(*Service)),
	)
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
