package gallery

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Container(rdb *redis.Client, logger *zap.SugaredLogger, root string) *Handler {
	repository := ProvideRepository(rdb, logger, root)
	handler := ProvideHandler(repository)

	return handler
}
