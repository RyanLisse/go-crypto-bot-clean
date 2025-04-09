package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLiteRepository implements the Repository interface for SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(config Config) (*SQLiteRepository, error) {
	// Ensure database path is set
	if config.DatabasePath == "" {
		return nil, fmt.Errorf("database path is required")
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
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

	return &SQLiteRepository{
		db: db,
	}, nil
}

// Execute executes a query without returning any rows
func (r *SQLiteRepository) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows
func (r *SQLiteRepository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
func (r *SQLiteRepository) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}

// BeginTx starts a transaction
func (r *SQLiteRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// Ping verifies a connection to the database is still alive
func (r *SQLiteRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// GetImplementationType returns the type of database implementation
func (r *SQLiteRepository) GetImplementationType() string {
	return "sqlite"
}

// SupportsSynchronization returns whether this implementation supports synchronization
func (r *SQLiteRepository) SupportsSynchronization() bool {
	return false
}

// Synchronize is a no-op for SQLite
func (r *SQLiteRepository) Synchronize(ctx context.Context) error {
	// No-op for SQLite
	return nil
}

// GetLastSyncTimestamp is a no-op for SQLite
func (r *SQLiteRepository) GetLastSyncTimestamp(ctx context.Context) (time.Time, error) {
	// Return current time for SQLite since it doesn't support synchronization
	return time.Now(), nil
}
