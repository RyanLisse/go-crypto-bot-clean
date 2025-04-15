package repo

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RiskProfileEntity is the GORM entity for risk profile
type RiskProfileEntity struct {
	ID                    string    `gorm:"column:id;primaryKey"`
	UserID                string    `gorm:"column:user_id;index;unique"`
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

// toRiskProfileEntity converts a risk profile model to entity
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
func NewGormRiskProfileRepository(db *gorm.DB) *GormRiskProfileRepository {
	return &GormRiskProfileRepository{
		db: db,
	}
}

// Save creates or updates a risk profile
func (r *GormRiskProfileRepository) Save(ctx context.Context, profile *model.RiskProfile) error {
	entity := toRiskProfileEntity(profile)
	entity.UpdatedAt = time.Now()

	// Use a transaction for the operation
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Try to insert or update based on user_id
	result := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(entity)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

// GetByUserID retrieves a risk profile for a specific user
func (r *GormRiskProfileRepository) GetByUserID(ctx context.Context, userID string) (*model.RiskProfile, error) {
	var entity RiskProfileEntity
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return a default risk profile if none exists
			return model.NewRiskProfile(userID), nil
		}
		return nil, result.Error
	}
	return entity.toDomain(), nil
}

// Delete removes a risk profile
func (r *GormRiskProfileRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&RiskProfileEntity{}, "id = ?", id)
	return result.Error
}

// Ensure GormRiskProfileRepository implements port.RiskProfileRepository
var _ port.RiskProfileRepository = (*GormRiskProfileRepository)(nil)
