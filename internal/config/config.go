package config

import (
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	PgMaxConn   string
}

func Load() *Config {
	return &Config{
		AppPort:     getEnv("APP_PORT", "5005"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		PgMaxConn:   getEnv("PG_MAX_CONN", "20"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
