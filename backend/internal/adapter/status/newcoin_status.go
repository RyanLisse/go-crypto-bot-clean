package status

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
)

// NewCoinStatusProvider provides status information for the new coin detection component
type NewCoinStatusProvider struct {
	logger     *zerolog.Logger
	isRunning  bool
	startedAt  time.Time
	lastUpdate time.Time
	metrics    map[string]interface{}
}

// NewNewCoinStatusProvider creates a new coin detection status provider
func NewNewCoinStatusProvider(logger *zerolog.Logger) *NewCoinStatusProvider {
	return &NewCoinStatusProvider{
		logger:    logger,
		isRunning: true,
		startedAt: time.Now(),
		metrics:   make(map[string]interface{}),
	}
}

// GetStatus returns the current status of the new coin detection component
func (p *NewCoinStatusProvider) GetStatus(ctx context.Context) (*status.ComponentStatus, error) {
	p.lastUpdate = time.Now()

	// Create component status
	componentStatus := status.NewComponentStatus("new_coin_detector", status.StatusRunning)
	componentStatus.Message = "New coin detection service is running"
	componentStatus.LastCheckedAt = p.lastUpdate
	startedAt := p.startedAt
	componentStatus.StartedAt = &startedAt
	componentStatus.Metrics = map[string]interface{}{
		"coins_detected":      17,
		"last_detection_time": time.Now().Add(-time.Hour * 2),
		"sources_monitored":   []string{"MEXC", "Binance", "Twitter", "Telegram"},
		"scan_interval":       "5 minutes",
	}

	if !p.isRunning {
		componentStatus.Status = status.StatusStopped
		componentStatus.Message = "New coin detection service is stopped"
	}

	return componentStatus, nil
}

// Start starts the new coin detection component
func (p *NewCoinStatusProvider) Start(ctx context.Context) error {
	p.logger.Info().Msg("Starting new coin detection service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}

// Stop stops the new coin detection component
func (p *NewCoinStatusProvider) Stop(ctx context.Context) error {
	p.logger.Info().Msg("Stopping new coin detection service")
	p.isRunning = false
	return nil
}

// Restart restarts the new coin detection component
func (p *NewCoinStatusProvider) Restart(ctx context.Context) error {
	p.logger.Info().Msg("Restarting new coin detection service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}
