package repositories

import (
	"context"
	"errors"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GORMBoughtCoinRepository implements the BoughtCoinRepository interface using GORM
type GORMBoughtCoinRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewGORMBoughtCoinRepository creates a new GORM-based BoughtCoinRepository
func NewGORMBoughtCoinRepository(db *gorm.DB, logger *zap.Logger) repositories.BoughtCoinRepository {
	return &GORMBoughtCoinRepository{
		db:     db,
		logger: logger,
	}
}

// FindAll returns all bought coins (GORM automatically handles non-deleted records)
func (r *GORMBoughtCoinRepository) FindAll(ctx context.Context) ([]*models.BoughtCoin, error) {
	var coins []*models.BoughtCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find all bought coins", zap.Error(result.Error))
		return nil, result.Error
	}

	// Set BuyPrice alias for backward compatibility
	for _, coin := range coins {
		coin.BuyPrice = coin.PurchasePrice
	}

	return coins, nil
}

// FindBySymbol returns a bought coin by symbol (GORM automatically handles non-deleted records)
func (r *GORMBoughtCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
	var coin models.BoughtCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&coin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found
		}
		r.logger.Error("Failed to find bought coin by symbol", zap.String("symbol", symbol), zap.Error(result.Error))
		return nil, result.Error
	}

	// Set BuyPrice alias for backward compatibility
	coin.BuyPrice = coin.PurchasePrice

	return &coin, nil
}

// Save saves a bought coin
func (r *GORMBoughtCoinRepository) Save(ctx context.Context, coin *models.BoughtCoin) error {
	// Set PurchasePrice from BuyPrice if needed
	if coin.PurchasePrice == 0 && coin.BuyPrice > 0 {
		coin.PurchasePrice = coin.BuyPrice
	}

	result := r.db.WithContext(ctx).Save(coin)
	if result.Error != nil {
		r.logger.Error("Failed to save bought coin", zap.String("symbol", coin.Symbol), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// Delete marks a bought coin as deleted using GORM's soft delete
func (r *GORMBoughtCoinRepository) Delete(ctx context.Context, symbol string) error {
	// GORM sets the DeletedAt field when db.Delete is called
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&models.BoughtCoin{})
	if result.Error != nil {
		r.logger.Error("Failed to soft delete bought coin by symbol", zap.String("symbol", symbol), zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Optionally return an error or log if the record to delete wasn't found
		r.logger.Warn("Soft delete attempted on non-existent bought coin symbol", zap.String("symbol", symbol))
		// return gorm.ErrRecordNotFound // Or return nil if not finding is ok
	}
	return nil
}

// UpdatePrice updates the current price of a bought coin (GORM automatically handles non-deleted records)
func (r *GORMBoughtCoinRepository) UpdatePrice(ctx context.Context, symbol string, price float64) error {
	// GORM automatically adds WHERE deleted_at IS NULL for updates via Model()
	result := r.db.WithContext(ctx).Model(&models.BoughtCoin{}).
		Where("symbol = ?", symbol).
		Updates(map[string]interface{}{
			"current_price": price,
			// updated_at is handled automatically by GORM
		})
	if result.Error != nil {
		r.logger.Error("Failed to update bought coin price", zap.String("symbol", symbol), zap.Float64("price", price), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// FindAllActive returns all active bought coins (GORM automatically handles non-deleted records)
func (r *GORMBoughtCoinRepository) FindAllActive(ctx context.Context) ([]*models.BoughtCoin, error) {
	// This function becomes identical to FindAll with GORM soft delete
	return r.FindAll(ctx)
}

// FindByID returns a bought coin by ID (GORM automatically handles non-deleted records)
func (r *GORMBoughtCoinRepository) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
	var coin models.BoughtCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&coin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found
		}
		r.logger.Error("Failed to find bought coin by ID", zap.Int64("id", id), zap.Error(result.Error))
		return nil, result.Error
	}

	// Set BuyPrice alias for backward compatibility
	coin.BuyPrice = coin.PurchasePrice

	return &coin, nil
}

// DeleteByID deletes a bought coin by ID using GORM's soft delete
func (r *GORMBoughtCoinRepository) DeleteByID(ctx context.Context, id int64) error {
	// GORM sets the DeletedAt field when db.Delete is called
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.BoughtCoin{})
	if result.Error != nil {
		r.logger.Error("Failed to soft delete bought coin by ID", zap.Int64("id", id), zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Warn("Soft delete attempted on non-existent bought coin ID", zap.Int64("id", id))
		// return gorm.ErrRecordNotFound // Or return nil
	}
	return nil
}

// HardDelete permanently deletes a bought coin
func (r *GORMBoughtCoinRepository) HardDelete(ctx context.Context, symbol string) error {
	// Use Unscoped() to bypass the soft delete hook for permanent deletion
	result := r.db.WithContext(ctx).Unscoped().Where("symbol = ?", symbol).Delete(&models.BoughtCoin{})
	if result.Error != nil {
		r.logger.Error("Failed to hard delete bought coin", zap.String("symbol", symbol), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// Count returns the count of bought coins (GORM automatically handles non-deleted records)
func (r *GORMBoughtCoinRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Model(&models.BoughtCoin{}).Count(&count)
	if result.Error != nil {
		r.logger.Error("Failed to count bought coins", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}
