package main

import (
	"context"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/config"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/logger"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/pg"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener"
	redis "github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/cache"
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

	redisClient := redis.NewRedisCache(cfg.RedisAddress, cfg.RedisPassword, int(cfg.RedisDB), cfg.RedisKeyExpiration)

	urlshortener.Process(appCtx, pool, srv, redisClient)

	logger.GetLogger().Fatal(srv.Start(":" + cfg.APIPort))
}
