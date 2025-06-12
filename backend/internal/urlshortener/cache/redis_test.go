package redis_test

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redis/go-redis/v9"
)

var _ = Describe("RedisCache", func() {
	var (
		key       string
		value     string
		errAction error
	)

	BeforeEach(func() {
		key = uuid.New().String()[:8]
		value = "test-value"
	})

	Describe("Set", func() {
		JustBeforeEach(func() {
			errAction = redisClient.Set(ctx, key, value)
		})

		JustAfterEach(func() {
			Expect(redisClient.Delete(ctx, key)).To(Succeed())
		})

		It("succeeds", func() {
			Expect(errAction).NotTo(HaveOccurred())
		})

		It("stores the values properly", func() {
			retrieved, err := redisClient.Get(ctx, key)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved).To(Equal(value))
		})
	})

	Describe("Get", func() {
		var (
			key         string
			value       string
			resultValue string
			errAction   error
		)

		BeforeEach(func() {
			key = uuid.New().String()[:8]
			value = "test-value"
		})

		JustBeforeEach(func() {
			resultValue, errAction = redisClient.Get(ctx, key)
		})

		It("fails with redis nil error", func() {
			Expect(errAction).To(MatchError(redis.Nil))
		})

		Context("a key is added", func() {
			BeforeEach(func() {
				Expect(redisClient.Set(ctx, key, value)).To(Succeed())
			})

			AfterEach(func() {
				Expect(redisClient.Delete(ctx, key)).To(Succeed())
			})

			It("succeeds", func() {
				Expect(errAction).NotTo(HaveOccurred())
				Expect(resultValue).To(Equal(value))
			})

			Context("and the key expires", func() {
				BeforeEach(func() {
					time.Sleep(expiration + 500*time.Microsecond)
				})

				It("fails with redis nil error", func() {
					Expect(errAction).To(MatchError(redis.Nil))
				})
			})
		})
	})
})
