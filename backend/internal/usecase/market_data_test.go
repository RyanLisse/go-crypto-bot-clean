package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockMarketRepository struct {
	mock.Mock
}

func (m *MockMarketRepository) SaveTicker(ctx context.Context, ticker *market.Ticker) error {
	args := m.Called(ctx, ticker)
	return args.Error(0)
}

func (m *MockMarketRepository) GetTicker(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

func (m *MockMarketRepository) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

func (m *MockMarketRepository) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
	args := m.Called(ctx, symbol, exchange, start, end)
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

func (m *MockMarketRepository) SaveCandle(ctx context.Context, candle *market.Candle) error {
	args := m.Called(ctx, candle)
	return args.Error(0)
}

func (m *MockMarketRepository) SaveCandles(ctx context.Context, candles []*market.Candle) error {
	args := m.Called(ctx, candles)
	return args.Error(0)
}

func (m *MockMarketRepository) GetCandle(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, openTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

func (m *MockMarketRepository) GetCandles(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, start, end, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Candle), args.Error(1)
}

func (m *MockMarketRepository) GetLatestCandle(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

func (m *MockMarketRepository) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

func (m *MockMarketRepository) GetLatestTickers(ctx context.Context, limit int) ([]*market.Ticker, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

func (m *MockMarketRepository) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
	args := m.Called(ctx, symbol, limit)
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

// GetOrderBook retrieves the order book for a symbol
func (m *MockMarketRepository) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error) {
	args := m.Called(ctx, symbol, exchange, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.OrderBook), args.Error(1)
}

type MockSymbolRepository struct {
	mock.Mock
}

func (m *MockSymbolRepository) Create(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *MockSymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Symbol), args.Error(1)
}

func (m *MockSymbolRepository) GetByExchange(ctx context.Context, exchange string) ([]*market.Symbol, error) {
	args := m.Called(ctx, exchange)
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *MockSymbolRepository) GetAll(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *MockSymbolRepository) Update(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *MockSymbolRepository) Delete(ctx context.Context, symbol string) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

type MockMarketCache struct {
	mock.Mock
}

func (m *MockMarketCache) CacheTicker(ticker *market.Ticker) {
	m.Called(ticker)
}

func (m *MockMarketCache) GetTicker(ctx context.Context, exchange, symbol string) (*market.Ticker, bool) {
	args := m.Called(ctx, exchange, symbol)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.Ticker), args.Bool(1)
}

func (m *MockMarketCache) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, bool) {
	args := m.Called(ctx, exchange)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).([]*market.Ticker), args.Bool(1)
}

func (m *MockMarketCache) GetLatestTickers(ctx context.Context) ([]*market.Ticker, bool) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).([]*market.Ticker), args.Bool(1)
}

func (m *MockMarketCache) CacheCandle(candle *market.Candle) {
	m.Called(candle)
}

func (m *MockMarketCache) CacheOrderBook(orderbook *market.OrderBook) {
	m.Called(orderbook)
}

func (m *MockMarketCache) GetCandle(ctx context.Context, exchange, symbol string, interval market.Interval, openTime time.Time) (*market.Candle, bool) {
	args := m.Called(ctx, exchange, symbol, interval, openTime)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.Candle), args.Bool(1)
}

func (m *MockMarketCache) GetLatestCandle(ctx context.Context, exchange, symbol string, interval market.Interval) (*market.Candle, bool) {
	args := m.Called(ctx, exchange, symbol, interval)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.Candle), args.Bool(1)
}

func (m *MockMarketCache) GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, bool) {
	args := m.Called(ctx, exchange, symbol)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*market.OrderBook), args.Bool(1)
}

func (m *MockMarketCache) Clear() {
	m.Called()
}

func (m *MockMarketCache) SetTickerExpiry(d time.Duration) {
	m.Called(d)
}

func (m *MockMarketCache) SetCandleExpiry(d time.Duration) {
	m.Called(d)
}

func (m *MockMarketCache) SetOrderbookExpiry(d time.Duration) {
	m.Called(d)
}

func (m *MockMarketCache) StartCleanupTask(ctx context.Context, interval time.Duration) {
	m.Called(ctx, interval)
}

// Test setup helper
func setupMarketDataUseCase(
	marketRepo *MockMarketRepository,
	symbolRepo *MockSymbolRepository,
	cache *MockMarketCache,
) *usecase.MarketDataUseCase {
	// Create a logger with a proper writer
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	return usecase.NewMarketDataUseCase(marketRepo, symbolRepo, cache, &logger)
}

// Tests
func TestGetLatestTickers_FromCache(t *testing.T) {
	// Arrange
	ctx := context.Background()
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	cache := new(MockMarketCache)

	tickers := []*market.Ticker{
		{Symbol: "BTCUSDT", Price: 50000.0, Exchange: "mexc"},
		{Symbol: "ETHUSDT", Price: 3000.0, Exchange: "mexc"},
	}

	cache.On("GetLatestTickers", ctx).Return(tickers, true)

	uc := setupMarketDataUseCase(marketRepo, symbolRepo, cache)

	// Act
	result, err := uc.GetLatestTickers(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "BTCUSDT", result[0].Symbol)
	assert.Equal(t, "ETHUSDT", result[1].Symbol)
	cache.AssertExpectations(t)
	marketRepo.AssertNotCalled(t, "GetAllTickers")
}

func TestGetLatestTickers_FromDatabase(t *testing.T) {
	// Arrange
	ctx := context.Background()
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	cache := new(MockMarketCache)

	tickers := []*market.Ticker{
		{Symbol: "BTCUSDT", Price: 50000.0, Exchange: "mexc"},
		{Symbol: "ETHUSDT", Price: 3000.0, Exchange: "mexc"},
	}

	// Empty cache
	cache.On("GetLatestTickers", ctx).Return([]*market.Ticker{}, false)

	// Database will return data
	marketRepo.On("GetAllTickers", ctx, "mexc").Return(tickers, nil)

	// Cache will be updated
	cache.On("CacheTicker", mock.Anything).Return()

	uc := setupMarketDataUseCase(marketRepo, symbolRepo, cache)

	// Act
	result, err := uc.GetLatestTickers(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "BTCUSDT", result[0].Symbol)
	assert.Equal(t, "ETHUSDT", result[1].Symbol)
	cache.AssertExpectations(t)
	marketRepo.AssertExpectations(t)
}

func TestGetTicker_FromCache(t *testing.T) {
	// Arrange
	ctx := context.Background()
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	cache := new(MockMarketCache)

	ticker := &market.Ticker{Symbol: "BTCUSDT", Price: 50000.0, Exchange: "mexc"}

	cache.On("GetTicker", ctx, "mexc", "BTCUSDT").Return(ticker, true)

	uc := setupMarketDataUseCase(marketRepo, symbolRepo, cache)

	// Act
	result, err := uc.GetTicker(ctx, "mexc", "BTCUSDT")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, 50000.0, result.Price)
	cache.AssertExpectations(t)
	marketRepo.AssertNotCalled(t, "GetTicker")
}

func TestGetTicker_FromDatabase(t *testing.T) {
	// Arrange
	ctx := context.Background()
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	cache := new(MockMarketCache)

	ticker := &market.Ticker{Symbol: "BTCUSDT", Price: 50000.0, Exchange: "mexc"}

	// Cache miss
	cache.On("GetTicker", ctx, "mexc", "BTCUSDT").Return((*market.Ticker)(nil), false)

	// Database will return data
	marketRepo.On("GetTicker", ctx, "BTCUSDT", "mexc").Return(ticker, nil)

	// Cache will be updated
	cache.On("CacheTicker", mock.Anything).Return()

	uc := setupMarketDataUseCase(marketRepo, symbolRepo, cache)

	// Act
	result, err := uc.GetTicker(ctx, "mexc", "BTCUSDT")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, 50000.0, result.Price)
	cache.AssertExpectations(t)
	marketRepo.AssertExpectations(t)
}

func TestGetCandles(t *testing.T) {
	// Arrange
	ctx := context.Background()
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	cache := new(MockMarketCache)

	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()

	candles := []*market.Candle{
		{
			Symbol:   "BTCUSDT",
			Exchange: "mexc",
			Interval: market.Interval1h,
			OpenTime: startTime,
			Open:     50000.0,
			High:     51000.0,
			Low:      49000.0,
			Close:    50500.0,
			Volume:   100.0,
		},
	}

	// No cache hit for recent data
	cache.On("GetCandle", ctx, "mexc", "BTCUSDT", market.Interval1h, mock.Anything).Return(nil, false)

	// Database will return data
	marketRepo.On("GetCandles", ctx, "BTCUSDT", "mexc", market.Interval1h, startTime, endTime, 10).Return(candles, nil)

	// Cache will be updated
	cache.On("CacheCandle", mock.Anything).Return()

	uc := setupMarketDataUseCase(marketRepo, symbolRepo, cache)

	// Act
	result, err := uc.GetCandles(ctx, "mexc", "BTCUSDT", market.Interval1h, startTime, endTime, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "BTCUSDT", result[0].Symbol)
	assert.Equal(t, 50500.0, result[0].Close)
	marketRepo.AssertExpectations(t)
}

func TestGetAllSymbols(t *testing.T) {
	// Arrange
	ctx := context.Background()
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	cache := new(MockMarketCache)

	symbols := []*market.Symbol{
		{
			Symbol:     "BTCUSDT",
			BaseAsset:  "BTC",
			QuoteAsset: "USDT",
			Status:     "TRADING",
			Exchange:   "mexc",
		},
		{
			Symbol:     "ETHUSDT",
			BaseAsset:  "ETH",
			QuoteAsset: "USDT",
			Status:     "TRADING",
			Exchange:   "mexc",
		},
	}

	symbolRepo.On("GetAll", ctx).Return(symbols, nil)

	uc := setupMarketDataUseCase(marketRepo, symbolRepo, cache)

	// Act
	result, err := uc.GetAllSymbols(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "BTCUSDT", result[0].Symbol)
	assert.Equal(t, "ETHUSDT", result[1].Symbol)
	symbolRepo.AssertExpectations(t)
}

func TestGetSymbolInfo(t *testing.T) {
	// Arrange
	ctx := context.Background()
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	cache := new(MockMarketCache)

	symbol := &market.Symbol{
		Symbol:     "BTCUSDT",
		BaseAsset:  "BTC",
		QuoteAsset: "USDT",
		Status:     "TRADING",
		Exchange:   "mexc",
	}

	symbolRepo.On("GetBySymbol", ctx, "BTCUSDT").Return(symbol, nil)

	uc := setupMarketDataUseCase(marketRepo, symbolRepo, cache)

	// Act
	result, err := uc.GetSymbolInfo(ctx, "BTCUSDT")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, "BTC", result.BaseAsset)
	assert.Equal(t, "USDT", result.QuoteAsset)
	symbolRepo.AssertExpectations(t)
}
