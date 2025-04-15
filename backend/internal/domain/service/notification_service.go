package service

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// NotificationService handles sending notifications based on system events
type NotificationService struct {
	logger   zerolog.Logger
	eventBus port.EventBus
}

// NewNotificationService creates a new NotificationService instance
func NewNotificationService(eventBus port.EventBus, logger zerolog.Logger) *NotificationService {
	service := &NotificationService{
		logger:   logger.With().Str("component", "NotificationService").Logger(),
		eventBus: eventBus,
	}
	// Subscribe to new coin events upon creation
	eventBus.Subscribe(service.HandleNewCoinEvent)
	service.logger.Info().Msg("NotificationService subscribed to NewCoinEvents")
	return service
}

// HandleNewCoinEvent processes new coin events and logs them (placeholder for actual notification logic)
func (s *NotificationService) HandleNewCoinEvent(event *model.NewCoinEvent) {
	s.logger.Info().
		Str("event_id", event.ID).
		Str("coin_id", event.CoinID).
		Str("event_type", event.EventType).
		Str("old_status", string(event.OldStatus)).
		Str("new_status", string(event.NewStatus)).
		Msg("Received NewCoinEvent")

	// TODO: Implement actual notification logic here (e.g., email, Slack, push notification)
	switch event.EventType {
	case "new_coin_detected":
		s.logger.Info().Str("coin_id", event.CoinID).Msgf("New coin detected! Status: %s", event.NewStatus)
	case "status_changed":
		s.logger.Info().Str("coin_id", event.CoinID).Msgf("Coin status changed from %s to %s", event.OldStatus, event.NewStatus)
		if event.NewStatus == model.StatusTrading {
			s.logger.Info().Str("coin_id", event.CoinID).Msg("Coin is now TRADING!")
		}
	}
}

// Stop unsubscribes the service from the event bus
func (s *NotificationService) Stop() {
	s.eventBus.Unsubscribe(s.HandleNewCoinEvent)
	s.logger.Info().Msg("NotificationService unsubscribed from NewCoinEvents")
}
