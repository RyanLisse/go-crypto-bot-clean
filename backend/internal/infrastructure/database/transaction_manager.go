package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TransactionManager provides support for executing functions within a database transaction.
// It encapsulates Begin, Commit, and Rollback logic to ensure consistency.
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new TransactionManager.
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// Exec executes the provided function fn within a transaction scope. If fn returns an error,
// the transaction is rolled back and the error is returned. Otherwise, the transaction
// is committed.
func (t *TransactionManager) Exec(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := t.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("transaction begin failed: %w", tx.Error)
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}
