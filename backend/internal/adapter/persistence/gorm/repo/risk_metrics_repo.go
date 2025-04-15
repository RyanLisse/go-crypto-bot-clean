package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"gorm.io/gorm"
)

// RiskMetricsEntity is the GORM entity for risk metrics
type RiskMetricsEntity struct {
	ID                   string    `gorm:"column:id;primaryKey"`
	UserID               string    `gorm:"column:user_id;index"`
	Date                 time.Time `gorm:"column:date;index"`
	PortfolioValue       float64   `gorm:"column:portfolio_value"`
	TotalExposure        float64   `gorm:"column:total_exposure"`
	MaxDrawdown          float64   `gorm:"column:max_drawdown"`
	DailyPnL             float64   `gorm:"column:daily_pnl"`
	WeeklyPnL            float64   `gorm:"column:weekly_pnl"`
	MonthlyPnL           float64   `gorm:"column:monthly_pnl"`
	HighestConcentration float64   `gorm:"column:highest_concentration"`
	VolatilityScore      float64   `gorm:"column:volatility_score"`
	LiquidityScore       float64   `gorm:"column:liquidity_score"`
	OverallRiskScore     float64   `gorm:"column:overall_risk_score"`
	AdditionalDataJSON   string    `gorm:"column:additional_data_json;type:text"` // JSON string of additional data
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

// TableName overrides the table name
func (RiskMetricsEntity) TableName() string {
	return "risk_metrics"
}

// toRiskMetricsEntity converts a risk metrics model to entity
func toRiskMetricsEntity(metrics *model.RiskMetrics) (*RiskMetricsEntity, error) {
	var additionalDataJSON string
	if metrics.AdditionalData != nil {
		data, err := json.Marshal(metrics.AdditionalData)
		if err != nil {
			return nil, err
		}
		additionalDataJSON = string(data)
	}

	return &RiskMetricsEntity{
		ID:                   metrics.ID,
		UserID:               metrics.UserID,
		Date:                 metrics.Date,
		PortfolioValue:       metrics.PortfolioValue,
		TotalExposure:        metrics.TotalExposure,
		MaxDrawdown:          metrics.MaxDrawdown,
		DailyPnL:             metrics.DailyPnL,
		WeeklyPnL:            metrics.WeeklyPnL,
		MonthlyPnL:           metrics.MonthlyPnL,
		HighestConcentration: metrics.HighestConcentration,
		VolatilityScore:      metrics.VolatilityScore,
		LiquidityScore:       metrics.LiquidityScore,
		OverallRiskScore:     metrics.OverallRiskScore,
		AdditionalDataJSON:   additionalDataJSON,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}

// toDomain converts a risk metrics entity to domain model
func (e *RiskMetricsEntity) toDomain() (*model.RiskMetrics, error) {
	var additionalData map[string]interface{}
	if e.AdditionalDataJSON != "" {
		if err := json.Unmarshal([]byte(e.AdditionalDataJSON), &additionalData); err != nil {
			return nil, err
		}
	}

	return &model.RiskMetrics{
		ID:                   e.ID,
		UserID:               e.UserID,
		Date:                 e.Date,
		PortfolioValue:       e.PortfolioValue,
		TotalExposure:        e.TotalExposure,
		MaxDrawdown:          e.MaxDrawdown,
		DailyPnL:             e.DailyPnL,
		WeeklyPnL:            e.WeeklyPnL,
		MonthlyPnL:           e.MonthlyPnL,
		HighestConcentration: e.HighestConcentration,
		VolatilityScore:      e.VolatilityScore,
		LiquidityScore:       e.LiquidityScore,
		OverallRiskScore:     e.OverallRiskScore,
		AdditionalData:       additionalData,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
	}, nil
}

// GormRiskMetricsRepository implements the RiskMetricsRepository using GORM
type GormRiskMetricsRepository struct {
	db *gorm.DB
}

// NewGormRiskMetricsRepository creates a new instance of GormRiskMetricsRepository
func NewGormRiskMetricsRepository(db *gorm.DB) *GormRiskMetricsRepository {
	return &GormRiskMetricsRepository{
		db: db,
	}
}

// Save creates or updates risk metrics
func (r *GormRiskMetricsRepository) Save(ctx context.Context, metrics *model.RiskMetrics) error {
	entity, err := toRiskMetricsEntity(metrics)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(entity).Error
}

// GetByID retrieves risk metrics by ID
func (r *GormRiskMetricsRepository) GetByID(ctx context.Context, id string) (*model.RiskMetrics, error) {
	var entity RiskMetricsEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return nil, err
	}
	return entity.toDomain()
}

// GetLatestByUserID retrieves the latest risk metrics for a user
func (r *GormRiskMetricsRepository) GetLatestByUserID(ctx context.Context, userID string) (*model.RiskMetrics, error) {
	var entity RiskMetricsEntity
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		First(&entity).Error; err != nil {
		return nil, err
	}
	return entity.toDomain()
}

// GetByUserID retrieves risk metrics for a specific user
func (r *GormRiskMetricsRepository) GetByUserID(ctx context.Context, userID string) (*model.RiskMetrics, error) {
	// This is an alias for GetLatestByUserID since we typically want the most recent metrics
	return r.GetLatestByUserID(ctx, userID)
}

// GetByUserIDAndDateRange retrieves risk metrics for a user within a date range
func (r *GormRiskMetricsRepository) GetByUserIDAndDateRange(
	ctx context.Context,
	userID string,
	startDate, endDate time.Time,
) ([]*model.RiskMetrics, error) {
	var entities []RiskMetricsEntity
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date ASC").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	metrics := make([]*model.RiskMetrics, 0, len(entities))
	for _, entity := range entities {
		m, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// GetByUserIDAndPeriod retrieves risk metrics for a user for a specific period
func (r *GormRiskMetricsRepository) GetByUserIDAndPeriod(
	ctx context.Context,
	userID string,
	period string, // "daily", "weekly", "monthly"
	limit int,
) ([]*model.RiskMetrics, error) {
	var entities []RiskMetricsEntity
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	switch period {
	case "daily":
		// Daily metrics - no additional filter needed
	case "weekly":
		// Filter for weekly metrics (e.g., every Monday)
		query = query.Where("EXTRACT(DOW FROM date) = 1")
	case "monthly":
		// Filter for monthly metrics (e.g., 1st day of month)
		query = query.Where("EXTRACT(DAY FROM date) = 1")
	default:
		// Default to daily if period is not recognized
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("date DESC").Find(&entities).Error; err != nil {
		return nil, err
	}

	metrics := make([]*model.RiskMetrics, 0, len(entities))
	for _, entity := range entities {
		m, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// Delete removes risk metrics by ID
func (r *GormRiskMetricsRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&RiskMetricsEntity{}, "id = ?", id).Error
}

// DeleteOlderThan removes risk metrics older than the specified date
func (r *GormRiskMetricsRepository) DeleteOlderThan(ctx context.Context, date time.Time) error {
	return r.db.WithContext(ctx).
		Delete(&RiskMetricsEntity{}, "date < ?", date).
		Error
}

// GetHistorical retrieves historical risk metrics for a user within a time range
func (r *GormRiskMetricsRepository) GetHistorical(ctx context.Context, userID string, from, to time.Time, interval string) ([]*model.RiskMetrics, error) {
	// This is an alias for GetByUserIDAndDateRange with some additional interval handling
	var entities []RiskMetricsEntity
	query := r.db.WithContext(ctx).Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to)

	// Apply interval filtering if specified
	switch interval {
	case "daily":
		// No additional filtering for daily
	case "weekly":
		// Filter for weekly data points
		query = query.Where("EXTRACT(DOW FROM date) = 1")
	case "monthly":
		// Filter for monthly data points
		query = query.Where("EXTRACT(DAY FROM date) = 1")
	}

	// Order by date ascending
	if err := query.Order("date ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	// Convert entities to domain models
	metrics := make([]*model.RiskMetrics, 0, len(entities))
	for _, entity := range entities {
		m, err := entity.toDomain()
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

// Ensure GormRiskMetricsRepository implements port.RiskMetricsRepository
var _ port.RiskMetricsRepository = (*GormRiskMetricsRepository)(nil)
