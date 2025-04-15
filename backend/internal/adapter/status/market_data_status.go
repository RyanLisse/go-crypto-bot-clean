package status

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
)

// MarketDataStatusProvider provides status information for the market data component
type MarketDataStatusProvider struct {
	logger     *zerolog.Logger
	isRunning  bool
	startedAt  time.Time
	lastUpdate time.Time
	metrics    map[string]interface{}
}

// NewMarketDataStatusProvider creates a new market data status provider
func NewMarketDataStatusProvider(logger *zerolog.Logger) *MarketDataStatusProvider {
	return &MarketDataStatusProvider{
		logger:    logger,
		isRunning: true,
		startedAt: time.Now(),
		metrics:   make(map[string]interface{}),
	}
}

// GetStatus returns the current status of the market data component
func (p *MarketDataStatusProvider) GetStatus(ctx context.Context) (*status.ComponentStatus, error) {
	p.lastUpdate = time.Now()

	// Create component status
	componentStatus := status.NewComponentStatus("market_data", status.StatusRunning)
	componentStatus.Message = "Market data service is running"
	componentStatus.LastCheckedAt = p.lastUpdate
	startedAt := p.startedAt
	componentStatus.StartedAt = &startedAt
	componentStatus.Metrics = map[string]interface{}{
		"symbols_count":      125,
		"tickers_updated":    time.Now().Add(-time.Minute * 2),
		"orderbooks_updated": time.Now().Add(-time.Second * 30),
		"candles_updated":    time.Now().Add(-time.Minute * 5),
	}

	if !p.isRunning {
		componentStatus.Status = status.StatusStopped
		componentStatus.Message = "Market data service is stopped"
	}

	return componentStatus, nil
}

// Start starts the market data component
func (p *MarketDataStatusProvider) Start(ctx context.Context) error {
	p.logger.Info().Msg("Starting market data service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}

// Stop stops the market data component
func (p *MarketDataStatusProvider) Stop(ctx context.Context) error {
	p.logger.Info().Msg("Stopping market data service")
	p.isRunning = false
	return nil
}

// Restart restarts the market data component
func (p *MarketDataStatusProvider) Restart(ctx context.Context) error {
	p.logger.Info().Msg("Restarting market data service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}
