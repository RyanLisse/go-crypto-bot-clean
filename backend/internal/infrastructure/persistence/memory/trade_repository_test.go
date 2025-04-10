package memory

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestTradeRepository_Store(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	trade := &models.Trade{
		ID:        "trade1",
		Symbol:    "BTC-USD",
		TradeTime: time.Now(),
	}

	err := repo.Store(ctx, trade)
	assert.NoError(t, err)

	// Test storing trade with empty ID
	emptyTrade := &models.Trade{
		Symbol: "BTC-USD",
	}
	err = repo.Store(ctx, emptyTrade)
	assert.Error(t, err)
}

func TestTradeRepository_GetByID(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	trade := &models.Trade{
		ID:        "trade1",
		Symbol:    "BTC-USD",
		TradeTime: time.Now(),
	}

	// Store the trade first
	err := repo.Store(ctx, trade)
	assert.NoError(t, err)

	// Test retrieving the trade
	retrieved, err := repo.GetByID(ctx, "trade1")
	assert.NoError(t, err)
	assert.Equal(t, trade.ID, retrieved.ID)
	assert.Equal(t, trade.Symbol, retrieved.Symbol)

	// Test retrieving non-existent trade
	_, err = repo.GetByID(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestTradeRepository_GetBySymbol(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	// Create test trades
	trades := []*models.Trade{
		{ID: "trade1", Symbol: "BTC-USD", TradeTime: time.Now()},
		{ID: "trade2", Symbol: "BTC-USD", TradeTime: time.Now()},
		{ID: "trade3", Symbol: "ETH-USD", TradeTime: time.Now()},
	}

	for _, trade := range trades {
		err := repo.Store(ctx, trade)
		assert.NoError(t, err)
	}

	// Test getting all BTC-USD trades
	btcTrades, err := repo.GetBySymbol(ctx, "BTC-USD", 0)
	assert.NoError(t, err)
	assert.Len(t, btcTrades, 2)

	// Test with limit
	limitedTrades, err := repo.GetBySymbol(ctx, "BTC-USD", 1)
	assert.NoError(t, err)
	assert.Len(t, limitedTrades, 1)
}

func TestTradeRepository_GetByTimeRange(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	now := time.Now()

	// Create test trades with different timestamps
	trades := []*models.Trade{
		{ID: "trade1", Symbol: "BTC-USD", TradeTime: now.Add(-2 * time.Hour)},
		{ID: "trade2", Symbol: "BTC-USD", TradeTime: now.Add(-1 * time.Hour)},
		{ID: "trade3", Symbol: "BTC-USD", TradeTime: now},
		{ID: "trade4", Symbol: "ETH-USD", TradeTime: now.Add(-1 * time.Hour)},
	}

	for _, trade := range trades {
		err := repo.Store(ctx, trade)
		assert.NoError(t, err)
	}

	// Test getting trades within time range
	rangeTrades, err := repo.GetByTimeRange(ctx, "BTC-USD", now.Add(-1*time.Hour).Add(-1*time.Minute), now.Add(1*time.Minute), 0)
	assert.NoError(t, err)
	assert.Len(t, rangeTrades, 2) // Should include trade2 and trade3

	// Test with limit
	limitedTrades, err := repo.GetByTimeRange(ctx, "BTC-USD", now.Add(-3*time.Hour), now.Add(1*time.Minute), 2)
	assert.NoError(t, err)
	assert.Len(t, limitedTrades, 2) // Should limit to 2 trades
}

func TestTradeRepository_GetByExchange(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	// Create test trades
	trades := []*models.Trade{
		{ID: "trade1", Symbol: "BTC-USD", Exchange: "Binance", TradeTime: time.Now()},
		{ID: "trade2", Symbol: "BTC-USD", Exchange: "Coinbase", TradeTime: time.Now()},
		{ID: "trade3", Symbol: "ETH-USD", Exchange: "Binance", TradeTime: time.Now()},
	}

	for _, trade := range trades {
		err := repo.Store(ctx, trade)
		assert.NoError(t, err)
	}

	// Test getting trades by exchange
	binanceTrades, err := repo.GetByExchange(ctx, "Binance", 0)
	assert.NoError(t, err)
	assert.Len(t, binanceTrades, 2)

	// Test with limit
	limitedTrades, err := repo.GetByExchange(ctx, "Binance", 1)
	assert.NoError(t, err)
	assert.Len(t, limitedTrades, 1)
}

func TestTradeRepository_GetByOrderID(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	// Create test trades
	trades := []*models.Trade{
		{ID: "trade1", Symbol: "BTC-USD", OrderID: "order1", TradeTime: time.Now()},
		{ID: "trade2", Symbol: "BTC-USD", OrderID: "order1", TradeTime: time.Now()},
		{ID: "trade3", Symbol: "ETH-USD", OrderID: "order2", TradeTime: time.Now()},
	}

	for _, trade := range trades {
		err := repo.Store(ctx, trade)
		assert.NoError(t, err)
	}

	// Test getting trades by order ID
	order1Trades, err := repo.GetByOrderID(ctx, "order1")
	assert.NoError(t, err)
	assert.Len(t, order1Trades, 2)

	order2Trades, err := repo.GetByOrderID(ctx, "order2")
	assert.NoError(t, err)
	assert.Len(t, order2Trades, 1)
}

func TestTradeRepository_Delete(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	// Create a trade
	trade := &models.Trade{
		ID:        "trade1",
		Symbol:    "BTC-USD",
		TradeTime: time.Now(),
	}
	err := repo.Store(ctx, trade)
	assert.NoError(t, err)

	// Delete the trade
	err = repo.Delete(ctx, "trade1")
	assert.NoError(t, err)

	// Verify the trade is deleted
	_, err = repo.GetByID(ctx, "trade1")
	assert.Error(t, err)

	// Test deleting non-existent trade
	err = repo.Delete(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestTradeRepository_DeleteOlderThan(t *testing.T) {
	repo := NewTradeRepository()
	ctx := context.Background()

	now := time.Now()

	// Create test trades with different timestamps
	trades := []*models.Trade{
		{ID: "trade1", Symbol: "BTC-USD", TradeTime: now.Add(-2 * time.Hour)},
		{ID: "trade2", Symbol: "BTC-USD", TradeTime: now.Add(-1 * time.Hour)},
		{ID: "trade3", Symbol: "BTC-USD", TradeTime: now},
	}

	for _, trade := range trades {
		err := repo.Store(ctx, trade)
		assert.NoError(t, err)
	}

	// Delete trades older than 1.5 hours ago
	err := repo.DeleteOlderThan(ctx, now.Add(-90*time.Minute))
	assert.NoError(t, err)

	// Verify only trade1 was deleted
	_, err = repo.GetByID(ctx, "trade1")
	assert.Error(t, err)

	// trade2 and trade3 should still exist
	_, err = repo.GetByID(ctx, "trade2")
	assert.NoError(t, err)
	_, err = repo.GetByID(ctx, "trade3")
	assert.NoError(t, err)
}
