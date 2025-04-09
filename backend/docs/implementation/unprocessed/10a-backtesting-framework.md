# Backtesting Framework Specification

## Overview

The backtesting framework allows traders and developers to test trading strategies against historical market data to evaluate their performance before deploying them in a live trading environment. This framework simulates the execution of trading strategies, tracks positions, calculates profits and losses, and generates performance metrics.

## Architecture

The backtesting framework consists of the following components:

1. **Data Provider**: Retrieves and prepares historical market data for backtesting
2. **Backtesting Engine**: Simulates the execution of trading strategies against historical data
3. **Position Tracker**: Tracks open positions, calculates P&L, and manages position lifecycle
4. **Performance Analyzer**: Calculates performance metrics and generates reports
5. **Strategy Adapter**: Adapts trading strategies to work with the backtesting engine

## Components

### Data Provider

The Data Provider is responsible for retrieving and preparing historical market data for backtesting.

```go
// DataProvider defines the interface for retrieving historical market data
type DataProvider interface {
    // GetKlines retrieves historical candlestick data for a symbol within a time range
    GetKlines(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error)
    
    // GetTickers retrieves historical ticker data for a symbol within a time range
    GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error)
    
    // GetOrderBook retrieves historical order book snapshots for a symbol at a specific time
    GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error)
}
```

#### Implementation Options

1. **SQLiteDataProvider**: Retrieves historical data from a local SQLite database
2. **CSVDataProvider**: Loads historical data from CSV files
3. **APIDataProvider**: Fetches historical data from exchange APIs (where available)

### Backtesting Engine

The Backtesting Engine simulates the execution of trading strategies against historical data.

```go
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
    Strategy           Strategy
}

// BacktestEngine defines the interface for the backtesting engine
type BacktestEngine interface {
    // Run executes a backtest with the given configuration
    Run(ctx context.Context, config *BacktestConfig) (*BacktestResult, error)
    
    // GetEvents returns all events generated during the backtest
    GetEvents() []*BacktestEvent
    
    // GetPositions returns all positions created during the backtest
    GetPositions() []*models.Position
    
    // GetTrades returns all trades executed during the backtest
    GetTrades() []*models.Order
}
```

### Position Tracker

The Position Tracker manages positions during the backtest, including opening, updating, and closing positions.

```go
// PositionTracker defines the interface for tracking positions during a backtest
type PositionTracker interface {
    // OpenPosition opens a new position
    OpenPosition(symbol string, side string, entryPrice float64, quantity float64, timestamp time.Time) (*models.Position, error)
    
    // ClosePosition closes an existing position
    ClosePosition(positionID string, exitPrice float64, timestamp time.Time) (*models.ClosedPosition, error)
    
    // UpdatePosition updates an existing position (e.g., for partial closes)
    UpdatePosition(positionID string, newQuantity float64, timestamp time.Time) (*models.Position, error)
    
    // GetOpenPositions returns all currently open positions
    GetOpenPositions() []*models.Position
    
    // GetClosedPositions returns all closed positions
    GetClosedPositions() []*models.ClosedPosition
    
    // CalculateUnrealizedPnL calculates the unrealized P&L for all open positions
    CalculateUnrealizedPnL(currentPrices map[string]float64) (float64, error)
}
```

### Performance Analyzer

The Performance Analyzer calculates performance metrics and generates reports based on backtest results.

```go
// PerformanceMetrics contains performance metrics for a backtest
type PerformanceMetrics struct {
    TotalReturn        float64
    AnnualizedReturn   float64
    SharpeRatio        float64
    SortinoRatio       float64
    MaxDrawdown        float64
    MaxDrawdownPercent float64
    WinRate            float64
    ProfitFactor       float64
    ExpectedPayoff     float64
    TotalTrades        int
    WinningTrades      int
    LosingTrades       int
    BreakEvenTrades    int
    AverageProfitTrade float64
    AverageLossTrade   float64
    LargestProfitTrade float64
    LargestLossTrade   float64
    AverageHoldingTime time.Duration
}

// PerformanceAnalyzer defines the interface for analyzing backtest performance
type PerformanceAnalyzer interface {
    // CalculateMetrics calculates performance metrics from backtest results
    CalculateMetrics(result *BacktestResult) (*PerformanceMetrics, error)
    
    // GenerateReport generates a detailed performance report
    GenerateReport(result *BacktestResult, metrics *PerformanceMetrics) (*BacktestReport, error)
    
    // GenerateEquityCurve generates an equity curve from backtest results
    GenerateEquityCurve(result *BacktestResult) ([]*EquityPoint, error)
    
    // GenerateDrawdownCurve generates a drawdown curve from backtest results
    GenerateDrawdownCurve(result *BacktestResult) ([]*DrawdownPoint, error)
}
```

### Strategy Adapter

The Strategy Adapter adapts trading strategies to work with the backtesting engine.

```go
// BacktestStrategy defines the interface for strategies used in backtesting
type BacktestStrategy interface {
    // Initialize initializes the strategy with backtest-specific parameters
    Initialize(ctx context.Context, config map[string]interface{}) error
    
    // OnTick is called for each new data point (candle, ticker, etc.)
    OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error)
    
    // OnOrderFilled is called when an order is filled during the backtest
    OnOrderFilled(ctx context.Context, order *models.Order) error
    
    // OnPositionClosed is called when a position is closed during the backtest
    OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error
}

// Signal represents a trading signal generated by a strategy
type Signal struct {
    Symbol    string
    Side      string
    Quantity  float64
    Price     float64
    Timestamp time.Time
    Reason    string
}
```

## Workflow

1. **Setup**: Configure the backtest with start/end dates, initial capital, symbols, and strategy parameters
2. **Data Loading**: Load historical data for the specified symbols and time range
3. **Simulation**: Process data chronologically, feeding it to the strategy and executing signals
4. **Position Management**: Track positions, calculate P&L, and manage position lifecycle
5. **Analysis**: Calculate performance metrics and generate reports
6. **Visualization**: Generate equity curves, drawdown charts, and other visualizations

## Implementation Details

### Slippage Models

The framework supports different slippage models to simulate real-world trading conditions:

1. **NoSlippage**: No slippage is applied (ideal conditions)
2. **FixedSlippage**: A fixed amount or percentage is added to each trade
3. **VariableSlippage**: Slippage varies based on volatility and volume
4. **OrderBookSlippage**: Simulates slippage based on order book depth

```go
// SlippageModel defines the interface for simulating slippage
type SlippageModel interface {
    // CalculateSlippage calculates the slippage for a trade
    CalculateSlippage(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64
}
```

### Commission Models

The framework supports different commission models to simulate trading costs:

1. **NoCommission**: No commission is applied (ideal conditions)
2. **FixedCommission**: A fixed amount or percentage is charged for each trade
3. **TieredCommission**: Commission rates vary based on trading volume

```go
// CommissionModel defines the interface for calculating trading commissions
type CommissionModel interface {
    // CalculateCommission calculates the commission for a trade
    CalculateCommission(symbol string, side string, quantity float64, price float64) float64
}
```

### Backtest Events

The framework generates events during the backtest to track the simulation process:

```go
// BacktestEventType defines the type of backtest event
type BacktestEventType string

const (
    EventStrategyInitialized BacktestEventType = "strategy_initialized"
    EventDataLoaded          BacktestEventType = "data_loaded"
    EventSignalGenerated     BacktestEventType = "signal_generated"
    EventOrderCreated        BacktestEventType = "order_created"
    EventOrderFilled         BacktestEventType = "order_filled"
    EventPositionOpened      BacktestEventType = "position_opened"
    EventPositionClosed      BacktestEventType = "position_closed"
    EventError               BacktestEventType = "error"
)

// BacktestEvent represents an event that occurred during the backtest
type BacktestEvent struct {
    Type      BacktestEventType
    Timestamp time.Time
    Symbol    string
    Data      interface{}
}
```

### Backtest Results

The framework generates comprehensive results from the backtest:

```go
// BacktestResult contains the results of a backtest
type BacktestResult struct {
    Config           *BacktestConfig
    StartTime        time.Time
    EndTime          time.Time
    InitialCapital   float64
    FinalCapital     float64
    Trades           []*models.Order
    Positions        []*models.Position
    ClosedPositions  []*models.ClosedPosition
    EquityCurve      []*EquityPoint
    DrawdownCurve    []*DrawdownPoint
    Events           []*BacktestEvent
    PerformanceMetrics *PerformanceMetrics
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
```

## CLI Integration

The backtesting framework will be integrated with the CLI to allow users to run backtests from the command line:

```
crypto-bot backtest --strategy=NewCoinStrategy --start=2023-01-01 --end=2023-12-31 --capital=10000 --symbols=BTCUSDT,ETHUSDT --interval=1h
```

## API Integration

The backtesting framework will also be exposed through the API to allow users to run backtests programmatically:

```
POST /api/v1/backtest
{
    "strategy": "NewCoinStrategy",
    "start_time": "2023-01-01T00:00:00Z",
    "end_time": "2023-12-31T23:59:59Z",
    "initial_capital": 10000,
    "symbols": ["BTCUSDT", "ETHUSDT"],
    "interval": "1h",
    "parameters": {
        "volume_threshold": 1000000,
        "price_change_threshold": 5.0
    }
}
```

## Visualization

The backtesting framework will generate visualizations to help users analyze the results:

1. **Equity Curve**: Shows the evolution of portfolio value over time
2. **Drawdown Curve**: Shows the drawdowns experienced during the backtest
3. **Trade Distribution**: Shows the distribution of winning and losing trades
4. **Monthly Returns**: Shows the returns for each month during the backtest
5. **Position Sizing**: Shows how position sizes evolved during the backtest

## Future Enhancements

1. **Monte Carlo Simulation**: Run multiple backtests with randomized parameters to assess strategy robustness
2. **Walk-Forward Analysis**: Test strategies on different time periods to assess consistency
3. **Parameter Optimization**: Automatically find optimal strategy parameters
4. **Multi-Strategy Backtesting**: Test multiple strategies simultaneously
5. **Portfolio Backtesting**: Test strategies across a portfolio of assets
6. **Machine Learning Integration**: Use ML to enhance strategy performance

## Conclusion

The backtesting framework provides a comprehensive solution for testing trading strategies against historical data. It simulates real-world trading conditions, tracks positions, calculates performance metrics, and generates detailed reports. This allows traders and developers to evaluate and refine their strategies before deploying them in a live trading environment.
