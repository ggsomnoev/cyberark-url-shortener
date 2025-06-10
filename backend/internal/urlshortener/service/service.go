package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
	"github.com/google/uuid"
)

type Store interface {
	Save(ctx context.Context, urlEntity model.URL) error
	FindByShortCode(ctx context.Context, code string) (model.URL, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) ResolveURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.store.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("failed to get original url:%w", err)
	}

	return url.Original, nil
}

func (s *Service) ShortenURL(ctx context.Context, orginalURL string) (string, error) {
	shortCode := generateShortCode()

	urlEntity := model.URL{
		ID:        uuid.New().String(),
		Original:  orginalURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
	}
	if err := s.store.Save(ctx, urlEntity); err != nil {
		return "", fmt.Errorf("failed to store url entity: %w", err)
	}

	return urlEntity.ShortCode, nil
}

func generateShortCode() string {
	return uuid.New().String()[:8]
}
