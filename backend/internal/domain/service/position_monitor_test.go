package service

import (
	"context"
	"testing"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockPositionUseCase struct {
	mock.Mock
}

func (m *MockPositionUseCase) CreatePosition(ctx context.Context, req model.PositionCreateRequest) (*model.Position, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetPositionByID(ctx context.Context, id string) (*model.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	args := m.Called(ctx, positionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetPositionsBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, from, to, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) UpdatePosition(ctx context.Context, id string, req model.PositionUpdateRequest) (*model.Position, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) UpdatePositionPrice(ctx context.Context, id string, currentPrice float64) (*model.Position, error) {
	args := m.Called(ctx, id, currentPrice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) ClosePosition(ctx context.Context, id string, exitPrice float64, exitOrderIDs []string) (*model.Position, error) {
	args := m.Called(ctx, id, exitPrice, exitOrderIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) SetStopLoss(ctx context.Context, id string, stopLoss float64) (*model.Position, error) {
	args := m.Called(ctx, id, stopLoss)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) SetTakeProfit(ctx context.Context, id string, takeProfit float64) (*model.Position, error) {
	args := m.Called(ctx, id, takeProfit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) DeletePosition(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockMarketDataService struct {
	mock.Mock
}

func (m *MockMarketDataService) RefreshTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

func (m *MockMarketDataService) GetHistoricalPrices(ctx context.Context, symbol string, startTime, endTime time.Time) ([]market.Ticker, error) {
	args := m.Called(ctx, symbol, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]market.Ticker), args.Error(1)
}

type MockTradeUseCase struct {
	mock.Mock
}

func (m *MockTradeUseCase) PlaceOrder(ctx context.Context, req model.OrderRequest) (*model.Order, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// Test setup
func setupPositionMonitorTest() (*MockPositionUseCase, *MockMarketDataService, *MockTradeUseCase, *PositionMonitor) {
	positionUC := new(MockPositionUseCase)
	marketService := new(MockMarketDataService)
	tradeUC := new(MockTradeUseCase)
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create a monitor with the mock service
	monitor := NewPositionMonitor(positionUC, marketService, tradeUC, &logger)
	monitor.SetInterval(100 * time.Millisecond) // Short interval for testing

	return positionUC, marketService, tradeUC, monitor
}

// Tests
func TestIsStopLossTriggered(t *testing.T) {
	_, _, _, monitor := setupPositionMonitorTest()

	// Test cases for long positions
	t.Run("Long position - stop-loss triggered", func(t *testing.T) {
		stopLoss := 9000.0
		position := &model.Position{
			Side:     model.PositionSideLong,
			StopLoss: &stopLoss,
		}
		triggered := monitor.isStopLossTriggered(position, 8900.0) // Price below stop-loss
		assert.True(t, triggered)
	})

	t.Run("Long position - stop-loss not triggered", func(t *testing.T) {
		stopLoss := 9000.0
		position := &model.Position{
			Side:     model.PositionSideLong,
			StopLoss: &stopLoss,
		}
		triggered := monitor.isStopLossTriggered(position, 9100.0) // Price above stop-loss
		assert.False(t, triggered)
	})

	// Test cases for short positions
	t.Run("Short position - stop-loss triggered", func(t *testing.T) {
		stopLoss := 11000.0
		position := &model.Position{
			Side:     model.PositionSideShort,
			StopLoss: &stopLoss,
		}
		triggered := monitor.isStopLossTriggered(position, 11100.0) // Price above stop-loss
		assert.True(t, triggered)
	})

	t.Run("Short position - stop-loss not triggered", func(t *testing.T) {
		stopLoss := 11000.0
		position := &model.Position{
			Side:     model.PositionSideShort,
			StopLoss: &stopLoss,
		}
		triggered := monitor.isStopLossTriggered(position, 10900.0) // Price below stop-loss
		assert.False(t, triggered)
	})

	// Test case for no stop-loss
	t.Run("No stop-loss", func(t *testing.T) {
		position := &model.Position{
			Side: model.PositionSideLong,
			// No stop-loss set
		}
		triggered := monitor.isStopLossTriggered(position, 9000.0)
		assert.False(t, triggered)
	})
}

func TestIsTakeProfitTriggered(t *testing.T) {
	_, _, _, monitor := setupPositionMonitorTest()

	// Test cases for long positions
	t.Run("Long position - take-profit triggered", func(t *testing.T) {
		takeProfit := 11000.0
		position := &model.Position{
			Side:       model.PositionSideLong,
			TakeProfit: &takeProfit,
		}
		triggered := monitor.isTakeProfitTriggered(position, 11100.0) // Price above take-profit
		assert.True(t, triggered)
	})

	t.Run("Long position - take-profit not triggered", func(t *testing.T) {
		takeProfit := 11000.0
		position := &model.Position{
			Side:       model.PositionSideLong,
			TakeProfit: &takeProfit,
		}
		triggered := monitor.isTakeProfitTriggered(position, 10900.0) // Price below take-profit
		assert.False(t, triggered)
	})

	// Test cases for short positions
	t.Run("Short position - take-profit triggered", func(t *testing.T) {
		takeProfit := 9000.0
		position := &model.Position{
			Side:       model.PositionSideShort,
			TakeProfit: &takeProfit,
		}
		triggered := monitor.isTakeProfitTriggered(position, 8900.0) // Price below take-profit
		assert.True(t, triggered)
	})

	t.Run("Short position - take-profit not triggered", func(t *testing.T) {
		takeProfit := 9000.0
		position := &model.Position{
			Side:       model.PositionSideShort,
			TakeProfit: &takeProfit,
		}
		triggered := monitor.isTakeProfitTriggered(position, 9100.0) // Price above take-profit
		assert.False(t, triggered)
	})

	// Test case for no take-profit
	t.Run("No take-profit", func(t *testing.T) {
		position := &model.Position{
			Side: model.PositionSideLong,
			// No take-profit set
		}
		triggered := monitor.isTakeProfitTriggered(position, 11000.0)
		assert.False(t, triggered)
	})
}

func TestHandleStopLossTrigger(t *testing.T) {
	positionUC, _, tradeUC, monitor := setupPositionMonitorTest()
	ctx := context.Background()

	// Create test position
	stopLoss := 9000.0
	position := &model.Position{
		ID:       "pos1",
		Symbol:   "BTCUSDT",
		Side:     model.PositionSideLong,
		Quantity: 0.1,
		StopLoss: &stopLoss,
	}

	// Create test order
	order := &model.Order{
		ID: "order1",
	}

	// Setup expectations
	tradeUC.On("PlaceOrder", ctx, mock.MatchedBy(func(req model.OrderRequest) bool {
		return req.Symbol == "BTCUSDT" && req.Side == model.OrderSideSell && req.Type == model.OrderTypeMarket && req.Quantity == 0.1
	})).Return(order, nil).Once()

	positionUC.On("ClosePosition", ctx, "pos1", 8900.0, []string{"order1"}).Return(position, nil).Once()

	// Execute
	monitor.handleStopLossTrigger(ctx, position, 8900.0)

	// Verify
	tradeUC.AssertExpectations(t)
	positionUC.AssertExpectations(t)
}

func TestHandleTakeProfitTrigger(t *testing.T) {
	positionUC, _, tradeUC, monitor := setupPositionMonitorTest()
	ctx := context.Background()

	// Create test position
	takeProfit := 11000.0
	position := &model.Position{
		ID:         "pos1",
		Symbol:     "BTCUSDT",
		Side:       model.PositionSideLong,
		Quantity:   0.1,
		TakeProfit: &takeProfit,
	}

	// Create test order
	order := &model.Order{
		ID: "order1",
	}

	// Setup expectations
	tradeUC.On("PlaceOrder", ctx, mock.MatchedBy(func(req model.OrderRequest) bool {
		return req.Symbol == "BTCUSDT" && req.Side == model.OrderSideSell && req.Type == model.OrderTypeMarket && req.Quantity == 0.1
	})).Return(order, nil).Once()

	positionUC.On("ClosePosition", ctx, "pos1", 11100.0, []string{"order1"}).Return(position, nil).Once()

	// Execute
	monitor.handleTakeProfitTrigger(ctx, position, 11100.0)

	// Verify
	tradeUC.AssertExpectations(t)
	positionUC.AssertExpectations(t)
}

func TestCheckPosition(t *testing.T) {
	positionUC, marketService, tradeUC, monitor := setupPositionMonitorTest()
	ctx := context.Background()

	// Create test position with stop-loss
	stopLoss := 9000.0
	position := &model.Position{
		ID:       "pos1",
		Symbol:   "BTCUSDT",
		Side:     model.PositionSideLong,
		Quantity: 0.1,
		Status:   model.PositionStatusOpen,
		StopLoss: &stopLoss,
	}

	// Create test ticker
	ticker := &market.Ticker{
		Symbol: "BTCUSDT",
		Price:  8900.0, // Below stop-loss
	}

	// Create test order
	order := &model.Order{
		ID: "order1",
	}

	// Setup expectations
	positionUC.On("GetPositionByID", ctx, "pos1").Return(position, nil).Once()
	marketService.On("RefreshTicker", ctx, "BTCUSDT").Return(ticker, nil).Once()
	positionUC.On("UpdatePositionPrice", ctx, "pos1", 8900.0).Return(position, nil).Once()

	tradeUC.On("PlaceOrder", ctx, mock.MatchedBy(func(req model.OrderRequest) bool {
		return req.Symbol == "BTCUSDT" && req.Side == model.OrderSideSell && req.Type == model.OrderTypeMarket
	})).Return(order, nil).Once()

	positionUC.On("ClosePosition", ctx, "pos1", 8900.0, []string{"order1"}).Return(position, nil).Once()

	// Execute
	err := monitor.CheckPosition(ctx, "pos1")

	// Verify
	assert.NoError(t, err)
	positionUC.AssertExpectations(t)
	marketService.AssertExpectations(t)
	tradeUC.AssertExpectations(t)
}

func TestCheckPositions(t *testing.T) {
	positionUC, marketService, tradeUC, monitor := setupPositionMonitorTest()

	// Create test positions
	stopLoss1 := 9000.0
	position1 := &model.Position{
		ID:       "pos1",
		Symbol:   "BTCUSDT",
		Side:     model.PositionSideLong,
		Quantity: 0.1,
		Status:   model.PositionStatusOpen,
		StopLoss: &stopLoss1,
	}

	takeProfit2 := 11000.0
	position2 := &model.Position{
		ID:         "pos2",
		Symbol:     "ETHUSDT",
		Side:       model.PositionSideLong,
		Quantity:   1.0,
		Status:     model.PositionStatusOpen,
		TakeProfit: &takeProfit2,
	}

	positions := []*model.Position{position1, position2}

	// Create test tickers
	ticker1 := &market.Ticker{
		Symbol: "BTCUSDT",
		Price:  8900.0, // Below stop-loss, triggered
	}

	ticker2 := &market.Ticker{
		Symbol: "ETHUSDT",
		Price:  11100.0, // Above take-profit, triggered
	}

	// Create test order
	order := &model.Order{
		ID: "order1",
	}

	// Setup expectations with mock.Anything for context
	positionUC.On("GetOpenPositions", mock.Anything).Return(positions, nil).Once()

	marketService.On("RefreshTicker", mock.Anything, "BTCUSDT").Return(ticker1, nil).Once()
	positionUC.On("UpdatePositionPrice", mock.Anything, "pos1", 8900.0).Return(position1, nil).Once()

	// For stop-loss on position1
	tradeUC.On("PlaceOrder", mock.Anything, mock.MatchedBy(func(req model.OrderRequest) bool {
		return req.Symbol == "BTCUSDT" && req.Side == model.OrderSideSell && req.Type == model.OrderTypeMarket
	})).Return(order, nil).Once()
	positionUC.On("ClosePosition", mock.Anything, "pos1", 8900.0, []string{"order1"}).Return(position1, nil).Once()

	marketService.On("RefreshTicker", mock.Anything, "ETHUSDT").Return(ticker2, nil).Once()
	positionUC.On("UpdatePositionPrice", mock.Anything, "pos2", 11100.0).Return(position2, nil).Once()

	// For take-profit on position2
	tradeUC.On("PlaceOrder", mock.Anything, mock.MatchedBy(func(req model.OrderRequest) bool {
		return req.Symbol == "ETHUSDT" && req.Side == model.OrderSideSell && req.Type == model.OrderTypeMarket
	})).Return(order, nil).Once()
	positionUC.On("ClosePosition", mock.Anything, "pos2", 11100.0, []string{"order1"}).Return(position2, nil).Once()

	// Execute
	monitor.checkPositions()

	// Wait for async operations to complete
	time.Sleep(100 * time.Millisecond)

	// Verify
	positionUC.AssertExpectations(t)
	marketService.AssertExpectations(t)
	tradeUC.AssertExpectations(t)
}
