package config

import (
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv"
)

type Config struct {
	APIPort            string
	DBConnectionURL    string
	DBMaxConns         int32
	DBMinConns         int32
	DBMaxConnIdleTime  time.Duration
	DBMaxConnLifeTime  time.Duration
	RedisAddress       string
	RedisPassword      string
	RedisDB            int32
	RedisKeyExpiration time.Duration
}

func Load() *Config {
	return &Config{
		APIPort:            getEnv("API_PORT", "5000"),
		DBConnectionURL:    getEnv("DB_CONNECTION_URL", ""),
		DBMinConns:         getInt32("DB_MIN_CONNS", 1),
		DBMaxConns:         getInt32("DB_MAX_CONNS", 5),
		DBMaxConnIdleTime:  getDuration("DB_MIN_CONN_IDLE_TIME", 1*time.Minute),
		DBMaxConnLifeTime:  getDuration("DB_MIN_CONN_LIFE_TIME", 5*time.Minute),
		RedisAddress:       getEnv("REDIS_ADDR", ""),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getInt32("REDIS_DB", 0),
		RedisKeyExpiration: getDuration("REDIS_EXPIRATION_SECONDS", 60*60*time.Second), // 1 hour
	}
}

func getEnv(key string, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defValue
}

func getInt32(key string, defValue int32) int32 {
	if value := os.Getenv(key); value != "" {
		v, err := strconv.Atoi(value)
		if err == nil {
			return int32(v)
		}
	}

	return defValue
}

func getDuration(key string, defValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		v, err := time.ParseDuration(value)
		if err == nil {
			return v
		}
	}

	return defValue
}
