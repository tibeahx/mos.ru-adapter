package redis

import (
	"github.com/tibeahx/mos.ru-adapter/internal/config"

	"github.com/redis/go-redis/v9"
)

type RC struct {
	Redis *redis.Client
}

func NewRC(cfg *config.Config) *RC {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisClientAddr,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	return &RC{Redis: client}
}
