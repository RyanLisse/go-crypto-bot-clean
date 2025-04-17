package mocks

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/stretchr/testify/mock"
)

// MarketDataService is a mock implementation of port.MarketDataService
type MarketDataService struct {
	mock.Mock
}

// GetTicker mocks the GetTicker method
func (m *MarketDataService) GetTicker(ctx context.Context, symbol string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticker), args.Error(1)
}

// GetCandles mocks the GetCandles method
func (m *MarketDataService) GetCandles(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Kline), args.Error(1)
}

// GetOrderBook mocks the GetOrderBook method
func (m *MarketDataService) GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderBook), args.Error(1)
}

// GetAllSymbols mocks the GetAllSymbols method
func (m *MarketDataService) GetAllSymbols(ctx context.Context) ([]*model.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Symbol), args.Error(1)
}

// GetSymbolInfo mocks the GetSymbolInfo method
func (m *MarketDataService) GetSymbolInfo(ctx context.Context, symbol string) (*model.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Symbol), args.Error(1)
}

// GetHistoricalPrices mocks the GetHistoricalPrices method
func (m *MarketDataService) GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time, interval model.KlineInterval) ([]*model.Kline, error) {
	args := m.Called(ctx, symbol, from, to, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Kline), args.Error(1)
}

// GetTickerLegacy mocks the GetTickerLegacy method
func (m *MarketDataService) GetTickerLegacy(ctx context.Context, symbol string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

// GetCandlesLegacy mocks the GetCandlesLegacy method
func (m *MarketDataService) GetCandlesLegacy(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, interval, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Candle), args.Error(1)
}

// GetOrderBookLegacy mocks the GetOrderBookLegacy method
func (m *MarketDataService) GetOrderBookLegacy(ctx context.Context, symbol string, depth int) (*market.OrderBook, error) {
	args := m.Called(ctx, symbol, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.OrderBook), args.Error(1)
}

// GetAllSymbolsLegacy mocks the GetAllSymbolsLegacy method
func (m *MarketDataService) GetAllSymbolsLegacy(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

// GetSymbolInfoLegacy mocks the GetSymbolInfoLegacy method
func (m *MarketDataService) GetSymbolInfoLegacy(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Symbol), args.Error(1)
}

// GetHistoricalPricesLegacy mocks the GetHistoricalPricesLegacy method
func (m *MarketDataService) GetHistoricalPricesLegacy(ctx context.Context, symbol string, from, to time.Time, interval string) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, from, to, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Candle), args.Error(1)
}
