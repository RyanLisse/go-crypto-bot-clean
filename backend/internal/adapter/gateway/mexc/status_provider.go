package mexc

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MEXCStatusProvider provides status information for the MEXC API
type MEXCStatusProvider struct {
	client        port.MEXCClient
	logger        *zerolog.Logger
	lastCheckTime time.Time
	isRunning     bool
	name          string
}

// NewMEXCStatusProvider creates a new MEXC status provider
func NewMEXCStatusProvider(client port.MEXCClient, logger *zerolog.Logger) *MEXCStatusProvider {
	return &MEXCStatusProvider{
		client:    client,
		logger:    logger,
		isRunning: true,
		name:      "mexc_api",
	}
}

// GetStatus returns the current status of the MEXC API
func (p *MEXCStatusProvider) GetStatus(ctx context.Context) (*status.ComponentStatus, error) {
	componentStatus := status.NewComponentStatus(p.name, status.StatusUnknown)
	now := time.Now()
	p.lastCheckTime = now

	// Check if the client is available by making a simple API call
	// We use the ping endpoint or a simple market data request that doesn't require authentication
	exchangeInfo, err := p.client.GetExchangeInfo(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("MEXC API health check failed")
		componentStatus.Status = status.StatusError
		componentStatus.Message = "API connection failed"
		componentStatus.LastError = err.Error()
		p.isRunning = false
		return componentStatus, nil
	}

	// API is responsive
	p.isRunning = true
	componentStatus.Status = status.StatusRunning
	componentStatus.Message = "API connection is healthy"

	// Add some metrics
	componentStatus.AddMetric("symbols_count", len(exchangeInfo.Symbols))
	componentStatus.AddMetric("last_check_time", now.Format(time.RFC3339))
	componentStatus.AddMetric("response_time_ms", time.Since(now).Milliseconds())

	// Add rate limit info if available
	componentStatus.AddMetric("rate_limits_count", len(exchangeInfo.Symbols))

	return componentStatus, nil
}

// GetName returns the name of the component
func (p *MEXCStatusProvider) GetName() string {
	return p.name
}

// IsRunning returns true if the component is running
func (p *MEXCStatusProvider) IsRunning() bool {
	return p.isRunning
}
