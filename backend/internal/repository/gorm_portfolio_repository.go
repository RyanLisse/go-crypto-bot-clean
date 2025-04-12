package repository

import (
	"context"

	"go-crypto-bot-clean/backend/internal/models"

	"gorm.io/gorm"
)

// GormPortfolioRepository implements the PortfolioRepository interface using GORM
type GormPortfolioRepository struct {
	db *gorm.DB
}

func NewGormPortfolioRepository(db *gorm.DB) *GormPortfolioRepository {
	return &GormPortfolioRepository{db: db}
}

func (r *GormPortfolioRepository) GetByUserID(ctx context.Context, userID string) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	err := r.db.Where("user_id = ?", userID).First(&portfolio).Error
	if err != nil {
		return nil, err
	}
	return &portfolio, nil
}

func (r *GormPortfolioRepository) AddPosition(ctx context.Context, position *models.Position) error {
	return r.db.Create(position).Error
}

func (r *GormPortfolioRepository) UpdatePosition(ctx context.Context, position *models.Position) error {
	return r.db.Save(position).Error
}

func (r *GormPortfolioRepository) GetPositionsByPortfolioID(ctx context.Context, portfolioID string) ([]*models.Position, error) {
	var positions []*models.Position
	err := r.db.Where("portfolio_id = ?", portfolioID).Find(&positions).Error
	if err != nil {
		return nil, err
	}
	return positions, nil
}
