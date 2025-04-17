package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockMEXCClient is a mock implementation of port.MEXCClient
type MockMEXCClient struct {
	mock.Mock
}

func (m *MockMEXCClient) GetNewListings(ctx context.Context) ([]*model.NewCoin, error) {
	args := m.Called(ctx)
	var coins []*model.NewCoin
	if arg0 := args.Get(0); arg0 != nil {
		coins = arg0.([]*model.NewCoin)
	}
	return coins, args.Error(1)
}

func (m *MockMEXCClient) GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) {
	args := m.Called(ctx, symbol)
	var info *model.SymbolInfo
	if arg0 := args.Get(0); arg0 != nil {
		info = arg0.(*model.SymbolInfo)
	}
	return info, args.Error(1)
}

func (m *MockMEXCClient) GetSymbolStatus(ctx context.Context, symbol string) (model.Status, error) {
	args := m.Called(ctx, symbol)
	var status model.Status
	if arg0 := args.Get(0); arg0 != nil {
		statusStr, ok := arg0.(string)
		if ok {
			if statusStr != "" {
				status = model.Status(statusStr)
			}
		} else {
			status, _ = arg0.(model.Status)
		}
	}
	return status, args.Error(1)
}

func (m *MockMEXCClient) GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error) {
	args := m.Called(ctx, symbol)
	var schedule model.TradingSchedule
	if arg0 := args.Get(0); arg0 != nil {
		schedule = arg0.(model.TradingSchedule)
	}
	return schedule, args.Error(1)
}

func (m *MockMEXCClient) GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error) {
	args := m.Called(ctx, symbol)
	var constraints *model.SymbolConstraints
	if arg0 := args.Get(0); arg0 != nil {
		constraints = arg0.(*model.SymbolConstraints)
	}
	return constraints, args.Error(1)
}

func (m *MockMEXCClient) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	args := m.Called(ctx)
	var info *model.ExchangeInfo
	if arg0 := args.Get(0); arg0 != nil {
		info = arg0.(*model.ExchangeInfo)
	}
	return info, args.Error(1)
}

func (m *MockMEXCClient) GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	var ticker *model.Ticker
	if arg0 := args.Get(0); arg0 != nil {
		ticker = arg0.(*model.Ticker)
	}
	return ticker, args.Error(1)
}

func (m *MockMEXCClient) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	var klines []*model.Kline
	if arg0 := args.Get(0); arg0 != nil {
		klines = arg0.([]*model.Kline)
	}
	return klines, args.Error(1)
}

func (m *MockMEXCClient) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, depth)
	var book *model.OrderBook
	if arg0 := args.Get(0); arg0 != nil {
		book = arg0.(*model.OrderBook)
	}
	return book, args.Error(1)
}

func (m *MockMEXCClient) GetAccount(ctx context.Context) (*model.Wallet, error) {
	args := m.Called(ctx)
	var wallet *model.Wallet
	if arg0 := args.Get(0); arg0 != nil {
		wallet = arg0.(*model.Wallet)
	}
	return wallet, args.Error(1)
}

func (m *MockMEXCClient) PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error) {
	args := m.Called(ctx, symbol, side, orderType, quantity, price, timeInForce)
	var order *model.Order
	if arg0 := args.Get(0); arg0 != nil {
		order = arg0.(*model.Order)
	}
	return order, args.Error(1)
}

func (m *MockMEXCClient) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	args := m.Called(ctx, symbol, orderID)
	return args.Error(0)
}

func (m *MockMEXCClient) GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error) {
	args := m.Called(ctx, symbol, orderID)
	var order *model.Order
	if arg0 := args.Get(0); arg0 != nil {
		order = arg0.(*model.Order)
	}
	return order, args.Error(1)
}

func (m *MockMEXCClient) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	args := m.Called(ctx, symbol)
	var orders []*model.Order
	if arg0 := args.Get(0); arg0 != nil {
		orders = arg0.([]*model.Order)
	}
	return orders, args.Error(1)
}

func (m *MockMEXCClient) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	args := m.Called(ctx, symbol, limit, offset)
	var orders []*model.Order
	if arg0 := args.Get(0); arg0 != nil {
		orders = arg0.([]*model.Order)
	}
	return orders, args.Error(1)
}
