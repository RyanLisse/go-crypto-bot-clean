# Core Business Logic Implementation

This document provides an overview of the core business logic components for the Go crypto trading bot. The implementation follows hexagonal architecture principles, with a clear separation between domain logic and external dependencies.

## Table of Contents

1. [Overview](#overview)
2. [Component Structure](#component-structure)
3. [Service Interfaces](#service-interfaces)
4. [Key Components](#key-components)
   - [New Coin Watcher](#new-coin-watcher)
   - [Trade Executor](#trade-executor)
   - [Account Manager](#account-manager)
5. [Related Documents](#related-documents)

## Overview

The core business logic of the trading bot consists of several key components:

1. **New Coin Watcher**: Monitors exchange for newly listed coins and triggers evaluation
2. **Trade Executor**: Handles purchase and sale decisions and execution
3. **Account Manager**: Manages wallet balances and account state
4. **Trading Strategies**: Implements various trading strategies (basic and advanced)
5. **Position Management**: Manages trade positions with stop-loss and take-profit
6. **Risk Controls**: Implements risk management and capital protection

These components work together to create a complete trading system that can detect new coins, make trading decisions, execute trades, and manage positions.

## Component Structure

The core business logic is organized in the following directory structure:

```
internal/domain/
├── models/              # Domain entities and value objects
│   ├── coin.go          # Coin and related types
│   ├── order.go         # Order models
│   ├── position.go      # Position models
│   └── wallet.go        # Wallet and balance models
│
├── service/             # Service interfaces
│   ├── coin_service.go  # Coin-related operations
│   ├── market_service.go # Market data operations
│   ├── trade_service.go # Trading operations
│   └── wallet_service.go # Wallet operations
│
├── core/                # Core business logic implementations
│   ├── newcoin/         # New coin detection
│   │   ├── watcher.go   # Coin monitoring
│   │   └── filter.go    # Filtering criteria
│   │
│   ├── trading/         # Trading logic
│   │   ├── executor.go  # Trade execution
│   │   ├── strategy.go  # Basic strategy interface
│   │   └── basic.go     # Basic trading strategy
│   │
│   ├── account/         # Account management
│   │   ├── manager.go   # Account operations
│   │   └── balance.go   # Balance tracking
│   │
│   ├── position/        # Position management
│   │   ├── manager.go   # Position lifecycle
│   │   └── adjustor.go  # Position adjustment
│   │
│   └── risk/            # Risk management
│       ├── calculator.go # Risk calculation
│       └── limits.go    # Risk limits
│
└── factory/             # Factory for creating service instances
    └── factory.go       # Service factory implementation
```

## Service Interfaces

The core business logic is defined through service interfaces that establish clear boundaries between components:

```go
// internal/domain/service/coin_service.go
package service

import (
    "context"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// CoinService defines operations for coin management
type CoinService interface {
    // GetBoughtCoins retrieves all purchased coins
    GetBoughtCoins(ctx context.Context, includeDeleted bool) ([]*models.BoughtCoin, error)
    
    // GetCoinByID retrieves a coin by its ID
    GetCoinByID(ctx context.Context, id int64) (*models.BoughtCoin, error)
    
    // GetCoinBySymbol retrieves a coin by its symbol
    GetCoinBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error)
    
    // SaveCoin saves a coin to the repository
    SaveCoin(ctx context.Context, coin *models.BoughtCoin) (*models.BoughtCoin, error)
    
    // DeleteCoin marks a coin as deleted
    DeleteCoin(ctx context.Context, id int64) error
}

// MarketService defines operations for market data
type MarketService interface {
    // GetTicker retrieves current price information for a symbol
    GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)
    
    // GetKlines retrieves historical candle data
    GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*models.Kline, error)
    
    // GetAllSymbols retrieves all trading symbols
    GetAllSymbols(ctx context.Context) ([]string, error)
}

// TradeService defines operations for trading
type TradeService interface {
    // ExecutePurchase performs a purchase operation
    ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error)
    
    // ExecuteSale performs a sale operation
    ExecuteSale(ctx context.Context, coinID int64) (*models.SoldCoin, error)
    
    // GetProfitLoss calculates profit/loss for a coin
    GetProfitLoss(ctx context.Context, coin *models.BoughtCoin) (float64, float64, error)
}

// WalletService defines operations for wallet management
type WalletService interface {
    // GetBalance retrieves the current wallet balance
    GetBalance(ctx context.Context) (float64, error)
    
    // UpdateBalance updates the wallet balance
    UpdateBalance(ctx context.Context, amount float64, reason string) error
    
    // GetTransactionHistory retrieves transaction history
    GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error)
}
```

## Key Components

### New Coin Watcher

The New Coin Watcher monitors the exchange for newly listed coins:

```go
// NewWatcher creates a new coin watcher
func NewWatcher(
    mexcClient exchange.Client,
    newCoinRepo repository.NewCoinRepository,
    decisionRepo repository.PurchaseDecisionRepository,
    tradeExecutor trading.Executor,
    checkInterval time.Duration,
) *Watcher {
    return &Watcher{
        mexcClient:          mexcClient,
        newCoinRepo:         newCoinRepo,
        purchaseDecisionRepo: decisionRepo,
        tradeExecutor:       tradeExecutor,
        checkInterval:       checkInterval,
        knownSymbols:        make(map[string]bool),
        stop:                make(chan struct{}),
    }
}

// Start begins watching for new coins
func (w *Watcher) Start(ctx context.Context) error {
    w.mutex.Lock()
    if w.running {
        w.mutex.Unlock()
        return errors.New("watcher is already running")
    }
    w.running = true
    w.mutex.Unlock()

    // Load known symbols
    if err := w.loadKnownSymbols(ctx); err != nil {
        return err
    }

    // Start watching
    go w.watchLoop(ctx)
    
    return nil
}
```

For more details on the New Coin Watcher implementation, see [06a-new-coin-watcher.md](06a-new-coin-watcher.md).

### Trade Executor

The Trade Executor handles the execution of trades based on strategy decisions:

```go
// NewExecutor creates a new trade executor
func NewExecutor(
    mexcClient exchange.Client,
    boughtCoinRepo repository.BoughtCoinRepository,
    soldCoinRepo repository.SoldCoinRepository,
    decisionRepo repository.PurchaseDecisionRepository,
    strategy trading.Strategy,
    accountManager account.Manager,
) *Executor {
    return &Executor{
        mexcClient:     mexcClient,
        boughtCoinRepo: boughtCoinRepo,
        soldCoinRepo:   soldCoinRepo,
        decisionRepo:   decisionRepo,
        strategy:       strategy,
        accountManager: accountManager,
    }
}

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
    balance, err := e.accountManager.GetBalance(ctx)
    if err != nil {
        return nil, fmt.Errorf("error getting balance: %w", err)
    }
    
    if balance < amount {
        return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", balance, amount)
    }
    
    // Get current price
    ticker, err := e.mexcClient.GetTicker(ctx, symbol)
    if err != nil {
        return nil, fmt.Errorf("error getting ticker: %w", err)
    }
    
    // Calculate quantity
    quantity := amount / ticker.Price
    
    // Execute the purchase
    // In a real implementation, this would call the exchange API
    orderID, err := e.mexcClient.PlaceBuyOrder(ctx, symbol, quantity)
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
    
    // Update account balance
    if err := e.accountManager.UpdateBalance(ctx, -amount, fmt.Sprintf("Purchase of %s", symbol)); err != nil {
        return nil, fmt.Errorf("error updating balance: %w", err)
    }
    
    return savedCoin, nil
}
```

For more details on the Trade Executor implementation, see [06b-trade-executor.md](06b-trade-executor.md).

### Account Manager

The Account Manager handles wallet balances and account state:

```go
// NewManager creates a new account manager
func NewManager(
    walletRepo repository.WalletRepository,
    transactionRepo repository.TransactionRepository,
    mexcClient exchange.Client,
) *Manager {
    return &Manager{
        walletRepo:      walletRepo,
        transactionRepo: transactionRepo,
        mexcClient:      mexcClient,
        cacheTTL:        cacheTTL,
        balanceCache:    nil,
        balanceCacheExp: time.Time{},
    }
}

// GetBalance retrieves the current wallet balance
func (m *Manager) GetBalance(ctx context.Context) (float64, error) {
    // Check cache first
    m.mutex.RLock()
    if m.balanceCache != nil && time.Now().Before(m.balanceCacheExp) {
        balance := *m.balanceCache
        m.mutex.RUnlock()
        return balance, nil
    }
    m.mutex.RUnlock()
    
    // Get balance from repository
    wallet, err := m.walletRepo.GetWallet(ctx)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            // Initialize wallet if not exists
            wallet = &models.Wallet{
                Balance:   0,
                UpdatedAt: time.Now(),
            }
            
            if _, err := m.walletRepo.SaveWallet(ctx, wallet); err != nil {
                return 0, fmt.Errorf("error initializing wallet: %w", err)
            }
        } else {
            return 0, fmt.Errorf("error getting wallet: %w", err)
        }
    }
    
    // Update cache
    m.mutex.Lock()
    m.balanceCache = &wallet.Balance
    m.balanceCacheExp = time.Now().Add(m.cacheTTL)
    m.mutex.Unlock()
    
    return wallet.Balance, nil
}
```

For more details on the Account Manager implementation, see [06c-account-manager.md](06c-account-manager.md).

## Related Documents

For more detailed information on specific components, refer to the following documents:

- [06a-new-coin-watcher.md](06a-new-coin-watcher.md) - Detailed implementation of the new coin detection system
- [06b-trade-executor.md](06b-trade-executor.md) - Implementation of trade execution logic
- [06c-account-manager.md](06c-account-manager.md) - Account and wallet management
- [08-position-management.md](08-position-management.md) - Position lifecycle management
- [09a-advanced-trading-strategies.md](09a-advanced-trading-strategies.md) - Advanced trading strategies
- [09b-risk-management.md](09b-risk-management.md) - Risk management and capital protection
