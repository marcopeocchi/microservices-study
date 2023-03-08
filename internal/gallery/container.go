package gallery

import (
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func Container(rdb *redis.Client, logger *zap.SugaredLogger, ch *amqp.Channel, root string) *Handler {
	repository := ProvideRepository(rdb, logger, ch, root)
	handler := ProvideHandler(repository)

	return handler
}
