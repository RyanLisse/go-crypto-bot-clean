package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// TursoRepository implements the Repository interface for TursoDB
type TursoRepository struct {
	db     *sql.DB
	config Config
}

// NewTursoRepository creates a new TursoDB repository
// This is a simplified implementation for demonstration purposes
// In a real-world scenario, you would use the TursoDB client library
func NewTursoRepository(config Config) (*TursoRepository, error) {
	// Validate configuration
	if !config.TursoEnabled {
		return nil, fmt.Errorf("turso is not enabled in configuration")
	}

	// Validate required configuration
	if config.TursoURL == "" {
		return nil, fmt.Errorf("turso URL is required")
	}
	if config.TursoAuthToken == "" {
		return nil, fmt.Errorf("turso auth token is required")
	}

	// For demonstration purposes, we'll use SQLite as a stand-in for TursoDB
	// In a real implementation, you would use the TursoDB client library
	dbPath := config.DatabasePath
	if dbPath == "" {
		dbPath = ":memory:"
	} else {
		dbPath = dbPath + ".turso"
	}

	// Open SQLite database (simulating TursoDB)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open TursoDB database: %w", err)
	}

	// Configure connection pool
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Enable WAL mode for better concurrency
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to TursoDB: %w", err)
	}

	return &TursoRepository{
		db:     db,
		config: config,
	}, nil
}

// Execute executes a query without returning any rows
func (r *TursoRepository) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows
func (r *TursoRepository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
func (r *TursoRepository) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}

// BeginTx starts a transaction
func (r *TursoRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}

// Close closes the database connection
func (r *TursoRepository) Close() error {
	return r.db.Close()
}

// Ping verifies a connection to the database is still alive
func (r *TursoRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// GetImplementationType returns the type of database implementation
func (r *TursoRepository) GetImplementationType() string {
	return "turso"
}

// SupportsSynchronization returns whether this implementation supports synchronization
func (r *TursoRepository) SupportsSynchronization() bool {
	return r.config.SyncEnabled
}

// Synchronize triggers synchronization with the cloud database
func (r *TursoRepository) Synchronize(ctx context.Context) error {
	if !r.config.SyncEnabled {
		return nil
	}

	// In a real implementation, you would call the TursoDB sync function
	// For demonstration purposes, we'll just log a message
	fmt.Println("Simulating TursoDB synchronization...")

	// Update sync timestamp
	_, err := r.Execute(ctx, "CREATE TABLE IF NOT EXISTS turso_sync (timestamp INTEGER)")
	if err != nil {
		return fmt.Errorf("failed to create sync table: %w", err)
	}

	// Delete old timestamps
	_, err = r.Execute(ctx, "DELETE FROM turso_sync")
	if err != nil {
		return fmt.Errorf("failed to clear sync table: %w", err)
	}

	// Insert new timestamp
	_, err = r.Execute(ctx, "INSERT INTO turso_sync (timestamp) VALUES (?)", time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to update sync timestamp: %w", err)
	}

	return nil
}

// GetLastSyncTimestamp gets the timestamp of the last synchronization
func (r *TursoRepository) GetLastSyncTimestamp(ctx context.Context) (time.Time, error) {
	if !r.config.SyncEnabled {
		return time.Now(), nil
	}

	// Create table if it doesn't exist
	_, err := r.Execute(ctx, "CREATE TABLE IF NOT EXISTS turso_sync (timestamp INTEGER)")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to create sync table: %w", err)
	}

	// Query the last sync timestamp
	var timestamp int64
	err = r.QueryRow(ctx, "SELECT timestamp FROM turso_sync ORDER BY timestamp DESC LIMIT 1").Scan(&timestamp)
	if err != nil {
		// If no timestamp is found, return current time
		if err == sql.ErrNoRows {
			return time.Now(), nil
		}
		return time.Time{}, fmt.Errorf("failed to get last sync timestamp: %w", err)
	}

	// Convert Unix timestamp to time.Time
	return time.Unix(timestamp, 0), nil
}
