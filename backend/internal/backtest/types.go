package backtest

import (
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/types"
)

// BacktestConfig represents the configuration for a backtest
type BacktestConfig struct {
	StartTime          time.Time
	EndTime            time.Time
	InitialCapital     float64
	Symbols            []string
	Interval           string
	CommissionRate     float64
	SlippageModel      SlippageModel
	EnableShortSelling bool
	DataProvider       DataProvider
	Strategy           types.Strategy
	Logger             Logger
}

// BacktestResult represents the result of a backtest run
type BacktestResult struct {
	Config             *BacktestConfig
	StartTime          time.Time
	EndTime            time.Time
	InitialCapital     float64
	FinalCapital       float64
	Trades             []*models.Trade
	Positions          []*models.Position
	ClosedPositions    []*models.ClosedPosition
	EquityCurve        []*EquityPoint
	DrawdownCurve      []*DrawdownPoint
	Events             []*BacktestEvent
	PerformanceMetrics *PerformanceMetrics
}

// EquityPoint represents a point in the equity curve
type EquityPoint struct {
	Timestamp time.Time
	Equity    float64
}

// DrawdownPoint represents a point in the drawdown curve
type DrawdownPoint struct {
	Timestamp time.Time
	Drawdown  float64
}

// BacktestEvent represents an event that occurred during backtesting
type BacktestEvent struct {
	Type      BacktestEventType
	Timestamp time.Time
	Symbol    string
	Data      interface{}
}

// BacktestEventType represents the type of backtest event
type BacktestEventType string

const (
	BacktestStarted    BacktestEventType = "started"
	BacktestCompleted  BacktestEventType = "completed"
	BacktestError      BacktestEventType = "error"
	SignalGenerated    BacktestEventType = "signal_generated"
	OrderCreated       BacktestEventType = "order_created"
	OrderFilled        BacktestEventType = "order_filled"
	PositionOpened     BacktestEventType = "position_opened"
	PositionClosed     BacktestEventType = "position_closed"
	EquityUpdated      BacktestEventType = "equity_updated"
	DrawdownCalculated BacktestEventType = "drawdown_calculated"
)

// BacktestRequestConfig represents the configuration for a backtest request
type BacktestRequestConfig struct {
	Strategy       string
	Symbol         string
	Timeframe      string
	StartTime      time.Time
	EndTime        time.Time
	InitialCapital float64
	RiskPerTrade   float64
}
