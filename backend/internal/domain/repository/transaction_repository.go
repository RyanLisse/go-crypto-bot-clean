package repository

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// Create creates a new transaction record
	Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	
	// FindByID retrieves a transaction by its ID
	FindByID(ctx context.Context, id int64) (*models.Transaction, error)
	
	// FindByTimeRange retrieves transactions within a time range
	FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error)
	
	// FindAll retrieves all transactions
	FindAll(ctx context.Context) ([]*models.Transaction, error)
}
