package notification

import (
	"context"

	"github.com/rs/zerolog"
)

// ConsoleNotificationService implements a simple console-based notification service
type ConsoleNotificationService struct {
	logger *zerolog.Logger
}

// NewConsoleNotificationService creates a new console notification service
func NewConsoleNotificationService(logger *zerolog.Logger) *ConsoleNotificationService {
	return &ConsoleNotificationService{
		logger: logger,
	}
}

// SendNotification sends a notification by logging it to the console
func (s *ConsoleNotificationService) SendNotification(ctx context.Context, userID, title, message string) error {
	s.logger.Info().
		Str("userID", userID).
		Str("title", title).
		Str("message", message).
		Msg("Notification sent")
	return nil
}

// SendAlert sends an alert by logging it to the console with warning level
func (s *ConsoleNotificationService) SendAlert(ctx context.Context, userID, title, message string) error {
	s.logger.Warn().
		Str("userID", userID).
		Str("title", title).
		Str("message", message).
		Msg("Alert sent")
	return nil
}

// SendError sends an error notification by logging it to the console with error level
func (s *ConsoleNotificationService) SendError(ctx context.Context, userID, title, message string) error {
	s.logger.Error().
		Str("userID", userID).
		Str("title", title).
		Str("message", message).
		Msg("Error notification sent")
	return nil
}
