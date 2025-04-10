# PortfolioService Implementation

This document covers the implementation of the `PortfolioService`, which is responsible for portfolio management, including tracking positions and calculating performance metrics.

## Interface Definition

```go
// internal/domain/service/portfolio.go
package service

import (
    "context"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// PortfolioService defines operations for portfolio management
type PortfolioService interface {
    GetPortfolioValue(ctx context.Context) (float64, error)
    GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error)
    GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error)
}
```

## Service Implementation

```go
// internal/core/portfolio/service.go
package portfolio

import (
    "context"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/service"
)

// Service implements the PortfolioService interface
type Service struct {
    exchangeService  service.ExchangeService
    boughtCoinRepo   service.BoughtCoinRepository
}

// NewService creates a new PortfolioService
func NewService(
    exchangeService service.ExchangeService,
    boughtCoinRepo service.BoughtCoinRepository,
) service.PortfolioService {
    return &Service{
        exchangeService: exchangeService,
        boughtCoinRepo:  boughtCoinRepo,
    }
}
```

## Portfolio Valuation Logic

```go
// internal/core/portfolio/valuation.go
package portfolio

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// GetPortfolioValue calculates the current total value of the portfolio
func (s *Service) GetPortfolioValue(ctx context.Context) (float64, error) {
    // Get wallet balance first
    wallet, err := s.exchangeService.GetWallet(ctx)
    if err != nil {
        return 0, fmt.Errorf("failed to get wallet: %w", err)
    }

    // Start with USDT balance
    totalValue := wallet.USDT

    // Get active trades
    activeCoins, err := s.boughtCoinRepo.FindActive(ctx)
    if err != nil {
        return 0, fmt.Errorf("failed to get active trades: %w", err)
    }

    // Fetch current prices and calculate value concurrently for better performance
    var wg sync.WaitGroup
    var mu sync.Mutex // To protect concurrent access to totalValue
    
    for _, coin := range activeCoins {
        wg.Add(1)
        go func(coin *models.BoughtCoin) {
            defer wg.Done()
            
            ticker, err := s.exchangeService.GetTicker(ctx, coin.Symbol)
            if err != nil {
                // Log error but continue with other coins
                fmt.Printf("Failed to get ticker for %s: %v\n", coin.Symbol, err)
                return
            }

            // Calculate coin value and add to total under lock
            coinValue := coin.Quantity * ticker.Price
            
            mu.Lock()
            totalValue += coinValue
            mu.Unlock()
        }(coin)
    }
    
    // Wait for all price fetches to complete
    wg.Wait()

    return totalValue, nil
}

// GetActiveTrades retrieves all active trade positions
func (s *Service) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
    activeCoins, err := s.boughtCoinRepo.FindActive(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get active trades: %w", err)
    }

    // Enrich with current prices and profit calculations concurrently
    var wg sync.WaitGroup
    
    for _, coin := range activeCoins {
        wg.Add(1)
        go func(coin *models.BoughtCoin) {
            defer wg.Done()
            
            ticker, err := s.exchangeService.GetTicker(ctx, coin.Symbol)
            if err != nil {
                // Log error but continue with other coins
                fmt.Printf("Failed to get ticker for %s: %v\n", coin.Symbol, err)
                return
            }

            // Calculate current value and profit
            coin.CurrentPrice = ticker.Price
            coin.CurrentValue = coin.Quantity * ticker.Price
            coin.ProfitPercentage = (ticker.Price - coin.PurchasePrice) / coin.PurchasePrice * 100
        }(coin)
    }
    
    // Wait for all enrichments to complete
    wg.Wait()

    return activeCoins, nil
}
```

## Performance Analytics Logic

```go
// internal/core/portfolio/analytics.go
package portfolio

import (
    "context"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// GetTradePerformance calculates performance metrics for a given time range
func (s *Service) GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error) {
    // Determine time range
    var startTime time.Time
    endTime := time.Now()

    switch timeRange {
    case "day":
        startTime = endTime.AddDate(0, 0, -1)
    case "week":
        startTime = endTime.AddDate(0, 0, -7)
    case "month":
        startTime = endTime.AddDate(0, -1, 0)
    case "year":
        startTime = endTime.AddDate(-1, 0, 0)
    case "all":
        startTime = time.Time{} // Zero time to get all history
    default:
        startTime = endTime.AddDate(0, 0, -7) // Default to one week
    }

    // Initialize metrics
    metrics := &models.PerformanceMetrics{
        TotalTrades:   0,
        WinningTrades: 0,
        LosingTrades:  0,
        StartValue:    0,
        CurrentValue:  0,
    }

    // Query completed trades (soft deleted)
    completedTrades, err := s.boughtCoinRepo.FindCompletedBetween(ctx, startTime, endTime)
    if err != nil {
        return nil, fmt.Errorf("failed to get completed trades: %w", err)
    }

    // Calculate metrics from completed trades
    totalProfit := 0.0
    totalHoldingTime := 0.0
    largestProfit := -1000.0 // Start with a very negative number
    largestLoss := 0.0       // Start with zero
    
    for _, trade := range completedTrades {
        metrics.TotalTrades++
        
        // Calculate profit percentage
        // In a real system, we would have the actual sell price recorded
        // For this implementation, we'll use the currentPrice if available
        profitPct := 0.0
        
        if trade.CurrentPrice > 0 {
            profitPct = (trade.CurrentPrice - trade.PurchasePrice) / trade.PurchasePrice * 100
        } else {
            // If current price isn't available, estimate from take profit levels
            // This is a simplified approach - in a real system, you'd record actual sell prices
            sellPrice := 0.0
            for _, level := range trade.TakeProfitLevels {
                if level.IsReached {
                    sellPrice = trade.PurchasePrice * (1 + level.Percentage/100)
                    break
                }
            }
            
            if sellPrice == 0 {
                // If no TP was reached, assume stoploss was hit
                sellPrice = trade.StopLossPrice
            }
            
            profitPct = (sellPrice - trade.PurchasePrice) / trade.PurchasePrice * 100
        }
        
        totalProfit += profitPct
        
        // Track win/loss statistics
        if profitPct >= 0 {
            metrics.WinningTrades++
            if profitPct > largestProfit {
                largestProfit = profitPct
            }
        } else {
            metrics.LosingTrades++
            if profitPct < largestLoss {
                largestLoss = profitPct
            }
        }
        
        // Calculate holding time
        var sellTime time.Time
        if trade.SellTime.IsZero() {
            // If sell time isn't recorded, estimate based on typical trade duration
            // In a real system, you'd record the actual sell time
            sellTime = endTime
        } else {
            sellTime = trade.SellTime
        }
        
        holdingHours := sellTime.Sub(trade.PurchaseTime).Hours()
        totalHoldingTime += holdingHours
    }
    
    // Calculate aggregated metrics
    if metrics.TotalTrades > 0 {
        metrics.AverageProfit = totalProfit / float64(metrics.TotalTrades)
        metrics.AverageHoldingTime = totalHoldingTime / float64(metrics.TotalTrades)
        metrics.WinRate = float64(metrics.WinningTrades) / float64(metrics.TotalTrades) * 100
    }
    
    metrics.LargestProfit = largestProfit
    metrics.LargestLoss = largestLoss
    
    // Calculate current portfolio value
    currentValue, err := s.GetPortfolioValue(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get current portfolio value: %w", err)
    }
    
    metrics.CurrentValue = currentValue
    
    // For start value, we need a more complex calculation that would ideally use historical data
    // For this implementation, we'll use a simplified approach
    metrics.StartValue = currentValue / (1 + metrics.AverageProfit/100)
    metrics.TotalProfit = currentValue - metrics.StartValue
    if metrics.StartValue > 0 {
        metrics.TotalProfitPct = (metrics.TotalProfit / metrics.StartValue) * 100
    }
    
    return metrics, nil
}
```

## Advanced Portfolio Analytics

```go
// internal/core/portfolio/advanced_analytics.go
package portfolio

import (
    "context"
    "fmt"
    "math"
    "time"

    "github.com/ryanlisse/cryptobot/internal/domain/models"
)

// GetTradeFrequency calculates how often trades are executed
func (s *Service) GetTradeFrequency(ctx context.Context, days int) (float64, error) {
    // Calculate start time
    startTime := time.Now().AddDate(0, 0, -days)
    
    // Get trades in the period
    trades, err := s.boughtCoinRepo.FindAllBetween(ctx, startTime, time.Now())
    if err != nil {
        return 0, fmt.Errorf("failed to get trades: %w", err)
    }
    
    // Calculate trades per day
    tradesPerDay := float64(len(trades)) / float64(days)
    
    return tradesPerDay, nil
}

// CalculateDrawdown calculates the maximum portfolio drawdown
func (s *Service) CalculateDrawdown(ctx context.Context, days int) (float64, error) {
    // In a real system, you would use historical portfolio values recorded over time
    // For this simplified implementation, we'll estimate based on trade history
    
    // Get all trades in the period
    startTime := time.Now().AddDate(0, 0, -days)
    trades, err := s.boughtCoinRepo.FindAllBetween(ctx, startTime, time.Now())
    if err != nil {
        return 0, fmt.Errorf("failed to get trades: %w", err)
    }
    
    // Calculate maximum drawdown
    maxDrawdown := 0.0
    
    for _, trade := range trades {
        // If the trade was a loss
        pnl := 0.0
        if trade.CurrentPrice > 0 {
            pnl = (trade.CurrentPrice - trade.PurchasePrice) / trade.PurchasePrice * 100
        }
        
        if pnl < 0 && math.Abs(pnl) > maxDrawdown {
            maxDrawdown = math.Abs(pnl)
        }
    }
    
    return maxDrawdown, nil
}

// CalculateRiskRewardRatio calculates the ratio of average profit to average loss
func (s *Service) CalculateRiskRewardRatio(ctx context.Context) (float64, error) {
    performance, err := s.GetTradePerformance(ctx, "all")
    if err != nil {
        return 0, err
    }
    
    if performance.LargestLoss == 0 {
        return math.Inf(1), nil // No losses, return infinity
    }
    
    riskRewardRatio := math.Abs(performance.LargestProfit / performance.LargestLoss)
    
    return riskRewardRatio, nil
}

// GetTopPerformingCoins returns the coins with the highest profit percentage
func (s *Service) GetTopPerformingCoins(ctx context.Context, limit int) ([]*models.BoughtCoin, error) {
    // Get all completed trades
    completedTrades, err := s.boughtCoinRepo.FindAllCompleted(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get completed trades: %w", err)
    }
    
    // Sort trades by profit percentage (this would typically be done in the repository)
    // For example, using a custom sorting function or SQL ORDER BY
    
    // For illustration, let's assume we have a sorted list
    if len(completedTrades) > limit {
        completedTrades = completedTrades[:limit]
    }
    
    return completedTrades, nil
}
```

## Unit Testing

```go
// internal/core/portfolio/service_test.go
package portfolio

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

// ... other required methods ...

type mockBoughtCoinRepo struct {
    activeCoins     []*models.BoughtCoin
    completedCoins  []*models.BoughtCoin
    allCoins        []*models.BoughtCoin
}

func (m *mockBoughtCoinRepo) FindActive(ctx context.Context) ([]*models.BoughtCoin, error) {
    return m.activeCoins, nil
}

func (m *mockBoughtCoinRepo) FindCompletedBetween(ctx context.Context, start, end time.Time) ([]*models.BoughtCoin, error) {
    // In a real implementation, you would filter by time range
    return m.completedCoins, nil
}

func (m *mockBoughtCoinRepo) FindAllBetween(ctx context.Context, start, end time.Time) ([]*models.BoughtCoin, error) {
    // In a real implementation, you would filter by time range
    return m.allCoins, nil
}

func (m *mockBoughtCoinRepo) FindAllCompleted(ctx context.Context) ([]*models.BoughtCoin, error) {
    return m.completedCoins, nil
}

// ... other required methods ...

func TestGetPortfolioValue(t *testing.T) {
    // Set up test data
    mockExchange := &mockExchangeService{
        tickers: map[string]*models.Ticker{
            "BTC/USDT": {
                Symbol: "BTC/USDT",
                Price:  50000.0,
            },
            "ETH/USDT": {
                Symbol: "ETH/USDT",
                Price:  3000.0,
            },
        },
        wallet: &models.Wallet{
            USDT: 1000.0,
            Assets: map[string]float64{
                "BTC": 0.1,
                "ETH": 1.0,
            },
        },
    }
    
    mockRepo := &mockBoughtCoinRepo{
        activeCoins: []*models.BoughtCoin{
            {
                Symbol:        "BTC/USDT",
                Quantity:      0.1,
                PurchasePrice: 45000.0,
            },
            {
                Symbol:        "ETH/USDT",
                Quantity:      1.0,
                PurchasePrice: 2800.0,
            },
        },
    }
    
    // Create service with mocks
    service := NewService(mockExchange, mockRepo)
    
    // Test GetPortfolioValue
    ctx := context.Background()
    value, err := service.GetPortfolioValue(ctx)
    
    if err != nil {
        t.Fatalf("GetPortfolioValue returned error: %v", err)
    }
    
    // Expected value: 1000 (USDT) + 0.1*50000 (BTC) + 1.0*3000 (ETH) = 1000 + 5000 + 3000 = 9000
    expectedValue := 9000.0
    if value != expectedValue {
        t.Errorf("Expected portfolio value %f, got %f", expectedValue, value)
    }
}

func TestGetTradePerformance(t *testing.T) {
    // Set up test data with a mix of profitable and losing trades
    now := time.Now()
    yesterday := now.AddDate(0, 0, -1)
    lastWeek := now.AddDate(0, 0, -7)
    
    mockExchange := &mockExchangeService{
        wallet: &models.Wallet{
            USDT: 1200.0, // Increased from initial 1000
        },
    }
    
    mockRepo := &mockBoughtCoinRepo{
        completedCoins: []*models.BoughtCoin{
            {
                Symbol:        "BTC/USDT",
                Quantity:      0.1,
                PurchasePrice: 45000.0,
                CurrentPrice:  50000.0, // 11.11% profit
                PurchaseTime:  lastWeek,
                SellTime:      yesterday,
                IsDeleted:     true,
            },
            {
                Symbol:        "ETH/USDT",
                Quantity:      1.0,
                PurchasePrice: 3000.0,
                CurrentPrice:  2700.0, // 10% loss
                PurchaseTime:  lastWeek,
                SellTime:      yesterday,
                IsDeleted:     true,
            },
        },
    }
    
    // Create service with mocks
    service := NewService(mockExchange, mockRepo)
    
    // Test GetTradePerformance for the week
    ctx := context.Background()
    performance, err := service.GetTradePerformance(ctx, "week")
    
    if err != nil {
        t.Fatalf("GetTradePerformance returned error: %v", err)
    }
    
    // Verify the metrics
    if performance.TotalTrades != 2 {
        t.Errorf("Expected 2 total trades, got %d", performance.TotalTrades)
    }
    
    if performance.WinningTrades != 1 {
        t.Errorf("Expected 1 winning trade, got %d", performance.WinningTrades)
    }
    
    if performance.LosingTrades != 1 {
        t.Errorf("Expected 1 losing trade, got %d", performance.LosingTrades)
    }
    
    // Expected win rate: 1/2 = 50%
    expectedWinRate := 50.0
    if math.Abs(performance.WinRate-expectedWinRate) > 0.01 {
        t.Errorf("Expected win rate %.2f%%, got %.2f%%", expectedWinRate, performance.WinRate)
    }
    
    // Expected average profit: (11.11 - 10) / 2 = 0.56%
    expectedAvgProfit := 0.56
    if math.Abs(performance.AverageProfit-expectedAvgProfit) > 0.1 {
        t.Errorf("Expected average profit %.2f%%, got %.2f%%", expectedAvgProfit, performance.AverageProfit)
    }
    
    // Expected largest profit: 11.11%
    expectedLargestProfit := 11.11
    if math.Abs(performance.LargestProfit-expectedLargestProfit) > 0.1 {
        t.Errorf("Expected largest profit %.2f%%, got %.2f%%", expectedLargestProfit, performance.LargestProfit)
    }
    
    // Expected largest loss: -10%
    expectedLargestLoss := -10.0
    if math.Abs(performance.LargestLoss-expectedLargestLoss) > 0.1 {
        t.Errorf("Expected largest loss %.2f%%, got %.2f%%", expectedLargestLoss, performance.LargestLoss)
    }
}
```

## Integration with API and Frontend

The portfolio service is designed to provide data to both the API layer and directly to the frontend for visualization. Here's how it might be integrated:

### API Integration

```go
// internal/api/handlers/portfolio.go
package handlers

import (
    "net/http"
    
    " "
    "github.com/ryanlisse/cryptobot/internal/domain/service"
)

// PortfolioHandler handles portfolio-related API requests
type PortfolioHandler struct {
    portfolioService service.PortfolioService
}

// NewPortfolioHandler creates a new portfolio handler
func NewPortfolioHandler(portfolioService service.PortfolioService) *PortfolioHandler {
    return &PortfolioHandler{
        portfolioService: portfolioService,
    }
}

// GetPortfolioSummary returns a summary of the portfolio
func (h *PortfolioHandler) GetPortfolioSummary(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Get current portfolio value
    value, err := h.portfolioService.GetPortfolioValue(ctx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Get active trades
    activeTrades, err := h.portfolioService.GetActiveTrades(ctx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Get performance metrics
    timeRange := c.DefaultQuery("timeRange", "week")
    performance, err := h.portfolioService.GetTradePerformance(ctx, timeRange)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Compile the response
    response := gin.H{
        "portfolio_value": value,
        "active_trades":   activeTrades,
        "performance":     performance,
    }
    
    c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers the portfolio handlers
func (h *PortfolioHandler) RegisterRoutes(router *gin.RouterGroup) {
    portfolioGroup := router.Group("/portfolio")
    {
        portfolioGroup.GET("", h.GetPortfolioSummary)
        // Additional portfolio routes can be added here
    }
}
```

### Frontend Example (Next.js)

```tsx
// nextjs-frontend/pages/dashboard.tsx
import { useEffect, useState } from 'react';
import { PortfolioSummary } from '../types/portfolio';
import { PortfolioValueChart } from '../components/PortfolioValueChart';
import { ActiveTradesTable } from '../components/ActiveTradesTable';
import { PerformanceMetricsCard } from '../components/PerformanceMetricsCard';

export default function Dashboard() {
  const [portfolioData, setPortfolioData] = useState<PortfolioSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [timeRange, setTimeRange] = useState('week');
  
  useEffect(() => {
    async function fetchPortfolioData() {
      setLoading(true);
      try {
        const response = await fetch(`/api/portfolio?timeRange=${timeRange}`);
        const data = await response.json();
        setPortfolioData(data);
      } catch (error) {
        console.error('Failed to fetch portfolio data:', error);
      } finally {
        setLoading(false);
      }
    }
    
    fetchPortfolioData();
    
    // Set up polling for real-time updates
    const intervalId = setInterval(fetchPortfolioData, 30000); // 30 seconds
    
    return () => clearInterval(intervalId);
  }, [timeRange]);
  
  if (loading) {
    return <div>Loading portfolio data...</div>;
  }
  
  if (!portfolioData) {
    return <div>No portfolio data available</div>;
  }
  
  return (
    <div className="dashboard">
      <h1>Portfolio Dashboard</h1>
      
      <div className="portfolio-value">
        <h2>Portfolio Value: ${portfolioData.portfolio_value.toFixed(2)}</h2>
        <PortfolioValueChart data={portfolioData.performance} />
      </div>
      
      <div className="time-range-selector">
        <button onClick={() => setTimeRange('day')} className={timeRange === 'day' ? 'active' : ''}>
          Day
        </button>
        <button onClick={() => setTimeRange('week')} className={timeRange === 'week' ? 'active' : ''}>
          Week
        </button>
        <button onClick={() => setTimeRange('month')} className={timeRange === 'month' ? 'active' : ''}>
          Month
        </button>
        <button onClick={() => setTimeRange('year')} className={timeRange === 'year' ? 'active' : ''}>
          Year
        </button>
        <button onClick={() => setTimeRange('all')} className={timeRange === 'all' ? 'active' : ''}>
          All Time
        </button>
      </div>
      
      <div className="performance-metrics">
        <PerformanceMetricsCard performance={portfolioData.performance} />
      </div>
      
      <div className="active-trades">
        <h2>Active Trades</h2>
        <ActiveTradesTable trades={portfolioData.active_trades} />
      </div>
    </div>
  );
}
```

The PortfolioService provides comprehensive functionality for tracking, analyzing, and visualizing the performance of the trading bot. It serves as a key component for monitoring the effectiveness of the trading strategies and helping users make informed decisions about their cryptocurrency investments.
