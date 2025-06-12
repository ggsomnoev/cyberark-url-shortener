package validator

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
)

var (
	shortCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{8,8}$`)

	ErrInvalidURL       = errors.New("invalid url format")
	ErrInvalidShortCode = errors.New("invalid short code format")
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateShortenRequest(req model.ShortenRequest) error {
	if err := v.ValidateURL(req.URL); err != nil {
		return err
	}

	// Optional param.
	if req.ShortCode == "" {
		return nil
	}

	if err := v.ValidateShortCode(req.ShortCode); err != nil {
		return err
	}

	return nil
}

func (v *Validator) ValidateURL(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

func (v *Validator) ValidateShortCode(shortCode string) error {
	if !shortCodeRegex.MatchString(shortCode) {
		return ErrInvalidShortCode
	}

	return nil
}
