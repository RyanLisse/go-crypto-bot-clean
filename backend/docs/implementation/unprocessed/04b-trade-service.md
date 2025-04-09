# TradeService Implementation

This document covers the implementation of the `TradeService`, which is responsible for the core trading logic, including executing purchases and managing stop-loss/take-profit strategies.

## Interface Definition

```go
// internal/domain/service/trade.go
package service

import (
    "context"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// TradeService defines operations for trading logic
type TradeService interface {
    EvaluatePurchaseDecision(ctx context.Context, symbol string) (*models.PurchaseDecision, error)
    ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error)
    CheckStopLoss(ctx context.Context, coin *models.BoughtCoin) (bool, error)
    CheckTakeProfit(ctx context.Context, coin *models.BoughtCoin) (bool, error)
    SellCoin(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error)
}
```

## Configuration

```go
// internal/core/trade/config.go
package trade

import (
    "errors"
)

// Config represents configuration for the TradeService
type Config struct {
    USDTPerTrade     float64   // Amount of USDT to use per trade
    StopLossPercent  float64   // Default stop loss percentage (e.g. 15.0 for 15%)
    TakeProfitLevels []float64 // Multiple take profit levels (e.g. [5.0, 10.0, 15.0, 20.0])
    SellPercentages  []float64 // Percentage of position to sell at each TP level
    EnableFallback   bool      // Whether to enable fallback to REST for price checks
}

// DefaultConfig returns sensible default values for trade configuration
func DefaultConfig() Config {
    return Config{
        USDTPerTrade:     25.0,  // $25 per trade
        StopLossPercent:  15.0,  // 15% stop loss
        TakeProfitLevels: []float64{5.0, 10.0, 15.0, 20.0},
        SellPercentages:  []float64{0.25, 0.25, 0.25, 0.25}, // Sell 25% at each level
        EnableFallback:   true,
    }
}

// Validate checks if the configuration is valid
func (c Config) Validate() error {
    if c.USDTPerTrade <= 0 {
        return errors.New("USDTPerTrade must be positive")
    }
    
    if c.StopLossPercent <= 0 || c.StopLossPercent >= 100 {
        return errors.New("StopLossPercent must be between 0 and 100")
    }
    
    if len(c.TakeProfitLevels) == 0 {
        return errors.New("TakeProfitLevels cannot be empty")
    }
    
    if len(c.TakeProfitLevels) != len(c.SellPercentages) {
        return errors.New("TakeProfitLevels and SellPercentages must have the same length")
    }
    
    totalSellPercentage := 0.0
    for _, pct := range c.SellPercentages {
        if pct <= 0 {
            return errors.New("all SellPercentages must be positive")
        }
        totalSellPercentage += pct
    }
    
    if totalSellPercentage > 1.01 { // Allow for small floating point errors
        return errors.New("sum of SellPercentages cannot exceed 100%")
    }
    
    return nil
}
```

## Service Implementation

```go
// internal/core/trade/service.go
package trade

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/service"
)

// Service implements the TradeService interface
type Service struct {
    exchangeService  service.ExchangeService
    boughtCoinRepo   service.BoughtCoinRepository
    logRepo          service.LogRepository
    config           Config
}

// NewService creates a new TradeService
func NewService(
    exchangeService service.ExchangeService,
    boughtCoinRepo service.BoughtCoinRepository,
    logRepo service.LogRepository,
    config Config,
) (service.TradeService, error) {
    // Validate config
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid trade service config: %w", err)
    }
    
    return &Service{
        exchangeService: exchangeService,
        boughtCoinRepo:  boughtCoinRepo,
        logRepo:         logRepo,
        config:          config,
    }, nil
}
```

## Purchase Decision Logic

```go
// internal/core/trade/purchase.go
package trade

import (
    "context"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// EvaluatePurchaseDecision determines whether to buy a given coin
func (s *Service) EvaluatePurchaseDecision(ctx context.Context, symbol string) (*models.PurchaseDecision, error) {
    decision := &models.PurchaseDecision{
        Symbol:    symbol,
        Timestamp: time.Now(),
        Status:    models.StatusPending,
    }

    // Check if we already own this coin
    existingCoin, err := s.boughtCoinRepo.FindBySymbol(ctx, symbol)
    if err == nil && existingCoin != nil && !existingCoin.IsDeleted {
        decision.Status = models.StatusRejected
        decision.Reason = "Already own this coin"
        return decision, nil
    }

    // Get ticker data
    ticker, err := s.exchangeService.GetTicker(ctx, symbol)
    if err != nil {
        decision.Status = models.StatusRejected
        decision.Reason = fmt.Sprintf("Failed to get ticker: %v", err)
        return decision, nil
    }

    decision.Price = ticker.Price

    // Check wallet balance
    wallet, err := s.exchangeService.GetWallet(ctx)
    if err != nil {
        decision.Status = models.StatusRejected
        decision.Reason = fmt.Sprintf("Failed to get wallet balance: %v", err)
        return decision, nil
    }

    if wallet.USDT < s.config.USDTPerTrade {
        decision.Status = models.StatusRejected
        decision.Reason = fmt.Sprintf("Insufficient USDT balance: %.2f (need %.2f)", 
            wallet.USDT, s.config.USDTPerTrade)
        return decision, nil
    }

    // Determine whether to buy based on specific criteria
    // This is where you would implement your trading strategy
    // For example, checking price action, volume, etc.
    
    // In this simplified example, we'll make a decision based on:
    // 1. 24h volume must be significant
    // 2. Price shouldn't have risen too much already
    
    if ticker.Volume < 10000 {
        decision.Status = models.StatusRejected
        decision.Reason = fmt.Sprintf("Insufficient 24h volume: %.2f USDT", ticker.Volume)
        return decision, nil
    }
    
    if ticker.PriceChangePct > 15 {
        decision.Status = models.StatusRejected
        decision.Reason = fmt.Sprintf("Price already increased by %.2f%% in 24h", ticker.PriceChangePct)
        return decision, nil
    }
    
    // Decision is to purchase
    decision.Status = models.StatusPurchased
    decision.Reason = "Meets purchase criteria: sufficient volume and reasonable price action"
    
    return decision, nil
}

// ExecutePurchase buys a coin using the configured strategy
func (s *Service) ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
    if amount <= 0 {
        amount = s.config.USDTPerTrade
    }

    // Get current price
    ticker, err := s.exchangeService.GetTicker(ctx, symbol)
    if err != nil {
        return nil, fmt.Errorf("failed to get ticker for %s: %w", symbol, err)
    }

    // Calculate quantity to buy
    quantity := amount / ticker.Price

    // Round quantity to appropriate precision (this would be exchange-specific)
    // For simplicity, we'll assume 8 decimal places for all coins
    quantity = float64(int64(quantity*100000000)) / 100000000

    // Place market buy order
    order := &models.Order{
        Symbol:   symbol,
        Side:     models.OrderSideBuy,
        Type:     models.OrderTypeMarket,
        Quantity: quantity,
    }

    executedOrder, err := s.exchangeService.PlaceOrder(ctx, order)
    if err != nil {
        return nil, fmt.Errorf("failed to place buy order for %s: %w", symbol, err)
    }

    // Create a BoughtCoin record
    takeProfitLevels := make([]models.TakeProfitLevel, len(s.config.TakeProfitLevels))
    for i, levelPct := range s.config.TakeProfitLevels {
        takeProfitLevels[i] = models.TakeProfitLevel{
            Percentage:   levelPct,
            SellQuantity: quantity * s.config.SellPercentages[i],
            IsReached:    false,
        }
    }

    boughtCoin := &models.BoughtCoin{
        Symbol:        symbol,
        PurchasePrice: executedOrder.AvgPrice,
        Quantity:      executedOrder.FilledQty,
        PurchaseTime:  time.Now(),
        IsDeleted:     false,
        StopLossPrice: executedOrder.AvgPrice * (1 - s.config.StopLossPercent/100),
        TakeProfitLevels: takeProfitLevels,
    }

    // Save to repository
    if err := s.boughtCoinRepo.Store(ctx, boughtCoin); err != nil {
        return nil, fmt.Errorf("failed to store bought coin record: %w", err)
    }

    // Log the purchase
    logEvent := &models.LogEvent{
        Timestamp: time.Now(),
        Level:     models.LogLevelInfo,
        Message:   fmt.Sprintf("Purchased %s: %.8f @ %.8f USDT", 
            symbol, executedOrder.FilledQty, executedOrder.AvgPrice),
        Context:   fmt.Sprintf("OrderID: %s, Total: %.2f USDT", 
            executedOrder.ID, executedOrder.FilledQty*executedOrder.AvgPrice),
    }
    s.logRepo.Store(ctx, logEvent)

    return boughtCoin, nil
}
```

## Stop-Loss and Take-Profit Management

```go
// internal/core/trade/stoploss_takeprofit.go
package trade

import (
    "context"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// CheckStopLoss checks if stop loss should be triggered for a coin
func (s *Service) CheckStopLoss(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
    // Get current price
    ticker, err := s.exchangeService.GetTicker(ctx, coin.Symbol)
    if err != nil {
        // If fallback is enabled, try one more time
        if s.config.EnableFallback {
            // Wait a bit and retry
            select {
            case <-ctx.Done():
                return false, ctx.Err()
            case <-time.After(500 * time.Millisecond):
                // Retry with a slight delay
            }
            
            ticker, err = s.exchangeService.GetTicker(ctx, coin.Symbol)
            if err != nil {
                return false, fmt.Errorf("failed to get ticker after fallback for %s: %w", 
                    coin.Symbol, err)
            }
        } else {
            return false, fmt.Errorf("failed to get ticker for %s: %w", coin.Symbol, err)
        }
    }

    // Check if current price is at or below stop loss
    if ticker.Price <= coin.StopLossPrice {
        logEvent := &models.LogEvent{
            Timestamp: time.Now(),
            Level:     models.LogLevelWarning,
            Message:   fmt.Sprintf("Stop loss triggered for %s @ %.8f", coin.Symbol, ticker.Price),
            Context:   fmt.Sprintf("Purchase price: %.8f, Stop loss: %.8f, Loss: %.2f%%", 
                coin.PurchasePrice, coin.StopLossPrice, 
                (ticker.Price-coin.PurchasePrice)/coin.PurchasePrice*100),
        }
        s.logRepo.Store(ctx, logEvent)
        
        return true, nil
    }

    return false, nil
}

// CheckTakeProfit checks if take profit should be triggered for a coin
func (s *Service) CheckTakeProfit(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
    // Get current price
    ticker, err := s.exchangeService.GetTicker(ctx, coin.Symbol)
    if err != nil {
        // If fallback is enabled, try one more time
        if s.config.EnableFallback {
            // Wait a bit and retry
            select {
            case <-ctx.Done():
                return false, ctx.Err()
            case <-time.After(500 * time.Millisecond):
                // Retry with a slight delay
            }
            
            ticker, err = s.exchangeService.GetTicker(ctx, coin.Symbol)
            if err != nil {
                return false, fmt.Errorf("failed to get ticker after fallback for %s: %w", 
                    coin.Symbol, err)
            }
        } else {
            return false, fmt.Errorf("failed to get ticker for %s: %w", coin.Symbol, err)
        }
    }

    // Calculate current profit percentage
    profitPercent := (ticker.Price - coin.PurchasePrice) / coin.PurchasePrice * 100

    // Check take profit levels (from highest to lowest)
    for i := len(coin.TakeProfitLevels) - 1; i >= 0; i-- {
        level := &coin.TakeProfitLevels[i]
        if !level.IsReached && profitPercent >= level.Percentage {
            logEvent := &models.LogEvent{
                Timestamp: time.Now(),
                Level:     models.LogLevelInfo,
                Message:   fmt.Sprintf("Take profit triggered for %s @ %.8f (Level: %.1f%%)", 
                    coin.Symbol, ticker.Price, level.Percentage),
                Context:   fmt.Sprintf("Purchase price: %.8f, Current profit: %.2f%%", 
                    coin.PurchasePrice, profitPercent),
            }
            s.logRepo.Store(ctx, logEvent)
            
            level.IsReached = true
            
            // Update the coin in repository
            if err := s.boughtCoinRepo.Update(ctx, coin); err != nil {
                return false, fmt.Errorf("failed to update take profit status: %w", err)
            }
            
            return true, nil
        }
    }

    return false, nil
}
```

## Selling Logic

```go
// internal/core/trade/sell.go
package trade

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// SellCoin sells a specified amount of a coin
func (s *Service) SellCoin(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error) {
    // Validate amount
    if amount <= 0 || amount > coin.Quantity {
        return nil, errors.New("invalid sell amount")
    }

    // Place market sell order
    order := &models.Order{
        Symbol:   coin.Symbol,
        Side:     models.OrderSideSell,
        Type:     models.OrderTypeMarket,
        Quantity: amount,
    }

    executedOrder, err := s.exchangeService.PlaceOrder(ctx, order)
    if err != nil {
        return nil, fmt.Errorf("failed to place sell order for %s: %w", coin.Symbol, err)
    }

    // Update coin quantity
    coin.Quantity -= amount
    
    // If all sold, mark as deleted
    if coin.Quantity <= 0.000001 { // Small epsilon to handle floating point errors
        coin.IsDeleted = true
    }

    // Update the coin in repository
    if err := s.boughtCoinRepo.Update(ctx, coin); err != nil {
        return nil, fmt.Errorf("failed to update coin after selling: %w", err)
    }

    // Calculate profit
    profitPercent := (executedOrder.AvgPrice - coin.PurchasePrice) / coin.PurchasePrice * 100
    profitAmount := amount * (executedOrder.AvgPrice - coin.PurchasePrice)

    // Log the sale
    logEvent := &models.LogEvent{
        Timestamp: time.Now(),
        Level:     models.LogLevelInfo,
        Message:   fmt.Sprintf("Sold %s: %.8f @ %.8f USDT", 
            coin.Symbol, executedOrder.FilledQty, executedOrder.AvgPrice),
        Context:   fmt.Sprintf("OrderID: %s, Total: %.2f USDT, Profit: %.2f%% (%.2f USDT)", 
            executedOrder.ID, 
            executedOrder.FilledQty*executedOrder.AvgPrice,
            profitPercent, profitAmount),
    }
    s.logRepo.Store(ctx, logEvent)

    return executedOrder, nil
}
```

## Unit Testing

```go
// internal/core/trade/service_test.go
package trade

import (
    "context"
    "testing"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// Mock implementations for testing
type mockExchangeService struct {
    tickers map[string]*models.Ticker
    wallet  *models.Wallet
    orders  []*models.Order
}

func (m *mockExchangeService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
    ticker, ok := m.tickers[symbol]
    if !ok {
        return nil, fmt.Errorf("ticker not found for %s", symbol)
    }
    return ticker, nil
}

func (m *mockExchangeService) GetWallet(ctx context.Context) (*models.Wallet, error) {
    return m.wallet, nil
}

func (m *mockExchangeService) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
    // Simplified mock that just records orders and returns a dummy execution
    m.orders = append(m.orders, order)
    
    // Create a dummy executed order
    executedOrder := *order
    executedOrder.ID = "order123"
    executedOrder.Status = models.OrderStatusFilled
    executedOrder.FilledQty = order.Quantity
    executedOrder.AvgPrice = m.tickers[order.Symbol].Price
    executedOrder.CreatedAt = time.Now()
    executedOrder.UpdatedAt = time.Now()
    
    return &executedOrder, nil
}

// ... other required mock methods ...

type mockBoughtCoinRepo struct {
    coins map[string]*models.BoughtCoin
}

func (m *mockBoughtCoinRepo) Store(ctx context.Context, coin *models.BoughtCoin) error {
    m.coins[coin.Symbol] = coin
    return nil
}

func (m *mockBoughtCoinRepo) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
    coin, ok := m.coins[symbol]
    if !ok {
        return nil, nil // Not found, but no error
    }
    return coin, nil
}

func (m *mockBoughtCoinRepo) Update(ctx context.Context, coin *models.BoughtCoin) error {
    m.coins[coin.Symbol] = coin
    return nil
}

// ... other mock methods ...

// Similarly, implement a mock LogRepository

func TestEvaluatePurchaseDecision(t *testing.T) {
    // Set up test data
    mockExchange := &mockExchangeService{
        tickers: map[string]*models.Ticker{
            "BTC/USDT": {
                Symbol:         "BTC/USDT",
                Price:          50000.0,
                PriceChange:    1000.0,
                PriceChangePct: 2.0, // 2% increase - acceptable
                Volume:         1000000.0, // High volume - acceptable
            },
            "ETH/USDT": {
                Symbol:         "ETH/USDT",
                Price:          3000.0,
                PriceChange:    600.0,
                PriceChangePct: 20.0, // 20% increase - too high
                Volume:         500000.0,
            },
            "LOW/USDT": {
                Symbol:         "LOW/USDT",
                Price:          1.0,
                PriceChange:    0.1,
                PriceChangePct: 10.0,
                Volume:         5000.0, // Low volume - too low
            },
        },
        wallet: &models.Wallet{
            USDT: 100.0, // Enough for a trade
        },
    }
    
    mockRepo := &mockBoughtCoinRepo{
        coins: map[string]*models.BoughtCoin{
            // We already own this coin
            "ETH/USDT": {
                Symbol:        "ETH/USDT",
                PurchasePrice: 2500.0,
                Quantity:      0.1,
                IsDeleted:     false,
            },
        },
    }
    
    mockLogRepo := &mockLogRepo{}
    
    // Create service with mocks
    service, err := NewService(mockExchange, mockRepo, mockLogRepo, DefaultConfig())
    if err != nil {
        t.Fatalf("Failed to create service: %v", err)
    }
    
    // Test cases
    testCases := []struct {
        name           string
        symbol         string
        expectedStatus models.PurchaseDecisionStatus
        expectedReason string
    }{
        {
            name:           "BTC should be purchased",
            symbol:         "BTC/USDT",
            expectedStatus: models.StatusPurchased,
        },
        {
            name:           "ETH already owned",
            symbol:         "ETH/USDT",
            expectedStatus: models.StatusRejected,
            expectedReason: "Already own this coin",
        },
        {
            name:           "LOW volume too low",
            symbol:         "LOW/USDT",
            expectedStatus: models.StatusRejected,
            expectedReason: "Insufficient",
        },
    }
    
    ctx := context.Background()
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            decision, err := service.EvaluatePurchaseDecision(ctx, tc.symbol)
            
            if err != nil {
                t.Fatalf("EvaluatePurchaseDecision returned error: %v", err)
            }
            
            if decision.Status != tc.expectedStatus {
                t.Errorf("Expected status %s, got %s", tc.expectedStatus, decision.Status)
            }
            
            if tc.expectedReason != "" && !strings.Contains(decision.Reason, tc.expectedReason) {
                t.Errorf("Expected reason to contain '%s', got '%s'", tc.expectedReason, decision.Reason)
            }
        })
    }
}

// Additional tests for ExecutePurchase, CheckStopLoss, CheckTakeProfit, SellCoin, etc.
```

## Integration with the Main Bot

This `TradeService` is designed to be integrated into the main bot logic. Here's an example of how it might be used in the bot's main loop:

```go
// internal/bot/bot.go
package bot

import (
    "context"
    "log"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/service"
)

// Bot represents the main trading bot
type Bot struct {
    tradeService     service.TradeService
    newCoinService   service.NewCoinService
    portfolioService service.PortfolioService
    boughtCoinRepo   service.BoughtCoinRepository
    
    // Control channels
    stopCh           chan struct{}
    isRunning        bool
}

// NewBot creates a new trading bot
func NewBot(
    tradeService service.TradeService,
    newCoinService service.NewCoinService,
    portfolioService service.PortfolioService,
    boughtCoinRepo service.BoughtCoinRepository,
) *Bot {
    return &Bot{
        tradeService:     tradeService,
        newCoinService:   newCoinService,
        portfolioService: portfolioService,
        boughtCoinRepo:   boughtCoinRepo,
        stopCh:           make(chan struct{}),
        isRunning:        false,
    }
}

// Start begins the bot's main loop
func (b *Bot) Start() {
    if b.isRunning {
        log.Println("Bot is already running")
        return
    }
    
    b.isRunning = true
    log.Println("Starting trading bot...")
    
    go b.runNewCoinWatcher()
    go b.runPositionManager()
}

// runNewCoinWatcher periodically checks for new coins and evaluates them
func (b *Bot) runNewCoinWatcher() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-b.stopCh:
            return
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            
            // Detect and process new coins
            newCoins, err := b.newCoinService.DetectNewCoins(ctx)
            if err != nil {
                log.Printf("Error detecting new coins: %v", err)
            } else if len(newCoins) > 0 {
                log.Printf("Detected %d new coins", len(newCoins))
                
                // Process new coins and make purchase decisions
                if err := b.newCoinService.ProcessNewCoins(ctx); err != nil {
                    log.Printf("Error processing new coins: %v", err)
                }
            }
            
            cancel()
        }
    }
}

// runPositionManager monitors active positions and manages stop-loss/take-profit
func (b *Bot) runPositionManager() {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-b.stopCh:
            return
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            
            // Get active positions
            activeCoins, err := b.portfolioService.GetActiveTrades(ctx)
            if err != nil {
                log.Printf("Error getting active trades: %v", err)
                cancel()
                continue
            }
            
            // Check stop-loss and take-profit for each position
            for _, coin := range activeCoins {
                // Check stop-loss
                stopLossTriggered, err := b.tradeService.CheckStopLoss(ctx, coin)
                if err != nil {
                    log.Printf("Error checking stop-loss for %s: %v", coin.Symbol, err)
                    continue
                }
                
                if stopLossTriggered {
                    log.Printf("Stop-loss triggered for %s, selling entire position", coin.Symbol)
                    if _, err := b.tradeService.SellCoin(ctx, coin, coin.Quantity); err != nil {
                        log.Printf("Error executing stop-loss for %s: %v", coin.Symbol, err)
                    }
                    continue // Skip take-profit check if stop-loss was triggered
                }
                
                // Check take-profit
                takeProfitTriggered, err := b.tradeService.CheckTakeProfit(ctx, coin)
                if err != nil {
                    log.Printf("Error checking take-profit for %s: %v", coin.Symbol, err)
                    continue
                }
                
                if takeProfitTriggered {
                    // Get the level that was triggered (the one that was just marked as reached)
                    var sellAmount float64
                    for _, level := range coin.TakeProfitLevels {
                        if level.IsReached {
                            sellAmount = level.SellQuantity
                            break
                        }
                    }
                    
                    if sellAmount > 0 {
                        log.Printf("Take-profit triggered for %s, selling %.8f", coin.Symbol, sellAmount)
                        if _, err := b.tradeService.SellCoin(ctx, coin, sellAmount); err != nil {
                            log.Printf("Error executing take-profit for %s: %v", coin.Symbol, err)
                        }
                    }
                }
            }
            
            cancel()
        }
    }
}

// Stop gracefully stops the bot
func (b *Bot) Stop() {
    if !b.isRunning {
        return
    }
    
    log.Println("Stopping trading bot...")
    close(b.stopCh)
    b.isRunning = false
}
```

The TradeService implementation provides a robust foundation for trading strategies with customizable stop-loss and take-profit levels. The service handles all order execution, position management, and decision-making, allowing the main bot logic to focus on higher-level coordination.
