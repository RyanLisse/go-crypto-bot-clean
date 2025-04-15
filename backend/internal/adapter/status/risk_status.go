package status

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
)

// RiskStatusProvider provides status information for the risk management component
type RiskStatusProvider struct {
	logger     *zerolog.Logger
	isRunning  bool
	startedAt  time.Time
	lastUpdate time.Time
	metrics    map[string]interface{}
}

// NewRiskStatusProvider creates a new risk management status provider
func NewRiskStatusProvider(logger *zerolog.Logger) *RiskStatusProvider {
	return &RiskStatusProvider{
		logger:    logger,
		isRunning: true,
		startedAt: time.Now(),
		metrics:   make(map[string]interface{}),
	}
}

// GetStatus returns the current status of the risk management component
func (p *RiskStatusProvider) GetStatus(ctx context.Context) (*status.ComponentStatus, error) {
	p.lastUpdate = time.Now()

	// Create component status
	componentStatus := status.NewComponentStatus("risk_management", status.StatusRunning)
	componentStatus.Message = "Risk management service is running"
	componentStatus.LastCheckedAt = p.lastUpdate
	startedAt := p.startedAt
	componentStatus.StartedAt = &startedAt
	componentStatus.Metrics = map[string]interface{}{
		"active_constraints":   8,
		"risk_checks_today":    124,
		"rejected_trades":      3,
		"current_risk_level":   "medium",
		"max_position_size":    1000.0,
		"max_drawdown":         "5%",
		"last_assessment_time": time.Now().Add(-time.Minute * 5),
	}

	if !p.isRunning {
		componentStatus.Status = status.StatusStopped
		componentStatus.Message = "Risk management service is stopped"
	}

	return componentStatus, nil
}

// Start starts the risk management component
func (p *RiskStatusProvider) Start(ctx context.Context) error {
	p.logger.Info().Msg("Starting risk management service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}

// Stop stops the risk management component
func (p *RiskStatusProvider) Stop(ctx context.Context) error {
	p.logger.Info().Msg("Stopping risk management service")
	p.isRunning = false
	return nil
}

// Restart restarts the risk management component
func (p *RiskStatusProvider) Restart(ctx context.Context) error {
	p.logger.Info().Msg("Restarting risk management service")
	p.isRunning = true
	p.startedAt = time.Now()
	return nil
}
