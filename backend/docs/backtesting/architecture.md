# Backtesting Framework Architecture

```mermaid
graph TD
    A[User/CLI] --> B[Backtest Engine]
    B --> C[Data Provider]
    B --> D[Strategy]
    B --> E[Position Tracker]
    B --> F[Performance Analyzer]
    
    C -->|Historical Data| B
    D -->|Trading Signals| B
    B -->|Position Updates| E
    B -->|Backtest Results| F
    F -->|Performance Metrics| A
    
    C1[SQLite Provider] --> C
    C2[CSV Provider] --> C
    C3[In-Memory Provider] --> C
    
    D1[Simple MA Strategy] --> D
    D2[Custom Strategies] --> D
    
    E -->|Position Management| B
    
    F1[Equity Curve] --> F
    F2[Drawdown Analysis] --> F
    F3[Trade Statistics] --> F
    
    G[Slippage Model] --> B
    G1[No Slippage] --> G
    G2[Fixed Slippage] --> G
    G3[Variable Slippage] --> G
    G4[OrderBook Slippage] --> G
```

## Component Descriptions

### Backtest Engine
The core component that orchestrates the backtesting process. It loads historical data, processes it chronologically, feeds it to the strategy, executes signals, tracks positions, and calculates performance metrics.

### Data Provider
Responsible for retrieving and preparing historical market data for backtesting. Multiple implementations are available for different data sources.

### Strategy
Implements the trading logic that generates buy and sell signals based on market data. Strategies implement a common interface to work with the backtesting engine.

### Position Tracker
Manages positions during the backtest, including opening, updating, and closing positions. It tracks open positions, calculates P&L, and maintains a history of closed positions.

### Performance Analyzer
Calculates performance metrics and generates reports based on backtest results. It provides insights into strategy performance, including risk-adjusted returns, drawdowns, and trade statistics.

### Slippage Model
Simulates the price slippage that occurs in real trading. Different models are available to simulate various market conditions and order execution scenarios.
