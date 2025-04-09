package request

// ConfigUpdateRequest represents a request to update bot configuration
type ConfigUpdateRequest struct {
	USDTPerTrade    *float64   `json:"usdt_per_trade,omitempty"`
	StopLossPercent *float64   `json:"stop_loss_percent,omitempty"`
	TakeProfitLevels []float64 `json:"take_profit_levels,omitempty"`
	SellPercentages  []float64 `json:"sell_percentages,omitempty"`
}
