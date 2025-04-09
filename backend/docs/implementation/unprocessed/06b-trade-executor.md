# Trade Executor Implementation

This document provides implementation details for the Trade Executor component, which is responsible for executing trading operations in the Go crypto trading bot.

## Table of Contents

1. [Overview](#overview)
2. [Component Structure](#component-structure)
3. [Core Implementation](#core-implementation)
4. [Trading Strategies](#trading-strategies)
5. [Integration with Other Components](#integration-with-other-components)
6. [Testing Guidance](#testing-guidance)

## Overview

The Trade Executor is responsible for:

- Executing buy and sell orders
- Implementing trading decision logic
- Managing trade lifecycle
- Interacting with the exchange API
- Maintaining trading records

This component works closely with trading strategies to make informed decisions about when to buy and sell cryptocurrency assets.

## Component Structure

The Trade Executor is implemented in the following file structure:

```
internal/domain/core/trading/
├── executor.go      # Main executor implementation
├── strategy.go      # Strategy interface definitions
└── basic.go         # Basic trading strategy implementation
```

## Core Implementation

### Executor Definition

```go
// internal/domain/core/trading/executor.go
package trading

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/domain/service"
    "github.com/ryanlisse/cryptobot/internal/platform/exchange"
)

// Executor handles trade execution
type Executor struct {
    exchangeClient     exchange.Client
    boughtCoinRepo     repository.BoughtCoinRepository
    soldCoinRepo       repository.SoldCoinRepository
    decisionRepo       repository.PurchaseDecisionRepository
    strategy           Strategy
    walletService      service.WalletService
    positionManager    service.PositionManager
}

// NewExecutor creates a new trade executor
func NewExecutor(
    exchangeClient exchange.Client,
    boughtCoinRepo repository.BoughtCoinRepository,
    soldCoinRepo repository.SoldCoinRepository,
    decisionRepo repository.PurchaseDecisionRepository,
    strategy Strategy,
    walletService service.WalletService,
    positionManager service.PositionManager,
) *Executor {
    return &Executor{
        exchangeClient:     exchangeClient,
        boughtCoinRepo:     boughtCoinRepo,
        soldCoinRepo:       soldCoinRepo,
        decisionRepo:       decisionRepo,
        strategy:           strategy,
        walletService:      walletService,
        positionManager:    positionManager,
    }
}
```

### Purchase Execution

```go
// ExecutePurchase executes a purchase for a coin
func (e *Executor) ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
    // Check if we already own this coin
    existingCoin, err := e.boughtCoinRepo.FindBySymbol(ctx, symbol)
    if err != nil && !errors.Is(err, repository.ErrNotFound) {
        return nil, fmt.Errorf("error checking existing holdings: %w", err)
    }
    
    if existingCoin != nil && !existingCoin.IsDeleted {
        return nil, fmt.Errorf("already own coin %s", symbol)
    }
    
    // Get current balance
    balance, err := e.walletService.GetBalance(ctx)
    if err != nil {
        return nil, fmt.Errorf("error getting balance: %w", err)
    }
    
    if balance < amount {
        return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", balance, amount)
    }
    
    // Get current price
    ticker, err := e.exchangeClient.GetTicker(ctx, symbol)
    if err != nil {
        return nil, fmt.Errorf("error getting ticker: %w", err)
    }
    
    // Calculate quantity
    quantity := amount / ticker.Price
    
    // Execute the purchase
    orderID, err := e.exchangeClient.PlaceBuyOrder(ctx, symbol, quantity)
    if err != nil {
        return nil, fmt.Errorf("error placing buy order: %w", err)
    }
    
    // Save the bought coin
    boughtCoin := &models.BoughtCoin{
        Symbol:        symbol,
        PurchasePrice: ticker.Price,
        Quantity:      quantity,
        PurchasedAt:   time.Now(),
        OrderID:       orderID,
    }
    
    savedCoin, err := e.boughtCoinRepo.Create(ctx, boughtCoin)
    if err != nil {
        return nil, fmt.Errorf("error saving bought coin: %w", err)
    }
    
    // Update wallet balance
    if err := e.walletService.UpdateBalance(ctx, -amount, fmt.Sprintf("Purchase of %s", symbol)); err != nil {
        return nil, fmt.Errorf("error updating balance: %w", err)
    }
    
    // Create position if position management is enabled
    if e.positionManager != nil {
        position := &models.Position{
            Symbol:        symbol,
            EntryPrice:    ticker.Price,
            Quantity:      quantity,
            OpenedAt:      time.Now(),
            CoinID:        savedCoin.ID,
            Status:        models.PositionStatusOpen,
            StopLoss:      ticker.Price * 0.95, // Default 5% stop loss
            TakeProfit:    ticker.Price * 1.10, // Default 10% take profit
        }
        
        if _, err := e.positionManager.CreatePosition(ctx, position); err != nil {
            log.Printf("Warning: Failed to create position for %s: %v", symbol, err)
        }
    }
    
    return savedCoin, nil
}
```

### Sale Execution

```go
// ExecuteSale sells a previously purchased coin
func (e *Executor) ExecuteSale(ctx context.Context, coinID int64) (*models.SoldCoin, error) {
    // Get the bought coin
    boughtCoin, err := e.boughtCoinRepo.FindByID(ctx, coinID)
    if err != nil {
        return nil, fmt.Errorf("error finding coin: %w", err)
    }
    
    if boughtCoin.IsDeleted {
        return nil, errors.New("coin has already been sold")
    }
    
    // Get current price
    ticker, err := e.exchangeClient.GetTicker(ctx, boughtCoin.Symbol)
    if err != nil {
        return nil, fmt.Errorf("error getting ticker: %w", err)
    }
    
    // Execute the sale
    orderID, err := e.exchangeClient.PlaceSellOrder(ctx, boughtCoin.Symbol, boughtCoin.Quantity)
    if err != nil {
        return nil, fmt.Errorf("error placing sell order: %w", err)
    }
    
    // Calculate profit/loss
    saleAmount := boughtCoin.Quantity * ticker.Price
    purchaseAmount := boughtCoin.Quantity * boughtCoin.PurchasePrice
    profitLoss := saleAmount - purchaseAmount
    profitLossPercentage := (profitLoss / purchaseAmount) * 100
    
    // Create sold coin record
    soldCoin := &models.SoldCoin{
        BoughtCoinID:         boughtCoin.ID,
        Symbol:               boughtCoin.Symbol,
        SellPrice:            ticker.Price,
        Quantity:             boughtCoin.Quantity,
        SoldAt:               time.Now(),
        ProfitLoss:           profitLoss,
        ProfitLossPercentage: profitLossPercentage,
        OrderID:              orderID,
    }
    
    savedSoldCoin, err := e.soldCoinRepo.Create(ctx, soldCoin)
    if err != nil {
        return nil, fmt.Errorf("error saving sold coin: %w", err)
    }
    
    // Mark the bought coin as deleted
    boughtCoin.IsDeleted = true
    if _, err := e.boughtCoinRepo.Update(ctx, boughtCoin); err != nil {
        log.Printf("Warning: Failed to mark coin as deleted: %v", err)
    }
    
    // Update wallet balance
    if err := e.walletService.UpdateBalance(ctx, saleAmount, 
        fmt.Sprintf("Sale of %s (P/L: %.2f%%)", boughtCoin.Symbol, profitLossPercentage)); err != nil {
        return nil, fmt.Errorf("error updating balance: %w", err)
    }
    
    // Close position if position management is enabled
    if e.positionManager != nil {
        if err := e.positionManager.ClosePositionByCoinID(ctx, boughtCoin.ID, ticker.Price); err != nil {
            log.Printf("Warning: Failed to close position for coin %d: %v", boughtCoin.ID, err)
        }
    }
    
    return savedSoldCoin, nil
}
```

### Decision Logic

```go
// EvaluateForPurchase evaluates whether to purchase a coin
func (e *Executor) EvaluateForPurchase(ctx context.Context, symbol string) (*models.PurchaseDecision, error) {
    // First check if we already own this coin
    existingCoin, err := e.boughtCoinRepo.FindBySymbol(ctx, symbol)
    if err != nil && !errors.Is(err, repository.ErrNotFound) {
        return nil, fmt.Errorf("error checking existing holdings: %w", err)
    }
    
    if existingCoin != nil && !existingCoin.IsDeleted {
        return &models.PurchaseDecision{
            Symbol:      symbol,
            EvaluatedAt: time.Now(),
            ShouldBuy:   false,
            Reason:      "Already own this coin",
        }, nil
    }
    
    // Let the strategy decide
    shouldBuy, reason, err := e.strategy.ShouldBuy(ctx, symbol)
    if err != nil {
        return nil, fmt.Errorf("strategy evaluation error: %w", err)
    }
    
    // Create decision record
    decision := &models.PurchaseDecision{
        Symbol:      symbol,
        EvaluatedAt: time.Now(),
        ShouldBuy:   shouldBuy,
        Reason:      reason,
    }
    
    // Save the decision
    savedDecision, err := e.decisionRepo.Create(ctx, decision)
    if err != nil {
        return nil, fmt.Errorf("error saving decision: %w", err)
    }
    
    return savedDecision, nil
}
```

## Trading Strategies

Trading strategies implement the decision-making logic for buying and selling coins:

```go
// internal/domain/core/trading/strategy.go
package trading

import (
    "context"
)

// Strategy defines the interface for trading strategies
type Strategy interface {
    // ShouldBuy determines if a coin should be purchased
    ShouldBuy(ctx context.Context, symbol string) (bool, string, error)
    
    // ShouldSell determines if a coin should be sold
    ShouldSell(ctx context.Context, coinID int64) (bool, string, error)
    
    // Name returns the strategy name
    Name() string
}
```

### Basic Strategy Implementation

```go
// internal/domain/core/trading/basic.go
package trading

import (
    "context"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/platform/exchange"
)

// BasicStrategy implements a simple trading strategy
type BasicStrategy struct {
    exchangeClient exchange.Client
    boughtCoinRepo repository.BoughtCoinRepository
    
    // Strategy parameters
    minVolume           float64
    maxPrice            float64
    profitTarget        float64
    stopLoss            float64
    maxHoldTimeMinutes  int
}

// NewBasicStrategy creates a new basic strategy
func NewBasicStrategy(
    exchangeClient exchange.Client,
    boughtCoinRepo repository.BoughtCoinRepository,
    minVolume float64,
    maxPrice float64,
    profitTarget float64,
    stopLoss float64,
    maxHoldTimeMinutes int,
) *BasicStrategy {
    return &BasicStrategy{
        exchangeClient:     exchangeClient,
        boughtCoinRepo:     boughtCoinRepo,
        minVolume:          minVolume,
        maxPrice:           maxPrice,
        profitTarget:       profitTarget,
        stopLoss:           stopLoss,
        maxHoldTimeMinutes: maxHoldTimeMinutes,
    }
}

// Name returns the strategy name
func (s *BasicStrategy) Name() string {
    return "BasicStrategy"
}

// ShouldBuy determines if a coin should be purchased
func (s *BasicStrategy) ShouldBuy(ctx context.Context, symbol string) (bool, string, error) {
    // Get ticker information
    ticker, err := s.exchangeClient.GetTicker(ctx, symbol)
    if err != nil {
        return false, "Error fetching ticker data", err
    }
    
    // Check volume
    if ticker.Volume < s.minVolume {
        return false, fmt.Sprintf("Insufficient volume: %.2f < %.2f", 
            ticker.Volume, s.minVolume), nil
    }
    
    // Check price
    if s.maxPrice > 0 && ticker.Price > s.maxPrice {
        return false, fmt.Sprintf("Price too high: %.8f > %.8f", 
            ticker.Price, s.maxPrice), nil
    }
    
    // Check recent price movement (last 24h)
    klines, err := s.exchangeClient.GetKlines(ctx, symbol, "1h", 24)
    if err != nil {
        return false, "Error fetching historical data", err
    }
    
    if len(klines) < 24 {
        return false, "Insufficient historical data", nil
    }
    
    // Implement basic analysis, e.g., check if price is trending up
    uptrend := isUptrend(klines)
    if !uptrend {
        return false, "Price not in uptrend", nil
    }
    
    return true, "Meets basic criteria for purchase", nil
}

// ShouldSell determines if a coin should be sold
func (s *BasicStrategy) ShouldSell(ctx context.Context, coinID int64) (bool, string, error) {
    // Get the bought coin
    boughtCoin, err := s.boughtCoinRepo.FindByID(ctx, coinID)
    if err != nil {
        return false, "Error retrieving coin data", err
    }
    
    // Get current price
    ticker, err := s.exchangeClient.GetTicker(ctx, boughtCoin.Symbol)
    if err != nil {
        return false, "Error fetching current price", err
    }
    
    // Calculate profit/loss percentage
    profitLoss := ((ticker.Price - boughtCoin.PurchasePrice) / boughtCoin.PurchasePrice) * 100
    
    // Check profit target
    if profitLoss >= s.profitTarget {
        return true, fmt.Sprintf("Reached profit target: %.2f%%", profitLoss), nil
    }
    
    // Check stop loss
    if profitLoss <= -s.stopLoss {
        return true, fmt.Sprintf("Hit stop loss: %.2f%%", profitLoss), nil
    }
    
    // Check max hold time
    holdTime := time.Since(boughtCoin.PurchasedAt)
    maxHoldTime := time.Duration(s.maxHoldTimeMinutes) * time.Minute
    
    if holdTime > maxHoldTime {
        return true, fmt.Sprintf("Exceeded max hold time: %s", holdTime), nil
    }
    
    return false, "Holding position", nil
}

// Helper function to detect uptrends
func isUptrend(klines []*models.Kline) bool {
    if len(klines) < 24 {
        return false
    }
    
    // Simple trend detection: closing price higher than opening for majority of recent candles
    upCount := 0
    for i := len(klines) - 12; i < len(klines); i++ {
        if klines[i].Close > klines[i].Open {
            upCount++
        }
    }
    
    // Consider it an uptrend if more than 50% of recent candles are positive
    return upCount > 6
}
```

## Integration with Other Components

The Trade Executor integrates with multiple components:

1. **Exchange Client**: For executing trades and retrieving market data
2. **Repository Layer**: For persisting trade data
3. **Wallet Service**: For managing account balances
4. **Position Manager**: For position management and risk control

### Factory Creation Example

```go
// internal/domain/factory/factory.go
package factory

// CreateTradeExecutor creates and configures a trade executor
func (f *ServiceFactory) CreateTradeExecutor() *trading.Executor {
    // Create basic strategy with config parameters
    strategy := trading.NewBasicStrategy(
        f.exchangeClient,
        f.repositories.BoughtCoinRepository,
        f.config.MinVolume,
        f.config.MaxPrice,
        f.config.ProfitTarget,
        f.config.StopLoss,
        f.config.MaxHoldTimeMinutes,
    )
    
    // Create wallet service
    walletService := f.CreateWalletService()
    
    // Create position manager if enabled
    var positionManager service.PositionManager
    if f.config.EnablePositionManagement {
        positionManager = f.CreatePositionManager()
    }
    
    return trading.NewExecutor(
        f.exchangeClient,
        f.repositories.BoughtCoinRepository,
        f.repositories.SoldCoinRepository,
        f.repositories.PurchaseDecisionRepository,
        strategy,
        walletService,
        positionManager,
    )
}
```

## Testing Guidance

### Unit Testing Example

```go
// internal/domain/core/trading/executor_test.go
package trading_test

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/ryanlisse/cryptobot/internal/domain/core/trading"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/domain/repository/mocks"
    "github.com/ryanlisse/cryptobot/internal/platform/exchange/mocks"
    "github.com/ryanlisse/cryptobot/internal/domain/service/mocks"
)

func TestExecutor_ExecutePurchase(t *testing.T) {
    // Create mocks
    mockExchange := new(mocks.ExchangeClient)
    mockBoughtRepo := new(mocks.BoughtCoinRepository)
    mockSoldRepo := new(mocks.SoldCoinRepository)
    mockDecisionRepo := new(mocks.PurchaseDecisionRepository)
    mockStrategy := new(mocks.Strategy)
    mockWalletService := new(mocks.WalletService)
    mockPositionManager := new(mocks.PositionManager)
    
    // Test data
    symbol := "BTCUSDT"
    amount := 100.0
    price := 50000.0
    quantity := amount / price
    orderID := "order123"
    
    // Setup expectations
    mockBoughtRepo.On("FindBySymbol", mock.Anything, symbol).
        Return(nil, repository.ErrNotFound)
    
    mockWalletService.On("GetBalance", mock.Anything).
        Return(1000.0, nil)
    
    mockExchange.On("GetTicker", mock.Anything, symbol).
        Return(&models.Ticker{Symbol: symbol, Price: price}, nil)
    
    mockExchange.On("PlaceBuyOrder", mock.Anything, symbol, quantity).
        Return(orderID, nil)
    
    mockBoughtRepo.On("Create", mock.Anything, mock.MatchedBy(func(coin *models.BoughtCoin) bool {
        return coin.Symbol == symbol &&
               coin.PurchasePrice == price &&
               coin.Quantity == quantity
    })).Return(&models.BoughtCoin{ID: 1, Symbol: symbol, PurchasePrice: price, Quantity: quantity}, nil)
    
    mockWalletService.On("UpdateBalance", mock.Anything, -amount, mock.Anything).
        Return(nil)
    
    mockPositionManager.On("CreatePosition", mock.Anything, mock.Anything).
        Return(&models.Position{ID: 1}, nil)
    
    // Create executor
    executor := trading.NewExecutor(
        mockExchange,
        mockBoughtRepo,
        mockSoldRepo,
        mockDecisionRepo,
        mockStrategy,
        mockWalletService,
        mockPositionManager,
    )
    
    // Execute test
    result, err := executor.ExecutePurchase(context.Background(), symbol, amount)
    
    // Verify results
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, symbol, result.Symbol)
    assert.Equal(t, price, result.PurchasePrice)
    assert.Equal(t, quantity, result.Quantity)
    
    // Verify mock expectations
    mockExchange.AssertExpectations(t)
    mockBoughtRepo.AssertExpectations(t)
    mockWalletService.AssertExpectations(t)
    mockPositionManager.AssertExpectations(t)
}
```

For more comprehensive testing examples and best practices, refer to the [Testing Strategy](../testing/overview.md) document.
