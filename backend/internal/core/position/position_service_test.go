package position

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// Mock dependencies
type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) GetByID(ctx context.Context, id string) (*models.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionRepository) GetBySymbol(ctx context.Context, symbol string) (*models.Position, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionRepository) GetAll(ctx context.Context) ([]*models.Position, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Position), args.Error(1)
}

func (m *MockPositionRepository) Save(ctx context.Context, position *models.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *MockPositionRepository) Update(ctx context.Context, position *models.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *MockPositionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockMarketService struct {
	mock.Mock
}

func (m *MockMarketService) GetCurrentPrice(ctx context.Context, symbol string) (float64, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(float64), args.Error(1)
}

// Test cases
func TestGetPosition_Success(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	expectedPosition := &models.Position{
		ID:         "pos123",
		Symbol:     "BTC/USDT",
		Amount:     1.0,
		EntryPrice: 50000.0,
		OpenTime:   time.Now(),
	}
	
	repo.On("GetByID", mock.Anything, "pos123").Return(expectedPosition, nil)
	
	// Execute test
	position, err := service.GetPosition(context.Background(), "pos123")
	
	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedPosition, position)
	repo.AssertExpectations(t)
}

func TestGetPosition_NotFound(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, errors.New("position not found"))
	
	// Execute test
	position, err := service.GetPosition(context.Background(), "nonexistent")
	
	// Assert results
	assert.Error(t, err)
	assert.Nil(t, position)
	repo.AssertExpectations(t)
}

func TestGetAllPositions_Success(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	expectedPositions := []*models.Position{
		{
			ID:         "pos123",
			Symbol:     "BTC/USDT",
			Amount:     1.0,
			EntryPrice: 50000.0,
			OpenTime:   time.Now(),
		},
		{
			ID:         "pos456",
			Symbol:     "ETH/USDT",
			Amount:     10.0,
			EntryPrice: 3000.0,
			OpenTime:   time.Now(),
		},
	}
	
	repo.On("GetAll", mock.Anything).Return(expectedPositions, nil)
	
	// Execute test
	positions, err := service.GetAllPositions(context.Background())
	
	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedPositions, positions)
	repo.AssertExpectations(t)
}

func TestOpenPosition_Success(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	symbol := "BTC/USDT"
	amount := 1.0
	entryPrice := 50000.0
	
	repo.On("Save", mock.Anything, mock.MatchedBy(func(p *models.Position) bool {
		return p.Symbol == symbol && 
			   p.Amount == amount && 
			   p.EntryPrice == entryPrice &&
			   p.StopLoss == 0.0 &&
			   p.TakeProfit == 0.0 &&
			   p.ID != ""
	})).Return(nil)
	
	// Execute test
	position, err := service.OpenPosition(context.Background(), symbol, amount, entryPrice)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, position)
	assert.Equal(t, symbol, position.Symbol)
	assert.Equal(t, amount, position.Amount)
	assert.Equal(t, entryPrice, position.EntryPrice)
	assert.NotEmpty(t, position.ID)
	repo.AssertExpectations(t)
}

func TestClosePosition_Success(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	positionID := "pos123"
	exitPrice := 55000.0
	
	position := &models.Position{
		ID:         positionID,
		Symbol:     "BTC/USDT",
		Amount:     1.0,
		EntryPrice: 50000.0,
		OpenTime:   time.Now().Add(-24 * time.Hour),
	}
	
	repo.On("GetByID", mock.Anything, positionID).Return(position, nil)
	repo.On("Delete", mock.Anything, positionID).Return(nil)
	
	// Execute test
	closedPosition, err := service.ClosePosition(context.Background(), positionID, exitPrice)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, closedPosition)
	assert.Equal(t, position.Symbol, closedPosition.Symbol)
	assert.Equal(t, position.Amount, closedPosition.Amount)
	assert.Equal(t, position.EntryPrice, closedPosition.EntryPrice)
	assert.Equal(t, exitPrice, closedPosition.ExitPrice)
	assert.Equal(t, 5000.0, closedPosition.ProfitLoss) // 55000 - 50000 = 5000
	assert.Equal(t, 0.1, closedPosition.ProfitLossPercentage) // 10% gain
	repo.AssertExpectations(t)
}

func TestSetStopLoss_Success(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	positionID := "pos123"
	stopLossPrice := 45000.0
	
	position := &models.Position{
		ID:         positionID,
		Symbol:     "BTC/USDT",
		Amount:     1.0,
		EntryPrice: 50000.0,
		OpenTime:   time.Now(),
	}
	
	repo.On("GetByID", mock.Anything, positionID).Return(position, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *models.Position) bool {
		return p.ID == positionID && p.StopLoss == stopLossPrice
	})).Return(nil)
	
	// Execute test
	err := service.SetStopLoss(context.Background(), positionID, stopLossPrice)
	
	// Assert results
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSetTakeProfit_Success(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	positionID := "pos123"
	takeProfitPrice := 60000.0
	
	position := &models.Position{
		ID:         positionID,
		Symbol:     "BTC/USDT",
		Amount:     1.0,
		EntryPrice: 50000.0,
		OpenTime:   time.Now(),
	}
	
	repo.On("GetByID", mock.Anything, positionID).Return(position, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *models.Position) bool {
		return p.ID == positionID && p.TakeProfit == takeProfitPrice
	})).Return(nil)
	
	// Execute test
	err := service.SetTakeProfit(context.Background(), positionID, takeProfitPrice)
	
	// Assert results
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetPositionPnL_Success(t *testing.T) {
	// Setup mocks and service
	repo := new(MockPositionRepository)
	market := new(MockMarketService)
	service := NewPositionService(repo, market)
	
	// Setup expectations
	positionID := "pos123"
	
	position := &models.Position{
		ID:         positionID,
		Symbol:     "BTC/USDT",
		Amount:     1.0,
		EntryPrice: 50000.0,
		OpenTime:   time.Now(),
	}
	
	currentPrice := 55000.0
	
	repo.On("GetByID", mock.Anything, positionID).Return(position, nil)
	market.On("GetCurrentPrice", mock.Anything, position.Symbol).Return(currentPrice, nil)
	
	// Execute test
	pnlValue, pnlPercentage, err := service.GetPositionPnL(context.Background(), positionID)
	
	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, 5000.0, pnlValue) // 55000 - 50000 = 5000
	assert.Equal(t, 0.1, pnlPercentage) // 10% gain
	repo.AssertExpectations(t)
	market.AssertExpectations(t)
}
