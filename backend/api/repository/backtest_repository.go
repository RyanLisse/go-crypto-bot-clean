package repository

import (
	"context"
	"errors"
	"fmt"

	"go-crypto-bot-clean/backend/api/models"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrBacktestNotFound = errors.New("backtest not found")
)

// BacktestRepository defines the interface for backtest operations
type BacktestRepository interface {
	// Create creates a new backtest
	Create(ctx context.Context, backtest *models.Backtest) error
	
	// GetByID gets a backtest by ID
	GetByID(ctx context.Context, id string) (*models.Backtest, error)
	
	// GetByUserID gets all backtests for a user
	GetByUserID(ctx context.Context, userID string) ([]*models.Backtest, error)
	
	// GetByStrategyID gets all backtests for a strategy
	GetByStrategyID(ctx context.Context, strategyID string) ([]*models.Backtest, error)
	
	// Update updates a backtest
	Update(ctx context.Context, backtest *models.Backtest) error
	
	// Delete deletes a backtest
	Delete(ctx context.Context, id string) error
	
	// AddTrade adds a trade to a backtest
	AddTrade(ctx context.Context, trade *models.BacktestTrade) error
	
	// GetTrades gets all trades for a backtest
	GetTrades(ctx context.Context, backtestID string) ([]*models.BacktestTrade, error)
	
	// AddEquityPoint adds an equity point to a backtest
	AddEquityPoint(ctx context.Context, equity *models.BacktestEquity) error
	
	// GetEquityCurve gets the equity curve for a backtest
	GetEquityCurve(ctx context.Context, backtestID string) ([]*models.BacktestEquity, error)
}

// GormBacktestRepository implements BacktestRepository using GORM
type GormBacktestRepository struct {
	db *gorm.DB
}

// NewGormBacktestRepository creates a new GORM backtest repository
func NewGormBacktestRepository(db *gorm.DB) BacktestRepository {
	return &GormBacktestRepository{
		db: db,
	}
}

// Create creates a new backtest
func (r *GormBacktestRepository) Create(ctx context.Context, backtest *models.Backtest) error {
	if err := r.db.WithContext(ctx).Create(backtest).Error; err != nil {
		return fmt.Errorf("failed to create backtest: %w", err)
	}
	return nil
}

// GetByID gets a backtest by ID
func (r *GormBacktestRepository) GetByID(ctx context.Context, id string) (*models.Backtest, error) {
	var backtest models.Backtest
	if err := r.db.WithContext(ctx).First(&backtest, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBacktestNotFound
		}
		return nil, fmt.Errorf("failed to get backtest by ID: %w", err)
	}
	return &backtest, nil
}

// GetByUserID gets all backtests for a user
func (r *GormBacktestRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Backtest, error) {
	var backtests []*models.Backtest
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&backtests).Error; err != nil {
		return nil, fmt.Errorf("failed to get backtests by user ID: %w", err)
	}
	return backtests, nil
}

// GetByStrategyID gets all backtests for a strategy
func (r *GormBacktestRepository) GetByStrategyID(ctx context.Context, strategyID string) ([]*models.Backtest, error) {
	var backtests []*models.Backtest
	if err := r.db.WithContext(ctx).Where("strategy_id = ?", strategyID).Find(&backtests).Error; err != nil {
		return nil, fmt.Errorf("failed to get backtests by strategy ID: %w", err)
	}
	return backtests, nil
}

// Update updates a backtest
func (r *GormBacktestRepository) Update(ctx context.Context, backtest *models.Backtest) error {
	// Check if backtest exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Backtest{}).Where("id = ?", backtest.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check backtest existence: %w", err)
	}
	if count == 0 {
		return ErrBacktestNotFound
	}

	// Update backtest
	if err := r.db.WithContext(ctx).Save(backtest).Error; err != nil {
		return fmt.Errorf("failed to update backtest: %w", err)
	}
	return nil
}

// Delete deletes a backtest
func (r *GormBacktestRepository) Delete(ctx context.Context, id string) error {
	// Check if backtest exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Backtest{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check backtest existence: %w", err)
	}
	if count == 0 {
		return ErrBacktestNotFound
	}

	// Use a transaction to delete backtest and related data
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete trades
		if err := tx.Where("backtest_id = ?", id).Delete(&models.BacktestTrade{}).Error; err != nil {
			return fmt.Errorf("failed to delete backtest trades: %w", err)
		}

		// Delete equity points
		if err := tx.Where("backtest_id = ?", id).Delete(&models.BacktestEquity{}).Error; err != nil {
			return fmt.Errorf("failed to delete backtest equity points: %w", err)
		}

		// Delete backtest (soft delete)
		if err := tx.Delete(&models.Backtest{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to delete backtest: %w", err)
		}

		return nil
	})
}

// AddTrade adds a trade to a backtest
func (r *GormBacktestRepository) AddTrade(ctx context.Context, trade *models.BacktestTrade) error {
	// Check if backtest exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Backtest{}).Where("id = ?", trade.BacktestID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check backtest existence: %w", err)
	}
	if count == 0 {
		return ErrBacktestNotFound
	}

	// Add trade
	if err := r.db.WithContext(ctx).Create(trade).Error; err != nil {
		return fmt.Errorf("failed to add backtest trade: %w", err)
	}
	return nil
}

// GetTrades gets all trades for a backtest
func (r *GormBacktestRepository) GetTrades(ctx context.Context, backtestID string) ([]*models.BacktestTrade, error) {
	// Check if backtest exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Backtest{}).Where("id = ?", backtestID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check backtest existence: %w", err)
	}
	if count == 0 {
		return nil, ErrBacktestNotFound
	}

	// Get trades
	var trades []*models.BacktestTrade
	if err := r.db.WithContext(ctx).Where("backtest_id = ?", backtestID).Order("entry_time ASC").Find(&trades).Error; err != nil {
		return nil, fmt.Errorf("failed to get backtest trades: %w", err)
	}
	return trades, nil
}

// AddEquityPoint adds an equity point to a backtest
func (r *GormBacktestRepository) AddEquityPoint(ctx context.Context, equity *models.BacktestEquity) error {
	// Check if backtest exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Backtest{}).Where("id = ?", equity.BacktestID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check backtest existence: %w", err)
	}
	if count == 0 {
		return ErrBacktestNotFound
	}

	// Add equity point
	if err := r.db.WithContext(ctx).Create(equity).Error; err != nil {
		return fmt.Errorf("failed to add backtest equity point: %w", err)
	}
	return nil
}

// GetEquityCurve gets the equity curve for a backtest
func (r *GormBacktestRepository) GetEquityCurve(ctx context.Context, backtestID string) ([]*models.BacktestEquity, error) {
	// Check if backtest exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Backtest{}).Where("id = ?", backtestID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check backtest existence: %w", err)
	}
	if count == 0 {
		return nil, ErrBacktestNotFound
	}

	// Get equity curve
	var equity []*models.BacktestEquity
	if err := r.db.WithContext(ctx).Where("backtest_id = ?", backtestID).Order("timestamp ASC").Find(&equity).Error; err != nil {
		return nil, fmt.Errorf("failed to get backtest equity curve: %w", err)
	}
	return equity, nil
}
