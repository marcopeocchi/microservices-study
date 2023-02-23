package gallery

import "github.com/redis/go-redis/v9"

func Container(rdb *redis.Client, root string) *Handler {
	repository := ProvideRepository(rdb, root)
	handler := ProvideHandler(repository)

	return handler
}
