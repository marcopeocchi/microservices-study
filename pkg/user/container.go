package user

import "gorm.io/gorm"

// Provide a dependency injection container using the singleton pattern
func Container(db *gorm.DB) *Handler {
	repository := provideRepository(db)
	service := provideService(repository)
	handler := provideHandler(service)

	return handler
}
