package ports

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/notification"
)

// Notifier defines the interface for a specific notification channel (e.g., email, slack).
type Notifier interface {
	// Send sends a notification.
	Send(ctx context.Context, recipient string, subject string, message string) error
	// Supports checks if this notifier handles the given channel type (e.g., "email", "slack").
	Supports(channel string) bool
}

// NotificationService defines the core interface for sending notifications.
type NotificationService interface {
	// SendNotification sends a message via all enabled notification channels for the given user.
	SendNotification(ctx context.Context, userID string, subject string, message string) error
}

// NotificationPreferenceRepository defines the interface for managing user notification preferences.
type NotificationPreferenceRepository interface {
	// GetUserPreferences retrieves all active notification preferences for a given user.
	// It should return only enabled preferences.
	GetUserPreferences(ctx context.Context, userID string) ([]notification.Preference, error)

	// TODO: Add methods for saving/updating preferences if needed by the application logic later.
	// SaveUserPreference(ctx context.Context, pref notification.Preference) error
	// DeleteUserPreference(ctx context.Context, userID, channel string) error
}
