package infra

import (
	"context"
	"message_app/internal/logger"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Cr *redis.Client
}

func NewRedisClient(logger *logger.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.L.Fatal(err.Error())
	}

	return client
}
