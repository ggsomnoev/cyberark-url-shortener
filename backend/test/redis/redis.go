package redis

import (
	"context"
	"time"

	redis "github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/cache"
)

const (
	testRedisAddr     = "localhost:6379"
	testRedisPassword = ""
	testRedisDB       = 0
)

func NewRedisTestClient(ctx context.Context, keyExpiration time.Duration) *redis.RedisCache {
	return redis.NewRedisCache(testRedisAddr, testRedisPassword, testRedisDB, keyExpiration)
}
