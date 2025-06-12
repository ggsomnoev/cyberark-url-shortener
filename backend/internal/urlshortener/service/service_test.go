package service_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/service"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/service/servicefakes"
)

var (
	ErrCacheGetFailed = errors.New("cache failed to retrieve value")
	ErrStoreFailed    = errors.New("store failed")
	ErrCacheSetFailed = errors.New("cache failed to store value")
)

var _ = Describe("Service", func() {
	When("creating", func() {
		It("should create an instance", func() {
			Expect(service.NewService(nil, nil)).NotTo(BeNil())
		})
	})

	Describe("Instance", func() {
		var (
			ctx   context.Context
			store *servicefakes.FakeStore
			cache *servicefakes.FakeCacheClient
			srv   *service.Service
		)

		BeforeEach(func() {
			ctx = context.Background()

			store = &servicefakes.FakeStore{}
			cache = &servicefakes.FakeCacheClient{}

			srv = service.NewService(store, cache)
		})

		Describe("ResolveURL", func() {
			var (
				shortCode, url string
				urlEntity      model.URL

				resolvedURL string
				errAction   error
			)

			BeforeEach(func() {
				shortCode = uuid.New().String()[:8]
				url = "https://www.example.com"

				urlEntity = model.URL{
					ID:        uuid.New(),
					Original:  url,
					ShortCode: shortCode,
				}

				cache.GetReturns(url, nil)
				store.FindByShortCodeReturns(urlEntity, nil)
				cache.SetReturns(nil)
			})

			JustBeforeEach(func() {
				resolvedURL, errAction = srv.ResolveURL(ctx, shortCode)
			})

			It("find the url value in the cache", func() {
				Expect(errAction).NotTo(HaveOccurred())

				Expect(resolvedURL).To(Equal(url))

				Expect(cache.GetCallCount()).To(Equal(1))
				actualCtx, actualShortCode := cache.GetArgsForCall(0)
				Expect(actualCtx).To(Equal(ctx))
				Expect(actualShortCode).To(Equal(shortCode))

				Expect(store.FindByShortCodeCallCount()).To(Equal(0))
			})

			When("the value is not found in the cache", func() {
				BeforeEach(func() {
					cache.GetReturns("", nil)
				})

				It("gets the value from the DB", func() {
					Expect(errAction).NotTo(HaveOccurred())
					Expect(resolvedURL).To(Equal(url))

					Expect(cache.GetCallCount()).To(Equal(1))
					Expect(store.FindByShortCodeCallCount()).To(Equal(1))

					actualCtx, actualShortCode := store.FindByShortCodeArgsForCall(0)
					Expect(actualCtx).To(Equal(ctx))
					Expect(actualShortCode).To(Equal(shortCode))
				})

				It("adds the resolved URL into the cache", func() {
					Expect(cache.SetCallCount()).To(Equal(1))

					actualCtx, actualShortCode, actualURL := cache.SetArgsForCall(0)
					Expect(actualCtx).To(Equal(ctx))
					Expect(actualShortCode).To(Equal(shortCode))
					Expect(actualURL).To(Equal(url))
				})

				When("cache fails to store the value", func() {
					BeforeEach(func() {
						cache.SetReturns(ErrCacheSetFailed)
					})

					It("fails with a cache error", func() {
						Expect(errAction).To(MatchError(ErrCacheSetFailed))
					})
				})
			})

			When("the cache fails", func() {
				BeforeEach(func() {
					cache.GetReturns("", ErrCacheGetFailed)
				})

				It("fails with a cache error", func() {
					Expect(errAction).To(MatchError(ErrCacheGetFailed))
				})
			})

			When("the url is not found in the cache and the store fails", func() {
				BeforeEach(func() {
					cache.GetReturns("", nil)
					store.FindByShortCodeReturns(model.URL{}, ErrStoreFailed)
				})

				It("fails with a cache error", func() {
					Expect(errAction).To(MatchError(ErrStoreFailed))
				})
			})
		})
		
		Describe("ShortenURL", func() {
			var (
				url string

				resolvedShortCode string
				errAction         error
			)

			BeforeEach(func() {
				url = "https://www.example.com"

				cache.SetReturns(nil)
				store.SaveReturns(nil)
			})

			JustBeforeEach(func() {
				resolvedShortCode, errAction = srv.ShortenURL(ctx, url)
			})

			It("succeeds", func() {
				Expect(errAction).NotTo(HaveOccurred())
				Expect(resolvedShortCode).To(HaveLen(8))

				Expect(cache.SetCallCount()).To(Equal(1))
				actualCtx, actualShortCode, actualURL := cache.SetArgsForCall(0)
				Expect(actualCtx).To(Equal(ctx))
				Expect(actualShortCode).To(HaveLen(8))
				Expect(actualURL).To(Equal(url))

				Expect(store.SaveCallCount()).To(Equal(1))
				actualCtx, actualURLEntity := store.SaveArgsForCall(0)
				Expect(actualCtx).To(Equal(ctx))
				Expect(actualURLEntity.ID).NotTo(BeNil())
				Expect(actualURLEntity.ShortCode).To(HaveLen(8))
				Expect(actualURLEntity.Original).To(Equal(url))
				Expect(actualURLEntity.CreatedAt).NotTo(BeNil())
			})

			When("storing value to the cache fails", func() {
				BeforeEach(func() {
					cache.SetReturns(ErrCacheSetFailed)
				})

				It("returns a cache store error", func() {
					Expect(errAction).To(MatchError(ErrCacheSetFailed))
					Expect(store.FindByShortCodeCallCount()).To(Equal(0))
				})
			})

			When("storing value to the db fails", func() {
				BeforeEach(func() {
					store.SaveReturns(ErrStoreFailed)
				})

				It("fails with a cache error", func() {
					Expect(errAction).To(MatchError(ErrStoreFailed))
				})
			})
		})
	})
})
