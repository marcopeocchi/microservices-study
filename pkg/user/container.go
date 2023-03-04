package user

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Provide a dependency injection container using the singleton pattern
func Container(db *gorm.DB, logger *zap.SugaredLogger) *Handler {
	repository := provideRepository(db, logger)
	service := provideService(repository, logger)
	handler := provideHandler(service, logger)

	return handler
}
