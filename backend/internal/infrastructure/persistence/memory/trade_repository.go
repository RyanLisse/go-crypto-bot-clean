package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"
)

// TradeRepository implements the ports.TradeRepository interface with in-memory storage
type TradeRepository struct {
	mu     sync.RWMutex
	trades map[string]*models.Trade
}

// NewTradeRepository creates a new in-memory trade repository
func NewTradeRepository() ports.TradeRepository {
	return &TradeRepository{
		trades: make(map[string]*models.Trade),
	}
}

// Store persists a trade in the repository
func (r *TradeRepository) Store(ctx context.Context, trade *models.Trade) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if trade.ID == "" {
		return errors.New("trade ID cannot be empty")
	}

	r.trades[trade.ID] = trade
	return nil
}

// GetByID retrieves a trade by its ID
func (r *TradeRepository) GetByID(ctx context.Context, id string) (*models.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	trade, exists := r.trades[id]
	if !exists {
		return nil, errors.New("trade not found")
	}

	return trade, nil
}

// GetBySymbol retrieves all trades for a given symbol with a limit
func (r *TradeRepository) GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Trade
	for _, trade := range r.trades {
		if symbol == "" || trade.Symbol == symbol {
			result = append(result, trade)
		}
	}

	// Apply limit if specified and if results exceed the limit
	if limit > 0 && len(result) > limit {
		return result[:limit], nil
	}

	return result, nil
}

// GetByTimeRange retrieves trades within a specific time range
func (r *TradeRepository) GetByTimeRange(ctx context.Context, symbol string, start, end time.Time, limit int) ([]*models.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Trade
	for _, trade := range r.trades {
		if (symbol == "" || trade.Symbol == symbol) &&
			(start.IsZero() || !trade.TradeTime.Before(start)) &&
			(end.IsZero() || !trade.TradeTime.After(end)) {
			result = append(result, trade)
		}
	}

	// Apply limit if specified and if results exceed the limit
	if limit > 0 && len(result) > limit {
		return result[:limit], nil
	}

	return result, nil
}

// GetByExchange retrieves trades from a specific exchange
func (r *TradeRepository) GetByExchange(ctx context.Context, exchange string, limit int) ([]*models.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Trade
	for _, trade := range r.trades {
		if exchange == "" || trade.Exchange == exchange {
			result = append(result, trade)
		}
	}

	// Apply limit if specified and if results exceed the limit
	if limit > 0 && len(result) > limit {
		return result[:limit], nil
	}

	return result, nil
}

// GetByOrderID retrieves trades associated with a specific order
func (r *TradeRepository) GetByOrderID(ctx context.Context, orderID string) ([]*models.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Trade
	for _, trade := range r.trades {
		if trade.OrderID == orderID {
			result = append(result, trade)
		}
	}

	return result, nil
}

// Delete removes a trade from the repository
func (r *TradeRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.trades[id]; !exists {
		return errors.New("trade not found")
	}

	delete(r.trades, id)
	return nil
}

// DeleteOlderThan removes trades older than the specified time
func (r *TradeRepository) DeleteOlderThan(ctx context.Context, before time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tradesToDelete := []string{}
	for id, trade := range r.trades {
		if trade.TradeTime.Before(before) {
			tradesToDelete = append(tradesToDelete, id)
		}
	}

	for _, id := range tradesToDelete {
		delete(r.trades, id)
	}

	return nil
}

// Update updates an existing trade (helper method, not part of the interface)
func (r *TradeRepository) Update(ctx context.Context, trade *models.Trade) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if trade.ID == "" {
		return errors.New("trade ID cannot be empty")
	}

	if _, exists := r.trades[trade.ID]; !exists {
		return errors.New("trade not found")
	}

	r.trades[trade.ID] = trade
	return nil
}
