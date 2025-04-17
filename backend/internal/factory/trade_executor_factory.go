package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/trade"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// TradeExecutorFactory creates trade executor components
type TradeExecutorFactory struct {
	config *config.Config
	logger *zerolog.Logger
}

// NewTradeExecutorFactory creates a new trade executor factory
func NewTradeExecutorFactory(config *config.Config, logger *zerolog.Logger) *TradeExecutorFactory {
	return &TradeExecutorFactory{
		config: config,
		logger: logger,
	}
}

// CreateTradeExecutor creates a new trade executor with rate limiting and error handling
func (f *TradeExecutorFactory) CreateTradeExecutor(tradeService port.TradeService) port.TradeExecutor {
	// Create logger for the executor
	executorLogger := f.logger.With().Str("component", "trade_executor").Logger()
	
	// Get default configuration
	executorConfig := trade.DefaultExecutorConfig()
	
	// Override with config values if available
	if f.config.RateLimit.Enabled {
		executorConfig.RequestsPerSecond = float64(f.config.RateLimit.DefaultLimit) / 60.0
		executorConfig.BurstSize = f.config.RateLimit.DefaultBurst
	}
	
	// Create and return the executor
	return trade.NewRateLimitedExecutor(tradeService, &executorLogger, executorConfig)
}
