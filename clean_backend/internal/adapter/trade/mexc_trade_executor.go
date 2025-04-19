package trade

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	portgateway "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	"github.com/rs/zerolog"
)

// Ensure mexcTradeExecutor implements port.TradeExecutor
var _ port.TradeExecutor = (*mexcTradeExecutor)(nil)

type mexcTradeExecutor struct {
	gateway portgateway.MEXCGateway
	logger  zerolog.Logger
	// Add retry logic configuration if needed
}

// NewMEXCTradeExecutor creates a new TradeExecutor using the MEXC gateway.
func NewMEXCTradeExecutor(gateway portgateway.MEXCGateway, logger *zerolog.Logger) port.TradeExecutor {
	return &mexcTradeExecutor{
		gateway: gateway,
		logger:  logger.With().Str("component", "MEXCTradeExecutor").Logger(),
	}
}

// --- Interface Implementations ---

func (e *mexcTradeExecutor) executeTrade(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity, price float64) (*model.Order, error) {
	// Delegate placing order to the gateway
	// Note: The gateway's PlaceOrder returns OrderResult, need conversion or adjustment
	// Assuming PlaceOrder is updated or we handle conversion here.
	orderResult, err := e.gateway.PlaceOrder(ctx, &model.Order{
		Symbol:   symbol,
		Side:     side,
		Type:     orderType,
		Quantity: quantity,
		Price:    price, // Price is ignored by gateway for MARKET orders
	})
	if err != nil {
		// Log error
		e.logger.Error().Err(err).Msg("Failed to place order via gateway")
		// Wrap error using fmt.Errorf
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	// Convert OrderResult to Order (or assume PlaceOrder returns *model.Order directly if refactored)
	order := &model.Order{
		Exchange:      "MEXC", // Assuming MEXC
		OrderID:       orderResult.OrderID,
		ClientOrderID: orderResult.ClientOrderID,
		Symbol:        orderResult.Symbol,
		Side:          orderResult.Side,
		Type:          orderResult.Type,
		Status:        orderResult.Status,
		Price:         orderResult.Price, // Or fetch actual fill price if market order
		Quantity:      orderResult.OrigQty,
		ExecutedQty:   orderResult.ExecutedQty,
		CreatedAt:     orderResult.TransactTime,
		UpdatedAt:     time.Now(), // Or use transact time?
	}

	return order, nil
}

func (e *mexcTradeExecutor) ExecuteMarketBuy(ctx context.Context, symbol string, quantity float64) (*model.Order, error) {
	return e.executeTrade(ctx, symbol, model.OrderSideBuy, model.OrderTypeMarket, quantity, 0)
}

func (e *mexcTradeExecutor) ExecuteMarketSell(ctx context.Context, symbol string, quantity float64) (*model.Order, error) {
	return e.executeTrade(ctx, symbol, model.OrderSideSell, model.OrderTypeMarket, quantity, 0)
}

func (e *mexcTradeExecutor) ExecuteLimitBuy(ctx context.Context, symbol string, quantity, price float64) (*model.Order, error) {
	return e.executeTrade(ctx, symbol, model.OrderSideBuy, model.OrderTypeLimit, quantity, price)
}

func (e *mexcTradeExecutor) ExecuteLimitSell(ctx context.Context, symbol string, quantity, price float64) (*model.Order, error) {
	return e.executeTrade(ctx, symbol, model.OrderSideSell, model.OrderTypeLimit, quantity, price)
}

func (e *mexcTradeExecutor) CancelOrder(ctx context.Context, symbol, orderId string) error {
	return e.gateway.CancelOrder(ctx, symbol, orderId)
}

func (e *mexcTradeExecutor) CancelOrderWithRetry(ctx context.Context, symbol, orderId string) error {
	// Implement basic retry logic here or delegate if gateway handles it
	var lastErr error
	for i := 0; i < 3; i++ { // Example: 3 attempts
		err := e.gateway.CancelOrder(ctx, symbol, orderId)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(100 * time.Millisecond) // Example delay
	}
	return fmt.Errorf("failed to cancel order after retries: %w", lastErr)
}

func (e *mexcTradeExecutor) ExecuteOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	// This method might be redundant if other Execute methods exist.
	// If needed, implement similar logic to executeTrade but return OrderResponse.
	order, err := e.executeTrade(ctx, request.Symbol, request.Side, request.Type, request.Quantity, request.Price)
	if err != nil {
		// Removed Message field from error response
		return &model.OrderResponse{IsSuccess: false /* Message: err.Error() */}, err
	}
	return &model.OrderResponse{Order: *order, IsSuccess: true}, nil
}

func (e *mexcTradeExecutor) GetOrderStatus(ctx context.Context, symbol, orderId string) (*model.Order, error) {
	// Gateway might need a GetOrderStatus method
	// Placeholder:
	return nil, fmt.Errorf("GetOrderStatus not implemented in MEXC gateway yet")
	// return e.gateway.GetOrderStatus(ctx, symbol, orderId)
}

func (e *mexcTradeExecutor) GetOrderStatusWithRetry(ctx context.Context, symbol, orderId string) (*model.Order, error) {
	// Implement basic retry logic here or delegate if gateway handles it
	var lastErr error
	for i := 0; i < 3; i++ { // Example: 3 attempts
		order, err := e.GetOrderStatus(ctx, symbol, orderId)
		if err == nil {
			return order, nil
		}
		lastErr = err
		time.Sleep(100 * time.Millisecond) // Example delay
	}
	return nil, fmt.Errorf("failed to get order status after retries: %w", lastErr)
}
