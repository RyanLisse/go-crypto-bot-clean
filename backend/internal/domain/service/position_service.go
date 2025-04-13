package service

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
)

// PositionService provides higher-level position operations
// and business logic for position management
type PositionService struct {
	positionUC    usecase.PositionUseCase
	marketService MarketDataServiceInterface
	logger        *zerolog.Logger
}

// NewPositionService creates a new position service
func NewPositionService(
	positionUC usecase.PositionUseCase,
	marketService MarketDataServiceInterface,
	logger *zerolog.Logger,
) *PositionService {
	return &PositionService{
		positionUC:    positionUC,
		marketService: marketService,
		logger:        logger,
	}
}

// CreatePosition creates a new position with validation
func (s *PositionService) CreatePosition(ctx context.Context, req model.PositionCreateRequest) (*model.Position, error) {
	// Log the request
	s.logger.Info().
		Str("symbol", req.Symbol).
		Str("side", string(req.Side)).
		Str("type", string(req.Type)).
		Float64("entryPrice", req.EntryPrice).
		Float64("quantity", req.Quantity).
		Msg("Creating new position")

	// Create the position
	position, err := s.positionUC.CreatePosition(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create position")
		return nil, err
	}

	s.logger.Info().
		Str("positionId", position.ID).
		Str("symbol", position.Symbol).
		Msg("Position created successfully")

	return position, nil
}

// UpdatePositionWithMarketData updates the position with the latest market data
func (s *PositionService) UpdatePositionWithMarketData(ctx context.Context, positionID string) (*model.Position, error) {
	// Get the position
	position, err := s.positionUC.GetPositionByID(ctx, positionID)
	if err != nil {
		s.logger.Error().Err(err).Str("positionId", positionID).Msg("Failed to get position")
		return nil, err
	}

	// Get latest ticker data
	ticker, err := s.marketService.RefreshTicker(ctx, position.Symbol)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", position.Symbol).Msg("Failed to get ticker data")
		return nil, err
	}

	// Update position with current price
	updatedPosition, err := s.positionUC.UpdatePositionPrice(ctx, positionID, ticker.Price)
	if err != nil {
		s.logger.Error().Err(err).Str("positionId", positionID).Msg("Failed to update position price")
		return nil, err
	}

	s.logger.Debug().
		Str("positionId", positionID).
		Float64("price", ticker.Price).
		Float64("pnl", updatedPosition.PnL).
		Float64("pnlPercent", updatedPosition.PnLPercent).
		Msg("Position updated with market data")

	return updatedPosition, nil
}

// ClosePosition closes a position with the current market price
func (s *PositionService) ClosePosition(ctx context.Context, positionID string, exitOrderIDs []string) (*model.Position, error) {
	// Get the position
	position, err := s.positionUC.GetPositionByID(ctx, positionID)
	if err != nil {
		s.logger.Error().Err(err).Str("positionId", positionID).Msg("Failed to get position")
		return nil, err
	}

	// Get latest ticker data
	ticker, err := s.marketService.RefreshTicker(ctx, position.Symbol)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", position.Symbol).Msg("Failed to get ticker data")
		return nil, err
	}

	// Close the position with the current price
	closedPosition, err := s.positionUC.ClosePosition(ctx, positionID, ticker.Price, exitOrderIDs)
	if err != nil {
		s.logger.Error().Err(err).Str("positionId", positionID).Msg("Failed to close position")
		return nil, err
	}

	s.logger.Info().
		Str("positionId", positionID).
		Float64("exitPrice", ticker.Price).
		Float64("pnl", closedPosition.PnL).
		Float64("pnlPercent", closedPosition.PnLPercent).
		Msg("Position closed successfully")

	return closedPosition, nil
}

// AnalyzePositionPerformance provides a detailed analysis of a position
func (s *PositionService) AnalyzePositionPerformance(ctx context.Context, positionID string) (map[string]interface{}, error) {
	// Get the position
	position, err := s.positionUC.GetPositionByID(ctx, positionID)
	if err != nil {
		s.logger.Error().Err(err).Str("positionId", positionID).Msg("Failed to get position")
		return nil, err
	}

	// Get historical market data for the position's duration
	var startTime, endTime time.Time
	startTime = position.OpenedAt

	if position.Status == model.PositionStatusClosed && position.ClosedAt != nil {
		endTime = *position.ClosedAt
	} else {
		endTime = time.Now()
	}

	tickers, err := s.marketService.GetHistoricalPrices(ctx, position.Symbol, startTime, endTime)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", position.Symbol).Msg("Failed to get historical ticker data")
		// Continue with partial analysis
	}

	// Perform analysis
	analysis := map[string]interface{}{
		"positionId":      position.ID,
		"symbol":          position.Symbol,
		"side":            position.Side,
		"type":            position.Type,
		"entryPrice":      position.EntryPrice,
		"currentPrice":    position.CurrentPrice,
		"quantity":        position.Quantity,
		"pnl":             position.PnL,
		"pnlPercent":      position.PnLPercent,
		"maxDrawdown":     position.MaxDrawdown,
		"maxProfit":       position.MaxProfit,
		"durationHours":   endTime.Sub(startTime).Hours(),
		"hasStopLoss":     position.StopLoss != nil,
		"hasTakeProfit":   position.TakeProfit != nil,
		"riskRewardRatio": position.RiskRewardRatio,
	}

	// Add market volatility if we have historical data
	if len(tickers) > 0 {
		var (
			highest    float64
			lowest     float64
			volatility float64
		)

		highest = tickers[0].Price
		lowest = tickers[0].Price

		for _, ticker := range tickers {
			if ticker.Price > highest {
				highest = ticker.Price
			}
			if ticker.Price < lowest {
				lowest = ticker.Price
			}
		}

		// Calculate volatility as percentage range
		if lowest > 0 {
			volatility = (highest - lowest) / lowest * 100
		}

		analysis["marketVolatility"] = volatility
		analysis["marketHighest"] = highest
		analysis["marketLowest"] = lowest
	}

	s.logger.Debug().
		Str("positionId", positionID).
		Msg("Position performance analysis completed")

	return analysis, nil
}

// GetOpenPositionsSummary returns a summary of all open positions
func (s *PositionService) GetOpenPositionsSummary(ctx context.Context) (map[string]interface{}, error) {
	// Get all open positions
	positions, err := s.positionUC.GetOpenPositions(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get open positions")
		return nil, err
	}

	var (
		totalPositions  = len(positions)
		totalValue      float64
		totalPnL        float64
		longPositions   int
		shortPositions  int
		positionsByType = make(map[string]int)
	)

	// Calculate summary statistics
	for _, position := range positions {
		totalValue += position.Quantity * position.EntryPrice
		totalPnL += position.PnL

		if position.Side == model.PositionSideLong {
			longPositions++
		} else {
			shortPositions++
		}

		positionsByType[string(position.Type)]++
	}

	summary := map[string]interface{}{
		"totalOpenPositions": totalPositions,
		"totalValue":         totalValue,
		"totalPnL":           totalPnL,
		"longPositions":      longPositions,
		"shortPositions":     shortPositions,
		"positionsByType":    positionsByType,
	}

	// Calculate average PnL if there are positions
	if totalPositions > 0 {
		summary["averagePnL"] = totalPnL / float64(totalPositions)
	}

	s.logger.Debug().
		Int("count", totalPositions).
		Float64("totalValue", totalValue).
		Float64("totalPnL", totalPnL).
		Msg("Generated open positions summary")

	return summary, nil
}
