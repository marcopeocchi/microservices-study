package listing

import (
	"fuu/v/pkg/domain"
	"sync"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	handler        *Handler
	repository     *Repository
	service        *Service
	handlerSingle  sync.Once
	repositoryOnce sync.Once
	serviceOnce    sync.Once
	ProviderSet    wire.ProviderSet = wire.NewSet(
		ProvideHandler,
		ProvideRepository,
		ProvideService,
		wire.Bind(new(domain.ListingHandler), new(*Handler)),
		wire.Bind(new(domain.ListingRepository), new(*Repository)),
		wire.Bind(new(domain.ListingService), new(*Service)),
	)
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

func ProvideRepository(db *gorm.DB) *Repository {
	repositoryOnce.Do(func() {
		repository = &Repository{
			db: db,
		}
	})

	return repository
}
