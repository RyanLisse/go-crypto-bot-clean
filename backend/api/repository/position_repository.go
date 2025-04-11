package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-crypto-bot-clean/backend/api/models"
)

type PositionRepository interface {
	Create(ctx context.Context, position *models.Position) error
	GetByID(ctx context.Context, id string) (*models.Position, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Position, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Position, error)
	Update(ctx context.Context, position *models.Position) error
	Delete(ctx context.Context, id string) error

	UpdatePrice(ctx context.Context, id string, newPrice float64) error
	MarkClosed(ctx context.Context, id string, closePrice float64, closeTime time.Time) error
}

type GormPositionRepository struct {
	db *gorm.DB
}

func NewGormPositionRepository(db *gorm.DB) PositionRepository {
	return &GormPositionRepository{db: db}
}

func (r *GormPositionRepository) Create(ctx context.Context, position *models.Position) error {
	return r.db.WithContext(ctx).Create(position).Error
}

func (r *GormPositionRepository) GetByID(ctx context.Context, id string) (*models.Position, error) {
	var position models.Position
	if err := r.db.WithContext(ctx).First(&position, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *GormPositionRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Position, error) {
	var positions []*models.Position
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *GormPositionRepository) GetByStatus(ctx context.Context, status string) ([]*models.Position, error) {
	var positions []*models.Position
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *GormPositionRepository) Update(ctx context.Context, position *models.Position) error {
	return r.db.WithContext(ctx).Save(position).Error
}

func (r *GormPositionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Position{}).Error
}
func (r *GormPositionRepository) UpdatePrice(ctx context.Context, id string, newPrice float64) error {
	return r.db.WithContext(ctx).Model(&models.Position{}).
		Where("id = ?", id).
		Update("current_price", newPrice).Error
}

func (r *GormPositionRepository) MarkClosed(ctx context.Context, id string, closePrice float64, closeTime time.Time) error {
	return r.db.WithContext(ctx).Model(&models.Position{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "closed",
			"close_price": closePrice,
			"close_time":  closeTime,
		}).Error
}
