package gorm

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Error definitions
var (
	// ErrPositionNotFound is returned when a position cannot be found
	ErrPositionNotFound = errors.New("position not found")
)

// Ensure PositionRepository implements port.PositionRepository
var _ port.PositionRepository = (*PositionRepository)(nil)

// PositionEntity is defined in entity.go

// PositionRepository implements the port.PositionRepository interface using GORM
type PositionRepository struct {
	db *gorm.DB
}

// NewPositionRepository creates a new instance of PositionRepository
func NewPositionRepository(db *gorm.DB) *PositionRepository {
	return &PositionRepository{
		db: db,
	}
}

// Create creates a new position in the database
func (r *PositionRepository) Create(ctx context.Context, position *model.Position) error {
	entity := r.toEntity(position)
	result := r.db.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("positionID", position.ID).Msg("Failed to create position")
		return result.Error
	}

	log.Debug().Str("positionID", position.ID).Msg("Position created successfully")
	return nil
}

// GetByID retrieves a position by its ID
func (r *PositionRepository) GetByID(ctx context.Context, id string) (*model.Position, error) {
	var entity PositionEntity
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Debug().Str("positionID", id).Msg("Position not found")
			return nil, ErrPositionNotFound
		}
		log.Error().Err(result.Error).Str("positionID", id).Msg("Failed to get position")
		return nil, result.Error
	}

	return r.toDomain(&entity), nil
}

// Update updates an existing position
func (r *PositionRepository) Update(ctx context.Context, position *model.Position) error {
	entity := r.toEntity(position)
	result := r.db.WithContext(ctx).Where("id = ?", position.ID).Updates(&entity)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("positionID", position.ID).Msg("Failed to update position")
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Debug().Str("positionID", position.ID).Msg("Position not found for update")
		return ErrPositionNotFound
	}

	log.Debug().Str("positionID", position.ID).Msg("Position updated successfully")
	return nil
}

// GetOpenPositions retrieves all open positions
func (r *PositionRepository) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.db.WithContext(ctx).Where("status = ?", string(model.PositionStatusOpen)).Find(&entities)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get open positions")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().Int("count", len(positions)).Msg("Retrieved open positions")
	return positions, nil
}

// GetOpenPositionsBySymbol retrieves all open positions for a specific symbol
func (r *PositionRepository) GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.db.WithContext(ctx).Where("status = ? AND symbol = ?",
		string(model.PositionStatusOpen), symbol).Find(&entities)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get open positions by symbol")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().Str("symbol", symbol).Int("count", len(positions)).Msg("Retrieved open positions by symbol")
	return positions, nil
}

// GetOpenPositionsByType retrieves all open positions for a specific type
func (r *PositionRepository) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.db.WithContext(ctx).Where("status = ? AND type = ?",
		string(model.PositionStatusOpen), string(positionType)).Find(&entities)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("type", string(positionType)).Msg("Failed to get open positions by type")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().Str("type", string(positionType)).Int("count", len(positions)).Msg("Retrieved open positions by type")
	return positions, nil
}

// GetBySymbol retrieves positions for a specific symbol with pagination
func (r *PositionRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).
		Limit(limit).Offset(offset).
		Order("opened_at DESC").
		Find(&entities)

	if result.Error != nil {
		log.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get positions by symbol")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().Str("symbol", symbol).Int("count", len(positions)).Msg("Retrieved positions by symbol")
	return positions, nil
}

// GetByUserID retrieves positions for a specific user with pagination
func (r *PositionRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error) {
	// Note: The Position model doesn't have a UserID field, so this is a placeholder implementation
	// that might need to be updated according to actual requirements
	log.Warn().Msg("GetByUserID called, but Position model doesn't have UserID field")
	return []*model.Position{}, nil
}

// GetActiveByUser retrieves active positions for a specific user
func (r *PositionRepository) GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.db.WithContext(ctx).Where("status = ? AND user_id = ?", string(model.PositionStatusOpen), userID).Find(&entities)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("userID", userID).Msg("Failed to get active positions by user")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().Str("userID", userID).Int("count", len(positions)).Msg("Retrieved active positions by user")
	return positions, nil
}

// GetClosedPositions retrieves closed positions within a time range with pagination
func (r *PositionRepository) GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.db.WithContext(ctx).
		Where("status = ? AND closed_at BETWEEN ? AND ?", string(model.PositionStatusClosed), from, to).
		Limit(limit).Offset(offset).
		Order("closed_at DESC").
		Find(&entities)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get closed positions")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().Int("count", len(positions)).Msg("Retrieved closed positions")
	return positions, nil
}

// Count counts positions based on provided filters
func (r *PositionRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&PositionEntity{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	result := query.Count(&count)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count positions")
		return 0, result.Error
	}

	log.Debug().Int64("count", count).Msg("Counted positions")
	return count, nil
}

// Delete deletes a position by its ID
func (r *PositionRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&PositionEntity{})
	if result.Error != nil {
		log.Error().Err(result.Error).Str("positionID", id).Msg("Failed to delete position")
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Debug().Str("positionID", id).Msg("Position not found for deletion")
		return ErrPositionNotFound
	}

	log.Debug().Str("positionID", id).Msg("Position deleted successfully")
	return nil
}

// GetBySymbolAndUser retrieves positions for a specific symbol and user with pagination
func (r *PositionRepository) GetBySymbolAndUser(ctx context.Context, symbol, userID string, page, limit int) ([]*model.Position, error) {
	var entities []PositionEntity

	// Calculate offset from page and limit
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND user_id = ?", symbol, userID).
		Limit(limit).
		Offset(offset).
		Order("opened_at DESC").
		Find(&entities)

	if result.Error != nil {
		log.Error().Err(result.Error).
			Str("symbol", symbol).
			Str("userID", userID).
			Int("page", page).
			Int("limit", limit).
			Msg("Failed to get positions by symbol and user")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().
		Str("symbol", symbol).
		Str("userID", userID).
		Int("page", page).
		Int("limit", limit).
		Int("count", len(positions)).
		Msg("Retrieved positions by symbol and user")
	return positions, nil
}

// GetOpenPositionsByUserID retrieves all open positions for a specific user
func (r *PositionRepository) GetOpenPositionsByUserID(ctx context.Context, userID string) ([]*model.Position, error) {
	var entities []PositionEntity
	result := r.db.WithContext(ctx).Where("status = ? AND user_id = ?", string(model.PositionStatusOpen), userID).Find(&entities)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("userID", userID).Msg("Failed to get open positions by user ID")
		return nil, result.Error
	}

	positions := make([]*model.Position, len(entities))
	for i, entity := range entities {
		positions[i] = r.toDomain(&entity)
	}

	log.Debug().Str("userID", userID).Int("count", len(positions)).Msg("Retrieved open positions by user ID")
	return positions, nil
}

// Helper methods for entity conversion
func (r *PositionRepository) toEntity(position *model.Position) *PositionEntity {
	// Convert slice fields to JSON strings
	entryOrderIDsJSON, err := json.Marshal(position.EntryOrderIDs)
	if err != nil {
		log.Error().Err(err).Str("positionID", position.ID).Msg("Failed to marshal EntryOrderIDs")
		entryOrderIDsJSON = []byte("[]")
	}

	exitOrderIDsJSON, err := json.Marshal(position.ExitOrderIDs)
	if err != nil {
		log.Error().Err(err).Str("positionID", position.ID).Msg("Failed to marshal ExitOrderIDs")
		exitOrderIDsJSON = []byte("[]")
	}

	openOrderIDsJSON, err := json.Marshal(position.OpenOrderIDs)
	if err != nil {
		log.Error().Err(err).Str("positionID", position.ID).Msg("Failed to marshal OpenOrderIDs")
		openOrderIDsJSON = []byte("[]")
	}

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
		EntryOrderIDs:   string(entryOrderIDsJSON),
		ExitOrderIDs:    string(exitOrderIDsJSON),
		OpenOrderIDs:    string(openOrderIDsJSON),
		Notes:           position.Notes,
		OpenedAt:        position.OpenedAt,
		ClosedAt:        position.ClosedAt,
		LastUpdatedAt:   position.LastUpdatedAt,
		MaxDrawdown:     position.MaxDrawdown,
		MaxProfit:       position.MaxProfit,
		RiskRewardRatio: position.RiskRewardRatio,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (r *PositionRepository) toDomain(entity *PositionEntity) *model.Position {
	// Parse JSON strings back to slices
	var entryOrderIDs, exitOrderIDs, openOrderIDs []string

	// Handle possible empty strings or invalid JSON
	if entity.EntryOrderIDs != "" {
		if err := json.Unmarshal([]byte(entity.EntryOrderIDs), &entryOrderIDs); err != nil {
			log.Error().Err(err).Str("positionID", entity.ID).Msg("Failed to unmarshal EntryOrderIDs")
			entryOrderIDs = []string{}
		}
	}

	if entity.ExitOrderIDs != "" {
		if err := json.Unmarshal([]byte(entity.ExitOrderIDs), &exitOrderIDs); err != nil {
			log.Error().Err(err).Str("positionID", entity.ID).Msg("Failed to unmarshal ExitOrderIDs")
			exitOrderIDs = []string{}
		}
	}

	if entity.OpenOrderIDs != "" {
		if err := json.Unmarshal([]byte(entity.OpenOrderIDs), &openOrderIDs); err != nil {
			log.Error().Err(err).Str("positionID", entity.ID).Msg("Failed to unmarshal OpenOrderIDs")
			openOrderIDs = []string{}
		}
	}

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
	}
}
