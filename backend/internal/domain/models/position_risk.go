package models

// PositionRisk represents the risk assessment for a trading position
type PositionRisk struct {
	Symbol      string  // Trading pair symbol
	ExposureUSD float64 // Current exposure in USD
	RiskLevel   string  // Risk level assessment (e.g., "LOW", "MEDIUM", "HIGH")
}
