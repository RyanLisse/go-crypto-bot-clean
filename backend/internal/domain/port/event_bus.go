package port

import (
	"context"

	"github.com/neo/crypto-bot/internal/domain/event"
)

// EventBus defines the interface for publishing domain events.
// Implementations could be in-memory, Kafka, NATS, etc.
type EventBus interface {
	// Publish sends a domain event to the event bus.
	Publish(ctx context.Context, event event.DomainEvent) error
}
