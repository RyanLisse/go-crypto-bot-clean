package database

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

// DefaultRepositoryFactory is the default implementation of RepositoryFactory
type DefaultRepositoryFactory struct{}

// NewDefaultRepositoryFactory creates a new default repository factory
func NewDefaultRepositoryFactory() *DefaultRepositoryFactory {
	return &DefaultRepositoryFactory{}
}

// RepositoryMode defines the mode of operation for the repository
type RepositoryMode string

const (
	// ModeGORM uses GORM with SQLite as the primary database
	ModeGORM RepositoryMode = "gorm"
	// ModeSQLite uses SQLite as the primary database
	ModeSQLite RepositoryMode = "sqlite"
	// ModeTurso uses TursoDB as the primary database
	ModeTurso RepositoryMode = "turso"
	// ModeShadow uses both SQLite and TursoDB in shadow mode
	ModeShadow RepositoryMode = "shadow"
)

// NewRepository creates a new repository based on the provided configuration and DB connection
func (f *DefaultRepositoryFactory) NewRepository(config Config, db *gorm.DB) (Repository, error) {
	// Determine repository mode
	mode := f.DetermineRepositoryMode(config)

	// Create repository based on mode
	switch mode {
	case ModeGORM:
		// Use the provided db instance
		return NewGormRepository(db)
	case ModeSQLite:
		// TODO: Refactor NewSQLiteRepository to accept a connection or handle initialization differently
		return NewSQLiteRepository(config)
	case ModeTurso:
		// TODO: Refactor NewTursoRepository to accept a connection or handle initialization differently
		return NewTursoRepository(config)
	case ModeShadow:
		// TODO: Refactor NewShadowRepository to accept connections or handle initialization differently
		return NewShadowRepository(config)
	default:
		// Default to GORM, using the provided db instance
		return NewGormRepository(db)
	}
}

// DetermineRepositoryMode determines the repository mode based on configuration
func (f *DefaultRepositoryFactory) DetermineRepositoryMode(config Config) RepositoryMode {
	// Check for explicit shadow mode
	if config.ShadowMode {
		return ModeShadow
	}

	// If TursoDB is enabled, use TursoDB mode
	if config.TursoEnabled {
		return ModeTurso
	}

	// Default to GORM mode
	return ModeGORM
}

// GetRepositoryFromEnv might need refactoring as it doesn't have access to the *gorm.DB instance
// This function's usage needs to be reviewed in the context of where the DB is initialized.
func GetRepositoryFromEnv(config Config /*, db *gorm.DB - Needs DB instance */) (Repository, error) {
	// Cannot call NewRepository without the db instance
	// return factory.NewRepository(config, db)
	return nil, fmt.Errorf("GetRepositoryFromEnv needs refactoring to receive DB instance")
}

// WithTransaction executes a function within a transaction
// This is a helper function that works with any Repository implementation
func WithTransaction(ctx context.Context, repo Repository, fn func(*sql.Tx) error) error {
	// Begin transaction
	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of error
	// This is a no-op if the transaction is committed
	defer tx.Rollback()

	// Execute the function
	if err := fn(tx); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
