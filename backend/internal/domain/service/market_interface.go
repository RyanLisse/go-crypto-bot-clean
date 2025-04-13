package service

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model/market"
)

// MarketDataServiceInterface defines the interface for market data services
type MarketDataServiceInterface interface {
	// RefreshTicker fetches the latest ticker data for a specific symbol
	RefreshTicker(ctx context.Context, symbol string) (*market.Ticker, error)
	
	// GetHistoricalPrices fetches historical price data for a specific symbol
	GetHistoricalPrices(ctx context.Context, symbol string, startTime, endTime time.Time) ([]market.Ticker, error)
}
