package backtest

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/backtest/types"
	"go-crypto-bot-clean/backend/internal/domain/models"

	"go.uber.org/zap"
)

// BacktestEventType represents different types of events during backtesting
type BacktestEventType int

const (
	// BacktestStarted indicates the backtest has started
	BacktestStarted BacktestEventType = iota
	// BacktestCompleted indicates the backtest has completed
	BacktestCompleted
	// BacktestError indicates an error occurred during backtesting
	BacktestError
	// BacktestProgress indicates progress update during backtesting
	BacktestProgress
)

// BacktestEvent represents an event that occurred during the backtest
type BacktestEvent struct {
	Type      BacktestEventType
	Timestamp time.Time
	Symbol    string
	Data      interface{}
}

// BacktestConfig contains configuration options for a backtest
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
	Strategy           BacktestStrategy
	Logger             *zap.Logger
}

// BacktestResult contains the results of a backtest
type BacktestResult struct {
	Config             *BacktestConfig
	StartTime          time.Time
	EndTime            time.Time
	InitialCapital     float64
	FinalCapital       float64
	Trades             []*models.Order
	Positions          []*models.Position
	ClosedPositions    []*models.ClosedPosition
	EquityCurve        []*EquityPoint
	DrawdownCurve      []*DrawdownPoint
	Events             []*BacktestEvent
	PerformanceMetrics *PerformanceMetrics
}

// Signal is an alias for types.Signal
type Signal = types.Signal

// BacktestStrategy defines the interface for backtesting strategies
type BacktestStrategy interface {
	// Initialize sets up the strategy with any required configuration
	Initialize(ctx context.Context, config interface{}) error

	// OnTick processes a new tick of market data and returns any trading signals
	OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error)

	// OnOrderFilled is called when an order has been filled
	OnOrderFilled(ctx context.Context, order *models.Order) error

	// ClosePositions is called at the end of backtesting to close any open positions
	ClosePositions(ctx context.Context) ([]*Signal, error)
}

// EquityPoint is an alias for types.EquityPoint
type EquityPoint = types.EquityPoint

// DrawdownPoint is an alias for types.DrawdownPoint
type DrawdownPoint = types.DrawdownPoint

// PerformanceMetrics is an alias for types.PerformanceMetrics
type PerformanceMetrics = types.PerformanceMetrics

// DataProvider defines the interface for historical data providers
type DataProvider interface {
	GetHistoricalData(ctx context.Context, symbol string, interval string, startTime time.Time, endTime time.Time) ([]*models.Kline, error)
}

// SlippageModel defines the interface for slippage calculation
type SlippageModel interface {
	CalculateSlippage(price float64, quantity float64, side models.OrderSide) float64
}
