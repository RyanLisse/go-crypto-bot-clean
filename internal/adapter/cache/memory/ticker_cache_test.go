package memory

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/rs/zerolog"
)

func TestTickerCache_CacheTicker(t *testing.T) {
	// Create a silent logger for testing
	logger := zerolog.New(os.DevNull)
	
	cache := NewTickerCache(&logger)
	
	// Create a test ticker
	ticker := &market.Ticker{
		Exchange:      "binance",
		Symbol:        "BTCUSDT",
		Price:         50000.0,
		Volume:        100.0,
		High:          51000.0,
		Low:           49000.0,
		ChangePercent: 2.5,
		UpdateTime:    time.Now(),
	}
	
	// Cache the ticker
	cache.CacheTicker(ticker)
	
	// Test GetTicker
	ctx := context.Background()
	retrievedTicker, found := cache.GetTicker(ctx, "binance", "BTCUSDT")
	
	if !found {
		t.Error("Expected ticker to be found in cache")
	}
	
	if retrievedTicker == nil {
		t.Fatal("Retrieved ticker should not be nil")
	}
	
	if retrievedTicker.Symbol != ticker.Symbol {
		t.Errorf("Expected ticker symbol %s, got %s", ticker.Symbol, retrievedTicker.Symbol)
	}
	
	if retrievedTicker.Price != ticker.Price {
		t.Errorf("Expected ticker price %f, got %f", ticker.Price, retrievedTicker.Price)
	}
	
	// Test GetAllTickers for an exchange
	allTickers, found := cache.GetAllTickers(ctx, "binance")
	
	if !found {
		t.Error("Expected tickers to be found for exchange")
	}
	
	if len(allTickers) != 1 {
		t.Errorf("Expected 1 ticker for exchange, got %d", len(allTickers))
	}
	
	if allTickers[0].Symbol != ticker.Symbol {
		t.Errorf("Expected ticker symbol %s, got %s", ticker.Symbol, allTickers[0].Symbol)
	}
	
	// Test GetLatestTickers
	latestTickers, found := cache.GetLatestTickers(ctx)
	
	if !found {
		t.Error("Expected latest tickers to be found")
	}
	
	if len(latestTickers) != 1 {
		t.Errorf("Expected 1 latest ticker, got %d", len(latestTickers))
	}
	
	if latestTickers[0].Symbol != ticker.Symbol {
		t.Errorf("Expected ticker symbol %s, got %s", ticker.Symbol, latestTickers[0].Symbol)
	}
}

func TestTickerCache_Clear(t *testing.T) {
	// Create a silent logger for testing
	logger := zerolog.New(os.DevNull)
	
	cache := NewTickerCache(&logger)
	
	// Create a test ticker
	ticker := &market.Ticker{
		Exchange:      "binance",
		Symbol:        "BTCUSDT",
		Price:         50000.0,
		UpdateTime:    time.Now(),
	}
	
	// Cache the ticker
	cache.CacheTicker(ticker)
	
	// Verify it's in the cache
	ctx := context.Background()
	_, found := cache.GetTicker(ctx, "binance", "BTCUSDT")
	
	if !found {
		t.Error("Expected ticker to be in cache before clearing")
	}
	
	// Clear the cache
	cache.Clear()
	
	// Verify it's no longer in the cache
	_, found = cache.GetTicker(ctx, "binance", "BTCUSDT")
	
	if found {
		t.Error("Expected ticker to be removed after clearing cache")
	}
	
	// Also check GetAllTickers and GetLatestTickers
	_, found = cache.GetAllTickers(ctx, "binance")
	if found {
		t.Error("Expected no tickers for exchange after clearing cache")
	}
	
	_, found = cache.GetLatestTickers(ctx)
	if found {
		t.Error("Expected no latest tickers after clearing cache")
	}
} 