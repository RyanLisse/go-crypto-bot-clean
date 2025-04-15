package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
)

// StatusNotifier implements the StatusNotifier interface
type StatusNotifier struct {
	logger           *zerolog.Logger
	notificationChan chan Notification
	enabled          bool
}

// Notification represents a status notification
type Notification struct {
	Type      string
	Component string
	OldStatus status.Status
	NewStatus status.Status
	Message   string
	Timestamp time.Time
}

// NewStatusNotifier creates a new status notifier
func NewStatusNotifier(logger *zerolog.Logger) *StatusNotifier {
	notifier := &StatusNotifier{
		logger:           logger,
		notificationChan: make(chan Notification, 100),
		enabled:          true,
	}

	// Start the notification processor
	go notifier.processNotifications()

	return notifier
}

// NotifyStatusChange sends a notification about a status change
func (n *StatusNotifier) NotifyStatusChange(ctx context.Context, component string, oldStatus, newStatus status.Status, message string) error {
	if !n.enabled {
		return nil
	}

	notification := Notification{
		Type:      "component",
		Component: component,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Message:   message,
		Timestamp: time.Now(),
	}

	// Try to send to channel with timeout
	select {
	case n.notificationChan <- notification:
		n.logger.Debug().
			Str("component", component).
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Msg("Component status change notification queued")
	case <-time.After(time.Second):
		n.logger.Warn().
			Str("component", component).
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Msg("Failed to queue component status change notification: channel full")
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// NotifySystemStatusChange sends a notification about a system status change
func (n *StatusNotifier) NotifySystemStatusChange(ctx context.Context, oldStatus, newStatus status.Status, message string) error {
	if !n.enabled {
		return nil
	}

	notification := Notification{
		Type:      "system",
		Component: "system",
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Message:   message,
		Timestamp: time.Now(),
	}

	// Try to send to channel with timeout
	select {
	case n.notificationChan <- notification:
		n.logger.Debug().
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Msg("System status change notification queued")
	case <-time.After(time.Second):
		n.logger.Warn().
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Msg("Failed to queue system status change notification: channel full")
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// Enable enables the notifier
func (n *StatusNotifier) Enable() {
	n.enabled = true
	n.logger.Info().Msg("Status notifier enabled")
}

// Disable disables the notifier
func (n *StatusNotifier) Disable() {
	n.enabled = false
	n.logger.Info().Msg("Status notifier disabled")
}

// processNotifications processes notifications in the background
func (n *StatusNotifier) processNotifications() {
	for notification := range n.notificationChan {
		n.processNotification(notification)
	}
}

// processNotification processes a single notification
func (n *StatusNotifier) processNotification(notification Notification) {
	// Log the notification
	logEvent := n.logger.Info().
		Str("type", notification.Type).
		Str("component", notification.Component).
		Str("old_status", string(notification.OldStatus)).
		Str("new_status", string(notification.NewStatus)).
		Time("timestamp", notification.Timestamp)

	if notification.Message != "" {
		logEvent = logEvent.Str("message", notification.Message)
	}

	logEvent.Msg("Status change notification")

	// Here you would implement sending to external notification systems
	// such as email, Slack, Telegram, etc.
	// For now, we just log it

	// Example of how you might format a message for external systems
	var message string
	if notification.Type == "system" {
		message = fmt.Sprintf("System status changed from %s to %s", notification.OldStatus, notification.NewStatus)
	} else {
		message = fmt.Sprintf("Component %s status changed from %s to %s", notification.Component, notification.OldStatus, notification.NewStatus)
	}

	if notification.Message != "" {
		message += fmt.Sprintf(": %s", notification.Message)
	}

	// TODO: Send to external notification systems
	// This would be implemented by integrating with specific notification services
	// For example:
	// - sendEmailNotification(message)
	// - sendSlackNotification(message)
	// - sendTelegramNotification(message)
}
