package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repositories"
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

// FindAll returns all bought coins
func (r *GORMBoughtCoinRepository) FindAll(ctx context.Context) ([]*models.BoughtCoin, error) {
	var coins []*models.BoughtCoin
	result := r.db.WithContext(ctx).Where("is_deleted = ?", false).Find(&coins)
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

// FindBySymbol returns a bought coin by symbol
func (r *GORMBoughtCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
	var coin models.BoughtCoin
	result := r.db.WithContext(ctx).Where("symbol = ? AND is_deleted = ?", symbol, false).First(&coin)
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

// Delete marks a bought coin as deleted
func (r *GORMBoughtCoinRepository) Delete(ctx context.Context, symbol string) error {
	result := r.db.WithContext(ctx).Model(&models.BoughtCoin{}).
		Where("symbol = ?", symbol).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		r.logger.Error("Failed to delete bought coin", zap.String("symbol", symbol), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// UpdatePrice updates the current price of a bought coin
func (r *GORMBoughtCoinRepository) UpdatePrice(ctx context.Context, symbol string, price float64) error {
	result := r.db.WithContext(ctx).Model(&models.BoughtCoin{}).
		Where("symbol = ? AND is_deleted = ?", symbol, false).
		Updates(map[string]interface{}{
			"current_price": price,
			"updated_at":    time.Now(),
		})
	if result.Error != nil {
		r.logger.Error("Failed to update bought coin price", zap.String("symbol", symbol), zap.Float64("price", price), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// FindAllActive returns all active bought coins
func (r *GORMBoughtCoinRepository) FindAllActive(ctx context.Context) ([]*models.BoughtCoin, error) {
	var coins []*models.BoughtCoin
	result := r.db.WithContext(ctx).Where("is_deleted = ?", false).Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find all active bought coins", zap.Error(result.Error))
		return nil, result.Error
	}

	// Set BuyPrice alias for backward compatibility
	for _, coin := range coins {
		coin.BuyPrice = coin.PurchasePrice
	}

	return coins, nil
}

// FindByID returns a bought coin by ID
func (r *GORMBoughtCoinRepository) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
	var coin models.BoughtCoin
	result := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&coin)
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

// DeleteByID deletes a bought coin by ID
func (r *GORMBoughtCoinRepository) DeleteByID(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&models.BoughtCoin{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		r.logger.Error("Failed to delete bought coin by ID", zap.Int64("id", id), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// HardDelete permanently deletes a bought coin
func (r *GORMBoughtCoinRepository) HardDelete(ctx context.Context, symbol string) error {
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&models.BoughtCoin{})
	if result.Error != nil {
		r.logger.Error("Failed to hard delete bought coin", zap.String("symbol", symbol), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// Count returns the count of bought coins
func (r *GORMBoughtCoinRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.BoughtCoin{}).Where("is_deleted = ?", false).Count(&count)
	if result.Error != nil {
		r.logger.Error("Failed to count bought coins", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}
