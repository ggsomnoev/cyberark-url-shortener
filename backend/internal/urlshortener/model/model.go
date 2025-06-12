package model

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID        uuid.UUID
	Original  string
	ShortCode string
	CreatedAt time.Time
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
