# New Coin Watcher Implementation

This document provides detailed implementation guidance for the New Coin Watcher component of the Go crypto trading bot. The New Coin Watcher is responsible for monitoring the MEXC exchange for newly listed coins and triggering evaluation and potential purchases.

## Table of Contents

1. [Overview](#overview)
2. [Component Structure](#component-structure)
3. [Implementation Details](#implementation-details)
4. [Integration with Other Components](#integration-with-other-components)
5. [Configuration Options](#configuration-options)
6. [Testing Guidance](#testing-guidance)

## Overview

The New Coin Watcher is a critical component that continuously monitors the MEXC exchange for newly listed coins. When a new coin is detected, the watcher records it and initiates an evaluation process to determine if the coin meets the criteria for purchase. This component works in tandem with trading strategies and the trade executor to automate the detection and purchase of promising new listings.

## Component Structure

The New Coin Watcher is implemented in the following file structure:

```
internal/domain/core/newcoin/
├── watcher.go       # Main watcher implementation
├── filter.go        # Filtering criteria for new coins
└── detector.go      # Detection logic for new coins
```

## Implementation Details

### Watcher Definition

```go
// internal/domain/core/newcoin/watcher.go
package newcoin

import (
    "context"
    "errors"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/domain/service"
    "github.com/ryanlisse/cryptobot/internal/platform/exchange"
)

// Watcher monitors MEXC for new coins
type Watcher struct {
    exchangeClient      exchange.Client
    newCoinRepo         repository.NewCoinRepository
    decisionRepo        repository.PurchaseDecisionRepository
    tradeService        service.TradeService
    checkInterval       time.Duration
    mutex               sync.RWMutex
    knownSymbols        map[string]bool
    running             bool
    stop                chan struct{}
    symbolFilter        SymbolFilter
}

// SymbolFilter defines an interface for filtering symbols
type SymbolFilter interface {
    // ShouldTrack determines if a symbol should be tracked
    ShouldTrack(symbol string) bool
}

// NewWatcher creates a new coin watcher
func NewWatcher(
    exchangeClient exchange.Client,
    newCoinRepo repository.NewCoinRepository,
    decisionRepo repository.PurchaseDecisionRepository,
    tradeService service.TradeService,
    checkInterval time.Duration,
    filter SymbolFilter,
) *Watcher {
    if filter == nil {
        filter = &DefaultSymbolFilter{}
    }
    
    return &Watcher{
        exchangeClient:      exchangeClient,
        newCoinRepo:         newCoinRepo,
        decisionRepo:        decisionRepo,
        tradeService:        tradeService,
        checkInterval:       checkInterval,
        knownSymbols:        make(map[string]bool),
        stop:                make(chan struct{}),
        symbolFilter:        filter,
    }
}
```

### Starting and Stopping the Watcher

```go
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
        w.running = false
        return err
    }

    // Start watching
    go w.watchLoop(ctx)
    
    return nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    if !w.running {
        return
    }
    
    close(w.stop)
    w.running = false
}

// IsRunning checks if the watcher is running
func (w *Watcher) IsRunning() bool {
    w.mutex.RLock()
    defer w.mutex.RUnlock()
    return w.running
}
```

### Main Watching Loop

```go
// watchLoop is the main loop for watching new coins
func (w *Watcher) watchLoop(ctx context.Context) {
    ticker := time.NewTicker(w.checkInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := w.checkForNewCoins(ctx); err != nil {
                log.Printf("Error checking for new coins: %v", err)
            }
        case <-w.stop:
            return
        case <-ctx.Done():
            return
        }
    }
}

// loadKnownSymbols loads known symbols from the repository
func (w *Watcher) loadKnownSymbols(ctx context.Context) error {
    newCoins, err := w.newCoinRepo.FindAll(ctx)
    if err != nil {
        return fmt.Errorf("error loading known coins: %w", err)
    }
    
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    for _, coin := range newCoins {
        w.knownSymbols[coin.Symbol] = true
    }
    
    return nil
}
```

### Checking for New Coins

```go
// checkForNewCoins fetches symbols from the exchange and compares with known symbols
func (w *Watcher) checkForNewCoins(ctx context.Context) error {
    symbols, err := w.exchangeClient.GetAllSymbols(ctx)
    if err != nil {
        return fmt.Errorf("error getting symbols: %w", err)
    }
    
    for _, symbol := range symbols {
        // Apply filter
        if !w.symbolFilter.ShouldTrack(symbol) {
            continue
        }
        
        // Check if this is a new symbol
        w.mutex.RLock()
        known := w.knownSymbols[symbol]
        w.mutex.RUnlock()
        
        if !known {
            w.handleNewSymbol(ctx, symbol)
        }
    }
    
    return nil
}

// handleNewSymbol processes a newly discovered symbol
func (w *Watcher) handleNewSymbol(ctx context.Context, symbol string) {
    // Record the new coin
    newCoin := &models.NewCoin{
        Symbol:    symbol,
        FoundAt:   time.Now(),
        Processed: false,
    }
    
    savedCoin, err := w.newCoinRepo.Create(ctx, newCoin)
    if err != nil {
        log.Printf("Error saving new coin %s: %v", symbol, err)
        return
    }
    
    // Add to known symbols
    w.mutex.Lock()
    w.knownSymbols[symbol] = true
    w.mutex.Unlock()
    
    // Trigger evaluation
    go w.evaluateNewCoin(context.Background(), savedCoin)
}

// evaluateNewCoin evaluates a new coin for potential purchase
func (w *Watcher) evaluateNewCoin(ctx context.Context, coin *models.NewCoin) {
    log.Printf("Evaluating new coin: %s", coin.Symbol)
    
    // Mark as processed
    coin.Processed = true
    if _, err := w.newCoinRepo.Update(ctx, coin); err != nil {
        log.Printf("Error updating coin status: %v", err)
    }
    
    // Evaluate for purchase
    decision := &models.PurchaseDecision{
        Symbol:       coin.Symbol,
        EvaluatedAt:  time.Now(),
        ShouldBuy:    false,
        Reason:       "Pending evaluation",
        NewCoinID:    coin.ID,
    }
    
    // Use the trade service to decide
    shouldBuy, reason, err := w.tradeService.ShouldBuyNewCoin(ctx, coin.Symbol)
    if err != nil {
        decision.ShouldBuy = false
        decision.Reason = fmt.Sprintf("Evaluation error: %v", err)
    } else {
        decision.ShouldBuy = shouldBuy
        decision.Reason = reason
    }
    
    // Save the decision
    savedDecision, err := w.decisionRepo.Create(ctx, decision)
    if err != nil {
        log.Printf("Error saving purchase decision: %v", err)
        return
    }
    
    // Execute purchase if recommended
    if savedDecision.ShouldBuy {
        w.executePurchase(ctx, coin.Symbol)
    }
}

// executePurchase initiates a purchase for a coin
func (w *Watcher) executePurchase(ctx context.Context, symbol string) {
    log.Printf("Executing purchase for symbol: %s", symbol)
    
    // Get purchase amount from configuration or use a default
    amount := 10.0 // Example fixed amount, in practice this should be configurable
    
    // Execute the purchase
    boughtCoin, err := w.tradeService.ExecutePurchase(ctx, symbol, amount)
    if err != nil {
        log.Printf("Error executing purchase for %s: %v", symbol, err)
        return
    }
    
    log.Printf("Successfully purchased %s: %f units at $%f", 
        boughtCoin.Symbol, boughtCoin.Quantity, boughtCoin.PurchasePrice)
}
```

### Symbol Filter Implementation

```go
// internal/domain/core/newcoin/filter.go
package newcoin

import (
    "strings"
)

// DefaultSymbolFilter implements basic symbol filtering
type DefaultSymbolFilter struct {
    baseAssets      []string
    excludedSymbols []string
}

// NewDefaultSymbolFilter creates a new default symbol filter
func NewDefaultSymbolFilter(baseAssets []string, excludedSymbols []string) *DefaultSymbolFilter {
    if baseAssets == nil {
        // Default to USDT pairs only
        baseAssets = []string{"USDT"}
    }
    
    if excludedSymbols == nil {
        excludedSymbols = []string{}
    }
    
    return &DefaultSymbolFilter{
        baseAssets:      baseAssets,
        excludedSymbols: excludedSymbols,
    }
}

// ShouldTrack determines if a symbol should be tracked
func (f *DefaultSymbolFilter) ShouldTrack(symbol string) bool {
    // Check if symbol is in excluded list
    for _, excluded := range f.excludedSymbols {
        if symbol == excluded {
            return false
        }
    }
    
    // Check if symbol has an acceptable base asset
    for _, baseAsset := range f.baseAssets {
        if strings.HasSuffix(symbol, baseAsset) {
            return true
        }
    }
    
    return false
}

// ConfigurableSymbolFilter allows for more advanced filtering
type ConfigurableSymbolFilter struct {
    DefaultSymbolFilter
    minVolume       float64
    maxPrice        float64
    minPrice        float64
    exchangeClient  exchange.Client
}

// NewConfigurableSymbolFilter creates a new configurable filter
func NewConfigurableSymbolFilter(
    baseAssets []string,
    excludedSymbols []string,
    minVolume float64,
    minPrice float64,
    maxPrice float64,
    exchangeClient exchange.Client,
) *ConfigurableSymbolFilter {
    baseFilter := NewDefaultSymbolFilter(baseAssets, excludedSymbols)
    
    return &ConfigurableSymbolFilter{
        DefaultSymbolFilter: *baseFilter,
        minVolume:          minVolume,
        maxPrice:           maxPrice,
        minPrice:           minPrice,
        exchangeClient:     exchangeClient,
    }
}

// ShouldTrack implements advanced filtering logic
func (f *ConfigurableSymbolFilter) ShouldTrack(ctx context.Context, symbol string) bool {
    // First apply basic filtering
    if !f.DefaultSymbolFilter.ShouldTrack(symbol) {
        return false
    }
    
    // Get market data
    ticker, err := f.exchangeClient.GetTicker(ctx, symbol)
    if err != nil {
        // If there's an error, be conservative and don't track
        return false
    }
    
    // Apply volume filter
    if ticker.Volume < f.minVolume {
        return false
    }
    
    // Apply price filters
    if f.minPrice > 0 && ticker.Price < f.minPrice {
        return false
    }
    
    if f.maxPrice > 0 && ticker.Price > f.maxPrice {
        return false
    }
    
    return true
}
```

## Integration with Other Components

The New Coin Watcher integrates with several other components:

1. **Exchange Client**: For retrieving available trading symbols and market data.
2. **Repository Layer**: For persisting new coin data and purchase decisions.
3. **Trade Service**: For evaluating coins and executing purchases.

### Factory Creation Example

```go
// internal/domain/factory/factory.go
package factory

// CreateNewCoinWatcher creates and configures a new coin watcher
func (f *ServiceFactory) CreateNewCoinWatcher() *newcoin.Watcher {
    tradeService := f.CreateTradeService()
    
    // Create appropriate filter based on config
    var filter newcoin.SymbolFilter
    if f.config.AdvancedFiltering {
        filter = newcoin.NewConfigurableSymbolFilter(
            f.config.BaseAssets,
            f.config.ExcludedSymbols,
            f.config.MinVolume,
            f.config.MinPrice,
            f.config.MaxPrice,
            f.exchangeClient,
        )
    } else {
        filter = newcoin.NewDefaultSymbolFilter(
            f.config.BaseAssets,
            f.config.ExcludedSymbols,
        )
    }
    
    return newcoin.NewWatcher(
        f.exchangeClient,
        f.repositories.NewCoinRepository,
        f.repositories.PurchaseDecisionRepository,
        tradeService,
        f.config.CheckInterval,
        filter,
    )
}
```

## Configuration Options

The New Coin Watcher supports several configuration options:

| Option | Description | Default |
|--------|-------------|---------|
| `CheckInterval` | Frequency of checking for new coins | 60 seconds |
| `BaseAssets` | Acceptable base currencies (e.g., USDT, BTC) | ["USDT"] |
| `ExcludedSymbols` | Symbols to ignore | [] |
| `AdvancedFiltering` | Whether to use advanced filtering | false |
| `MinVolume` | Minimum 24h trading volume | 0 |
| `MinPrice` | Minimum acceptable price | 0 |
| `MaxPrice` | Maximum acceptable price | 0 |

Example configuration in YAML:

```yaml
newCoinWatcher:
  checkInterval: 60s
  baseAssets: ["USDT"]
  excludedSymbols: ["BTCUSDT", "ETHUSDT"]
  advancedFiltering: true
  minVolume: 100000
  minPrice: 0.00001
  maxPrice: 1.0
```

## Testing Guidance

### Unit Testing

```go
// internal/domain/core/newcoin/watcher_test.go
package newcoin_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/ryanlisse/cryptobot/internal/domain/core/newcoin"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository/mocks"
    "github.com/ryanlisse/cryptobot/internal/platform/exchange/mocks"
)

func TestWatcher_Start(t *testing.T) {
    // Create mocks
    mockExchange := new(mocks.ExchangeClient)
    mockNewCoinRepo := new(mocks.NewCoinRepository)
    mockDecisionRepo := new(mocks.PurchaseDecisionRepository)
    mockTradeService := new(mocks.TradeService)
    
    // Setup expectations
    mockNewCoinRepo.On("FindAll", mock.Anything).Return([]*models.NewCoin{
        {ID: 1, Symbol: "BTCUSDT", FoundAt: time.Now(), Processed: true},
    }, nil)
    
    // Create watcher
    watcher := newcoin.NewWatcher(
        mockExchange,
        mockNewCoinRepo,
        mockDecisionRepo,
        mockTradeService,
        1*time.Second,
        nil,
    )
    
    // Test start
    err := watcher.Start(context.Background())
    assert.NoError(t, err)
    assert.True(t, watcher.IsRunning())
    
    // Clean up
    watcher.Stop()
    assert.False(t, watcher.IsRunning())
    
    // Verify expectations
    mockNewCoinRepo.AssertExpectations(t)
}

func TestWatcher_CheckForNewCoins(t *testing.T) {
    // Create mocks
    mockExchange := new(mocks.ExchangeClient)
    mockNewCoinRepo := new(mocks.NewCoinRepository)
    mockDecisionRepo := new(mocks.PurchaseDecisionRepository)
    mockTradeService := new(mocks.TradeService)
    
    // Setup expectations
    mockNewCoinRepo.On("FindAll", mock.Anything).Return([]*models.NewCoin{
        {ID: 1, Symbol: "BTCUSDT", FoundAt: time.Now(), Processed: true},
    }, nil)
    
    mockExchange.On("GetAllSymbols", mock.Anything).Return([]string{
        "BTCUSDT", "ETHUSDT", "NEWUSDT", // NEWUSDT is new
    }, nil)
    
    mockNewCoinRepo.On("Create", mock.Anything, mock.MatchedBy(func(coin *models.NewCoin) bool {
        return coin.Symbol == "NEWUSDT"
    })).Return(&models.NewCoin{
        ID: 2, Symbol: "NEWUSDT", FoundAt: time.Now(),
    }, nil)
    
    // For async evaluation
    mockNewCoinRepo.On("Update", mock.Anything, mock.Anything).Return(&models.NewCoin{}, nil)
    mockTradeService.On("ShouldBuyNewCoin", mock.Anything, "NEWUSDT").Return(false, "Test reason", nil)
    mockDecisionRepo.On("Create", mock.Anything, mock.Anything).Return(&models.PurchaseDecision{}, nil)
    
    // Create watcher with longer check interval to control execution
    watcher := newcoin.NewWatcher(
        mockExchange,
        mockNewCoinRepo,
        mockDecisionRepo,
        mockTradeService,
        1*time.Hour, // long interval so we control execution
        nil,
    )
    
    // Start watcher and manually trigger check
    err := watcher.Start(context.Background())
    assert.NoError(t, err)
    
    // Use exposed method for testing
    err = watcher.CheckForNewCoins(context.Background()) // This would normally be private
    assert.NoError(t, err)
    
    // Give time for async processing
    time.Sleep(100 * time.Millisecond)
    
    // Clean up
    watcher.Stop()
    
    // Verify expectations
    mockExchange.AssertExpectations(t)
    mockNewCoinRepo.AssertExpectations(t)
    mockTradeService.AssertExpectations(t)
    mockDecisionRepo.AssertExpectations(t)
}
```

### Integration Testing

When testing the entire component together with other services:

1. Use a test database
2. Mock the exchange client
3. Create realistic test scenarios with multiple symbols

Example integration test structure:

```go
func TestNewCoinWatcherIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(db)
    
    // Create repositories with real DB
    repos := createRepositories(db)
    
    // Mock exchange
    mockExchange := new(mocks.ExchangeClient)
    
    // Create real services with mocked dependencies
    tradeService := service.NewTradeService(...)
    
    // Create watcher
    watcher := newcoin.NewWatcher(
        mockExchange,
        repos.NewCoinRepository,
        repos.PurchaseDecisionRepository,
        tradeService,
        1*time.Second,
        nil,
    )
    
    // Setup test scenario
    mockExchange.On("GetAllSymbols", mock.Anything).Return([]string{
        "BTCUSDT", "ETHUSDT", "NEWUSDT", "ANOTHERNEWUSDT",
    }, nil)
    
    // Setup ticker data for advanced filtering
    mockExchange.On("GetTicker", mock.Anything, mock.Anything).Return(&models.Ticker{
        Symbol: "NEWUSDT",
        Price: 0.5,
        Volume: 1000000,
    }, nil)
    
    // Start watcher
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err := watcher.Start(ctx)
    assert.NoError(t, err)
    
    // Wait for processing
    time.Sleep(3 * time.Second)
    
    // Verify database state
    newCoins, err := repos.NewCoinRepository.FindAll(context.Background())
    assert.NoError(t, err)
    assert.Len(t, newCoins, 2) // Should have found 2 new coins
    
    // Check decisions
    decisions, err := repos.PurchaseDecisionRepository.FindAll(context.Background())
    assert.NoError(t, err)
    assert.Len(t, decisions, 2)
    
    // Stop watcher
    watcher.Stop()
}
```

For more details on testing, see the general [Testing Strategy](../testing/overview.md) document.
