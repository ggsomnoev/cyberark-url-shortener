package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/logger"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
	"github.com/labstack/echo/v4"
)

//counterfeiter:generate . Service
type Service interface {
	ResolveURL(context.Context, string) (string, error)
	ShortenURL(context.Context, string) (string, error)
}

//counterfeiter:generate . Validator
type Validator interface {
	ValidateShortenRequest(model.ShortenRequest) error
	ValidateShortCode(string) error
}

func RegisterRoutes(ctx context.Context, srv *echo.Echo, urlShortener Service, validator Validator) {
	if srv != nil {
		srv.POST("/api/shorten", handleURLShorten(ctx, urlShortener, validator))
		srv.GET("/:code", handleRedirects(ctx, urlShortener, validator))
	} else {
		logger.GetLogger().Warn("Running routes without a webapi server. Did not register routes")
	}
}

func handleURLShorten(ctx context.Context, urlShortener Service, validator Validator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var urlEntity model.ShortenRequest
		if err := c.Bind(&urlEntity); err != nil {
			c.JSON(http.StatusBadRequest, echo.Map{
				"message": "failed to resolve URL parameter",
				"error":   err.Error(),
			})
		}

		if err := validator.ValidateShortenRequest(urlEntity); err != nil {
			c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid request data",
				"error":   err.Error(),
			})
		}

		shortCode, err := urlShortener.ShortenURL(ctx, urlEntity.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "could not shorten URL",
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.ShortenResponse{
			ShortURL: fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, shortCode),
		})
	}
}

func handleRedirects(ctx context.Context, urlShortener Service, validator Validator) echo.HandlerFunc {
	return func(c echo.Context) error {
		shortCode := c.Param("code")

		if err := validator.ValidateShortCode(shortCode); err != nil {
			c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid request data",
				"error":   err.Error(),
			})
		}

		originalURL, err := urlShortener.ResolveURL(ctx, shortCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "could not find requested URL",
				"error":   err.Error(),
			})
		}

		return c.Redirect(http.StatusFound, originalURL)
	}
}
