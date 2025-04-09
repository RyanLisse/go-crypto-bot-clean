package models

import (
	"time"
)

// TimeFrame represents a time period for analytics
type TimeFrame string

const (
	TimeFrameDay     TimeFrame = "DAY"
	TimeFrameWeek    TimeFrame = "WEEK"
	TimeFrameMonth   TimeFrame = "MONTH"
	TimeFrameQuarter TimeFrame = "QUARTER"
	TimeFrameYear    TimeFrame = "YEAR"
	TimeFrameAll     TimeFrame = "ALL"
)

// TradeAnalytics represents analytics data for trading performance
type TradeAnalytics struct {
	TimeFrame           TimeFrame `json:"time_frame"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	
	// Overall performance
	TotalTrades         int       `json:"total_trades"`
	WinningTrades       int       `json:"winning_trades"`
	LosingTrades        int       `json:"losing_trades"`
	WinRate             float64   `json:"win_rate"`
	TotalProfit         float64   `json:"total_profit"`
	TotalLoss           float64   `json:"total_loss"`
	NetProfit           float64   `json:"net_profit"`
	ProfitFactor        float64   `json:"profit_factor"`
	AverageProfit       float64   `json:"average_profit"`
	AverageLoss         float64   `json:"average_loss"`
	LargestProfit       float64   `json:"largest_profit"`
	LargestLoss         float64   `json:"largest_loss"`
	
	// Risk metrics
	MaxDrawdown         float64   `json:"max_drawdown"`
	MaxDrawdownPercent  float64   `json:"max_drawdown_percent"`
	SharpeRatio         float64   `json:"sharpe_ratio"`
	SortinoRatio        float64   `json:"sortino_ratio"`
	RiskRewardRatio     float64   `json:"risk_reward_ratio"`
	
	// Time metrics
	AverageHoldingTime  string    `json:"average_holding_time"`
	AverageHoldingTimeWinning string `json:"average_holding_time_winning"`
	AverageHoldingTimeLosing  string `json:"average_holding_time_losing"`
	
	// Trade frequency
	TradesPerDay        float64   `json:"trades_per_day"`
	TradesPerWeek       float64   `json:"trades_per_week"`
	TradesPerMonth      float64   `json:"trades_per_month"`
	
	// Performance by reason
	PerformanceByReason map[string]ReasonPerformance `json:"performance_by_reason"`
	
	// Performance by symbol
	PerformanceBySymbol map[string]SymbolPerformance `json:"performance_by_symbol"`
	
	// Performance by strategy
	PerformanceByStrategy map[string]StrategyPerformance `json:"performance_by_strategy"`
	
	// Balance history
	BalanceHistory     []BalancePoint `json:"balance_history"`
	
	// Equity curve
	EquityCurve        []EquityPoint  `json:"equity_curve"`
}

// ReasonPerformance tracks performance metrics for a specific decision reason
type ReasonPerformance struct {
	Reason          string  `json:"reason"`
	TotalTrades     int     `json:"total_trades"`
	WinningTrades   int     `json:"winning_trades"`
	LosingTrades    int     `json:"losing_trades"`
	WinRate         float64 `json:"win_rate"`
	TotalProfit     float64 `json:"total_profit"`
	AverageProfit   float64 `json:"average_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
}

// SymbolPerformance tracks performance metrics for a specific trading symbol
type SymbolPerformance struct {
	Symbol          string  `json:"symbol"`
	TotalTrades     int     `json:"total_trades"`
	WinningTrades   int     `json:"winning_trades"`
	LosingTrades    int     `json:"losing_trades"`
	WinRate         float64 `json:"win_rate"`
	TotalProfit     float64 `json:"total_profit"`
	AverageProfit   float64 `json:"average_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
}

// StrategyPerformance tracks performance metrics for a specific trading strategy
type StrategyPerformance struct {
	Strategy        string  `json:"strategy"`
	TotalTrades     int     `json:"total_trades"`
	WinningTrades   int     `json:"winning_trades"`
	LosingTrades    int     `json:"losing_trades"`
	WinRate         float64 `json:"win_rate"`
	TotalProfit     float64 `json:"total_profit"`
	AverageProfit   float64 `json:"average_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
}

// BalancePoint represents a point in the balance history
type BalancePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Balance   float64   `json:"balance"`
}

// EquityPoint represents a point in the equity curve
type EquityPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Equity    float64   `json:"equity"`
}

// TradePerformance represents the performance of a single trade
type TradePerformance struct {
	TradeID         string    `json:"trade_id"`
	Symbol          string    `json:"symbol"`
	EntryTime       time.Time `json:"entry_time"`
	ExitTime        time.Time `json:"exit_time"`
	EntryPrice      float64   `json:"entry_price"`
	ExitPrice       float64   `json:"exit_price"`
	Quantity        float64   `json:"quantity"`
	ProfitLoss      float64   `json:"profit_loss"`
	ProfitLossPercent float64 `json:"profit_loss_percent"`
	HoldingTime     string    `json:"holding_time"`
	HoldingTimeMs   int64     `json:"holding_time_ms"`
	EntryReason     string    `json:"entry_reason"`
	ExitReason      string    `json:"exit_reason"`
	Strategy        string    `json:"strategy"`
	StopLoss        float64   `json:"stop_loss"`
	TakeProfit      float64   `json:"take_profit"`
	RiskRewardRatio float64   `json:"risk_reward_ratio"`
	ExpectedValue   float64   `json:"expected_value"`
	ActualRR        float64   `json:"actual_rr"`
}
