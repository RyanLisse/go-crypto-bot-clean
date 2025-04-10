package repositories

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"

	"gorm.io/gorm"
)

// OrderRepository implements the ports.OrderRepository interface
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository instance
func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// Create persists a new order in the repository
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// List retrieves orders based on symbol and status
func (r *OrderRepository) List(ctx context.Context, symbol string, status models.OrderStatus) ([]*models.Order, error) {
	var orders []*models.Order
	query := r.db.WithContext(ctx)
	
	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	
	return orders, nil
}

// Update updates an existing order in the repository
func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete removes an order from the repository
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Order{}, "id = ?", id).Error
}
