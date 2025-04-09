# Account Manager Implementation

This document provides implementation details for the Account Manager component of the Go crypto trading bot. The Account Manager is responsible for handling wallet balances, transaction history, and account state management.

## Table of Contents

1. [Overview](#overview)
2. [Component Structure](#component-structure)
3. [Implementation Details](#implementation-details)
4. [Transaction Management](#transaction-management)
5. [Integration with Other Components](#integration-with-other-components)
6. [Testing Guidance](#testing-guidance)

## Overview

The Account Manager serves as the central component for managing the financial aspects of the trading bot, including:

- Tracking wallet balance
- Recording transactions
- Managing deposits and withdrawals
- Providing balance reports
- Ensuring data consistency between local state and exchange

This component implements caching mechanisms to reduce unnecessary database and API calls while maintaining accurate balance information.

## Component Structure

The Account Manager is implemented in the following file structure:

```
internal/domain/core/account/
├── manager.go        # Main account manager implementation
├── balance.go        # Balance tracking logic
└── transaction.go    # Transaction recording and management
```

## Implementation Details

### Manager Definition

```go
// internal/domain/core/account/manager.go
package account

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/platform/exchange"
)

// Default cache time-to-live
const defaultCacheTTL = 5 * time.Minute

// Manager handles account and wallet operations
type Manager struct {
    walletRepo         repository.WalletRepository
    transactionRepo    repository.TransactionRepository
    exchangeClient     exchange.Client
    
    // Cache mechanism
    cacheTTL           time.Duration
    mutex              sync.RWMutex
    balanceCache       *float64
    balanceCacheExp    time.Time
}

// NewManager creates a new account manager
func NewManager(
    walletRepo repository.WalletRepository,
    transactionRepo repository.TransactionRepository,
    exchangeClient exchange.Client,
    cacheTTL time.Duration,
) *Manager {
    if cacheTTL <= 0 {
        cacheTTL = defaultCacheTTL
    }
    
    return &Manager{
        walletRepo:      walletRepo,
        transactionRepo: transactionRepo,
        exchangeClient:  exchangeClient,
        cacheTTL:        cacheTTL,
        balanceCache:    nil,
        balanceCacheExp: time.Time{},
    }
}
```

### Balance Management

```go
// GetBalance retrieves the current wallet balance
func (m *Manager) GetBalance(ctx context.Context) (float64, error) {
    // Check cache first for better performance
    m.mutex.RLock()
    if m.balanceCache != nil && time.Now().Before(m.balanceCacheExp) {
        balance := *m.balanceCache
        m.mutex.RUnlock()
        return balance, nil
    }
    m.mutex.RUnlock()
    
    // Get balance from repository if cache is invalid
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

// UpdateBalance updates the wallet balance
func (m *Manager) UpdateBalance(ctx context.Context, amount float64, reason string) error {
    // Get current balance
    wallet, err := m.walletRepo.GetWallet(ctx)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            // Initialize wallet with the provided amount if not exists
            wallet = &models.Wallet{
                Balance:   amount,
                UpdatedAt: time.Now(),
            }
        } else {
            return fmt.Errorf("error getting wallet: %w", err)
        }
    } else {
        // Update existing wallet
        wallet.Balance += amount
        wallet.UpdatedAt = time.Now()
    }
    
    // Save updated wallet
    updatedWallet, err := m.walletRepo.SaveWallet(ctx, wallet)
    if err != nil {
        return fmt.Errorf("error saving wallet: %w", err)
    }
    
    // Record transaction
    transaction := &models.Transaction{
        Amount:    amount,
        Balance:   updatedWallet.Balance,
        Reason:    reason,
        Timestamp: time.Now(),
    }
    
    if _, err := m.transactionRepo.Create(ctx, transaction); err != nil {
        return fmt.Errorf("error recording transaction: %w", err)
    }
    
    // Update cache
    m.mutex.Lock()
    m.balanceCache = &updatedWallet.Balance
    m.balanceCacheExp = time.Now().Add(m.cacheTTL)
    m.mutex.Unlock()
    
    return nil
}

// SyncWithExchange synchronizes the local wallet balance with the exchange
func (m *Manager) SyncWithExchange(ctx context.Context) error {
    // Get balance from exchange
    exchangeBalance, err := m.exchangeClient.GetAccountBalance(ctx)
    if err != nil {
        return fmt.Errorf("error getting exchange balance: %w", err)
    }
    
    // Get local wallet
    wallet, err := m.walletRepo.GetWallet(ctx)
    if err != nil && !errors.Is(err, repository.ErrNotFound) {
        return fmt.Errorf("error getting local wallet: %w", err)
    }
    
    var currentBalance float64
    if wallet != nil {
        currentBalance = wallet.Balance
    }
    
    // If there's a discrepancy, update local wallet and record transaction
    if currentBalance != exchangeBalance {
        difference := exchangeBalance - currentBalance
        reason := "Balance sync with exchange"
        
        if err := m.UpdateBalance(ctx, difference, reason); err != nil {
            return fmt.Errorf("error updating balance: %w", err)
        }
    }
    
    return nil
}
```

### Balance Summary

```go
// GetBalanceSummary generates a summary of the wallet and transactions
func (m *Manager) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
    // Get current balance
    balance, err := m.GetBalance(ctx)
    if err != nil {
        return nil, fmt.Errorf("error getting balance: %w", err)
    }
    
    // Calculate time period
    endTime := time.Now()
    startTime := endTime.AddDate(0, 0, -days)
    
    // Get transactions for period
    transactions, err := m.transactionRepo.FindByTimeRange(ctx, startTime, endTime)
    if err != nil {
        return nil, fmt.Errorf("error getting transactions: %w", err)
    }
    
    // Calculate metrics
    var deposits, withdrawals float64
    
    for _, tx := range transactions {
        if tx.Amount > 0 {
            deposits += tx.Amount
        } else {
            withdrawals += -tx.Amount
        }
    }
    
    // Construct summary
    summary := &models.BalanceSummary{
        CurrentBalance: balance,
        Deposits:       deposits,
        Withdrawals:    withdrawals,
        NetChange:      deposits - withdrawals,
        TransactionCount: len(transactions),
        Period:         days,
        GeneratedAt:    time.Now(),
    }
    
    return summary, nil
}
```

## Transaction Management

The transaction functionality is implemented to record all financial activities:

```go
// internal/domain/core/account/transaction.go
package account

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
)

// TransactionManager handles transaction recording and querying
type TransactionManager struct {
    transactionRepo repository.TransactionRepository
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(
    transactionRepo repository.TransactionRepository,
) *TransactionManager {
    return &TransactionManager{
        transactionRepo: transactionRepo,
    }
}

// RecordTransaction creates a new transaction record
func (m *TransactionManager) RecordTransaction(
    ctx context.Context, 
    amount float64, 
    balance float64, 
    reason string,
) (*models.Transaction, error) {
    if reason == "" {
        reason = "Unspecified"
    }
    
    transaction := &models.Transaction{
        Amount:    amount,
        Balance:   balance,
        Reason:    reason,
        Timestamp: time.Now(),
    }
    
    return m.transactionRepo.Create(ctx, transaction)
}

// GetTransactionHistory retrieves transaction history for a specified period
func (m *TransactionManager) GetTransactionHistory(
    ctx context.Context, 
    startTime, 
    endTime time.Time,
) ([]*models.Transaction, error) {
    if endTime.IsZero() {
        endTime = time.Now()
    }
    
    if startTime.IsZero() || startTime.After(endTime) {
        return nil, errors.New("invalid time range")
    }
    
    return m.transactionRepo.FindByTimeRange(ctx, startTime, endTime)
}

// AnalyzeTransactions performs analysis on transaction data
func (m *TransactionManager) AnalyzeTransactions(
    ctx context.Context, 
    startTime, 
    endTime time.Time,
) (*models.TransactionAnalysis, error) {
    transactions, err := m.GetTransactionHistory(ctx, startTime, endTime)
    if err != nil {
        return nil, fmt.Errorf("error getting transactions: %w", err)
    }
    
    if len(transactions) == 0 {
        return &models.TransactionAnalysis{
            StartTime:   startTime,
            EndTime:     endTime,
            TotalCount:  0,
            TotalVolume: 0,
        }, nil
    }
    
    var buys, sells int
    var buyVolume, sellVolume float64
    
    for _, tx := range transactions {
        if isBuyTransaction(tx.Reason) {
            buys++
            buyVolume += tx.Amount
        } else if isSellTransaction(tx.Reason) {
            sells++
            sellVolume += -tx.Amount
        }
    }
    
    // Create analysis result
    analysis := &models.TransactionAnalysis{
        StartTime:   startTime,
        EndTime:     endTime,
        TotalCount:  len(transactions),
        BuyCount:    buys,
        SellCount:   sells,
        TotalVolume: buyVolume + sellVolume,
        BuyVolume:   buyVolume,
        SellVolume:  sellVolume,
    }
    
    return analysis, nil
}

// Helper to determine transaction type from reason
func isBuyTransaction(reason string) bool {
    return strings.Contains(strings.ToLower(reason), "purchase") ||
           strings.Contains(strings.ToLower(reason), "buy") ||
           strings.Contains(strings.ToLower(reason), "deposit")
}

func isSellTransaction(reason string) bool {
    return strings.Contains(strings.ToLower(reason), "sale") ||
           strings.Contains(strings.ToLower(reason), "sell") ||
           strings.Contains(strings.ToLower(reason), "withdrawal")
}
```

## Integration with Other Components

The Account Manager integrates with several other components:

1. **Exchange Client**: For retrieving account balance information from the exchange
2. **Repository Layer**: For persisting wallet and transaction data
3. **Trade Executor**: For updating the wallet when trades are executed

### Model Definitions

These are the relevant model definitions:

```go
// internal/domain/models/wallet.go
package models

import (
    "time"
)

// Wallet represents the user's wallet
type Wallet struct {
    ID        int64     `json:"id" db:"id"`
    Balance   float64   `json:"balance" db:"balance"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction represents a financial transaction
type Transaction struct {
    ID        int64     `json:"id" db:"id"`
    Amount    float64   `json:"amount" db:"amount"`
    Balance   float64   `json:"balance" db:"balance"`
    Reason    string    `json:"reason" db:"reason"`
    Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// BalanceSummary provides an overview of wallet activity
type BalanceSummary struct {
    CurrentBalance   float64   `json:"current_balance"`
    Deposits         float64   `json:"deposits"`
    Withdrawals      float64   `json:"withdrawals"`
    NetChange        float64   `json:"net_change"`
    TransactionCount int       `json:"transaction_count"`
    Period           int       `json:"period_days"`
    GeneratedAt      time.Time `json:"generated_at"`
}

// TransactionAnalysis provides analysis of transaction history
type TransactionAnalysis struct {
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    TotalCount  int       `json:"total_count"`
    BuyCount    int       `json:"buy_count"`
    SellCount   int       `json:"sell_count"`
    TotalVolume float64   `json:"total_volume"`
    BuyVolume   float64   `json:"buy_volume"`
    SellVolume  float64   `json:"sell_volume"`
}
```

### Factory Creation Example

```go
// internal/domain/factory/factory.go
package factory

// CreateAccountManager creates and configures an account manager
func (f *ServiceFactory) CreateAccountManager() *account.Manager {
    return account.NewManager(
        f.repositories.WalletRepository,
        f.repositories.TransactionRepository,
        f.exchangeClient,
        f.config.CacheTTL,
    )
}
```

## Testing Guidance

### Unit Testing

```go
// internal/domain/core/account/manager_test.go
package account_test

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/ryanlisse/cryptobot/internal/domain/core/account"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/domain/repository/mocks"
    "github.com/ryanlisse/cryptobot/internal/platform/exchange/mocks"
)

func TestManager_GetBalance(t *testing.T) {
    // Create mocks
    mockWalletRepo := new(mocks.WalletRepository)
    mockTxRepo := new(mocks.TransactionRepository)
    mockExchange := new(mocks.ExchangeClient)
    
    // Test wallet
    wallet := &models.Wallet{
        ID:        1,
        Balance:   100.0,
        UpdatedAt: time.Now(),
    }
    
    // Setup expectations
    mockWalletRepo.On("GetWallet", mock.Anything).Return(wallet, nil)
    
    // Create manager with short cache TTL for testing
    manager := account.NewManager(
        mockWalletRepo,
        mockTxRepo,
        mockExchange,
        100*time.Millisecond,
    )
    
    // Test getting balance
    balance, err := manager.GetBalance(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, wallet.Balance, balance)
    
    // Test that repository is only called once due to caching
    balance, err = manager.GetBalance(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, wallet.Balance, balance)
    
    // Verify mock was called exactly once
    mockWalletRepo.AssertNumberOfCalls(t, "GetWallet", 1)
    
    // Test cache expiration
    time.Sleep(200 * time.Millisecond)
    
    balance, err = manager.GetBalance(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, wallet.Balance, balance)
    
    // Verify repository was called again after cache expired
    mockWalletRepo.AssertNumberOfCalls(t, "GetWallet", 2)
}

func TestManager_UpdateBalance(t *testing.T) {
    // Create mocks
    mockWalletRepo := new(mocks.WalletRepository)
    mockTxRepo := new(mocks.TransactionRepository)
    mockExchange := new(mocks.ExchangeClient)
    
    // Initial wallet
    wallet := &models.Wallet{
        ID:        1,
        Balance:   100.0,
        UpdatedAt: time.Now(),
    }
    
    updatedWallet := &models.Wallet{
        ID:        1,
        Balance:   150.0,
        UpdatedAt: time.Now(),
    }
    
    // Setup expectations
    mockWalletRepo.On("GetWallet", mock.Anything).Return(wallet, nil)
    mockWalletRepo.On("SaveWallet", mock.Anything, mock.MatchedBy(func(w *models.Wallet) bool {
        return w.Balance == 150.0
    })).Return(updatedWallet, nil)
    
    mockTxRepo.On("Create", mock.Anything, mock.MatchedBy(func(tx *models.Transaction) bool {
        return tx.Amount == 50.0 && tx.Balance == 150.0
    })).Return(&models.Transaction{ID: 1}, nil)
    
    // Create manager
    manager := account.NewManager(
        mockWalletRepo,
        mockTxRepo,
        mockExchange,
        time.Minute,
    )
    
    // Test updating balance
    err := manager.UpdateBalance(context.Background(), 50.0, "Test deposit")
    assert.NoError(t, err)
    
    // Verify mocks
    mockWalletRepo.AssertExpectations(t)
    mockTxRepo.AssertExpectations(t)
    
    // Test that cached balance is updated
    balance, err := manager.GetBalance(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, 150.0, balance)
    
    // Repository should not be called again since we have a valid cache
    mockWalletRepo.AssertNumberOfCalls(t, "GetWallet", 1)
}
```

### Integration Testing

When performing integration tests on the Account Manager, consider these scenarios:

1. Creating an initial wallet when none exists
2. Processing multiple transactions and verifying balance consistency
3. Syncing with exchange balances
4. Testing cache invalidation during concurrent operations

For a comprehensive integration test, use a real SQLite database in a test fixture:

```go
func TestAccountManager_Integration(t *testing.T) {
    // Setup test database
    db, cleanup := setupTestDatabase(t)
    defer cleanup()
    
    // Create repos with real database
    walletRepo := sqlite.NewWalletRepository(db)
    txRepo := sqlite.NewTransactionRepository(db)
    
    // Create mock exchange
    mockExchange := new(mocks.ExchangeClient)
    mockExchange.On("GetAccountBalance", mock.Anything).Return(100.0, nil)
    
    // Create manager
    manager := account.NewManager(
        walletRepo,
        txRepo,
        mockExchange,
        time.Minute,
    )
    
    ctx := context.Background()
    
    // First sync should create wallet and set initial balance
    err := manager.SyncWithExchange(ctx)
    assert.NoError(t, err)
    
    // Check balance
    balance, err := manager.GetBalance(ctx)
    assert.NoError(t, err)
    assert.Equal(t, 100.0, balance)
    
    // Record some transactions
    err = manager.UpdateBalance(ctx, -20.0, "Test purchase")
    assert.NoError(t, err)
    
    balance, err = manager.GetBalance(ctx)
    assert.NoError(t, err)
    assert.Equal(t, 80.0, balance)
    
    err = manager.UpdateBalance(ctx, 30.0, "Test sale")
    assert.NoError(t, err)
    
    balance, err = manager.GetBalance(ctx)
    assert.NoError(t, err)
    assert.Equal(t, 110.0, balance)
    
    // Check transaction history
    txManager := account.NewTransactionManager(txRepo)
    
    history, err := txManager.GetTransactionHistory(ctx, time.Time{}, time.Now())
    assert.NoError(t, err)
    assert.Len(t, history, 3) // Initial sync + 2 updates
    
    // Get balance summary
    summary, err := manager.GetBalanceSummary(ctx, 1)
    assert.NoError(t, err)
    assert.Equal(t, 110.0, summary.CurrentBalance)
    assert.Equal(t, 130.0, summary.Deposits)
    assert.Equal(t, 20.0, summary.Withdrawals)
    assert.Equal(t, 110.0, summary.NetChange)
    assert.Equal(t, 3, summary.TransactionCount)
}
```

For more details on testing, see the general [Testing Strategy](../testing/overview.md) document.
