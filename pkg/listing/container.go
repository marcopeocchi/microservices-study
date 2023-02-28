package listing

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Container(db *gorm.DB, rdb *redis.Client) *Handler {
	repository := ProvideRepository(db, rdb)
	service := ProvideService(repository)
	handler := ProvideHandler(service)

	return handler
}
