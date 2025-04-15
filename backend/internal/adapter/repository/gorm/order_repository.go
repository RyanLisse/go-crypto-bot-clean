package gorm

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// OrderRepository implements port.OrderRepository using GORM
type OrderRepository struct {
	BaseRepository
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *gorm.DB, logger *zerolog.Logger) port.OrderRepository {
	l := logger.With().Str("component", "order_repository").Logger()
	return &OrderRepository{
		BaseRepository: NewBaseRepository(db, &l),
	}
}

// Create stores a new order in the database
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.BaseRepository.Create(ctx, order)
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	err := r.BaseRepository.FindOne(ctx, &order, "id = ?", id)
	if err != nil {
		return nil, err
	}
	if order.ID == "" {
		return nil, nil
	}
	return &order, nil
}

// GetByClientOrderID retrieves an order by its client order ID
func (r *OrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error) {
	var order model.Order
	err := r.BaseRepository.FindOne(ctx, &order, "client_order_id = ?", clientOrderID)
	if err != nil {
		return nil, err
	}
	if order.ID == "" {
		return nil, nil
	}
	return &order, nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.BaseRepository.Save(ctx, order)
}

// GetBySymbol retrieves orders for a specific symbol with pagination
func (r *OrderRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.BaseRepository.FindAllWithPagination(ctx, &orders, offset/limit+1, limit, "symbol = ?", symbol)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByUserID retrieves orders for a specific user with pagination
func (r *OrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.BaseRepository.FindAllWithPagination(ctx, &orders, offset/limit+1, limit, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByStatus retrieves orders with a specific status with pagination
func (r *OrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.BaseRepository.FindAllWithPagination(ctx, &orders, offset/limit+1, limit, "status = ?", status)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Count returns the total number of orders matching the given filters
func (r *OrderRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	err := r.BaseRepository.Count(ctx, &model.Order{}, &count, filters)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes an order from the database
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.DeleteByID(ctx, &model.Order{}, id)
}
