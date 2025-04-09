package models

import "time"

// ClosedPosition represents a trading position that has been closed
type ClosedPosition struct {
	ID                   string    `json:"id"`
	Symbol               string    `json:"symbol"`
	Side                 OrderSide `json:"side"` // BUY or SELL
	Quantity             float64   `json:"quantity"`
	Amount               float64   `json:"amount"` // Alias for Quantity for backward compatibility
	EntryPrice           float64   `json:"entry_price"`
	ExitPrice            float64   `json:"exit_price"`
	OpenTime             time.Time `json:"open_time"`
	CloseTime            time.Time `json:"close_time"`
	HoldingTimeMs        int64     `json:"holding_time_ms"` // Holding time in milliseconds
	ProfitLoss           float64   `json:"profit_loss"`
	Profit               float64   `json:"profit"` // Alias for ProfitLoss for backward compatibility
	ProfitLossPercentage float64   `json:"profit_loss_percentage"`
	ExitReason           string    `json:"exit_reason"` // e.g., "take_profit", "stop_loss", "manual"
	EntryReason          string    `json:"entry_reason,omitempty"`
	Strategy             string    `json:"strategy,omitempty"`
	InitialStopLoss      float64   `json:"initial_stop_loss,omitempty"`
	InitialTakeProfit    float64   `json:"initial_take_profit,omitempty"`
	RiskRewardRatio      float64   `json:"risk_reward_ratio,omitempty"`
	ActualRR             float64   `json:"actual_rr,omitempty"` // Actual risk/reward achieved
	ExpectedValue        float64   `json:"expected_value,omitempty"`
	MaxPrice             float64   `json:"max_price,omitempty"`    // Highest price during position lifetime
	MinPrice             float64   `json:"min_price,omitempty"`    // Lowest price during position lifetime
	MaxDrawdown          float64   `json:"max_drawdown,omitempty"` // Maximum drawdown during position lifetime
	MaxDrawdownPercent   float64   `json:"max_drawdown_percent,omitempty"`
	MaxProfit            float64   `json:"max_profit,omitempty"` // Maximum profit during position lifetime
	MaxProfitPercent     float64   `json:"max_profit_percent,omitempty"`
	Tags                 []string  `json:"tags,omitempty"`
	Notes                string    `json:"notes,omitempty"`
	Orders               []Order   `json:"orders,omitempty"` // All orders related to this position
}
