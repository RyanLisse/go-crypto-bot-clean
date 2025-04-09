# Domain Models for Crypto Trading Bot

This document outlines the core domain models for the Go cryptocurrency trading bot, mapping the data structures from the existing Python implementation to idiomatic Go.

## 1. Key Domain Entities

### BoughtCoin

```go
package models

import (
    "time"
)

// BoughtCoin represents a cryptocurrency that has been purchased and is being tracked
type BoughtCoin struct {
    ID            uint      `json:"id"`
    Symbol        string    `json:"symbol"`          // Trading pair (e.g., "BTCUSDT")
    PurchasePrice float64   `json:"purchase_price"`  // Price at purchase time
    Quantity      float64   `json:"quantity"`        // Amount purchased
    PurchaseTime  time.Time `json:"purchase_time"`   // When the purchase was made
    IsDeleted     bool      `json:"is_deleted"`      // Soft delete flag
    
    // Additional fields
    StopLossPrice float64   `json:"stop_loss_price"` // Price at which to trigger stop loss
    TakeProfitLevels []TakeProfitLevel `json:"take_profit_levels"` // Multiple TP levels
}

// TakeProfitLevel represents a target level for taking profit on a trade
type TakeProfitLevel struct {
    Percentage   float64  `json:"percentage"`   // Percentage above purchase price
    SellQuantity float64  `json:"sell_quantity"` // Amount to sell at this level
    IsReached    bool     `json:"is_reached"`    // Whether this level has been hit
}

// CalculateCurrentProfit computes the current profit percentage
func (c *BoughtCoin) CalculateCurrentProfit(currentPrice float64) float64 {
    return (currentPrice - c.PurchasePrice) / c.PurchasePrice * 100
}

// ShouldTriggerStopLoss checks if the stop loss should be triggered
func (c *BoughtCoin) ShouldTriggerStopLoss(currentPrice float64) bool {
    return currentPrice <= c.StopLossPrice
}

// GetNextUnreachedProfitLevel returns the next take profit level that hasn't been reached
func (c *BoughtCoin) GetNextUnreachedProfitLevel() *TakeProfitLevel {
    for i := range c.TakeProfitLevels {
        if !c.TakeProfitLevels[i].IsReached {
            return &c.TakeProfitLevels[i]
        }
    }
    return nil
}
```

### NewCoin

```go
package models

import (
    "time"
)

// NewCoin represents a newly detected coin on the exchange
type NewCoin struct {
    ID          uint      `json:"id"`
    Symbol      string    `json:"symbol"`          // Trading pair
    DetectedAt  time.Time `json:"detected_at"`     // When it was first seen
    LastChecked time.Time `json:"last_checked"`    // Last time we checked status
    IsActive    bool      `json:"is_active"`       // Whether it's still considered "new"
}
```

### PurchaseDecision

```go
package models

import (
    "time"
)

// PurchaseDecisionStatus represents the outcome of a purchase decision
type PurchaseDecisionStatus string

const (
    StatusPending   PurchaseDecisionStatus = "pending"
    StatusPurchased PurchaseDecisionStatus = "purchased"
    StatusRejected  PurchaseDecisionStatus = "rejected"
)

// PurchaseDecision tracks the logic around purchasing decisions
type PurchaseDecision struct {
    ID          uint                 `json:"id"`
    Symbol      string               `json:"symbol"`          // Trading pair
    Timestamp   time.Time            `json:"timestamp"`       // When the decision was made
    Status      PurchaseDecisionStatus `json:"status"`
    Reason      string               `json:"reason"`          // Why purchased/rejected
    Price       float64              `json:"price,omitempty"` // Price when decision was made
}
```

### LogEvent

```go
package models

import (
    "time"
)

// LogLevel represents the severity of a log event
type LogLevel string

const (
    LogLevelDebug   LogLevel = "debug"
    LogLevelInfo    LogLevel = "info"
    LogLevelWarning LogLevel = "warning"
    LogLevelError   LogLevel = "error"
)

// LogEvent represents a log entry in the database
type LogEvent struct {
    ID        uint      `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Level     LogLevel  `json:"level"`
    Message   string    `json:"message"`
    Context   string    `json:"context"` // JSON string with additional context
}
```

### Account / Wallet

```go
package models

// Wallet represents the user's cryptocurrency wallet
type Wallet struct {
    USDT            float64            `json:"usdt"`              // Available USDT balance
    Assets          map[string]float64 `json:"assets"`            // Map of coin symbol to amount
    ReservedBalance map[string]float64 `json:"reserved_balance"`  // Funds in open orders
}
```

### Order

```go
package models

import (
    "time"
)

// OrderType represents the type of order
type OrderType string

const (
    OrderTypeMarket OrderType = "market"
    OrderTypeLimit  OrderType = "limit"
)

// OrderSide represents buy or sell side
type OrderSide string

const (
    OrderSideBuy  OrderSide = "buy"
    OrderSideSell OrderSide = "sell"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
    OrderStatusNew       OrderStatus = "new"
    OrderStatusPartial   OrderStatus = "partially_filled"
    OrderStatusFilled    OrderStatus = "filled"
    OrderStatusCanceled  OrderStatus = "canceled"
    OrderStatusRejected  OrderStatus = "rejected"
)

// Order represents a buy or sell order on the exchange
type Order struct {
    ID            string      `json:"id"`            // Exchange order ID
    Symbol        string      `json:"symbol"`        // Trading pair
    Side          OrderSide   `json:"side"`          // Buy or sell
    Type          OrderType   `json:"type"`          // Market or limit
    Quantity      float64     `json:"quantity"`      // Amount to buy/sell
    Price         float64     `json:"price"`         // Price for limit orders
    Status        OrderStatus `json:"status"`        // Current status
    CreatedAt     time.Time   `json:"created_at"`    // When order was created
    UpdatedAt     time.Time   `json:"updated_at"`    // Last status update
    FilledQty     float64     `json:"filled_qty"`    // Amount filled so far
    AvgPrice      float64     `json:"avg_price"`     // Average execution price
    Fee           float64     `json:"fee"`           // Exchange fee
    FeeCurrency   string      `json:"fee_currency"`  // Currency of the fee
    
    // Optional fields for advanced orders
    StopPrice     *float64    `json:"stop_price,omitempty"`     // For stop orders
    IcebergQty    *float64    `json:"iceberg_qty,omitempty"`    // For iceberg orders
}
```

## 2. Domain Service Interfaces

These interfaces define the core business operations and follow the Ports and Adapters pattern.

### Repository Interfaces

```go
package service

import (
    "context"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// BoughtCoinRepository defines operations for bought coins storage
type BoughtCoinRepository interface {
    // CRUD operations
    Store(ctx context.Context, coin *models.BoughtCoin) error
    FindByID(ctx context.Context, id uint) (*models.BoughtCoin, error)
    FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error)
    FindActive(ctx context.Context) ([]*models.BoughtCoin, error)
    SoftDelete(ctx context.Context, id uint) error
    Restore(ctx context.Context, id uint) error
    Update(ctx context.Context, coin *models.BoughtCoin) error
}

// NewCoinRepository defines operations for new coins storage
type NewCoinRepository interface {
    Store(ctx context.Context, coin *models.NewCoin) error
    FindAll(ctx context.Context) ([]*models.NewCoin, error)
    FindActive(ctx context.Context) ([]*models.NewCoin, error)
    Update(ctx context.Context, coin *models.NewCoin) error
    Archive(ctx context.Context, id uint) error
}

// PurchaseDecisionRepository defines operations for purchase decisions
type PurchaseDecisionRepository interface {
    Store(ctx context.Context, decision *models.PurchaseDecision) error
    FindByID(ctx context.Context, id uint) (*models.PurchaseDecision, error)
    FindBySymbol(ctx context.Context, symbol string) ([]*models.PurchaseDecision, error)
    Update(ctx context.Context, decision *models.PurchaseDecision) error
}

// LogRepository defines operations for logging events
type LogRepository interface {
    Store(ctx context.Context, logEvent *models.LogEvent) error
    FindByLevel(ctx context.Context, level models.LogLevel) ([]*models.LogEvent, error)
    FindByTimeRange(ctx context.Context, start, end time.Time) ([]*models.LogEvent, error)
}
```

### Exchange Service Interface

```go
package service

import (
    "context"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// ExchangeService defines operations to interact with the crypto exchange
type ExchangeService interface {
    // Market data
    GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)
    GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error)
    GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*models.Kline, error)
    
    // Account operations
    GetWallet(ctx context.Context) (*models.Wallet, error)
    
    // Order operations
    PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error)
    CancelOrder(ctx context.Context, orderID string, symbol string) error
    GetOrder(ctx context.Context, orderID string, symbol string) (*models.Order, error)
    GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error)
    
    // WebSocket subscriptions
    SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error
    UnsubscribeFromTickers(ctx context.Context, symbols []string) error
    
    // New coins detection
    GetNewCoins(ctx context.Context) ([]*models.NewCoin, error)
}
```

### Core Service Interfaces

```go
package service

import (
    "context"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// NewCoinService defines operations related to new coin detection
type NewCoinService interface {
    DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error)
    ProcessNewCoins(ctx context.Context) error
    ArchiveOldCoins(ctx context.Context, daysOld int) error
}

// TradeService defines operations for trading logic
type TradeService interface {
    EvaluatePurchaseDecision(ctx context.Context, symbol string) (*models.PurchaseDecision, error)
    ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error)
    CheckStopLoss(ctx context.Context, coin *models.BoughtCoin) (bool, error)
    CheckTakeProfit(ctx context.Context, coin *models.BoughtCoin) (bool, error)
    SellCoin(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error)
}

// PortfolioService defines operations for portfolio management
type PortfolioService interface {
    GetPortfolioValue(ctx context.Context) (float64, error)
    GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error)
    GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error)
}
```

## 3. Additional Models for WebSocket and API

```go
package models

// Ticker represents real-time price information for a trading pair
type Ticker struct {
    Symbol        string    `json:"symbol"`
    Price         float64   `json:"price"`
    PriceChange   float64   `json:"price_change"`    // 24h price change
    PriceChangePct float64  `json:"price_change_pct"` // 24h price change percent
    Volume        float64   `json:"volume"`          // 24h volume
    High24h       float64   `json:"high_24h"`
    Low24h        float64   `json:"low_24h"`
    Timestamp     time.Time `json:"timestamp"`
}

// Kline represents a candlestick chart data point
type Kline struct {
    Symbol    string    `json:"symbol"`
    Interval  string    `json:"interval"`  // e.g., "1m", "5m", "1h"
    OpenTime  time.Time `json:"open_time"`
    CloseTime time.Time `json:"close_time"`
    Open      float64   `json:"open"`
    High      float64   `json:"high"`
    Low       float64   `json:"low"`
    Close     float64   `json:"close"`
    Volume    float64   `json:"volume"`
}

// PerformanceMetrics represents trading performance statistics
type PerformanceMetrics struct {
    TotalTrades        int       `json:"total_trades"`
    WinningTrades      int       `json:"winning_trades"`
    LosingTrades       int       `json:"losing_trades"`
    WinRate            float64   `json:"win_rate"`            // Percentage
    AverageProfit      float64   `json:"average_profit"`      // Percentage
    AverageLoss        float64   `json:"average_loss"`        // Percentage
    LargestProfit      float64   `json:"largest_profit"`      // Percentage
    LargestLoss        float64   `json:"largest_loss"`        // Percentage
    AverageHoldingTime float64   `json:"average_holding_time"` // In hours
    StartValue         float64   `json:"start_value"`
    CurrentValue       float64   `json:"current_value"`
    TotalProfit        float64   `json:"total_profit"`        // Absolute value
    TotalProfitPct     float64   `json:"total_profit_pct"`    // Percentage
}
```

## 4. Best Practices for Go Domain Models

1. **Use Proper Types**:
   - Prefer strongly typed enums (using const strings) over raw strings
   - Use pointers for optional fields (`*float64` instead of `float64`)
   - Use the appropriate numeric types (float64 for prices, int64 for IDs)

2. **Immutability and Value Objects**:
   - Consider making simple value objects immutable
   - Use constructor functions to ensure valid state

3. **Validation**:
   - Add validation methods to ensure domain objects are in a valid state
   - Consider implementing a `Validate() error` method on complex types

4. **Behavior and Data Together**:
   - Encapsulate behavior with data (methods on struct types)
   - Place business logic in domain methods where appropriate

5. **Error Handling**:
   - Define domain-specific errors in `internal/domain/errs`
   - Use error wrapping for context preservation

Example for validation and constructors:

```go
// NewBoughtCoin creates a new BoughtCoin with validation
func NewBoughtCoin(symbol string, price, quantity float64) (*BoughtCoin, error) {
    if symbol == "" {
        return nil, errors.New("symbol cannot be empty")
    }
    if price <= 0 {
        return nil, errors.New("price must be positive")
    }
    if quantity <= 0 {
        return nil, errors.New("quantity must be positive")
    }
    
    return &BoughtCoin{
        Symbol:        symbol,
        PurchasePrice: price,
        Quantity:      quantity,
        PurchaseTime:  time.Now(),
        IsDeleted:     false,
        StopLossPrice: price * 0.85, // 15% stop loss
        TakeProfitLevels: []TakeProfitLevel{
            {Percentage: 5, SellQuantity: quantity * 0.25, IsReached: false},
            {Percentage: 10, SellQuantity: quantity * 0.25, IsReached: false},
            {Percentage: 15, SellQuantity: quantity * 0.25, IsReached: false},
            {Percentage: 20, SellQuantity: quantity * 0.25, IsReached: false},
        },
    }, nil
}
```
