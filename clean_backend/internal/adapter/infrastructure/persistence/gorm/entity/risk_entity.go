package entity

import (
	"time"
)

// RiskAssessmentEntity is the GORM entity for risk assessment
type RiskAssessmentEntity struct {
	ID             string     `gorm:"column:id;primaryKey"`
	UserID         string     `gorm:"column:user_id;index"`
	Type           string     `gorm:"column:type;index"`   // RiskType as string
	Level          string     `gorm:"column:level;index"`  // RiskLevel as string
	Status         string     `gorm:"column:status;index"` // RiskStatus as string
	Symbol         string     `gorm:"column:symbol;index"`
	PositionID     string     `gorm:"column:position_id;index"`
	OrderID        string     `gorm:"column:order_id;index"`
	Score          float64    `gorm:"column:score"`
	Message        string     `gorm:"column:message;type:text"`
	Recommendation string     `gorm:"column:recommendation;type:text"`
	MetadataJSON   string     `gorm:"column:metadata_json;type:text"` // JSON string of metadata
	CreatedAt      time.Time  `gorm:"column:created_at;index"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	ResolvedAt     *time.Time `gorm:"column:resolved_at"`
}

// TableName overrides the table name
func (RiskAssessmentEntity) TableName() string {
	return "risk_assessments"
}

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
