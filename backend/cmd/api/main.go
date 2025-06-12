package main

import (
	"context"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/config"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/logger"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/pg"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/webapi"
)

func main() {
	appCtx := context.Background()
	srv := webapi.NewWebAPI()

	cfg := config.Load()

	dbConfig := pg.PoolConfig{
		MinConns:        cfg.DBMinConns,
		MaxConns:        cfg.DBMaxConns,
		MaxConnLifetime: cfg.DBMaxConnLifeTime,
		MaxConnIdleTime: cfg.DBMaxConnIdleTime,
	}

	pool, err := pg.InitDBConnection(appCtx, cfg.DBConnectionURL, dbConfig)
	if err != nil {
		panic(err)
	}

	urlshortener.Process(appCtx, pool, srv)

	logger.GetLogger().Fatal(srv.Start(":" + cfg.APIPort))
}
