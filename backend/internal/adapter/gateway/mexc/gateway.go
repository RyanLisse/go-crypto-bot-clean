package mexc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MEXCGateway implements the gateway to MEXC exchange
// It adapts the MEXC client to the domain interfaces
type MEXCGateway struct {
	client port.MEXCClient
	logger *zerolog.Logger
}

// NewMEXCGateway creates a new MEXC gateway
func NewMEXCGateway(client port.MEXCClient, logger *zerolog.Logger) *MEXCGateway {
	return &MEXCGateway{
		client: client,
		logger: logger,
	}
}

// GetTicker fetches current ticker data for a symbol
func (g *MEXCGateway) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	g.logger.Debug().Str("symbol", symbol).Msg("Fetching ticker from MEXC")

	// Get market data from MEXC client
	marketData, err := g.client.GetMarketData(ctx, symbol)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get market data from MEXC")
		return nil, fmt.Errorf("failed to get market data from MEXC: %w", err)
	}

	// Convert to domain model
	ticker := &market.Ticker{
		Symbol:        symbol,
		Exchange:      "mexc",
		Price:         marketData.LastPrice,
		Volume:        marketData.Volume,
		High24h:       marketData.HighPrice,
		Low24h:        marketData.LowPrice,
		PriceChange:   marketData.PriceChange,
		PercentChange: marketData.PriceChangePercent,
		LastUpdated:   time.Now(),
	}

	return ticker, nil
}

// GetCandles fetches historical candle data
func (g *MEXCGateway) GetCandles(ctx context.Context, symbol string, interval market.Interval, limit int) ([]*market.Candle, error) {
	g.logger.Debug().
		Str("symbol", symbol).
		Str("interval", string(interval)).
		Int("limit", limit).
		Msg("Fetching candles from MEXC")

	// Convert market.Interval to model.KlineInterval
	modelInterval := model.KlineInterval(interval)

	// Get klines from MEXC client
	klines, err := g.client.GetKlines(ctx, symbol, modelInterval, limit)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Str("interval", string(interval)).
			Msg("Failed to get klines from MEXC")
		return nil, fmt.Errorf("failed to get klines from MEXC: %w", err)
	}

	// Convert to domain model
	candles := make([]*market.Candle, len(klines))
	for i, kline := range klines {
		candles[i] = &market.Candle{
			Symbol:    symbol,
			Exchange:  "mexc",
			Interval:  interval,
			OpenTime:  kline.OpenTime,
			CloseTime: kline.CloseTime,
			Open:      kline.Open,
			High:      kline.High,
			Low:       kline.Low,
			Close:     kline.Close,
			Volume:    kline.Volume,
			Complete:  true, // Assume historical candles are complete
		}
	}

	return candles, nil
}

// GetOrderBook fetches current order book data
func (g *MEXCGateway) GetOrderBook(ctx context.Context, symbol string, limit int) (*market.OrderBook, error) {
	g.logger.Debug().
		Str("symbol", symbol).
		Int("limit", limit).
		Msg("Fetching order book from MEXC")

	// Get order book from MEXC client
	orderBook, err := g.client.GetOrderBook(ctx, symbol, limit)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Msg("Failed to get order book from MEXC")
		return nil, fmt.Errorf("failed to get order book from MEXC: %w", err)
	}

	// Convert to domain model
	domainOrderBook := &market.OrderBook{
		Symbol:      symbol,
		Exchange:    "mexc",
		LastUpdated: time.Now(),
		Bids:        make([]market.OrderBookEntry, len(orderBook.Bids)),
		Asks:        make([]market.OrderBookEntry, len(orderBook.Asks)),
	}

	// Convert bids
	for i, bid := range orderBook.Bids {
		domainOrderBook.Bids[i] = market.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	// Convert asks
	for i, ask := range orderBook.Asks {
		domainOrderBook.Asks[i] = market.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	return domainOrderBook, nil
}

// GetSymbols fetches available trading symbols
func (g *MEXCGateway) GetSymbols(ctx context.Context) ([]*market.Symbol, error) {
	g.logger.Debug().Msg("Fetching symbols from MEXC")

	// Get exchange info from MEXC client
	exchangeInfo, err := g.client.GetExchangeInfo(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get exchange info from MEXC")
		return nil, fmt.Errorf("failed to get exchange info from MEXC: %w", err)
	}

	// Convert to domain model
	symbols := make([]*market.Symbol, len(exchangeInfo.Symbols))
	for i, symbol := range exchangeInfo.Symbols {
		symbols[i] = &market.Symbol{
			Symbol:              symbol.Symbol,
			BaseAsset:           symbol.BaseAsset,
			QuoteAsset:          symbol.QuoteAsset,
			Status:              symbol.Status,
			BaseAssetPrecision:  symbol.BaseAssetPrecision,
			QuoteAssetPrecision: symbol.QuoteAssetPrecision,
			MinNotional:         parseStringToFloat64(symbol.MinNotional),
			MinQty:              parseStringToFloat64(symbol.MinLotSize),
			MaxQty:              parseStringToFloat64(symbol.MaxLotSize),
			QtyPrecision:        symbol.BaseAssetPrecision,
			PricePrecision:      symbol.QuoteAssetPrecision,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}
	}

	return symbols, nil
}

// GetSymbolInfo fetches detailed information about a trading symbol
func (g *MEXCGateway) GetSymbolInfo(ctx context.Context, symbol string) (*market.Symbol, error) {
	g.logger.Debug().Str("symbol", symbol).Msg("Fetching symbol info from MEXC")

	// Get symbol info from MEXC client
	symbolInfo, err := g.client.GetSymbolInfo(ctx, symbol)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get symbol info from MEXC")
		return nil, fmt.Errorf("failed to get symbol info from MEXC: %w", err)
	}

	// Convert to domain model
	domainSymbol := &market.Symbol{
		Symbol:              symbolInfo.Symbol,
		BaseAsset:           symbolInfo.BaseAsset,
		QuoteAsset:          symbolInfo.QuoteAsset,
		Status:              symbolInfo.Status,
		BaseAssetPrecision:  symbolInfo.BaseAssetPrecision,
		QuoteAssetPrecision: symbolInfo.QuoteAssetPrecision,
		MinNotional:         parseStringToFloat64(symbolInfo.MinNotional),
		MinQty:              parseStringToFloat64(symbolInfo.MinLotSize),
		MaxQty:              parseStringToFloat64(symbolInfo.MaxLotSize),
		QtyPrecision:        symbolInfo.BaseAssetPrecision,
		PricePrecision:      symbolInfo.QuoteAssetPrecision,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return domainSymbol, nil
}

// Helper function to parse string to float64
func parseStringToFloat64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return val
}

// GetNewCoins fetches information about newly listed coins
func (g *MEXCGateway) GetNewCoins(ctx context.Context) ([]*model.NewCoin, error) {
	g.logger.Debug().Msg("Fetching new coins from MEXC")

	// Get new listings from MEXC client
	newListings, err := g.client.GetNewListings(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get new listings from MEXC")
		return nil, fmt.Errorf("failed to get new listings from MEXC: %w", err)
	}

	return newListings, nil
}

// GetAccount fetches account information
func (g *MEXCGateway) GetAccount(ctx context.Context) (*model.Wallet, error) {
	g.logger.Debug().Msg("Fetching account information from MEXC")

	// Get account from MEXC client
	wallet, err := g.client.GetAccount(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get account from MEXC")
		return nil, fmt.Errorf("failed to get account from MEXC: %w", err)
	}

	return wallet, nil
}

// PlaceOrder places a new order
func (g *MEXCGateway) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	g.logger.Debug().
		Str("symbol", symbol).
		Str("side", string(side)).
		Str("type", string(orderType)).
		Float64("quantity", quantity).
		Float64("price", price).
		Str("timeInForce", string(timeInForce)).
		Msg("Placing order on MEXC")

	// Place order using MEXC client
	order, err := g.client.PlaceOrder(ctx, symbol, side, orderType, quantity, price, timeInForce)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Str("side", string(side)).
			Msg("Failed to place order on MEXC")
		return nil, fmt.Errorf("failed to place order on MEXC: %w", err)
	}

	return order, nil
}

// CancelOrder cancels an existing order
func (g *MEXCGateway) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	g.logger.Debug().
		Str("symbol", symbol).
		Str("orderID", orderID).
		Msg("Cancelling order on MEXC")

	// Cancel order using MEXC client
	err := g.client.CancelOrder(ctx, symbol, orderID)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Str("orderID", orderID).
			Msg("Failed to cancel order on MEXC")
		return fmt.Errorf("failed to cancel order on MEXC: %w", err)
	}

	return nil
}

// GetOrderStatus gets the status of an order
func (g *MEXCGateway) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	g.logger.Debug().
		Str("symbol", symbol).
		Str("orderID", orderID).
		Msg("Getting order status from MEXC")

	// Get order status using MEXC client
	order, err := g.client.GetOrderStatus(ctx, symbol, orderID)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Str("orderID", orderID).
			Msg("Failed to get order status from MEXC")
		return nil, fmt.Errorf("failed to get order status from MEXC: %w", err)
	}

	return order, nil
}

// GetOpenOrders gets all open orders
func (g *MEXCGateway) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	g.logger.Debug().
		Str("symbol", symbol).
		Msg("Getting open orders from MEXC")

	// Get open orders using MEXC client
	orders, err := g.client.GetOpenOrders(ctx, symbol)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Msg("Failed to get open orders from MEXC")
		return nil, fmt.Errorf("failed to get open orders from MEXC: %w", err)
	}

	return orders, nil
}

// GetOrderHistory gets order history
func (g *MEXCGateway) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	g.logger.Debug().
		Str("symbol", symbol).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting order history from MEXC")

	// Get order history using MEXC client
	orders, err := g.client.GetOrderHistory(ctx, symbol, limit, offset)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("symbol", symbol).
			Msg("Failed to get order history from MEXC")
		return nil, fmt.Errorf("failed to get order history from MEXC: %w", err)
	}

	return orders, nil
}
