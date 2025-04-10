package memory

import (
	"context"
	"errors"
	"sync"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"
)

// OrderRepository implements the ports.OrderRepository interface with in-memory storage
type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*models.Order
}

// NewOrderRepository creates a new in-memory order repository
func NewOrderRepository() ports.OrderRepository {
	return &OrderRepository{
		orders: make(map[string]*models.Order),
	}
}

// Create stores a new order in memory
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.ID == "" {
		return errors.New("order ID cannot be empty")
	}

	r.orders[order.ID] = order
	return nil
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	return order, nil
}

// List retrieves orders based on symbol and status
func (r *OrderRepository) List(ctx context.Context, symbol string, status models.OrderStatus) ([]*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Order
	for _, order := range r.orders {
		if (symbol == "" || order.Symbol == symbol) && (status == "" || order.Status == status) {
			result = append(result, order)
		}
	}

	return result, nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.ID == "" {
		return errors.New("order ID cannot be empty")
	}

	if _, exists := r.orders[order.ID]; !exists {
		return errors.New("order not found")
	}

	r.orders[order.ID] = order
	return nil
}

// Delete removes an order
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[id]; !exists {
		return errors.New("order not found")
	}

	delete(r.orders, id)
	return nil
}
