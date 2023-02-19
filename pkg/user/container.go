package user

import "gorm.io/gorm"

func Container(db *gorm.DB) *Handler {
	repository := provideRepository(db)
	service := provideService(repository)
	handler := provideHandler(service)

	return handler
}
