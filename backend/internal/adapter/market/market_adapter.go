package market

import (
	"context"
	"time"

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
func (a *MarketDataServiceAdapter) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	return a.service.GetTicker(ctx, symbol)
}

// GetCandles adapts the GetCandles method
func (a *MarketDataServiceAdapter) GetCandles(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error) {
	return a.service.GetCandles(ctx, symbol, interval, limit)
}

// GetOrderBook adapts the GetOrderBook method
func (a *MarketDataServiceAdapter) GetOrderBook(ctx context.Context, symbol string, depth int) (*market.OrderBook, error) {
	// This is a simplified implementation since the actual OrderBook method might differ
	a.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Adapter: Getting order book")

	// Try to get from the service
	ticker, err := a.service.GetTicker(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Create a simple order book with the ticker price
	orderBook := &market.OrderBook{
		Symbol: symbol,
		Bids: []market.OrderBookEntry{
			{
				Price:    ticker.Price * 0.99, // Simulate a bid slightly below current price
				Quantity: 1.0,
			},
		},
		Asks: []market.OrderBookEntry{
			{
				Price:    ticker.Price * 1.01, // Simulate an ask slightly above current price
				Quantity: 1.0,
			},
		},
	}

	return orderBook, nil
}

// GetAllSymbols adapts the GetAllSymbols method
func (a *MarketDataServiceAdapter) GetAllSymbols(ctx context.Context) ([]*market.Symbol, error) {
	// This is a simplified implementation
	a.logger.Debug().Msg("Adapter: Getting all symbols")

	// Return a minimal implementation with a few common symbols
	return []*market.Symbol{
		{
			Symbol:     "BTCUSDT",
			BaseAsset:  "BTC",
			QuoteAsset: "USDT",
			Status:     "TRADING",
		},
		{
			Symbol:     "ETHUSDT",
			BaseAsset:  "ETH",
			QuoteAsset: "USDT",
			Status:     "TRADING",
		},
		{
			Symbol:     "BNBUSDT",
			BaseAsset:  "BNB",
			QuoteAsset: "USDT",
			Status:     "TRADING",
		},
	}, nil
}

// GetSymbolInfo adapts the GetSymbolInfo method
func (a *MarketDataServiceAdapter) GetSymbolInfo(ctx context.Context, symbol string) (*market.Symbol, error) {
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
	return &market.Symbol{
		Symbol:     symbol,
		BaseAsset:  symbol[:len(symbol)-4], // Assume the last 4 chars are the quote asset
		QuoteAsset: symbol[len(symbol)-4:],
		Status:     "TRADING",
	}, nil
}

// GetHistoricalPrices adapts the GetHistoricalPrices method
func (a *MarketDataServiceAdapter) GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*market.Candle, error) {
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
