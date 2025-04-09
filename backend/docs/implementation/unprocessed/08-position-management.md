# Position Management Implementation

This document outlines the implementation details for the Position Management component of the Go crypto trading bot. Position management is essential for tracking trades, applying risk management rules, and executing take-profit and stop-loss strategies.

## Table of Contents

1. [Overview](#overview)
2. [Core Concepts](#core-concepts)
3. [Component Structure](#component-structure)
4. [Implementation Details](#implementation-details)
5. [Position Lifecycle](#position-lifecycle)
6. [Risk Management](#risk-management)
7. [Integration with Other Components](#integration-with-other-components)
8. [Testing Strategy](#testing-strategy)

## Overview

The Position Management system tracks and manages trading positions from opening to closing. It implements:

- Position creation when trades are executed
- Stop-loss and take-profit functionality
- Position monitoring and adjustments
- Risk management rules
- Position performance tracking

This component works closely with the Trade Executor and Account Manager to ensure proper trade execution and accurate financial reporting.

## Core Concepts

### Position

A Position represents an active trade that has been executed but not yet closed:

```go
// Position represents an active trading position
type Position struct {
    ID           int64     `json:"id" db:"id"`
    Symbol       string    `json:"symbol" db:"symbol"`
    EntryPrice   float64   `json:"entry_price" db:"entry_price"`
    Quantity     float64   `json:"quantity" db:"quantity"`
    OpenedAt     time.Time `json:"opened_at" db:"opened_at"`
    ClosedAt     time.Time `json:"closed_at,omitempty" db:"closed_at"`
    ExitPrice    float64   `json:"exit_price,omitempty" db:"exit_price"`
    ProfitLoss   float64   `json:"profit_loss" db:"profit_loss"`
    PLPercentage float64   `json:"pl_percentage" db:"pl_percentage"`
    Status       string    `json:"status" db:"status"`
    StopLoss     float64   `json:"stop_loss" db:"stop_loss"`
    TakeProfit   float64   `json:"take_profit" db:"take_profit"`
    TrailingStop float64   `json:"trailing_stop,omitempty" db:"trailing_stop"`
    CoinID       int64     `json:"coin_id" db:"coin_id"`
    Notes        string    `json:"notes,omitempty" db:"notes"`
}
```

### Position Status

Positions can be in one of the following states:

- `OPEN`: Position is active
- `CLOSED`: Position has been closed (either manually or via take-profit/stop-loss)
- `PENDING`: Position creation has been initiated but not confirmed
- `CANCELED`: Position was canceled before it was opened
- `PARTIALLY_CLOSED`: Only a portion of the position has been closed

## Component Structure

The Position Management system is structured as follows:

```
internal/domain/core/position/
├── manager.go       # Main position manager implementation
├── monitor.go       # Position monitoring and automated actions
├── adjustor.go      # Position adjustment logic
└── risk.go          # Risk management rules
```

## Implementation Details

### Position Manager Interface

```go
// internal/domain/service/position_service.go
package service

import (
    "context"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// PositionManager defines the operations for position management
type PositionManager interface {
    // CreatePosition creates a new position
    CreatePosition(ctx context.Context, position *models.Position) (*models.Position, error)
    
    // GetOpenPositions returns all currently open positions
    GetOpenPositions(ctx context.Context) ([]*models.Position, error)
    
    // GetPositionByID retrieves a position by its ID
    GetPositionByID(ctx context.Context, id int64) (*models.Position, error)
    
    // GetPositionsForSymbol returns all positions for a given symbol
    GetPositionsForSymbol(ctx context.Context, symbol string) ([]*models.Position, error)
    
    // ClosePosition closes a position with the provided exit price
    ClosePosition(ctx context.Context, positionID int64, exitPrice float64) error
    
    // ClosePositionByCoinID closes a position associated with a specific coin
    ClosePositionByCoinID(ctx context.Context, coinID int64, exitPrice float64) error
    
    // UpdateStopLoss updates the stop-loss level for a position
    UpdateStopLoss(ctx context.Context, positionID int64, newStopLoss float64) error
    
    // UpdateTakeProfit updates the take-profit level for a position
    UpdateTakeProfit(ctx context.Context, positionID int64, newTakeProfit float64) error
    
    // EnableTrailingStop enables a trailing stop for a position
    EnableTrailingStop(ctx context.Context, positionID int64, trailingDistance float64) error
    
    // GetPositionStats returns statistics about positions in a date range
    GetPositionStats(ctx context.Context, startTime, endTime time.Time) (*models.PositionStats, error)
}
```

### Position Manager Implementation

```go
// internal/domain/core/position/manager.go
package position

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
)

// Manager implements the PositionManager interface
type Manager struct {
    positionRepo repository.PositionRepository
    coinRepo     repository.BoughtCoinRepository
}

// NewManager creates a new position manager
func NewManager(
    positionRepo repository.PositionRepository,
    coinRepo repository.BoughtCoinRepository,
) *Manager {
    return &Manager{
        positionRepo: positionRepo,
        coinRepo:     coinRepo,
    }
}

// CreatePosition creates a new position
func (m *Manager) CreatePosition(ctx context.Context, position *models.Position) (*models.Position, error) {
    // Set initial values
    if position.Status == "" {
        position.Status = models.PositionStatusOpen
    }
    
    position.OpenedAt = time.Now()
    
    // Validate the position
    if position.Symbol == "" {
        return nil, errors.New("position symbol cannot be empty")
    }
    
    if position.Quantity <= 0 {
        return nil, errors.New("position quantity must be greater than zero")
    }
    
    if position.EntryPrice <= 0 {
        return nil, errors.New("position entry price must be greater than zero")
    }
    
    // Create the position
    return m.positionRepo.Create(ctx, position)
}

// GetOpenPositions returns all currently open positions
func (m *Manager) GetOpenPositions(ctx context.Context) ([]*models.Position, error) {
    return m.positionRepo.FindByStatus(ctx, models.PositionStatusOpen)
}

// ClosePosition closes a position with the provided exit price
func (m *Manager) ClosePosition(ctx context.Context, positionID int64, exitPrice float64) error {
    // Get the position
    position, err := m.positionRepo.FindByID(ctx, positionID)
    if err != nil {
        return fmt.Errorf("error finding position: %w", err)
    }
    
    // Check if already closed
    if position.Status == models.PositionStatusClosed {
        return errors.New("position is already closed")
    }
    
    // Calculate profit/loss
    positionValue := position.Quantity * position.EntryPrice
    exitValue := position.Quantity * exitPrice
    profitLoss := exitValue - positionValue
    plPercentage := (profitLoss / positionValue) * 100
    
    // Update position
    position.Status = models.PositionStatusClosed
    position.ClosedAt = time.Now()
    position.ExitPrice = exitPrice
    position.ProfitLoss = profitLoss
    position.PLPercentage = plPercentage
    
    _, err = m.positionRepo.Update(ctx, position)
    return err
}

// ClosePositionByCoinID closes a position associated with a specific coin
func (m *Manager) ClosePositionByCoinID(ctx context.Context, coinID int64, exitPrice float64) error {
    // Find position by coin ID
    position, err := m.positionRepo.FindByCoinID(ctx, coinID)
    if err != nil {
        return fmt.Errorf("error finding position: %w", err)
    }
    
    return m.ClosePosition(ctx, position.ID, exitPrice)
}

// UpdateStopLoss updates the stop-loss level for a position
func (m *Manager) UpdateStopLoss(ctx context.Context, positionID int64, newStopLoss float64) error {
    position, err := m.positionRepo.FindByID(ctx, positionID)
    if err != nil {
        return fmt.Errorf("error finding position: %w", err)
    }
    
    position.StopLoss = newStopLoss
    _, err = m.positionRepo.Update(ctx, position)
    return err
}

// UpdateTakeProfit updates the take-profit level for a position
func (m *Manager) UpdateTakeProfit(ctx context.Context, positionID int64, newTakeProfit float64) error {
    position, err := m.positionRepo.FindByID(ctx, positionID)
    if err != nil {
        return fmt.Errorf("error finding position: %w", err)
    }
    
    position.TakeProfit = newTakeProfit
    _, err = m.positionRepo.Update(ctx, position)
    return err
}

// GetPositionStats returns statistics about positions in a date range
func (m *Manager) GetPositionStats(ctx context.Context, startTime, endTime time.Time) (*models.PositionStats, error) {
    positions, err := m.positionRepo.FindByTimeRange(ctx, startTime, endTime)
    if err != nil {
        return nil, fmt.Errorf("error finding positions: %w", err)
    }
    
    stats := &models.PositionStats{
        TotalPositions: len(positions),
        WinningPositions: 0,
        LosingPositions: 0,
        TotalProfitLoss: 0,
        AverageProfitLoss: 0,
        LargestWin: 0,
        LargestLoss: 0,
        WinRate: 0,
        StartTime: startTime,
        EndTime: endTime,
    }
    
    if len(positions) == 0 {
        return stats, nil
    }
    
    // Calculate statistics
    for _, p := range positions {
        if p.Status != models.PositionStatusClosed {
            continue
        }
        
        stats.TotalProfitLoss += p.ProfitLoss
        
        if p.ProfitLoss > 0 {
            stats.WinningPositions++
            if p.ProfitLoss > stats.LargestWin {
                stats.LargestWin = p.ProfitLoss
            }
        } else {
            stats.LosingPositions++
            if p.ProfitLoss < stats.LargestLoss {
                stats.LargestLoss = p.ProfitLoss
            }
        }
    }
    
    closedPositions := stats.WinningPositions + stats.LosingPositions
    if closedPositions > 0 {
        stats.AverageProfitLoss = stats.TotalProfitLoss / float64(closedPositions)
        stats.WinRate = float64(stats.WinningPositions) / float64(closedPositions) * 100
    }
    
    return stats, nil
}
```

## Position Lifecycle

Positions follow a clear lifecycle in the system:

1. **Creation**: When a trade is executed via the Trade Executor
   - Initial stop-loss and take-profit levels are set
   - Position is saved to the database

2. **Monitoring**: The Position Monitor continuously checks positions
   - Compares current market prices against stop-loss/take-profit levels
   - Implements trailing stops if enabled
   - Signals when action is needed

3. **Adjustment**: Position parameters may be adjusted during its lifetime
   - Stop-loss may be moved to break-even after certain profit threshold
   - Take-profit may be adjusted based on market conditions
   - Position size can be partially reduced to lock in profits

4. **Closing**: Position is closed when a condition is met
   - Take-profit level is reached
   - Stop-loss level is reached
   - Manual decision to close
   - Maximum holding time is reached

### Position Monitor Implementation

```go
// internal/domain/core/position/monitor.go
package position

import (
    "context"
    "log"
    "sync"
    "time"

    "github.com/ryanlisse/cryptobot/internal/platform/exchange"
)

// Monitor watches positions and triggers actions when conditions are met
type Monitor struct {
    manager        *Manager
    exchangeClient exchange.Client
    checkInterval  time.Duration
    running        bool
    mutex          sync.RWMutex
    stop           chan struct{}
}

// NewMonitor creates a new position monitor
func NewMonitor(
    manager *Manager,
    exchangeClient exchange.Client,
    checkInterval time.Duration,
) *Monitor {
    if checkInterval <= 0 {
        checkInterval = 30 * time.Second
    }
    
    return &Monitor{
        manager:        manager,
        exchangeClient: exchangeClient,
        checkInterval:  checkInterval,
        stop:           make(chan struct{}),
    }
}

// Start begins monitoring positions
func (m *Monitor) Start(ctx context.Context) error {
    m.mutex.Lock()
    if m.running {
        m.mutex.Unlock()
        return errors.New("monitor is already running")
    }
    m.running = true
    m.mutex.Unlock()
    
    go m.monitorLoop(ctx)
    return nil
}

// Stop stops the monitor
func (m *Monitor) Stop() {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    if !m.running {
        return
    }
    
    close(m.stop)
    m.running = false
}

// monitorLoop is the main monitoring loop
func (m *Monitor) monitorLoop(ctx context.Context) {
    ticker := time.NewTicker(m.checkInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := m.checkPositions(ctx); err != nil {
                log.Printf("Error checking positions: %v", err)
            }
        case <-m.stop:
            return
        case <-ctx.Done():
            return
        }
    }
}

// checkPositions checks all open positions against current market prices
func (m *Monitor) checkPositions(ctx context.Context) error {
    positions, err := m.manager.GetOpenPositions(ctx)
    if err != nil {
        return fmt.Errorf("error getting open positions: %w", err)
    }
    
    for _, position := range positions {
        // Get current price
        ticker, err := m.exchangeClient.GetTicker(ctx, position.Symbol)
        if err != nil {
            log.Printf("Error getting ticker for %s: %v", position.Symbol, err)
            continue
        }
        
        currentPrice := ticker.Price
        
        // Check stop loss
        if position.StopLoss > 0 && currentPrice <= position.StopLoss {
            log.Printf("Stop loss triggered for position %d (%s) at price %.8f",
                position.ID, position.Symbol, currentPrice)
            
            if err := m.manager.ClosePosition(ctx, position.ID, currentPrice); err != nil {
                log.Printf("Error closing position at stop loss: %v", err)
            }
            continue
        }
        
        // Check take profit
        if position.TakeProfit > 0 && currentPrice >= position.TakeProfit {
            log.Printf("Take profit triggered for position %d (%s) at price %.8f",
                position.ID, position.Symbol, currentPrice)
            
            if err := m.manager.ClosePosition(ctx, position.ID, currentPrice); err != nil {
                log.Printf("Error closing position at take profit: %v", err)
            }
            continue
        }
        
        // Update trailing stop if enabled
        if position.TrailingStop > 0 {
            m.updateTrailingStop(ctx, position, currentPrice)
        }
    }
    
    return nil
}

// updateTrailingStop updates the stop loss based on trailing stop settings
func (m *Monitor) updateTrailingStop(ctx context.Context, position *models.Position, currentPrice float64) {
    // Calculate what the stop loss would be at the current price
    trailingDistance := position.TrailingStop
    potentialStopLoss := currentPrice * (1 - trailingDistance)
    
    // Only update if the new stop loss would be higher than the current one
    if potentialStopLoss > position.StopLoss {
        if err := m.manager.UpdateStopLoss(ctx, position.ID, potentialStopLoss); err != nil {
            log.Printf("Error updating trailing stop: %v", err)
        } else {
            log.Printf("Updated trailing stop for position %d to %.8f", 
                position.ID, potentialStopLoss)
        }
    }
}
```

## Risk Management

The Position Management system implements several risk management strategies:

### Position Size Limits

```go
// internal/domain/core/position/risk.go
package position

import (
    "context"
    "errors"
)

// RiskManager handles risk-related aspects of position management
type RiskManager struct {
    manager          *Manager
    maxPositionSize  float64 // Maximum position size as percentage of portfolio
    maxOpenPositions int     // Maximum number of concurrent open positions
    maxLossPerTrade  float64 // Maximum loss per trade as percentage
}

// NewRiskManager creates a new risk manager
func NewRiskManager(
    manager *Manager,
    maxPositionSize float64,
    maxOpenPositions int,
    maxLossPerTrade float64,
) *RiskManager {
    return &RiskManager{
        manager:          manager,
        maxPositionSize:  maxPositionSize,
        maxOpenPositions: maxOpenPositions,
        maxLossPerTrade:  maxLossPerTrade,
    }
}

// CalculatePositionSize calculates the appropriate position size
func (r *RiskManager) CalculatePositionSize(
    ctx context.Context,
    symbol string,
    price float64,
    accountBalance float64,
) (float64, error) {
    // Calculate maximum position value based on account balance
    maxPositionValue := accountBalance * (r.maxPositionSize / 100.0)
    
    // Calculate position size based on current price
    positionSize := maxPositionValue / price
    
    return positionSize, nil
}

// CheckPositionCreationAllowed verifies if creating a new position is allowed
func (r *RiskManager) CheckPositionCreationAllowed(
    ctx context.Context,
    symbol string,
    price float64,
    accountBalance float64,
) error {
    // Check max open positions
    openPositions, err := r.manager.GetOpenPositions(ctx)
    if err != nil {
        return fmt.Errorf("error getting open positions: %w", err)
    }
    
    if len(openPositions) >= r.maxOpenPositions {
        return errors.New("maximum number of open positions reached")
    }
    
    // Check existing positions for this symbol
    symbolPositions, err := r.manager.GetPositionsForSymbol(ctx, symbol)
    if err != nil {
        return fmt.Errorf("error getting symbol positions: %w", err)
    }
    
    for _, pos := range symbolPositions {
        if pos.Status == models.PositionStatusOpen {
            return errors.New("position already exists for this symbol")
        }
    }
    
    return nil
}

// CalculateStopLossLevel determines appropriate stop loss for a position
func (r *RiskManager) CalculateStopLossLevel(
    entryPrice float64,
    riskPercent float64,
) float64 {
    if riskPercent <= 0 {
        riskPercent = r.maxLossPerTrade
    }
    
    stopLoss := entryPrice * (1 - (riskPercent / 100.0))
    return stopLoss
}

// CalculateTakeProfitLevel determines appropriate take profit for a position
func (r *RiskManager) CalculateTakeProfitLevel(
    entryPrice float64,
    stopLoss float64,
    riskRewardRatio float64,
) float64 {
    if riskRewardRatio <= 0 {
        riskRewardRatio = 2.0 // Default 1:2 risk-reward ratio
    }
    
    riskAmount := entryPrice - stopLoss
    rewardAmount := riskAmount * riskRewardRatio
    
    takeProfit := entryPrice + rewardAmount
    return takeProfit
}
```

## Integration with Other Components

The Position Management system integrates with several other components:

1. **Trade Executor**: Creates positions when trades are executed
2. **Exchange Client**: Provides real-time price data
3. **Account Manager**: Supplies portfolio value for position sizing
4. **Repository Layer**: Persists position data

### Factory Creation Example

```go
// internal/domain/factory/factory.go
package factory

// CreatePositionManager creates and configures a position manager
func (f *ServiceFactory) CreatePositionManager() *position.Manager {
    return position.NewManager(
        f.repositories.PositionRepository,
        f.repositories.BoughtCoinRepository,
    )
}

// CreatePositionMonitor creates and configures a position monitor
func (f *ServiceFactory) CreatePositionMonitor() *position.Monitor {
    manager := f.CreatePositionManager()
    
    return position.NewMonitor(
        manager,
        f.exchangeClient,
        f.config.PositionCheckInterval,
    )
}

// CreateRiskManager creates and configures a risk manager
func (f *ServiceFactory) CreateRiskManager() *position.RiskManager {
    manager := f.CreatePositionManager()
    
    return position.NewRiskManager(
        manager,
        f.config.MaxPositionSize,
        f.config.MaxOpenPositions,
        f.config.MaxLossPerTrade,
    )
}
```

## Testing Strategy

### Unit Testing

```go
// internal/domain/core/position/manager_test.go
package position_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/ryanlisse/cryptobot/internal/domain/core/position"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository/mocks"
)

func TestManager_CreatePosition(t *testing.T) {
    // Create mocks
    mockPositionRepo := new(mocks.PositionRepository)
    mockCoinRepo := new(mocks.BoughtCoinRepository)
    
    // Test position
    testPosition := &models.Position{
        Symbol:     "BTCUSDT",
        EntryPrice: 50000.0,
        Quantity:   0.1,
        StopLoss:   48000.0,
        TakeProfit: 55000.0,
        CoinID:     1,
    }
    
    expectedPosition := &models.Position{
        ID:         1,
        Symbol:     "BTCUSDT",
        EntryPrice: 50000.0,
        Quantity:   0.1,
        OpenedAt:   time.Now(),
        Status:     models.PositionStatusOpen,
        StopLoss:   48000.0,
        TakeProfit: 55000.0,
        CoinID:     1,
    }
    
    // Setup expectations
    mockPositionRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *models.Position) bool {
        return p.Symbol == testPosition.Symbol &&
               p.EntryPrice == testPosition.EntryPrice &&
               p.Quantity == testPosition.Quantity &&
               p.Status == models.PositionStatusOpen
    })).Return(expectedPosition, nil)
    
    // Create manager
    manager := position.NewManager(mockPositionRepo, mockCoinRepo)
    
    // Test creating a position
    result, err := manager.CreatePosition(context.Background(), testPosition)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedPosition.ID, result.ID)
    assert.Equal(t, expectedPosition.Symbol, result.Symbol)
    assert.Equal(t, expectedPosition.Status, result.Status)
    
    // Verify mocks
    mockPositionRepo.AssertExpectations(t)
}

func TestManager_ClosePosition(t *testing.T) {
    // Create mocks
    mockPositionRepo := new(mocks.PositionRepository)
    mockCoinRepo := new(mocks.BoughtCoinRepository)
    
    // Test position
    positionID := int64(1)
    exitPrice := 55000.0
    
    openPosition := &models.Position{
        ID:         positionID,
        Symbol:     "BTCUSDT",
        EntryPrice: 50000.0,
        Quantity:   0.1,
        OpenedAt:   time.Now().Add(-24 * time.Hour),
        Status:     models.PositionStatusOpen,
        StopLoss:   48000.0,
        TakeProfit: 55000.0,
    }
    
    // Setup expectations
    mockPositionRepo.On("FindByID", mock.Anything, positionID).
        Return(openPosition, nil)
    
    mockPositionRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *models.Position) bool {
        return p.ID == positionID &&
               p.Status == models.PositionStatusClosed &&
               p.ExitPrice == exitPrice &&
               p.ProfitLoss > 0
    })).Return(openPosition, nil)
    
    // Create manager
    manager := position.NewManager(mockPositionRepo, mockCoinRepo)
    
    // Test closing a position
    err := manager.ClosePosition(context.Background(), positionID, exitPrice)
    
    // Assertions
    assert.NoError(t, err)
    
    // Verify mocks
    mockPositionRepo.AssertExpectations(t)
}
```

For more comprehensive testing examples and best practices, refer to the [Testing Strategy](../testing/overview.md) document.
