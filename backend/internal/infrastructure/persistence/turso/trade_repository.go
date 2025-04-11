package turso

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"

	"gorm.io/gorm"
)

// TradeRepository implements the ports.TradeRepository interface for TursoDB
type TradeRepository struct {
	*BaseRepository
}

// NewTradeRepository creates a new trade repository instance
func NewTradeRepository(db *gorm.DB) ports.TradeRepository {
	return &TradeRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Store persists a trade in the repository
func (r *TradeRepository) Store(ctx context.Context, trade *models.Trade) error {
	return r.BaseRepository.Create(ctx, trade)
}

// GetByID retrieves a trade by its ID
func (r *TradeRepository) GetByID(ctx context.Context, id string) (*models.Trade, error) {
	var trade models.Trade
	err := r.FindByID(ctx, &trade, id)
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

// GetBySymbol retrieves all trades for a given symbol
func (r *TradeRepository) GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.Trade, error) {
	var trades []*models.Trade
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Limit(limit).Find(&trades)
	if result.Error != nil {
		return nil, result.Error
	}
	return trades, nil
}

// GetByTimeRange retrieves trades within a specific time range
func (r *TradeRepository) GetByTimeRange(ctx context.Context, symbol string, start, end time.Time, limit int) ([]*models.Trade, error) {
	var trades []*models.Trade
	result := r.db.WithContext(ctx).
		Where("symbol = ? AND trade_time BETWEEN ? AND ?", symbol, start, end).
		Limit(limit).
		Find(&trades)
	if result.Error != nil {
		return nil, result.Error
	}
	return trades, nil
}

// GetByExchange retrieves trades from a specific exchange
func (r *TradeRepository) GetByExchange(ctx context.Context, exchange string, limit int) ([]*models.Trade, error) {
	var trades []*models.Trade
	result := r.db.WithContext(ctx).
		Where("exchange = ?", exchange).
		Limit(limit).
		Find(&trades)
	if result.Error != nil {
		return nil, result.Error
	}
	return trades, nil
}

// GetByOrderID retrieves trades associated with a specific order
func (r *TradeRepository) GetByOrderID(ctx context.Context, orderID string) ([]*models.Trade, error) {
	var trades []*models.Trade
	result := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Find(&trades)
	if result.Error != nil {
		return nil, result.Error
	}
	return trades, nil
}

// Delete removes a trade from the repository
func (r *TradeRepository) Delete(ctx context.Context, id string) error {
	trade := &models.Trade{ID: id}
	return r.BaseRepository.Delete(ctx, trade)
}

// DeleteOlderThan removes trades older than the specified time
func (r *TradeRepository) DeleteOlderThan(ctx context.Context, before time.Time) error {
	result := r.db.WithContext(ctx).
		Where("trade_time < ?", before).
		Delete(&models.Trade{})
	return result.Error
}
