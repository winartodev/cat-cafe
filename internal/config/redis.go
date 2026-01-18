package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (r *RedisConfig) SetupRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
