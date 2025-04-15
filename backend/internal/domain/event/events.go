package event

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
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

// NewCoinTradable represents a newly tradable coin event
type NewCoinTradable struct {
	BaseEvent
	Symbol        string    `json:"symbol"`
	TradableAt    time.Time `json:"tradable_at"`
	InitialPrice  float64   `json:"initial_price,omitempty"`
	InitialVolume float64   `json:"initial_volume,omitempty"`
	Price         float64   `json:"price"`
	Volume        float64   `json:"volume"`
	QuoteAsset    string    `json:"quote_asset"`
}

// NewNewCoinTradable creates a new NewCoinTradable event.
func NewNewCoinTradable(coin *model.NewCoin, price *float64, volume *float64) *NewCoinTradable {
	// Ensure BecameTradableAt is not nil before dereferencing
	tradableAt := time.Now().UTC() // Default to now if nil, though it shouldn't be
	if coin.BecameTradableAt != nil {
		tradableAt = *coin.BecameTradableAt
	}

	// Handle nil values for price and volume
	var initialPrice, initialVolume, currentPrice, currentVolume float64
	if price != nil {
		initialPrice = *price
		currentPrice = *price
	}
	if volume != nil {
		initialVolume = *volume
		currentVolume = *volume
	}

	return &NewCoinTradable{
		BaseEvent:     NewBaseEvent(NewCoinTradableEvent, coin.Symbol),
		Symbol:        coin.Symbol,
		TradableAt:    tradableAt,
		InitialPrice:  initialPrice,
		InitialVolume: initialVolume,
		Price:         currentPrice,
		Volume:        currentVolume,
		QuoteAsset:    coin.QuoteAsset,
	}
}

// Ensure NewCoinTradable implements DomainEvent (compile-time check)
var _ DomainEvent = (*NewCoinTradable)(nil)
