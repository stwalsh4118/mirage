package config

import (
	"os"
	"strconv"
)

const (
	DefaultHTTPPort = "8080"
	// Poller defaults
	DefaultPollIntervalSeconds = 5
	DefaultPollJitterFraction  = 0.2
)

// AppConfig holds runtime configuration for the API service.
type AppConfig struct {
	Environment     string
	HTTPPort        string
	DatabaseURL     string
	RailwayAPIToken string
	// Status poller settings
	PollIntervalSeconds int
	PollJitterFraction  float64
}

// LoadFromEnv loads configuration from environment variables with defaults.
func LoadFromEnv() (AppConfig, error) {
	cfg := AppConfig{
		Environment:         getEnv("APP_ENV", "development"),
		HTTPPort:            getEnv("HTTP_PORT", DefaultHTTPPort),
		DatabaseURL:         firstNonEmpty(os.Getenv("DATABASE_URL"), os.Getenv("DB_URL")),
		RailwayAPIToken:     os.Getenv("RAILWAY_API_TOKEN"),
		PollIntervalSeconds: getEnvInt("POLL_INTERVAL_SECONDS", DefaultPollIntervalSeconds),
		PollJitterFraction:  getEnvFloat("POLL_JITTER_FRACTION", DefaultPollJitterFraction),
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

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
