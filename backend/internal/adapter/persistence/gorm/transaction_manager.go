package gorm

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TransactionManager implements port.TransactionManager using GORM
type TransactionManager struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewTransactionManager creates a new TransactionManager
func NewTransactionManager(db *gorm.DB, logger *zerolog.Logger) *TransactionManager {
	return &TransactionManager{
		db:     db,
		logger: logger,
	}
}

// WithTransaction executes the given function within a transaction
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Start a new transaction
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		tm.logger.Error().Err(tx.Error).Msg("Failed to begin transaction")
		return tx.Error
	}

	// Create a new context with the transaction
	txCtx := context.WithValue(ctx, port.TxContextKey, tx)

	// Execute the function
	err := fn(txCtx)
	if err != nil {
		// Rollback the transaction if there was an error
		if rbErr := tx.Rollback().Error; rbErr != nil {
			tm.logger.Error().Err(rbErr).Msg("Failed to rollback transaction")
			// Return the original error, not the rollback error
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tm.logger.Error().Err(err).Msg("Failed to commit transaction")
		return err
	}

	return nil
}

// GetDB returns the database connection
func (tm *TransactionManager) GetDB(ctx context.Context) *gorm.DB {
	// Check if there's a transaction in the context
	if tx, ok := ctx.Value(port.TxContextKey).(*gorm.DB); ok {
		return tx
	}
	// Otherwise, return the regular DB connection
	return tm.db.WithContext(ctx)
}

// Ensure TransactionManager implements port.TransactionManager
var _ port.TransactionManager = (*TransactionManager)(nil)
