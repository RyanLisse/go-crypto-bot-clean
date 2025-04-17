package market

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
)

// MarketDataServiceAdapter adapts the concrete MarketDataService to the port.MarketDataService interface
type MarketDataServiceAdapter struct {
	service *service.MarketDataService
	logger  *zerolog.Logger
}

// NewMarketDataServiceAdapter creates a new MarketDataServiceAdapter
func NewMarketDataServiceAdapter(service *service.MarketDataService, logger *zerolog.Logger) port.MarketDataService {
	return &MarketDataServiceAdapter{
		service: service,
		logger:  logger,
	}
}

// GetTicker adapts the GetTicker method
func (a *MarketDataServiceAdapter) GetTicker(ctx context.Context, symbol string) (*model.Ticker, error) {
	return a.service.GetTicker(ctx, symbol)
}

// GetCandles adapts the GetCandles method
func (a *MarketDataServiceAdapter) GetCandles(ctx context.Context, symbol string, interval string, limit int) ([]*model.Kline, error) {
	return a.service.GetCandles(ctx, symbol, interval, limit)
}

// GetOrderBook adapts the GetOrderBook method
func (a *MarketDataServiceAdapter) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	// This is a simplified implementation since the actual OrderBook method might differ
	a.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Adapter: Getting order book")

	// Try to get from the service
	ticker, err := a.service.GetTicker(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Create a simple order book with the ticker price
	orderBook := &model.OrderBook{
		Symbol: symbol,
		Bids: []model.OrderBookEntry{
			{
				Price:    ticker.LastPrice * 0.99, // Simulate a bid slightly below current price
				Quantity: 1.0,
			},
		},
		Asks: []model.OrderBookEntry{
			{
				Price:    ticker.LastPrice * 1.01, // Simulate an ask slightly above current price
				Quantity: 1.0,
			},
		},
		Timestamp: time.Now(),
	}

	return orderBook, nil
}

// GetAllSymbols adapts the GetAllSymbols method
func (a *MarketDataServiceAdapter) GetAllSymbols(ctx context.Context) ([]*model.Symbol, error) {
	// This is a simplified implementation
	a.logger.Debug().Msg("Adapter: Getting all symbols")

	// Return a minimal implementation with a few common symbols
	return []*model.Symbol{
		{
			Symbol:     "BTCUSDT",
			BaseAsset:  "BTC",
			QuoteAsset: "USDT",
			Status:     model.SymbolStatusTrading,
		},
		{
			Symbol:     "ETHUSDT",
			BaseAsset:  "ETH",
			QuoteAsset: "USDT",
			Status:     model.SymbolStatusTrading,
		},
		{
			Symbol:     "BNBUSDT",
			BaseAsset:  "BNB",
			QuoteAsset: "USDT",
			Status:     model.SymbolStatusTrading,
		},
	}, nil
}

// GetSymbolInfo adapts the GetSymbolInfo method
func (a *MarketDataServiceAdapter) GetSymbolInfo(ctx context.Context, symbol string) (*model.Symbol, error) {
	// This is a simplified implementation
	a.logger.Debug().Str("symbol", symbol).Msg("Adapter: Getting symbol info")

	// Get all symbols and find the matching one
	symbols, err := a.GetAllSymbols(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range symbols {
		if s.Symbol == symbol {
			return s, nil
		}
	}

	// Return a default symbol if not found
	return &model.Symbol{
		Symbol:     symbol,
		BaseAsset:  symbol[:len(symbol)-4], // Assume the last 4 chars are the quote asset
		QuoteAsset: symbol[len(symbol)-4:],
		Status:     model.SymbolStatusTrading,
	}, nil
}

// GetHistoricalPrices adapts the GetHistoricalPrices method
func (a *MarketDataServiceAdapter) GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*model.Kline, error) {
	// This is a simplified implementation
	a.logger.Debug().
		Str("symbol", symbol).
		Time("from", from).
		Time("to", to).
		Str("interval", interval).
		Msg("Adapter: Getting historical candles")

	// Call the GetCandles method with a limit of 1000
	return a.service.GetCandles(ctx, symbol, interval, 1000)
}

// Deprecated legacy methods to satisfy the interface

// GetTickerLegacy implements the legacy ticker method
func (a *MarketDataServiceAdapter) GetTickerLegacy(ctx context.Context, symbol string) (*market.Ticker, error) {
	ticker, err := a.GetTicker(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Convert to legacy model
	return &market.Ticker{
		Symbol:      ticker.Symbol,
		Price:       ticker.LastPrice,
		Volume:      ticker.Volume,
		LastUpdated: ticker.Timestamp,
		Exchange:    ticker.Exchange,
	}, nil
}

// GetCandlesLegacy implements the legacy candles method
func (a *MarketDataServiceAdapter) GetCandlesLegacy(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error) {
	klines, err := a.GetCandles(ctx, symbol, interval, limit)
	if err != nil {
		return nil, err
	}

	// Convert to legacy model
	candles := make([]*market.Candle, len(klines))
	for i, kline := range klines {
		// Convert KlineInterval to market.Interval
		var marketInterval market.Interval
		switch kline.Interval {
		case model.KlineInterval1m:
			marketInterval = market.Interval1m
		case model.KlineInterval5m:
			marketInterval = market.Interval5m
		case model.KlineInterval15m:
			marketInterval = market.Interval15m
		case model.KlineInterval1h:
			marketInterval = market.Interval1h
		case model.KlineInterval4h:
			marketInterval = market.Interval4h
		case model.KlineInterval1d:
			marketInterval = market.Interval1d
		default:
			marketInterval = market.Interval1h // Default to 1h
		}

		candles[i] = &market.Candle{
			Symbol:    kline.Symbol,
			Interval:  marketInterval,
			OpenTime:  kline.OpenTime,
			CloseTime: kline.CloseTime,
			Open:      kline.Open,
			High:      kline.High,
			Low:       kline.Low,
			Close:     kline.Close,
			Volume:    kline.Volume,
		}
	}

	return candles, nil
}

// GetOrderBookLegacy implements the legacy order book method
func (a *MarketDataServiceAdapter) GetOrderBookLegacy(ctx context.Context, symbol string, depth int) (*market.OrderBook, error) {
	orderBook, err := a.GetOrderBook(ctx, symbol, depth)
	if err != nil {
		return nil, err
	}

	// Convert to legacy model
	legacyBids := make([]market.OrderBookEntry, len(orderBook.Bids))
	for i, bid := range orderBook.Bids {
		legacyBids[i] = market.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	legacyAsks := make([]market.OrderBookEntry, len(orderBook.Asks))
	for i, ask := range orderBook.Asks {
		legacyAsks[i] = market.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	return &market.OrderBook{
		Symbol:      orderBook.Symbol,
		Bids:        legacyBids,
		Asks:        legacyAsks,
		LastUpdated: orderBook.Timestamp,
		Exchange:    "", // No exchange in model.OrderBook, set empty
	}, nil
}

// GetAllSymbolsLegacy implements the legacy get all symbols method
func (a *MarketDataServiceAdapter) GetAllSymbolsLegacy(ctx context.Context) ([]*market.Symbol, error) {
	symbols, err := a.GetAllSymbols(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to legacy model
	legacySymbols := make([]*market.Symbol, len(symbols))
	for i, symbol := range symbols {
		legacySymbols[i] = &market.Symbol{
			Symbol:     symbol.Symbol,
			BaseAsset:  symbol.BaseAsset,
			QuoteAsset: symbol.QuoteAsset,
			Status:     string(symbol.Status),
		}
	}

	return legacySymbols, nil
}

// GetSymbolInfoLegacy implements the legacy get symbol info method
func (a *MarketDataServiceAdapter) GetSymbolInfoLegacy(ctx context.Context, symbol string) (*market.Symbol, error) {
	symbolInfo, err := a.GetSymbolInfo(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Convert to legacy model
	return &market.Symbol{
		Symbol:     symbolInfo.Symbol,
		BaseAsset:  symbolInfo.BaseAsset,
		QuoteAsset: symbolInfo.QuoteAsset,
		Status:     string(symbolInfo.Status),
	}, nil
}

// GetHistoricalPricesLegacy implements the legacy get historical prices method
func (a *MarketDataServiceAdapter) GetHistoricalPricesLegacy(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*market.Candle, error) {
	klines, err := a.GetHistoricalPrices(ctx, symbol, from, to, interval)
	if err != nil {
		return nil, err
	}

	// Convert to legacy model
	candles := make([]*market.Candle, len(klines))
	for i, kline := range klines {
		// Convert KlineInterval to market.Interval
		var marketInterval market.Interval
		switch kline.Interval {
		case model.KlineInterval1m:
			marketInterval = market.Interval1m
		case model.KlineInterval5m:
			marketInterval = market.Interval5m
		case model.KlineInterval15m:
			marketInterval = market.Interval15m
		case model.KlineInterval1h:
			marketInterval = market.Interval1h
		case model.KlineInterval4h:
			marketInterval = market.Interval4h
		case model.KlineInterval1d:
			marketInterval = market.Interval1d
		default:
			marketInterval = market.Interval1h // Default to 1h
		}

		candles[i] = &market.Candle{
			Symbol:    kline.Symbol,
			Interval:  marketInterval,
			OpenTime:  kline.OpenTime,
			CloseTime: kline.CloseTime,
			Open:      kline.Open,
			High:      kline.High,
			Low:       kline.Low,
			Close:     kline.Close,
			Volume:    kline.Volume,
		}
	}

	return candles, nil
}
