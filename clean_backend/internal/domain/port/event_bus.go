package port

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// EventBus defines the interface for a domain event bus.
type EventBus interface {
	// Publish publishes an event to all subscribers.
	// Using specific event type for now, could be made generic later.
	Publish(event *model.NewCoinEvent)

	// Subscribe adds a callback function to be invoked when events are published.
	Subscribe(callback func(*model.NewCoinEvent)) error

	// Unsubscribe removes a previously subscribed callback.
	Unsubscribe(callback func(*model.NewCoinEvent)) error
}
