package management

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/domain/interfaces"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// MockPositionRepository is a mock implementation of the PositionRepository
type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) FindAll(ctx context.Context, filter interfaces.PositionFilter) ([]*models.Position, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*models.Position), args.Error(1)
}

func (m *MockPositionRepository) FindByID(ctx context.Context, id string) (*models.Position, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionRepository) FindBySymbol(ctx context.Context, symbol string) ([]*models.Position, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*models.Position), args.Error(1)
}

func (m *MockPositionRepository) Create(ctx context.Context, position *models.Position) (string, error) {
	args := m.Called(ctx, position)
	return args.String(0), args.Error(1)
}

func (m *MockPositionRepository) Update(ctx context.Context, position *models.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *MockPositionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPositionRepository) AddOrder(ctx context.Context, positionID string, order *models.Order) error {
	args := m.Called(ctx, positionID, order)
	return args.Error(0)
}

func (m *MockPositionRepository) UpdateOrder(ctx context.Context, positionID string, order *models.Order) error {
	args := m.Called(ctx, positionID, order)
	return args.Error(0)
}

// MockOrderService is a mock implementation of the OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) ExecuteOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) CancelOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockOrderService) GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*models.Order), args.Error(1)
}

// MockPriceService is a mock implementation of the PriceService
type MockPriceService struct {
	mock.Mock
}

func (m *MockPriceService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPriceService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*models.Ticker), args.Error(1)
}

func (m *MockPriceService) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	return args.Get(0).([]*models.Kline), args.Error(1)
}

func (m *MockPriceService) GetPriceHistory(ctx context.Context, symbol string, startTime, endTime time.Time) ([]float64, error) {
	args := m.Called(ctx, symbol, startTime, endTime)
	return args.Get(0).([]float64), args.Error(1)
}

func TestPositionManager_EnterPosition(t *testing.T) {
	// Create mocks
	positionRepo := new(MockPositionRepository)
	orderService := new(MockOrderService)
	priceService := new(MockPriceService)
	logger, _ := zap.NewDevelopment()

	// Create position manager
	pm := NewPositionManager(positionRepo, orderService, priceService, logger)

	// Setup test data
	ctx := context.Background()
	order := &models.Order{
		Symbol:   "BTCUSDT",
		Side:     "BUY",
		Type:     "MARKET",
		Quantity: 1.0,
		Price:    50000.0,
		Status:   "FILLED",
	}

	// Setup expectations
	positionRepo.On("Create", ctx, mock.AnythingOfType("*models.Position")).Return("pos_123", nil)

	// Execute test
	position, err := pm.EnterPosition(ctx, order)

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, position)
	assert.Equal(t, "pos_123", position.ID)
	assert.Equal(t, "BTCUSDT", position.Symbol)
	assert.Equal(t, 1.0, position.Quantity)
	assert.Equal(t, 50000.0, position.EntryPrice)
	assert.Equal(t, models.PositionStatusOpen, position.Status)
	assert.Len(t, position.Orders, 1)

	// Verify expectations
	positionRepo.AssertExpectations(t)
}

func TestPositionManager_ExitPosition(t *testing.T) {
	// Create mocks
	positionRepo := new(MockPositionRepository)
	orderService := new(MockOrderService)
	priceService := new(MockPriceService)
	logger, _ := zap.NewDevelopment()

	// Create position manager
	pm := NewPositionManager(positionRepo, orderService, priceService, logger)

	// Setup test data
	ctx := context.Background()
	positionID := "pos_123"
	position := &models.Position{
		ID:         positionID,
		Symbol:     "BTCUSDT",
		Quantity:   1.0,
		EntryPrice: 50000.0,
		Status:     "OPEN",
		Orders:     []models.Order{},
	}
	sellOrder := &models.Order{
		ID:        "order_456",
		Symbol:    "BTCUSDT",
		Side:      "SELL",
		Type:      "MARKET",
		Quantity:  1.0,
		Price:     55000.0,
		Status:    "FILLED",
		CreatedAt: time.Now(),
	}

	// Setup expectations
	positionRepo.On("FindByID", ctx, positionID).Return(position, nil)
	orderService.On("ExecuteOrder", ctx, mock.AnythingOfType("*models.Order")).Return(sellOrder, nil)
	positionRepo.On("AddOrder", ctx, positionID, sellOrder).Return(nil)
	positionRepo.On("Update", ctx, mock.AnythingOfType("*models.Position")).Return(nil)

	// Execute test
	err := pm.ExitPosition(ctx, positionID, 55000.0)

	// Verify results
	assert.NoError(t, err)

	// Verify expectations
	positionRepo.AssertExpectations(t)
	orderService.AssertExpectations(t)
}

func TestPositionManager_UpdateStopLoss(t *testing.T) {
	// Create mocks
	positionRepo := new(MockPositionRepository)
	orderService := new(MockOrderService)
	priceService := new(MockPriceService)
	logger, _ := zap.NewDevelopment()

	// Create position manager
	pm := NewPositionManager(positionRepo, orderService, priceService, logger)

	// Setup test data
	ctx := context.Background()
	positionID := "pos_123"
	position := &models.Position{
		ID:         positionID,
		Symbol:     "BTCUSDT",
		Quantity:   1.0,
		EntryPrice: 50000.0,
		StopLoss:   45000.0,
		Status:     "OPEN",
	}

	// Setup expectations
	positionRepo.On("FindByID", ctx, positionID).Return(position, nil)
	positionRepo.On("Update", ctx, mock.AnythingOfType("*models.Position")).Return(nil)

	// Execute test
	err := pm.UpdateStopLoss(ctx, positionID, 47000.0)

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, 47000.0, position.StopLoss)

	// Verify expectations
	positionRepo.AssertExpectations(t)
}
