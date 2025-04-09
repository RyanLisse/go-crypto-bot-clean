package unit

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/risk"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/risk/controls"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/risk/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockBalanceHistoryRepository mocks the BalanceHistoryRepository interface
type MockBalanceHistoryRepository struct {
	mock.Mock
}

func (m *MockBalanceHistoryRepository) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time) ([]controls.BalanceRecord, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).([]controls.BalanceRecord), args.Error(1)
}

func (m *MockBalanceHistoryRepository) SaveBalanceRecord(ctx context.Context, record controls.BalanceRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

// MockPositionSizer mocks the PositionSizer interface
type MockPositionSizer struct {
	mock.Mock
}

func (m *MockPositionSizer) CalculatePositionSize(ctx context.Context, symbol string, accountBalance, riskPercent, stopLossPercent float64) (float64, error) {
	args := m.Called(ctx, symbol, accountBalance, riskPercent, stopLossPercent)
	return args.Get(0).(float64), args.Error(1)
}

// MockDrawdownMonitor mocks the DrawdownMonitor interface
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

// MockExposureMonitor mocks the ExposureMonitor interface
type MockExposureMonitor struct {
	mock.Mock
}

func (m *MockExposureMonitor) CalculateTotalExposure(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExposureMonitor) CheckExposureLimit(ctx context.Context, newOrderValue, maxExposurePercent float64) (bool, error) {
	args := m.Called(ctx, newOrderValue, maxExposurePercent)
	return args.Bool(0), args.Error(1)
}

func (m *MockExposureMonitor) GetAccountBalance(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

// MockDailyLimitMonitor mocks the DailyLimitMonitor interface
type MockDailyLimitMonitor struct {
	mock.Mock
}

func (m *MockDailyLimitMonitor) CalculateDailyPnL(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDailyLimitMonitor) CheckDailyLossLimit(ctx context.Context, dailyLossLimitPercent float64) (bool, error) {
	args := m.Called(ctx, dailyLossLimitPercent)
	return args.Bool(0), args.Error(1)
}

// TestRiskManager_CalculatePositionSize tests the CalculatePositionSize method
func TestRiskManager_CalculatePositionSize(t *testing.T) {
	// Create mocks
	balanceRepo := new(MockBalanceHistoryRepository)
	positionSizer := new(MockPositionSizer)
	drawdownMonitor := new(MockDrawdownMonitor)
	exposureMonitor := new(MockExposureMonitor)
	dailyLimitMonitor := new(MockDailyLimitMonitor)
	logger, _ := zap.NewDevelopment()

	// Create risk manager
	riskManager := service.NewRiskManager(
		balanceRepo,
		positionSizer,
		drawdownMonitor,
		exposureMonitor,
		dailyLimitMonitor,
		logger,
	)

	// Set up expectations
	ctx := context.Background()
	symbol := "BTCUSDT"
	accountBalance := 10000.0
	expectedQuantity := 0.1

	// Default risk percent from DefaultRiskParameters
	riskPercent := 1.0
	stopLossPercent := 5.0

	positionSizer.On("CalculatePositionSize", ctx, symbol, accountBalance, riskPercent, stopLossPercent).
		Return(expectedQuantity, nil)

	// Call the method
	quantity, err := riskManager.CalculatePositionSize(ctx, symbol, accountBalance)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedQuantity, quantity)
	positionSizer.AssertExpectations(t)
}

// TestRiskManager_IsTradeAllowed tests the IsTradeAllowed method
func TestRiskManager_IsTradeAllowed(t *testing.T) {
	// Create mocks
	balanceRepo := new(MockBalanceHistoryRepository)
	positionSizer := new(MockPositionSizer)
	drawdownMonitor := new(MockDrawdownMonitor)
	exposureMonitor := new(MockExposureMonitor)
	dailyLimitMonitor := new(MockDailyLimitMonitor)
	logger, _ := zap.NewDevelopment()

	// Create risk manager
	riskManager := service.NewRiskManager(
		balanceRepo,
		positionSizer,
		drawdownMonitor,
		exposureMonitor,
		dailyLimitMonitor,
		logger,
	)

	// Set up expectations
	ctx := context.Background()
	symbol := "BTCUSDT"
	orderValue := 1000.0
	accountBalance := 10000.0
	drawdownPercent := 5.0 // Below the default max of 20%

	// Default parameters from DefaultRiskParameters
	minAccountBalance := 100.0
	maxDrawdownPercent := 20.0
	dailyLossLimitPercent := 5.0
	maxExposurePercent := 50.0

	exposureMonitor.On("GetAccountBalance", ctx).Return(accountBalance, nil)
	drawdownMonitor.On("CalculateDrawdown", ctx, 90).Return(drawdownPercent, nil)
	dailyLimitMonitor.On("CheckDailyLossLimit", ctx, dailyLossLimitPercent).Return(true, nil)
	exposureMonitor.On("CheckExposureLimit", ctx, orderValue, maxExposurePercent).Return(true, nil)

	// Call the method
	allowed, reason, err := riskManager.IsTradeAllowed(ctx, symbol, orderValue)

	// Assert results
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Empty(t, reason)
	exposureMonitor.AssertExpectations(t)
	drawdownMonitor.AssertExpectations(t)
	dailyLimitMonitor.AssertExpectations(t)
}

// TestRiskManager_IsTradeAllowed_AccountBalanceTooLow tests the IsTradeAllowed method with low account balance
func TestRiskManager_IsTradeAllowed_AccountBalanceTooLow(t *testing.T) {
	// Create mocks
	balanceRepo := new(MockBalanceHistoryRepository)
	positionSizer := new(MockPositionSizer)
	drawdownMonitor := new(MockDrawdownMonitor)
	exposureMonitor := new(MockExposureMonitor)
	dailyLimitMonitor := new(MockDailyLimitMonitor)
	logger, _ := zap.NewDevelopment()

	// Create risk manager
	riskManager := service.NewRiskManager(
		balanceRepo,
		positionSizer,
		drawdownMonitor,
		exposureMonitor,
		dailyLimitMonitor,
		logger,
	)

	// Set up expectations
	ctx := context.Background()
	symbol := "BTCUSDT"
	orderValue := 1000.0
	accountBalance := 50.0 // Below the default min of 100

	// Default parameters from DefaultRiskParameters
	minAccountBalance := 100.0

	exposureMonitor.On("GetAccountBalance", ctx).Return(accountBalance, nil)

	// Call the method
	allowed, reason, err := riskManager.IsTradeAllowed(ctx, symbol, orderValue)

	// Assert results
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Contains(t, reason, "Account balance below minimum")
	exposureMonitor.AssertExpectations(t)
}

// TestRiskManager_UpdateRiskParameters tests the UpdateRiskParameters method
func TestRiskManager_UpdateRiskParameters(t *testing.T) {
	// Create mocks
	balanceRepo := new(MockBalanceHistoryRepository)
	positionSizer := new(MockPositionSizer)
	drawdownMonitor := new(MockDrawdownMonitor)
	exposureMonitor := new(MockExposureMonitor)
	dailyLimitMonitor := new(MockDailyLimitMonitor)
	logger, _ := zap.NewDevelopment()

	// Create risk manager
	riskManager := service.NewRiskManager(
		balanceRepo,
		positionSizer,
		drawdownMonitor,
		exposureMonitor,
		dailyLimitMonitor,
		logger,
	)

	// Set up new parameters
	ctx := context.Background()
	newParams := risk.RiskParameters{
		MaxDrawdownPercent:    15.0,
		RiskPerTradePercent:   2.0,
		MaxExposurePercent:    40.0,
		DailyLossLimitPercent: 3.0,
		MinAccountBalance:     200.0,
	}

	// Update parameters
	err := riskManager.UpdateRiskParameters(ctx, newParams)
	assert.NoError(t, err)

	// Get risk status to verify parameters were updated
	exposureMonitor.On("GetAccountBalance", ctx).Return(1000.0, nil)
	drawdownMonitor.On("CalculateDrawdown", ctx, 90).Return(5.0, nil)
	exposureMonitor.On("CalculateTotalExposure", ctx).Return(200.0, nil)
	dailyLimitMonitor.On("CalculateDailyPnL", ctx).Return(50.0, nil)

	status, err := riskManager.GetRiskStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, newParams, status.RiskParameters)
}

// TestRiskManager_GetRiskStatus tests the GetRiskStatus method
func TestRiskManager_GetRiskStatus(t *testing.T) {
	// Create mocks
	balanceRepo := new(MockBalanceHistoryRepository)
	positionSizer := new(MockPositionSizer)
	drawdownMonitor := new(MockDrawdownMonitor)
	exposureMonitor := new(MockExposureMonitor)
	dailyLimitMonitor := new(MockDailyLimitMonitor)
	logger, _ := zap.NewDevelopment()

	// Create risk manager
	riskManager := service.NewRiskManager(
		balanceRepo,
		positionSizer,
		drawdownMonitor,
		exposureMonitor,
		dailyLimitMonitor,
		logger,
	)

	// Set up expectations
	ctx := context.Background()
	accountBalance := 10000.0
	totalExposure := 2000.0
	drawdown := 5.0
	dailyPnL := 100.0

	exposureMonitor.On("GetAccountBalance", ctx).Return(accountBalance, nil)
	drawdownMonitor.On("CalculateDrawdown", ctx, 90).Return(drawdown, nil)
	exposureMonitor.On("CalculateTotalExposure", ctx).Return(totalExposure, nil)
	dailyLimitMonitor.On("CalculateDailyPnL", ctx).Return(dailyPnL, nil)

	// Call the method
	status, err := riskManager.GetRiskStatus(ctx)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, accountBalance, status.AccountBalance)
	assert.Equal(t, totalExposure, status.TotalExposure)
	assert.Equal(t, drawdown, status.CurrentDrawdown)
	assert.Equal(t, dailyPnL, status.TodayPnL)
	assert.True(t, status.TradingEnabled)
	assert.Empty(t, status.DisabledReason)
	exposureMonitor.AssertExpectations(t)
	drawdownMonitor.AssertExpectations(t)
	dailyLimitMonitor.AssertExpectations(t)
}
