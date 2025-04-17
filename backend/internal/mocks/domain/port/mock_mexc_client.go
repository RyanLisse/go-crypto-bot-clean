package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MEXCClient is a mock implementation of port.MEXCClient
type MEXCClient struct {
	mock.Mock
}

// GetNewListings mocks the GetNewListings method
func (m *MEXCClient) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

// GetSymbolInfo mocks the GetSymbolInfo method
func (m *MEXCClient) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SymbolInfo), args.Error(1)
}

// GetSymbolStatus mocks the GetSymbolStatus method
func (m *MEXCClient) GetSymbolStatus(ctx context.Context, symbol string) (model.CoinStatus, error) {
	args := m.Called(ctx, symbol)
	var status model.CoinStatus
	arg0 := args.Get(0)
	statusStr, ok := arg0.(string)
	if ok {
		if statusStr != "" {
			status = model.CoinStatus(statusStr)
		}
	} else {
		status, _ = arg0.(model.CoinStatus)
	}
	return status, args.Error(1)
}

// GetTradingSchedule mocks the GetTradingSchedule method
func (m *MEXCClient) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(model.TradingSchedule), args.Error(1)
}

// GetSymbolConstraints mocks the GetSymbolConstraints method
func (m *MEXCClient) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SymbolConstraints), args.Error(1)
}

// GetExchangeInfo mocks the GetExchangeInfo method
func (m *MEXCClient) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExchangeInfo), args.Error(1)
}

// GetMarketData mocks the GetMarketData method
func (m *MEXCClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticker), args.Error(1)
}

// GetKlines mocks the GetKlines method
func (m *MEXCClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Kline), args.Error(1)
}

// GetOrderBook mocks the GetOrderBook method
func (m *MEXCClient) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderBook), args.Error(1)
}

// GetAccount mocks the GetAccount method
func (m *MEXCClient) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// PlaceOrder mocks the PlaceOrder method
func (m *MEXCClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	args := m.Called(ctx, symbol, side, orderType, quantity, price, timeInForce)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// CancelOrder mocks the CancelOrder method
func (m *MEXCClient) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

// GetOrderStatus mocks the GetOrderStatus method
func (m *MEXCClient) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// GetOpenOrders mocks the GetOpenOrders method
func (m *MEXCClient) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

// GetOrderHistory mocks the GetOrderHistory method
func (m *MEXCClient) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, symbol, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}
