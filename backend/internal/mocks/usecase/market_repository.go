package mocks

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/stretchr/testify/mock"
)

// MockMarketRepository is a mock implementation of port.MarketRepository
type MockMarketRepository struct {
	mock.Mock
}

// SaveTicker mocks the SaveTicker method
func (m *MockMarketRepository) SaveTicker(ctx context.Context, ticker *market.Ticker) error {
	args := m.Called(ctx, ticker)
	return args.Error(0)
}

// GetTicker mocks the GetTicker method
func (m *MockMarketRepository) GetTicker(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

// GetAllTickers mocks the GetAllTickers method
func (m *MockMarketRepository) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

// GetTickerHistory mocks the GetTickerHistory method
func (m *MockMarketRepository) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
	args := m.Called(ctx, symbol, exchange, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

// SaveCandle mocks the SaveCandle method
func (m *MockMarketRepository) SaveCandle(ctx context.Context, candle *market.Candle) error {
	args := m.Called(ctx, candle)
	return args.Error(0)
}

// SaveCandles mocks the SaveCandles method
func (m *MockMarketRepository) SaveCandles(ctx context.Context, candles []*market.Candle) error {
	args := m.Called(ctx, candles)
	return args.Error(0)
}

// GetCandle mocks the GetCandle method
func (m *MockMarketRepository) GetCandle(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, openTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

// GetCandles mocks the GetCandles method
func (m *MockMarketRepository) GetCandles(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, start, end, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Candle), args.Error(1)
}

// GetLatestCandle mocks the GetLatestCandle method
func (m *MockMarketRepository) GetLatestCandle(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

// PurgeOldData mocks the PurgeOldData method
func (m *MockMarketRepository) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

// GetLatestTickers mocks the GetLatestTickers method
func (m *MockMarketRepository) GetLatestTickers(ctx context.Context, limit int) ([]*market.Ticker, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

// GetTickersBySymbol mocks the GetTickersBySymbol method
func (m *MockMarketRepository) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
	args := m.Called(ctx, symbol, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

// GetOrderBook mocks the GetOrderBook method
func (m *MockMarketRepository) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error) {
	args := m.Called(ctx, symbol, exchange, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.OrderBook), args.Error(1)
}
