package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// NewCoinEntity represents a new coin in the database
type NewCoinEntity struct {
	ID                    string    `gorm:"primaryKey;type:varchar(50)"`
	Symbol                string    `gorm:"uniqueIndex;type:varchar(20)"`
	Name                  string    `gorm:"type:varchar(100)"`
	Status                string    `gorm:"type:varchar(20);index"`
	ExpectedListingTime   time.Time `gorm:"index"`
	BecameTradableAt      *time.Time
	BaseAsset             string `gorm:"type:varchar(10)"`
	QuoteAsset            string `gorm:"type:varchar(10)"`
	MinPrice              float64
	MaxPrice              float64
	MinQty                float64
	MaxQty                float64
	PriceScale            int
	QtyScale              int
	IsProcessedForAutobuy bool      `gorm:"default:false"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime"`
}

// NewCoinEventEntity represents a new coin event in the database
type NewCoinEventEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	CoinID    string    `gorm:"index;type:varchar(50)"`
	EventType string    `gorm:"type:varchar(50)"`
	OldStatus string    `gorm:"type:varchar(20)"`
	NewStatus string    `gorm:"type:varchar(20)"`
	Data      []byte    `gorm:"type:json"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// GormNewCoinRepository implements port.NewCoinRepository using GORM
type GormNewCoinRepository struct {
	BaseRepository
}

// NewGormNewCoinRepository creates a new GormNewCoinRepository
func NewGormNewCoinRepository(db *gorm.DB, logger *zerolog.Logger) *GormNewCoinRepository {
	return &GormNewCoinRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Save creates or updates a new coin in the database
func (r *GormNewCoinRepository) Save(ctx context.Context, coin *model.NewCoin) error {
	entity := r.toEntity(coin)

	// Use upsert to handle both create and update
	return r.Upsert(ctx, entity, []string{"id"}, []string{
		"symbol", "name", "status", "expected_listing_time", "became_tradable_at",
		"base_asset", "quote_asset", "min_price", "max_price", "min_qty", "max_qty",
		"price_scale", "qty_scale", "is_processed_for_autobuy", "updated_at",
	})
}

// GetByID retrieves a coin by its ID
func (r *GormNewCoinRepository) GetByID(ctx context.Context, id string) (*model.NewCoin, error) {
	var entity NewCoinEntity
	// Use FindOne from BaseRepository, assuming it handles not found correctly
	err := r.FindOne(ctx, &entity, "id = ?", id)
	if err != nil {
		// If FindOne returns gorm.ErrRecordNotFound, return nil, nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err // Return other errors
	}
	return r.toDomain(&entity), nil
}

// GetBySymbol retrieves a coin by its trading symbol
func (r *GormNewCoinRepository) GetBySymbol(ctx context.Context, symbol string) (*model.NewCoin, error) {
	var entity NewCoinEntity
	err := r.FindOne(ctx, &entity, "symbol = ?", symbol)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, nil // Not found
	}

	return r.toDomain(&entity), nil
}

// GetRecent retrieves recently listed coins
func (r *GormNewCoinRepository) GetRecent(ctx context.Context, limit int) ([]*model.NewCoin, error) {
	var entities []NewCoinEntity

   err := r.GetDB(ctx).
       Where("status = ?", string(model.CoinStatusTrading)).
		Order("became_tradable_at DESC").
		Limit(limit).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(entities), nil
}

// GetByStatus retrieves coins with a specific status
func (r *GormNewCoinRepository) GetByStatus(ctx context.Context, status model.CoinStatus) ([]*model.NewCoin, error) {
	var entities []NewCoinEntity

	err := r.GetDB(ctx).
		Where("status = ?", string(status)).
		Order("expected_listing_time DESC").
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(entities), nil
}

// Update updates an existing coin
func (r *GormNewCoinRepository) Update(ctx context.Context, coin *model.NewCoin) error {
	entity := r.toEntity(coin)

	result := r.GetDB(ctx).Save(entity)
	return result.Error
}

// FindRecentlyListed retrieves coins expected to list soon or recently became tradable
func (r *GormNewCoinRepository) FindRecentlyListed(ctx context.Context, thresholdTime time.Time) ([]*model.NewCoin, error) {
	var entities []NewCoinEntity

	err := r.GetDB(ctx).
       Where("(status = ? AND expected_listing_time <= ?) OR (status = ? AND became_tradable_at >= ?)",
           string(model.CoinStatusListed), time.Now(),
           string(model.CoinStatusTrading), thresholdTime).
		Order("CASE WHEN status = 'listed' THEN 0 ELSE 1 END, expected_listing_time ASC").
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(entities), nil
}

// SaveEvent saves a new coin event
func (r *GormNewCoinRepository) SaveEvent(ctx context.Context, event *model.NewCoinEvent) error {
	entity := r.toEventEntity(event)

	return r.Create(ctx, entity)
}

// GetEvents retrieves events for a specific coin
func (r *GormNewCoinRepository) GetEvents(ctx context.Context, coinID string, limit, offset int) ([]*model.NewCoinEvent, error) {
	var entities []NewCoinEventEntity

	err := r.GetDB(ctx).
		Where("coin_id = ?", coinID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return r.toEventDomainSlice(entities), nil
}

// Helper methods for entity conversion

// toEntity converts a domain model to a database entity
func (r *GormNewCoinRepository) toEntity(coin *model.NewCoin) *NewCoinEntity {
	if coin.ID == "" {
		coin.ID = uuid.New().String()
	}

	return &NewCoinEntity{
		ID:                    coin.ID,
		Symbol:                coin.Symbol,
		Name:                  coin.Name,
		Status:                string(coin.Status),
		ExpectedListingTime:   coin.ExpectedListingTime,
		BecameTradableAt:      coin.BecameTradableAt,
		BaseAsset:             coin.BaseAsset,
		QuoteAsset:            coin.QuoteAsset,
		MinPrice:              coin.MinPrice,
		MaxPrice:              coin.MaxPrice,
		MinQty:                coin.MinQty,
		MaxQty:                coin.MaxQty,
		PriceScale:            coin.PriceScale,
		QtyScale:              coin.QtyScale,
		IsProcessedForAutobuy: coin.IsProcessedForAutobuy,
		CreatedAt:             coin.CreatedAt,
		UpdatedAt:             coin.UpdatedAt,
	}
}

// toDomain converts a database entity to a domain model
func (r *GormNewCoinRepository) toDomain(entity *NewCoinEntity) *model.NewCoin {
	if entity == nil {
		return nil
	}

	return &model.NewCoin{
		ID:                    entity.ID,
		Symbol:                entity.Symbol,
		Name:                  entity.Name,
       Status:                model.CoinStatus(entity.Status),
		ExpectedListingTime:   entity.ExpectedListingTime,
		BecameTradableAt:      entity.BecameTradableAt,
		BaseAsset:             entity.BaseAsset,
		QuoteAsset:            entity.QuoteAsset,
		MinPrice:              entity.MinPrice,
		MaxPrice:              entity.MaxPrice,
		MinQty:                entity.MinQty,
		MaxQty:                entity.MaxQty,
		PriceScale:            entity.PriceScale,
		QtyScale:              entity.QtyScale,
		IsProcessedForAutobuy: entity.IsProcessedForAutobuy,
		CreatedAt:             entity.CreatedAt,
		UpdatedAt:             entity.UpdatedAt,
	}
}

// toDomainSlice converts a slice of database entities to domain models
func (r *GormNewCoinRepository) toDomainSlice(entities []NewCoinEntity) []*model.NewCoin {
	coins := make([]*model.NewCoin, len(entities))
	for i, entity := range entities {
		coins[i] = r.toDomain(&entity)
	}
	return coins
}

// toEventEntity converts a domain event to a database entity
func (r *GormNewCoinRepository) toEventEntity(event *model.NewCoinEvent) *NewCoinEventEntity {
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

	return &NewCoinEventEntity{
		ID:        event.ID,
		CoinID:    event.CoinID,
		EventType: event.EventType,
		OldStatus: string(event.OldStatus),
		NewStatus: string(event.NewStatus),
		Data:      data,
		CreatedAt: event.CreatedAt,
	}
}

// toEventDomain converts a database entity to a domain event
func (r *GormNewCoinRepository) toEventDomain(entity *NewCoinEventEntity) *model.NewCoinEvent {
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

// toEventDomainSlice converts a slice of database entities to domain events
func (r *GormNewCoinRepository) toEventDomainSlice(entities []NewCoinEventEntity) []*model.NewCoinEvent {
	events := make([]*model.NewCoinEvent, len(entities))
	for i, entity := range entities {
		events[i] = r.toEventDomain(&entity)
	}
	return events
}

// Ensure GormNewCoinRepository implements port.NewCoinRepository
var _ port.NewCoinRepository = (*GormNewCoinRepository)(nil)
