package model

import "time"

type URL struct {
	ID        string
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
