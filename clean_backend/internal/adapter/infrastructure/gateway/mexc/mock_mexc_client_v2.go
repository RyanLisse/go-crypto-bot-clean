package mexc

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MockMEXCClientV2 implements the port.MEXCClient interface with mock data
type MockMEXCClientV2 struct {
	logger *zerolog.Logger
}

// NewMockMEXCClientV2 creates a new mock MEXC client
func NewMockMEXCClientV2(logger *zerolog.Logger) port.MEXCClient {
	return &MockMEXCClientV2{
		logger: logger,
	}
}

// GetMarketData retrieves mock market data for a symbol
func (c *MockMEXCClientV2) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	c.logger.Debug().
		Str("symbol", symbol).
		Str("component", "MockMEXCClientV2").
		Str("method", "GetMarketData").
		Str("data_source", "MOCK").
		Msg("Returning mock market data")

	// Return mock data
	return &model.Ticker{
		Symbol:    symbol,
		Exchange:  "MEXC",
		LastPrice: 1000.0,
		Volume:    1000000.0,
		Timestamp: time.Now(),
	}, nil
}

// GetKlines retrieves mock kline/candlestick data for a symbol
func (c *MockMEXCClientV2) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	c.logger.Debug().
		Str("symbol", symbol).
		Str("interval", string(interval)).
		Int("limit", limit).
		Str("component", "MockMEXCClientV2").
		Str("method", "GetKlines").
		Str("data_source", "MOCK").
		Msg("Returning mock klines")

	// Return mock data
	klines := make([]*model.Kline, 0, limit)
	for i := 0; i < limit; i++ {
		klines = append(klines, &model.Kline{
			Symbol:    symbol,
			Exchange:  "MEXC",
			Interval:  interval,
			OpenTime:  time.Now().Add(-time.Duration(i) * time.Hour),
			CloseTime: time.Now().Add(-time.Duration(i-1) * time.Hour),
			Open:      1000.0,
			High:      1010.0,
			Low:       990.0,
			Close:     1005.0,
			Volume:    1000000.0,
		})
	}
	return klines, nil
}

// GetOrderBook retrieves a mock order book for a symbol
func (c *MockMEXCClientV2) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	c.logger.Debug().
		Str("symbol", symbol).
		Int("depth", depth).
		Str("component", "MockMEXCClientV2").
		Str("method", "GetOrderBook").
		Str("data_source", "MOCK").
		Msg("Returning mock order book")

	// Return mock data
	return &model.OrderBook{
		Symbol:    symbol,
		Exchange:  "MEXC",
		Timestamp: time.Now(),
		Bids:      []model.OrderBookEntry{{Price: 990.0, Quantity: 1.0}, {Price: 980.0, Quantity: 2.0}},
		Asks:      []model.OrderBookEntry{{Price: 1010.0, Quantity: 1.0}, {Price: 1020.0, Quantity: 2.0}},
	}, nil
}

// GetSymbols retrieves mock symbols
func (c *MockMEXCClientV2) GetSymbols(ctx context.Context) ([]*model.Symbol, error) {
	c.logger.Debug().
		Str("component", "MockMEXCClientV2").
		Str("method", "GetSymbols").
		Str("data_source", "MOCK").
		Msg("Returning mock symbols")

	// Return mock data
	return []*model.Symbol{
		{
			Symbol:             "BTCUSDT",
			BaseAsset:          "BTC",
			QuoteAsset:         "USDT",
			Exchange:           "MEXC",
			MinPrice:           0.01,
			MaxPrice:           100000.0,
			TickSize:           0.01,
			MinQuantity:        0.0001,
			MaxQuantity:        1000.0,
			StepSize:           0.0001,
			MinNotional:        10.0,
			Status:             model.SymbolStatusTrading,
			PricePrecision:     2,
			QuantityPrecision:  4,
			AllowedOrderTypes:  []string{"LIMIT", "MARKET"},
		},
		{
			Symbol:             "ETHUSDT",
			BaseAsset:          "ETH",
			QuoteAsset:         "USDT",
			Exchange:           "MEXC",
			MinPrice:           0.01,
			MaxPrice:           100000.0,
			TickSize:           0.01,
			MinQuantity:        0.001,
			MaxQuantity:        1000.0,
			StepSize:           0.001,
			MinNotional:        10.0,
			Status:             model.SymbolStatusTrading,
			PricePrecision:     2,
			QuantityPrecision:  3,
			AllowedOrderTypes:  []string{"LIMIT", "MARKET"},
		},
	}, nil
}

// GetSymbol retrieves mock information about a specific symbol
func (c *MockMEXCClientV2) GetSymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	c.logger.Debug().
		Str("symbol", symbol).
		Str("component", "MockMEXCClientV2").
		Str("method", "GetSymbol").
		Str("data_source", "MOCK").
		Msg("Returning mock symbol")

	// Return mock data
	return &model.Symbol{
		Symbol:             symbol,
		BaseAsset:          "BTC",
		QuoteAsset:         "USDT",
		Exchange:           "MEXC",
		MinPrice:           0.01,
		MaxPrice:           100000.0,
		TickSize:           0.01,
		MinQuantity:        0.0001,
		MaxQuantity:        1000.0,
		StepSize:           0.0001,
		MinNotional:        10.0,
		Status:             model.SymbolStatusTrading,
		PricePrecision:     2,
		QuantityPrecision:  4,
		AllowedOrderTypes:  []string{"LIMIT", "MARKET"},
	}, nil
}

// GetServerTime retrieves a mock server time
func (c *MockMEXCClientV2) GetServerTime(ctx context.Context) (time.Time, error) {
	c.logger.Debug().
		Str("component", "MockMEXCClientV2").
		Str("method", "GetServerTime").
		Str("data_source", "MOCK").
		Msg("Returning mock server time")

	// Return mock data
	return time.Now(), nil
}

// GetExchangeInfo retrieves mock exchange information
func (c *MockMEXCClientV2) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	c.logger.Debug().
		Str("component", "MockMEXCClientV2").
		Str("method", "GetExchangeInfo").
		Str("data_source", "MOCK").
		Msg("Returning mock exchange info")

	// Return mock data
	return &model.ExchangeInfo{
		Symbols: []model.SymbolInfo{
			{
				Symbol:               "BTCUSDT",
				Status:               "TRADING",
				BaseAsset:            "BTC",
				BaseAssetPrecision:   8,
				QuoteAsset:           "USDT",
				QuoteAssetPrecision:  8,
				OrderTypes:           []string{"LIMIT", "MARKET"},
				IsSpotTradingAllowed: true,
				Permissions:          []string{"SPOT"},
			},
		},
	}, nil
}

// PlaceOrder places a mock order on the exchange
func (c *MockMEXCClientV2) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	c.logger.Debug().
		Str("symbol", symbol).
		Str("side", string(side)).
		Str("orderType", string(orderType)).
		Float64("quantity", quantity).
		Float64("price", price).
		Str("timeInForce", string(timeInForce)).
		Str("component", "MockMEXCClientV2").
		Str("method", "PlaceOrder").
		Str("data_source", "MOCK").
		Msg("Placing mock order")

	// Return mock data
	orderID := fmt.Sprintf("mock-order-%d", time.Now().UnixNano())
	clientOrderID := fmt.Sprintf("mock-client-order-%d", time.Now().UnixNano())
	now := time.Now()
	
	return &model.Order{
		ID:            orderID,
		OrderID:       orderID,
		ClientOrderID: clientOrderID,
		Symbol:        symbol,
		Side:          side,
		Type:          orderType,
		Status:        model.OrderStatusNew,
		Price:         price,
		Quantity:      quantity,
		ExecutedQty:   0,
		TimeInForce:   timeInForce,
		CreatedAt:     now,
		UpdatedAt:     now,
		Exchange:      "MEXC",
	}, nil
}
