package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// Common errors
var (
	ErrInvalidOrderRequest = errors.New("invalid order request")
	ErrOrderNotFound       = errors.New("order not found")
	ErrInsufficientBalance = errors.New("insufficient balance for the order")
	ErrSymbolNotSupported  = errors.New("trading symbol not supported")
)

// MexcTradeService implements the TradeService interface for the MEXC exchange
type MexcTradeService struct {
	mexcClient    port.MEXCClient // Changed from mexcAPI to mexcClient
	marketService *MarketDataService
	symbolRepo    port.SymbolRepository
	orderRepo     port.OrderRepository
	logger        *zerolog.Logger
}

// NewMexcTradeService creates a new MexcTradeService
func NewMexcTradeService(
	mexcClient port.MEXCClient, // Changed from mexcAPI to mexcClient
	marketService *MarketDataService,
	symbolRepo port.SymbolRepository,
	orderRepo port.OrderRepository,
	logger *zerolog.Logger,
) *MexcTradeService {
	return &MexcTradeService{
		mexcClient:    mexcClient, // Changed from mexcAPI to mexcClient
		marketService: marketService,
		symbolRepo:    symbolRepo,
		orderRepo:     orderRepo,
		logger:        logger,
	}
}

// PlaceOrder creates and submits a new order to the MEXC exchange
func (s *MexcTradeService) PlaceOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	// Validate request
	if request == nil {
		return nil, ErrInvalidOrderRequest
	}

	// Check if symbol exists
	symbol, err := s.symbolRepo.GetBySymbol(ctx, request.Symbol)
	if err != nil || symbol == nil {
		s.logger.Error().Err(err).Str("symbol", request.Symbol).Msg("Symbol not found")
		return nil, ErrSymbolNotSupported
	}

	// For limit orders, verify price is set
	if request.Type == model.OrderTypeLimit && request.Price == 0 {
		return nil, errors.New("limit orders require a price")
	}

	// Place order with the exchange
	timeInForce := model.TimeInForceGTC // Default for limit orders
	if request.Type == model.OrderTypeMarket {
		timeInForce = "" // Not used for market orders
	}

	// Submit order to exchange
	order, err := s.mexcClient.PlaceOrder( // Changed from mexcAPI to mexcClient
		ctx,
		request.Symbol,
		request.Side,
		request.Type,
		request.Quantity,
		request.Price,
		timeInForce,
	)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", request.Symbol).
			Str("side", string(request.Side)).
			Str("type", string(request.Type)).
			Float64("quantity", request.Quantity).
			Float64("price", request.Price).
			Msg("Failed to place order")
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	// Save order to database
	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		s.logger.Error().Err(err).Str("orderID", order.OrderID).Msg("Failed to save order")
		// We continue because the order was placed successfully on the exchange
	}

	// Create and return OrderResponse
	response := &model.OrderResponse{
		Order:     *order,
		IsSuccess: true,
	}

	return response, nil
}

// CancelOrder cancels an existing order
func (s *MexcTradeService) CancelOrder(ctx context.Context, symbol, orderID string) error {
	// Verify order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		s.logger.Error().Err(err).Str("orderID", orderID).Msg("Failed to retrieve order")
		return ErrOrderNotFound
	}

	// If order is already in a terminal state, return early
	if order.IsComplete() {
		s.logger.Info().Str("orderID", orderID).Str("status", string(order.Status)).Msg("Order already in terminal state")
		return fmt.Errorf("order is already %s", order.Status)
	}

	// Call exchange API to cancel the order
	err = s.mexcClient.CancelOrder(ctx, symbol, orderID) // Changed from mexcAPI to mexcClient
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Str("orderID", orderID).Msg("Failed to cancel order")
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Update order status in database
	order.Status = model.OrderStatusCanceled
	order.UpdatedAt = time.Now()
	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		s.logger.Error().Err(err).Str("orderID", orderID).Msg("Failed to update order status")
		// We continue because the order was canceled successfully on the exchange
	}

	return nil
}

// GetOrderStatus retrieves the current status of an order
func (s *MexcTradeService) GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	// First check our local database
	localOrder, err := s.orderRepo.GetByID(ctx, orderID)
	if err == nil && localOrder.IsComplete() {
		// If we have the order and it's in a terminal state, return it directly
		return localOrder, nil
	}

	// Get latest status from exchange
	order, err := s.mexcClient.GetOrderStatus(ctx, symbol, orderID) // Changed from mexcAPI to mexcClient
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Str("orderID", orderID).Msg("Failed to get order status")
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}

	// Update order in database
	if localOrder != nil {
		// Update existing order
		err = s.orderRepo.Update(ctx, order)
	} else {
		// Save new order
		err = s.orderRepo.Create(ctx, order)
	}

	if err != nil {
		s.logger.Error().Err(err).Str("orderID", orderID).Msg("Failed to update order in database")
		// Continue because we still want to return the order from the exchange
	}

	return order, nil
}

// GetOpenOrders retrieves all open orders
func (s *MexcTradeService) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	// Implement this method to get open orders from MEXC API
	// For now, we'll just return open orders from our database
	limit := 100
	offset := 0
	orders, err := s.orderRepo.GetBySymbol(ctx, symbol, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get open orders")
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	// Filter orders by status
	openOrders := make([]*model.Order, 0)
	for _, order := range orders {
		if order.Status == model.OrderStatusNew || order.Status == model.OrderStatusPartiallyFilled {
			openOrders = append(openOrders, order)
		}
	}

	return openOrders, nil
}

// GetOrderHistory retrieves historical orders for a symbol
func (s *MexcTradeService) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	// For now, implement as a simple database query for completed orders
	orders, err := s.orderRepo.GetBySymbol(ctx, symbol, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get order history")
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}

	return orders, nil
}

// CalculateRequiredQuantity calculates the required quantity for an order based on amount
func (s *MexcTradeService) CalculateRequiredQuantity(ctx context.Context, symbol string, side model.OrderSide, amount float64) (float64, error) {
	// Get current ticker to determine price
	ticker, err := s.marketService.RefreshTicker(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get ticker: %w", err)
	}

	// Get symbol information for minimum quantity
	symbolInfo, err := s.symbolRepo.GetBySymbol(ctx, symbol)
	if err != nil || symbolInfo == nil {
		return 0, ErrSymbolNotSupported
	}

	// Calculate quantity based on current price and amount
	price := ticker.Price
	if price <= 0 {
		return 0, errors.New("invalid price from ticker")
	}

	// Calculate quantity: amount / price
	quantity := amount / price

	// Round to the precision required by the exchange
	// For now, we'll just return the raw quantity
	// In a real implementation, you'd apply proper rounding based on the symbol's specifications

	// Check minimum quantity
	if quantity < symbolInfo.MinQty {
		return 0, fmt.Errorf("calculated quantity %f is below minimum allowed %f", quantity, symbolInfo.MinQty)
	}

	return quantity, nil
}
