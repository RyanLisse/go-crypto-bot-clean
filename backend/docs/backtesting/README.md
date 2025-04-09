# Backtesting Framework Documentation

The backtesting framework allows traders and developers to test trading strategies against historical market data to evaluate their performance before deploying them in a live trading environment.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Components](#components)
4. [Usage](#usage)
5. [Examples](#examples)
6. [Performance Metrics](#performance-metrics)
7. [Best Practices](#best-practices)
8. [Extending the Framework](#extending-the-framework)

## Overview

The backtesting framework simulates the execution of trading strategies against historical market data. It tracks positions, calculates profits and losses, and generates performance metrics to help evaluate strategy effectiveness.

Key features:
- Historical data loading from various sources (SQLite, CSV, in-memory)
- Position tracking and P&L calculation
- Performance metrics calculation (Sharpe ratio, drawdown, etc.)
- Slippage models for realistic trade simulation
- Strategy interface for testing different strategies
- CLI command for running backtests

## Architecture

The backtesting framework consists of the following components:

1. **Data Provider**: Retrieves and prepares historical market data for backtesting
2. **Backtesting Engine**: Simulates the execution of trading strategies against historical data
3. **Position Tracker**: Tracks open positions, calculates P&L, and manages position lifecycle
4. **Performance Analyzer**: Calculates performance metrics and generates reports
5. **Strategy Adapter**: Adapts trading strategies to work with the backtesting engine

![Backtesting Architecture](architecture.png)

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

#### Available Implementations

1. **SQLiteDataProvider**: Retrieves historical data from a SQLite database
2. **CSVDataProvider**: Loads historical data from CSV files
3. **InMemoryDataProvider**: Stores data in memory (primarily for testing)

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
    Strategy           BacktestStrategy
}

// Engine implements the backtesting engine
type Engine struct {
    // ...
}

// Run executes a backtest with the given configuration
func (e *Engine) Run(ctx context.Context) (*BacktestResult, error) {
    // ...
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

### Strategy Interface

The Strategy Interface defines how trading strategies interact with the backtesting engine.

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
```

## Usage

### Running a Backtest via CLI

The backtesting framework can be run via the CLI:

```bash
go run cmd/backtest/main.go backtest --strategy=simple_ma --symbols=BTCUSDT --start=2023-01-01 --end=2023-12-31 --capital=10000 --interval=1h
```

Available options:
- `--strategy`: Strategy to backtest (e.g., simple_ma)
- `--symbols`: Symbols to backtest (comma-separated)
- `--start`: Start date (YYYY-MM-DD)
- `--end`: End date (YYYY-MM-DD)
- `--capital`: Initial capital
- `--interval`: Candle interval (1m, 5m, 15m, 1h, 4h, 1d)
- `--data-dir`: Directory containing historical data CSV files
- `--short-period`: Short period for MA strategy
- `--long-period`: Long period for MA strategy
- `--output`: Output file for backtest results

### Running a Backtest Programmatically

You can also run a backtest programmatically:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/ryanlisse/go-crypto-bot/internal/backtest"
    "github.com/ryanlisse/go-crypto-bot/internal/backtest/strategies"
    "go.uber.org/zap"
)

func main() {
    // Create logger
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()

    // Create data provider
    dataProvider := backtest.NewInMemoryDataProvider()
    // Load data into the provider...

    // Create strategy
    strategy := strategies.NewSimpleMAStrategy(10, 50, logger)

    // Create slippage model
    slippageModel := backtest.NewFixedSlippage(0.1) // 0.1% slippage

    // Create backtest config
    config := &backtest.BacktestConfig{
        StartTime:          time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
        EndTime:            time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
        InitialCapital:     10000,
        Symbols:            []string{"BTCUSDT"},
        Interval:           "1h",
        CommissionRate:     0.001, // 0.1% commission
        SlippageModel:      slippageModel,
        EnableShortSelling: false,
        DataProvider:       dataProvider,
        Strategy:           strategy,
        Logger:             logger,
    }

    // Create backtest engine
    engine := backtest.NewEngine(config)

    // Run backtest
    result, err := engine.Run(context.Background())
    if err != nil {
        logger.Fatal("Backtest failed", zap.Error(err))
    }

    // Print results
    fmt.Printf("Initial Capital: $%.2f\n", result.InitialCapital)
    fmt.Printf("Final Capital: $%.2f\n", result.FinalCapital)
    fmt.Printf("Total Return: %.2f%%\n", result.PerformanceMetrics.TotalReturn)
    fmt.Printf("Sharpe Ratio: %.2f\n", result.PerformanceMetrics.SharpeRatio)
    fmt.Printf("Max Drawdown: %.2f%%\n", result.PerformanceMetrics.MaxDrawdownPercent)
    fmt.Printf("Win Rate: %.2f%%\n", result.PerformanceMetrics.WinRate)
    fmt.Printf("Total Trades: %d\n", result.PerformanceMetrics.TotalTrades)
}
```

## Examples

### Simple Moving Average Strategy

The Simple Moving Average (SMA) strategy is a basic trend-following strategy that generates buy and sell signals based on the crossover of two moving averages.

```go
// SimpleMAStrategy implements a simple moving average crossover strategy
type SimpleMAStrategy struct {
    backtest.BaseStrategy
    ShortPeriod int
    LongPeriod  int
    shortMA     map[string][]float64
    longMA      map[string][]float64
    prices      map[string][]float64
    position    map[string]bool // true if long position is open
    logger      *zap.Logger
}

// OnTick is called for each new data point (candle, ticker, etc.)
func (s *SimpleMAStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*backtest.Signal, error) {
    // Check if data is a Kline
    kline, ok := data.(*models.Kline)
    if !ok {
        return nil, nil
    }

    // Add price to the price array
    s.prices[symbol] = append(s.prices[symbol], kline.Close)

    // Calculate moving averages
    shortMA := calculateSMA(s.prices[symbol], s.ShortPeriod)
    longMA := calculateSMA(s.prices[symbol], s.LongPeriod)

    // Add moving averages to the arrays
    s.shortMA[symbol] = append(s.shortMA[symbol], shortMA)
    s.longMA[symbol] = append(s.longMA[symbol], longMA)

    // Generate signals
    var signals []*backtest.Signal

    // We need enough data to calculate both moving averages
    if len(s.prices[symbol]) >= s.LongPeriod {
        // Check for crossover
        if len(s.shortMA[symbol]) >= 2 && len(s.longMA[symbol]) >= 2 {
            prevShortMA := s.shortMA[symbol][len(s.shortMA[symbol])-2]
            prevLongMA := s.longMA[symbol][len(s.longMA[symbol])-2]
            currentShortMA := s.shortMA[symbol][len(s.shortMA[symbol])-1]
            currentLongMA := s.longMA[symbol][len(s.longMA[symbol])-1]

            // Buy signal: short MA crosses above long MA
            if prevShortMA <= prevLongMA && currentShortMA > currentLongMA && !s.position[symbol] {
                signals = append(signals, &backtest.Signal{
                    Symbol:    symbol,
                    Side:      "BUY",
                    Quantity:  1.0, // Fixed quantity for simplicity
                    Price:     kline.Close,
                    Timestamp: timestamp,
                    Reason:    "MA crossover (short > long)",
                })

                s.position[symbol] = true
            }

            // Sell signal: short MA crosses below long MA
            if prevShortMA >= prevLongMA && currentShortMA < currentLongMA && s.position[symbol] {
                signals = append(signals, &backtest.Signal{
                    Symbol:    symbol,
                    Side:      "SELL",
                    Quantity:  1.0, // Fixed quantity for simplicity
                    Price:     kline.Close,
                    Timestamp: timestamp,
                    Reason:    "MA crossover (short < long)",
                })

                s.position[symbol] = false
            }
        }
    }

    return signals, nil
}
```

## Performance Metrics

The backtesting framework calculates the following performance metrics:

| Metric | Description |
|--------|-------------|
| Total Return | The total percentage return of the strategy |
| Annualized Return | The annualized percentage return of the strategy |
| Sharpe Ratio | Risk-adjusted return (higher is better) |
| Sortino Ratio | Downside risk-adjusted return (higher is better) |
| Max Drawdown | The maximum peak-to-trough decline in portfolio value |
| Win Rate | Percentage of trades that were profitable |
| Profit Factor | Gross profit divided by gross loss |
| Expected Payoff | Average profit/loss per trade |
| Total Trades | Total number of trades executed |
| Winning Trades | Number of profitable trades |
| Losing Trades | Number of unprofitable trades |
| Average Profit Trade | Average profit of winning trades |
| Average Loss Trade | Average loss of losing trades |
| Largest Profit Trade | Largest profit from a single trade |
| Largest Loss Trade | Largest loss from a single trade |
| Average Holding Time | Average duration of trades |

## Best Practices

### Data Quality

- Use high-quality historical data from reliable sources
- Ensure data is clean and free of errors
- Use sufficient historical data to cover different market conditions
- Be aware of survivorship bias in your data

### Strategy Development

- Start with simple strategies and gradually add complexity
- Test strategies on multiple symbols and timeframes
- Use proper position sizing and risk management
- Avoid overfitting by testing on out-of-sample data

### Realistic Simulation

- Include transaction costs (commissions, slippage)
- Consider market impact for larger positions
- Account for liquidity constraints
- Be conservative in your assumptions

### Performance Evaluation

- Don't rely solely on total return
- Consider risk-adjusted metrics (Sharpe, Sortino)
- Analyze drawdowns and recovery periods
- Evaluate consistency across different market conditions

## Extending the Framework

### Adding a New Strategy

To add a new strategy, implement the `BacktestStrategy` interface:

```go
type MyStrategy struct {
    backtest.BaseStrategy
    // Strategy-specific fields
}

func NewMyStrategy() *MyStrategy {
    return &MyStrategy{
        BaseStrategy: backtest.BaseStrategy{
            Name: "MyStrategy",
        },
        // Initialize strategy-specific fields
    }
}

func (s *MyStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
    // Initialize strategy with configuration
    return nil
}

func (s *MyStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*backtest.Signal, error) {
    // Generate signals based on your strategy logic
    return signals, nil
}

func (s *MyStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
    // Handle order fills
    return nil
}

func (s *MyStrategy) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
    // Handle position closes
    return nil
}
```

### Adding a New Data Provider

To add a new data provider, implement the `DataProvider` interface:

```go
type MyDataProvider struct {
    // Provider-specific fields
}

func NewMyDataProvider() *MyDataProvider {
    return &MyDataProvider{
        // Initialize provider-specific fields
    }
}

func (p *MyDataProvider) GetKlines(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error) {
    // Retrieve klines from your data source
    return klines, nil
}

func (p *MyDataProvider) GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error) {
    // Retrieve tickers from your data source
    return tickers, nil
}

func (p *MyDataProvider) GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error) {
    // Retrieve order book from your data source
    return orderBook, nil
}
```

### Adding a New Slippage Model

To add a new slippage model, implement the `SlippageModel` interface:

```go
type MySlippageModel struct {
    // Model-specific fields
}

func NewMySlippageModel() *MySlippageModel {
    return &MySlippageModel{
        // Initialize model-specific fields
    }
}

func (s *MySlippageModel) CalculateSlippage(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
    // Calculate slippage based on your model
    return slippage
}
```
