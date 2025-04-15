package service

import (
	"context"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
)

// PositionMonitor monitors open positions for stop-loss and take-profit triggers
type PositionMonitor struct {
	positionUC    usecase.PositionUseCase
	marketService MarketDataServiceInterface
	tradeUC       usecase.TradeUseCase
	logger        *zerolog.Logger
	interval      time.Duration
	stopChan      chan struct{}
	wg            sync.WaitGroup
	running       bool
	mutex         sync.Mutex
}

// NewPositionMonitor creates a new PositionMonitor
func NewPositionMonitor(
	positionUC usecase.PositionUseCase,
	marketService MarketDataServiceInterface,
	tradeUC usecase.TradeUseCase,
	logger *zerolog.Logger,
) *PositionMonitor {
	return &PositionMonitor{
		positionUC:    positionUC,
		marketService: marketService,
		tradeUC:       tradeUC,
		logger:        logger,
		interval:      15 * time.Second, // Default check interval
		stopChan:      make(chan struct{}),
	}
}

// SetInterval sets the monitoring interval
func (m *PositionMonitor) SetInterval(interval time.Duration) {
	m.interval = interval
}

// Start starts the position monitor
func (m *PositionMonitor) Start() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.running {
		m.logger.Warn().Msg("Position monitor is already running")
		return
	}

	m.running = true
	m.stopChan = make(chan struct{})
	m.wg.Add(1)

	go m.monitorPositions()

	m.logger.Info().
		Dur("interval", m.interval).
		Msg("Position monitor started")
}

// Stop stops the position monitor
func (m *PositionMonitor) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.running {
		m.logger.Warn().Msg("Position monitor is not running")
		return
	}

	close(m.stopChan)
	m.wg.Wait()
	m.running = false

	m.logger.Info().Msg("Position monitor stopped")
}

// monitorPositions continuously monitors open positions
func (m *PositionMonitor) monitorPositions() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.checkPositions()
		}
	}
}

// checkPositions checks all open positions for stop-loss and take-profit triggers
func (m *PositionMonitor) checkPositions() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all open positions
	positions, err := m.positionUC.GetOpenPositions(ctx)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to get open positions")
		return
	}

	if len(positions) == 0 {
		return
	}

	m.logger.Debug().Int("count", len(positions)).Msg("Checking positions for SL/TP triggers")

	// Check each position
	for _, position := range positions {
		// Skip positions without stop-loss or take-profit
		if position.StopLoss == nil && position.TakeProfit == nil {
			continue
		}

		// Get current price for the symbol
		ticker, err := m.marketService.RefreshTicker(ctx, position.Symbol)
		if err != nil {
			m.logger.Error().
				Err(err).
				Str("symbol", position.Symbol).
				Str("positionId", position.ID).
				Msg("Failed to get current price")
			continue
		}

		if ticker == nil {
			m.logger.Warn().
				Str("symbol", position.Symbol).
				Str("positionId", position.ID).
				Msg("No ticker data available")
			continue
		}

		// Update position with current price
		updatedPosition, err := m.positionUC.UpdatePositionPrice(ctx, position.ID, ticker.Price)
		if err != nil {
			m.logger.Error().
				Err(err).
				Str("positionId", position.ID).
				Msg("Failed to update position price")
			continue
		}

		// Check for stop-loss trigger
		if m.isStopLossTriggered(updatedPosition, ticker.Price) {
			m.handleStopLossTrigger(ctx, updatedPosition, ticker.Price)
			continue
		}

		// Check for take-profit trigger
		if m.isTakeProfitTriggered(updatedPosition, ticker.Price) {
			m.handleTakeProfitTrigger(ctx, updatedPosition, ticker.Price)
			continue
		}
	}
}

// isStopLossTriggered checks if a position's stop-loss has been triggered
func (m *PositionMonitor) isStopLossTriggered(position *model.Position, currentPrice float64) bool {
	if position.StopLoss == nil {
		return false
	}

	if position.Side == model.PositionSideLong {
		// For long positions, stop-loss is triggered when price falls below stop-loss level
		return currentPrice <= *position.StopLoss
	} else {
		// For short positions, stop-loss is triggered when price rises above stop-loss level
		return currentPrice >= *position.StopLoss
	}
}

// isTakeProfitTriggered checks if a position's take-profit has been triggered
func (m *PositionMonitor) isTakeProfitTriggered(position *model.Position, currentPrice float64) bool {
	if position.TakeProfit == nil {
		return false
	}

	if position.Side == model.PositionSideLong {
		// For long positions, take-profit is triggered when price rises above take-profit level
		return currentPrice >= *position.TakeProfit
	} else {
		// For short positions, take-profit is triggered when price falls below take-profit level
		return currentPrice <= *position.TakeProfit
	}
}

// handleStopLossTrigger handles a triggered stop-loss
func (m *PositionMonitor) handleStopLossTrigger(ctx context.Context, position *model.Position, currentPrice float64) {
	m.logger.Info().
		Str("positionId", position.ID).
		Str("symbol", position.Symbol).
		Float64("stopLoss", *position.StopLoss).
		Float64("currentPrice", currentPrice).
		Msg("Stop-loss triggered")

	// Create a market order to close the position
	side := model.OrderSideSell
	if position.Side == model.PositionSideShort {
		side = model.OrderSideBuy
	}

	// Place a market order to close the position
	orderRequest := model.OrderRequest{
		Symbol:   position.Symbol,
		Side:     side,
		Type:     model.OrderTypeMarket,
		Quantity: position.Quantity,
	}

	order, err := m.tradeUC.PlaceOrder(ctx, orderRequest)
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("positionId", position.ID).
			Msg("Failed to place order for stop-loss")
		return
	}

	// Close the position
	_, err = m.positionUC.ClosePosition(ctx, position.ID, currentPrice, []string{order.ID})
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("positionId", position.ID).
			Msg("Failed to close position after stop-loss trigger")
		return
	}

	m.logger.Info().
		Str("positionId", position.ID).
		Str("symbol", position.Symbol).
		Float64("exitPrice", currentPrice).
		Float64("pnl", position.PnL).
		Msg("Position closed by stop-loss")
}

// handleTakeProfitTrigger handles a triggered take-profit
func (m *PositionMonitor) handleTakeProfitTrigger(ctx context.Context, position *model.Position, currentPrice float64) {
	m.logger.Info().
		Str("positionId", position.ID).
		Str("symbol", position.Symbol).
		Float64("takeProfit", *position.TakeProfit).
		Float64("currentPrice", currentPrice).
		Msg("Take-profit triggered")

	// Create a market order to close the position
	side := model.OrderSideSell
	if position.Side == model.PositionSideShort {
		side = model.OrderSideBuy
	}

	// Place a market order to close the position
	orderRequest := model.OrderRequest{
		Symbol:   position.Symbol,
		Side:     side,
		Type:     model.OrderTypeMarket,
		Quantity: position.Quantity,
	}

	order, err := m.tradeUC.PlaceOrder(ctx, orderRequest)
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("positionId", position.ID).
			Msg("Failed to place order for take-profit")
		return
	}

	// Close the position
	_, err = m.positionUC.ClosePosition(ctx, position.ID, currentPrice, []string{order.ID})
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("positionId", position.ID).
			Msg("Failed to close position after take-profit trigger")
		return
	}

	m.logger.Info().
		Str("positionId", position.ID).
		Str("symbol", position.Symbol).
		Float64("exitPrice", currentPrice).
		Float64("pnl", position.PnL).
		Msg("Position closed by take-profit")
}

// CheckPosition checks a specific position for stop-loss and take-profit triggers
// This can be called manually to check a position outside the regular monitoring cycle
func (m *PositionMonitor) CheckPosition(ctx context.Context, positionID string) error {
	// Get the position
	position, err := m.positionUC.GetPositionByID(ctx, positionID)
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("positionId", positionID).
			Msg("Failed to get position")
		return err
	}

	// Skip if position is closed or has no stop-loss/take-profit
	if position.Status == model.PositionStatusClosed || (position.StopLoss == nil && position.TakeProfit == nil) {
		return nil
	}

	// Get current price for the symbol
	ticker, err := m.marketService.RefreshTicker(ctx, position.Symbol)
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("symbol", position.Symbol).
			Str("positionId", position.ID).
			Msg("Failed to get current price")
		return err
	}

	if ticker == nil {
		m.logger.Warn().
			Str("symbol", position.Symbol).
			Str("positionId", position.ID).
			Msg("No ticker data available")
		return nil
	}

	// Update position with current price
	updatedPosition, err := m.positionUC.UpdatePositionPrice(ctx, position.ID, ticker.Price)
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("positionId", position.ID).
			Msg("Failed to update position price")
		return err
	}

	// Check for stop-loss trigger
	if m.isStopLossTriggered(updatedPosition, ticker.Price) {
		m.handleStopLossTrigger(ctx, updatedPosition, ticker.Price)
		return nil
	}

	// Check for take-profit trigger
	if m.isTakeProfitTriggered(updatedPosition, ticker.Price) {
		m.handleTakeProfitTrigger(ctx, updatedPosition, ticker.Price)
		return nil
	}

	return nil
}
