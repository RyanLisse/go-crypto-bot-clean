//go:build turso

// Package turso provides a database adapter for Turso
package turso

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/tursodatabase/go-libsql"
)

type TursoDB struct {
	db        *sql.DB
	connector *libsql.Connector
	logger    *zerolog.Logger
	dbPath    string
}

// NewTursoDB creates a new TursoDB instance with embedded replica support
func NewTursoDB(primaryURL, authToken string, syncInterval time.Duration, logger *zerolog.Logger) (*TursoDB, error) {
	// Create a persistent directory for the local database
	dir := "./data/turso"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	// Use a persistent path for the local database
	dbPath := filepath.Join(dir, "local.db")
	logger.Info().Str("path", dbPath).Msg("Using local database for Turso")

	// Create connector with sync interval
	connectorOptions := []libsql.Option{
		libsql.WithAuthToken(authToken),
	}

	// Only add sync interval if it's greater than 0
	if syncInterval > 0 {
		connectorOptions = append(connectorOptions, libsql.WithSyncInterval(syncInterval))
		logger.Info().Dur("interval", syncInterval).Msg("Configured automatic sync interval")
	} else {
		logger.Info().Msg("Automatic sync disabled, manual sync will be required")
	}

	// Create connector
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryURL, connectorOptions...)
	if err != nil {
		return nil, fmt.Errorf("error creating connector: %w", err)
	}

	// Open database connection
	db := sql.OpenDB(connector)
	if err := db.Ping(); err != nil {
		connector.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Perform initial sync
	logger.Info().Msg("Performing initial sync with Turso primary database")
	_, syncErr := connector.Sync()
	if syncErr != nil {
		logger.Warn().Err(syncErr).Msg("Initial sync failed, will retry later")
	} else {
		logger.Info().Msg("Initial sync completed successfully")
	}

	return &TursoDB{
		db:        db,
		connector: connector,
		logger:    logger,
		dbPath:    dbPath,
	}, nil
}

// DB returns the underlying *sql.DB instance
func (t *TursoDB) DB() *sql.DB {
	return t.db
}

// Sync manually syncs the local database with the primary database
// Returns any error that occurred
func (t *TursoDB) Sync() error {
	t.logger.Debug().Msg("Syncing with Turso primary database")
	_, err := t.connector.Sync()
	if err != nil {
		t.logger.Error().Err(err).Msg("Failed to sync with Turso primary database")
		return err
	}

	t.logger.Debug().Msg("Sync completed successfully")
	return nil
}

// Close closes the database connection
func (t *TursoDB) Close() error {
	t.logger.Info().Msg("Closing Turso database connection")

	if err := t.db.Close(); err != nil {
		t.logger.Error().Err(err).Msg("Error closing database")
		return fmt.Errorf("error closing database: %w", err)
	}

	if err := t.connector.Close(); err != nil {
		t.logger.Error().Err(err).Msg("Error closing connector")
		return fmt.Errorf("error closing connector: %w", err)
	}

	t.logger.Info().Msg("Turso database connection closed successfully")
	return nil
}
