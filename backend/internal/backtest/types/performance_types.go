package types

import "time"

// PerformanceMetrics contains performance metrics from a backtest
type PerformanceMetrics struct {
	TotalTrades        int
	WinningTrades      int
	LosingTrades       int
	BreakEvenTrades    int
	TotalReturn        float64
	AnnualizedReturn   float64
	AverageProfitTrade float64
	AverageLossTrade   float64
	ExpectedPayoff     float64
	ProfitFactor       float64
	SharpeRatio        float64
	SortinoRatio       float64
	CalmarRatio        float64
	OmegaRatio         float64
	InformationRatio   float64
	MaxDrawdown        float64
	MaxDrawdownPercent float64
	LargestProfitTrade float64
	LargestLossTrade   float64
	AverageHoldingTime time.Duration
	WinRate            float64
}

// EquityPoint represents a point on the equity curve
type EquityPoint struct {
	Timestamp time.Time
	Equity    float64
}

// DrawdownPoint represents a point on the drawdown curve
type DrawdownPoint struct {
	Timestamp time.Time
	Drawdown  float64
}
