package response

import "time"

// ConfigResponse represents the bot configuration
type ConfigResponse struct {
	USDTPerTrade     float64   `json:"usdt_per_trade"`
	StopLossPercent  float64   `json:"stop_loss_percent"`
	TakeProfitLevels []float64 `json:"take_profit_levels"`
	SellPercentages  []float64 `json:"sell_percentages"`
	UpdatedAt        time.Time `json:"updated_at"`
}
