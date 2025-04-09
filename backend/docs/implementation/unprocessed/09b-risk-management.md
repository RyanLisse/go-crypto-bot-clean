# Risk Management Implementation

This document outlines the implementation approach for the risk management component of the Go crypto trading bot. Risk management is a critical part of any trading system, ensuring capital preservation and managing potential losses.

## Table of Contents

1. [Overview](#overview)
2. [Core Concepts](#core-concepts)
3. [Component Structure](#component-structure)
4. [Implementation Details](#implementation-details)
5. [Position Sizing](#position-sizing)
6. [Stop-Loss Strategies](#stop-loss-strategies)
7. [Integration with Other Components](#integration-with-other-components)
8. [Testing Strategy](#testing-strategy)

## Overview

The Risk Management system is responsible for:

- Calculating appropriate position sizes based on account balance and risk parameters
- Implementing various stop-loss strategies to limit potential losses
- Managing overall portfolio risk exposure
- Preventing excessive trading during unfavorable market conditions
- Providing risk metrics and alerts

This component works closely with the Position Management and Trade Executor components to ensure trading decisions adhere to predefined risk parameters.

## Core Concepts

### Risk-Reward Ratio

The ratio between the potential profit and potential loss of a trade:

```go
// RiskRewardRatio calculates the ratio between potential profit and potential loss
func RiskRewardRatio(entryPrice, targetPrice, stopLossPrice float64) float64 {
    potentialProfit := math.Abs(targetPrice - entryPrice)
    potentialLoss := math.Abs(entryPrice - stopLossPrice)
    
    if potentialLoss == 0 {
        return 0 // Avoid division by zero
    }
    
    return potentialProfit / potentialLoss
}
```

### Maximum Drawdown

The maximum observed loss from a peak to a trough of the account balance:

```go
// CalculateMaxDrawdown finds the maximum drawdown percentage from a series of balances
func CalculateMaxDrawdown(balances []float64) float64 {
    if len(balances) == 0 {
        return 0
    }
    
    maxBalance := balances[0]
    maxDrawdown := 0.0
    
    for _, balance := range balances {
        if balance > maxBalance {
            maxBalance = balance
        }
        
        drawdown := (maxBalance - balance) / maxBalance
        if drawdown > maxDrawdown {
            maxDrawdown = drawdown
        }
    }
    
    return maxDrawdown
}
```

### Position Size

The amount of capital allocated to a single trade, determined by risk parameters:

```go
// PositionSize represents the calculation of an appropriate position size
type PositionSize struct {
    Symbol          string  `json:"symbol"`
    EntryPrice      float64 `json:"entry_price"`
    StopLossPrice   float64 `json:"stop_loss_price"`
    AccountBalance  float64 `json:"account_balance"`
    RiskPercentage  float64 `json:"risk_percentage"`
    Quantity        float64 `json:"quantity"`
    PositionValue   float64 `json:"position_value"`
    RiskAmount      float64 `json:"risk_amount"`
    MaxPositionSize float64 `json:"max_position_size"`
}
```

## Component Structure

The Risk Management system is structured as follows:

```
internal/domain/risk/
├── manager.go         # Main risk manager implementation
├── position_sizing.go # Position size calculation logic
├── stop_loss.go       # Stop-loss strategy implementations
└── metrics.go         # Risk metrics calculations
```

## Implementation Details

### Risk Manager Interface

```go
// internal/domain/service/risk_service.go
package service

import (
    "context"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// RiskManager defines the operations for risk management
type RiskManager interface {
    // CalculatePositionSize determines the appropriate position size based on risk parameters
    CalculatePositionSize(ctx context.Context, symbol string, entryPrice, stopLossPrice float64) (*models.PositionSize, error)
    
    // CheckRiskParameters validates if a trade meets the risk parameters
    CheckRiskParameters(ctx context.Context, trade *models.TradeRequest) (bool, string, error)
    
    // UpdateStopLoss updates a stop-loss based on the selected strategy
    UpdateStopLoss(ctx context.Context, positionID int64, strategy string, params map[string]interface{}) (float64, error)
    
    // CalculateRiskMetrics returns the current risk metrics for the portfolio
    CalculateRiskMetrics(ctx context.Context) (*models.RiskMetrics, error)
    
    // GetMaxDrawdown calculates the maximum drawdown over a period
    GetMaxDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, error)
    
    // SetRiskParameters updates the global risk parameters
    SetRiskParameters(ctx context.Context, params map[string]interface{}) error
    
    // GetRiskParameters returns the current risk parameters
    GetRiskParameters(ctx context.Context) (map[string]interface{}, error)
}
```

### Risk Manager Implementation

```go
// internal/domain/risk/manager.go
package risk

import (
    "context"
    "errors"
    "fmt"
    "math"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
)

// Manager implements the RiskManager interface
type Manager struct {
    accountRepo       repository.AccountRepository
    positionRepo      repository.PositionRepository
    transactionRepo   repository.TransactionRepository
    
    maxRiskPerTrade   float64 // Maximum percentage of account to risk per trade
    maxOpenTrades     int     // Maximum number of concurrent open trades
    minRiskRewardRatio float64 // Minimum acceptable risk-reward ratio
    maxDailyLoss      float64 // Maximum percentage loss allowed per day
    maxDrawdown       float64 // Maximum allowable drawdown before halting trading
    maxPositionSize   float64 // Maximum size of a position as a percentage of account
}

// NewManager creates a new risk manager
func NewManager(
    accountRepo repository.AccountRepository,
    positionRepo repository.PositionRepository,
    transactionRepo repository.TransactionRepository,
    params map[string]interface{},
) *Manager {
    // Set default risk parameters
    maxRiskPerTrade := 0.01   // 1% per trade
    maxOpenTrades := 5
    minRiskRewardRatio := 2.0 // 1:2 risk-reward ratio
    maxDailyLoss := 0.05      // 5% max daily loss
    maxDrawdown := 0.20       // 20% max drawdown
    maxPositionSize := 0.20   // 20% max position size
    
    // Override with provided parameters if available
    if val, ok := params["maxRiskPerTrade"].(float64); ok {
        maxRiskPerTrade = val
    }
    if val, ok := params["maxOpenTrades"].(int); ok {
        maxOpenTrades = val
    }
    if val, ok := params["minRiskRewardRatio"].(float64); ok {
        minRiskRewardRatio = val
    }
    if val, ok := params["maxDailyLoss"].(float64); ok {
        maxDailyLoss = val
    }
    if val, ok := params["maxDrawdown"].(float64); ok {
        maxDrawdown = val
    }
    if val, ok := params["maxPositionSize"].(float64); ok {
        maxPositionSize = val
    }
    
    return &Manager{
        accountRepo:       accountRepo,
        positionRepo:      positionRepo,
        transactionRepo:   transactionRepo,
        maxRiskPerTrade:   maxRiskPerTrade,
        maxOpenTrades:     maxOpenTrades,
        minRiskRewardRatio: minRiskRewardRatio,
        maxDailyLoss:      maxDailyLoss,
        maxDrawdown:       maxDrawdown,
        maxPositionSize:   maxPositionSize,
    }
}

// CalculatePositionSize determines the appropriate position size based on risk parameters
func (m *Manager) CalculatePositionSize(
    ctx context.Context,
    symbol string,
    entryPrice,
    stopLossPrice float64,
) (*models.PositionSize, error) {
    // Get account balance
    account, err := m.accountRepo.GetAccount(ctx)
    if err != nil {
        return nil, fmt.Errorf("error getting account: %w", err)
    }
    
    // Calculate risk amount
    riskAmount := account.Balance * m.maxRiskPerTrade
    
    // Calculate risk per unit
    riskPerUnit := math.Abs(entryPrice - stopLossPrice)
    if riskPerUnit == 0 {
        return nil, errors.New("stop loss must be different from entry price")
    }
    
    // Calculate quantity based on risk amount
    quantity := riskAmount / riskPerUnit
    
    // Calculate position value
    positionValue := quantity * entryPrice
    
    // Check if position value exceeds max position size
    maxPositionValue := account.Balance * m.maxPositionSize
    if positionValue > maxPositionValue {
        // Adjust quantity to respect max position size
        quantity = maxPositionValue / entryPrice
        positionValue = quantity * entryPrice
        riskAmount = quantity * riskPerUnit
    }
    
    return &models.PositionSize{
        Symbol:          symbol,
        EntryPrice:      entryPrice,
        StopLossPrice:   stopLossPrice,
        AccountBalance:  account.Balance,
        RiskPercentage:  m.maxRiskPerTrade * 100,
        Quantity:        quantity,
        PositionValue:   positionValue,
        RiskAmount:      riskAmount,
        MaxPositionSize: maxPositionValue,
    }, nil
}

// CheckRiskParameters validates if a trade meets the risk parameters
func (m *Manager) CheckRiskParameters(
    ctx context.Context,
    trade *models.TradeRequest,
) (bool, string, error) {
    // Check risk-reward ratio if take profit is specified
    if trade.TakeProfit > 0 && trade.StopLoss > 0 {
        var riskRewardRatio float64
        
        if trade.Direction == "buy" {
            // For long positions
            potentialProfit := trade.TakeProfit - trade.Price
            potentialLoss := trade.Price - trade.StopLoss
            riskRewardRatio = potentialProfit / potentialLoss
        } else {
            // For short positions
            potentialProfit := trade.Price - trade.TakeProfit
            potentialLoss := trade.StopLoss - trade.Price
            riskRewardRatio = potentialProfit / potentialLoss
        }
        
        if riskRewardRatio < m.minRiskRewardRatio {
            return false, fmt.Sprintf("risk-reward ratio %.2f is below minimum %.2f", 
                riskRewardRatio, m.minRiskRewardRatio), nil
        }
    }
    
    // Check if max open trades would be exceeded
    openPositions, err := m.positionRepo.FindByStatus(ctx, models.PositionStatusOpen)
    if err != nil {
        return false, "", fmt.Errorf("error checking open positions: %w", err)
    }
    
    if len(openPositions) >= m.maxOpenTrades {
        return false, fmt.Sprintf("maximum number of open trades (%d) would be exceeded", 
            m.maxOpenTrades), nil
    }
    
    // Check daily loss limit
    dayStart := time.Now().Truncate(24 * time.Hour)
    transactions, err := m.transactionRepo.FindByTimeRange(ctx, dayStart, time.Now())
    if err != nil {
        return false, "", fmt.Errorf("error checking daily transactions: %w", err)
    }
    
    // Calculate daily P/L
    dailyPL := 0.0
    for _, tx := range transactions {
        if tx.Type == "trade_pl" {
            dailyPL += tx.Amount
        }
    }
    
    account, err := m.accountRepo.GetAccount(ctx)
    if err != nil {
        return false, "", fmt.Errorf("error getting account: %w", err)
    }
    
    dailyLossPercentage := math.Abs(dailyPL) / account.Balance
    if dailyPL < 0 && dailyLossPercentage >= m.maxDailyLoss {
        return false, fmt.Sprintf("daily loss limit of %.2f%% reached (current: %.2f%%)", 
            m.maxDailyLoss*100, dailyLossPercentage*100), nil
    }
    
    // Check drawdown limit
    drawdown, err := m.GetMaxDrawdown(ctx, time.Now().AddDate(0, -1, 0), time.Now())
    if err != nil {
        return false, "", fmt.Errorf("error checking drawdown: %w", err)
    }
    
    if drawdown >= m.maxDrawdown {
        return false, fmt.Sprintf("maximum drawdown of %.2f%% reached (current: %.2f%%)", 
            m.maxDrawdown*100, drawdown*100), nil
    }
    
    return true, "", nil
}

// Additional methods would be implemented here...
```

## Position Sizing

Position sizing determines how much capital to allocate to each trade, based on predefined risk parameters.

### Fixed Percentage Risk

This approach risks a fixed percentage of the account balance on each trade:

```go
// internal/domain/risk/position_sizing.go
package risk

import (
    "errors"
    "math"
)

// CalculateFixedPercentRisk calculates position size based on fixed percent risk
func CalculateFixedPercentRisk(
    accountBalance float64,
    riskPercentage float64,
    entryPrice float64,
    stopLossPrice float64,
) (float64, error) {
    if stopLossPrice >= entryPrice {
        return 0, errors.New("stop loss must be below entry price for long positions")
    }
    
    riskPerUnit := entryPrice - stopLossPrice
    if riskPerUnit == 0 {
        return 0, errors.New("stop loss must be different from entry price")
    }
    
    riskAmount := accountBalance * riskPercentage
    quantity := riskAmount / riskPerUnit
    
    return quantity, nil
}

// CalculateFixedPositionSize calculates position based on a fixed position size
func CalculateFixedPositionSize(
    accountBalance float64,
    positionSizePercentage float64,
    entryPrice float64,
) (float64, error) {
    if positionSizePercentage <= 0 || positionSizePercentage > 1 {
        return 0, errors.New("position size percentage must be between 0 and 1")
    }
    
    positionValue := accountBalance * positionSizePercentage
    quantity := positionValue / entryPrice
    
    return quantity, nil
}

// CalculateKellyPositionSize uses the Kelly Criterion for position sizing
func CalculateKellyPositionSize(
    accountBalance float64,
    winRate float64, // Historical win rate (0-1)
    avgWinLossRatio float64, // Average win/loss ratio
    maxKellyPercentage float64, // Maximum percentage of Kelly to use (0-1)
) (float64, error) {
    if winRate <= 0 || winRate >= 1 {
        return 0, errors.New("win rate must be between 0 and 1 exclusive")
    }
    
    if avgWinLossRatio <= 0 {
        return 0, errors.New("win/loss ratio must be positive")
    }
    
    if maxKellyPercentage <= 0 || maxKellyPercentage > 1 {
        return 0, errors.New("max Kelly percentage must be between 0 and 1")
    }
    
    // Kelly formula: f* = (bp - q) / b
    // Where:
    // - f* is the fraction of the current bankroll to wager
    // - b is the net odds received on the wager (b to 1)
    // - p is the probability of winning
    // - q is the probability of losing (1 - p)
    
    kellyPercentage := (winRate * avgWinLossRatio - (1 - winRate)) / avgWinLossRatio
    
    // Apply the maximum Kelly percentage to reduce risk
    adjustedKelly := math.Min(kellyPercentage, maxKellyPercentage)
    
    // Ensure Kelly is not negative
    if adjustedKelly <= 0 {
        return 0, errors.New("Kelly criterion suggests not taking this trade")
    }
    
    return accountBalance * adjustedKelly, nil
}
```

## Stop-Loss Strategies

Different stop-loss strategies to protect capital and limit potential losses:

```go
// internal/domain/risk/stop_loss.go
package risk

import (
    "errors"
    "math"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// CalculateFixedStopLoss calculates a fixed percentage stop-loss
func CalculateFixedStopLoss(
    entryPrice float64,
    percentageRisk float64,
    isLong bool,
) (float64, error) {
    if percentageRisk <= 0 {
        return 0, errors.New("percentage risk must be positive")
    }
    
    if isLong {
        return entryPrice * (1 - percentageRisk), nil
    }
    
    return entryPrice * (1 + percentageRisk), nil
}

// CalculateATRStopLoss calculates a stop-loss based on Average True Range
func CalculateATRStopLoss(
    entryPrice float64,
    atrValue float64,
    atrMultiplier float64,
    isLong bool,
) (float64, error) {
    if atrValue <= 0 {
        return 0, errors.New("ATR value must be positive")
    }
    
    if atrMultiplier <= 0 {
        return 0, errors.New("ATR multiplier must be positive")
    }
    
    atrDistance := atrValue * atrMultiplier
    
    if isLong {
        return entryPrice - atrDistance, nil
    }
    
    return entryPrice + atrDistance, nil
}

// CalculateSwingStopLoss calculates a stop-loss based on recent swing points
func CalculateSwingStopLoss(
    entryPrice float64,
    candles []*models.Candle,
    lookbackPeriod int,
    buffer float64,
    isLong bool,
) (float64, error) {
    if len(candles) < lookbackPeriod {
        return 0, errors.New("not enough candles for swing calculation")
    }
    
    if isLong {
        // For long positions, find recent swing low
        lowestLow := candles[0].Low
        for i := 1; i < lookbackPeriod; i++ {
            if candles[i].Low < lowestLow {
                lowestLow = candles[i].Low
            }
        }
        
        // Apply buffer
        return lowestLow * (1 - buffer), nil
    } else {
        // For short positions, find recent swing high
        highestHigh := candles[0].High
        for i := 1; i < lookbackPeriod; i++ {
            if candles[i].High > highestHigh {
                highestHigh = candles[i].High
            }
        }
        
        // Apply buffer
        return highestHigh * (1 + buffer), nil
    }
}

// CalculateTrailingStopLoss updates a trailing stop-loss based on price movement
func CalculateTrailingStopLoss(
    currentPrice float64,
    currentStopLoss float64,
    trailPercentage float64,
    isLong bool,
) (float64, error) {
    if trailPercentage <= 0 {
        return 0, errors.New("trail percentage must be positive")
    }
    
    if isLong {
        // For long positions, move stop loss up if price increases
        potentialStopLoss := currentPrice * (1 - trailPercentage)
        if potentialStopLoss > currentStopLoss {
            return potentialStopLoss, nil
        }
    } else {
        // For short positions, move stop loss down if price decreases
        potentialStopLoss := currentPrice * (1 + trailPercentage)
        if potentialStopLoss < currentStopLoss {
            return potentialStopLoss, nil
        }
    }
    
    // No update needed
    return currentStopLoss, nil
}

// CalculateBreakEvenStopLoss moves the stop-loss to break-even after a certain profit threshold
func CalculateBreakEvenStopLoss(
    entryPrice float64,
    currentPrice float64,
    currentStopLoss float64,
    profitThresholdPercentage float64,
    bufferPercentage float64,
    isLong bool,
) (float64, error) {
    if profitThresholdPercentage <= 0 {
        return 0, errors.New("profit threshold percentage must be positive")
    }
    
    if isLong {
        // Calculate current profit percentage
        profitPercentage := (currentPrice - entryPrice) / entryPrice
        
        // Check if profit threshold is reached
        if profitPercentage >= profitThresholdPercentage {
            // Move stop loss to break-even plus a small buffer
            return entryPrice * (1 + bufferPercentage), nil
        }
    } else {
        // Calculate current profit percentage for short
        profitPercentage := (entryPrice - currentPrice) / entryPrice
        
        // Check if profit threshold is reached
        if profitPercentage >= profitThresholdPercentage {
            // Move stop loss to break-even plus a small buffer
            return entryPrice * (1 - bufferPercentage), nil
        }
    }
    
    // No update needed
    return currentStopLoss, nil
}
```

## Integration with Other Components

The Risk Management system integrates with several other components:

1. **Position Management**: Provides position size recommendations and stop-loss updates
2. **Trade Executor**: Validates trades against risk parameters before execution
3. **Account Manager**: Monitors account balance and drawdown limits
4. **Dashboard**: Displays risk metrics and alerts

### Factory Creation Example

```go
// internal/domain/factory/factory.go
package factory

// CreateRiskManager creates and configures a risk manager
func (f *ServiceFactory) CreateRiskManager() *risk.Manager {
    riskParams := map[string]interface{}{
        "maxRiskPerTrade":    f.config.MaxRiskPerTrade,
        "maxOpenTrades":      f.config.MaxOpenTrades,
        "minRiskRewardRatio": f.config.MinRiskRewardRatio,
        "maxDailyLoss":       f.config.MaxDailyLoss,
        "maxDrawdown":        f.config.MaxDrawdown,
        "maxPositionSize":    f.config.MaxPositionSize,
    }
    
    return risk.NewManager(
        f.repositories.AccountRepository,
        f.repositories.PositionRepository,
        f.repositories.TransactionRepository,
        riskParams,
    )
}
```

## Testing Strategy

### Unit Testing

```go
// internal/domain/risk/manager_test.go
package risk_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/risk"
    "github.com/ryanlisse/cryptobot/internal/domain/repository/mocks"
)

func TestCalculatePositionSize(t *testing.T) {
    // Create mocks
    accountRepo := new(mocks.AccountRepository)
    positionRepo := new(mocks.PositionRepository)
    transactionRepo := new(mocks.TransactionRepository)
    
    // Configure mocks
    accountRepo.On("GetAccount", mock.Anything).Return(&models.Account{
        Balance: 10000.0,
    }, nil)
    
    // Create risk manager with parameters
    riskParams := map[string]interface{}{
        "maxRiskPerTrade": 0.01,  // 1% risk per trade
        "maxPositionSize": 0.20,  // 20% max position size
    }
    
    manager := risk.NewManager(accountRepo, positionRepo, transactionRepo, riskParams)
    
    // Test case
    posSize, err := manager.CalculatePositionSize(
        context.Background(),
        "BTCUSDT",
        50000.0,  // Entry price
        48500.0,  // Stop loss price (3% below entry)
    )
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, posSize)
    
    // Check risk amount (should be 1% of account)
    expectedRiskAmount := 10000.0 * 0.01 // = 100.0
    assert.InDelta(t, expectedRiskAmount, posSize.RiskAmount, 0.01)
    
    // Check position value doesn't exceed max position size
    maxPositionValue := 10000.0 * 0.20 // = 2000.0
    assert.LessOrEqual(t, posSize.PositionValue, maxPositionValue)
    
    // Verify mocks
    accountRepo.AssertExpectations(t)
}

func TestCheckRiskParameters(t *testing.T) {
    // Create mocks
    accountRepo := new(mocks.AccountRepository)
    positionRepo := new(mocks.PositionRepository)
    transactionRepo := new(mocks.TransactionRepository)
    
    // Configure mocks
    accountRepo.On("GetAccount", mock.Anything).Return(&models.Account{
        Balance: 10000.0,
    }, nil)
    
    positionRepo.On("FindByStatus", mock.Anything, "OPEN").Return([]*models.Position{
        {ID: 1, Symbol: "ETHUSDT"},
        {ID: 2, Symbol: "ADAUSDT"},
    }, nil)
    
    dayStart := time.Now().Truncate(24 * time.Hour)
    transactionRepo.On("FindByTimeRange", mock.Anything, dayStart, mock.Anything).Return([]*models.Transaction{
        {Type: "trade_pl", Amount: -100.0}, // Loss of $100
    }, nil)
    
    // Create risk manager with parameters
    riskParams := map[string]interface{}{
        "maxRiskPerTrade":    0.01,  // 1% risk per trade
        "maxOpenTrades":      5,     // Max 5 open trades
        "minRiskRewardRatio": 2.0,   // Minimum 1:2 risk-reward ratio
        "maxDailyLoss":       0.05,  // 5% max daily loss
    }
    
    manager := risk.NewManager(accountRepo, positionRepo, transactionRepo, riskParams)
    
    // Test case with valid trade
    validTrade := &models.TradeRequest{
        Symbol:     "BTCUSDT",
        Direction:  "buy",
        Price:      50000.0,
        StopLoss:   48500.0,  // 3% below entry
        TakeProfit: 55000.0,  // 10% above entry (R:R ratio = 3.33)
    }
    
    allowed, _, err := manager.CheckRiskParameters(context.Background(), validTrade)
    
    // Assertions
    assert.NoError(t, err)
    assert.True(t, allowed)
    
    // Test case with invalid risk-reward ratio
    invalidTrade := &models.TradeRequest{
        Symbol:     "BTCUSDT",
        Direction:  "buy",
        Price:      50000.0,
        StopLoss:   48500.0,  // 3% below entry
        TakeProfit: 51000.0,  // 2% above entry (R:R ratio = 0.67)
    }
    
    allowed, reason, err := manager.CheckRiskParameters(context.Background(), invalidTrade)
    
    // Assertions
    assert.NoError(t, err)
    assert.False(t, allowed)
    assert.Contains(t, reason, "risk-reward ratio")
    
    // Verify mocks
    accountRepo.AssertExpectations(t)
    positionRepo.AssertExpectations(t)
    transactionRepo.AssertExpectations(t)
}
```

For more comprehensive testing examples and best practices, refer to the [Testing Strategy](../testing/overview.md) document.
