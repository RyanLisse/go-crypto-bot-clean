package service

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MarketDataServiceAdapter adapts port.MarketDataService to MarketDataServiceInterface
type MarketDataServiceAdapter struct {
	marketDataService port.MarketDataService
	logger            *zerolog.Logger
}

// NewMarketDataServiceAdapter creates a new adapter
func NewMarketDataServiceAdapter(marketDataService port.MarketDataService, logger *zerolog.Logger) MarketDataServiceInterface {
	return &MarketDataServiceAdapter{
		marketDataService: marketDataService,
		logger:            logger,
	}
}

// RefreshTicker implements MarketDataServiceInterface
func (a *MarketDataServiceAdapter) RefreshTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	// Get the latest ticker for the symbol
	ticker, err := a.marketDataService.GetTicker(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

// GetHistoricalPrices implements MarketDataServiceInterface
func (a *MarketDataServiceAdapter) GetHistoricalPrices(ctx context.Context, symbol string, startTime, endTime time.Time) ([]market.Ticker, error) {
	// Convert candles to tickers
	candles, err := a.marketDataService.GetHistoricalPrices(ctx, symbol, startTime, endTime, string(market.Interval1h))
	if err != nil {
		return nil, err
	}

	// Convert candles to tickers
	tickers := make([]market.Ticker, 0, len(candles))
	for _, candle := range candles {
		tickers = append(tickers, market.Ticker{
			Symbol:      candle.Symbol,
			Price:       candle.Close,
			Volume:      candle.Volume,
			LastUpdated: candle.CloseTime,
		})
	}

	return tickers, nil
}
