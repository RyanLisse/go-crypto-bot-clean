package port

import "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"

// EventBus defines the interface for an event bus system
type EventBus interface {
	// Publish sends an event to all subscribers
	Publish(event *model.NewCoinEvent)

	// Subscribe adds a listener for new coin events
	Subscribe(listener func(*model.NewCoinEvent))

	// Unsubscribe removes a listener
	Unsubscribe(listener func(*model.NewCoinEvent))
}
