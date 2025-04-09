package types

import (
	"time"
)

// BacktestRequest represents the parameters for running a backtest
type BacktestRequest struct {
	UserID         string                 `json:"userId" validate:"required"`
	Strategy       string                 `json:"strategy" validate:"required"`
	Symbol         string                 `json:"symbol" validate:"required"`
	Timeframe      string                 `json:"timeframe" validate:"required"`
	StartTime      time.Time              `json:"startTime" validate:"required"`
	EndTime        time.Time              `json:"endTime" validate:"required"`
	InitialCapital float64                `json:"initialCapital" validate:"required,gt=0"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
}

// BacktestTrade represents a trade executed during backtesting
type BacktestTrade struct {
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"` // "BUY" or "SELL"
	EntryPrice    float64   `json:"entryPrice"`
	ExitPrice     float64   `json:"exitPrice,omitempty"`
	Quantity      float64   `json:"quantity"`
	EntryTime     time.Time `json:"entryTime"`
	ExitTime      time.Time `json:"exitTime,omitempty"`
	ProfitLoss    float64   `json:"profitLoss"`
	ProfitLossPct float64   `json:"profitLossPct"`
}

// BacktestMetrics represents the performance metrics from a backtest
type BacktestMetrics struct {
	TotalTrades      int       `json:"totalTrades"`
	WinningTrades    int       `json:"winningTrades"`
	LosingTrades     int       `json:"losingTrades"`
	WinRate          float64   `json:"winRate"`
	AverageWin       float64   `json:"averageWin"`
	AverageLoss      float64   `json:"averageLoss"`
	LargestWin       float64   `json:"largestWin"`
	LargestLoss      float64   `json:"largestLoss"`
	ProfitFactor     float64   `json:"profitFactor"`
	SharpeRatio      float64   `json:"sharpeRatio"`
	MaxDrawdown      float64   `json:"maxDrawdown"`
	MaxDrawdownPct   float64   `json:"maxDrawdownPct"`
	AnnualizedReturn float64   `json:"annualizedReturn"`
	TotalReturn      float64   `json:"totalReturn"`
	TotalReturnPct   float64   `json:"totalReturnPct"`
	DailyReturns     []float64 `json:"dailyReturns,omitempty"`
}

// BacktestResult represents the complete results of a backtest
type BacktestResult struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"userId"`
	Strategy       string                 `json:"strategy"`
	Symbol         string                 `json:"symbol"`
	Timeframe      string                 `json:"timeframe"`
	StartTime      time.Time              `json:"startTime"`
	EndTime        time.Time              `json:"endTime"`
	InitialCapital float64                `json:"initialCapital"`
	FinalCapital   float64                `json:"finalCapital"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
	Metrics        BacktestMetrics        `json:"metrics"`
	Trades         []BacktestTrade        `json:"trades"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}
