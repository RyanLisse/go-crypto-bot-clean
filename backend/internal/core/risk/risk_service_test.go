package risk

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// DummySizer implements PositionSizer for tests.
type DummySizer struct{}

func (d *DummySizer) Calculate(ctx context.Context, accountBalance, entryPrice, stopLossPrice float64) (float64, error) {
	return 1.0, nil
}

// Mock dependencies.
type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetAccountBalance(ctx context.Context) (models.Balance, error) {
	args := m.Called(ctx)
	return args.Get(0).(models.Balance), args.Error(1)
}

func (m *MockAccountService) GetPortfolioValue(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAccountService) GetPositionRisk(ctx context.Context, symbol string) (models.PositionRisk, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(models.PositionRisk), args.Error(1)
}

func (m *MockAccountService) GetAllPositionRisks(ctx context.Context) (map[string]models.PositionRisk, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]models.PositionRisk), args.Error(1)
}

func (m *MockAccountService) ValidateAPIKeys(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockAccountService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockAccountService) GetCurrentExposure(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAccountService) UpdateBalance(ctx context.Context, amount float64, reason string) error {
	args := m.Called(ctx, amount, reason)
	return args.Error(0)
}

func (m *MockAccountService) SyncWithExchange(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAccountService) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
	args := m.Called(ctx, days)
	return args.Get(0).(*models.BalanceSummary), args.Error(1)
}

func (m *MockAccountService) GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockAccountService) AnalyzeTransactions(ctx context.Context, startTime, endTime time.Time) (*models.TransactionAnalysis, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(*models.TransactionAnalysis), args.Error(1)
}

func (m *MockAccountService) SubscribeToBalanceUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	args := m.Called(ctx, callback)
	return args.Error(0)
}

type MockPositionService struct {
	mock.Mock
}

func (m *MockPositionService) OpenPosition(ctx context.Context, symbol string, amount float64, entryPrice float64) (*models.Position, error) {
	args := m.Called(ctx, symbol, amount, entryPrice)
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionService) ClosePosition(ctx context.Context, positionID string, exitPrice float64) (*models.ClosedPosition, error) {
	args := m.Called(ctx, positionID, exitPrice)
	return args.Get(0).(*models.ClosedPosition), args.Error(1)
}

func (m *MockPositionService) GetPosition(ctx context.Context, positionID string) (*models.Position, error) {
	args := m.Called(ctx, positionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *MockPositionService) GetAllPositions(ctx context.Context) ([]*models.Position, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Position), args.Error(1)
}

func (m *MockPositionService) SetStopLoss(ctx context.Context, positionID string, price float64) error {
	args := m.Called(ctx, positionID, price)
	return args.Error(0)
}

func (m *MockPositionService) SetTakeProfit(ctx context.Context, positionID string, price float64) error {
	args := m.Called(ctx, positionID, price)
	return args.Error(0)
}

func (m *MockPositionService) GetPositionPnL(ctx context.Context, positionID string) (float64, float64, error) {
	args := m.Called(ctx, positionID)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

// TestCalculatePositionSize tests the position size calculation.
func TestCalculatePositionSize(t *testing.T) {
	accountSvc := new(MockAccountService)
	positionSvc := new(MockPositionService)

	// Set up mock for the new method
	accountSvc.On("SubscribeToBalanceUpdates", mock.Anything, mock.Anything).Return(nil).Once()

	riskService := NewRiskService(accountSvc, positionSvc, &DefaultRiskConfig{}, &DummySizer{})

	ctx := context.Background()
	portfolioValue := 10000.0
	maxRiskPerTrade := 0.02
	entryPrice := 50000.0
	stopLossPrice := 48000.0
	currentExposure := 3000.0

	// Expect two calls to GetPortfolioValue.
	accountSvc.On("GetPortfolioValue", ctx).Return(portfolioValue, nil).Times(2)
	accountSvc.On("GetCurrentExposure", ctx).Return(currentExposure, nil).Once()

	positionSize, riskAmount, err := riskService.CalculatePositionSize(ctx, entryPrice, stopLossPrice)

	expectedRiskAmount := portfolioValue * maxRiskPerTrade
	expectedPositionSize := (expectedRiskAmount / (entryPrice - stopLossPrice)) * entryPrice

	assert.NoError(t, err)
	assert.InDelta(t, expectedPositionSize, positionSize, 0.001)
	assert.InDelta(t, expectedRiskAmount, riskAmount, 0.001)
	accountSvc.AssertExpectations(t)
	positionSvc.AssertExpectations(t)
}

// TestCalculatePortfolioRisk tests the portfolio risk calculation.
func TestCalculatePortfolioRisk(t *testing.T) {
	accountSvc := new(MockAccountService)
	positionSvc := new(MockPositionService)

	// Set up mock for the new method
	accountSvc.On("SubscribeToBalanceUpdates", mock.Anything, mock.Anything).Return(nil).Once()

	riskService := NewRiskService(accountSvc, positionSvc, &DefaultRiskConfig{}, &DummySizer{})

	ctx := context.Background()
	positionRisks := map[string]models.PositionRisk{
		"BTC/USDT": {
			Symbol:      "BTC/USDT",
			ExposureUSD: 5000.0,
			RiskLevel:   "high",
		},
		"ETH/USDT": {
			Symbol:      "ETH/USDT",
			ExposureUSD: 3000.0,
			RiskLevel:   "medium",
		},
	}
	portfolioValue := 10000.0
	currentExposure := 3000.0

	accountSvc.On("GetPortfolioValue", ctx).Return(portfolioValue, nil).Times(2)
	accountSvc.On("GetCurrentExposure", ctx).Return(currentExposure, nil).Once()
	accountSvc.On("GetAllPositionRisks", ctx).Return(positionRisks, nil).Once()

	portfolioRisk, exposureBySymbol, err := riskService.CalculatePortfolioRisk(ctx)

	totalExposure := 5000.0 + 3000.0
	expectedPortfolioRisk := totalExposure / portfolioValue
	expectedExposureBySymbol := map[string]float64{
		"BTC/USDT": 5000.0 / portfolioValue,
		"ETH/USDT": 3000.0 / portfolioValue,
	}

	assert.NoError(t, err)
	assert.InDelta(t, expectedPortfolioRisk, portfolioRisk, 0.001)
	assert.Equal(t, expectedExposureBySymbol, exposureBySymbol)
	accountSvc.AssertExpectations(t)
	positionSvc.AssertExpectations(t)
}

// TestCheckRiskLimits tests the risk limits check.
func TestCheckRiskLimits(t *testing.T) {
	accountSvc := new(MockAccountService)
	positionSvc := new(MockPositionService)

	// Set up mock for the new method
	accountSvc.On("SubscribeToBalanceUpdates", mock.Anything, mock.Anything).Return(nil).Once()

	riskService := NewRiskService(accountSvc, positionSvc, &DefaultRiskConfig{}, &DummySizer{})

	ctx := context.Background()
	symbol := "BTC/USDT"
	amount := 1.0
	entryPrice := 50000.0
	stopLossPrice := 48000.0

	portfolioValue := 100000.0
	currentExposure := 30000.0

	positionRisks := map[string]models.PositionRisk{
		"ETH/USDT": {
			Symbol:      "ETH/USDT",
			ExposureUSD: 30000.0,
			RiskLevel:   "medium",
		},
		"SOL/USDT": {
			Symbol:      "SOL/USDT",
			ExposureUSD: 20000.0,
			RiskLevel:   "low",
		},
	}

	accountSvc.On("GetPortfolioValue", ctx).Return(portfolioValue, nil).Times(2)
	accountSvc.On("GetCurrentExposure", ctx).Return(currentExposure, nil).Once()
	accountSvc.On("GetAllPositionRisks", ctx).Return(positionRisks, nil).Once()

	withinLimits, currentRisk, maxRisk, err := riskService.CheckRiskLimits(ctx, symbol, amount, entryPrice, stopLossPrice)

	totalExposure := 30000.0 + 20000.0
	expectedCurrentRisk := totalExposure / portfolioValue

	assert.NoError(t, err)
	assert.True(t, withinLimits)
	assert.InDelta(t, expectedCurrentRisk, currentRisk, 0.001)
	assert.InDelta(t, 0.8, maxRisk, 0.001)
	accountSvc.AssertExpectations(t)
	positionSvc.AssertExpectations(t)
}

// TestCalculateRiskRewardRatio tests the risk-reward ratio calculation.
func TestCalculateRiskRewardRatio(t *testing.T) {
	accountSvc := new(MockAccountService)
	positionSvc := new(MockPositionService)

	// Set up mock for the new method
	accountSvc.On("SubscribeToBalanceUpdates", mock.Anything, mock.Anything).Return(nil).Once()

	riskService := NewRiskService(accountSvc, positionSvc, &DefaultRiskConfig{}, &DummySizer{})

	testCases := []struct {
		name            string
		entryPrice      float64
		takeProfitPrice float64
		stopLossPrice   float64
		expectedRatio   float64
	}{
		{"2:1 Risk-Reward Ratio", 100.0, 120.0, 90.0, 2.0},
		{"1:1 Risk-Reward Ratio", 100.0, 110.0, 90.0, 1.0},
		{"3:1 Risk-Reward Ratio", 100.0, 130.0, 90.0, 3.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ratio := riskService.CalculateRiskRewardRatio(tc.entryPrice, tc.takeProfitPrice, tc.stopLossPrice)
			assert.InDelta(t, tc.expectedRatio, ratio, 0.001)
		})
	}

	accountSvc.AssertExpectations(t)
	positionSvc.AssertExpectations(t)
}

// TestCalculateDrawdown tests the maximum drawdown calculation.
func TestCalculateDrawdown(t *testing.T) {
	accountSvc := new(MockAccountService)
	positionSvc := new(MockPositionService)

	// Set up mock for the new method
	accountSvc.On("SubscribeToBalanceUpdates", mock.Anything, mock.Anything).Return(nil).Once()

	riskService := NewRiskService(accountSvc, positionSvc, &DefaultRiskConfig{}, &DummySizer{})

	balances := []float64{10000, 11000, 9500, 10500, 8000, 8500, 9000}
	expectedDrawdown := 0.2727 // (11000 - 8000) / 11000

	drawdown := riskService.CalculateMaxDrawdown(balances)
	t.Logf("Expected: %v, Actual: %v", expectedDrawdown, drawdown)
	assert.InDelta(t, expectedDrawdown, drawdown, 0.001)

	accountSvc.AssertExpectations(t)
	positionSvc.AssertExpectations(t)
}

// TestIsTradeAllowed tests the full trade allowance check.
func TestIsTradeAllowed(t *testing.T) {
	accountSvc := new(MockAccountService)
	positionSvc := new(MockPositionService)

	// Set up mock for the new method
	accountSvc.On("SubscribeToBalanceUpdates", mock.Anything, mock.Anything).Return(nil).Once()

	riskService := NewRiskService(accountSvc, positionSvc, &DefaultRiskConfig{}, &DummySizer{})

	ctx := context.Background()
	symbol := "BTC/USDT"
	amount := 0.1
	entryPrice := 50000.0
	takeProfitPrice := 55000.0
	stopLossPrice := 48000.0

	portfolioValue := 100000.0
	currentExposure := 0.0

	positionRisks := map[string]models.PositionRisk{
		"ETH/USDT": {
			Symbol:      "ETH/USDT",
			ExposureUSD: 30000.0,
			RiskLevel:   "medium",
		},
	}

	// Expect two calls to GetPortfolioValue within CheckRiskLimits.
	accountSvc.On("GetPortfolioValue", ctx).Return(portfolioValue, nil).Times(2)
	accountSvc.On("GetCurrentExposure", ctx).Return(currentExposure, nil).Once()
	accountSvc.On("GetAllPositionRisks", ctx).Return(positionRisks, nil).Once()

	// Note: The initial checkGlobalLimits call in IsTradeAllowed has been removed.
	log.Println("TestIsTradeAllowed called")
	allowed, reason, err := riskService.IsTradeAllowed(ctx, symbol, amount, entryPrice, takeProfitPrice, stopLossPrice)

	// Expected risk-reward ratio: (55000-50000) / (50000-48000) = 5000/2000 = 2.5.
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Empty(t, reason)
	accountSvc.AssertExpectations(t)
	positionSvc.AssertExpectations(t)
}
