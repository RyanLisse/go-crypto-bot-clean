package model

// RiskParameters represents risk management parameters for a user
type RiskParameters struct {
	UserID                          string  `json:"user_id"`
	MaxConcentrationPercentage      float64 `json:"max_concentration_percentage"`
	MinLiquidityThresholdUSD        float64 `json:"min_liquidity_threshold_usd"`
	MaxPositionSizePercentage       float64 `json:"max_position_size_percentage"`
	MaxDrawdownPercentage           float64 `json:"max_drawdown_percentage"`
	VolatilityMultiplier            float64 `json:"volatility_multiplier"`
	DefaultMaxConcentrationPct      float64 `json:"default_max_concentration_pct"`
	DefaultMaxPositionSizePct       float64 `json:"default_max_position_size_pct"`
	DefaultMinLiquidityThresholdUSD float64 `json:"default_min_liquidity_threshold_usd"`
	DefaultMaxDrawdownPct           float64 `json:"default_max_drawdown_pct"`
	DefaultVolatilityMultiplier     float64 `json:"default_volatility_multiplier"`
}

// NewDefaultRiskParameters creates a new set of default risk parameters
func NewDefaultRiskParameters(userID string) *RiskParameters {
	return &RiskParameters{
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
	}
}
