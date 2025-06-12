package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client     *redis.Client
	Expiration time.Duration
}

func NewRedisCache(address string, password string, db int, expiration time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		Client:     client,
		Expiration: expiration,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, value string) error {
	return r.Client.Set(ctx, key, value, r.Expiration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
