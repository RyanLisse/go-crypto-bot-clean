package repositories

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GormTransactionRepository implements the TransactionRepository interface using GORM
type GormTransactionRepository struct {
	db     *gorm.DB
	logger *logging.LoggerWrapper
}

// NewGormTransactionRepository creates a new GormTransactionRepository
func NewGormTransactionRepository(db *gorm.DB, logger *logging.LoggerWrapper) *GormTransactionRepository {
	return &GormTransactionRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new transaction record
func (r *GormTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	// Set timestamps if not already set
	now := time.Now()
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = now
	}
	if transaction.UpdatedAt.IsZero() {
		transaction.UpdatedAt = now
	}

	result := r.db.WithContext(ctx).Create(transaction)
	if result.Error != nil {
		r.logger.Error("Failed to create transaction", zap.Error(result.Error))
		return nil, result.Error
	}

	return transaction, nil
}

// FindByID retrieves a transaction by its ID
func (r *GormTransactionRepository) FindByID(ctx context.Context, id int64) (*models.Transaction, error) {
	var transaction models.Transaction
	result := r.db.WithContext(ctx).First(&transaction, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to find transaction by ID", zap.Error(result.Error), zap.Int64("id", id))
		return nil, result.Error
	}
	return &transaction, nil
}

// FindByTimeRange retrieves transactions within a time range
func (r *GormTransactionRepository) FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	result := r.db.WithContext(ctx).
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Order("created_at DESC").
		Find(&transactions)

	if result.Error != nil {
		r.logger.Error("Failed to find transactions in time range",
			zap.Error(result.Error),
			zap.Time("startTime", startTime),
			zap.Time("endTime", endTime))
		return nil, result.Error
	}
	return transactions, nil
}

// FindAll retrieves all transactions
func (r *GormTransactionRepository) FindAll(ctx context.Context) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&transactions)

	if result.Error != nil {
		r.logger.Error("Failed to find all transactions", zap.Error(result.Error))
		return nil, result.Error
	}
	return transactions, nil
}
