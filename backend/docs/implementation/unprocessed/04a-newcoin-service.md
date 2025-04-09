# NewCoinService Implementation

This document covers the implementation of the `NewCoinService`, which is responsible for detecting new coin listings on the MEXC exchange and processing them according to the trading bot's strategy.

## Interface Definition

```go
// internal/domain/service/newcoin.go
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
```

## Configuration

```go
// internal/core/newcoin/config.go
package newcoin

import (
    "time"
)

// Config holds configuration parameters for the NewCoinService
type Config struct {
    DetectionInterval time.Duration // How often to check for new coins
    USDTThreshold     float64       // Minimum USDT trading volume to consider
    EnableRetry       bool          // Whether to retry API calls if they fail
    MaxRetries        int           // Maximum number of retries for API calls
    RetryDelay        time.Duration // Delay between retries
}

// DefaultConfig returns sensible defaults for NewCoinService
func DefaultConfig() Config {
    return Config{
        DetectionInterval: 1 * time.Second,
        USDTThreshold:     1000.0, // Coins must have at least $1000 in 24h volume
        EnableRetry:       true,
        MaxRetries:        3,
        RetryDelay:        500 * time.Millisecond,
    }
}
```

## Service Implementation

```go
// internal/core/newcoin/service.go
package newcoin

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/service"
)

// Service implements the NewCoinService interface
type Service struct {
    exchangeService  service.ExchangeService
    newCoinRepo      service.NewCoinRepository
    purchaseDecisionTracker service.PurchaseDecisionRepository
    logRepo          service.LogRepository
    config           Config
    mu               sync.RWMutex  // Protects lastDetectionTime
    lastDetectionTime time.Time
}

// NewService creates a new NewCoinService
func NewService(
    exchangeService service.ExchangeService,
    newCoinRepo service.NewCoinRepository,
    purchaseDecisionTracker service.PurchaseDecisionRepository,
    logRepo service.LogRepository,
    config Config,
) service.NewCoinService {
    return &Service{
        exchangeService:  exchangeService,
        newCoinRepo:      newCoinRepo,
        purchaseDecisionTracker: purchaseDecisionTracker,
        logRepo:          logRepo,
        config:           config,
        lastDetectionTime: time.Now().Add(-24 * time.Hour), // Start with a day ago
    }
}
```

## Detecting New Coins

```go
// internal/core/newcoin/detect.go
package newcoin

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// DetectNewCoins identifies new coin listings on MEXC
func (s *Service) DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
    // Fetch all available coins from exchange
    // Using the MEXC V3 API, reference: https://api.mexc.com
    var newCoins []*models.NewCoin
    var err error
    
    if s.config.EnableRetry {
        // With retry logic 
        attempt := 0
        for attempt <= s.config.MaxRetries {
            newCoins, err = s.exchangeService.GetNewCoins(ctx)
            if err == nil {
                break
            }
            attempt++
            if attempt > s.config.MaxRetries {
                return nil, fmt.Errorf("failed to fetch new coins after %d attempts: %w", 
                    s.config.MaxRetries, err)
            }
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(s.config.RetryDelay):
                // Wait before retrying
            }
        }
    } else {
        // Without retry
        newCoins, err = s.exchangeService.GetNewCoins(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch new coins: %w", err)
        }
    }
    
    // Update detection time
    s.mu.Lock()
    s.lastDetectionTime = time.Now()
    s.mu.Unlock()
    
    // Filter and store new coins that we haven't seen before
    result := make([]*models.NewCoin, 0, len(newCoins))
    
    for _, coin := range newCoins {
        // Check if coin already exists in our database
        existingCoin, err := s.newCoinRepo.FindBySymbol(ctx, coin.Symbol)
        if err == nil && existingCoin != nil {
            // Update last checked time
            existingCoin.LastChecked = time.Now()
            if err := s.newCoinRepo.Update(ctx, existingCoin); err != nil {
                log.Printf("Failed to update existing coin %s: %v", coin.Symbol, err)
            }
            continue
        }
        
        // Get ticker to check volume
        ticker, err := s.exchangeService.GetTicker(ctx, coin.Symbol)
        if err != nil {
            log.Printf("Failed to get ticker for %s: %v", coin.Symbol, err)
            continue
        }
        
        // Apply volume filter
        if ticker.Volume < s.config.USDTThreshold {
            log.Printf("Skipping low volume coin %s (volume: %.2f USDT)", 
                coin.Symbol, ticker.Volume)
            continue
        }
        
        // Store new coin
        if err := s.newCoinRepo.Store(ctx, coin); err != nil {
            log.Printf("Failed to store new coin %s: %v", coin.Symbol, err)
            continue
        }
        
        // Log event
        logEvent := &models.LogEvent{
            Timestamp: time.Now(),
            Level:     models.LogLevelInfo,
            Message:   fmt.Sprintf("New coin detected: %s", coin.Symbol),
            Context:   fmt.Sprintf("Volume: %.2f USDT", ticker.Volume),
        }
        s.logRepo.Store(ctx, logEvent)
        
        result = append(result, coin)
    }
    
    return result, nil
}
```

## Processing and Evaluating New Coins

```go
// internal/core/newcoin/process.go
package newcoin

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// ProcessNewCoins evaluates new coins and makes trading decisions
func (s *Service) ProcessNewCoins(ctx context.Context) error {
    // Get active new coins
    activeCoins, err := s.newCoinRepo.FindActive(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch active new coins: %w", err)
    }
    
    for _, coin := range activeCoins {
        // Create a decision record (tracking both buys and non-buys)
        decision := &models.PurchaseDecision{
            Symbol:    coin.Symbol,
            Timestamp: time.Now(),
            Status:    models.StatusPending,
        }
        
        // Get ticker data for analysis
        ticker, err := s.exchangeService.GetTicker(ctx, coin.Symbol)
        if err != nil {
            decision.Status = models.StatusRejected
            decision.Reason = fmt.Sprintf("Failed to get ticker: %v", err)
            s.purchaseDecisionTracker.Store(ctx, decision)
            continue
        }
        
        // Save ticker price for reference
        decision.Price = ticker.Price
        
        // Apply decision criteria (simplified example)
        hourlyVolume := ticker.Volume / 24 // Rough estimate for hourly volume
        
        if hourlyVolume < 100 {
            decision.Status = models.StatusRejected
            decision.Reason = fmt.Sprintf("Insufficient hourly volume: %.2f USDT", hourlyVolume)
        } else if ticker.PriceChangePct > 20 {
            decision.Status = models.StatusRejected
            decision.Reason = fmt.Sprintf("Price increase too high: %.2f%%", ticker.PriceChangePct)
        } else {
            // Decision to purchase (in a real system, the purchase would be executed via TradeService)
            decision.Status = models.StatusPurchased
            decision.Reason = "Meets volume and price criteria"
            
            // Log decision to purchase
            logEvent := &models.LogEvent{
                Timestamp: time.Now(),
                Level:     models.LogLevelInfo,
                Message:   fmt.Sprintf("Purchase decision for %s: PURCHASE", coin.Symbol),
                Context:   fmt.Sprintf("Price: %.8f, Volume: %.2f", ticker.Price, ticker.Volume),
            }
            s.logRepo.Store(ctx, logEvent)
        }
        
        // Store the decision
        if err := s.purchaseDecisionTracker.Store(ctx, decision); err != nil {
            log.Printf("Failed to store purchase decision for %s: %v", coin.Symbol, err)
        }
    }
    
    return nil
}

// ArchiveOldCoins marks old coins as inactive
func (s *Service) ArchiveOldCoins(ctx context.Context, daysOld int) error {
    cutoffTime := time.Now().Add(-time.Duration(daysOld) * 24 * time.Hour)
    
    activeCoins, err := s.newCoinRepo.FindActive(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch active coins: %w", err)
    }
    
    for _, coin := range activeCoins {
        if coin.DetectedAt.Before(cutoffTime) {
            if err := s.newCoinRepo.Archive(ctx, coin.ID); err != nil {
                log.Printf("Failed to archive old coin %s: %v", coin.Symbol, err)
                continue
            }
            
            // Log event
            logEvent := &models.LogEvent{
                Timestamp: time.Now(),
                Level:     models.LogLevelInfo,
                Message:   fmt.Sprintf("Archived old coin: %s", coin.Symbol),
                Context:   fmt.Sprintf("Detected at: %s", coin.DetectedAt.Format(time.RFC3339)),
            }
            s.logRepo.Store(ctx, logEvent)
        }
    }
    
    return nil
}
```

## NewCoin Service Testing

```go
// internal/core/newcoin/service_test.go
package newcoin

import (
    "context"
    "testing"
    "time"
    
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/service"
)

// Mock implementations of dependencies
type mockExchangeService struct {
    newCoins  []*models.NewCoin
    tickers   map[string]*models.Ticker
    callCount int
}

func (m *mockExchangeService) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
    m.callCount++
    return m.newCoins, nil
}

func (m *mockExchangeService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
    ticker, ok := m.tickers[symbol]
    if !ok {
        ticker = &models.Ticker{
            Symbol:  symbol,
            Price:   1000.0,
            Volume:  5000.0,
        }
    }
    return ticker, nil
}

// Other required methods of the ExchangeService interface...

type mockNewCoinRepo struct {
    coins map[string]*models.NewCoin
}

func (m *mockNewCoinRepo) Store(ctx context.Context, coin *models.NewCoin) error {
    m.coins[coin.Symbol] = coin
    return nil
}

func (m *mockNewCoinRepo) FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error) {
    coin, ok := m.coins[symbol]
    if !ok {
        return nil, nil // Not found, no error
    }
    return coin, nil
}

func (m *mockNewCoinRepo) FindActive(ctx context.Context) ([]*models.NewCoin, error) {
    activeCoins := make([]*models.NewCoin, 0)
    for _, coin := range m.coins {
        if coin.IsActive {
            activeCoins = append(activeCoins, coin)
        }
    }
    return activeCoins, nil
}

// Other required methods...

// Similarly, create mock implementations for PurchaseDecisionRepository and LogRepository

func TestDetectNewCoins(t *testing.T) {
    // Initialize test data
    now := time.Now()
    testCoins := []*models.NewCoin{
        {Symbol: "BTC/USDT", DetectedAt: now, IsActive: true},
        {Symbol: "ETH/USDT", DetectedAt: now, IsActive: true},
        {Symbol: "NEW/USDT", DetectedAt: now, IsActive: true},
    }
    
    testTickers := map[string]*models.Ticker{
        "BTC/USDT": {Symbol: "BTC/USDT", Price: 50000.0, Volume: 10000.0},
        "ETH/USDT": {Symbol: "ETH/USDT", Price: 3000.0, Volume: 5000.0},
        "NEW/USDT": {Symbol: "NEW/USDT", Price: 1.0, Volume: 100.0}, // Low volume
    }
    
    // Create mocks
    mockExchange := &mockExchangeService{
        newCoins: testCoins,
        tickers:  testTickers,
    }
    
    mockRepo := &mockNewCoinRepo{
        coins: make(map[string]*models.NewCoin),
    }
    
    // Add a known coin to the repository
    mockRepo.coins["BTC/USDT"] = testCoins[0]
    
    // Create service with mocks
    service := NewService(
        mockExchange,
        mockRepo,
        nil, // You'd add a mock purchase decision repo here
        nil, // You'd add a mock log repo here
        DefaultConfig(),
    )
    
    // Run the test
    ctx := context.Background()
    newCoins, err := service.DetectNewCoins(ctx)
    
    // Verify
    if err != nil {
        t.Errorf("DetectNewCoins returned error: %v", err)
    }
    
    // We should detect ETH but not BTC (already known) or NEW (too low volume)
    if len(newCoins) != 1 {
        t.Errorf("Expected 1 new coin, got %d", len(newCoins))
    }
    
    if len(newCoins) > 0 && newCoins[0].Symbol != "ETH/USDT" {
        t.Errorf("Expected coin ETH/USDT, got %s", newCoins[0].Symbol)
    }
    
    // Verify that the exchange was called
    if mockExchange.callCount != 1 {
        t.Errorf("Expected 1 call to GetNewCoins, got %d", mockExchange.callCount)
    }
}
```

This implementation provides a robust service for detecting new coins on the MEXC exchange. The service includes:

1. **Configuration Options**: Easily configurable parameters for detection thresholds, retry behavior, and intervals.
2. **Retry Mechanism**: Built-in retry logic for API calls to handle transient failures.
3. **Volume Filtering**: Filters out coins with insufficient trading volume.
4. **Decision Tracking**: Records all purchase decisions, both positive and negative, for analysis.
5. **Logging**: Comprehensive logging of new coin detections and purchase decisions.
6. **Coin Management**: Automatically archives old coins that are no longer of interest.
