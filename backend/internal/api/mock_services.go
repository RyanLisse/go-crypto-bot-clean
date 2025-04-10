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
