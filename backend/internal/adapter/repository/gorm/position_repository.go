package gorm

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Error definitions
var (
	// ErrPositionNotFound is returned when a position cannot be found
	ErrPositionNotFound = errors.New("position not found")
)

// Ensure PositionRepository implements port.PositionRepository
var _ port.PositionRepository = (*PositionRepository)(nil)

// PositionEntity represents the database model for a position
type PositionEntity struct {
	ID              string `gorm:"primaryKey"`
	Symbol          string `gorm:"index"`
	Side            string `gorm:"index"`
	Status          string `gorm:"index"`
	Type            string `gorm:"index"`
	EntryPrice      float64
	Quantity        float64
	CurrentPrice    float64
	PnL             float64
	PnLPercent      float64
	StopLoss        *float64
	TakeProfit      *float64
	StrategyID      *string
	EntryOrderIDs   string // Stored as JSON array
	ExitOrderIDs    string // Stored as JSON array
	OpenOrderIDs    string // Stored as JSON array
	Notes           string
	OpenedAt        time.Time
	ClosedAt        *time.Time
	LastUpdatedAt   time.Time
	MaxDrawdown     float64
	MaxProfit       float64
	RiskRewardRatio float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PositionRepository implements the port.PositionRepository interface using GORM
type PositionRepository struct {
	base BaseRepository
}

// NewPositionRepository creates a new instance of PositionRepository
func NewPositionRepository(db *gorm.DB, logger *zerolog.Logger) *PositionRepository {
	l := logger.With().Str("component", "position_repository").Logger()
	return &PositionRepository{
		base: NewBaseRepository(db, &l),
	}
}

// Create creates a new position in the database
func (r *PositionRepository) Create(ctx context.Context, position *model.Position) error {
	entity := r.toEntity(position)
	result := r.base.GetDB(ctx).Create(&entity)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("positionID", position.ID).Msg("Failed to create position")
		return result.Error
	}

	r.base.logger.Debug().Str("positionID", position.ID).Msg("Position created successfully")
	return nil
}

// GetByID retrieves a position by its ID
func (r *PositionRepository) GetByID(ctx context.Context, id string) (*model.Position, error) {
	var entity PositionEntity
	result := r.base.GetDB(ctx).Where("id = ?", id).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.base.logger.Debug().Str("positionID", id).Msg("Position not found")
			return nil, ErrPositionNotFound
		}
		r.base.logger.Error().Err(result.Error).Str("positionID", id).Msg("Failed to get position")
		return nil, result.Error
	}

	return r.toDomain(&entity), nil
}

// Update updates an existing position
func (r *PositionRepository) Update(ctx context.Context, position *model.Position) error {
	entity := r.toEntity(position)
	result := r.base.GetDB(ctx).Where("id = ?", position.ID).Updates(&entity)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("positionID", position.ID).Msg("Failed to update position")
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.base.logger.Debug().Str("positionID", position.ID).Msg("Position not found for update")
		return ErrPositionNotFound
	}

	r.base.logger.Debug().Str("positionID", position.ID).Msg("Position updated successfully")
	return nil
}

// GetOpenPositions retrieves all open positions
func (r *PositionRepository) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.base.GetDB(ctx).Where("status = ?", string(model.PositionStatusOpen)).Find(&entities)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Msg("Failed to get open positions")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().Int("count", len(positions)).Msg("Retrieved open positions")
	return positions, nil
}

// GetOpenPositionsBySymbol retrieves all open positions for a specific symbol
func (r *PositionRepository) GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.base.GetDB(ctx).Where("status = ? AND symbol = ?",
		string(model.PositionStatusOpen), symbol).Find(&entities)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get open positions by symbol")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().Str("symbol", symbol).Int("count", len(positions)).Msg("Retrieved open positions by symbol")
	return positions, nil
}

// GetOpenPositionsByType retrieves all open positions for a specific type
func (r *PositionRepository) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.base.GetDB(ctx).Where("status = ? AND type = ?",
		string(model.PositionStatusOpen), string(positionType)).Find(&entities)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("type", string(positionType)).Msg("Failed to get open positions by type")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().Str("type", string(positionType)).Int("count", len(positions)).Msg("Retrieved open positions by type")
	return positions, nil
}

// GetBySymbol retrieves positions for a specific symbol with pagination
func (r *PositionRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.base.GetDB(ctx).Where("symbol = ?", symbol).
		Limit(limit).Offset(offset).
		Order("opened_at DESC").
		Find(&entities)

	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get positions by symbol")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().Str("symbol", symbol).Int("count", len(positions)).Msg("Retrieved positions by symbol")
	return positions, nil
}

// GetByUserID retrieves positions for a specific user with pagination
func (r *PositionRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error) {
	// Note: The Position model doesn't have a UserID field, so this is a placeholder implementation
	// that might need to be updated according to actual requirements
	r.base.logger.Warn().Msg("GetByUserID called, but Position model doesn't have UserID field")
	return []*model.Position{}, nil
}

// GetActiveByUser retrieves active positions for a specific user
func (r *PositionRepository) GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.base.GetDB(ctx).Where("status = ? AND user_id = ?", string(model.PositionStatusOpen), userID).Find(&entities)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("userID", userID).Msg("Failed to get active positions by user")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().Str("userID", userID).Int("count", len(positions)).Msg("Retrieved active positions by user")
	return positions, nil
}

// GetClosedPositions retrieves closed positions within a time range with pagination
func (r *PositionRepository) GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.base.GetDB(ctx).Where("status = ? AND closed_at BETWEEN ? AND ?",
		string(model.PositionStatusClosed), from, to).
		Limit(limit).Offset(offset).
		Order("closed_at DESC").
		Find(&entities)

	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Msg("Failed to get closed positions")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().Int("count", len(positions)).Msg("Retrieved closed positions")
	return positions, nil
}

// Count returns the number of positions matching the given filters
func (r *PositionRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.base.GetDB(ctx)

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	result := query.Model(&PositionEntity{}).Count(&count)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Msg("Failed to count positions")
		return 0, result.Error
	}

	r.base.logger.Debug().Int64("count", count).Msg("Counted positions")
	return count, nil
}

// Delete deletes a position by its ID
func (r *PositionRepository) Delete(ctx context.Context, id string) error {
	result := r.base.GetDB(ctx).Where("id = ?", id).Delete(&PositionEntity{})
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("positionID", id).Msg("Failed to delete position")
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.base.logger.Debug().Str("positionID", id).Msg("Position not found for deletion")
		return ErrPositionNotFound
	}

	r.base.logger.Debug().Str("positionID", id).Msg("Position deleted successfully")
	return nil
}

// GetBySymbolAndUser retrieves positions for a specific symbol and user with pagination
func (r *PositionRepository) GetBySymbolAndUser(ctx context.Context, symbol, userID string, page, limit int) ([]*model.Position, error) {
	var entities []PositionEntity
	offset := (page - 1) * limit

	result := r.base.GetDB(ctx).Where("symbol = ? AND user_id = ?", symbol, userID).
		Limit(limit).Offset(offset).
		Order("opened_at DESC").
		Find(&entities)

	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).
			Str("symbol", symbol).
			Str("userID", userID).
			Msg("Failed to get positions by symbol and user")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().
		Str("symbol", symbol).
		Str("userID", userID).
		Int("count", len(positions)).
		Msg("Retrieved positions by symbol and user")
	return positions, nil
}

// GetOpenPositionsByUserID retrieves all open positions for a specific user
func (r *PositionRepository) GetOpenPositionsByUserID(ctx context.Context, userID string) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.base.GetDB(ctx).Where("status = ? AND user_id = ?",
		string(model.PositionStatusOpen), userID).Find(&entities)
	if result.Error != nil {
		r.base.logger.Error().Err(result.Error).Str("userID", userID).Msg("Failed to get open positions by user ID")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	r.base.logger.Debug().Str("userID", userID).Int("count", len(positions)).Msg("Retrieved open positions by user ID")
	return positions, nil
}

// toEntity converts a domain model to a database entity
func (r *PositionRepository) toEntity(position *model.Position) *PositionEntity {
	entryOrderIDs, _ := json.Marshal(position.EntryOrderIDs)
	exitOrderIDs, _ := json.Marshal(position.ExitOrderIDs)
	openOrderIDs, _ := json.Marshal(position.OpenOrderIDs)

	return &PositionEntity{
		ID:              position.ID,
		Symbol:          position.Symbol,
		Side:            string(position.Side),
		Status:          string(position.Status),
		Type:            string(position.Type),
		EntryPrice:      position.EntryPrice,
		Quantity:        position.Quantity,
		CurrentPrice:    position.CurrentPrice,
		PnL:             position.PnL,
		PnLPercent:      position.PnLPercent,
		StopLoss:        position.StopLoss,
		TakeProfit:      position.TakeProfit,
		StrategyID:      position.StrategyID,
		EntryOrderIDs:   string(entryOrderIDs),
		ExitOrderIDs:    string(exitOrderIDs),
		OpenOrderIDs:    string(openOrderIDs),
		Notes:           position.Notes,
		OpenedAt:        position.OpenedAt,
		ClosedAt:        position.ClosedAt,
		LastUpdatedAt:   position.LastUpdatedAt,
		MaxDrawdown:     position.MaxDrawdown,
		MaxProfit:       position.MaxProfit,
		RiskRewardRatio: position.RiskRewardRatio,
		CreatedAt:       position.CreatedAt,
		UpdatedAt:       position.UpdatedAt,
	}
}

// toDomain converts a database entity to a domain model
func (r *PositionRepository) toDomain(entity *PositionEntity) *model.Position {
	var entryOrderIDs []string
	var exitOrderIDs []string
	var openOrderIDs []string

	_ = json.Unmarshal([]byte(entity.EntryOrderIDs), &entryOrderIDs)
	_ = json.Unmarshal([]byte(entity.ExitOrderIDs), &exitOrderIDs)
	_ = json.Unmarshal([]byte(entity.OpenOrderIDs), &openOrderIDs)

	return &model.Position{
		ID:              entity.ID,
		Symbol:          entity.Symbol,
		Side:            model.PositionSide(entity.Side),
		Status:          model.PositionStatus(entity.Status),
		Type:            model.PositionType(entity.Type),
		EntryPrice:      entity.EntryPrice,
		Quantity:        entity.Quantity,
		CurrentPrice:    entity.CurrentPrice,
		PnL:             entity.PnL,
		PnLPercent:      entity.PnLPercent,
		StopLoss:        entity.StopLoss,
		TakeProfit:      entity.TakeProfit,
		StrategyID:      entity.StrategyID,
		EntryOrderIDs:   entryOrderIDs,
		ExitOrderIDs:    exitOrderIDs,
		OpenOrderIDs:    openOrderIDs,
		Notes:           entity.Notes,
		OpenedAt:        entity.OpenedAt,
		ClosedAt:        entity.ClosedAt,
		LastUpdatedAt:   entity.LastUpdatedAt,
		MaxDrawdown:     entity.MaxDrawdown,
		MaxProfit:       entity.MaxProfit,
		RiskRewardRatio: entity.RiskRewardRatio,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}
