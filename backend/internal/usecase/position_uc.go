package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Common errors
var (
	ErrPositionNotFound       = errors.New("position not found")
	ErrInvalidPositionData    = errors.New("invalid position data")
	ErrPositionSymbolNotFound = errors.New("symbol not found")
)

type PositionUseCase interface {
	// Create operations
	CreatePosition(ctx context.Context, req model.PositionCreateRequest) (*model.Position, error)

	// Read operations
	GetPositionByID(ctx context.Context, id string) (*model.Position, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error)
	GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error)
	GetPositionsBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error)
	GetOpenPositions(ctx context.Context) ([]*model.Position, error)
	GetClosedPositions(ctx context.Context, fromTime, toTime time.Time, limit, offset int) ([]*model.Position, error)
	GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error)

	// Update operations
	UpdatePosition(ctx context.Context, id string, req model.PositionUpdateRequest) (*model.Position, error)
	UpdatePositionPrice(ctx context.Context, id string, currentPrice float64) (*model.Position, error)
	ClosePosition(ctx context.Context, id string, exitPrice float64, exitOrderIDs []string) (*model.Position, error)
	SetStopLoss(ctx context.Context, id string, stopLoss float64) (*model.Position, error)
	SetTakeProfit(ctx context.Context, id string, takeProfit float64) (*model.Position, error)

	// Delete operations
	DeletePosition(ctx context.Context, id string) error
}

// Add implementation for GetPositionsBySymbol
func (uc *positionUseCase) GetPositionsBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	positions, err := uc.positionRepo.GetBySymbol(ctx, symbol, limit, offset)
	if err != nil {
		uc.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get positions by symbol")
		return nil, err
	}
	return positions, nil
}

type positionUseCase struct {
	positionRepo port.PositionRepository
	marketRepo   port.MarketRepository
	symbolRepo   port.SymbolRepository
	logger       zerolog.Logger
}

// Add implementation for GetByUserID
func (uc *positionUseCase) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error) {
	positions, err := uc.positionRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get positions by user ID")
		return nil, err
	}
	return positions, nil
}

// NewPositionUseCase creates a new PositionUseCase
func NewPositionUseCase(
	positionRepo port.PositionRepository,
	marketRepo port.MarketRepository,
	symbolRepo port.SymbolRepository,
	logger zerolog.Logger,
) PositionUseCase {
	return &positionUseCase{
		positionRepo: positionRepo,
		marketRepo:   marketRepo,
		symbolRepo:   symbolRepo,
		logger:       logger.With().Str("component", "position_usecase").Logger(),
	}
}

// CreatePosition creates a new position
func (uc *positionUseCase) CreatePosition(ctx context.Context, req model.PositionCreateRequest) (*model.Position, error) {
	// Validate symbol exists
	symbol, err := uc.symbolRepo.GetBySymbol(ctx, req.Symbol)
	if err != nil {
		uc.logger.Error().Err(err).Str("symbol", req.Symbol).Msg("Failed to validate symbol")
		return nil, err
	}
	if symbol == nil {
		uc.logger.Warn().Str("symbol", req.Symbol).Msg("Symbol not found")
		return nil, ErrSymbolNotFound
	}

	// Create position model
	position := &model.Position{
		ID:            uuid.New().String(),
		Symbol:        req.Symbol,
		Side:          req.Side,
		Status:        model.PositionStatusOpen,
		Type:          req.Type,
		EntryPrice:    req.EntryPrice,
		Quantity:      req.Quantity,
		CurrentPrice:  req.EntryPrice, // Initially set to entry price
		StopLoss:      req.StopLoss,
		TakeProfit:    req.TakeProfit,
		StrategyID:    req.StrategyID,
		EntryOrderIDs: req.OrderIDs,
		Notes:         req.Notes,
		OpenedAt:      time.Now(),
		LastUpdatedAt: time.Now(),
	}

	// Calculate initial PnL (will be 0 since current price = entry price)
	position.UpdateCurrentPrice(req.EntryPrice)

	// Calculate risk/reward ratio if stop-loss and take-profit are set
	position.CalculateRiskRewardRatio()

	// Save to repository
	err = uc.positionRepo.Create(ctx, position)
	if err != nil {
		uc.logger.Error().Err(err).Str("symbol", req.Symbol).Msg("Failed to create position")
		return nil, err
	}

	uc.logger.Info().
		Str("id", position.ID).
		Str("symbol", position.Symbol).
		Str("side", string(position.Side)).
		Float64("entry_price", position.EntryPrice).
		Float64("quantity", position.Quantity).
		Msg("Position created")

	return position, nil
}

// GetPositionByID retrieves a position by its ID
func (uc *positionUseCase) GetPositionByID(ctx context.Context, id string) (*model.Position, error) {
	position, err := uc.positionRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get position")
		return nil, err
	}
	return position, nil
}

// Add implementation for GetActiveByUser
func (uc *positionUseCase) GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error) {
	positions, err := uc.positionRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get active positions by user")
		return nil, err
	}
	return positions, nil
}

// GetOpenPositions retrieves all currently open positions
func (uc *positionUseCase) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	// Since we don't have a direct repository method to get all open positions,
	// we'll need to use a different approach

	// We can use Count to get the total number of positions, then fetch them all
	// with appropriate filtering
	filters := map[string]interface{}{
		"status": model.PositionStatusOpen,
	}

	total, err := uc.positionRepo.Count(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to get count of open positions")
		return nil, err
	}

	// If there are no open positions, return early
	if total == 0 {
		return []*model.Position{}, nil
	}

	// For simplicity in this implementation, we'll get positions by symbol
	// but without specifying a particular symbol, which should return all positions
	// Then we'll filter for open ones in memory
	// This is not optimal but a workaround for the missing repository method
	allPositions, err := uc.positionRepo.GetBySymbol(ctx, "", 0, int(total))
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to get positions")
		return nil, err
	}

	// Filter for open positions
	var openPositions []*model.Position
	for _, pos := range allPositions {
		if pos.Status == model.PositionStatusOpen {
			openPositions = append(openPositions, pos)
		}
	}

	uc.logger.Debug().Int("count", len(openPositions)).Msg("Retrieved open positions")
	return openPositions, nil
}

// GetClosedPositions retrieves closed positions within a given time range with pagination
func (uc *positionUseCase) GetClosedPositions(ctx context.Context, fromTime, toTime time.Time, limit, offset int) ([]*model.Position, error) {
	// Since we don't have a direct repository method for this,
	// we'll implement a similar approach as GetOpenPositions
	filters := map[string]interface{}{
		"status": model.PositionStatusClosed,
	}

	total, err := uc.positionRepo.Count(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to get count of closed positions")
		return nil, err
	}

	// If there are no closed positions, return early
	if total == 0 {
		return []*model.Position{}, nil
	}

	// Get all positions and filter in memory
	// This is not optimal but a workaround for the missing repository method
	allPositions, err := uc.positionRepo.GetBySymbol(ctx, "", offset, limit)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to get positions")
		return nil, err
	}

	// Filter for closed positions within the time range
	var closedPositions []*model.Position
	for _, pos := range allPositions {
		if pos.Status == model.PositionStatusClosed &&
			pos.ClosedAt != nil &&
			!pos.ClosedAt.Before(fromTime) &&
			!pos.ClosedAt.After(toTime) {
			closedPositions = append(closedPositions, pos)
		}
	}

	uc.logger.Debug().
		Int("count", len(closedPositions)).
		Time("fromTime", fromTime).
		Time("toTime", toTime).
		Msg("Retrieved closed positions")

	return closedPositions, nil
}

// GetOpenPositionsByType retrieves open positions of a specific type
func (uc *positionUseCase) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	// Get all open positions first
	openPositions, err := uc.GetOpenPositions(ctx)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("type", string(positionType)).
			Msg("Failed to get open positions for filtering by type")
		return nil, err
	}

	// Filter by position type
	var filteredPositions []*model.Position
	for _, pos := range openPositions {
		if pos.Type == positionType {
			filteredPositions = append(filteredPositions, pos)
		}
	}

	uc.logger.Debug().
		Int("count", len(filteredPositions)).
		Str("type", string(positionType)).
		Msg("Retrieved open positions by type")

	return filteredPositions, nil
}

// UpdatePosition updates a position based on the provided request
func (uc *positionUseCase) UpdatePosition(ctx context.Context, id string, req model.PositionUpdateRequest) (*model.Position, error) {
	// Get the current position
	position, err := uc.positionRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get position for update")
		return nil, err
	}

	// Update fields that are provided in the request
	if req.CurrentPrice != nil {
		position.UpdateCurrentPrice(*req.CurrentPrice)
	}

	if req.StopLoss != nil {
		position.StopLoss = req.StopLoss
		position.CalculateRiskRewardRatio()
	}

	if req.TakeProfit != nil {
		position.TakeProfit = req.TakeProfit
		position.CalculateRiskRewardRatio()
	}

	if req.Notes != nil {
		position.Notes = *req.Notes
	}

	if req.Status != nil {
		status := model.PositionStatus(*req.Status)
		if status == model.PositionStatusClosed && position.Status != model.PositionStatusClosed {
			// If closing the position and it's not already closed, update fields
			position.Status = model.PositionStatusClosed
			now := time.Now()
			position.ClosedAt = &now
		} else {
			position.Status = status
		}
	}

	if req.ClosedAt != nil {
		closedAt, err := time.Parse(time.RFC3339, *req.ClosedAt)
		if err != nil {
			uc.logger.Error().Err(err).Str("closedAt", *req.ClosedAt).Msg("Invalid closedAt time format")
			return nil, ErrInvalidPositionData
		}
		position.ClosedAt = &closedAt

		// If setting closed time, ensure status is CLOSED
		if position.Status != model.PositionStatusClosed {
			position.Status = model.PositionStatusClosed
		}
	}

	if req.ExitOrderIDs != nil {
		position.ExitOrderIDs = *req.ExitOrderIDs
	}

	// Update last updated timestamp
	position.LastUpdatedAt = time.Now()

	// Save the updated position
	err = uc.positionRepo.Update(ctx, position)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to update position")
		return nil, err
	}

	uc.logger.Info().
		Str("id", position.ID).
		Str("symbol", position.Symbol).
		Str("status", string(position.Status)).
		Float64("pnl", position.PnL).
		Msg("Position updated")

	return position, nil
}

// UpdatePositionPrice updates a position's current price and recalculates PnL
func (uc *positionUseCase) UpdatePositionPrice(ctx context.Context, id string, currentPrice float64) (*model.Position, error) {
	// Get the current position
	position, err := uc.positionRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get position for price update")
		return nil, err
	}

	// Update price and recalculate PnL
	position.UpdateCurrentPrice(currentPrice)

	// Save the updated position
	err = uc.positionRepo.Update(ctx, position)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to update position price")
		return nil, err
	}

	uc.logger.Debug().
		Str("id", position.ID).
		Float64("price", currentPrice).
		Float64("pnl", position.PnL).
		Float64("pnlPercent", position.PnLPercent).
		Msg("Position price updated")

	return position, nil
}

// ClosePosition closes a position with the specified exit price and order IDs
func (uc *positionUseCase) ClosePosition(ctx context.Context, id string, exitPrice float64, exitOrderIDs []string) (*model.Position, error) {
	// Get the current position
	position, err := uc.positionRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get position for closing")
		return nil, err
	}

	// Check if already closed
	if position.Status == model.PositionStatusClosed {
		uc.logger.Warn().Str("id", id).Msg("Position is already closed")
		return position, nil
	}

	// Close the position
	position.Close(exitPrice, exitOrderIDs)

	// Save the updated position
	err = uc.positionRepo.Update(ctx, position)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to close position")
		return nil, err
	}

	uc.logger.Info().
		Str("id", position.ID).
		Str("symbol", position.Symbol).
		Float64("exitPrice", exitPrice).
		Float64("pnl", position.PnL).
		Float64("pnlPercent", position.PnLPercent).
		Int("exitOrderIDsCount", len(exitOrderIDs)).
		Msg("Position closed")

	return position, nil
}

// SetStopLoss sets a stop-loss for a position
func (uc *positionUseCase) SetStopLoss(ctx context.Context, id string, stopLoss float64) (*model.Position, error) {
	// Get the current position
	position, err := uc.positionRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get position for setting stop-loss")
		return nil, err
	}

	// Check if position is closed
	if position.Status == model.PositionStatusClosed {
		uc.logger.Warn().Str("id", id).Msg("Cannot set stop-loss on closed position")
		return nil, errors.New("cannot set stop-loss on closed position")
	}

	// Validate stop-loss based on position side
	if position.Side == model.PositionSideLong && stopLoss >= position.EntryPrice {
		return nil, errors.New("stop-loss must be below entry price for long positions")
	}
	if position.Side == model.PositionSideShort && stopLoss <= position.EntryPrice {
		return nil, errors.New("stop-loss must be above entry price for short positions")
	}

	// Set stop-loss
	stopLossCopy := stopLoss
	position.StopLoss = &stopLossCopy
	position.LastUpdatedAt = time.Now()

	// Recalculate risk/reward ratio if take-profit is also set
	position.CalculateRiskRewardRatio()

	// Save the updated position
	err = uc.positionRepo.Update(ctx, position)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to set stop-loss")
		return nil, err
	}

	uc.logger.Info().
		Str("id", position.ID).
		Float64("stopLoss", stopLoss).
		Float64("entryPrice", position.EntryPrice).
		Str("side", string(position.Side)).
		Msg("Stop-loss set for position")

	return position, nil
}

// SetTakeProfit sets a take-profit for a position
func (uc *positionUseCase) SetTakeProfit(ctx context.Context, id string, takeProfit float64) (*model.Position, error) {
	// Get the current position
	position, err := uc.positionRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get position for setting take-profit")
		return nil, err
	}

	// Check if position is closed
	if position.Status == model.PositionStatusClosed {
		uc.logger.Warn().Str("id", id).Msg("Cannot set take-profit on closed position")
		return nil, errors.New("cannot set take-profit on closed position")
	}

	// Validate take-profit based on position side
	if position.Side == model.PositionSideLong && takeProfit <= position.EntryPrice {
		return nil, errors.New("take-profit must be above entry price for long positions")
	}
	if position.Side == model.PositionSideShort && takeProfit >= position.EntryPrice {
		return nil, errors.New("take-profit must be below entry price for short positions")
	}

	// Set take-profit
	takeProfitCopy := takeProfit
	position.TakeProfit = &takeProfitCopy
	position.LastUpdatedAt = time.Now()

	// Recalculate risk/reward ratio if stop-loss is also set
	position.CalculateRiskRewardRatio()

	// Save the updated position
	err = uc.positionRepo.Update(ctx, position)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to set take-profit")
		return nil, err
	}

	uc.logger.Info().
		Str("id", position.ID).
		Float64("takeProfit", takeProfit).
		Float64("entryPrice", position.EntryPrice).
		Str("side", string(position.Side)).
		Msg("Take-profit set for position")

	return position, nil
}

// DeletePosition deletes a position
func (uc *positionUseCase) DeletePosition(ctx context.Context, id string) error {
	// Check if the position exists
	_, err := uc.positionRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get position for deletion")
		return err
	}

	// Delete the position
	err = uc.positionRepo.Delete(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to delete position")
		return err
	}

	uc.logger.Info().Str("id", id).Msg("Position deleted")
	return nil
}
