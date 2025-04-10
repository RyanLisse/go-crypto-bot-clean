package repositories

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"

	"gorm.io/gorm"
)

// PositionRepository implements the ports.PositionRepository interface
type PositionRepository struct {
	db *gorm.DB
}

// NewPositionRepository creates a new position repository instance
func NewPositionRepository(db *gorm.DB) ports.PositionRepository {
	return &PositionRepository{
		db: db,
	}
}

// Create persists a new position in the repository
func (r *PositionRepository) Create(ctx context.Context, position *models.Position) error {
	return r.db.WithContext(ctx).Create(position).Error
}

// GetByID retrieves a position by its ID
func (r *PositionRepository) GetByID(ctx context.Context, id string) (*models.Position, error) {
	var position models.Position
	if err := r.db.WithContext(ctx).First(&position, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

// List retrieves positions based on status
func (r *PositionRepository) List(ctx context.Context, status models.PositionStatus) ([]*models.Position, error) {
	var positions []*models.Position
	query := r.db.WithContext(ctx)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Find(&positions).Error; err != nil {
		return nil, err
	}
	
	return positions, nil
}

// Update updates an existing position in the repository
func (r *PositionRepository) Update(ctx context.Context, position *models.Position) error {
	return r.db.WithContext(ctx).Save(position).Error
}

// Delete removes a position from the repository
func (r *PositionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Position{}, "id = ?", id).Error
}

// GetOpenPositionBySymbol retrieves an open position for a specific symbol
func (r *PositionRepository) GetOpenPositionBySymbol(ctx context.Context, symbol string) (*models.Position, error) {
	var position models.Position
	err := r.db.WithContext(ctx).
		Where("symbol = ? AND status = ?", symbol, models.PositionStatusOpen).
		First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}
