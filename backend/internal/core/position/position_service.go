package position

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PositionService defines the interface for position management
type PositionService interface {
	// GetPosition retrieves a position by ID
	GetPosition(ctx context.Context, positionID string) (*models.Position, error)

	// GetAllPositions retrieves all open positions
	GetAllPositions(ctx context.Context) ([]*models.Position, error)

	// OpenPosition opens a new position
	OpenPosition(ctx context.Context, symbol string, amount float64, entryPrice float64) (*models.Position, error)

	// ClosePosition closes an existing position
	ClosePosition(ctx context.Context, positionID string, exitPrice float64) (*models.ClosedPosition, error)

	// SetStopLoss sets a stop-loss level for a position
	SetStopLoss(ctx context.Context, positionID string, price float64) error

	// SetTakeProfit sets a take-profit level for a position
	SetTakeProfit(ctx context.Context, positionID string, price float64) error

	// GetPositionPnL calculates the current profit/loss for a position
	GetPositionPnL(ctx context.Context, positionID string) (float64, float64, error)
}

// PositionRepository defines the interface for position data access
type PositionRepository interface {
	GetByID(ctx context.Context, id string) (*models.Position, error)
	GetBySymbol(ctx context.Context, symbol string) (*models.Position, error)
	GetAll(ctx context.Context) ([]*models.Position, error)
	Save(ctx context.Context, position *models.Position) error
	Update(ctx context.Context, position *models.Position) error
	Delete(ctx context.Context, id string) error
}

// MarketService defines the interface for market data access
type MarketService interface {
	GetCurrentPrice(ctx context.Context, symbol string) (float64, error)
}

// positionService implements the PositionService interface
type positionService struct {
	positionRepo PositionRepository
	marketSvc    MarketService
}

// NewPositionService creates a new instance of PositionService
func NewPositionService(positionRepo PositionRepository, marketSvc MarketService) PositionService {
	return &positionService{
		positionRepo: positionRepo,
		marketSvc:    marketSvc,
	}
}

// GetPosition retrieves a position by ID
func (s *positionService) GetPosition(ctx context.Context, positionID string) (*models.Position, error) {
	return s.positionRepo.GetByID(ctx, positionID)
}

// GetAllPositions retrieves all open positions
func (s *positionService) GetAllPositions(ctx context.Context) ([]*models.Position, error) {
	return s.positionRepo.GetAll(ctx)
}

// OpenPosition opens a new position
func (s *positionService) OpenPosition(ctx context.Context, symbol string, amount float64, entryPrice float64) (*models.Position, error) {
	// Create a new position
	position := &models.Position{
		ID:         uuid.New().String(),
		Symbol:     symbol,
		Amount:     amount,
		EntryPrice: entryPrice,
		OpenTime:   time.Now(),
		StopLoss:   0.0,
		TakeProfit: 0.0,
	}

	// Save the position
	err := s.positionRepo.Save(ctx, position)
	if err != nil {
		return nil, err
	}

	return position, nil
}

// ClosePosition closes an existing position
func (s *positionService) ClosePosition(ctx context.Context, positionID string, exitPrice float64) (*models.ClosedPosition, error) {
	// Get the position
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return nil, err
	}

	if position == nil {
		return nil, errors.New("position not found")
	}

	// Calculate profit/loss
	profitLoss := (exitPrice - position.EntryPrice) * position.Amount
	profitLossPercentage := (exitPrice - position.EntryPrice) / position.EntryPrice

	// Create closed position record
	closedPosition := &models.ClosedPosition{
		ID:                   position.ID,
		Symbol:               position.Symbol,
		Amount:               position.Amount,
		EntryPrice:           position.EntryPrice,
		ExitPrice:            exitPrice,
		OpenTime:             position.OpenTime,
		CloseTime:            time.Now(),
		ProfitLoss:           profitLoss,
		ProfitLossPercentage: profitLossPercentage,
		ExitReason:           "manual", // Default reason
	}

	// Delete the open position
	err = s.positionRepo.Delete(ctx, positionID)
	if err != nil {
		return nil, err
	}

	return closedPosition, nil
}

// SetStopLoss sets a stop-loss level for a position
func (s *positionService) SetStopLoss(ctx context.Context, positionID string, price float64) error {
	// Get the position
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return err
	}

	if position == nil {
		return errors.New("position not found")
	}

	// Update stop-loss
	position.StopLoss = price

	// Save the updated position
	return s.positionRepo.Update(ctx, position)
}

// SetTakeProfit sets a take-profit level for a position
func (s *positionService) SetTakeProfit(ctx context.Context, positionID string, price float64) error {
	// Get the position
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return err
	}

	if position == nil {
		return errors.New("position not found")
	}

	// Update take-profit
	position.TakeProfit = price

	// Save the updated position
	return s.positionRepo.Update(ctx, position)
}

// GetPositionPnL calculates the current profit/loss for a position
func (s *positionService) GetPositionPnL(ctx context.Context, positionID string) (float64, float64, error) {
	// Get the position
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return 0, 0, err
	}

	if position == nil {
		return 0, 0, errors.New("position not found")
	}

	// Get current market price
	currentPrice, err := s.marketSvc.GetCurrentPrice(ctx, position.Symbol)
	if err != nil {
		return 0, 0, err
	}

	// Calculate profit/loss
	profitLoss := (currentPrice - position.EntryPrice) * position.Amount
	profitLossPercentage := (currentPrice - position.EntryPrice) / position.EntryPrice

	return profitLoss, profitLossPercentage, nil
}
