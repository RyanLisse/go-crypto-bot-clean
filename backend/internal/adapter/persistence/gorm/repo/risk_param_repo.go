package repo

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RiskParameterEntity is the GORM entity for risk parameters
type RiskParameterEntity struct {
	ID                              uint      `gorm:"primaryKey"`
	UserID                          string    `gorm:"column:user_id;index"`
	MaxConcentrationPercentage      float64   `gorm:"column:max_concentration_percentage"`
	MinLiquidityThresholdUSD        float64   `gorm:"column:min_liquidity_threshold_usd"`
	MaxPositionSizePercentage       float64   `gorm:"column:max_position_size_percentage"`
	MaxDrawdownPercentage           float64   `gorm:"column:max_drawdown_percentage"`
	VolatilityMultiplier            float64   `gorm:"column:volatility_multiplier"`
	DefaultMaxConcentrationPct      float64   `gorm:"column:default_max_concentration_pct"`
	DefaultMaxPositionSizePct       float64   `gorm:"column:default_max_position_size_pct"`
	DefaultMinLiquidityThresholdUSD float64   `gorm:"column:default_min_liquidity_threshold_usd"`
	DefaultMaxDrawdownPct           float64   `gorm:"column:default_max_drawdown_pct"`
	DefaultVolatilityMultiplier     float64   `gorm:"column:default_volatility_multiplier"`
	CreatedAt                       time.Time `gorm:"column:created_at"`
	UpdatedAt                       time.Time `gorm:"column:updated_at"`
}

// TableName overrides the table name
func (RiskParameterEntity) TableName() string {
	return "risk_parameters"
}

// toEntity converts a risk parameter model to entity
func toEntity(model *model.RiskParameters) *RiskParameterEntity {
	return &RiskParameterEntity{
		UserID:                          model.UserID,
		MaxConcentrationPercentage:      model.MaxConcentrationPercentage,
		MinLiquidityThresholdUSD:        model.MinLiquidityThresholdUSD,
		MaxPositionSizePercentage:       model.MaxPositionSizePercentage,
		MaxDrawdownPercentage:           model.MaxDrawdownPercentage,
		VolatilityMultiplier:            model.VolatilityMultiplier,
		DefaultMaxConcentrationPct:      model.DefaultMaxConcentrationPct,
		DefaultMaxPositionSizePct:       model.DefaultMaxPositionSizePct,
		DefaultMinLiquidityThresholdUSD: model.DefaultMinLiquidityThresholdUSD,
		DefaultMaxDrawdownPct:           model.DefaultMaxDrawdownPct,
		DefaultVolatilityMultiplier:     model.DefaultVolatilityMultiplier,
	}
}

// toDomain converts a risk parameter entity to domain model
func (e *RiskParameterEntity) toDomain() *model.RiskParameters {
	return &model.RiskParameters{
		UserID:                          e.UserID,
		MaxConcentrationPercentage:      e.MaxConcentrationPercentage,
		MinLiquidityThresholdUSD:        e.MinLiquidityThresholdUSD,
		MaxPositionSizePercentage:       e.MaxPositionSizePercentage,
		MaxDrawdownPercentage:           e.MaxDrawdownPercentage,
		VolatilityMultiplier:            e.VolatilityMultiplier,
		DefaultMaxConcentrationPct:      e.DefaultMaxConcentrationPct,
		DefaultMaxPositionSizePct:       e.DefaultMaxPositionSizePct,
		DefaultMinLiquidityThresholdUSD: e.DefaultMinLiquidityThresholdUSD,
		DefaultMaxDrawdownPct:           e.DefaultMaxDrawdownPct,
		DefaultVolatilityMultiplier:     e.DefaultVolatilityMultiplier,
	}
}

// GormRiskParameterRepository implements the RiskParameterRepository using GORM
type GormRiskParameterRepository struct {
	db *gorm.DB
}

// NewGormRiskParameterRepository creates a new instance of GormRiskParameterRepository
func NewGormRiskParameterRepository(db *gorm.DB) *GormRiskParameterRepository {
	return &GormRiskParameterRepository{
		db: db,
	}
}

// GetParameters retrieves risk parameters for a user
func (r *GormRiskParameterRepository) GetParameters(ctx context.Context, userID string) (*model.RiskParameters, error) {
	var entity RiskParameterEntity
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return default parameters if none exist for this user
			return &model.RiskParameters{
				UserID:                          userID,
				MaxConcentrationPercentage:      30.0,   // Default 30% max concentration
				MinLiquidityThresholdUSD:        100000, // Default $100k min liquidity
				MaxPositionSizePercentage:       10.0,   // Default 10% max position size
				MaxDrawdownPercentage:           20.0,   // Default 20% max drawdown
				VolatilityMultiplier:            1.5,    // Default volatility multiplier
				DefaultMaxConcentrationPct:      30.0,
				DefaultMaxPositionSizePct:       10.0,
				DefaultMinLiquidityThresholdUSD: 100000,
				DefaultMaxDrawdownPct:           20.0,
				DefaultVolatilityMultiplier:     1.5,
			}, nil
		}
		return nil, result.Error
	}
	return entity.toDomain(), nil
}

// SaveParameters saves risk parameters for a user
func (r *GormRiskParameterRepository) SaveParameters(ctx context.Context, params *model.RiskParameters) error {
	entity := toEntity(params)

	// Use a transaction to perform the upsert
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

// Ensure GormRiskParameterRepository implements port.RiskParameterRepository
var _ port.RiskParameterRepository = (*GormRiskParameterRepository)(nil)
