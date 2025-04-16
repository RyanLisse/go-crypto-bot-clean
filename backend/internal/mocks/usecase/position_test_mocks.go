package usecase

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for position use case tests
type PositionMockMarketRepository struct {
	mock.Mock
}

func (m *PositionMockMarketRepository) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*model.OrderBook, error) {
	args := m.Called(ctx, symbol, exchange, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OrderBook), args.Error(1)
}

func (m *PositionMockMarketRepository) GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error) {
	args := m.Called(ctx, symbol, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticker), args.Error(1)
}

func (m *PositionMockMarketRepository) GetLatestCandle(ctx context.Context, symbol, exchange, interval string) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

func (m *PositionMockMarketRepository) SaveTicker(ctx context.Context, ticker *model.Ticker) error {
	args := m.Called(ctx, ticker)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) SaveCandle(ctx context.Context, candle *market.Candle) error {
	args := m.Called(ctx, candle)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) SaveOrderBook(ctx context.Context, orderBook *model.OrderBook) error {
	args := m.Called(ctx, orderBook)
	return args.Error(0)
}

// Add missing method for MarketRepository interface
func (m *PositionMockMarketRepository) GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error) {
	args := m.Called(ctx, exchange)
	return args.Get(0).([]*model.Ticker), args.Error(1)
}

// Add additional required methods
func (m *PositionMockMarketRepository) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*model.Ticker, error) {
	args := m.Called(ctx, symbol, exchange, start, end)
	return args.Get(0).([]*model.Ticker), args.Error(1)
}

func (m *PositionMockMarketRepository) SaveCandles(ctx context.Context, candles []*market.Candle) error {
	args := m.Called(ctx, candles)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) GetCandle(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, openTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

func (m *PositionMockMarketRepository) GetCandles(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, start, end, limit)
	return args.Get(0).([]*market.Candle), args.Error(1)
}

func (m *PositionMockMarketRepository) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*model.Ticker), args.Error(1)
}

func (m *PositionMockMarketRepository) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*model.Ticker, error) {
	args := m.Called(ctx, symbol, limit)
	return args.Get(0).([]*model.Ticker), args.Error(1)
}

type PositionMockSymbolRepository struct {
	mock.Mock
}

func (m *PositionMockSymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Symbol), args.Error(1)
}

func (m *PositionMockSymbolRepository) GetAll(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *PositionMockSymbolRepository) Save(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

// Add missing method for SymbolRepository interface
func (m *PositionMockSymbolRepository) Create(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

// Add additional required methods
func (m *PositionMockSymbolRepository) GetByExchange(ctx context.Context, exchange string) ([]*market.Symbol, error) {
	args := m.Called(ctx, exchange)
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *PositionMockSymbolRepository) Update(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *PositionMockSymbolRepository) Delete(ctx context.Context, symbol string) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

type PositionMockRepository struct {
	mock.Mock
}

func (m *PositionMockRepository) Create(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *PositionMockRepository) GetByID(ctx context.Context, id string) (*model.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *PositionMockRepository) Update(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *PositionMockRepository) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	args := m.Called(ctx, positionType)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, from, to, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *PositionMockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *PositionMockRepository) GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetBySymbolAndUser(ctx context.Context, symbol, userID string, page, limit int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, userID, page, limit)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetOpenPositionsByUserID(ctx context.Context, userID string) ([]*model.Position, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Position), args.Error(1)
}
