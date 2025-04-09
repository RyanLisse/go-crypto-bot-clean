package database

import (
	"context"
	"database/sql"
	"fmt"
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
	// ModeSQLite uses SQLite as the primary database
	ModeSQLite RepositoryMode = "sqlite"
	// ModeTurso uses TursoDB as the primary database
	ModeTurso RepositoryMode = "turso"
	// ModeShadow uses both SQLite and TursoDB in shadow mode
	ModeShadow RepositoryMode = "shadow"
)

// NewRepository creates a new repository based on the provided configuration
func (f *DefaultRepositoryFactory) NewRepository(config Config) (Repository, error) {
	// Determine repository mode
	mode := f.DetermineRepositoryMode(config)

	// Create repository based on mode
	switch mode {
	case ModeSQLite:
		return NewSQLiteRepository(config)
	case ModeTurso:
		return NewTursoRepository(config)
	case ModeShadow:
		return NewShadowRepository(config)
	default:
		return NewSQLiteRepository(config)
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

	// Default to SQLite mode
	return ModeSQLite
}

// GetRepositoryFromEnv creates a repository based on environment configuration
func GetRepositoryFromEnv(config Config) (Repository, error) {
	factory := NewDefaultRepositoryFactory()
	return factory.NewRepository(config)
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
