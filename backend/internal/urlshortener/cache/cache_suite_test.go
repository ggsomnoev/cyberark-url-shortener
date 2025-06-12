package redis_test

import (
	"context"
	"testing"
	"time"

	redis "github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/cache"
	testclient "github.com/ggsomnoev/cyberark-url-shortener/test/redis"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx         context.Context
	redisClient *redis.RedisCache
	expiration  time.Duration
)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache Suite")
}

var _ = BeforeEach(func() {
	ctx = context.Background()
	expiration = 1 * time.Second
	redisClient = testclient.NewRedisTestClient(ctx, expiration)
})

var _ = AfterEach(func() {
	Expect(redisClient.Client.Close()).To(Succeed())
})
