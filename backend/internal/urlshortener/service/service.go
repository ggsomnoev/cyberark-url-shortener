package service

import (
	"context"
	"fmt"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//counterfeiter:generate . Store
type Store interface {
	Save(ctx context.Context, urlEntity model.URL) error
	FindByShortCode(ctx context.Context, code string) (model.URL, error)
}

//counterfeiter:generate . CacheClient
type CacheClient interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
}

type Service struct {
	store Store
	cache CacheClient
}

func NewService(store Store, cache CacheClient) *Service {
	return &Service{
		store: store,
		cache: cache,
	}
}

func (s *Service) ResolveURL(ctx context.Context, shortCode string) (string, error) {
	originalURL, err := s.cache.Get(ctx, shortCode)
	if err != nil && err != redis.Nil {
		return "", fmt.Errorf("failed to get original url from the cache: %w", err)
	}

	if originalURL != "" {
		return originalURL, nil
	}

	url, err := s.store.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("failed to get original url: %w", err)
	}

	if err = s.cache.Set(ctx, url.ShortCode, url.Original); err != nil {
		return "", fmt.Errorf("failed to add to the cache - short_code: %s, original_url: %s: %w", url.ShortCode, url.Original, err)
	}

	return url.Original, nil
}

func (s *Service) ShortenURL(ctx context.Context, orginalURL string) (string, error) {
	shortCode := generateShortCode()

	urlEntity := model.URL{
		ID:        uuid.New(),
		Original:  orginalURL,
		ShortCode: shortCode,
	}

	if err := s.cache.Set(ctx, shortCode, orginalURL); err != nil {
		return "", fmt.Errorf("failed to add to the cache - short_code: %s, original_url: %s: %w", shortCode, orginalURL, err)
	}

	if err := s.store.Save(ctx, urlEntity); err != nil {
		return "", fmt.Errorf("failed to store url entity: %w", err)
	}

	return urlEntity.ShortCode, nil
}

func generateShortCode() string {
	return uuid.New().String()[:8]
}
