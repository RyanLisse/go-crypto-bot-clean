package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MarketDataService provides higher-level market data operations
// and coordinates between gateway and data storage
type MarketDataService struct {
	marketRepo  port.MarketRepository
	symbolRepo  port.SymbolRepository
	cache       port.MarketCache
	mexcClient  port.MEXCClient // Changed mexcAPI to mexcClient
	logger      *zerolog.Logger
	refreshLock sync.Mutex
}

// NewMarketDataService creates a new MarketDataService
func NewMarketDataService(
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
	cache port.MarketCache,
	mexcClient port.MEXCClient, // Changed mexcAPI to mexcClient
	logger *zerolog.Logger,
) *MarketDataService {
	return &MarketDataService{
		marketRepo:  marketRepo,
		symbolRepo:  symbolRepo,
		cache:       cache,
		mexcClient:  mexcClient, // Changed mexcAPI to mexcClient
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

	// Try to get from MEXC Client using GetMarketData
	ticker, err := s.mexcClient.GetMarketData(ctx, symbol) // Changed mexcAPI to mexcClient
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

	klines, err := s.mexcClient.GetKlines(ctx, symbol, modelInterval, limit) // Changed mexcAPI to mexcClient
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

// GetHistoricalTickerPrices fetches historical ticker price data for a specific symbol
func (s *MarketDataService) GetHistoricalTickerPrices(
	ctx context.Context,
	symbol string,
	startTime, endTime time.Time,
) ([]market.Ticker, error) {
	s.logger.Debug().
		Str("symbol", symbol).
		Time("startTime", startTime).
		Time("endTime", endTime).
		Msg("Getting historical ticker prices")

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

// GetTicker implements the port.MarketDataService interface
func (s *MarketDataService) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	return s.RefreshTicker(ctx, symbol)
}

// GetCandles implements the port.MarketDataService interface
func (s *MarketDataService) GetCandles(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error) {
	marketInterval := market.Interval(interval)
	candles, err := s.RefreshCandles(ctx, symbol, marketInterval, limit)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	result := make([]*market.Candle, len(candles))
	for i := range candles {
		result[i] = &candles[i]
	}

	return result, nil
}

// GetOrderBook implements the port.MarketDataService interface
func (s *MarketDataService) GetOrderBook(ctx context.Context, symbol string, depth int) (*market.OrderBook, error) {
	s.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Getting order book")

	// Try to get from cache first
	cachedOrderBook, exists := s.cache.GetOrderBook(ctx, "mexc", symbol)
	if exists {
		return cachedOrderBook, nil
	}

	// Try to get from MEXC Client
	if s.mexcClient != nil { // Changed mexcAPI to mexcClient
		orderBook, err := s.mexcClient.GetOrderBook(ctx, symbol, depth) // Changed mexcAPI to mexcClient
		if err == nil {
			// Convert the model to market model if needed
			marketOrderBook := &market.OrderBook{
				Symbol:      orderBook.Symbol,
				Bids:        make([]market.OrderBookEntry, len(orderBook.Bids)),
				Asks:        make([]market.OrderBookEntry, len(orderBook.Asks)),
				LastUpdated: time.Now(), // Changed Timestamp to LastUpdated
				Exchange:    "mexc",
			}

			for i, bid := range orderBook.Bids {
				marketOrderBook.Bids[i] = market.OrderBookEntry{
					Price:    bid.Price,
					Quantity: bid.Quantity,
				}
			}

			for i, ask := range orderBook.Asks {
				marketOrderBook.Asks[i] = market.OrderBookEntry{
					Price:    ask.Price,
					Quantity: ask.Quantity,
				}
			}

			// Cache the result
			s.cache.CacheOrderBook(marketOrderBook)

			return marketOrderBook, nil
		}

		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch order book from exchange API")
	}

	// Fall back to database
	dbOrderBook, err := s.marketRepo.GetOrderBook(ctx, symbol, "mexc", depth)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch order book from database")
		return nil, err
	}

	return dbOrderBook, nil
}

// GetAllSymbols implements the port.MarketDataService interface
func (s *MarketDataService) GetAllSymbols(ctx context.Context) ([]*market.Symbol, error) {
	s.logger.Debug().Msg("Getting all symbols")

	// Try to get the symbols from the database
	symbols, err := s.symbolRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to fetch symbols from database")

		// If database fetch fails and we have MEXC Client available, try to get from there
		if s.mexcClient != nil { // Changed mexcAPI to mexcClient
			exchangeSymbols, err := s.mexcClient.GetExchangeInfo(ctx) // Changed mexcAPI to mexcClient
			if err == nil {
				// Convert symbols and return
				result := make([]*market.Symbol, len(exchangeSymbols.Symbols))
				for i, sym := range exchangeSymbols.Symbols {
					result[i] = &market.Symbol{
						Symbol:              sym.Symbol,
						BaseAsset:           sym.BaseAsset,
						QuoteAsset:          sym.QuoteAsset,
						Status:              sym.Status,
						BaseAssetPrecision:  sym.BaseAssetPrecision,
						QuoteAssetPrecision: sym.QuoteAssetPrecision,
						MinNotional:         parseStringToFloat64(sym.MinNotional),
						MinLotSize:          parseStringToFloat64(sym.MinLotSize),
						MaxLotSize:          parseStringToFloat64(sym.MaxLotSize),
						StepSize:            parseStringToFloat64(sym.StepSize),
						TickSize:            parseStringToFloat64(sym.TickSize),
					}
				}
				return result, nil
			}
			s.logger.Error().Err(err).Msg("Failed to fetch symbols from exchange API")
		}

		return nil, err
	}

	return symbols, nil
}

// GetSymbolInfo implements the port.MarketDataService interface
func (s *MarketDataService) GetSymbolInfo(ctx context.Context, symbol string) (*market.Symbol, error) {
	s.logger.Debug().Str("symbol", symbol).Msg("Getting symbol info")

	// Get symbol info from the repository
	symbolInfo, err := s.symbolRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch symbol info from database")

		// If database fetch fails and we have MEXC Client available, try to get from there
		if s.mexcClient != nil { // Changed mexcAPI to mexcClient
			exchangeSymbol, err := s.mexcClient.GetSymbolInfo(ctx, symbol) // Changed mexcAPI to mexcClient
			if err == nil {
				// Convert and return
				return &market.Symbol{
					Symbol:              exchangeSymbol.Symbol,
					BaseAsset:           exchangeSymbol.BaseAsset,
					QuoteAsset:          exchangeSymbol.QuoteAsset,
					Status:              exchangeSymbol.Status,
					BaseAssetPrecision:  exchangeSymbol.BaseAssetPrecision,
					QuoteAssetPrecision: exchangeSymbol.QuoteAssetPrecision,
					MinNotional:         parseStringToFloat64(exchangeSymbol.MinNotional),
					MinLotSize:          parseStringToFloat64(exchangeSymbol.MinLotSize),
					MaxLotSize:          parseStringToFloat64(exchangeSymbol.MaxLotSize),
					StepSize:            parseStringToFloat64(exchangeSymbol.StepSize),
					TickSize:            parseStringToFloat64(exchangeSymbol.TickSize),
				}, nil
			}
			s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to fetch symbol info from exchange API")
		}

		return nil, err
	}

	return symbolInfo, nil
}

// GetHistoricalPrices implements the port.MarketDataService interface
func (s *MarketDataService) GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*market.Candle, error) {
	marketInterval := market.Interval(interval)
	candles, err := s.marketRepo.GetCandles(ctx, symbol, "mexc", marketInterval, from, to, 0)
	if err != nil {
		return nil, err
	}
	return candles, nil
}

// Helper to parse string to float64, returns 0 if parsing fails
func parseStringToFloat64(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
