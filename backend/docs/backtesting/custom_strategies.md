# Creating Custom Strategies

This guide explains how to create custom trading strategies for the backtesting framework.

## Strategy Interface

All strategies must implement the `BacktestStrategy` interface:

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

## BaseStrategy

The framework provides a `BaseStrategy` struct that implements the common functionality of the `BacktestStrategy` interface. You can embed this struct in your custom strategy to avoid implementing all methods from scratch:

```go
// BaseStrategy provides a base implementation of the BacktestStrategy interface
type BaseStrategy struct {
    Name string
    // Other common fields
}

// Initialize initializes the strategy with backtest-specific parameters
func (s *BaseStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
    // Default implementation
    return nil
}

// OnTick is called for each new data point (candle, ticker, etc.)
func (s *BaseStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*Signal, error) {
    // Default implementation (no signals)
    return nil, nil
}

// OnOrderFilled is called when an order is filled during the backtest
func (s *BaseStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
    // Default implementation
    return nil
}

// OnPositionClosed is called when a position is closed during the backtest
func (s *BaseStrategy) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
    // Default implementation
    return nil
}
```

## Step-by-Step Guide

### 1. Create a New Strategy Struct

Start by creating a new struct for your strategy that embeds the `BaseStrategy`:

```go
// MyStrategy implements a custom trading strategy
type MyStrategy struct {
    backtest.BaseStrategy
    // Strategy-specific fields
    Parameter1 int
    Parameter2 float64
    // State variables
    prices     map[string][]float64
    indicators map[string][]float64
    positions  map[string]bool
    logger     *zap.Logger
}

// NewMyStrategy creates a new instance of MyStrategy
func NewMyStrategy(param1 int, param2 float64, logger *zap.Logger) *MyStrategy {
    if logger == nil {
        logger, _ = zap.NewDevelopment()
    }

    return &MyStrategy{
        BaseStrategy: backtest.BaseStrategy{
            Name: "MyStrategy",
        },
        Parameter1: param1,
        Parameter2: param2,
        prices:     make(map[string][]float64),
        indicators: make(map[string][]float64),
        positions:  make(map[string]bool),
        logger:     logger,
    }
}
```

### 2. Implement the Initialize Method

Override the `Initialize` method to handle strategy-specific configuration:

```go
// Initialize initializes the strategy with backtest-specific parameters
func (s *MyStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
    // Call the base implementation first
    err := s.BaseStrategy.Initialize(ctx, config)
    if err != nil {
        return err
    }

    // Override parameters if provided in config
    if config != nil {
        if param1, ok := config["parameter1"].(int); ok {
            s.Parameter1 = param1
        }
        if param2, ok := config["parameter2"].(float64); ok {
            s.Parameter2 = param2
        }
    }

    s.logger.Info("Initialized MyStrategy",
        zap.Int("parameter1", s.Parameter1),
        zap.Float64("parameter2", s.Parameter2),
    )

    return nil
}
```

### 3. Implement the OnTick Method

The `OnTick` method is the heart of your strategy. It receives market data and generates trading signals:

```go
// OnTick is called for each new data point (candle, ticker, etc.)
func (s *MyStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*backtest.Signal, error) {
    // Check if data is a Kline
    kline, ok := data.(*models.Kline)
    if !ok {
        return nil, nil
    }

    // Initialize arrays for this symbol if they don't exist
    if _, ok := s.prices[symbol]; !ok {
        s.prices[symbol] = make([]float64, 0)
        s.indicators[symbol] = make([]float64, 0)
        s.positions[symbol] = false
    }

    // Add price to the price array
    s.prices[symbol] = append(s.prices[symbol], kline.Close)

    // Calculate indicators
    // This is where you implement your strategy's logic
    // For example, calculating moving averages, RSI, MACD, etc.
    
    // Example: Calculate a simple indicator (e.g., price momentum)
    var indicator float64
    if len(s.prices[symbol]) > 1 {
        indicator = s.prices[symbol][len(s.prices[symbol])-1] - s.prices[symbol][len(s.prices[symbol])-2]
    }
    s.indicators[symbol] = append(s.indicators[symbol], indicator)

    // Generate signals based on indicators
    var signals []*backtest.Signal

    // Example: Generate buy signal when indicator is positive and we don't have a position
    if len(s.indicators[symbol]) > 0 && s.indicators[symbol][len(s.indicators[symbol])-1] > 0 && !s.positions[symbol] {
        s.logger.Info("Buy signal generated",
            zap.String("symbol", symbol),
            zap.Time("timestamp", timestamp),
            zap.Float64("price", kline.Close),
            zap.Float64("indicator", s.indicators[symbol][len(s.indicators[symbol])-1]),
        )

        signals = append(signals, &backtest.Signal{
            Symbol:    symbol,
            Side:      "BUY",
            Quantity:  1.0, // Fixed quantity for simplicity
            Price:     kline.Close,
            Timestamp: timestamp,
            Reason:    "Positive momentum",
        })

        s.positions[symbol] = true
    }

    // Example: Generate sell signal when indicator is negative and we have a position
    if len(s.indicators[symbol]) > 0 && s.indicators[symbol][len(s.indicators[symbol])-1] < 0 && s.positions[symbol] {
        s.logger.Info("Sell signal generated",
            zap.String("symbol", symbol),
            zap.Time("timestamp", timestamp),
            zap.Float64("price", kline.Close),
            zap.Float64("indicator", s.indicators[symbol][len(s.indicators[symbol])-1]),
        )

        signals = append(signals, &backtest.Signal{
            Symbol:    symbol,
            Side:      "SELL",
            Quantity:  1.0, // Fixed quantity for simplicity
            Price:     kline.Close,
            Timestamp: timestamp,
            Reason:    "Negative momentum",
        })

        s.positions[symbol] = false
    }

    return signals, nil
}
```

### 4. Implement the OnOrderFilled Method

The `OnOrderFilled` method is called when an order is filled during the backtest:

```go
// OnOrderFilled is called when an order is filled during the backtest
func (s *MyStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
    s.logger.Info("Order filled",
        zap.String("symbol", order.Symbol),
        zap.String("side", string(order.Side)),
        zap.Float64("quantity", order.Quantity),
        zap.Float64("price", order.Price),
        zap.Time("time", order.Time),
    )
    
    // You can update strategy state based on filled orders
    // For example, tracking position sizes, average entry prices, etc.
    
    return nil
}
```

### 5. Implement the OnPositionClosed Method

The `OnPositionClosed` method is called when a position is closed during the backtest:

```go
// OnPositionClosed is called when a position is closed during the backtest
func (s *MyStrategy) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
    s.logger.Info("Position closed",
        zap.String("symbol", position.Symbol),
        zap.Float64("entry_price", position.EntryPrice),
        zap.Float64("exit_price", position.ExitPrice),
        zap.Float64("amount", position.Amount),
        zap.Float64("profit", position.ProfitLoss),
        zap.Time("open_time", position.OpenTime),
        zap.Time("close_time", position.CloseTime),
    )
    
    // You can update strategy state based on closed positions
    // For example, tracking performance metrics, adjusting parameters, etc.
    
    return nil
}
```

### 6. Register Your Strategy

To use your strategy in the backtesting framework, you need to register it in the CLI command:

```go
// In internal/cli/backtest_cmd.go

// Create strategy
var strategy backtest.BacktestStrategy
switch strategyName {
case "simple_ma":
    strategy = strategies.NewSimpleMAStrategy(shortPeriod, longPeriod, logger)
case "my_strategy":
    strategy = strategies.NewMyStrategy(param1, param2, logger)
default:
    return fmt.Errorf("unknown strategy: %s", strategyName)
}
```

## Example: RSI Strategy

Here's an example of a strategy that uses the Relative Strength Index (RSI) indicator:

```go
// RSIStrategy implements a strategy based on the Relative Strength Index
type RSIStrategy struct {
    backtest.BaseStrategy
    Period       int
    OverboughtLevel float64
    OversoldLevel  float64
    prices       map[string][]float64
    rsi          map[string][]float64
    position     map[string]bool
    logger       *zap.Logger
}

// NewRSIStrategy creates a new RSIStrategy
func NewRSIStrategy(period int, overboughtLevel, oversoldLevel float64, logger *zap.Logger) *RSIStrategy {
    if logger == nil {
        logger, _ = zap.NewDevelopment()
    }

    return &RSIStrategy{
        BaseStrategy: backtest.BaseStrategy{
            Name: "RSIStrategy",
        },
        Period:         period,
        OverboughtLevel: overboughtLevel,
        OversoldLevel:   oversoldLevel,
        prices:         make(map[string][]float64),
        rsi:            make(map[string][]float64),
        position:       make(map[string]bool),
        logger:         logger,
    }
}

// OnTick is called for each new data point
func (s *RSIStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*backtest.Signal, error) {
    // Check if data is a Kline
    kline, ok := data.(*models.Kline)
    if !ok {
        return nil, nil
    }

    // Initialize arrays for this symbol if they don't exist
    if _, ok := s.prices[symbol]; !ok {
        s.prices[symbol] = make([]float64, 0)
        s.rsi[symbol] = make([]float64, 0)
        s.position[symbol] = false
    }

    // Add price to the price array
    s.prices[symbol] = append(s.prices[symbol], kline.Close)

    // Calculate RSI
    rsiValue := calculateRSI(s.prices[symbol], s.Period)
    s.rsi[symbol] = append(s.rsi[symbol], rsiValue)

    // Generate signals
    var signals []*backtest.Signal

    // We need enough data to calculate RSI
    if len(s.prices[symbol]) >= s.Period {
        // Buy signal: RSI crosses below oversold level
        if len(s.rsi[symbol]) >= 2 && 
           s.rsi[symbol][len(s.rsi[symbol])-2] <= s.OversoldLevel && 
           s.rsi[symbol][len(s.rsi[symbol])-1] > s.OversoldLevel && 
           !s.position[symbol] {
            
            signals = append(signals, &backtest.Signal{
                Symbol:    symbol,
                Side:      "BUY",
                Quantity:  1.0,
                Price:     kline.Close,
                Timestamp: timestamp,
                Reason:    fmt.Sprintf("RSI crossed above oversold level (%.2f)", s.OversoldLevel),
            })

            s.position[symbol] = true
        }

        // Sell signal: RSI crosses above overbought level
        if len(s.rsi[symbol]) >= 2 && 
           s.rsi[symbol][len(s.rsi[symbol])-2] >= s.OverboughtLevel && 
           s.rsi[symbol][len(s.rsi[symbol])-1] < s.OverboughtLevel && 
           s.position[symbol] {
            
            signals = append(signals, &backtest.Signal{
                Symbol:    symbol,
                Side:      "SELL",
                Quantity:  1.0,
                Price:     kline.Close,
                Timestamp: timestamp,
                Reason:    fmt.Sprintf("RSI crossed below overbought level (%.2f)", s.OverboughtLevel),
            })

            s.position[symbol] = false
        }
    }

    return signals, nil
}

// calculateRSI calculates the Relative Strength Index
func calculateRSI(prices []float64, period int) float64 {
    if len(prices) < period+1 {
        return 50.0 // Default value when not enough data
    }

    var gains, losses float64
    for i := len(prices) - period; i < len(prices); i++ {
        change := prices[i] - prices[i-1]
        if change >= 0 {
            gains += change
        } else {
            losses -= change
        }
    }

    if losses == 0 {
        return 100.0
    }

    rs := gains / losses
    return 100.0 - (100.0 / (1.0 + rs))
}
```

## Best Practices for Strategy Development

1. **Keep it simple**: Start with simple strategies and gradually add complexity.
2. **Avoid overfitting**: Test your strategy on different time periods and symbols.
3. **Use proper risk management**: Implement position sizing and stop-loss mechanisms.
4. **Handle edge cases**: Ensure your strategy can handle missing data, extreme market conditions, etc.
5. **Log important events**: Use logging to track strategy decisions and performance.
6. **Separate logic from implementation**: Keep your trading logic separate from the backtesting framework.
7. **Test thoroughly**: Write unit tests for your strategy's core logic.
8. **Document your strategy**: Include comments explaining the strategy's logic and parameters.

## Common Indicators

Here are some common indicators you might want to implement in your strategies:

- Moving Averages (Simple, Exponential, Weighted)
- Relative Strength Index (RSI)
- Moving Average Convergence Divergence (MACD)
- Bollinger Bands
- Average True Range (ATR)
- Stochastic Oscillator
- Ichimoku Cloud
- Fibonacci Retracement

You can implement these indicators yourself or use a technical analysis library like [ta-lib](https://github.com/markcheno/go-talib) or [techan](https://github.com/sdcoffey/techan).
