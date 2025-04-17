package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/rs/zerolog"
)

// MockTradeExecutor is a simplified mock for testing
type MockTradeExecutor struct{}

func (m *MockTradeExecutor) ExecuteOrder(ctx context.Context, order *model.OrderRequest) (*model.OrderResponse, error) {
	// Simulate a successful order execution
	return &model.OrderResponse{
		Order: model.Order{
			ID:           "mock-id",
			OrderID:      "mock-order-123",
			UserID:       order.UserID,
			Symbol:       order.Symbol,
			Side:         order.Side,
			Type:         order.Type,
			Status:       model.OrderStatusFilled,
			Price:        order.Price,
			Quantity:     order.Quantity,
			ExecutedQty:  order.Quantity,
			AvgFillPrice: order.Price,
			Exchange:     "mock-exchange",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		IsSuccess: true,
	}, nil
}

func (m *MockTradeExecutor) CancelOrderWithRetry(ctx context.Context, symbol, orderID string) error {
	return nil
}

func (m *MockTradeExecutor) GetOrderStatusWithRetry(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	return nil, nil
}

// MockOrderRepository is a simplified mock for testing
type MockOrderRepository struct{}

func (m *MockOrderRepository) Save(ctx context.Context, order *model.Order) error {
	return nil
}

func (m *MockOrderRepository) Get(ctx context.Context, id string) (*model.Order, error) {
	return nil, nil
}

func (m *MockOrderRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Order, error) {
	return nil, nil
}

func (m *MockOrderRepository) List(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	return nil, nil
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	return nil
}

// Define a mock getCurrentPrice function for this example
func init() {
	// Override the default getCurrentPrice function
	origPrice := 49000.0
	service.GetCurrentPrice = func(ctx context.Context, symbol string) (float64, error) {
		// Simulate price increasing over time
		time.Sleep(500 * time.Millisecond)
		origPrice += 100.0 // Increase by $100 each time
		fmt.Printf("Current price for %s: $%.2f\n", symbol, origPrice)
		return origPrice, nil
	}
}

func main() {
	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create a mock service with our mocks
	mockTradeExecutor := &MockTradeExecutor{}
	mockOrderRepo := &MockOrderRepository{}
	sniperService := service.NewSniperShotService(&logger, mockTradeExecutor, mockOrderRepo)

	// Create a channel to wait for callback
	callbackCalled := make(chan float64, 1)

	// Define a callback function
	priceAlertCallback := func(price float64) {
		fmt.Printf("\n🔔 PRICE ALERT! Target price reached: $%.2f\n\n", price)
		callbackCalled <- price
	}

	// Create a trigger condition
	triggerCondition := service.NewTriggerCondition(49500.0, ">=").
		WithTimeout(10).                 // 10 seconds timeout
		WithPriceBuffer(0.005).          // 0.5% buffer
		WithCheckInterval(500).          // Check every 500ms
		WithCallback(priceAlertCallback) // Add our callback

	// Create a sniper shot request
	request := &service.SniperShotRequest{
		UserID:    "example-user",
		Symbol:    "BTCUSDT",
		Side:      model.OrderSideBuy,
		Quantity:  0.1,
		Price:     50000.0,
		Type:      model.OrderTypeLimit,
		TimeLimit: 15 * time.Second,
		Condition: &triggerCondition,
	}

	fmt.Println("Starting sniper shot with trigger condition...")
	fmt.Printf("Waiting for price to reach $%.2f...\n\n", triggerCondition.TargetPrice)

	// Execute the sniper shot in a separate goroutine
	go func() {
		result, err := sniperService.ExecuteSniper(context.Background(), request)
		if err != nil {
			fmt.Printf("❌ Error executing sniper shot: %v\n", err)
			return
		}

		fmt.Printf("\n✅ Sniper shot executed successfully!\n")
		fmt.Printf("Order ID: %s\n", result.Order.OrderID)
		fmt.Printf("Status: %s\n", result.Order.Status)
		fmt.Printf("Executed Quantity: %.8f\n", result.Order.ExecutedQty)
		fmt.Printf("Average Fill Price: $%.2f\n", result.Order.AvgFillPrice)
		fmt.Printf("Execution Time: %v\n", result.Latency)
	}()

	// Wait for the callback to be called
	select {
	case price := <-callbackCalled:
		fmt.Printf("Callback received price: $%.2f\n", price)
	case <-time.After(12 * time.Second):
		fmt.Println("❌ Timeout waiting for callback")
	}

	// Wait a bit to see the final results
	time.Sleep(2 * time.Second)
	fmt.Println("\nExample completed.")
}
