package types

import (
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// BacktestRequest represents a request to run a backtest
type BacktestRequest struct {
	UserID         string                 `json:"user_id"`
	StrategyName   string                 `json:"strategy_name"`
	Symbol         string                 `json:"symbol"`
	Timeframe      string                 `json:"timeframe"`
	StartTime      time.Time              `json:"start_time"`
	EndTime        time.Time              `json:"end_time"`
	InitialCapital float64                `json:"initial_capital"`
	PositionSize   float64                `json:"position_size"`
	MaxPositions   int                    `json:"max_positions"`
	Parameters     map[string]interface{} `json:"parameters"`
}

// BacktestConfig holds configuration for a backtest run
type BacktestConfig struct {
	InitialCapital float64                `json:"initial_capital"`
	PositionSize   float64                `json:"position_size"`
	MaxPositions   int                    `json:"max_positions"`
	Commission     float64                `json:"commission"`
	Slippage       float64                `json:"slippage"`
	DataFeed       string                 `json:"data_feed"`
	Parameters     map[string]interface{} `json:"parameters"`
}

// BacktestTrade represents a trade executed during backtesting
type BacktestTrade struct {
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	EntryPrice    float64   `json:"entry_price"`
	ExitPrice     float64   `json:"exit_price"`
	Quantity      float64   `json:"quantity"`
	EntryTime     time.Time `json:"entry_time"`
	ExitTime      time.Time `json:"exit_time"`
	PnL           float64   `json:"pnl"`
	PnLPercent    float64   `json:"pnl_percent"`
	Commission    float64   `json:"commission"`
	Slippage      float64   `json:"slippage"`
	NetPnL        float64   `json:"net_pnl"`
	NetPnLPercent float64   `json:"net_pnl_percent"`
}

// BacktestMetrics represents performance metrics from a backtest
type BacktestMetrics struct {
	TotalTrades   int     `json:"total_trades"`
	WinningTrades int     `json:"winning_trades"`
	LosingTrades  int     `json:"losing_trades"`
	WinRate       float64 `json:"win_rate"`
	AverageWin    float64 `json:"average_win"`
	AverageLoss   float64 `json:"average_loss"`
	LargestWin    float64 `json:"largest_win"`
	LargestLoss   float64 `json:"largest_loss"`
}
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
}

// BacktestStrategy is an interface that strategies must implement for backtesting
type BacktestStrategy interface {
	Reset()
	SetConfig(config *BacktestConfig)
	GetConfig() *BacktestConfig
	OnHistoricalCandleUpdate(candle *models.Candle) (*Signal, error)
	OnHistoricalTickerUpdate(ticker *models.Ticker) (*Signal, error)
	OnHistoricalTradeUpdate(trade *models.Trade) (*Signal, error)
	OnHistoricalMarketDepthUpdate(depth *models.OrderBook) (*Signal, error)
}

// PerformanceAnalyzer is an interface for analyzing backtest performance
type PerformanceAnalyzer interface {
	AddTrade(trade *BacktestTrade)
	AddEquityPoint(point *EquityPoint)
	CalculateMetrics() *BacktestMetrics
	GetTrades() []*BacktestTrade
	GetEquityCurve() []*EquityPoint
	GetDrawdownCurve() []*DrawdownPoint
}

// EquityPoint represents a point on the equity curve
type EquityPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Equity    float64   `json:"equity"`
}

// DrawdownPoint represents a point on the drawdown curve
type DrawdownPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Drawdown  float64   `json:"drawdown"`
}
