package models

import (
	"encoding/json"
	"time"
)

// BacktestResult represents a stored backtest result in the database
type BacktestResult struct {
	ID                 string          `json:"id" gorm:"primaryKey"`
	UserID             string          `json:"user_id"`
	Strategy           string          `json:"strategy"`
	Symbol             string          `json:"symbol"`
	Timeframe          string          `json:"timeframe"`
	StartTime          time.Time       `json:"start_time"`
	EndTime            time.Time       `json:"end_time"`
	InitialCapital     float64         `json:"initial_capital"`
	FinalCapital       float64         `json:"final_capital"`
	TotalTrades        int             `json:"total_trades"`
	WinningTrades      int             `json:"winning_trades"`
	LosingTrades       int             `json:"losing_trades"`
	WinRate            float64         `json:"win_rate"`
	ProfitFactor       float64         `json:"profit_factor"`
	MaxDrawdown        float64         `json:"max_drawdown"`
	SharpeRatio        float64         `json:"sharpe_ratio"`
	EquityCurve        json.RawMessage `json:"equity_curve"`
	DrawdownCurve      json.RawMessage `json:"drawdown_curve"`
	Trades             json.RawMessage `json:"trades"`
	PerformanceMetrics json.RawMessage `json:"performance_metrics"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// TableName returns the table name for the BacktestResult model
func (BacktestResult) TableName() string {
	return "backtest_results"
}

// ToBacktestResult converts a backtest.BacktestResult to a database model
func ToBacktestResult(result interface{}, userID string) (*BacktestResult, error) {
	// Implementation will be added after backtest package is updated
	return nil, nil
}

// ToResponse converts the database model to a response format
func (br *BacktestResult) ToResponse() interface{} {
	// Implementation will be added after response types are defined
	return nil
}

// BacktestTrade represents a single trade executed during backtesting
type BacktestTrade struct {
	ID               string          `json:"id" gorm:"primaryKey"`
	BacktestResultID string          `json:"backtest_result_id" gorm:"index"`
	Symbol           string          `json:"symbol"`
	EntryTime        time.Time       `json:"entry_time"`
	ExitTime         time.Time       `json:"exit_time"`
	EntryPrice       float64         `json:"entry_price"`
	ExitPrice        float64         `json:"exit_price"`
	Size             float64         `json:"size"`
	ProfitLoss       float64         `json:"profit_loss"`
	ProfitLossPerc   float64         `json:"profit_loss_perc"`
	Fees             float64         `json:"fees"`
	SignalConfidence float64         `json:"signal_confidence"`
	Metadata         json.RawMessage `json:"metadata"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
