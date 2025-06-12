package urlshortener

import (
	"context"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/handler"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/service"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Process(ctx context.Context, pool *pgxpool.Pool, srv *echo.Echo) {
	urlShortenerStore := store.NewStore(pool)
	urlShortenerService := service.NewService(urlShortenerStore)

	handler.RegisterRoutes(ctx, srv, urlShortenerService)
}
