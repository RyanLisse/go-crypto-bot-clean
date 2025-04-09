package database

import (
	"os"
	"strconv"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/config"
)

// LoadConfig loads database configuration from config file and environment variables
func LoadConfig(cfg *config.Config) Config {
	dbConfig := Config{
		// Default values
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		SyncInterval:    5 * time.Minute,
	}

	// Load from config file if available
	if cfg != nil {
		// SQLite configuration
		if cfg.Database.Path != "" {
			dbConfig.DatabasePath = cfg.Database.Path
		}
		if cfg.Database.MaxOpenConns > 0 {
			dbConfig.MaxOpenConns = cfg.Database.MaxOpenConns
		}
		if cfg.Database.MaxIdleConns > 0 {
			dbConfig.MaxIdleConns = cfg.Database.MaxIdleConns
		}
		if cfg.Database.ConnMaxLifetimeSeconds > 0 {
			dbConfig.ConnMaxLifetime = time.Duration(cfg.Database.ConnMaxLifetimeSeconds) * time.Second
		}

		// TursoDB configuration
		if cfg.Database.Turso.Enabled {
			dbConfig.TursoEnabled = true
			dbConfig.TursoURL = cfg.Database.Turso.URL
			dbConfig.TursoAuthToken = cfg.Database.Turso.AuthToken
			dbConfig.SyncEnabled = cfg.Database.Turso.SyncEnabled
			if cfg.Database.Turso.SyncIntervalSeconds > 0 {
				dbConfig.SyncInterval = time.Duration(cfg.Database.Turso.SyncIntervalSeconds) * time.Second
			}
		}

		// Shadow mode configuration
		dbConfig.ShadowMode = cfg.Database.ShadowMode
	}

	// Override with environment variables if set (for development and testing)
	if path := os.Getenv("DB_PATH"); path != "" {
		dbConfig.DatabasePath = path
	}

	if maxOpenStr := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenStr != "" {
		if maxOpen, err := strconv.Atoi(maxOpenStr); err == nil {
			dbConfig.MaxOpenConns = maxOpen
		}
	}

	if maxIdleStr := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleStr != "" {
		if maxIdle, err := strconv.Atoi(maxIdleStr); err == nil {
			dbConfig.MaxIdleConns = maxIdle
		}
	}

	if maxLifeStr := os.Getenv("DB_CONN_MAX_LIFETIME_SECONDS"); maxLifeStr != "" {
		if maxLife, err := strconv.Atoi(maxLifeStr); err == nil {
			dbConfig.ConnMaxLifetime = time.Duration(maxLife) * time.Second
		}
	}

	// TursoDB specific environment variables (override config file)
	if os.Getenv("TURSO_ENABLED") == "true" {
		dbConfig.TursoEnabled = true
	}

	if tursoURL := os.Getenv("TURSO_URL"); tursoURL != "" {
		dbConfig.TursoURL = tursoURL
	}

	if tursoToken := os.Getenv("TURSO_AUTH_TOKEN"); tursoToken != "" {
		dbConfig.TursoAuthToken = tursoToken
	}

	if os.Getenv("TURSO_SYNC_ENABLED") == "true" {
		dbConfig.SyncEnabled = true
	}

	if syncIntervalStr := os.Getenv("TURSO_SYNC_INTERVAL_SECONDS"); syncIntervalStr != "" {
		if syncInterval, err := strconv.Atoi(syncIntervalStr); err == nil {
			dbConfig.SyncInterval = time.Duration(syncInterval) * time.Second
		}
	}

	// Shadow mode configuration
	if os.Getenv("DB_SHADOW_MODE") == "true" {
		dbConfig.ShadowMode = true
	}

	return dbConfig
}
