package model

import (
	"time"

	"github.com/google/uuid"
)

// RiskLevel represents the severity level of a risk
type RiskLevel string

const (
	// RiskLevelLow represents a low risk level
	RiskLevelLow RiskLevel = "LOW"
	// RiskLevelMedium represents a medium risk level
	RiskLevelMedium RiskLevel = "MEDIUM"
	// RiskLevelHigh represents a high risk level
	RiskLevelHigh RiskLevel = "HIGH"
	// RiskLevelCritical represents a critical risk level
	RiskLevelCritical RiskLevel = "CRITICAL"
)

// RiskType represents the type of risk being managed
type RiskType string

const (
	// RiskTypePosition is for position-related risks
	RiskTypePosition RiskType = "POSITION"
	// RiskTypeVolatility is for market volatility risks
	RiskTypeVolatility RiskType = "VOLATILITY"
	// RiskTypeLiquidity is for market liquidity risks
	RiskTypeLiquidity RiskType = "LIQUIDITY"
	// RiskTypeExposure is for total market exposure risks
	RiskTypeExposure RiskType = "EXPOSURE"
	// RiskTypeConcentration is for concentration risks (too much in one asset)
	RiskTypeConcentration RiskType = "CONCENTRATION"
	// RiskTypeDrawdown is for drawdown risks in positions
	RiskTypeDrawdown RiskType = "DRAWDOWN"
)

// RiskStatus represents the current status of a risk assessment
type RiskStatus string

const (
	// RiskStatusActive means the risk assessment is currently active
	RiskStatusActive RiskStatus = "ACTIVE"
	// RiskStatusResolved means the risk has been resolved
	RiskStatusResolved RiskStatus = "RESOLVED"
	// RiskStatusIgnored means the risk has been acknowledged but ignored
	RiskStatusIgnored RiskStatus = "IGNORED"
)

// RiskProfile defines the risk tolerance and parameters for a user
type RiskProfile struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"userId"`
	MaxPositionSize       float64   `json:"maxPositionSize"`       // Maximum size of a single position in quote currency
	MaxTotalExposure      float64   `json:"maxTotalExposure"`      // Maximum total exposure across all positions
	MaxDrawdown           float64   `json:"maxDrawdown"`           // Maximum allowed drawdown percentage (0-1)
	MaxLeverage           float64   `json:"maxLeverage"`           // Maximum allowed leverage
	MaxConcentration      float64   `json:"maxConcentration"`      // Maximum portfolio concentration in a single asset (0-1)
	MinLiquidity          float64   `json:"minLiquidity"`          // Minimum required liquidity for trading
	VolatilityThreshold   float64   `json:"volatilityThreshold"`   // Threshold for high volatility warning
	DailyLossLimit        float64   `json:"dailyLossLimit"`        // Maximum allowed loss in a day
	WeeklyLossLimit       float64   `json:"weeklyLossLimit"`       // Maximum allowed loss in a week
	EnableAutoRiskControl bool      `json:"enableAutoRiskControl"` // Whether to enable automatic risk control
	EnableNotifications   bool      `json:"enableNotifications"`   // Whether to enable risk notifications
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// NewRiskProfile creates a new risk profile with default values
func NewRiskProfile(userID string) *RiskProfile {
	now := time.Now()
	return &RiskProfile{
		ID:                    uuid.New().String(),
		UserID:                userID,
		MaxPositionSize:       1000.0,  // Default $1000 max position
		MaxTotalExposure:      5000.0,  // Default $5000 max total exposure
		MaxDrawdown:           0.1,     // Default 10% max drawdown
		MaxLeverage:           3.0,     // Default 3x max leverage
		MaxConcentration:      0.2,     // Default 20% max concentration in one asset
		MinLiquidity:          10000.0, // Default $10000 min liquidity
		VolatilityThreshold:   0.05,    // Default 5% volatility threshold
		DailyLossLimit:        100.0,   // Default $100 daily loss limit
		WeeklyLossLimit:       500.0,   // Default $500 weekly loss limit
		EnableAutoRiskControl: true,
		EnableNotifications:   true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

// RiskAssessment represents a single risk assessment event
type RiskAssessment struct {
	ID             string      `json:"id"`
	UserID         string      `json:"userId"`
	Type           RiskType    `json:"type"`
	Level          RiskLevel   `json:"level"`
	Status         RiskStatus  `json:"status"`
	Symbol         string      `json:"symbol,omitempty"`     // Optional, for symbol-specific risks
	PositionID     string      `json:"positionId,omitempty"` // Optional, for position-specific risks
	OrderID        string      `json:"orderId,omitempty"`    // Optional, for order-specific risks
	Score          float64     `json:"score"`                // Numerical risk score (0-100)
	Message        string      `json:"message"`              // Human-readable description of the risk
	Recommendation string      `json:"recommendation"`       // Recommended action to mitigate the risk
	Metadata       interface{} `json:"metadata,omitempty"`   // Additional context-specific data
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`            // Time of last update
	ResolvedAt     *time.Time  `json:"resolvedAt,omitempty"` // When the risk was resolved, if applicable
}

// NewRiskAssessment creates a new risk assessment
func NewRiskAssessment(userID string, riskType RiskType, level RiskLevel, message string) *RiskAssessment {
	now := time.Now()
	return &RiskAssessment{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      riskType,
		Level:     level,
		Status:    RiskStatusActive,
		Message:   message,
		Score:     0, // Default score
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Resolve marks a risk assessment as resolved
func (r *RiskAssessment) Resolve() {
	r.Status = RiskStatusResolved
	now := time.Now()
	r.ResolvedAt = &now
	r.UpdatedAt = now
}

// Ignore marks a risk assessment as ignored
func (r *RiskAssessment) Ignore() {
	r.Status = RiskStatusIgnored
	r.UpdatedAt = time.Now()
}
