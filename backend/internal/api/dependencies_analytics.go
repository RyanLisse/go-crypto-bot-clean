package api

import (
	"context"
	"math/rand"
	"time"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"

	"go.uber.org/zap"
)

// InitializeAnalyticsDependencies initializes the Analytics dependencies
func (d *Dependencies) InitializeAnalyticsDependencies() {
	// Use the logger from the dependencies
	logger := d.logger

	// Check if we're in development mode
	if d.Config.App.Environment == "development" {
		logger.Info("Using mock Analytics service for development mode")
		// Use our MockAnalyticsService
		mockService := &MockAnalyticsService{}
		d.AnalyticsService = mockService
		d.AnalyticsHandler = handlers.NewAnalyticsHandler(mockService, logger)
		return
	}

	// Create MEXC client for real data
	mexcClient, err := rest.NewClient(d.Config.Mexc.APIKey, d.Config.Mexc.SecretKey)
	if err != nil {
		logger.Error("Failed to create MEXC client, falling back to mock service", zap.Error(err))
		// Fall back to mock service
		mockService := &mockAnalyticsService{}
		d.AnalyticsService = mockService
		d.AnalyticsHandler = handlers.NewAnalyticsHandler(mockService, logger)
		return
	}

	// Create real analytics service adapter
	mockService := &mockAnalyticsService{}
	analyticsAdapter := NewRealAnalyticsServiceAdapter(mockService, mexcClient, logger)
	d.AnalyticsService = analyticsAdapter
	d.AnalyticsHandler = handlers.NewAnalyticsHandler(analyticsAdapter, logger)
}

// mockAnalyticsService is a mock implementation of the analytics.TradeAnalyticsService interface
type mockAnalyticsService struct{}

// GetTradeAnalytics returns mock trade analytics
func (s *mockAnalyticsService) GetTradeAnalytics(ctx context.Context, timeFrame models.TimeFrame, startTime, endTime time.Time) (*models.TradeAnalytics, error) {
	return &models.TradeAnalytics{
		TimeFrame:             timeFrame,
		StartTime:             startTime,
		EndTime:               endTime,
		TotalTrades:           10,
		WinningTrades:         7,
		LosingTrades:          3,
		WinRate:               0.7,
		TotalProfit:           100.0,
		TotalLoss:             30.0,
		NetProfit:             70.0,
		ProfitFactor:          3.33,
		LargestProfit:         30.0,
		LargestLoss:           15.0,
		AverageProfit:         14.28,
		AverageLoss:           10.0,
		AverageHoldingTime:    "2h 30m",
		PerformanceBySymbol:   make(map[string]models.SymbolPerformance),
		PerformanceByReason:   make(map[string]models.ReasonPerformance),
		PerformanceByStrategy: make(map[string]models.StrategyPerformance),
	}, nil
}

// GetTradePerformance returns mock trade performance
func (s *mockAnalyticsService) GetTradePerformance(ctx context.Context, tradeID string) (*models.TradePerformance, error) {
	return &models.TradePerformance{
		Symbol:      "BTC/USDT",
		EntryPrice:  20000.0,
		ExitPrice:   22000.0,
		Quantity:    0.1,
		ProfitLoss:  200.0,
		EntryTime:   time.Now().Add(-24 * time.Hour),
		ExitTime:    time.Now(),
		HoldingTime: "24h 0m 0s",
		ExitReason:  "take_profit",
		Strategy:    "newcoin",
	}, nil
}

// GetAllTradePerformance returns mock trade performances
func (s *mockAnalyticsService) GetAllTradePerformance(ctx context.Context, startTime, endTime time.Time) ([]*models.TradePerformance, error) {
	return []*models.TradePerformance{
		{
			Symbol:      "BTC/USDT",
			EntryPrice:  20000.0,
			ExitPrice:   22000.0,
			Quantity:    0.1,
			ProfitLoss:  200.0,
			EntryTime:   time.Now().Add(-48 * time.Hour),
			ExitTime:    time.Now().Add(-24 * time.Hour),
			HoldingTime: "24h 0m 0s",
			ExitReason:  "take_profit",
			Strategy:    "newcoin",
		},
		{
			Symbol:      "ETH/USDT",
			EntryPrice:  1500.0,
			ExitPrice:   1650.0,
			Quantity:    1.0,
			ProfitLoss:  150.0,
			EntryTime:   time.Now().Add(-24 * time.Hour),
			ExitTime:    time.Now(),
			HoldingTime: "24h 0m 0s",
			ExitReason:  "take_profit",
			Strategy:    "newcoin",
		},
	}, nil
}

// GetWinRate returns a mock win rate
func (s *mockAnalyticsService) GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	return 0.7, nil
}

// GetProfitFactor returns a mock profit factor
func (s *mockAnalyticsService) GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	return 3.33, nil
}

// GetDrawdown returns mock drawdown values
func (s *mockAnalyticsService) GetDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, float64, error) {
	return 15.0, 5.0, nil
}

// GetBalanceHistory returns mock balance history
func (s *mockAnalyticsService) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time, interval time.Duration) ([]models.BalancePoint, error) {
	// Generate mock balance points for the given time range
	var points []models.BalancePoint
	currentTime := startTime
	balance := 1000.0 // Starting balance

	for currentTime.Before(endTime) {
		// Add some random variation to the balance
		balanceChange := (rand.Float64() * 100) - 30 // Random change between -30 and +70
		balance += balanceChange

		points = append(points, models.BalancePoint{
			Timestamp: currentTime,
			Balance:   balance,
		})

		currentTime = currentTime.Add(interval)
	}

	return points, nil
}

// GetPerformanceBySymbol returns mock performance by symbol
func (s *mockAnalyticsService) GetPerformanceBySymbol(ctx context.Context, startTime, endTime time.Time) (map[string]models.SymbolPerformance, error) {
	return map[string]models.SymbolPerformance{
		"BTC/USDT": {
			Symbol:        "BTC/USDT",
			TotalTrades:   5,
			WinningTrades: 4,
			LosingTrades:  1,
			TotalProfit:   400.0,
			WinRate:       0.8,
		},
		"ETH/USDT": {
			Symbol:        "ETH/USDT",
			TotalTrades:   5,
			WinningTrades: 3,
			LosingTrades:  2,
			TotalProfit:   250.0,
			WinRate:       0.6,
		},
	}, nil
}

// GetPerformanceByReason returns mock performance by reason
func (s *mockAnalyticsService) GetPerformanceByReason(ctx context.Context, startTime, endTime time.Time) (map[string]models.ReasonPerformance, error) {
	return map[string]models.ReasonPerformance{
		"take_profit": {
			Reason:        "take_profit",
			TotalTrades:   7,
			WinningTrades: 7,
			LosingTrades:  0,
			TotalProfit:   700.0,
			WinRate:       1.0,
		},
		"stop_loss": {
			Reason:        "stop_loss",
			TotalTrades:   3,
			WinningTrades: 0,
			LosingTrades:  3,
			TotalProfit:   -150.0,
			WinRate:       0.0,
		},
	}, nil
}

// GetPerformanceByStrategy returns mock performance by strategy
func (s *mockAnalyticsService) GetPerformanceByStrategy(ctx context.Context, startTime, endTime time.Time) (map[string]models.StrategyPerformance, error) {
	return map[string]models.StrategyPerformance{
		"newcoin": {
			Strategy:      "newcoin",
			TotalTrades:   8,
			WinningTrades: 6,
			LosingTrades:  2,
			TotalProfit:   500.0,
			WinRate:       0.75,
		},
		"trend_following": {
			Strategy:      "trend_following",
			TotalTrades:   2,
			WinningTrades: 1,
			LosingTrades:  1,
			TotalProfit:   50.0,
			WinRate:       0.5,
		},
	}, nil
}
