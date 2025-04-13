package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) Create(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *MockPositionRepository) GetByID(ctx context.Context, id string) (*model.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionRepository) Update(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *MockPositionRepository) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionRepository) GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionRepository) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	args := m.Called(ctx, positionType)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionRepository) GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, from, to, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPositionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Test setup helper
func setupPositionUseCase(
	positionRepo *MockPositionRepository,
	marketRepo *MockMarketRepository,
	symbolRepo *MockSymbolRepository,
) usecase.PositionUseCase {
	// Create a null logger for testing
	logger := zerolog.Nop()
	return usecase.NewPositionUseCase(positionRepo, marketRepo, symbolRepo, logger)
}

// Tests
func TestCreatePosition(t *testing.T) {
	// Setup
	positionRepo := new(MockPositionRepository)
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Valid symbol
	symbol := &market.Symbol{
		Symbol:     "BTCUSDT",
		BaseAsset:  "BTC",
		QuoteAsset: "USDT",
	}

	// Input data
	createReq := model.PositionCreateRequest{
		Symbol:     "BTCUSDT",
		Side:       model.PositionSideLong,
		Type:       model.PositionTypeManual,
		EntryPrice: 50000.0,
		Quantity:   0.1,
		OrderIDs:   []string{"order1"},
		Notes:      "Test position",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		symbolRepo.On("GetBySymbol", ctx, "BTCUSDT").Return(symbol, nil).Once()
		positionRepo.On("Create", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.Symbol == "BTCUSDT" && p.Side == model.PositionSideLong
		})).Return(nil).Once()

		// Execute
		position, err := positionUC.CreatePosition(ctx, createReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, position)
		assert.Equal(t, "BTCUSDT", position.Symbol)
		assert.Equal(t, model.PositionSideLong, position.Side)
		assert.Equal(t, model.PositionTypeManual, position.Type)
		assert.Equal(t, 50000.0, position.EntryPrice)
		assert.Equal(t, 0.1, position.Quantity)
		assert.Equal(t, model.PositionStatusOpen, position.Status)
		assert.Equal(t, []string{"order1"}, position.EntryOrderIDs)

		// Verify mocks
		symbolRepo.AssertExpectations(t)
		positionRepo.AssertExpectations(t)
	})

	t.Run("Symbol Not Found", func(t *testing.T) {
		// Setup expectations
		symbolRepo.On("GetBySymbol", ctx, "BTCUSDT").Return(nil, nil).Once()

		// Execute
		position, err := positionUC.CreatePosition(ctx, createReq)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrSymbolNotFound, err)
		assert.Nil(t, position)

		// Verify mocks
		symbolRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		// Setup expectations
		symbolRepo.On("GetBySymbol", ctx, "BTCUSDT").Return(symbol, nil).Once()
		positionRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error")).Once()

		// Execute
		position, err := positionUC.CreatePosition(ctx, createReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, position)

		// Verify mocks
		symbolRepo.AssertExpectations(t)
		positionRepo.AssertExpectations(t)
	})
}

func TestGetPositionByID(t *testing.T) {
	// Setup
	positionRepo := new(MockPositionRepository)
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Test data
	position := &model.Position{
		ID:           "pos1",
		Symbol:       "BTCUSDT",
		Side:         model.PositionSideLong,
		Status:       model.PositionStatusOpen,
		EntryPrice:   50000.0,
		Quantity:     0.1,
		CurrentPrice: 55000.0,
		PnL:          500.0,
	}

	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos1").Return(position, nil).Once()

		// Execute
		result, err := positionUC.GetPositionByID(ctx, "pos1")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, position, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found")).Once()

		// Execute
		result, err := positionUC.GetPositionByID(ctx, "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})
}

func TestUpdatePositionPrice(t *testing.T) {
	// Setup
	positionRepo := new(MockPositionRepository)
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Test position
	now := time.Now()
	position := &model.Position{
		ID:            "pos1",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusOpen,
		EntryPrice:    50000.0,
		Quantity:      0.1,
		CurrentPrice:  50000.0,
		PnL:           0.0,
		PnLPercent:    0.0,
		OpenedAt:      now,
		LastUpdatedAt: now,
	}

	t.Run("Update Long Position Price", func(t *testing.T) {
		// Clone position to avoid test interference
		pos := *position
		testPos := &pos

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos1").Return(testPos, nil).Once()
		positionRepo.On("Update", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.ID == "pos1" && p.CurrentPrice == 55000.0
		})).Return(nil).Once()

		// Execute
		result, err := positionUC.UpdatePositionPrice(ctx, "pos1", 55000.0)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 55000.0, result.CurrentPrice)
		assert.Equal(t, 500.0, result.PnL)       // (55000-50000)*0.1
		assert.Equal(t, 10.0, result.PnLPercent) // (55000-50000)/50000*100

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Update Short Position Price", func(t *testing.T) {
		// Create short position
		shortPos := &model.Position{
			ID:            "pos2",
			Symbol:        "BTCUSDT",
			Side:          model.PositionSideShort,
			Status:        model.PositionStatusOpen,
			EntryPrice:    50000.0,
			Quantity:      0.1,
			CurrentPrice:  50000.0,
			PnL:           0.0,
			PnLPercent:    0.0,
			OpenedAt:      now,
			LastUpdatedAt: now,
		}

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos2").Return(shortPos, nil).Once()
		positionRepo.On("Update", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.ID == "pos2" && p.CurrentPrice == 45000.0
		})).Return(nil).Once()

		// Execute
		result, err := positionUC.UpdatePositionPrice(ctx, "pos2", 45000.0)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 45000.0, result.CurrentPrice)
		assert.Equal(t, 500.0, result.PnL)       // (50000-45000)*0.1
		assert.Equal(t, 10.0, result.PnLPercent) // (50000-45000)/50000*100

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Position Not Found", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found")).Once()

		// Execute
		result, err := positionUC.UpdatePositionPrice(ctx, "nonexistent", 55000.0)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})
}

func TestClosePosition(t *testing.T) {
	// Setup
	positionRepo := new(MockPositionRepository)
	marketRepo := new(MockMarketRepository)
	symbolRepo := new(MockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Test position
	now := time.Now()
	position := &model.Position{
		ID:            "pos1",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusOpen,
		EntryPrice:    50000.0,
		Quantity:      0.1,
		CurrentPrice:  50000.0,
		PnL:           0.0,
		PnLPercent:    0.0,
		OpenedAt:      now,
		LastUpdatedAt: now,
	}

	t.Run("Close Position Successfully", func(t *testing.T) {
		// Clone position to avoid test interference
		pos := *position
		testPos := &pos

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos1").Return(testPos, nil).Once()
		positionRepo.On("Update", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.ID == "pos1" &&
				p.Status == model.PositionStatusClosed &&
				p.CurrentPrice == 55000.0 &&
				p.ExitOrderIDs[0] == "exit1"
		})).Return(nil).Once()

		// Execute
		result, err := positionUC.ClosePosition(ctx, "pos1", 55000.0, []string{"exit1"})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, model.PositionStatusClosed, result.Status)
		assert.Equal(t, 55000.0, result.CurrentPrice)
		assert.NotNil(t, result.ClosedAt)
		assert.Equal(t, []string{"exit1"}, result.ExitOrderIDs)
		assert.Equal(t, 500.0, result.PnL) // (55000-50000)*0.1

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Position Already Closed", func(t *testing.T) {
		// Create closed position
		closedTime := now.Add(-1 * time.Hour)
		closedPos := &model.Position{
			ID:            "pos2",
			Symbol:        "BTCUSDT",
			Side:          model.PositionSideLong,
			Status:        model.PositionStatusClosed,
			EntryPrice:    50000.0,
			Quantity:      0.1,
			CurrentPrice:  55000.0,
			PnL:           500.0,
			PnLPercent:    10.0,
			OpenedAt:      now.Add(-2 * time.Hour),
			ClosedAt:      &closedTime,
			LastUpdatedAt: closedTime,
			ExitOrderIDs:  []string{"existing"},
		}

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos2").Return(closedPos, nil).Once()

		// Execute - trying to close an already closed position
		result, err := positionUC.ClosePosition(ctx, "pos2", 60000.0, []string{"new"})

		// Assert
		assert.NoError(t, err) // This is not an error, just a no-op
		assert.NotNil(t, result)
		assert.Equal(t, model.PositionStatusClosed, result.Status)
		assert.Equal(t, 55000.0, result.CurrentPrice)              // Original price, not updated
		assert.Equal(t, []string{"existing"}, result.ExitOrderIDs) // Original orders not updated

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Position Not Found", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found")).Once()

		// Execute
		result, err := positionUC.ClosePosition(ctx, "nonexistent", 55000.0, []string{"exit1"})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})
}
