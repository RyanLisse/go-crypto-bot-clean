package database

import (
	"context"
	"database/sql"
	"time"
)

// Repository defines the interface for database operations
// This abstraction allows us to switch between SQLite and TursoDB implementations
type Repository interface {
	// Core database operations
	Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Transaction support
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// Connection management
	Close() error
	Ping(ctx context.Context) error

	// Database-specific operations
	GetImplementationType() string
	SupportsSynchronization() bool

	// For TursoDB, these will trigger synchronization
	// For SQLite, these will be no-ops
	Synchronize(ctx context.Context) error
	GetLastSyncTimestamp(ctx context.Context) (time.Time, error)
}

// Config holds configuration for database connections
type Config struct {
	// Common configuration
	DatabasePath    string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration

	// TursoDB specific configuration
	TursoEnabled   bool
	TursoURL       string
	TursoAuthToken string
	SyncEnabled    bool
	SyncInterval   time.Duration

	// Shadow mode configuration
	ShadowMode bool // When true, writes to both SQLite and TursoDB, but reads from SQLite
}

// RepositoryFactory creates database repositories
type RepositoryFactory interface {
	// Create a new repository instance
	NewRepository(config Config) (Repository, error)
}
