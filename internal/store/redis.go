package storage

import (
	"test-task/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Redis    *redis.Client
	addr     string
	username string
	password string
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisClientAddr,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	return &RedisClient{
		Redis:    client,
		addr:     client.Options().Addr,
		username: client.Options().Username,
		password: client.Options().Password,
	}
}
