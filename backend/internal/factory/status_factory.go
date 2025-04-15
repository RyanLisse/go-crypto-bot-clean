package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/notification"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/status"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/system"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// StatusFactory creates status-related components
type StatusFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewStatusFactory creates a new StatusFactory
func NewStatusFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *StatusFactory {
	return &StatusFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateSystemInfoProvider creates a system info provider
func (f *StatusFactory) CreateSystemInfoProvider() port.SystemInfoProvider {
	return system.NewSystemInfoProvider(f.logger, "/")
}

// CreateStatusRepository creates a status repository
func (f *StatusFactory) CreateStatusRepository() port.SystemStatusRepository {
	return repo.NewStatusRepository(f.db, f.logger)
}

// CreateStatusNotifier creates a status notifier
func (f *StatusFactory) CreateStatusNotifier() port.StatusNotifier {
	return notification.NewStatusNotifier(f.logger)
}

// CreateAlertNotifier creates an alert notifier
func (f *StatusFactory) CreateAlertNotifier() *notification.AlertNotifier {
	notifier := notification.NewAlertNotifier(f.logger, 100)

	// Configure email alerts if enabled
	if f.cfg.Notifications.Email.Enabled {
		emailConfig := notification.EmailConfig{
			Enabled:       f.cfg.Notifications.Email.Enabled,
			SMTPServer:    f.cfg.Notifications.Email.SMTPServer,
			SMTPPort:      f.cfg.Notifications.Email.SMTPPort,
			Username:      f.cfg.Notifications.Email.Username,
			Password:      f.cfg.Notifications.Email.Password,
			FromAddress:   f.cfg.Notifications.Email.FromAddress,
			ToAddresses:   f.cfg.Notifications.Email.ToAddresses,
			MinLevel:      notification.AlertLevel(f.cfg.Notifications.Email.MinLevel),
			SubjectPrefix: f.cfg.Notifications.Email.SubjectPrefix,
		}
		emailSubscriber := notification.NewEmailSubscriber(emailConfig, f.logger)
		notifier.AddSubscriber(emailSubscriber)
	}

	// Configure webhook alerts if enabled
	if f.cfg.Notifications.Webhook.Enabled {
		webhookConfig := notification.WebhookConfig{
			Enabled:   f.cfg.Notifications.Webhook.Enabled,
			URL:       f.cfg.Notifications.Webhook.URL,
			Method:    f.cfg.Notifications.Webhook.Method,
			Headers:   f.cfg.Notifications.Webhook.Headers,
			MinLevel:  notification.AlertLevel(f.cfg.Notifications.Webhook.MinLevel),
			Timeout:   f.cfg.Notifications.Webhook.Timeout,
			BatchSize: f.cfg.Notifications.Webhook.BatchSize,
		}
		webhookSubscriber := notification.NewWebhookSubscriber(webhookConfig, f.logger)
		notifier.AddSubscriber(webhookSubscriber)
	}

	return notifier
}

// CreateAlertHandler creates an alert handler
func (f *StatusFactory) CreateAlertHandler() *handler.AlertHandler {
	alertNotifier := f.CreateAlertNotifier()
	return handler.NewAlertHandler(alertNotifier, f.logger)
}

// CreateStatusUseCase creates a status use case
func (f *StatusFactory) CreateStatusUseCase() usecase.StatusUseCase {
	systemInfo := f.CreateSystemInfoProvider()
	statusRepo := f.CreateStatusRepository()
	notifier := f.CreateStatusNotifier()

	config := usecase.StatusUseCaseConfig{
		Version:        f.cfg.Version,
		UpdateInterval: 30, // 30 seconds
	}

	return usecase.NewStatusUseCase(systemInfo, statusRepo, notifier, f.logger, config)
}

// CreateStatusHandler creates a status handler
func (f *StatusFactory) CreateStatusHandler() *handler.StatusHandler {
	statusUseCase := f.CreateStatusUseCase()
	return handler.NewStatusHandler(statusUseCase, f.logger)
}

// RegisterStatusProviders registers status providers with the status use case
func (f *StatusFactory) RegisterStatusProviders(
	statusUseCase usecase.StatusUseCase,
	mexcFactory *MarketFactory,
) {
	// Register market data status provider
	marketDataProvider := status.NewMarketDataStatusProvider(f.logger)
	statusUseCase.RegisterProvider(marketDataProvider)

	// Register trading status provider
	tradingProvider := status.NewTradingStatusProvider(f.logger)
	statusUseCase.RegisterProvider(tradingProvider)

	// Register new coin detection status provider
	newCoinProvider := status.NewNewCoinStatusProvider(f.logger)
	statusUseCase.RegisterProvider(newCoinProvider)

	// Register risk management status provider
	riskProvider := status.NewRiskStatusProvider(f.logger)
	statusUseCase.RegisterProvider(riskProvider)

	// Register MEXC API status provider if available
	if mexcFactory != nil {
		mexcProvider := mexcFactory.CreateMEXCStatusProvider()
		statusUseCase.RegisterProvider(mexcProvider)
	}
}
