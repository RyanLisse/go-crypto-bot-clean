package mexc

import (
	"context"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
)

// MarketDataProvider implements market data fetching from MEXC exchange
type MarketDataProvider struct {
	marketService *service.MarketDataService
	client        port.MEXCClient
	logger        *zerolog.Logger
}

// NewMarketDataProvider creates a new MEXC market data provider
func NewMarketDataProvider(marketService *service.MarketDataService, client port.MEXCClient, logger *zerolog.Logger) *MarketDataProvider {
	return &MarketDataProvider{
		marketService: marketService,
		client:        client,
		logger:        logger,
	}
}

// GetTicker fetches current ticker data for a symbol
func (p *MarketDataProvider) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	p.logger.Debug().Str("symbol", symbol).Msg("Getting ticker from MEXC")

	// Use the MEXC client to get ticker data
	modelTicker, err := p.client.GetMarketData(ctx, symbol)
	if err != nil {
		p.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get ticker from MEXC")
		return nil, err
	}

	// Convert model.Ticker to market.Ticker
	ticker := &market.Ticker{
		Symbol:        modelTicker.Symbol,
		Exchange:      "mexc",
		Price:         modelTicker.LastPrice,
		PriceChange:   modelTicker.PriceChange,
		PercentChange: modelTicker.PriceChangePercent,
		High24h:       modelTicker.HighPrice,
		Low24h:        modelTicker.LowPrice,
		Volume:        modelTicker.Volume,
		LastUpdated:   time.Now(),
	}

	return ticker, nil
}

// GetCandles fetches historical candle data
func (p *MarketDataProvider) GetCandles(ctx context.Context, symbol string, interval market.Interval, limit int) ([]*market.Candle, error) {
	p.logger.Debug().Str("symbol", symbol).Str("interval", string(interval)).Int("limit", limit).Msg("Getting candles from MEXC")

	// Convert market.Interval to model.KlineInterval
	klineInterval := convertIntervalToKlineInterval(interval)

	// Use the MEXC client to get candle data
	modelKlines, err := p.client.GetKlines(ctx, symbol, model.KlineInterval(klineInterval), limit)
	if err != nil {
		p.logger.Error().Err(err).Str("symbol", symbol).Str("interval", string(interval)).Msg("Failed to get candles from MEXC")
		return nil, err
	}

	// Convert model.Kline to market.Candle
	candles := make([]*market.Candle, 0, len(modelKlines))
	for _, kline := range modelKlines {
		candle := &market.Candle{
			Symbol:    kline.Symbol,
			Interval:  interval,
			OpenTime:  kline.OpenTime,
			CloseTime: kline.CloseTime,
			Open:      kline.Open,
			High:      kline.High,
			Low:       kline.Low,
			Close:     kline.Close,
			Volume:    kline.Volume,
		}
		candles = append(candles, candle)
	}

	return candles, nil
}

// Helper function to convert market.Interval to string for MEXC API
func convertIntervalToKlineInterval(interval market.Interval) string {
	switch interval {
	case market.Interval1m:
		return "1m"
	case market.Interval5m:
		return "5m"
	case market.Interval15m:
		return "15m"
	case market.Interval30m:
		return "30m"
	case market.Interval1h:
		return "1h"
	case market.Interval4h:
		return "4h"
	case market.Interval1d:
		return "1d"
	case market.Interval1w:
		return "1w"
	default:
		return "1h"
	}
}

// GetOrderBook fetches current order book data
func (p *MarketDataProvider) GetOrderBook(ctx context.Context, symbol string, limit int) (*market.OrderBook, error) {
	p.logger.Debug().Str("symbol", symbol).Int("limit", limit).Msg("Getting order book from MEXC")

	// Use the MEXC client to get order book data
	modelOrderBook, err := p.client.GetOrderBook(ctx, symbol, limit)
	if err != nil {
		p.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get order book from MEXC")
		return nil, err
	}

	// Convert model.OrderBook to market.OrderBook
	orderBook := &market.OrderBook{
		Symbol:      symbol,
		Bids:        make([]market.OrderBookEntry, len(modelOrderBook.Bids)),
		Asks:        make([]market.OrderBookEntry, len(modelOrderBook.Asks)),
		LastUpdated: time.Now(),
		Exchange:    "mexc",
	}

	// Convert bids
	for i, bid := range modelOrderBook.Bids {
		orderBook.Bids[i] = market.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	// Convert asks
	for i, ask := range modelOrderBook.Asks {
		orderBook.Asks[i] = market.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	return orderBook, nil
}

// GetSymbols fetches available trading symbols
func (p *MarketDataProvider) GetSymbols(ctx context.Context) ([]*market.Symbol, error) {
	p.logger.Debug().Msg("Getting symbols from MEXC")

	// Use the MEXC client to get exchange info
	exchangeInfo, err := p.client.GetExchangeInfo(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to get exchange info from MEXC")
		return nil, err
	}

	// Convert model.SymbolInfo to market.Symbol
	symbols := make([]*market.Symbol, 0, len(exchangeInfo.Symbols))
	for _, symbolInfo := range exchangeInfo.Symbols {
		// Only include symbols that are actively trading
		if symbolInfo.Status == "TRADING" {
			// Parse min/max values
			minPrice, _ := strconv.ParseFloat(symbolInfo.TickSize, 64)
			minQty, _ := strconv.ParseFloat(symbolInfo.MinLotSize, 64)
			maxQty, _ := strconv.ParseFloat(symbolInfo.MaxLotSize, 64)

			symbol := &market.Symbol{
				Symbol:         symbolInfo.Symbol,
				BaseAsset:      symbolInfo.BaseAsset,
				QuoteAsset:     symbolInfo.QuoteAsset,
				Status:         symbolInfo.Status,
				MinPrice:       minPrice,
				PricePrecision: symbolInfo.QuoteAssetPrecision,
				MinQty:         minQty,
				MaxQty:         maxQty,
				QtyPrecision:   symbolInfo.BaseAssetPrecision,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			symbols = append(symbols, symbol)
		}
	}

	return symbols, nil
}
