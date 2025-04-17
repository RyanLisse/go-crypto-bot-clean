package service

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	mocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/port"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestExecuteSnipe tests the ExecuteSnipe method
func TestExecuteSnipe(t *testing.T) {
	// Create mocks
	mockMexcClient := &mocks.MEXCClient{}
	mockSymbolRepo := &mocks.SymbolRepository{}
	mockOrderRepo := &mocks.OrderRepository{}
	mockMarketService := &mocks.MarketDataService{}
	
	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t))
	
	// Create test symbol
	symbol := "BTC_USDT"
	
	// Create test ticker
	ticker := &market.Ticker{
		Symbol: symbol,
		Price:  50000.0,
	}
	
	// Create test order
	order := &model.Order{
		OrderID:  "test-order-id",
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     model.OrderTypeMarket,
		Quantity: 0.002, // 100 USDT / 50000 BTC price
		Price:    50000.0,
		Status:   model.OrderStatusFilled,
	}
	
	// Set up mock expectations
	mockSymbolRepo.On("GetBySymbol", mock.Anything, symbol).Return(&model.Symbol{
		Symbol: symbol,
		Status: "TRADING",
	}, nil)
	
	mockMarketService.On("GetTicker", mock.Anything, symbol).Return(ticker, nil)
	
	mockMexcClient.On("PlaceOrder", 
		mock.Anything, 
		symbol, 
		model.OrderSideBuy, 
		model.OrderTypeMarket, 
		mock.AnythingOfType("float64"), 
		mock.AnythingOfType("float64"),
		model.TimeInForceGTC,
	).Return(order, nil)
	
	mockOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Return(nil)
	
	// Create sniper service
	sniperService := NewMexcSniperService(
		mockMexcClient,
		mockSymbolRepo,
		mockOrderRepo,
		mockMarketService,
		nil, // No listing detection service for this test
		&logger,
	)
	
	// Start the service
	err := sniperService.Start()
	assert.NoError(t, err)
	
	// Execute snipe
	ctx := context.Background()
	result, err := sniperService.ExecuteSnipe(ctx, symbol)
	
	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, symbol, result.Symbol)
	assert.Equal(t, model.OrderSideBuy, result.Side)
	assert.Equal(t, model.OrderTypeMarket, result.Type)
	
	// Wait for async operations to complete
	time.Sleep(100 * time.Millisecond)
	
	// Verify mock expectations
	mockSymbolRepo.AssertExpectations(t)
	mockMarketService.AssertExpectations(t)
	mockMexcClient.AssertExpectations(t)
	// Don't verify mockOrderRepo as it's called asynchronously
}
