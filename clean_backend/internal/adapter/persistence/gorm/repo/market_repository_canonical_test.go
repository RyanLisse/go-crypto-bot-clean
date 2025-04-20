package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestCanonicalRepository(t *testing.T) (*MarketRepositoryCanonical, func()) {
	// Setup test DB
	db, cleanup := setupTestDB(t)

	// Create a logger
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create the canonical repository
	repo := NewMarketRepositoryCanonical(db, &logger)

	return repo, cleanup
}

func TestMarketRepositoryCanonical(t *testing.T) {
	// Setup canonical repository
	canonicalRepo, cleanup := setupTestCanonicalRepository(t)
	defer cleanup()

	ctx := context.Background()

	// Test SaveTicker and GetTicker
	t.Run("SaveTicker and GetTicker", func(t *testing.T) {
		ticker := &model.Ticker{
			ID:        "test-ticker",
			Symbol:    "BTCUSDT",
			Exchange:  "test-exchange",
			LastPrice: 50000.0,
			Timestamp: time.Now().Round(time.Millisecond),
		}

		// Save the ticker
		err := canonicalRepo.SaveTicker(ctx, ticker)
		require.NoError(t, err)

		// Retrieve the ticker
		retrievedTicker, err := canonicalRepo.GetTicker(ctx, "BTCUSDT", "test-exchange")
		require.NoError(t, err)

		// Verify the ticker was saved correctly
		assert.Equal(t, ticker.ID, retrievedTicker.ID)
		assert.Equal(t, ticker.Symbol, retrievedTicker.Symbol)
		assert.Equal(t, ticker.Exchange, retrievedTicker.Exchange)
		assert.Equal(t, ticker.LastPrice, retrievedTicker.LastPrice)
		assert.Equal(t, ticker.Timestamp.Unix(), retrievedTicker.Timestamp.Unix())
	})

	// Test SaveKline and GetKline
	t.Run("SaveKline and GetKline", func(t *testing.T) {
		now := time.Now().Round(time.Millisecond)
		kline := &model.Kline{
			Symbol:    "BTCUSDT",
			Exchange:  "test-exchange",
			Interval:  model.KlineInterval("1h"),
			OpenTime:  now,
			CloseTime: now.Add(1 * time.Hour),
			Open:      49000.0,
			Close:     50000.0,
			Complete:  true,
		}

		// Save the kline
		err := canonicalRepo.SaveKline(ctx, kline)
		require.NoError(t, err)

		// Retrieve the kline
		retrievedKline, err := canonicalRepo.GetKline(ctx, "BTCUSDT", "test-exchange", model.KlineInterval("1h"), now)
		require.NoError(t, err)

		// Verify the kline was saved correctly
		assert.Equal(t, kline.Symbol, retrievedKline.Symbol)
		assert.Equal(t, kline.Exchange, retrievedKline.Exchange)
		assert.Equal(t, kline.Interval, retrievedKline.Interval)
		assert.Equal(t, kline.OpenTime.Unix(), retrievedKline.OpenTime.Unix())
		assert.Equal(t, kline.CloseTime.Unix(), retrievedKline.CloseTime.Unix())
		assert.Equal(t, kline.Open, retrievedKline.Open)
		assert.Equal(t, kline.Close, retrievedKline.Close)
		assert.Equal(t, kline.Complete, retrievedKline.Complete)
	})

	// Test Create and GetBySymbol
	t.Run("Create and GetBySymbol", func(t *testing.T) {
		symbol := &model.Symbol{
			Symbol:            "ETHUSDT",
			BaseAsset:         "ETH",
			QuoteAsset:        "USDT",
			Exchange:          "test-exchange",
			Status:            model.SymbolStatus("TRADING"),
			MinPrice:          0.01,
			MaxPrice:          100000.0,
			PricePrecision:    2,
			MinQuantity:       0.0001,
			MaxQuantity:       1000.0,
			QuantityPrecision: 4,
			AllowedOrderTypes: []string{"LIMIT", "MARKET"},
		}

		// Save the symbol
		err := canonicalRepo.Create(ctx, symbol)
		require.NoError(t, err)

		// Retrieve the symbol
		retrievedSymbol, err := canonicalRepo.GetBySymbol(ctx, "ETHUSDT")
		require.NoError(t, err)

		// Verify the symbol was saved correctly
		assert.Equal(t, symbol.Symbol, retrievedSymbol.Symbol)
		assert.Equal(t, symbol.BaseAsset, retrievedSymbol.BaseAsset)
		assert.Equal(t, symbol.QuoteAsset, retrievedSymbol.QuoteAsset)
		assert.Equal(t, symbol.Exchange, retrievedSymbol.Exchange)
		assert.Equal(t, symbol.Status, retrievedSymbol.Status)
		assert.Equal(t, symbol.MinPrice, retrievedSymbol.MinPrice)
		assert.Equal(t, symbol.MaxPrice, retrievedSymbol.MaxPrice)
		assert.Equal(t, symbol.PricePrecision, retrievedSymbol.PricePrecision)
		assert.Equal(t, symbol.MinQuantity, retrievedSymbol.MinQuantity)
		assert.Equal(t, symbol.MaxQuantity, retrievedSymbol.MaxQuantity)
		assert.Equal(t, symbol.QuantityPrecision, retrievedSymbol.QuantityPrecision)
		assert.ElementsMatch(t, symbol.AllowedOrderTypes, retrievedSymbol.AllowedOrderTypes)
	})
}
