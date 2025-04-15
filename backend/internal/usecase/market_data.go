package usecase

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MarketDataUseCase implements use cases for market data
type MarketDataUseCase struct {
	marketRepo port.MarketRepository
	symbolRepo port.SymbolRepository
	cache      port.MarketCache
	logger     *zerolog.Logger
}

// NewMarketDataUseCase creates a new MarketDataUseCase
func NewMarketDataUseCase(
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
	cache port.MarketCache,
	logger *zerolog.Logger,
) *MarketDataUseCase {
	return &MarketDataUseCase{
		marketRepo: marketRepo,
		symbolRepo: symbolRepo,
		cache:      cache,
		logger:     logger,
	}
}

// GetLatestTickers returns the latest tickers from cache or database
func (uc *MarketDataUseCase) GetLatestTickers(ctx context.Context) ([]market.Ticker, error) {
	// Try to get from cache first
	tickers, exists := uc.cache.GetLatestTickers(ctx)
	if exists && len(tickers) > 0 {
		uc.logger.Debug().Int("count", len(tickers)).Msg("Retrieved tickers from cache")

		// Convert from pointer slice to value slice
		result := make([]market.Ticker, len(tickers))
		for i, ticker := range tickers {
			result[i] = *ticker
		}
		return result, nil
	}

	// Fall back to database - using default exchange for now
	const defaultExchange = "mexc"
	tickers, err := uc.marketRepo.GetAllTickers(ctx, defaultExchange)
	if err != nil {
		return nil, err
	}

	// Store in cache for next time
	for _, ticker := range tickers {
		tickerCopy := *ticker // Create a copy to avoid storing references that might change
		uc.cache.CacheTicker(&tickerCopy)
	}

	// Convert from pointer slice to value slice
	result := make([]market.Ticker, len(tickers))
	for i, ticker := range tickers {
		result[i] = *ticker
	}

	uc.logger.Debug().Int("count", len(result)).Msg("Retrieved tickers from database")
	return result, nil
}

// GetTicker returns the latest ticker for a specific symbol from cache or database
func (uc *MarketDataUseCase) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, error) {
	// Try to get from cache first
	ticker, exists := uc.cache.GetTicker(ctx, exchange, symbol)
	if exists && ticker != nil {
		uc.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Msg("Retrieved ticker from cache")
		return ticker, nil
	}

	// Fall back to database
	ticker, err := uc.marketRepo.GetTicker(ctx, symbol, exchange)
	if err != nil {
		return nil, err
	}

	if ticker != nil {
		// Store in cache for next time
		tickerCopy := *ticker // Create a copy to avoid storing references that might change
		uc.cache.CacheTicker(&tickerCopy)
	}

	uc.logger.Debug().
		Str("exchange", exchange).
		Str("symbol", symbol).
		Bool("found", ticker != nil).
		Msg("Retrieved ticker from database")
	return ticker, nil
}

// GetCandles returns historical candles for a specific symbol
func (uc *MarketDataUseCase) GetCandles(
	ctx context.Context,
	exchange string,
	symbol string,
	interval market.Interval,
	startTime time.Time,
	endTime time.Time,
	limit int,
) ([]market.Candle, error) {
	// Try to get from cache if the time range is recent
	// Only use cache for last 24h data and if limit is reasonable
	now := time.Now()
	isRecent := startTime.After(now.Add(-24 * time.Hour))
	useCache := isRecent && limit <= 200

	if useCache {
		// Check each candle individually in cache
		var result []market.Candle
		currentTime := startTime

		for currentTime.Before(endTime) || currentTime.Equal(endTime) {
			candle, exists := uc.cache.GetCandle(ctx, exchange, symbol, interval, currentTime)
			if exists && candle != nil {
				result = append(result, *candle)
			}

			// Move to next interval
			switch interval {
			case market.Interval1m:
				currentTime = currentTime.Add(1 * time.Minute)
			case market.Interval5m:
				currentTime = currentTime.Add(5 * time.Minute)
			case market.Interval15m:
				currentTime = currentTime.Add(15 * time.Minute)
			case market.Interval30m:
				currentTime = currentTime.Add(30 * time.Minute)
			case market.Interval1h:
				currentTime = currentTime.Add(1 * time.Hour)
			case market.Interval4h:
				currentTime = currentTime.Add(4 * time.Hour)
			case market.Interval1d:
				currentTime = currentTime.Add(24 * time.Hour)
			default:
				currentTime = currentTime.Add(1 * time.Hour)
			}
		}

		if len(result) > 0 {
			// Ensure we don't return more than the limit
			if limit > 0 && len(result) > limit {
				result = result[len(result)-limit:]
			}
			uc.logger.Debug().
				Str("exchange", exchange).
				Str("symbol", symbol).
				Str("interval", string(interval)).
				Int("count", len(result)).
				Msg("Retrieved candles from cache")
			return result, nil
		}
	}

	// Fall back to database
	candles, err := uc.marketRepo.GetCandles(ctx, symbol, exchange, interval, startTime, endTime, limit)
	if err != nil {
		return nil, err
	}

	// Store in cache for next time if it's recent data
	if isRecent && len(candles) > 0 {
		for _, candle := range candles {
			candleCopy := *candle // Create a copy to avoid storing references that might change
			uc.cache.CacheCandle(&candleCopy)
		}
	}

	// Convert from pointer slice to value slice
	result := make([]market.Candle, len(candles))
	for i, candle := range candles {
		result[i] = *candle
	}

	uc.logger.Debug().
		Str("exchange", exchange).
		Str("symbol", symbol).
		Str("interval", string(interval)).
		Int("count", len(result)).
		Msg("Retrieved candles from database")
	return result, nil
}

// GetAllSymbols returns all available trading symbols
func (uc *MarketDataUseCase) GetAllSymbols(ctx context.Context) ([]market.Symbol, error) {
	// Try to get from cache first

	// For now, since the Cache interface doesn't have a direct method to get all symbols,
	// we'll just fetch from the database

	// Fall back to database
	symbols, err := uc.symbolRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert from pointer slice to value slice
	result := make([]market.Symbol, len(symbols))
	for i, symbol := range symbols {
		result[i] = *symbol
	}

	uc.logger.Debug().Int("count", len(result)).Msg("Retrieved symbols from database")
	return result, nil
}

// GetSymbolInfo returns detailed information for a specific trading symbol
func (uc *MarketDataUseCase) GetSymbolInfo(ctx context.Context, symbol string) (*market.Symbol, error) {
	// We'll fetch directly from the database since we don't have a specific cache method for symbols
	symbolInfo, err := uc.symbolRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	uc.logger.Debug().
		Str("symbol", symbol).
		Bool("found", symbolInfo != nil).
		Msg("Retrieved symbol info from database")
	return symbolInfo, nil
}
