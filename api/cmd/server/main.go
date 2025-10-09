package main

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/stwalsh4118/mirageapi/internal/auth"
	"github.com/stwalsh4118/mirageapi/internal/config"
	"github.com/stwalsh4118/mirageapi/internal/jobs"
	"github.com/stwalsh4118/mirageapi/internal/logging"
	"github.com/stwalsh4118/mirageapi/internal/railway"
	"github.com/stwalsh4118/mirageapi/internal/server"
	"github.com/stwalsh4118/mirageapi/internal/store"
	"github.com/stwalsh4118/mirageapi/internal/vault"
	"gorm.io/gorm"
)

func main() {
	// Load .env if present (dev convenience)
	if err := godotenv.Load(); err == nil {
		log.Info().Msg("loaded .env file")
	}
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logging.Setup(cfg.Environment)

	// Initialize Clerk SDK with secret key
	if cfg.ClerkSecretKey != "" {
		auth.InitClerk(cfg.ClerkSecretKey)
	} else {
		log.Warn().Msg("CLERK_SECRET_KEY not set - authentication will not work")
	}

	// Initialize DB: prefer DATABASE_URL if provided, else fallback to SQLite (SQLITE_PATH or default)
	var db *gorm.DB
	if cfg.DatabaseURL != "" {
		d, derr := store.OpenFromURL(cfg.DatabaseURL)
		if derr != nil {
			log.Fatal().Err(derr).Str("database_url", cfg.DatabaseURL).Msg("failed to init database from url")
		}
		db = d
	} else {
		sqlitePath := os.Getenv("SQLITE_PATH")
		d, derr := store.Open(sqlitePath)
		if derr != nil {
			log.Fatal().Err(derr).Str("path", sqlitePath).Msg("failed to init sqlite database")
		}
		db = d
	}

	rw := railway.NewFromConfig(cfg)

	// Initialize Vault client if enabled
	var vaultClient *vault.Client
	if cfg.VaultEnabled {
		vaultCfg := vault.Config{
			Address:    cfg.VaultAddr,
			Token:      cfg.VaultToken,
			RoleID:     cfg.VaultRoleID,
			SecretID:   cfg.VaultSecretID,
			Namespace:  cfg.VaultNamespace,
			SkipVerify: cfg.VaultSkipVerify,
			MountPath:  cfg.VaultMountPath,
		}
		vc, err := vault.NewClient(vaultCfg)
		if err != nil {
			log.Warn().Err(err).Msg("failed to initialize Vault client, continuing without Vault")
		} else {
			vaultClient = vc
			log.Info().Msg("Vault client initialized successfully")
		}
	}

	// Start status poller (Phase 1)
	pollInterval := time.Duration(cfg.PollIntervalSeconds) * time.Second
	if pollInterval <= 0 {
		log.Warn().Int("poll_interval_seconds", cfg.PollIntervalSeconds).Msg("poll interval invalid; skipping status poller startup")
	} else {
		parent := context.Background()
		ctx, cancel := context.WithCancel(parent)
		pollStop := jobs.StartStatusPoller(
			ctx,
			db,
			rw,
			pollInterval,
			cfg.PollJitterFraction,
			nil, // use default log publisher
		)
		defer func() {
			// Ensure poller goroutine and its ticker are stopped
			if pollStop != nil {
				pollStop()
			}
			cancel()
		}()
	}

	engine := server.NewHTTPServer(cfg, db, rw, vaultClient)

	port := cfg.HTTPPort
	if port == "" {
		port = config.DefaultHTTPPort
	}

	addr := ":" + port
	log.Info().Str("addr", addr).Str("env", cfg.Environment).Msg("starting api")
	if err := engine.Run(addr); err != nil {
		log.Error().Err(err).Msg("server exited with error")
	}
}
