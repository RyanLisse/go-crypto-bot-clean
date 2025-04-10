package services

import (
	"context"
	"errors"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"

	"github.com/google/uuid"
)

// PositionService handles business logic related to positions
type PositionService struct {
	positionRepository ports.PositionRepository
	orderService       *OrderService
}

// NewPositionService creates a new PositionService
func NewPositionService(positionRepository ports.PositionRepository, orderService *OrderService) *PositionService {
	return &PositionService{
		positionRepository: positionRepository,
		orderService:       orderService,
	}
}

// GetPositionByID retrieves a position by its ID
func (s *PositionService) GetPositionByID(ctx context.Context, id string) (*models.Position, error) {
	return s.positionRepository.GetByID(ctx, id)
}

// ListPositions retrieves positions based on status
func (s *PositionService) ListPositions(ctx context.Context, status models.PositionStatus) ([]*models.Position, error) {
	return s.positionRepository.List(ctx, status)
}

// OpenPosition opens a new position
func (s *PositionService) OpenPosition(ctx context.Context, symbol string, side models.PositionSide, quantity float64, price float64) (*models.Position, error) {
	// Check if there's already an open position for this symbol
	existingPosition, err := s.positionRepository.GetOpenPositionBySymbol(ctx, symbol)
	if err == nil && existingPosition != nil {
		return nil, errors.New("an open position already exists for this symbol")
	}

	position := &models.Position{
		ID:         uuid.New().String(),
		Symbol:     symbol,
		Side:       side,
		Quantity:   quantity,
		EntryPrice: price,
		Status:     models.PositionStatusOpen,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.positionRepository.Create(ctx, position); err != nil {
		return nil, err
	}

	return position, nil
}

// ClosePosition closes an open position
func (s *PositionService) ClosePosition(ctx context.Context, id string) (*models.Position, error) {
	position, err := s.positionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if position.Status != models.PositionStatusOpen {
		return nil, errors.New("only open positions can be closed")
	}

	// In a real implementation, we would create a market order to close the position
	// and update the position with the actual exit price from the order execution

	// For now, we'll just update the position status
	position.Status = models.PositionStatusClosed
	position.ExitPrice = 0 // This would be set to the actual exit price
	position.UpdatedAt = time.Now()

	if err := s.positionRepository.Update(ctx, position); err != nil {
		return nil, err
	}

	return position, nil
}
