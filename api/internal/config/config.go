package config

import (
	"os"
)

const (
	DefaultHTTPPort = "8080"
)

// AppConfig holds runtime configuration for the API service.
type AppConfig struct {
	Environment     string
	HTTPPort        string
	DatabaseURL     string
	RailwayAPIToken string
}

// LoadFromEnv loads configuration from environment variables with defaults.
func LoadFromEnv() (AppConfig, error) {
	cfg := AppConfig{
		Environment:     getEnv("APP_ENV", "development"),
		HTTPPort:        getEnv("HTTP_PORT", DefaultHTTPPort),
		DatabaseURL:     firstNonEmpty(os.Getenv("DATABASE_URL"), os.Getenv("DB_URL")),
		RailwayAPIToken: os.Getenv("RAILWAY_API_TOKEN"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
