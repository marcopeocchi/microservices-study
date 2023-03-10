package listing

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func Container(db *gorm.DB, rdb *redis.Client, conn *grpc.ClientConn, logger *zap.SugaredLogger) *Handler {
	repository := ProvideRepository(db, rdb, conn, logger)
	service := ProvideService(repository)
	handler := ProvideHandler(service)

	return handler
}
