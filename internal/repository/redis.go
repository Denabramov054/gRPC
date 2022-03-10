package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string // "" for no password
	DB       int    // 0 for default
}

var ctx = context.Background()

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}
