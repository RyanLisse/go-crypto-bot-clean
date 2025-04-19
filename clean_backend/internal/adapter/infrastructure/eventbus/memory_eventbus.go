package eventbus

import (
	"sync"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// Ensure MemoryEventBus implements the port.EventBus interface.
var _ port.EventBus = (*MemoryEventBus)(nil)

type MemoryEventBus struct {
	subscribers []func(*model.NewCoinEvent) // Using specific event type for now
	mu          sync.RWMutex
	logger      zerolog.Logger
}

func NewMemoryEventBus(logger zerolog.Logger) *MemoryEventBus {
	return &MemoryEventBus{
		subscribers: make([]func(*model.NewCoinEvent), 0),
		logger:      logger.With().Str("component", "MemoryEventBus").Logger(),
	}
}

func (b *MemoryEventBus) Publish(event *model.NewCoinEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	b.logger.Info().Str("eventType", event.EventType).Str("coinID", event.CoinID).Msg("Publishing event")
	for _, sub := range b.subscribers {
		go sub(event) // Run subscriber in a goroutine to avoid blocking publisher
	}
}

func (b *MemoryEventBus) Subscribe(callback func(*model.NewCoinEvent)) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers = append(b.subscribers, callback)
	b.logger.Info().Int("totalSubscribers", len(b.subscribers)).Msg("New event subscriber added")
	return nil
}

func (b *MemoryEventBus) Unsubscribe(callback func(*model.NewCoinEvent)) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Find and remove the subscriber
	// Note: Comparing functions directly might not be reliable.
	// A more robust implementation would assign IDs to subscribers.
	// For this placeholder, we'll skip the actual removal.
	b.logger.Warn().Msg("EventBus Unsubscribe is not fully implemented for function comparison")

	// Example of potential removal (needs reliable function comparison):
	// newSubscribers := make([]func(*model.NewCoinEvent), 0)
	// for _, sub := range b.subscribers {
	// 	if &sub != &callback { // This comparison is likely incorrect
	// 		newSubscribers = append(newSubscribers, sub)
	// 	}
	// }
	// b.subscribers = newSubscribers

	return nil
}
