package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// GormRepository implements the Repository interface using GORM
type GormRepository struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

// NewGormRepository creates a new GORM repository using an existing *gorm.DB connection
func NewGormRepository(db *gorm.DB) (*GormRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("gorm DB instance cannot be nil")
	}

	// Get the underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database from gorm DB: %w", err)
	}

	// PRAGMA settings might be better placed during initial gorm.Open configuration
	// or potentially removed if default GORM/SQLite behavior is sufficient.
	// db.Exec(\"PRAGMA foreign_keys = ON\")
	// db.Exec(\"PRAGMA journal_mode = WAL\")

	return &GormRepository{
		db:    db,
		sqlDB: sqlDB,
	}, nil
}

// Execute executes a query without returning any rows
func (r *GormRepository) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result := r.db.WithContext(ctx).Exec(query, args...)
	return nil, result.Error // GORM doesn't return sql.Result directly
}

// Query executes a query that returns rows
func (r *GormRepository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.WithContext(ctx).Raw(query, args...).Rows()
}

// QueryRow executes a query that returns a single row
func (r *GormRepository) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	// GORM doesn't have direct QueryRow equivalent, so we use the underlying sqlDB
	return r.sqlDB.QueryRowContext(ctx, query, args...)
}

// BeginTx starts a transaction
func (r *GormRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.sqlDB.BeginTx(ctx, opts)
}

// Close closes the database connection
func (r *GormRepository) Close() error {
	// return r.sqlDB.Close() // Connection is likely closed elsewhere
	return nil // Or perhaps log a warning
}

// Ping verifies a connection to the database is still alive
func (r *GormRepository) Ping(ctx context.Context) error {
	return r.sqlDB.PingContext(ctx)
}

// GetImplementationType returns the type of database implementation
func (r *GormRepository) GetImplementationType() string {
	return "gorm-sqlite"
}

// SupportsSynchronization returns whether this implementation supports synchronization
func (r *GormRepository) SupportsSynchronization() bool {
	return false
}

// Synchronize is a no-op for GORM SQLite
func (r *GormRepository) Synchronize(ctx context.Context) error {
	return nil
}

// GetLastSyncTimestamp is a no-op for GORM SQLite
func (r *GormRepository) GetLastSyncTimestamp(ctx context.Context) (time.Time, error) {
	return time.Time{}, nil
}

// AutoMigrate runs GORM's auto-migration for the given models
func (r *GormRepository) AutoMigrate(models ...interface{}) error {
	return r.db.AutoMigrate(models...)
}
