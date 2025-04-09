package risk

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/risk/controls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockBalanceHistoryRepository struct {
	mock.Mock
}

func (m *MockBalanceHistoryRepository) AddBalanceRecord(ctx context.Context, balance float64) error {
	args := m.Called(ctx, balance)
	return args.Error(0)
}

func (m *MockBalanceHistoryRepository) GetHistory(ctx context.Context, days int) ([]controls.BalanceHistory, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]controls.BalanceHistory), args.Error(1)
}

func (m *MockBalanceHistoryRepository) GetHighestBalance(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockBalanceHistoryRepository) GetBalanceAt(ctx context.Context, timestamp time.Time) (float64, error) {
	args := m.Called(ctx, timestamp)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockBalanceHistoryRepository) GetLatestBalance(ctx context.Context) (*controls.BalanceHistory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*controls.BalanceHistory), args.Error(1)
}

type MockPositionSizer struct {
	mock.Mock
}

func (m *MockPositionSizer) CalculatePositionSize(
	ctx context.Context,
	symbol string,
	accountBalance float64,
	riskPercent float64,
	stopLossPercent float64,
) (float64, error) {
	args := m.Called(ctx, symbol, accountBalance, riskPercent, stopLossPercent)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPositionSizer) CalculateOrderValue(
	ctx context.Context,
	symbol string,
	quantity float64,
) (float64, error) {
	args := m.Called(ctx, symbol, quantity)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPositionSizer) CalculateMaxQuantity(
	ctx context.Context,
	symbol string,
	maxOrderValue float64,
) (float64, error) {
	args := m.Called(ctx, symbol, maxOrderValue)
	return args.Get(0).(float64), args.Error(1)
}

type MockDrawdownMonitor struct {
	mock.Mock
}

func (m *MockDrawdownMonitor) CalculateDrawdown(ctx context.Context, days int) (float64, error) {
	args := m.Called(ctx, days)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDrawdownMonitor) CheckDrawdownLimit(ctx context.Context, maxDrawdownPercent float64) (bool, error) {
	args := m.Called(ctx, maxDrawdownPercent)
	return args.Bool(0), args.Error(1)
}

func (m *MockDrawdownMonitor) RecordBalance(ctx context.Context, balance float64) error {
	args := m.Called(ctx, balance)
	return args.Error(0)
}

type MockExposureMonitor struct {
	mock.Mock
	accountSvc controls.AccountService
}

func (m *MockExposureMonitor) GetAccountBalance(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExposureMonitor) CalculateTotalExposure(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExposureMonitor) CalculateExposurePercent(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExposureMonitor) CheckExposureLimit(
	ctx context.Context,
	newOrderValue float64,
	maxExposurePercent float64,
) (bool, error) {
	args := m.Called(ctx, newOrderValue, maxExposurePercent)
	return args.Bool(0), args.Error(1)
}

type MockDailyLimitMonitor struct {
	mock.Mock
}

func (m *MockDailyLimitMonitor) CalculateDailyPnL(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDailyLimitMonitor) CheckDailyLossLimit(
	ctx context.Context,
	dailyLossLimitPercent float64,
) (bool, error) {
	args := m.Called(ctx, dailyLossLimitPercent)
	return args.Bool(0), args.Error(1)
}

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetBalance(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

// Tests
func TestCalculatePositionSize(t *testing.T) {
	// Create mocks
	mockBalanceRepo := new(MockBalanceHistoryRepository)
	mockPositionSizer := new(MockPositionSizer)
	mockDrawdownMonitor := new(MockDrawdownMonitor)
	mockExposureMonitor := new(MockExposureMonitor)
	mockDailyLimitMonitor := new(MockDailyLimitMonitor)
	mockLogger := new(MockLogger)

	// Create risk manager
	riskManager := NewRiskManager(
		mockBalanceRepo,
		mockPositionSizer,
		mockDrawdownMonitor,
		mockExposureMonitor,
		mockDailyLimitMonitor,
		mockLogger,
	)

	// Set up expectations
	ctx := context.Background()
	symbol := "BTC/USDT"
	accountBalance := 1000.0
	expectedQuantity := 0.01

	mockPositionSizer.On("CalculatePositionSize", ctx, symbol, accountBalance, 1.0, 5.0).
		Return(expectedQuantity, nil)

	// Call the method
	quantity, err := riskManager.CalculatePositionSize(ctx, symbol, accountBalance)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedQuantity, quantity)
	mockPositionSizer.AssertExpectations(t)
}

func TestCalculateDrawdown(t *testing.T) {
	// Create mocks
	mockBalanceRepo := new(MockBalanceHistoryRepository)
	mockPositionSizer := new(MockPositionSizer)
	mockDrawdownMonitor := new(MockDrawdownMonitor)
	mockExposureMonitor := new(MockExposureMonitor)
	mockDailyLimitMonitor := new(MockDailyLimitMonitor)
	mockLogger := new(MockLogger)

	// Create risk manager
	riskManager := NewRiskManager(
		mockBalanceRepo,
		mockPositionSizer,
		mockDrawdownMonitor,
		mockExposureMonitor,
		mockDailyLimitMonitor,
		mockLogger,
	)

	// Set up expectations
	ctx := context.Background()
	expectedDrawdown := 0.15

	mockDrawdownMonitor.On("CalculateDrawdown", ctx, 90).
		Return(expectedDrawdown, nil)

	// Call the method
	drawdown, err := riskManager.CalculateDrawdown(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedDrawdown, drawdown)
	mockDrawdownMonitor.AssertExpectations(t)
}

func TestCheckExposureLimit(t *testing.T) {
	// Create mocks
	mockBalanceRepo := new(MockBalanceHistoryRepository)
	mockPositionSizer := new(MockPositionSizer)
	mockDrawdownMonitor := new(MockDrawdownMonitor)
	mockExposureMonitor := new(MockExposureMonitor)
	mockDailyLimitMonitor := new(MockDailyLimitMonitor)
	mockLogger := new(MockLogger)

	// Create risk manager
	riskManager := NewRiskManager(
		mockBalanceRepo,
		mockPositionSizer,
		mockDrawdownMonitor,
		mockExposureMonitor,
		mockDailyLimitMonitor,
		mockLogger,
	)

	// Set up expectations
	ctx := context.Background()
	newOrderValue := 100.0
	expectedAllowed := true

	mockExposureMonitor.On("CheckExposureLimit", ctx, newOrderValue, 50.0).
		Return(expectedAllowed, nil)

	// Call the method
	allowed, err := riskManager.CheckExposureLimit(ctx, newOrderValue)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAllowed, allowed)
	mockExposureMonitor.AssertExpectations(t)
}

func TestCheckDailyLossLimit(t *testing.T) {
	// Create mocks
	mockBalanceRepo := new(MockBalanceHistoryRepository)
	mockPositionSizer := new(MockPositionSizer)
	mockDrawdownMonitor := new(MockDrawdownMonitor)
	mockExposureMonitor := new(MockExposureMonitor)
	mockDailyLimitMonitor := new(MockDailyLimitMonitor)
	mockLogger := new(MockLogger)

	// Create risk manager
	riskManager := NewRiskManager(
		mockBalanceRepo,
		mockPositionSizer,
		mockDrawdownMonitor,
		mockExposureMonitor,
		mockDailyLimitMonitor,
		mockLogger,
	)

	// Set up expectations
	ctx := context.Background()
	expectedAllowed := true

	mockDailyLimitMonitor.On("CheckDailyLossLimit", ctx, 5.0).
		Return(expectedAllowed, nil)

	// Call the method
	allowed, err := riskManager.CheckDailyLossLimit(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAllowed, allowed)
	mockDailyLimitMonitor.AssertExpectations(t)
}

func TestGetRiskStatus(t *testing.T) {
	// Create mocks
	mockBalanceRepo := new(MockBalanceHistoryRepository)
	mockPositionSizer := new(MockPositionSizer)
	mockDrawdownMonitor := new(MockDrawdownMonitor)
	mockExposureMonitor := new(MockExposureMonitor)
	mockDailyLimitMonitor := new(MockDailyLimitMonitor)
	mockLogger := new(MockLogger)
	mockAccountSvc := new(MockAccountService)

	// Set up the account service in the exposure monitor
	mockExposureMonitor.accountSvc = mockAccountSvc

	// Create risk manager
	riskManager := NewRiskManager(
		mockBalanceRepo,
		mockPositionSizer,
		mockDrawdownMonitor,
		mockExposureMonitor,
		mockDailyLimitMonitor,
		mockLogger,
	)

	// Set up expectations
	ctx := context.Background()
	accountBalance := 1000.0
	drawdown := 0.15
	totalExposure := 500.0
	todayPnL := 50.0

	mockExposureMonitor.On("GetAccountBalance", ctx).Return(accountBalance, nil)
	mockDrawdownMonitor.On("CalculateDrawdown", ctx, 90).Return(drawdown, nil)
	mockExposureMonitor.On("CalculateTotalExposure", ctx).Return(totalExposure, nil)
	mockDailyLimitMonitor.On("CalculateDailyPnL", ctx).Return(todayPnL, nil)
	mockDrawdownMonitor.On("CheckDrawdownLimit", ctx, 20.0).Return(true, nil)
	mockDailyLimitMonitor.On("CheckDailyLossLimit", ctx, 5.0).Return(true, nil)

	// Call the method
	status, err := riskManager.GetRiskStatus(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, drawdown, status.CurrentDrawdown)
	assert.Equal(t, totalExposure, status.TotalExposure)
	assert.Equal(t, todayPnL, status.TodayPnL)
	assert.Equal(t, accountBalance, status.AccountBalance)
	assert.True(t, status.TradingEnabled)
	assert.Empty(t, status.DisabledReason)
	mockAccountSvc.AssertExpectations(t)
	mockDrawdownMonitor.AssertExpectations(t)
	mockExposureMonitor.AssertExpectations(t)
	mockDailyLimitMonitor.AssertExpectations(t)
}

func TestUpdateRiskParameters(t *testing.T) {
	// Create mocks
	mockBalanceRepo := new(MockBalanceHistoryRepository)
	mockPositionSizer := new(MockPositionSizer)
	mockDrawdownMonitor := new(MockDrawdownMonitor)
	mockExposureMonitor := new(MockExposureMonitor)
	mockDailyLimitMonitor := new(MockDailyLimitMonitor)
	mockLogger := new(MockLogger)

	// Create risk manager
	riskManager := NewRiskManager(
		mockBalanceRepo,
		mockPositionSizer,
		mockDrawdownMonitor,
		mockExposureMonitor,
		mockDailyLimitMonitor,
		mockLogger,
	)

	// Set up expectations
	ctx := context.Background()
	params := RiskParameters{
		MaxDrawdownPercent:    15.0,
		RiskPerTradePercent:   0.5,
		MaxExposurePercent:    40.0,
		DailyLossLimitPercent: 3.0,
		MinAccountBalance:     200.0,
	}

	mockLogger.On("Info", "Updated risk parameters", mock.Anything).Return()

	// Call the method
	err := riskManager.UpdateRiskParameters(ctx, params)

	// Assert
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)

	// Verify that the parameters were updated
	riskManager.lock.RLock()
	defer riskManager.lock.RUnlock()
	assert.Equal(t, params.MaxDrawdownPercent, riskManager.riskParams.MaxDrawdownPercent)
	assert.Equal(t, params.RiskPerTradePercent, riskManager.riskParams.RiskPerTradePercent)
	assert.Equal(t, params.MaxExposurePercent, riskManager.riskParams.MaxExposurePercent)
	assert.Equal(t, params.DailyLossLimitPercent, riskManager.riskParams.DailyLossLimitPercent)
	assert.Equal(t, params.MinAccountBalance, riskManager.riskParams.MinAccountBalance)
}

func TestIsTradeAllowed(t *testing.T) {
	// Create mocks
	mockBalanceRepo := new(MockBalanceHistoryRepository)
	mockPositionSizer := new(MockPositionSizer)
	mockDrawdownMonitor := new(MockDrawdownMonitor)
	mockExposureMonitor := new(MockExposureMonitor)
	mockDailyLimitMonitor := new(MockDailyLimitMonitor)
	mockLogger := new(MockLogger)
	mockAccountSvc := new(MockAccountService)

	// Set up the account service in the exposure monitor
	mockExposureMonitor.accountSvc = mockAccountSvc

	// Create risk manager
	riskManager := NewRiskManager(
		mockBalanceRepo,
		mockPositionSizer,
		mockDrawdownMonitor,
		mockExposureMonitor,
		mockDailyLimitMonitor,
		mockLogger,
	)

	// Test cases
	testCases := []struct {
		name            string
		symbol          string
		orderValue      float64
		accountBalance  float64
		drawdownOK      bool
		exposureOK      bool
		dailyLossOK     bool
		expectedAllowed bool
		expectedReason  string
	}{
		{
			name:            "All checks pass",
			symbol:          "BTC/USDT",
			orderValue:      100.0,
			accountBalance:  1000.0,
			drawdownOK:      true,
			exposureOK:      true,
			dailyLossOK:     true,
			expectedAllowed: true,
			expectedReason:  "",
		},
		{
			name:            "Account balance below minimum",
			symbol:          "BTC/USDT",
			orderValue:      100.0,
			accountBalance:  50.0,
			drawdownOK:      true,
			exposureOK:      true,
			dailyLossOK:     true,
			expectedAllowed: false,
			expectedReason:  "Account balance below minimum: 50.00 < 100.00",
		},
		{
			name:            "Drawdown limit reached",
			symbol:          "BTC/USDT",
			orderValue:      100.0,
			accountBalance:  1000.0,
			drawdownOK:      false,
			exposureOK:      true,
			dailyLossOK:     true,
			expectedAllowed: false,
			expectedReason:  "Maximum drawdown limit reached",
		},
		{
			name:            "Exposure limit would be exceeded",
			symbol:          "BTC/USDT",
			orderValue:      100.0,
			accountBalance:  1000.0,
			drawdownOK:      true,
			exposureOK:      false,
			dailyLossOK:     true,
			expectedAllowed: false,
			expectedReason:  "Maximum exposure limit would be exceeded",
		},
		{
			name:            "Daily loss limit reached",
			symbol:          "BTC/USDT",
			orderValue:      100.0,
			accountBalance:  1000.0,
			drawdownOK:      true,
			exposureOK:      true,
			dailyLossOK:     false,
			expectedAllowed: false,
			expectedReason:  "Daily loss limit reached",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			ctx := context.Background()

			mockExposureMonitor.On("GetAccountBalance", ctx).Return(tc.accountBalance, nil).Once()

			if tc.accountBalance >= 100.0 {
				mockDrawdownMonitor.On("CheckDrawdownLimit", ctx, 20.0).Return(tc.drawdownOK, nil).Once()

				if tc.drawdownOK && tc.orderValue > 0 {
					mockExposureMonitor.On("CheckExposureLimit", ctx, tc.orderValue, 50.0).Return(tc.exposureOK, nil).Once()
				}

				if tc.drawdownOK && (tc.orderValue <= 0 || tc.exposureOK) {
					mockDailyLimitMonitor.On("CheckDailyLossLimit", ctx, 5.0).Return(tc.dailyLossOK, nil).Once()
				}
			}

			// Call the method
			allowed, reason, err := riskManager.IsTradeAllowed(ctx, tc.symbol, tc.orderValue)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedAllowed, allowed)
			assert.Equal(t, tc.expectedReason, reason)
			mockAccountSvc.AssertExpectations(t)
			mockDrawdownMonitor.AssertExpectations(t)
			mockExposureMonitor.AssertExpectations(t)
			mockDailyLimitMonitor.AssertExpectations(t)
		})
	}
}
