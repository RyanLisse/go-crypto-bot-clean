package factory

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/turso"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
)

// NewDB creates a new database connection based on configuration
func NewDB(cfg *config.Config, logger *zerolog.Logger) (*sql.DB, error) {
	// Always try to use Turso first if configured
	if cfg.Database.Turso.URL != "" && cfg.Database.Turso.AuthToken != "" {
		// Get sync interval from environment or use default
		syncIntervalStr := os.Getenv("TURSO_SYNC_INTERVAL_SECONDS")
		syncInterval := 5 * time.Minute // Default sync interval

		if syncIntervalStr != "" {
			syncIntervalSec, err := strconv.Atoi(syncIntervalStr)
			if err == nil && syncIntervalSec > 0 {
				syncInterval = time.Duration(syncIntervalSec) * time.Second
			}
		}

		// Check if sync is enabled
		syncEnabled := true
		if os.Getenv("TURSO_SYNC_ENABLED") == "false" {
			syncEnabled = false
			syncInterval = 0 // Disable automatic sync
		}

		logger.Info().Str("url", cfg.Database.Turso.URL).Bool("sync_enabled", syncEnabled).Dur("sync_interval", syncInterval).Msg("Initializing Turso database")

		tursoDB, err := turso.NewTursoDB(
			cfg.Database.Turso.URL,
			cfg.Database.Turso.AuthToken,
			syncInterval,
			logger,
		)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to create Turso database, falling back to SQLite")
		} else if db := tursoDB.DB(); db != nil {
			// Configure SQLite connection
			db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
			db.SetMaxIdleConns(1)
			db.SetConnMaxLifetime(time.Hour)
			return db, nil
		}
	}

	// If Turso is not configured or failed, return an error
	// We're committed to using Turso for production
	return nil, fmt.Errorf("turso database is required but not properly configured; check your environment variables")
}
