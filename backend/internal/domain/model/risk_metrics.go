package model

import (
	"time"

	"github.com/google/uuid"
)

// RiskMetrics represents aggregated risk metrics for a user
type RiskMetrics struct {
	ID                   string                 `json:"id" gorm:"primaryKey"`
	UserID               string                 `json:"user_id" gorm:"index"`
	Date                 time.Time              `json:"date" gorm:"index"`
	PortfolioValue       float64                `json:"portfolio_value"`
	TotalExposure        float64                `json:"total_exposure"`         // Total exposure in quote currency
	MaxDrawdown          float64                `json:"max_drawdown"`           // Maximum drawdown
	DailyPnL             float64                `json:"daily_pnl"`              // Daily profit/loss
	WeeklyPnL            float64                `json:"weekly_pnl"`             // Weekly profit/loss
	MonthlyPnL           float64                `json:"monthly_pnl"`            // Monthly profit/loss
	HighestConcentration float64                `json:"highest_concentration"`  // Highest concentration in a single asset
	VolatilityScore      float64                `json:"volatility_score"`       // Volatility score
	LiquidityScore       float64                `json:"liquidity_score"`        // Liquidity score
	OverallRiskScore     float64                `json:"overall_risk_score"`     // Overall risk score
	CurrentDrawdown      float64                `json:"current_drawdown"`       // Current drawdown from peak
	DailyPnLPercent      float64                `json:"daily_pnl_percent"`      // Daily profit/loss as percentage
	PortfolioVolatility  float64                `json:"portfolio_volatility"`   // Portfolio volatility
	ActiveRiskCount      int                    `json:"active_risk_count"`      // Number of active risks
	HighRiskCount        int                    `json:"high_risk_count"`        // Number of high/critical risks
	Timestamp            time.Time              `json:"timestamp" gorm:"index"` // When these metrics were calculated
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	AdditionalData       map[string]interface{} `json:"additional_data,omitempty"`
}

// NewRiskMetrics creates a new risk metrics record
func NewRiskMetrics(userID string) *RiskMetrics {
	now := time.Now()
	return &RiskMetrics{
		ID:        uuid.New().String(),
		UserID:    userID,
		Timestamp: now,
		CreatedAt: now,
		Date:      now,
	}
}
