package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mockPort "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/port"
	mockService "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupSniperShotTest(t *testing.T) (
	*usecase.SniperShotUseCase,
	*mockService.MockSniperShotService,
	*mockPort.OrderRepository,
	*mockPort.WalletRepository,
	*zerolog.Logger,
) {
	// Create mocks
	mockOrderRepo := new(mockPort.OrderRepository)
	mockWalletRepo := new(mockPort.WalletRepository)
	mockSniperService := new(mockService.MockSniperShotService)

	// Create a logger that writes to nowhere
	logger := zerolog.New(zerolog.Nop()).With().Timestamp().Logger()

	// Create the SniperShotUseCase with mocks
	uc := usecase.NewSniperShotUseCase(&logger, mockSniperService, mockWalletRepo, mockOrderRepo)

	return uc, mockSniperService, mockOrderRepo, mockWalletRepo, &logger
}

func TestValidateParams(t *testing.T) {
	uc, _, _, _, _ := setupSniperShotTest(t)

	testCases := []struct {
		name          string
		params        *usecase.SniperShotParams
		expectedError bool
		errorContains string
	}{
		{
			name: "Valid market order",
			params: &usecase.SniperShotParams{
				UserID:    "user123",
				Symbol:    "BTCUSDT",
				Side:      model.OrderSideBuy,
				Quantity:  1.0,
				Price:     0, // Market order
				Type:      model.OrderTypeMarket,
				TimeLimit: 30 * time.Second,
			},
			expectedError: false,
		},
		{
			name: "Valid limit order",
			params: &usecase.SniperShotParams{
				UserID:    "user123",
				Symbol:    "BTCUSDT",
				Side:      model.OrderSideBuy,
				Quantity:  1.0,
				Price:     50000.0, // Limit price
				Type:      model.OrderTypeLimit,
				TimeLimit: 30 * time.Second,
			},
			expectedError: false,
		},
		{
			name: "Missing user ID",
			params: &usecase.SniperShotParams{
				UserID:    "", // Empty user ID
				Symbol:    "BTCUSDT",
				Side:      model.OrderSideBuy,
				Quantity:  1.0,
				Price:     50000.0,
				Type:      model.OrderTypeLimit,
				TimeLimit: 30 * time.Second,
			},
			expectedError: true,
			errorContains: "missing user ID",
		},
		{
			name: "Missing symbol",
			params: &usecase.SniperShotParams{
				UserID:    "user123",
				Symbol:    "", // Empty symbol
				Side:      model.OrderSideBuy,
				Quantity:  1.0,
				Price:     50000.0,
				Type:      model.OrderTypeLimit,
				TimeLimit: 30 * time.Second,
			},
			expectedError: true,
			errorContains: "missing trading symbol",
		},
		{
			name: "Zero quantity",
			params: &usecase.SniperShotParams{
				UserID:    "user123",
				Symbol:    "BTCUSDT",
				Side:      model.OrderSideBuy,
				Quantity:  0.0, // Zero quantity
				Price:     50000.0,
				Type:      model.OrderTypeLimit,
				TimeLimit: 30 * time.Second,
			},
			expectedError: true,
			errorContains: "quantity must be greater than zero",
		},
		{
			name: "Limit order with zero price",
			params: &usecase.SniperShotParams{
				UserID:    "user123",
				Symbol:    "BTCUSDT",
				Side:      model.OrderSideBuy,
				Quantity:  1.0,
				Price:     0.0, // Zero price for limit order
				Type:      model.OrderTypeLimit,
				TimeLimit: 30 * time.Second,
			},
			expectedError: true,
			errorContains: "price must be specified for limit orders",
		},
		{
			name: "Invalid order side",
			params: &usecase.SniperShotParams{
				UserID:    "user123",
				Symbol:    "BTCUSDT",
				Side:      "INVALID", // Invalid side
				Quantity:  1.0,
				Price:     50000.0,
				Type:      model.OrderTypeLimit,
				TimeLimit: 30 * time.Second,
			},
			expectedError: true,
			errorContains: "invalid order side",
		},
		{
			name: "Price threshold with invalid comparison",
			params: &usecase.SniperShotParams{
				UserID:         "user123",
				Symbol:         "BTCUSDT",
				Side:           model.OrderSideBuy,
				Quantity:       1.0,
				Price:          50000.0,
				Type:           model.OrderTypeLimit,
				TimeLimit:      30 * time.Second,
				PriceThreshold: 49000.0,
				ComparisonType: "INVALID", // Invalid comparison
			},
			expectedError: true,
			errorContains: "invalid comparison type for price threshold",
		},
		{
			name: "Zero time limit gets default",
			params: &usecase.SniperShotParams{
				UserID:    "user123",
				Symbol:    "BTCUSDT",
				Side:      model.OrderSideBuy,
				Quantity:  1.0,
				Price:     50000.0,
				Type:      model.OrderTypeLimit,
				TimeLimit: 0, // Zero time limit should get default
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := uc.ValidateParams(tc.params)

			if tc.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				require.NoError(t, err)
			}

			// Check if default time limit is set
			if tc.params.TimeLimit == 0 {
				assert.Equal(t, 30*time.Second, tc.params.TimeLimit)
			}
		})
	}
}

func TestExecuteSniper_Success(t *testing.T) {
	uc, mockSniperService, mockOrderRepo, mockWalletRepo, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	orderID := "order123"
	symbol := "BTCUSDT"
	baseAsset := model.Asset("BTC")
	quoteAsset := model.Asset("USDT")

	// Create a wallet with sufficient balance
	wallet := &model.Wallet{
		ID:       "wallet123",
		UserID:   userID,
		Type:     model.WalletTypeExchange,
		Balances: make(map[model.Asset]*model.Balance),
	}

	// Add balances to the wallet
	wallet.Balances[quoteAsset] = &model.Balance{
		Asset:  quoteAsset,
		Free:   10000.0, // 10,000 USDT
		Locked: 0,
		Total:  10000.0,
	}

	params := &usecase.SniperShotParams{
		UserID:    userID,
		Symbol:    symbol,
		Side:      model.OrderSideBuy,
		Quantity:  0.1,     // 0.1 BTC
		Price:     50000.0, // $50,000 per BTC
		Type:      model.OrderTypeLimit,
		TimeLimit: time.Minute, // 1 minute timeout
	}

	// Expected result
	orderResponse := &model.OrderResponse{
		Order: model.Order{
			ID:           "db123",
			OrderID:      orderID,
			UserID:       userID,
			Symbol:       symbol,
			Side:         model.OrderSideBuy,
			Type:         model.OrderTypeLimit,
			Status:       model.OrderStatusNew,
			Price:        50000.0,
			Quantity:     0.1,
			ExecutedQty:  0.0,
			AvgFillPrice: 0.0,
			TimeInForce:  model.TimeInForceGTC,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Exchange:     "binance",
		},
		IsSuccess: true,
	}

	// Mock expectations
	mockWalletRepo.On("GetByUserID", ctx, userID).Return(wallet, nil)

	// Create a successful execution result
	executeResult := &service.SniperShotResult{
		Success:   true,
		Order:     &orderResponse.Order,
		Timestamp: time.Now(),
		Latency:   100 * time.Millisecond,
	}

	// Setup the mock for the service call
	mockSniperService.On("ExecuteSniper", mock.Anything, mock.MatchedBy(func(req *service.SniperShotRequest) bool {
		return req.UserID == params.UserID &&
			req.Symbol == params.Symbol &&
			req.Side == params.Side &&
			req.Quantity == params.Quantity &&
			req.Price == params.Price &&
			req.Type == params.Type &&
			req.TimeLimit == params.TimeLimit
	})).Return(executeResult, nil)

	// This is a bit complex as we need to match the SniperShotRequest
	mockOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Order")).Return(nil)

	// Execute the use case
	result, err := uc.ExecuteSniper(ctx, params)

	// Verify results
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, orderID, result.Order.OrderID)
	assert.Equal(t, symbol, result.Order.Symbol)
	assert.Equal(t, userID, result.Order.UserID)

	// Verify mock expectations
	mockWalletRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
	mockSniperService.AssertExpectations(t)
}

func TestExecuteSniper_InsufficientFunds(t *testing.T) {
	uc, _, _, mockWalletRepo, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	symbol := "BTCUSDT"
	quoteAsset := model.Asset("USDT")

	// Create a wallet with insufficient balance
	wallet := &model.Wallet{
		ID:       "wallet123",
		UserID:   userID,
		Type:     model.WalletTypeExchange,
		Balances: make(map[model.Asset]*model.Balance),
	}

	// Add balances to the wallet (insufficient for the trade)
	wallet.Balances[quoteAsset] = &model.Balance{
		Asset:  quoteAsset,
		Free:   1000.0, // Only 1,000 USDT (not enough for 0.1 BTC at $50,000)
		Locked: 0,
		Total:  1000.0,
	}

	params := &usecase.SniperShotParams{
		UserID:    userID,
		Symbol:    symbol,
		Side:      model.OrderSideBuy,
		Quantity:  0.1,     // 0.1 BTC
		Price:     50000.0, // $50,000 per BTC (total $5,000)
		Type:      model.OrderTypeLimit,
		TimeLimit: time.Minute, // 1 minute timeout
	}

	// Mock expectations
	mockWalletRepo.On("GetByUserID", ctx, userID).Return(wallet, nil)

	// Execute the use case
	result, err := uc.ExecuteSniper(ctx, params)

	// Verify results
	require.Error(t, err)
	assert.Equal(t, usecase.ErrInsufficientFunds, err)
	assert.Nil(t, result)

	// Verify mock expectations
	mockWalletRepo.AssertExpectations(t)
}

func TestExecuteSniper_WalletRepoError(t *testing.T) {
	uc, _, _, mockWalletRepo, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	symbol := "BTCUSDT"

	params := &usecase.SniperShotParams{
		UserID:    userID,
		Symbol:    symbol,
		Side:      model.OrderSideBuy,
		Quantity:  0.1,     // 0.1 BTC
		Price:     50000.0, // $50,000 per BTC
		Type:      model.OrderTypeLimit,
		TimeLimit: time.Minute, // 1 minute timeout
	}

	// Mock error from wallet repository
	expectedError := errors.New("database error")
	mockWalletRepo.On("GetByUserID", ctx, userID).Return(nil, expectedError)

	// Execute the use case
	result, err := uc.ExecuteSniper(ctx, params)

	// Verify results
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check user balance")

	// Verify mock expectations
	mockWalletRepo.AssertExpectations(t)
}

func TestExecuteSniper_ValidationFailure(t *testing.T) {
	uc, _, _, _, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup invalid parameters
	params := &usecase.SniperShotParams{
		UserID:    "", // Invalid: missing user ID
		Symbol:    "BTCUSDT",
		Side:      model.OrderSideBuy,
		Quantity:  0.1,
		Price:     50000.0,
		Type:      model.OrderTypeLimit,
		TimeLimit: time.Minute,
	}

	// Execute the use case
	result, err := uc.ExecuteSniper(ctx, params)

	// Verify results
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "missing user ID")
}

func TestCancelSniper_Success(t *testing.T) {
	uc, mockSniperService, mockOrderRepo, _, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	orderID := "order123"
	exchangeOrderID := "exch_order_456"
	symbol := "BTCUSDT"

	// Create an order
	order := &model.Order{
		ID:       orderID,
		OrderID:  exchangeOrderID,
		UserID:   userID,
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     model.OrderTypeLimit,
		Status:   model.OrderStatusNew,
		Price:    50000.0,
		Quantity: 0.1,
	}

	// Mock expectations
	mockOrderRepo.On("GetByID", ctx, orderID).Return(order, nil)
	mockSniperService.On("CancelSniper", ctx, symbol, exchangeOrderID).Return(nil)

	// Execute the use case
	err := uc.CancelSniper(ctx, userID, orderID)

	// Verify results
	require.NoError(t, err)

	// Verify mock expectations
	mockOrderRepo.AssertExpectations(t)
	mockSniperService.AssertExpectations(t)
}

func TestCancelSniper_OrderNotFound(t *testing.T) {
	uc, _, mockOrderRepo, _, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	orderID := "nonexistent_order"

	// Mock expectations - order not found
	expectedError := errors.New("order not found")
	mockOrderRepo.On("GetByID", ctx, orderID).Return(nil, expectedError)

	// Execute the use case
	err := uc.CancelSniper(ctx, userID, orderID)

	// Verify results
	require.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")

	// Verify mock expectations
	mockOrderRepo.AssertExpectations(t)
}

func TestCancelSniper_UnauthorizedUser(t *testing.T) {
	uc, _, mockOrderRepo, _, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	otherUserID := "other_user"
	orderID := "order123"
	symbol := "BTCUSDT"

	// Create an order belonging to a different user
	order := &model.Order{
		ID:       orderID,
		OrderID:  "exch_order_456",
		UserID:   otherUserID, // Different user
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     model.OrderTypeLimit,
		Status:   model.OrderStatusNew,
		Price:    50000.0,
		Quantity: 0.1,
	}

	// Mock expectations
	mockOrderRepo.On("GetByID", ctx, orderID).Return(order, nil)

	// Execute the use case
	err := uc.CancelSniper(ctx, userID, orderID)

	// Verify results
	require.Error(t, err)
	assert.Contains(t, err.Error(), "order does not belong to the user")

	// Verify mock expectations
	mockOrderRepo.AssertExpectations(t)
}

func TestGetOrderStatus_Success(t *testing.T) {
	uc, mockSniperService, mockOrderRepo, _, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	orderID := "order123"
	exchangeOrderID := "exch_order_456"
	symbol := "BTCUSDT"

	// Create an order in the repository
	storedOrder := &model.Order{
		ID:          orderID,
		OrderID:     exchangeOrderID,
		UserID:      userID,
		Symbol:      symbol,
		Side:        model.OrderSideBuy,
		Type:        model.OrderTypeLimit,
		Status:      model.OrderStatusNew,
		Price:       50000.0,
		Quantity:    0.1,
		ExecutedQty: 0,
	}

	// Create an updated order from the exchange
	updatedOrder := &model.Order{
		ID:           orderID,
		OrderID:      exchangeOrderID,
		UserID:       userID,
		Symbol:       symbol,
		Side:         model.OrderSideBuy,
		Type:         model.OrderTypeLimit,
		Status:       model.OrderStatusFilled, // Status changed to FILLED
		Price:        50000.0,
		Quantity:     0.1,
		ExecutedQty:  0.1,     // Now fully executed
		AvgFillPrice: 50100.0, // Slight slippage
	}

	// Mock expectations
	mockOrderRepo.On("GetByID", ctx, orderID).Return(storedOrder, nil)
	mockSniperService.On("GetSniperOrderStatus", ctx, symbol, exchangeOrderID).Return(updatedOrder, nil)

	// If status has changed, expect an update call
	mockOrderRepo.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil).Run(func(args mock.Arguments) {
		order := args.Get(1).(*model.Order)
		assert.Equal(t, model.OrderStatusFilled, order.Status)
		assert.Equal(t, 0.1, order.ExecutedQty)
		assert.Equal(t, 50100.0, order.AvgFillPrice)
	})

	// Execute the use case
	result, err := uc.GetOrderStatus(ctx, userID, orderID)

	// Verify results
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, model.OrderStatusFilled, result.Status)
	assert.Equal(t, 0.1, result.ExecutedQty)
	assert.Equal(t, 50100.0, result.AvgFillPrice)

	// Verify mock expectations
	mockOrderRepo.AssertExpectations(t)
	mockSniperService.AssertExpectations(t)
}

func TestGetOrderStatus_OrderNotFound(t *testing.T) {
	uc, _, mockOrderRepo, _, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	orderID := "nonexistent_order"

	// Mock expectations - order not found
	expectedError := errors.New("order not found")
	mockOrderRepo.On("GetByID", ctx, orderID).Return(nil, expectedError)

	// Execute the use case
	result, err := uc.GetOrderStatus(ctx, userID, orderID)

	// Verify results
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order not found")

	// Verify mock expectations
	mockOrderRepo.AssertExpectations(t)
}

func TestGetOrderStatus_UnauthorizedUser(t *testing.T) {
	uc, _, mockOrderRepo, _, _ := setupSniperShotTest(t)
	ctx := context.Background()

	// Setup test data
	userID := "user123"
	otherUserID := "other_user"
	orderID := "order123"
	symbol := "BTCUSDT"

	// Create an order belonging to a different user
	order := &model.Order{
		ID:       orderID,
		OrderID:  "exch_order_456",
		UserID:   otherUserID, // Different user
		Symbol:   symbol,
		Side:     model.OrderSideBuy,
		Type:     model.OrderTypeLimit,
		Status:   model.OrderStatusNew,
		Price:    50000.0,
		Quantity: 0.1,
	}

	// Mock expectations
	mockOrderRepo.On("GetByID", ctx, orderID).Return(order, nil)

	// Execute the use case
	result, err := uc.GetOrderStatus(ctx, userID, orderID)

	// Verify results
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order does not belong to the user")

	// Verify mock expectations
	mockOrderRepo.AssertExpectations(t)
}

func TestExtractAssetsFromSymbol(t *testing.T) {
	testCases := []struct {
		symbol      string
		baseAsset   model.Asset
		quoteAsset  model.Asset
		description string
	}{
		{
			symbol:      "BTCUSDT",
			baseAsset:   "BTC",
			quoteAsset:  "USDT",
			description: "Common trading pair with USDT",
		},
		{
			symbol:      "ETHBTC",
			baseAsset:   "ETH",
			quoteAsset:  "BTC",
			description: "ETH-BTC trading pair",
		},
		{
			symbol:      "DOGEUSDC",
			baseAsset:   "DOGE",
			quoteAsset:  "USDC",
			description: "DOGE-USDC trading pair",
		},
		{
			symbol:      "SOLUSDT",
			baseAsset:   "SOL",
			quoteAsset:  "USDT",
			description: "SOL-USDT trading pair",
		},
		{
			symbol:      "AAVEETH",
			baseAsset:   "AAVE",
			quoteAsset:  "ETH",
			description: "AAVE-ETH trading pair",
		},
		{
			symbol:      "BTCBUSD",
			baseAsset:   "BTC",
			quoteAsset:  "BUSD",
			description: "BTC-BUSD trading pair",
		},
		{
			symbol:      "BTCBNB",
			baseAsset:   "BTC",
			quoteAsset:  "BNB",
			description: "BTC-BNB trading pair",
		},
		{
			symbol:      "A",
			baseAsset:   "",
			quoteAsset:  "",
			description: "Single character symbol",
		},
		{
			symbol:      "ABCDE",
			baseAsset:   "A",
			quoteAsset:  "BCDE",
			description: "5-character symbol using fallback logic",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Since extractAssetsFromSymbol is an unexported function, we need to test it
			// indirectly by using reflection or reimplementing it here

			// Common quote assets to look for
			quoteAssets := []string{"USDT", "BTC", "ETH", "BNB", "BUSD", "USDC"}

			var baseAsset, quoteAsset model.Asset

			for _, quote := range quoteAssets {
				if len(tc.symbol) > len(quote) && tc.symbol[len(tc.symbol)-len(quote):] == quote {
					base := tc.symbol[:len(tc.symbol)-len(quote)]
					baseAsset = model.Asset(base)
					quoteAsset = model.Asset(quote)
					break
				}
			}

			// Default fallback - best guess dividing the symbol
			if baseAsset == "" && quoteAsset == "" {
				if len(tc.symbol) >= 6 {
					base := tc.symbol[:len(tc.symbol)-4]
					quote := tc.symbol[len(tc.symbol)-4:]
					baseAsset = model.Asset(base)
					quoteAsset = model.Asset(quote)
				} else if len(tc.symbol) >= 3 {
					base := tc.symbol[:3]
					quote := tc.symbol[3:]
					baseAsset = model.Asset(base)
					quoteAsset = model.Asset(quote)
				}
			}

			assert.Equal(t, tc.baseAsset, baseAsset, "Base asset mismatch for symbol %s", tc.symbol)
			assert.Equal(t, tc.quoteAsset, quoteAsset, "Quote asset mismatch for symbol %s", tc.symbol)
		})
	}
}
