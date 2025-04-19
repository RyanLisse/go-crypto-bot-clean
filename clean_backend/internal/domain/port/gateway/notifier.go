package gateway

import "context"

// Notifier defines the port for sending notifications.
type Notifier interface {
	// Send sends a notification message.
	Send(ctx context.Context, recipient string, subject string, message string) error
	// SendUrgent is similar to Send but might use a higher priority channel.
	SendUrgent(ctx context.Context, recipient string, subject string, message string) error
	// Add other notification methods as needed (e.g., SendToTopic)
}
