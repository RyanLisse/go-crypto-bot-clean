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

// TursoDB provides a database adapter for Turso
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

	// Create connector with auth token
	connector, err := libsql.NewEmbeddedReplicaConnector(
		dbPath,
		primaryURL,
		libsql.WithAuthToken(authToken),
	)
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
	result, syncErr := connector.Sync()
	if syncErr != nil {
		logger.Warn().Err(syncErr).Msg("Initial sync failed, will retry later")
	} else {
		logger.Info().Int("frames_synced", result.FramesSynced).Msg("Initial sync completed successfully")
	}

	// Configure connection pool
	db.SetMaxOpenConns(10) // Adjust based on your needs
	db.SetMaxIdleConns(5)  // Adjust based on your needs
	db.SetConnMaxLifetime(time.Hour)

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
	result, err := t.connector.Sync()
	if err != nil {
		t.logger.Error().Err(err).Msg("Failed to sync with Turso primary database")
		return fmt.Errorf("failed to sync with primary database: %w", err)
	}

	t.logger.Debug().Int("frames_synced", result.FramesSynced).Msg("Sync completed successfully")
	return nil
}

// Close closes the database connection
func (t *TursoDB) Close() error {
	if t.db != nil {
		if err := t.db.Close(); err != nil {
			t.logger.Error().Err(err).Msg("Error closing database connection")
			return err
		}
	}

	if t.connector != nil {
		if err := t.connector.Close(); err != nil {
			t.logger.Error().Err(err).Msg("Error closing connector")
			return err
		}
	}

	t.logger.Info().Msg("Database connection closed")
	return nil
}
