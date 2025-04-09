package models

import "time"

type Position struct {
	ID               string            `json:"id"`
	Symbol           string            `json:"symbol"`
	Side             OrderSide         `json:"side"` // BUY or SELL
	Quantity         float64           `json:"quantity"`
	Amount           float64           `json:"amount"` // Alias for Quantity for backward compatibility
	EntryPrice       float64           `json:"entry_price"`
	CurrentPrice     float64           `json:"current_price"`
	OpenTime         time.Time         `json:"open_time"`
	OpenedAt         time.Time         `json:"opened_at"` // Alias for OpenTime for backward compatibility
	CloseTime        *time.Time        `json:"close_time,omitempty"`
	StopLoss         float64           `json:"stop_loss"`
	TakeProfit       float64           `json:"take_profit"`
	TrailingStop     *float64          `json:"trailing_stop,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	PnL              float64           `json:"pnl"`
	PnLPercentage    float64           `json:"pnl_percentage"`
	Status           string            `json:"status"` // open, closed
	Orders           []Order           `json:"orders"` // Entry and scaling orders
	EntryReason      string            `json:"entry_reason,omitempty"`
	ExitReason       string            `json:"exit_reason,omitempty"`
	Strategy         string            `json:"strategy,omitempty"`
	RiskRewardRatio  float64           `json:"risk_reward_ratio,omitempty"`
	ExpectedProfit   float64           `json:"expected_profit,omitempty"`
	MaxRisk          float64           `json:"max_risk,omitempty"`
	TakeProfitLevels []TakeProfitLevel `json:"take_profit_levels,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Notes            string            `json:"notes,omitempty"`
}
