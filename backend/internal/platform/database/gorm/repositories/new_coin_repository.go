package repositories

import (
	"context"
	"errors"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GORMNewCoinRepository implements the NewCoinRepository interface using GORM
type GORMNewCoinRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewGORMNewCoinRepository creates a new GORM-based NewCoinRepository
func NewGORMNewCoinRepository(db *gorm.DB, logger *zap.Logger) repositories.NewCoinRepository {
	return &GORMNewCoinRepository{
		db:     db,
		logger: logger,
	}
}

// FindAll returns all new coins (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindAll(ctx context.Context) ([]*models.NewCoin, error) {
	var coins []*models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find all new coins", zap.Error(result.Error))
		return nil, result.Error
	}
	return coins, nil
}

// FindBySymbol returns a new coin by symbol (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error) {
	var coin models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&coin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found
		}
		r.logger.Error("Failed to find new coin by symbol", zap.String("symbol", symbol), zap.Error(result.Error))
		return nil, result.Error
	}
	return &coin, nil
}

// Save saves a new coin
func (r *GORMNewCoinRepository) Save(ctx context.Context, coin *models.NewCoin) error {
	result := r.db.WithContext(ctx).Save(coin)
	if result.Error != nil {
		r.logger.Error("Failed to save new coin", zap.String("symbol", coin.Symbol), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// Delete marks a new coin as deleted using GORM's soft delete
func (r *GORMNewCoinRepository) Delete(ctx context.Context, symbol string) error {
	// GORM sets the DeletedAt field when db.Delete is called
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&models.NewCoin{})
	if result.Error != nil {
		r.logger.Error("Failed to soft delete new coin by symbol", zap.String("symbol", symbol), zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Warn("Soft delete attempted on non-existent new coin symbol", zap.String("symbol", symbol))
		// return gorm.ErrRecordNotFound // Or return nil
	}
	return nil
}

// FindByID returns a new coin by ID (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindByID(ctx context.Context, id int64) (*models.NewCoin, error) {
	var coin models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&coin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found
		}
		r.logger.Error("Failed to find new coin by ID", zap.Int64("id", id), zap.Error(result.Error))
		return nil, result.Error
	}
	return &coin, nil
}

// MarkAsProcessed updates the is_processed flag (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) MarkAsProcessed(ctx context.Context, id int64) error {
	// GORM automatically adds WHERE deleted_at IS NULL for updates via Model()
	result := r.db.WithContext(ctx).Model(&models.NewCoin{}).
		Where("id = ?", id).
		Update("is_processed", true) // Use Update for single field
	if result.Error != nil {
		r.logger.Error("Failed to mark new coin as processed", zap.Int64("id", id), zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// FindByDateRange returns new coins found within a date range (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.NewCoin, error) {
	var coins []*models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).
		Where("found_at BETWEEN ? AND ?", startDate, endDate).
		Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find new coins by date range",
			zap.Time("startDate", startDate),
			zap.Time("endDate", endDate),
			zap.Error(result.Error))
		return nil, result.Error
	}
	return coins, nil
}

// FindByDate returns new coins found on a specific date (Relies on FindByDateRange, which handles soft delete)
func (r *GORMNewCoinRepository) FindByDate(ctx context.Context, date time.Time) ([]*models.NewCoin, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	return r.FindByDateRange(ctx, startOfDay, endOfDay)
}

// FindUpcoming returns upcoming new coins (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindUpcoming(ctx context.Context) ([]*models.NewCoin, error) {
	var coins []*models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).
		Where("is_upcoming = ?", true).
		Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find upcoming new coins", zap.Error(result.Error))
		return nil, result.Error
	}
	return coins, nil
}

// FindUpcomingByDate returns upcoming new coins for a specific date (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindUpcomingByDate(ctx context.Context, date time.Time) ([]*models.NewCoin, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var coins []*models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).
		Where("is_upcoming = ? AND first_open_time BETWEEN ? AND ?",
			true, startOfDay, endOfDay).
		Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find upcoming new coins by date",
			zap.Time("date", date),
			zap.Error(result.Error))
		return nil, result.Error
	}
	return coins, nil
}

// FindTradable returns tradable new coins (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindTradable(ctx context.Context) ([]*models.NewCoin, error) {
	var coins []*models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).
		Where("status = ?", "1").
		Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find tradable new coins", zap.Error(result.Error))
		return nil, result.Error
	}
	return coins, nil
}

// FindTradableToday returns tradable new coins for today (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) FindTradableToday(ctx context.Context) ([]*models.NewCoin, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var coins []*models.NewCoin
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).
		Where("status = ? AND became_tradable_at BETWEEN ? AND ?",
			"1", startOfDay, endOfDay).
		Find(&coins)
	if result.Error != nil {
		r.logger.Error("Failed to find tradable new coins for today", zap.Error(result.Error))
		return nil, result.Error
	}
	return coins, nil
}

// SaveAll saves multiple new coins
func (r *GORMNewCoinRepository) SaveAll(ctx context.Context, coins []*models.NewCoin) error {
	if len(coins) == 0 {
		return nil
	}

	// Use a transaction for better performance
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, coin := range coins {
			if err := tx.Save(coin).Error; err != nil {
				r.logger.Error("Failed to save new coin in batch",
					zap.String("symbol", coin.Symbol),
					zap.Error(err))
				return err
			}
		}
		return nil
	})

	if err != nil {
		r.logger.Error("Failed to save all new coins", zap.Error(err))
		return err
	}

	return nil
}

// Count returns the count of new coins (GORM automatically handles non-deleted records)
func (r *GORMNewCoinRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	// GORM automatically adds WHERE deleted_at IS NULL
	result := r.db.WithContext(ctx).Model(&models.NewCoin{}).Count(&count)
	if result.Error != nil {
		r.logger.Error("Failed to count new coins", zap.Error(result.Error))
		return 0, result.Error
	}
	return count, nil
}
