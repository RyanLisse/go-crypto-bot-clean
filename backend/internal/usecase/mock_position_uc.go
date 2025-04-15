package usecase

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// MockPositionUseCase is a mock implementation of the PositionUseCase interface
type MockPositionUseCase struct{}

// GetOpenPositions retrieves all open positions
func (m *MockPositionUseCase) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	return []*model.Position{
		{
			ID:     "mock-position-1",
			Symbol: "BTCUSDT",
			Side:   model.PositionSideLong,
			Status: model.PositionStatusOpen,
		},
	}, nil
}

// GetPositionByID retrieves a position by ID
func (m *MockPositionUseCase) GetPositionByID(ctx context.Context, positionID string) (*model.Position, error) {
	return &model.Position{
		ID:     positionID,
		Symbol: "BTCUSDT",
		Side:   model.PositionSideLong,
		Status: model.PositionStatusOpen,
	}, nil
}

// ClosePosition closes an open position
func (m *MockPositionUseCase) ClosePosition(ctx context.Context, positionID string, exitPrice float64, exitOrderIDs []string) (*model.Position, error) {
	return &model.Position{
		ID:     positionID,
		Symbol: "BTCUSDT",
		Side:   model.PositionSideLong,
		Status: model.PositionStatusClosed,
	}, nil
}

// SetStopLoss sets a stop-loss for a position
func (m *MockPositionUseCase) SetStopLoss(ctx context.Context, positionID string, stopLoss float64) (*model.Position, error) {
	return &model.Position{
		ID:       positionID,
		Symbol:   "BTCUSDT",
		Side:     model.PositionSideLong,
		Status:   model.PositionStatusOpen,
		StopLoss: &stopLoss,
	}, nil
}

// SetTakeProfit sets a take-profit for a position
func (m *MockPositionUseCase) SetTakeProfit(ctx context.Context, positionID string, takeProfit float64) (*model.Position, error) {
	return &model.Position{
		ID:         positionID,
		Symbol:     "BTCUSDT",
		Side:       model.PositionSideLong,
		Status:     model.PositionStatusOpen,
		TakeProfit: &takeProfit,
	}, nil
}

// CreatePosition creates a new position
func (m *MockPositionUseCase) CreatePosition(ctx context.Context, req model.PositionCreateRequest) (*model.Position, error) {
	return &model.Position{
		ID:         "mock-position-new",
		Symbol:     req.Symbol,
		Side:       req.Side,
		Status:     model.PositionStatusOpen,
		EntryPrice: req.EntryPrice,
		Quantity:   req.Quantity,
	}, nil
}

// GetByUserID retrieves positions for a user
func (m *MockPositionUseCase) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error) {
	return []*model.Position{
		{
			ID:     "mock-position-history-1",
			Symbol: "BTCUSDT",
			Side:   model.PositionSideLong,
			Status: model.PositionStatusClosed,
		},
	}, nil
}

// GetActiveByUser retrieves active positions for a user
func (m *MockPositionUseCase) GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error) {
	return []*model.Position{
		{
			ID:     "mock-active-position-1",
			Symbol: "BTCUSDT",
			Side:   model.PositionSideLong,
			Status: model.PositionStatusOpen,
		},
	}, nil
}

// GetPositionsBySymbol retrieves positions for a symbol
func (m *MockPositionUseCase) GetPositionsBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	return []*model.Position{
		{
			ID:     "mock-symbol-position-1",
			Symbol: symbol,
			Side:   model.PositionSideLong,
			Status: model.PositionStatusOpen,
		},
	}, nil
}

// GetClosedPositions retrieves closed positions within a time range
func (m *MockPositionUseCase) GetClosedPositions(ctx context.Context, fromTime, toTime time.Time, limit, offset int) ([]*model.Position, error) {
	return []*model.Position{
		{
			ID:     "mock-closed-position-1",
			Symbol: "BTCUSDT",
			Side:   model.PositionSideLong,
			Status: model.PositionStatusClosed,
		},
	}, nil
}

// GetOpenPositionsByType retrieves open positions by type
func (m *MockPositionUseCase) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	return []*model.Position{
		{
			ID:     "mock-type-position-1",
			Symbol: "BTCUSDT",
			Side:   model.PositionSideLong,
			Status: model.PositionStatusOpen,
			Type:   positionType,
		},
	}, nil
}

// UpdatePosition updates a position
func (m *MockPositionUseCase) UpdatePosition(ctx context.Context, id string, req model.PositionUpdateRequest) (*model.Position, error) {
	return &model.Position{
		ID:     id,
		Symbol: "BTCUSDT",
		Side:   model.PositionSideLong,
		Status: model.PositionStatusOpen,
	}, nil
}

// UpdatePositionPrice updates a position's current price
func (m *MockPositionUseCase) UpdatePositionPrice(ctx context.Context, id string, currentPrice float64) (*model.Position, error) {
	return &model.Position{
		ID:           id,
		Symbol:       "BTCUSDT",
		Side:         model.PositionSideLong,
		Status:       model.PositionStatusOpen,
		CurrentPrice: currentPrice,
	}, nil
}

// DeletePosition deletes a position
func (m *MockPositionUseCase) DeletePosition(ctx context.Context, id string) error {
	return nil
}

// Ensure MockPositionUseCase implements PositionUseCase
var _ PositionUseCase = (*MockPositionUseCase)(nil)
