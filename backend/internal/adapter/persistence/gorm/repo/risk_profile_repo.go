package repo

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"gorm.io/gorm"
)

// RiskProfileEntity is the GORM entity for risk profiles
type RiskProfileEntity struct {
	ID                    string    `gorm:"column:id;primaryKey"`
	UserID                string    `gorm:"column:user_id;index;uniqueIndex"`
	MaxPositionSize       float64   `gorm:"column:max_position_size"`
	MaxTotalExposure      float64   `gorm:"column:max_total_exposure"`
	MaxDrawdown           float64   `gorm:"column:max_drawdown"`
	MaxLeverage           float64   `gorm:"column:max_leverage"`
	MaxConcentration      float64   `gorm:"column:max_concentration"`
	MinLiquidity          float64   `gorm:"column:min_liquidity"`
	VolatilityThreshold   float64   `gorm:"column:volatility_threshold"`
	DailyLossLimit        float64   `gorm:"column:daily_loss_limit"`
	WeeklyLossLimit       float64   `gorm:"column:weekly_loss_limit"`
	EnableAutoRiskControl bool      `gorm:"column:enable_auto_risk_control"`
	EnableNotifications   bool      `gorm:"column:enable_notifications"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

// TableName overrides the table name
func (RiskProfileEntity) TableName() string {
	return "risk_profiles"
}

// toEntity converts a risk profile model to entity
func toRiskProfileEntity(model *model.RiskProfile) *RiskProfileEntity {
	return &RiskProfileEntity{
		ID:                    model.ID,
		UserID:                model.UserID,
		MaxPositionSize:       model.MaxPositionSize,
		MaxTotalExposure:      model.MaxTotalExposure,
		MaxDrawdown:           model.MaxDrawdown,
		MaxLeverage:           model.MaxLeverage,
		MaxConcentration:      model.MaxConcentration,
		MinLiquidity:          model.MinLiquidity,
		VolatilityThreshold:   model.VolatilityThreshold,
		DailyLossLimit:        model.DailyLossLimit,
		WeeklyLossLimit:       model.WeeklyLossLimit,
		EnableAutoRiskControl: model.EnableAutoRiskControl,
		EnableNotifications:   model.EnableNotifications,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
	}
}

// toDomain converts a risk profile entity to domain model
func (e *RiskProfileEntity) toDomain() *model.RiskProfile {
	return &model.RiskProfile{
		ID:                    e.ID,
		UserID:                e.UserID,
		MaxPositionSize:       e.MaxPositionSize,
		MaxTotalExposure:      e.MaxTotalExposure,
		MaxDrawdown:           e.MaxDrawdown,
		MaxLeverage:           e.MaxLeverage,
		MaxConcentration:      e.MaxConcentration,
		MinLiquidity:          e.MinLiquidity,
		VolatilityThreshold:   e.VolatilityThreshold,
		DailyLossLimit:        e.DailyLossLimit,
		WeeklyLossLimit:       e.WeeklyLossLimit,
		EnableAutoRiskControl: e.EnableAutoRiskControl,
		EnableNotifications:   e.EnableNotifications,
		CreatedAt:             e.CreatedAt,
		UpdatedAt:             e.UpdatedAt,
	}
}

// GormRiskProfileRepository implements the RiskProfileRepository using GORM
type GormRiskProfileRepository struct {
	db *gorm.DB
}

// NewGormRiskProfileRepository creates a new instance of GormRiskProfileRepository
func NewGormRiskProfileRepository(db *gorm.DB) port.RiskProfileRepository {
	return &GormRiskProfileRepository{
		db: db,
	}
}

// Save creates or updates a risk profile
func (r *GormRiskProfileRepository) Save(ctx context.Context, profile *model.RiskProfile) error {
	entity := toRiskProfileEntity(profile)
	return r.db.WithContext(ctx).Save(entity).Error
}

// GetByUserID retrieves a risk profile for a specific user
func (r *GormRiskProfileRepository) GetByUserID(ctx context.Context, userID string) (*model.RiskProfile, error) {
	var entity RiskProfileEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return a default risk profile if none exists for this user
			return model.NewRiskProfile(userID), nil
		}
		return nil, err
	}
	return entity.toDomain(), nil
}

// Delete removes a risk profile
func (r *GormRiskProfileRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&RiskProfileEntity{}).Error
}

// Ensure GormRiskProfileRepository implements port.RiskProfileRepository
var _ port.RiskProfileRepository = (*GormRiskProfileRepository)(nil)
