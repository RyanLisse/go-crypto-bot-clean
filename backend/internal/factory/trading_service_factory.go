package factory

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TradingServiceFactory creates the main trading service
type TradingServiceFactory struct {
	config *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewTradingServiceFactory creates a new trading service factory
func NewTradingServiceFactory(config *config.Config, logger *zerolog.Logger, db *gorm.DB) *TradingServiceFactory {
	return &TradingServiceFactory{
		config: config,
		logger: logger,
		db:     db,
	}
}

// CreateTradingService creates the main trading service with all dependencies
func (f *TradingServiceFactory) CreateTradingService(
	tradeExecutor port.TradeExecutor,
	tradeHistoryRepo port.TradeHistoryRepository,
) (*service.TradingService, error) {
	// Create logger for the service
	serviceLogger := f.logger.With().Str("component", "trading_service").Logger()
	
	// Create trade history factory
	tradeHistoryFactory := NewTradeHistoryFactory(f.config, &serviceLogger, f.db)
	
	// Create CSV writer
	csvWriter, err := tradeHistoryFactory.CreateTradeHistoryWriter()
	if err != nil {
		return nil, err
	}
	
	// Create telegram factory
	telegramFactory := NewTelegramFactory(f.config, &serviceLogger)
	
	// Create telegram notifier
	telegramNotifier := telegramFactory.CreateTelegramNotifier()
	
	// Create service config
	serviceConfig := service.TradingServiceConfig{
		ShutdownTimeout: 30 * time.Second,
		RecoveryEnabled: true,
	}
	
	// Create and return the service
	return service.NewTradingService(
		&serviceLogger,
		tradeExecutor,
		tradeHistoryRepo,
		csvWriter,
		telegramNotifier,
		serviceConfig,
	), nil
}
