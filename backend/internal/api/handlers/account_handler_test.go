package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// MockAccountService extends the mockExchangeService to add account-specific testing functions
type MockAccountService struct {
	mockExchangeService
	MockWallet   *models.Wallet
	ShouldErr    bool
	ErrorMessage string
}

// AnalyzeTransactions implements the account.AccountService interface
func (m *MockAccountService) AnalyzeTransactions(ctx context.Context, startTime, endTime time.Time) (*models.TransactionAnalysis, error) {
	if m.ShouldErr {
		return nil, errors.New(m.getErrorMessage())
	}
	return &models.TransactionAnalysis{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// GetAccountBalance implements the account.AccountService interface
func (m *MockAccountService) GetAccountBalance(ctx context.Context) (models.Balance, error) {
	if m.ShouldErr {
		return models.Balance{}, errors.New(m.getErrorMessage())
	}
	return models.Balance{
		Fiat: 1000.0,
		Available: map[string]float64{
			"BTC": 0.1,
			"ETH": 1.0,
		},
	}, nil
}

// Helper method to get error message
func (m *MockAccountService) getErrorMessage() string {
	if m.ErrorMessage != "" {
		return m.ErrorMessage
	}
	return "mock error"
}

// GetAllPositionRisks implements the account.AccountService interface
func (m *MockAccountService) GetAllPositionRisks(ctx context.Context) (map[string]models.PositionRisk, error) {
	if m.ShouldErr {
		return nil, errors.New(m.getErrorMessage())
	}
	return map[string]models.PositionRisk{}, nil
}

// GetPositionRisk implements the account.AccountService interface
func (m *MockAccountService) GetPositionRisk(ctx context.Context, symbol string) (models.PositionRisk, error) {
	if m.ShouldErr {
		return models.PositionRisk{}, errors.New(m.getErrorMessage())
	}
	return models.PositionRisk{}, nil
}

// GetCurrentExposure implements the account.AccountService interface
func (m *MockAccountService) GetCurrentExposure(ctx context.Context) (float64, error) {
	if m.ShouldErr {
		return 0, errors.New(m.getErrorMessage())
	}
	return 1000.0, nil
}

// GetPortfolioValue implements the account.AccountService interface
func (m *MockAccountService) GetPortfolioValue(ctx context.Context) (float64, error) {
	if m.ShouldErr {
		return 0, errors.New(m.getErrorMessage())
	}
	return 5000.0, nil
}

// ValidateAPIKeys implements the account.AccountService interface
func (m *MockAccountService) ValidateAPIKeys(ctx context.Context) (bool, error) {
	if m.ShouldErr {
		return false, errors.New(m.getErrorMessage())
	}
	return true, nil
}

// UpdateBalance implements the account.AccountService interface
func (m *MockAccountService) UpdateBalance(ctx context.Context, amount float64, reason string) error {
	if m.ShouldErr {
		return errors.New(m.getErrorMessage())
	}
	return nil
}

// SyncWithExchange implements the account.AccountService interface
func (m *MockAccountService) SyncWithExchange(ctx context.Context) error {
	if m.ShouldErr {
		return errors.New(m.getErrorMessage())
	}
	return nil
}

// GetBalanceSummary implements the account.AccountService interface
func (m *MockAccountService) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
	if m.ShouldErr {
		return nil, errors.New(m.getErrorMessage())
	}
	return &models.BalanceSummary{}, nil
}

// GetTransactionHistory implements the account.AccountService interface
func (m *MockAccountService) GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	if m.ShouldErr {
		return nil, errors.New(m.getErrorMessage())
	}
	return []*models.Transaction{}, nil
}

// SubscribeToBalanceUpdates implements the account.AccountService interface
func (m *MockAccountService) SubscribeToBalanceUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	if m.ShouldErr {
		return errors.New(m.getErrorMessage())
	}
	return nil
}

// GetListenKey implements the account.AccountService interface
func (m *MockAccountService) GetListenKey(ctx context.Context) (string, error) {
	if m.ShouldErr {
		return "", errors.New(m.getErrorMessage())
	}
	return "test-listen-key", nil
}

// RenewListenKey implements the account.AccountService interface
func (m *MockAccountService) RenewListenKey(ctx context.Context, listenKey string) error {
	if m.ShouldErr {
		return errors.New(m.getErrorMessage())
	}
	return nil
}

// CloseListenKey implements the account.AccountService interface
func (m *MockAccountService) CloseListenKey(ctx context.Context, listenKey string) error {
	if m.ShouldErr {
		return errors.New(m.getErrorMessage())
	}
	return nil
}

// Override GetWallet to return test-specific values
func (m *MockAccountService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	if m.ShouldErr {
		return nil, errors.New("mock service error")
	}
	return m.MockWallet, nil
}

func TestGetAccount(t *testing.T) {

	// Define a wallet struct for testing
	testWallet := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"BTC": {
				Asset:  "BTC",
				Free:   0.5,
				Locked: 0.0,
				Total:  0.5,
			},
			"ETH": {
				Asset:  "ETH",
				Free:   10.0,
				Locked: 0.0,
				Total:  10.0,
			},
		},
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name              string
		mockWallet        *models.Wallet
		mockError         bool
		expectedStatus    int
		expectedErrorCode string
	}{
		{
			name:           "returns account successfully",
			mockWallet:     testWallet,
			mockError:      false,
			expectedStatus: http.StatusOK,
		},
		{
			name:              "handles service error",
			mockWallet:        nil,
			mockError:         true,
			expectedStatus:    http.StatusInternalServerError,
			expectedErrorCode: "GET_ACCOUNT_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock service
			mockService := &MockAccountService{
				MockWallet:   tt.mockWallet,
				ShouldErr:    tt.mockError,
				ErrorMessage: "mock service error",
			}

			// Create handler with mock service
			logger, _ := zap.NewDevelopment()
			handler := NewAccountHandler(nil, mockService, logger)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Perform request
			handler.GetAccount(w, req)

			// Assert response status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.mockError {
				// For successful response, check wallet data is in response
				assert.Contains(t, w.Body.String(), "balances")
				assert.Contains(t, w.Body.String(), "BTC")
				assert.Contains(t, w.Body.String(), "ETH")
				assert.Contains(t, w.Body.String(), "updatedAt")
			} else {
				// For error response, check error details
				assert.Contains(t, w.Body.String(), tt.expectedErrorCode)
				assert.Contains(t, w.Body.String(), "Failed to get account info")
				assert.Contains(t, w.Body.String(), "mock service error")
			}
		})
	}
}

func TestGetBalances(t *testing.T) {

	// Define a wallet struct for testing
	testWallet := &models.Wallet{
		Balances: map[string]*models.AssetBalance{
			"BTC": {
				Asset:  "BTC",
				Free:   0.5,
				Locked: 0.0,
				Total:  0.5,
			},
			"ETH": {
				Asset:  "ETH",
				Free:   10.0,
				Locked: 0.0,
				Total:  10.0,
			},
		},
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name              string
		mockWallet        *models.Wallet
		mockError         bool
		expectedStatus    int
		expectedErrorCode string
	}{
		{
			name:           "returns balances successfully",
			mockWallet:     testWallet,
			mockError:      false,
			expectedStatus: http.StatusOK,
		},
		{
			name:              "handles service error",
			mockWallet:        nil,
			mockError:         true,
			expectedStatus:    http.StatusInternalServerError,
			expectedErrorCode: "GET_BALANCES_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock service
			mockService := &MockAccountService{
				MockWallet:   tt.mockWallet,
				ShouldErr:    tt.mockError,
				ErrorMessage: "mock service error",
			}

			// Create handler with mock service
			logger, _ := zap.NewDevelopment()
			handler := NewAccountHandler(nil, mockService, logger)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Perform request
			handler.GetBalances(w, req)

			// Assert response status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.mockError {
				// Should include individual balances
				assert.Contains(t, w.Body.String(), "BTC")
				assert.Contains(t, w.Body.String(), "ETH")
				// Verify we just got balances, not the full wallet
				assert.NotContains(t, w.Body.String(), "updatedAt")
			} else {
				// For error response, check error details
				assert.Contains(t, w.Body.String(), tt.expectedErrorCode)
				assert.Contains(t, w.Body.String(), "Failed to get account balances")
				assert.Contains(t, w.Body.String(), "mock service error")
			}
		})
	}
}
