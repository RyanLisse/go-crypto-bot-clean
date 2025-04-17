package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// EventEntity represents an event in the database
type EventEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	CoinID    string    `gorm:"index;type:varchar(50)"`
	EventType string    `gorm:"type:varchar(50)"`
	OldStatus string    `gorm:"type:varchar(20)"`
	NewStatus string    `gorm:"type:varchar(20)"`
	Data      []byte    `gorm:"type:json"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// GormEventRepository implements port.EventRepository using GORM
type GormEventRepository struct {
	BaseRepository
}

// NewGormEventRepository creates a new GormEventRepository
func NewGormEventRepository(db *gorm.DB, logger *zerolog.Logger) *GormEventRepository {
	return &GormEventRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveEvent stores a new event
func (r *GormEventRepository) SaveEvent(ctx context.Context, event *model.NewCoinEvent) error {
	entity := r.toEntity(event)
	return r.Create(ctx, entity)
}

// GetEvents retrieves events for a specific coin
func (r *GormEventRepository) GetEvents(ctx context.Context, coinID string, limit, offset int) ([]*model.NewCoinEvent, error) {
	var entities []EventEntity

	err := r.GetDB(ctx).
		Where("coin_id = ?", coinID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(entities), nil
}

// Helper methods for entity conversion

// toEntity converts a domain event to a database entity
func (r *GormEventRepository) toEntity(event *model.NewCoinEvent) *EventEntity {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	var data []byte
	if event.Data != nil {
		var err error
		data, err = json.Marshal(event.Data)
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to marshal event data")
		}
	}

	return &EventEntity{
		ID:        event.ID,
		CoinID:    event.CoinID,
		EventType: event.EventType,
		OldStatus: string(event.OldStatus),
		NewStatus: string(event.NewStatus),
		Data:      data,
		CreatedAt: event.CreatedAt,
	}
}

// toDomain converts a database entity to a domain event
func (r *GormEventRepository) toDomain(entity *EventEntity) *model.NewCoinEvent {
	if entity == nil {
		return nil
	}

	var data interface{}
	if len(entity.Data) > 0 {
		if err := json.Unmarshal(entity.Data, &data); err != nil {
			r.logger.Error().Err(err).Msg("Failed to unmarshal event data")
		}
	}

	return &model.NewCoinEvent{
		ID:        entity.ID,
		CoinID:    entity.CoinID,
		EventType: entity.EventType,
       OldStatus: model.CoinStatus(entity.OldStatus),
       NewStatus: model.CoinStatus(entity.NewStatus),
		Data:      data,
		CreatedAt: entity.CreatedAt,
	}
}

// toDomainSlice converts a slice of database entities to domain events
func (r *GormEventRepository) toDomainSlice(entities []EventEntity) []*model.NewCoinEvent {
	events := make([]*model.NewCoinEvent, len(entities))
	for i, entity := range entities {
		events[i] = r.toDomain(&entity)
	}
	return events
}

// Ensure GormEventRepository implements port.EventRepository
var _ port.EventRepository = (*GormEventRepository)(nil)
