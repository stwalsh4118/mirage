package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	DefaultHTTPPort = "8080"
	// Poller defaults
	DefaultPollIntervalSeconds = 0
	DefaultPollJitterFraction  = 0.2
	// CORS defaults
	DefaultAllowedOrigins = "http://localhost:3000,http://127.0.0.1:3000,http://localhost:3002"
)

// AppConfig holds runtime configuration for the API service.
type AppConfig struct {
	Environment      string
	HTTPPort         string
	DatabaseURL      string
	RailwayAPIToken  string
	RailwayProjectID string
	RailwayEndpoint  string
	// CORS configuration
	AllowedOrigins []string
	// Status poller settings
	PollIntervalSeconds int
	PollJitterFraction  float64
	// Clerk authentication
	ClerkSecretKey     string
	ClerkWebhookSecret string
	// Vault configuration
	VaultEnabled    bool
	VaultAddr       string
	VaultToken      string
	VaultRoleID     string
	VaultSecretID   string
	VaultNamespace  string
	VaultSkipVerify bool
	VaultMountPath  string
}

// LoadFromEnv loads configuration from environment variables with defaults.
func LoadFromEnv() (AppConfig, error) {
	cfg := AppConfig{
		Environment:         getEnv("APP_ENV", "development"),
		HTTPPort:            getEnv("HTTP_PORT", DefaultHTTPPort),
		DatabaseURL:         firstNonEmpty(os.Getenv("DATABASE_URL"), os.Getenv("DB_URL")),
		RailwayAPIToken:     os.Getenv("RAILWAY_API_TOKEN"),
		RailwayProjectID:    os.Getenv("RAILWAY_PROJECT_ID"),
		RailwayEndpoint:     os.Getenv("RAILWAY_GRAPHQL_ENDPOINT"),
		AllowedOrigins:      parseAllowedOrigins(getEnv("ALLOWED_ORIGINS", DefaultAllowedOrigins)),
		PollIntervalSeconds: getEnvInt("POLL_INTERVAL_SECONDS", DefaultPollIntervalSeconds),
		PollJitterFraction:  getEnvFloat("POLL_JITTER_FRACTION", DefaultPollJitterFraction),
		ClerkSecretKey:      os.Getenv("CLERK_SECRET_KEY"),
		ClerkWebhookSecret:  os.Getenv("CLERK_WEBHOOK_SECRET"),
		VaultEnabled:        getEnvBool("VAULT_ENABLED", false),
		VaultAddr:           os.Getenv("VAULT_ADDR"),
		VaultToken:          os.Getenv("VAULT_TOKEN"),
		VaultRoleID:         os.Getenv("VAULT_ROLE_ID"),
		VaultSecretID:       os.Getenv("VAULT_SECRET_ID"),
		VaultNamespace:      os.Getenv("VAULT_NAMESPACE"),
		VaultSkipVerify:     getEnvBool("VAULT_SKIP_VERIFY", false),
		VaultMountPath:      getEnv("VAULT_MOUNT_PATH", "mirage"),
	}

	// Clamp and validate poller configuration
	if cfg.PollIntervalSeconds <= 0 {
		old := cfg.PollIntervalSeconds
		cfg.PollIntervalSeconds = DefaultPollIntervalSeconds
		log.Warn().Int("old", old).Int("new", cfg.PollIntervalSeconds).Msg("invalid PollIntervalSeconds; using default")
	}
	if cfg.PollJitterFraction < 0 || !(cfg.PollJitterFraction >= 0) { // also guards NaN
		old := cfg.PollJitterFraction
		cfg.PollJitterFraction = 0
		log.Warn().Float64("old", old).Float64("new", cfg.PollJitterFraction).Msg("invalid PollJitterFraction; clamped to 0")
	} else if cfg.PollJitterFraction >= 1 {
		old := cfg.PollJitterFraction
		cfg.PollJitterFraction = 0.999
		log.Warn().Float64("old", old).Float64("new", cfg.PollJitterFraction).Msg("PollJitterFraction too high; clamped below 1")
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

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return fallback
}

// parseAllowedOrigins parses a comma-separated list of origins and returns a slice.
// Trims whitespace from each origin.
func parseAllowedOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}
	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
