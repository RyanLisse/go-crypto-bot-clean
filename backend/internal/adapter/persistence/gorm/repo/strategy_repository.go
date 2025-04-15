package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// StrategyConfigEntity represents a strategy configuration in the database
type StrategyConfigEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	Name      string    `gorm:"index;type:varchar(100)"`
	Config    []byte    `gorm:"type:json"`
	Active    bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// GormStrategyRepository implements port.StrategyRepository using GORM
type GormStrategyRepository struct {
	BaseRepository
}

// NewGormStrategyRepository creates a new GormStrategyRepository
func NewGormStrategyRepository(db *gorm.DB, logger *zerolog.Logger) *GormStrategyRepository {
	return &GormStrategyRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveConfig saves a strategy configuration
func (r *GormStrategyRepository) SaveConfig(ctx context.Context, strategyID string, config map[string]interface{}) error {
	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		r.logger.Error().Err(err).Str("strategy_id", strategyID).Msg("Failed to marshal strategy config")
		return err
	}

	// Create entity
	entity := &StrategyConfigEntity{
		ID:        strategyID,
		Name:      getStrategyName(config),
		Config:    configJSON,
		UpdatedAt: time.Now(),
	}

	// Save entity
	return r.Upsert(ctx, entity, []string{"id"}, []string{
		"name", "config", "active", "updated_at",
	})
}

// GetConfig retrieves a strategy configuration
func (r *GormStrategyRepository) GetConfig(ctx context.Context, strategyID string) (map[string]interface{}, error) {
	var entity StrategyConfigEntity
	err := r.FindOne(ctx, &entity, "id = ?", strategyID)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, nil // Not found
	}

	// Parse config
	var config map[string]interface{}
	if len(entity.Config) > 0 {
		if err := json.Unmarshal(entity.Config, &config); err != nil {
			r.logger.Error().Err(err).Str("strategy_id", strategyID).Msg("Failed to unmarshal strategy config")
			return nil, err
		}
	}

	return config, nil
}

// ListStrategies lists all strategy IDs
func (r *GormStrategyRepository) ListStrategies(ctx context.Context) ([]string, error) {
	var entities []StrategyConfigEntity
	err := r.GetDB(ctx).
		Where("active = ?", true).
		Order("name ASC").
		Find(&entities).Error
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to list strategies")
		return nil, err
	}

	// Extract IDs
	ids := make([]string, len(entities))
	for i, entity := range entities {
		ids[i] = entity.ID
	}

	return ids, nil
}

// DeleteStrategy deletes a strategy
func (r *GormStrategyRepository) DeleteStrategy(ctx context.Context, strategyID string) error {
	// Soft delete by setting active to false
	return r.Update(ctx, &StrategyConfigEntity{ID: strategyID}, map[string]interface{}{
		"active":     false,
		"updated_at": time.Now(),
	})
}

// Helper function to extract strategy name from config
func getStrategyName(config map[string]interface{}) string {
	if name, ok := config["name"].(string); ok {
		return name
	}
	return "Unnamed Strategy"
}

// Ensure GormStrategyRepository implements port.StrategyRepository
var _ port.StrategyRepository = (*GormStrategyRepository)(nil)
