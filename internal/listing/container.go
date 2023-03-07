package listing

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Container(db *gorm.DB, rdb *redis.Client, logger *zap.SugaredLogger) *Handler {
	repository := ProvideRepository(db, rdb, logger)
	service := ProvideService(repository)
	handler := ProvideHandler(service)

	return handler
}
