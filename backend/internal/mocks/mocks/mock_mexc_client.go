package mocks

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/stretchr/testify/mock"
)

// MockMEXCClient is a mock implementation of the MEXCClient interface
type MockMEXCClient struct {
	mock.Mock
}

// Market data methods
// GetTicker returns a mock ticker
func (m *MockMEXCClient) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol)
	var ticker *market.Ticker
	if args.Get(0) != nil {
		ticker = args.Get(0).(*market.Ticker)
	}
	return ticker, args.Error(1)
}

// GetCandles returns mock candles
func (m *MockMEXCClient) GetCandles(ctx context.Context, symbol string, interval market.Interval, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, interval, limit)
	var candles []*market.Candle
	if args.Get(0) != nil {
		candles = args.Get(0).([]*market.Candle)
	}
	return candles, args.Error(1)
}

// GetOrderBook returns a mock order book
func (m *MockMEXCClient) GetOrderBook(ctx context.Context, symbol string, limit int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, limit)
	var orderBook *model.OrderBook
	if args.Get(0) != nil {
		orderBook = args.Get(0).(*model.OrderBook)
	}
	return orderBook, args.Error(1)
}

// GetMarketData retrieves ticker data
func (m *MockMEXCClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	var ticker *model.Ticker
	if args.Get(0) != nil {
		ticker = args.Get(0).(*model.Ticker)
	}
	return ticker, args.Error(1)
}

// GetKlines returns mock klines
func (m *MockMEXCClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	var klines []*model.Kline
	if args.Get(0) != nil {
		klines = args.Get(0).([]*model.Kline)
	}
	return klines, args.Error(1)
}

// GetHistoricalCandles returns mock historical candles
func (m *MockMEXCClient) GetHistoricalCandles(ctx context.Context, symbol string, interval market.Interval, startTime, endTime time.Time) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, interval, startTime, endTime)
	var candles []*market.Candle
	if args.Get(0) != nil {
		candles = args.Get(0).([]*market.Candle)
	}
	return candles, args.Error(1)
}

// Symbol related methods
// GetSymbols returns mock symbols
func (m *MockMEXCClient) GetSymbols(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	var symbols []*market.Symbol
	if args.Get(0) != nil {
		symbols = args.Get(0).([]*market.Symbol)
	}
	return symbols, args.Error(1)
}

// GetSymbolInfo returns mock symbol info
func (m *MockMEXCClient) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) {
	args := m.Called(ctx, symbol)
	var info *model.SymbolInfo
	if args.Get(0) != nil {
		info = args.Get(0).(*model.SymbolInfo)
	}
	return info, args.Error(1)
}

// GetSymbolStatus checks if a symbol is currently tradeable
func (m *MockMEXCClient) GetSymbolStatus(ctx context.Context, symbol string) (model.Status, error) {
	args := m.Called(ctx, symbol)
	var status model.Status
	if args.Get(0) != nil {
		status = args.Get(0).(model.Status)
	}
	return status, args.Error(1)
}

// GetSymbolConstraints retrieves trading constraints for a symbol
func (m *MockMEXCClient) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	args := m.Called(ctx, symbol)
	var constraints *model.SymbolConstraints
	if args.Get(0) != nil {
		constraints = args.Get(0).(*model.SymbolConstraints)
	}
	return constraints, args.Error(1)
}

// GetTradingSchedule retrieves the listing and trading schedule for a symbol
func (m *MockMEXCClient) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	args := m.Called(ctx, symbol)
	var schedule model.TradingSchedule
	if args.Get(0) != nil {
		schedule = args.Get(0).(model.TradingSchedule)
	}
	return schedule, args.Error(1)
}

// Exchange info methods
// GetExchangeInfo returns mock exchange info
func (m *MockMEXCClient) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	args := m.Called(ctx)
	var info *model.ExchangeInfo
	if args.Get(0) != nil {
		info = args.Get(0).(*model.ExchangeInfo)
	}
	return info, args.Error(1)
}

// GetNewListings returns mock new coin listings
func (m *MockMEXCClient) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	args := m.Called(ctx)
	var coins []*model.NewCoin
	if args.Get(0) != nil {
		coins = args.Get(0).([]*model.NewCoin)
	}
	return coins, args.Error(1)
}

// GetNewCoins returns mock new coins (alias for GetNewListings for backward compatibility)
func (m *MockMEXCClient) GetNewCoins(ctx context.Context) ([]*model.NewCoin, error) {
	args := m.Called(ctx)
	var coins []*model.NewCoin
	if args.Get(0) != nil {
		coins = args.Get(0).([]*model.NewCoin)
	}
	return coins, args.Error(1)
}

// Order related methods
// PlaceOrder places a mock order
func (m *MockMEXCClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	args := m.Called(ctx, symbol, side, orderType, quantity, price, timeInForce)
	var order *model.Order
	if args.Get(0) != nil {
		order = args.Get(0).(*model.Order)
	}
	return order, args.Error(1)
}

// CancelOrder cancels a mock order
func (m *MockMEXCClient) CancelOrder(ctx context.Context, symbol, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

// GetOrderStatus returns a mock order status
func (m *MockMEXCClient) GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	var order *model.Order
	if args.Get(0) != nil {
		order = args.Get(0).(*model.Order)
	}
	return order, args.Error(1)
}

// GetOpenOrders returns mock open orders
func (m *MockMEXCClient) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	args := m.Called(ctx, symbol)
	var orders []*model.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]*model.Order)
	}
	return orders, args.Error(1)
}

// GetOrderHistory returns mock order history
func (m *MockMEXCClient) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, symbol, limit, offset)
	var orders []*model.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]*model.Order)
	}
	return orders, args.Error(1)
}

// Account related methods
// GetAccount returns a mock account
func (m *MockMEXCClient) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	var wallet *model.Wallet
	if args.Get(0) != nil {
		wallet = args.Get(0).(*model.Wallet)
	}
	return wallet, args.Error(1)
}

// Trading constraints methods

// Streaming methods
// StreamTicker streams mock tickers
func (m *MockMEXCClient) StreamTicker(ctx context.Context, symbols []string, callback func(*market.Ticker) error) error {
	args := m.Called(ctx, symbols, callback)
	return args.Error(0)
}

// StreamCandles streams mock candles
func (m *MockMEXCClient) StreamCandles(ctx context.Context, symbol string, interval market.Interval, callback func(*market.Candle) error) error {
	args := m.Called(ctx, symbol, interval, callback)
	return args.Error(0)
}
