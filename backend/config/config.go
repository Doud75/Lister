package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	RateLimitEnabled bool
	RedisURL         string
}

func Load() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	rateLimitEnabled := true
	if val := os.Getenv("RATE_LIMIT_ENABLED"); val == "false" {
		rateLimitEnabled = false
	}

	return Config{
		DatabaseURL:      dbURL,
		JWTSecret:        jwtSecret,
		RateLimitEnabled: rateLimitEnabled,
		RedisURL:         os.Getenv("REDIS_URL"),
	}
}
