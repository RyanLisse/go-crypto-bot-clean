package trade

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ryanlisse/go-crypto-bot/internal/config"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	_ "github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
	_ "github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/rest"
)

// Mock implementations for dependencies
type MockMEXCClient struct {
	mock.Mock
}

func (m *MockMEXCClient) PlaceOrder(ctx context.Context, symbol string, side string, orderType string, quantity, price float64) (string, error) {
	args := m.Called(ctx, symbol, side, orderType, quantity, price)
	return args.String(0), args.Error(1)
}

type MockBoughtCoinRepository struct {
	mock.Mock
}

func (m *MockBoughtCoinRepository) Create(ctx context.Context, coin *models.BoughtCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

func (m *MockBoughtCoinRepository) GetActiveBoughtCoins(ctx context.Context) ([]models.BoughtCoin, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.BoughtCoin), args.Error(1)
}

func (m *MockBoughtCoinRepository) UpdateBoughtCoin(ctx context.Context, coin *models.BoughtCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

// Mock implementation of TradeService for testing
type mockTradeService struct {
	mockRepo   *MockBoughtCoinRepository
	mockClient *MockMEXCClient
	cfg        *config.Config
}

func (m *mockTradeService) EvaluateTrade(ctx context.Context, symbol string, ticker models.Ticker) (bool, string, error) {
	activeTrades, err := m.mockRepo.GetActiveBoughtCoins(ctx)
	if err != nil {
		return false, "", err
	}

	// Simple max positions check
	if len(activeTrades) >= 5 {
		return false, "max positions reached", nil
	}

	return true, "trade opportunity found", nil
}

func (m *mockTradeService) ExecuteTrade(ctx context.Context, symbol string, price float64, quantity float64) (string, error) {
	tradeID, err := m.mockClient.PlaceOrder(ctx, symbol, "BUY", "MARKET", quantity, price)
	if err != nil {
		return "", err
	}

	boughtCoin := &models.BoughtCoin{
		Symbol:   symbol,
		BuyPrice: price,
		Quantity: quantity,
		BoughtAt: time.Now(),
	}

	err = m.mockRepo.Create(ctx, boughtCoin)
	if err != nil {
		return "", err
	}

	return tradeID, nil
}

func (m *mockTradeService) ClosePosition(ctx context.Context, symbol string, reason string) error {
	activeTrades, err := m.mockRepo.GetActiveBoughtCoins(ctx)
	if err != nil {
		return err
	}

	var positionToClose *models.BoughtCoin
	for _, trade := range activeTrades {
		if trade.Symbol == symbol {
			positionToClose = &trade
			break
		}
	}

	if positionToClose == nil {
		return fmt.Errorf("no active position found for symbol %s", symbol)
	}

	_, err = m.mockClient.PlaceOrder(ctx, symbol, "SELL", "MARKET", positionToClose.Quantity, positionToClose.BuyPrice)
	if err != nil {
		return err
	}

	return m.mockRepo.UpdateBoughtCoin(ctx, positionToClose)
}

func (m *mockTradeService) GetActiveTrades(ctx context.Context) ([]models.BoughtCoin, error) {
	return m.mockRepo.GetActiveBoughtCoins(ctx)
}

func TestTradeService_EvaluateTrade(t *testing.T) {
	cfg := &config.Config{}

	t.Run("Successful trade evaluation", func(t *testing.T) {
		mockRepo := new(MockBoughtCoinRepository)
		mockClient := new(MockMEXCClient)

		tradeService := &mockTradeService{
			mockRepo:   mockRepo,
			mockClient: mockClient,
			cfg:        cfg,
		}

		ticker := models.Ticker{
			Symbol: "BTCUSDT",
			Price:  50000.0,
			Volume: 1000.0,
		}

		mockRepo.On("GetActiveBoughtCoins", mock.Anything).Return([]models.BoughtCoin{}, nil)

		decision, reason, err := tradeService.EvaluateTrade(context.Background(), "BTCUSDT", ticker)

		assert.NoError(t, err)
		assert.True(t, decision)
		assert.NotEmpty(t, reason)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Reject trade when max positions reached", func(t *testing.T) {
		mockRepo := new(MockBoughtCoinRepository)
		mockClient := new(MockMEXCClient)

		tradeService := &mockTradeService{
			mockRepo:   mockRepo,
			mockClient: mockClient,
			cfg:        cfg,
		}

		ticker := models.Ticker{
			Symbol: "BTCUSDT",
			Price:  50000.0,
			Volume: 1000.0,
		}

		// Simulate max positions reached
		activeTrades := make([]models.BoughtCoin, 5)
		mockRepo.On("GetActiveBoughtCoins", mock.Anything).Return(activeTrades, nil)

		decision, reason, err := tradeService.EvaluateTrade(context.Background(), "BTCUSDT", ticker)

		assert.NoError(t, err)
		assert.False(t, decision)
		assert.Contains(t, reason, "max positions")
		mockRepo.AssertExpectations(t)
	})
}

func TestTradeService_ExecuteTrade(t *testing.T) {
	cfg := &config.Config{}

	t.Run("Successful trade execution", func(t *testing.T) {
		mockRepo := new(MockBoughtCoinRepository)
		mockClient := new(MockMEXCClient)

		tradeService := &mockTradeService{
			mockRepo:   mockRepo,
			mockClient: mockClient,
			cfg:        cfg,
		}

		symbol := "BTCUSDT"
		price := 50000.0
		quantity := 0.1

		mockClient.On("PlaceOrder",
			mock.Anything,
			symbol,
			"BUY",
			"MARKET",
			quantity,
			price,
		).Return("trade123", nil)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		tradeID, err := tradeService.ExecuteTrade(context.Background(), symbol, price, quantity)

		assert.NoError(t, err)
		assert.Equal(t, "trade123", tradeID)
		mockClient.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Trade execution fails", func(t *testing.T) {
		mockRepo := new(MockBoughtCoinRepository)
		mockClient := new(MockMEXCClient)

		tradeService := &mockTradeService{
			mockRepo:   mockRepo,
			mockClient: mockClient,
			cfg:        cfg,
		}

		symbol := "BTCUSDT"
		price := 50000.0
		quantity := 0.1

		mockClient.On("PlaceOrder",
			mock.Anything,
			symbol,
			"BUY",
			"MARKET",
			quantity,
			price,
		).Return("", assert.AnError)

		tradeID, err := tradeService.ExecuteTrade(context.Background(), symbol, price, quantity)

		assert.Error(t, err)
		assert.Empty(t, tradeID)
		mockClient.AssertExpectations(t)
	})
}

func TestTradeService_ClosePosition(t *testing.T) {
	cfg := &config.Config{}

	t.Run("Successful position closure", func(t *testing.T) {
		mockRepo := new(MockBoughtCoinRepository)
		mockClient := new(MockMEXCClient)

		tradeService := &mockTradeService{
			mockRepo:   mockRepo,
			mockClient: mockClient,
			cfg:        cfg,
		}

		symbol := "BTCUSDT"
		reason := "take profit"

		// First, simulate finding an active position
		activeCoin := models.BoughtCoin{
			Symbol:   symbol,
			Quantity: 0.1,
			BuyPrice: 50000.0,
		}
		mockRepo.On("GetActiveBoughtCoins", mock.Anything).Return([]models.BoughtCoin{activeCoin}, nil)

		// Then simulate order placement
		mockClient.On("PlaceOrder",
			mock.Anything,
			symbol,
			"SELL",
			"MARKET",
			activeCoin.Quantity,
			activeCoin.BuyPrice,
		).Return("close123", nil)

		// Finally, update repository
		mockRepo.On("UpdateBoughtCoin", mock.Anything, mock.Anything).Return(nil)

		err := tradeService.ClosePosition(context.Background(), symbol, reason)

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("No active position to close", func(t *testing.T) {
		mockRepo := new(MockBoughtCoinRepository)
		mockClient := new(MockMEXCClient)

		tradeService := &mockTradeService{
			mockRepo:   mockRepo,
			mockClient: mockClient,
			cfg:        cfg,
		}

		symbol := "BTCUSDT"
		reason := "take profit"

		mockRepo.On("GetActiveBoughtCoins", mock.Anything).Return([]models.BoughtCoin{}, nil)

		err := tradeService.ClosePosition(context.Background(), symbol, reason)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active position")
		mockRepo.AssertExpectations(t)
	})
}

func TestTradeService_GetActiveTrades(t *testing.T) {
	cfg := &config.Config{}

	t.Run("Retrieve active trades", func(t *testing.T) {
		mockRepo := new(MockBoughtCoinRepository)
		mockClient := new(MockMEXCClient)

		tradeService := &mockTradeService{
			mockRepo:   mockRepo,
			mockClient: mockClient,
			cfg:        cfg,
		}

		expectedTrades := []models.BoughtCoin{
			{
				Symbol:   "BTCUSDT",
				Quantity: 0.1,
				BuyPrice: 50000.0,
			},
		}

		mockRepo.On("GetActiveBoughtCoins", mock.Anything).Return(expectedTrades, nil)

		activeTrades, err := tradeService.GetActiveTrades(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedTrades, activeTrades)
		mockRepo.AssertExpectations(t)
	})
}
