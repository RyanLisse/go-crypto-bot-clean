package services

import (
	"context"
	"errors"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"

	"github.com/google/uuid"
)

// OrderService handles business logic related to orders
type OrderService struct {
	orderRepository ports.OrderRepository
	tradeRepository ports.TradeRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepository ports.OrderRepository, tradeRepository ports.TradeRepository) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		tradeRepository: tradeRepository,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if order.Symbol == "" {
		return errors.New("symbol is required")
	}
	if order.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if order.Type == models.OrderTypeLimit && order.Price <= 0 {
		return errors.New("price must be positive for limit orders")
	}

	// Generate a unique ID
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return s.orderRepository.Create(ctx, order)
}

// GetOrderByID retrieves an order by its ID
func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	return s.orderRepository.GetByID(ctx, id)
}

// ListOrders retrieves orders based on symbol and status
func (s *OrderService) ListOrders(ctx context.Context, symbol string, status models.OrderStatus) ([]*models.Order, error) {
	return s.orderRepository.List(ctx, symbol, status)
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.orderRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusNew && order.Status != models.OrderStatusPartiallyFilled {
		return errors.New("only new or partially filled orders can be canceled")
	}

	order.Status = models.OrderStatusCanceled
	order.UpdatedAt = time.Now()

	return s.orderRepository.Update(ctx, order)
}
