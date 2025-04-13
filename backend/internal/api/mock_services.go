package api

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/core/status" // Ensure this is imported
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// MockStatusService is a mock implementation of the status service
type MockStatusService struct{}

// Mock methods to satisfy the interface expected by handlers.StatusHandler
func (m *MockStatusService) GetStatus() (*status.SystemStatus, error) {
	// Return a default status or nil, nil for mock purposes
	return &status.SystemStatus{OverallStatus: "mocked"}, nil
}

// Add context parameter to match handler usage (even if concrete type differs)
func (m *MockStatusService) StartProcesses(ctx context.Context) error {
	// No-op for mock
	return nil
}

// Signature matches handler usage
func (m *MockStatusService) StopProcesses() error {
	// No-op for mock
	return nil
}

// Removed duplicate method definitions below

// MockPortfolioService is a mock implementation of the portfolio service
type MockPortfolioService struct{}

// MockAccountService is a mock implementation of the account service
type MockAccountService struct{}

// GetPortfolioValue returns a mock portfolio value
func (m *MockPortfolioService) GetPortfolioValue(ctx context.Context) (float64, error) {
	return 10000.0, nil // Mock value of 10,000 USDT
}

// GetActiveTrades returns mock active trades
func (m *MockPortfolioService) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
	// Create some mock active trades
	return []*models.BoughtCoin{
		{
			ID:            1,
			Symbol:        "BTCUSDT",
			PurchasePrice: 75000.0,
			CurrentPrice:  79000.0,
			Quantity:      0.1,
			BoughtAt:      time.Now().Add(-24 * time.Hour),
			StopLoss:      70000.0,
			TakeProfit:    85000.0,
		},
		{
			ID:            2,
			Symbol:        "ETHUSDT",
			PurchasePrice: 3500.0,
			CurrentPrice:  3800.0,
			Quantity:      1.0,
			BoughtAt:      time.Now().Add(-48 * time.Hour),
			StopLoss:      3200.0,
			TakeProfit:    4000.0,
		},
	}, nil
}

// GetTradePerformance returns mock trade performance metrics
func (m *MockPortfolioService) GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error) {
	return &models.PerformanceMetrics{
		TotalTrades:           10,
		WinningTrades:         7,
		LosingTrades:          3,
		WinRate:               70.0,
		TotalProfitLoss:       1500.0,
		AverageProfitPerTrade: 150.0,
		LargestProfit:         500.0,
		LargestLoss:           -200.0,
	}, nil
}

// GetPositions returns mock positions
func (m *MockPortfolioService) GetPositions(ctx context.Context) ([]models.Position, error) {
	return []models.Position{
		{
			Symbol:        "BTCUSDT",
			Quantity:      0.1,
			EntryPrice:    75000.0,
			CurrentPrice:  79000.0,
			PnL:           400.0,
			PnLPercentage: 5.33,
			OpenTime:      time.Now().Add(-24 * time.Hour),
			Status:        models.PositionStatusOpen,
			Side:          models.OrderSideBuy,
			StopLoss:      70000.0,
			TakeProfit:    85000.0,
		},
		{
			Symbol:        "ETHUSDT",
			Quantity:      1.0,
			EntryPrice:    3500.0,
			CurrentPrice:  3800.0,
			PnL:           300.0,
			PnLPercentage: 8.57,
			OpenTime:      time.Now().Add(-48 * time.Hour),
			Status:        models.PositionStatusOpen,
			Side:          models.OrderSideBuy,
			StopLoss:      3200.0,
			TakeProfit:    4000.0,
		},
	}, nil
}

// GetAccountBalance returns a mock account balance
func (m *MockAccountService) GetAccountBalance(ctx context.Context) (models.Balance, error) {
	return models.Balance{
		Fiat: 10000.0,
		Assets: map[string]float64{
			"BTC": 0.1,
			"ETH": 1.0,
			"SOL": 10.0,
		},
		Available: map[string]float64{
			"BTC": 0.1,
			"ETH": 1.0,
			"SOL": 10.0,
		},
		Locked: map[string]float64{
			"BTC": 0.0,
			"ETH": 0.0,
			"SOL": 0.0,
		},
		UpdatedAt: time.Now(),
	}, nil
}

// GetWallet returns a mock wallet
func (m *MockAccountService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	return &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"USDT": {
				Asset:  "USDT",
				Free:   10000.0,
				Locked: 0.0,
				Total:  10000.0,
			},
			"BTC": {
				Asset:  "BTC",
				Free:   0.1,
				Locked: 0.0,
				Total:  0.1,
			},
			"ETH": {
				Asset:  "ETH",
				Free:   1.0,
				Locked: 0.0,
				Total:  1.0,
			},
			"SOL": {
				Asset:  "SOL",
				Free:   10.0,
				Locked: 0.0,
				Total:  10.0,
			},
		},
		UpdatedAt: time.Now(),
	}, nil
}

// ValidateAPIKeys validates API keys
func (m *MockAccountService) ValidateAPIKeys(ctx context.Context) (bool, error) {
	return true, nil
}

// GetCurrentExposure returns the current exposure
func (m *MockAccountService) GetCurrentExposure(ctx context.Context) (float64, error) {
	return 5000.0, nil
}

// GetListenKey returns a mock listen key for WebSocket authentication
func (m *MockAccountService) GetListenKey(ctx context.Context) (string, error) {
	return "mock-listen-key-12345", nil
}

// RenewListenKey renews a listen key
func (m *MockAccountService) RenewListenKey(ctx context.Context, listenKey string) error {
	return nil
}

// CloseListenKey closes a listen key
func (m *MockAccountService) CloseListenKey(ctx context.Context, listenKey string) error {
	return nil
}

// AnalyzeTransactions analyzes the transaction history
func (m *MockAccountService) AnalyzeTransactions(ctx context.Context, startTime, endTime time.Time) (*models.TransactionAnalysis, error) {
	return &models.TransactionAnalysis{
		StartTime:   startTime,
		EndTime:     endTime,
		TotalCount:  10,
		BuyCount:    7,
		SellCount:   3,
		TotalVolume: 1500.0,
		BuyVolume:   1000.0,
		SellVolume:  500.0,
	}, nil
}

// GetBalanceSummary returns a mock balance summary
func (m *MockAccountService) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
	return &models.BalanceSummary{
		CurrentBalance:   10000.0,
		Deposits:         1000.0,
		Withdrawals:      500.0,
		NetChange:        500.0,
		TransactionCount: 10,
		Period:           days,
		GeneratedAt:      time.Now(),
	}, nil
}

// GetTransactionHistory returns mock transaction history
func (m *MockAccountService) GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	return []*models.Transaction{
		{
			ID:        "1",
			Amount:    1000.0,
			Balance:   10000.0,
			Reason:    "deposit",
			Timestamp: time.Now().Add(-48 * time.Hour),
		},
		{
			ID:        "2",
			Amount:    -500.0,
			Balance:   9500.0,
			Reason:    "withdrawal",
			Timestamp: time.Now().Add(-24 * time.Hour),
		},
	}, nil
}

// GetPositionRisk returns mock position risk
func (m *MockAccountService) GetPositionRisk(ctx context.Context, symbol string) (models.PositionRisk, error) {
	return models.PositionRisk{
		Symbol:      symbol,
		ExposureUSD: 5000.0,
		RiskLevel:   "LOW",
	}, nil
}

// GetAllPositionRisks returns all mock position risks
func (m *MockAccountService) GetAllPositionRisks(ctx context.Context) (map[string]models.PositionRisk, error) {
	return map[string]models.PositionRisk{
		"BTCUSDT": {
			Symbol:      "BTCUSDT",
			ExposureUSD: 5000.0,
			RiskLevel:   "LOW",
		},
		"ETHUSDT": {
			Symbol:      "ETHUSDT",
			ExposureUSD: 3000.0,
			RiskLevel:   "MEDIUM",
		},
	}, nil
}

// SubscribeToBalanceUpdates subscribes to balance updates
func (m *MockAccountService) SubscribeToBalanceUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	// Mock implementation just returns nil
	return nil
}

// UpdateBalance updates the balance
func (m *MockAccountService) UpdateBalance(ctx context.Context, amount float64, reason string) error {
	// Mock implementation just returns nil
	return nil
}

// SyncWithExchange syncs with the exchange
func (m *MockAccountService) SyncWithExchange(ctx context.Context) error {
	// Mock implementation just returns nil
	return nil
}

// GetPortfolioValue returns the portfolio value
func (m *MockAccountService) GetPortfolioValue(ctx context.Context) (float64, error) {
	return 15000.0, nil
}

// MockNewCoinService is a mock implementation of the new coin service
type MockNewCoinService struct{}

// GetNewCoins returns mock new coins
func (m *MockNewCoinService) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	// Create time pointers
	foundAt1 := time.Now().Add(-24 * time.Hour)
	firstOpenTime1 := time.Now().Add(24 * time.Hour)
	foundAt2 := time.Now().Add(-48 * time.Hour)
	firstOpenTime2 := time.Now().Add(12 * time.Hour)

	return []*models.NewCoin{
		{
			ID:            1,
			Symbol:        "NEWBTC/USDT",
			FoundAt:       foundAt1,
			FirstOpenTime: &firstOpenTime1,
			QuoteVolume:   1000000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
		{
			ID:            2,
			Symbol:        "NEWETH/USDT",
			FoundAt:       foundAt2,
			FirstOpenTime: &firstOpenTime2,
			QuoteVolume:   500000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
	}, nil
}

// GetUpcomingCoinsForTodayAndTomorrow returns mock upcoming coins for today and tomorrow
func (m *MockNewCoinService) GetUpcomingCoinsForTodayAndTomorrow(ctx context.Context) ([]*models.NewCoin, error) {
	// Create time pointers
	foundAt1 := time.Now().Add(-24 * time.Hour)
	firstOpenTime1 := time.Now().Add(12 * time.Hour)
	foundAt2 := time.Now().Add(-48 * time.Hour)
	firstOpenTime2 := time.Now().Add(24 * time.Hour)

	return []*models.NewCoin{
		{
			ID:            1,
			Symbol:        "NEWBTC/USDT",
			FoundAt:       foundAt1,
			FirstOpenTime: &firstOpenTime1,
			QuoteVolume:   1000000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
		{
			ID:            2,
			Symbol:        "NEWETH/USDT",
			FoundAt:       foundAt2,
			FirstOpenTime: &firstOpenTime2,
			QuoteVolume:   500000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
	}, nil
}

// GetUpcomingCoinsByDate returns mock upcoming coins by date
func (m *MockNewCoinService) GetUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]*models.NewCoin, error) {
	// Create time pointers
	foundAt1 := time.Now().Add(-24 * time.Hour)
	firstOpenTime1 := date
	foundAt2 := time.Now().Add(-48 * time.Hour)
	firstOpenTime2 := date.Add(12 * time.Hour)

	return []*models.NewCoin{
		{
			ID:            1,
			Symbol:        "NEWBTC/USDT",
			FoundAt:       foundAt1,
			FirstOpenTime: &firstOpenTime1,
			QuoteVolume:   1000000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
		{
			ID:            2,
			Symbol:        "NEWETH/USDT",
			FoundAt:       foundAt2,
			FirstOpenTime: &firstOpenTime2,
			QuoteVolume:   500000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
	}, nil
}

// MockAnalyticsService is a mock implementation of the analytics service
type MockAnalyticsService struct{}

// GetBalanceHistory returns mock balance history data
func (m *MockAnalyticsService) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time, interval time.Duration) ([]models.BalancePoint, error) {
	// Calculate number of points based on time range and interval
	duration := endTime.Sub(startTime)
	numPoints := int(duration/interval) + 1
	if numPoints > 100 {
		numPoints = 100 // Cap at 100 points to avoid excessive data
	}

	// Create mock balance history entries
	history := make([]models.BalancePoint, numPoints)
	baseBalance := 10000.0

	for i := 0; i < numPoints; i++ {
		// Calculate timestamp for this entry
		timestamp := startTime.Add(time.Duration(i) * interval)

		// Add some variation to the balance (simple upward trend with some noise)
		progress := float64(i) / float64(numPoints-1) // 0.0 to 1.0
		variation := progress * 2000.0                // Up to $2000 increase

		// Add some randomness
		random := (float64(i%5) - 2.0) * 100.0 // -200 to +200 random noise

		balance := baseBalance + variation + random

		history[i] = models.BalancePoint{
			Timestamp: timestamp,
			Balance:   balance,
		}
	}

	return history, nil
}

// GetTradePerformance returns mock trade performance data
func (m *MockAnalyticsService) GetTradePerformance(ctx context.Context, tradeID string) (*models.TradePerformance, error) {
	return &models.TradePerformance{
		TradeID:           tradeID,
		Symbol:            "BTC/USDT",
		EntryPrice:        30000.0,
		ExitPrice:         32000.0,
		Quantity:          0.1,
		ProfitLoss:        200.0,
		ProfitLossPercent: 6.67,
		EntryTime:         time.Now().Add(-48 * time.Hour),
		ExitTime:          time.Now().Add(-24 * time.Hour),
		HoldingTime:       "24h 0m 0s",
		HoldingTimeMs:     86400000,
		EntryReason:       "signal",
		ExitReason:        "take_profit",
		Strategy:          "trend_following",
	}, nil
}

// GetWinRate returns a mock win rate
func (m *MockAnalyticsService) GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	return 0.7, nil // 70% win rate
}

// GetProfitFactor returns a mock profit factor
func (m *MockAnalyticsService) GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	return 2.5, nil // 2.5 profit factor
}

// GetDrawdown returns mock drawdown values
func (m *MockAnalyticsService) GetDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, float64, error) {
	return 1500.0, 15.0, nil // $1500 max drawdown, 15% max drawdown percentage
}

// GetPerformanceBySymbol returns mock performance by symbol
func (m *MockAnalyticsService) GetPerformanceBySymbol(ctx context.Context, startTime, endTime time.Time) (map[string]models.SymbolPerformance, error) {
	return map[string]models.SymbolPerformance{
		"BTC/USDT": {
			Symbol:        "BTC/USDT",
			TotalTrades:   10,
			WinningTrades: 7,
			LosingTrades:  3,
			WinRate:       0.7,
			TotalProfit:   500.0,
			AverageProfit: 50.0,
			ProfitFactor:  2.5,
		},
		"ETH/USDT": {
			Symbol:        "ETH/USDT",
			TotalTrades:   8,
			WinningTrades: 5,
			LosingTrades:  3,
			WinRate:       0.625,
			TotalProfit:   300.0,
			AverageProfit: 37.5,
			ProfitFactor:  2.0,
		},
	}, nil
}

// GetPerformanceByReason returns mock performance by reason
func (m *MockAnalyticsService) GetPerformanceByReason(ctx context.Context, startTime, endTime time.Time) (map[string]models.ReasonPerformance, error) {
	return map[string]models.ReasonPerformance{
		"signal": {
			Reason:        "signal",
			TotalTrades:   12,
			WinningTrades: 8,
			LosingTrades:  4,
			WinRate:       0.67,
			TotalProfit:   600.0,
			AverageProfit: 50.0,
			ProfitFactor:  2.5,
		},
		"manual": {
			Reason:        "manual",
			TotalTrades:   6,
			WinningTrades: 4,
			LosingTrades:  2,
			WinRate:       0.67,
			TotalProfit:   200.0,
			AverageProfit: 33.3,
			ProfitFactor:  2.0,
		},
	}, nil
}

// GetPerformanceByStrategy returns mock performance by strategy
func (m *MockAnalyticsService) GetPerformanceByStrategy(ctx context.Context, startTime, endTime time.Time) (map[string]models.StrategyPerformance, error) {
	return map[string]models.StrategyPerformance{
		"trend_following": {
			Strategy:      "trend_following",
			TotalTrades:   10,
			WinningTrades: 7,
			LosingTrades:  3,
			WinRate:       0.7,
			TotalProfit:   500.0,
			AverageProfit: 50.0,
			ProfitFactor:  2.5,
		},
		"mean_reversion": {
			Strategy:      "mean_reversion",
			TotalTrades:   8,
			WinningTrades: 5,
			LosingTrades:  3,
			WinRate:       0.625,
			TotalProfit:   300.0,
			AverageProfit: 37.5,
			ProfitFactor:  2.0,
		},
	}, nil
}

// GetTradeAnalytics returns mock trade analytics data
func (m *MockAnalyticsService) GetTradeAnalytics(ctx context.Context, timeFrame models.TimeFrame, startTime, endTime time.Time) (*models.TradeAnalytics, error) {
	return &models.TradeAnalytics{
		TimeFrame:     timeFrame,
		StartTime:     startTime,
		EndTime:       endTime,
		TotalTrades:   20,
		WinningTrades: 14,
		LosingTrades:  6,
		WinRate:       70.0,
		TotalProfit:   1000.0,
		TotalLoss:     -300.0,
		NetProfit:     700.0,
		ProfitFactor:  3.33,
		AverageProfit: 71.43,
		AverageLoss:   -50.0,
		LargestProfit: 200.0,
		LargestLoss:   -100.0,
	}, nil
}

// GetAllTradePerformance returns mock trade performance data
func (m *MockAnalyticsService) GetAllTradePerformance(ctx context.Context, startTime, endTime time.Time) ([]*models.TradePerformance, error) {
	return []*models.TradePerformance{
		{
			TradeID:           "trade-1",
			Symbol:            "BTC/USDT",
			EntryPrice:        30000.0,
			ExitPrice:         32000.0,
			Quantity:          0.1,
			ProfitLoss:        200.0,
			ProfitLossPercent: 6.67,
			EntryTime:         time.Now().Add(-48 * time.Hour),
			ExitTime:          time.Now().Add(-24 * time.Hour),
			HoldingTime:       "24h 0m 0s",
			HoldingTimeMs:     86400000,
			EntryReason:       "signal",
			ExitReason:        "take_profit",
			Strategy:          "trend_following",
		},
		{
			TradeID:           "trade-2",
			Symbol:            "ETH/USDT",
			EntryPrice:        2000.0,
			ExitPrice:         1900.0,
			Quantity:          1.0,
			ProfitLoss:        -100.0,
			ProfitLossPercent: -5.0,
			EntryTime:         time.Now().Add(-24 * time.Hour),
			ExitTime:          time.Now().Add(-12 * time.Hour),
			HoldingTime:       "12h 0m 0s",
			HoldingTimeMs:     43200000,
			EntryReason:       "signal",
			ExitReason:        "stop_loss",
			Strategy:          "mean_reversion",
		},
	}, nil
}
