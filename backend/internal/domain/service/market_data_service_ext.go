package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/cache/standard"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MarketDataServiceWithErrorHandling extends MarketDataService with error handling capabilities
type MarketDataServiceWithErrorHandling struct {
	marketRepo  port.MarketRepository
	symbolRepo  port.SymbolRepository
	cache       port.ExtendedMarketCache
	mexcClient  port.MEXCClient
	baseService *MarketDataService
	logger      *zerolog.Logger
}

// NewMarketDataServiceWithErrorHandling creates a new MarketDataServiceWithErrorHandling
func NewMarketDataServiceWithErrorHandling(
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
	cache port.ExtendedMarketCache,
	mexcClient port.MEXCClient,
	logger *zerolog.Logger,
) *MarketDataServiceWithErrorHandling {
	return &MarketDataServiceWithErrorHandling{
		marketRepo: marketRepo,
		symbolRepo: symbolRepo,
		cache:      cache,
		mexcClient: mexcClient,
		logger:     logger,
	}
}

// NewMarketDataServiceWithErrorHandlingWithService creates a new MarketDataServiceWithErrorHandling using a base service
func NewMarketDataServiceWithErrorHandlingWithService(
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
	cache port.ExtendedMarketCache,
	baseService *MarketDataService,
	mexcClient port.MEXCClient,
	logger *zerolog.Logger,
) *MarketDataServiceWithErrorHandling {
	return &MarketDataServiceWithErrorHandling{
		marketRepo:  marketRepo,
		symbolRepo:  symbolRepo,
		cache:       cache,
		baseService: baseService,
		mexcClient:  mexcClient,
		logger:      logger,
	}
}

// GetTickerWithErrorHandling gets a ticker with error handling
func (s *MarketDataServiceWithErrorHandling) GetTickerWithErrorHandling(ctx context.Context, exchange, symbol string) (*market.Ticker, error) {
	// First try to get from cache
	ticker, err := s.cache.GetTickerWithError(ctx, exchange, symbol)
	if err == nil {
		s.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Msg("Retrieved ticker from cache")
		return ticker, nil
	}

	// Convert cache error for better understanding of what happened
	cacheErr := standard.ConvertCacheError(err)

	if apperror.Is(cacheErr, apperror.ErrNotFound) {
		// Cache miss, fetch using the base service
		s.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Msg("Ticker not found in cache, fetching from API")

		// Use the base service if available, otherwise directly call MEXC
		var newTicker *market.Ticker
		var apiErr error

		if s.baseService != nil {
			// Use the base service - it already handles API interaction
			newTicker, apiErr = s.baseService.RefreshTicker(ctx, symbol)
		} else if s.mexcClient != nil && exchange == "mexc" {
			// Use the MEXC client directly if base service is not available
			apiRes, err := s.mexcClient.GetMarketData(ctx, symbol)
			if err != nil {
				return nil, apperror.NewExternalService("MEXC API", "Failed to get ticker", err)
			}

			// Convert API response to ticker model
			newTicker = &market.Ticker{
				Symbol:        symbol,
				Exchange:      exchange,
				Price:         apiRes.LastPrice,
				Volume:        apiRes.Volume,
				High24h:       apiRes.HighPrice,
				Low24h:        apiRes.LowPrice,
				PriceChange:   apiRes.PriceChange,
				PercentChange: apiRes.PriceChangePercent,
				LastUpdated:   time.Now(),
			}
		} else {
			return nil, apperror.NewNotFound("exchange", exchange, fmt.Errorf("unsupported exchange"))
		}

		if apiErr != nil {
			return nil, apiErr
		}

		// Cache the new ticker for future use
		s.cache.CacheTicker(newTicker)

		// Also persist to database for historical records
		err = s.marketRepo.SaveTicker(ctx, newTicker)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Str("exchange", exchange).
				Str("symbol", symbol).
				Msg("Failed to save ticker to database")
			// Continue despite DB error since we have the data
		}

		return newTicker, nil
	}

	// For other errors, return as is
	return nil, cacheErr
}

// GetLatestCandleWithErrorHandling gets the latest candle with error handling
func (s *MarketDataServiceWithErrorHandling) GetLatestCandleWithErrorHandling(
	ctx context.Context,
	exchange string,
	symbol string,
	interval market.Interval,
) (*market.Candle, error) {
	// First try to get from cache
	candle, err := s.cache.GetLatestCandleWithError(ctx, exchange, symbol, interval)
	if err == nil {
		s.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Str("interval", string(interval)).
			Msg("Retrieved latest candle from cache")
		return candle, nil
	}

	// Convert cache error for better understanding of what happened
	cacheErr := standard.ConvertCacheError(err)

	if apperror.Is(cacheErr, apperror.ErrNotFound) {
		// Cache miss, fetch using base service or API
		s.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Str("interval", string(interval)).
			Msg("Latest candle not found in cache, fetching from API")

		if s.baseService != nil {
			// Use the base service's RefreshCandles, which already has error handling
			candles, err := s.baseService.RefreshCandles(ctx, symbol, interval, 1)
			if err != nil {
				return nil, err
			}

			if len(candles) == 0 {
				return nil, apperror.NewNotFound("candle",
					fmt.Sprintf("%s:%s:%s", exchange, symbol, interval), nil)
			}

			// Get first candle and return pointer
			latestCandle := &candles[0]
			return latestCandle, nil
		} else if s.mexcClient != nil && exchange == "mexc" {
			// Fall back to direct API usage if needed
			// Convert market.Interval to model.KlineInterval
			modelInterval := model.KlineInterval(interval)

			klines, err := s.mexcClient.GetKlines(ctx, symbol, modelInterval, 1)
			if err != nil {
				return nil, apperror.NewExternalService("MEXC API", "Failed to get candles", err)
			}

			if len(klines) == 0 {
				return nil, apperror.NewNotFound("candle",
					fmt.Sprintf("%s:%s:%s", exchange, symbol, interval), nil)
			}

			// Convert model.Kline to market.Candle
			kline := klines[0]
			latestCandle := &market.Candle{
				Symbol:    symbol,
				Interval:  interval,
				OpenTime:  kline.OpenTime,
				Open:      kline.Open,
				High:      kline.High,
				Low:       kline.Low,
				Close:     kline.Close,
				Volume:    kline.Volume,
				CloseTime: kline.CloseTime,
				Exchange:  exchange,
			}

			// Cache for future use
			s.cache.CacheCandle(latestCandle)

			return latestCandle, nil
		} else {
			return nil, apperror.NewNotFound("exchange", exchange, fmt.Errorf("unsupported exchange"))
		}
	}

	// For other errors, return as is
	return nil, cacheErr
}

// GetOrderBookWithErrorHandling gets an order book with error handling
func (s *MarketDataServiceWithErrorHandling) GetOrderBookWithErrorHandling(
	ctx context.Context,
	exchange string,
	symbol string,
) (*market.OrderBook, error) {
	// First try to get from cache
	orderBook, err := s.cache.GetOrderBookWithError(ctx, exchange, symbol)
	if err == nil {
		s.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Msg("Retrieved order book from cache")
		return orderBook, nil
	}

	// Convert cache error for better understanding of what happened
	cacheErr := standard.ConvertCacheError(err)

	if apperror.Is(cacheErr, apperror.ErrNotFound) {
		// Cache miss, fetch using base service or API
		s.logger.Debug().
			Str("exchange", exchange).
			Str("symbol", symbol).
			Msg("Order book not found in cache, fetching from API")

		if s.baseService != nil {
			// Use the base service's GetOrderBook, which already has error handling
			orderBook, err := s.baseService.GetOrderBook(ctx, symbol, 10)
			if err != nil {
				return nil, err
			}

			return orderBook, nil
		} else if s.mexcClient != nil && exchange == "mexc" {
			modelOrderBook, err := s.mexcClient.GetOrderBook(ctx, symbol, 10) // Default depth 10
			if err != nil {
				return nil, apperror.NewExternalService("MEXC API", "Failed to get order book", err)
			}

			// Convert model.OrderBook to market.OrderBook
			marketOrderBook := &market.OrderBook{
				Symbol:       modelOrderBook.Symbol,
				LastUpdateID: modelOrderBook.LastUpdateID,
				Exchange:     exchange,
				LastUpdated:  time.Now(),
			}

			// Convert bids and asks
			marketOrderBook.Bids = make([]market.OrderBookEntry, len(modelOrderBook.Bids))
			for i, bid := range modelOrderBook.Bids {
				marketOrderBook.Bids[i] = market.OrderBookEntry{
					Price:    bid.Price,
					Quantity: bid.Quantity,
				}
			}

			marketOrderBook.Asks = make([]market.OrderBookEntry, len(modelOrderBook.Asks))
			for i, ask := range modelOrderBook.Asks {
				marketOrderBook.Asks[i] = market.OrderBookEntry{
					Price:    ask.Price,
					Quantity: ask.Quantity,
				}
			}

			// Cache for future use
			s.cache.CacheOrderBook(marketOrderBook)

			return marketOrderBook, nil
		} else {
			return nil, apperror.NewNotFound("exchange", exchange, fmt.Errorf("unsupported exchange"))
		}
	}

	// For other errors, return as is
	return nil, cacheErr
}

// Implement port.MarketDataService interface
func (s *MarketDataServiceWithErrorHandling) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	if s.baseService != nil {
		return s.baseService.GetTicker(ctx, symbol)
	}
	return nil, fmt.Errorf("baseService not available")
}

func (s *MarketDataServiceWithErrorHandling) GetCandles(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error) {
	if s.baseService != nil {
		return s.baseService.GetCandles(ctx, symbol, interval, limit)
	}
	return nil, fmt.Errorf("baseService not available")
}

func (s *MarketDataServiceWithErrorHandling) GetOrderBook(ctx context.Context, symbol string, depth int) (*market.OrderBook, error) {
	if s.baseService != nil {
		return s.baseService.GetOrderBook(ctx, symbol, depth)
	}
	return nil, fmt.Errorf("baseService not available")
}

func (s *MarketDataServiceWithErrorHandling) GetAllSymbols(ctx context.Context) ([]*market.Symbol, error) {
	if s.baseService != nil {
		return s.baseService.GetAllSymbols(ctx)
	}
	return nil, fmt.Errorf("baseService not available")
}

func (s *MarketDataServiceWithErrorHandling) GetSymbolInfo(ctx context.Context, symbol string) (*market.Symbol, error) {
	if s.baseService != nil {
		return s.baseService.GetSymbolInfo(ctx, symbol)
	}
	return nil, fmt.Errorf("baseService not available")
}

func (s *MarketDataServiceWithErrorHandling) GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*market.Candle, error) {
	if s.baseService != nil {
		return s.baseService.GetHistoricalPrices(ctx, symbol, from, to, interval)
	}
	return nil, fmt.Errorf("baseService not available")
}
