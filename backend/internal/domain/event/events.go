package event

import (
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
)

// EventType defines the type of a domain event.
type EventType string

const (
	// NewCoinTradableEvent signifies that a newly listed coin has become tradable.
	NewCoinTradableEvent EventType = "NewCoinTradable"
	// Add other event types here as needed...
)

// DomainEvent represents a base interface for domain events.
// It includes basic metadata common to all events.
type DomainEvent interface {
	Type() EventType
	OccurredAt() time.Time
	AggregateID() string // ID of the aggregate root that generated the event (e.g., Symbol)
}

// BaseEvent provides common fields for domain events.
type BaseEvent struct {
	eventType   EventType
	occurredAt  time.Time
	aggregateID string
}

// NewBaseEvent creates a new BaseEvent.
func NewBaseEvent(eventType EventType, aggregateID string) BaseEvent {
	return BaseEvent{
		eventType:   eventType,
		occurredAt:  time.Now().UTC(),
		aggregateID: aggregateID,
	}
}

func (e BaseEvent) Type() EventType       { return e.eventType }
func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }
func (e BaseEvent) AggregateID() string   { return e.aggregateID }

// --- Specific Event Definitions ---

// NewCoinTradable is published when a new coin becomes available for trading.
type NewCoinTradable struct {
	BaseEvent
	Symbol        string    `json:"symbol"`
	TradableAt    time.Time `json:"tradable_at"`
	InitialPrice  *float64  `json:"initial_price,omitempty"`  // Optional: Price at the time it became tradable
	InitialVolume *float64  `json:"initial_volume,omitempty"` // Optional: Volume at the time it became tradable
}

// NewNewCoinTradable creates a new NewCoinTradable event.
func NewNewCoinTradable(coin *model.NewCoin, price *float64, volume *float64) *NewCoinTradable {
	// Ensure BecameTradableAt is not nil before dereferencing
	tradableAt := time.Now().UTC() // Default to now if nil, though it shouldn't be
	if coin.BecameTradableAt != nil {
		tradableAt = *coin.BecameTradableAt
	}

	return &NewCoinTradable{
		BaseEvent:     NewBaseEvent(NewCoinTradableEvent, coin.Symbol),
		Symbol:        coin.Symbol,
		TradableAt:    tradableAt,
		InitialPrice:  price,  // Pass along if available
		InitialVolume: volume, // Pass along if available
	}
}

// Ensure NewCoinTradable implements DomainEvent (compile-time check)
var _ DomainEvent = (*NewCoinTradable)(nil)
