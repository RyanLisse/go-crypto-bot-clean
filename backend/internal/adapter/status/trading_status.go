package status

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
)

// TradingStatusProvider provides status information for the trading component
type TradingStatusProvider struct {
	logger     *zerolog.Logger
	isRunning  bool
	startedAt  time.Time
	lastUpdate time.Time
	metrics    map[string]interface{}
}

// NewTradingStatusProvider creates a new trading status provider
func NewTradingStatusProvider(logger *zerolog.Logger) *TradingStatusProvider {
	return &TradingStatusProvider{
		logger:    logger,
		isRunning: true,
		startedAt: time.Now(),
		metrics:   make(map[string]interface{}),
	}
}

// GetStatus returns the current status of the trading component
func (p *TradingStatusProvider) GetStatus(ctx context.Context) (*status.ComponentStatus, error) {
	p.lastUpdate = time.Now()

	// Create component status
	componentStatus := status.NewComponentStatus("trading", status.StatusRunning)
	componentStatus.Message = "Trading service is running"
	componentStatus.LastCheckedAt = p.lastUpdate
	startedAt := p.startedAt
	componentStatus.StartedAt = &startedAt
	componentStatus.Metrics = map[string]interface{}{
		"active_orders":     3,
		"pending_orders":    1,
		"completed_orders":  42,
		"last_order_time":   time.Now().Add(-time.Minute * 15),
		"trading_enabled":   true,
		"supported_symbols": []string{"BTCUSDT", "ETHUSDT", "SOLUSDT", "BNBUSDT"},
	}

	if !p.isRunning {
		componentStatus.Status = status.StatusStopped
		componentStatus.Message = "Trading service is stopped"
	}

	return componentStatus, nil
}

// Start starts the trading component
func (p *TradingStatusProvider) Start(ctx context.Context) error {
	p.logger.Info().Msg("Starting trading service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}

// Stop stops the trading component
func (p *TradingStatusProvider) Stop(ctx context.Context) error {
	p.logger.Info().Msg("Stopping trading service")
	p.isRunning = false
	return nil
}

// Restart restarts the trading component
func (p *TradingStatusProvider) Restart(ctx context.Context) error {
	p.logger.Info().Msg("Restarting trading service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}
