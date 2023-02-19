package listing

import "gorm.io/gorm"

func Container(db *gorm.DB) *Handler {
	repository := ProvideRepository(db)
	service := ProvideService(repository)
	handler := ProvideHandler(service)

	return handler
}
