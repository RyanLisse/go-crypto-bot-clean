package response

import "time"

// PortfolioSummaryResponse represents the overall portfolio status
type PortfolioSummaryResponse struct {
	TotalValue      float64         `json:"total_value"`
	ActiveTradeCount int            `json:"active_trade_count"`
	ActiveTrades    []TradeResponse `json:"active_trades"`
	Performance     PerformanceResponse `json:"performance"`
	Timestamp       time.Time       `json:"timestamp"`
}

// ActiveTradesResponse represents a list of active trades
type ActiveTradesResponse struct {
	Trades    []TradeResponse `json:"trades"`
	Count     int             `json:"count"`
	Timestamp time.Time       `json:"timestamp"`
}

// PerformanceResponse represents trading performance metrics
type PerformanceResponse struct {
	TotalTrades         int     `json:"total_trades"`
	WinningTrades       int     `json:"winning_trades"`
	LosingTrades        int     `json:"losing_trades"`
	WinRate             float64 `json:"win_rate"`
	TotalProfitLoss     float64 `json:"total_profit_loss"`
	AverageProfitPerTrade float64 `json:"average_profit_per_trade"`
	LargestProfit       float64 `json:"largest_profit"`
	LargestLoss         float64 `json:"largest_loss"`
	TimeRange           string  `json:"time_range"`
}

// TotalValueResponse represents the total portfolio value
type TotalValueResponse struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
