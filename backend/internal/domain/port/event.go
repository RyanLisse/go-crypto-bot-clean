package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// EventRepository defines the interface for event persistence
type EventRepository interface {
	// SaveEvent stores a new event
	SaveEvent(ctx context.Context, event *model.NewCoinEvent) error
	// GetEvents retrieves events for a specific coin
	GetEvents(ctx context.Context, coinID string, limit, offset int) ([]*model.NewCoinEvent, error)
}
