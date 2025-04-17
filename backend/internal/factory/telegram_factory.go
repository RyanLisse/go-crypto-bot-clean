package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/notification"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
)

// TelegramFactory creates Telegram notification components
type TelegramFactory struct {
	config *config.Config
	logger *zerolog.Logger
}

// NewTelegramFactory creates a new Telegram factory
func NewTelegramFactory(config *config.Config, logger *zerolog.Logger) *TelegramFactory {
	return &TelegramFactory{
		config: config,
		logger: logger,
	}
}

// CreateTelegramNotifier creates a new Telegram notifier
func (f *TelegramFactory) CreateTelegramNotifier() *notification.TelegramNotifier {
	// Create logger for the notifier
	notifierLogger := f.logger.With().Str("component", "telegram_notifier").Logger()
	
	// Create and return the notifier
	return notification.NewTelegramNotifier(f.config.Telegram, &notifierLogger)
}
