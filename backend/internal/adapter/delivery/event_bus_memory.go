package delivery

import (
	"fmt"
	"sync"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// InMemoryEventBus implements port.EventBus using in-memory channels
type InMemoryEventBus struct {
	listeners []func(*model.NewCoinEvent)
	mu        sync.RWMutex
	logger    zerolog.Logger
}

// NewInMemoryEventBus creates a new InMemoryEventBus
func NewInMemoryEventBus(logger zerolog.Logger) *InMemoryEventBus {
	return &InMemoryEventBus{
		listeners: make([]func(*model.NewCoinEvent), 0),
		logger:    logger.With().Str("component", "InMemoryEventBus").Logger(),
	}
}

// Publish sends an event to all registered listeners asynchronously
func (b *InMemoryEventBus) Publish(event *model.NewCoinEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	b.logger.Info().Str("event_type", event.EventType).Str("coin_id", event.CoinID).Msg("Publishing event")
	for _, listener := range b.listeners {
		go func(l func(*model.NewCoinEvent)) {
			defer func() {
				if r := recover(); r != nil {
					b.logger.Error().Interface("panic", r).Msg("Recovered from panic in event listener")
				}
			}()
			l(event)
		}(listener)
	}
}

// Subscribe adds a listener for new coin events
func (b *InMemoryEventBus) Subscribe(listener func(*model.NewCoinEvent)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners = append(b.listeners, listener)
	b.logger.Info().Msg("New event listener subscribed")
}

// Unsubscribe removes a listener
func (b *InMemoryEventBus) Unsubscribe(listener func(*model.NewCoinEvent)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, l := range b.listeners {
		// Compare function pointers to identify the listener
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			b.logger.Info().Msg("Event listener unsubscribed")
			break
		}
	}
}
