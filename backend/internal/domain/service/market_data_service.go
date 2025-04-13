package service

import (
	"context"
	"sync"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
)

// MarketDataService provides higher-level market data operations
// and coordinates between gateway and data storage
type MarketDataService struct {
	marketRepo  port.MarketRepository
	symbolRepo  port.SymbolRepository
	cache       port.MarketCache
	mexcAPI     port.MexcAPI
	logger      *zerolog.Logger
	refreshLock sync.Mutex
}

// NewMarketDataService creates a new MarketDataService
func NewMarketDataService(
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
	cache port.MarketCache,
	mexcAPI port.MexcAPI,
	logger *zerolog.Logger,
) *MarketDataService {
	return &MarketDataService{
		marketRepo:  marketRepo,
		symbolRepo:  symbolRepo,
		cache:       cache,
		mexcAPI:     mexcAPI,
		logger:      logger,
		refreshLock: sync.Mutex{},
	}
}

// RefreshSymbols fetches all trading symbols from the exchange, updates the database
// and returns the updated list
func (s *MarketDataService) RefreshSymbols(ctx context.Context) ([]market.Symbol, error) {
	// Make sure only one refresh is running at a time
	s.refreshLock.Lock()
	defer s.refreshLock.Unlock()

	s.logger.Debug().Msg("Refreshing symbols from exchange")

	// For now, we'll just get symbols from the database
	// In a real implementation, we would fetch from the exchange API
	dbSymbols, dbErr := s.symbolRepo.GetAll(ctx)
	if dbErr != nil {
		s.logger.Error().Err(dbErr).Msg("Failed to fetch symbols from database")
		return nil, dbErr
	}

	// Convert pointer slice to value slice
	symbols := make([]market.Symbol, len(dbSymbols))
	for i, symbol := range dbSymbols {
		symbols[i] = *symbol
	}

	return symbols, nil
}

// RefreshTicker fetches the latest ticker data for a specific symbol from the exchange,
// updates the database, and returns the updated ticker
func (s *MarketDataService) RefreshTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	s.logger.Debug().Str("symbol", symbol).Msg("Refreshing ticker from exchange")

	// Try to get from cache first
	cachedTicker, exists := s.cache.GetTicker(ctx, "mexc", symbol)
	if exists {
		return cachedTicker, nil
	}

	// Try to get from MEXC API using GetMarketData
	ticker, err := s.mexcAPI.GetMarketData(ctx, symbol)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch ticker from exchange API")

		// Fall back to database
		dbTicker, dbErr := s.marketRepo.GetTicker(ctx, symbol, "mexc")
		if dbErr != nil {
			s.logger.Error().Err(dbErr).Str("symbol", symbol).Msg("Failed to fetch ticker from database")
			return nil, dbErr
		}

		return dbTicker, nil
	}

	// Convert model.Ticker to market.Ticker
	marketTicker := &market.Ticker{
		Symbol:      ticker.Symbol,
		Price:       ticker.LastPrice,
		Volume:      ticker.Volume,
		High24h:     ticker.HighPrice,
		Low24h:      ticker.LowPrice,
		LastUpdated: time.Now(),
		Exchange:    "mexc",
	}

	// Update cache
	s.cache.CacheTicker(marketTicker)

	// Update database
	err = s.marketRepo.SaveTicker(ctx, marketTicker)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to save ticker to database")
	}

	return marketTicker, nil
}

// RefreshCandles fetches the latest candle data for a specific symbol and interval,
// updates the database, and returns the updated candles
func (s *MarketDataService) RefreshCandles(
	ctx context.Context,
	symbol string,
	interval market.Interval,
	limit int,
) ([]market.Candle, error) {
	s.logger.Debug().
		Str("symbol", symbol).
		Str("interval", string(interval)).
		Int("limit", limit).
		Msg("Refreshing candles from exchange")

	// Try to get candles from MEXC API using GetKlines
	endTime := time.Now()
	startTime := endTime.Add(-time.Hour * 24) // Default to last 24 hours

	// Convert market.Interval to model.KlineInterval
	modelInterval := model.KlineInterval(interval)

	klines, err := s.mexcAPI.GetKlines(ctx, symbol, modelInterval, limit)
	if err != nil {
		s.logger.Error().Err(err).
			Str("symbol", symbol).
			Str("interval", string(interval)).
			Msg("Failed to fetch candles from exchange API")

		// Fall back to database
		dbCandles, dbErr := s.marketRepo.GetCandles(ctx, symbol, "mexc", interval, startTime, endTime, limit)
		if dbErr != nil {
			s.logger.Error().Err(dbErr).
				Str("symbol", symbol).
				Str("interval", string(interval)).
				Msg("Failed to fetch candles from database")
			return nil, dbErr
		}

		// Convert pointer slice to value slice
		candles := make([]market.Candle, len(dbCandles))
		for i, candle := range dbCandles {
			candles[i] = *candle
		}

		return candles, nil
	}

	// Convert model.Kline to market.Candle
	candles := make([]market.Candle, len(klines))
	for i, kline := range klines {
		candles[i] = market.Candle{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  kline.OpenTime,
			Open:      kline.Open,
			High:      kline.High,
			Low:       kline.Low,
			Close:     kline.Close,
			Volume:    kline.Volume,
			CloseTime: kline.CloseTime,
			Exchange:  "mexc",
		}
	}

	// Update cache and database
	for i := range candles {
		s.cache.CacheCandle(&candles[i])

		err = s.marketRepo.SaveCandle(ctx, &candles[i])
		if err != nil {
			s.logger.Error().Err(err).
				Str("symbol", symbol).
				Str("interval", string(interval)).
				Time("openTime", candles[i].OpenTime).
				Msg("Failed to save candle to database")
		}
	}

	return candles, nil
}

// GetHistoricalPrices fetches historical price data for a specific symbol
func (s *MarketDataService) GetHistoricalPrices(
	ctx context.Context,
	symbol string,
	startTime, endTime time.Time,
) ([]market.Ticker, error) {
	s.logger.Debug().
		Str("symbol", symbol).
		Time("startTime", startTime).
		Time("endTime", endTime).
		Msg("Getting historical prices")

	// Fetch from database
	tickers, err := s.marketRepo.GetTickerHistory(ctx, symbol, "mexc", startTime, endTime)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Msg("Failed to fetch historical ticker data")
		return nil, err
	}

	// Convert pointer slice to value slice
	result := make([]market.Ticker, len(tickers))
	for i, ticker := range tickers {
		result[i] = *ticker
	}

	return result, nil
}
