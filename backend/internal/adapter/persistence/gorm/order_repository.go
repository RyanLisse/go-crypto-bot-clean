package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure OrderRepository implements the port.OrderRepository interface
var _ port.OrderRepository = (*OrderRepository)(nil)

// OrderEntity represents the database model for order
type OrderEntity struct {
	ID            string    `gorm:"primaryKey;size:36"`     // Internal ID
	OrderID       string    `gorm:"size:36;index;not null"` // Exchange order ID
	ClientOrderID string    `gorm:"size:36;uniqueIndex;not null"`
	Symbol        string    `gorm:"size:20;index;not null"`
	Side          string    `gorm:"size:10;not null"`
	Type          string    `gorm:"size:20;not null"`
	Status        string    `gorm:"size:20;index;not null"`
	TimeInForce   string    `gorm:"size:10"`
	Price         float64   `gorm:"type:decimal(18,8)"`
	Quantity      float64   `gorm:"type:decimal(18,8);not null"`
	ExecutedQty   float64   `gorm:"type:decimal(18,8);default:0"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

// OrderRepository implements the port.OrderRepository interface
type OrderRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *gorm.DB, logger *zerolog.Logger) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger,
	}
}

// Create persists a new order to the database
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	r.logger.Debug().
		Str("orderID", order.OrderID).
		Str("clientOrderID", order.ClientOrderID).
		Str("symbol", order.Symbol).
		Str("side", string(order.Side)).
		Str("type", string(order.Type)).
		Msg("Creating order")

	entity := r.toEntity(order)
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	r.logger.Debug().
		Str("orderID", id).
		Msg("Getting order by ID")

	var entity OrderEntity
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found: %w", result.Error)
		}
		return nil, result.Error
	}

	return r.toDomain(&entity), nil
}

// GetByClientOrderID retrieves an order by its client order ID
func (r *OrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error) {
	r.logger.Debug().
		Str("clientOrderID", clientOrderID).
		Msg("Getting order by client order ID")

	var entity OrderEntity
	result := r.db.WithContext(ctx).Where("client_order_id = ?", clientOrderID).First(&entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found: %w", result.Error)
		}
		return nil, result.Error
	}

	return r.toDomain(&entity), nil
}

// Update updates an existing order in the database
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	r.logger.Debug().
		Str("orderID", order.OrderID).
		Str("status", string(order.Status)).
		Float64("executedQty", order.ExecutedQty).
		Msg("Updating order")

	// Only use the entity for logging purposes
	result := r.db.WithContext(ctx).Model(&OrderEntity{}).
		Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"status":       string(order.Status),
			"executed_qty": order.ExecutedQty,
			"updated_at":   time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found: %s", order.ID)
	}

	return nil
}

// GetBySymbol retrieves orders for a specific symbol with pagination
func (r *OrderRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	r.logger.Debug().
		Str("symbol", symbol).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting orders by symbol")

	var entities []OrderEntity
	result := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities)

	if result.Error != nil {
		return nil, result.Error
	}

	orders := make([]*model.Order, len(entities))
	for i, entity := range entities {
		orders[i] = r.toDomain(&entity)
	}

	return orders, nil
}

// GetByUserID retrieves orders for a specific user with pagination
func (r *OrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	r.logger.Debug().
		Str("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting orders by user ID")

	// Since the Order model doesn't have a UserID field, we can implement this
	// when the application evolves to include user association with orders
	// For now, return an empty slice as a placeholder
	return []*model.Order{}, nil
}

// GetByStatus retrieves orders with a specific status with pagination
func (r *OrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error) {
	r.logger.Debug().
		Str("status", string(status)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting orders by status")

	var entities []OrderEntity
	result := r.db.WithContext(ctx).
		Where("status = ?", string(status)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities)

	if result.Error != nil {
		return nil, result.Error
	}

	orders := make([]*model.Order, len(entities))
	for i, entity := range entities {
		orders[i] = r.toDomain(&entity)
	}

	return orders, nil
}

// Count counts orders based on provided filters
func (r *OrderRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	r.logger.Debug().
		Interface("filters", filters).
		Msg("Counting orders")

	var count int64
	query := r.db.WithContext(ctx).Model(&OrderEntity{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	result := query.Count(&count)
	return count, result.Error
}

// Delete removes an order from the database
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug().
		Str("orderID", id).
		Msg("Deleting order")

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&OrderEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found: %s", id)
	}

	return nil
}

// Helper methods for entity conversion

// toEntity converts a domain model to a database entity
func (r *OrderRepository) toEntity(order *model.Order) *OrderEntity {
	return &OrderEntity{
		ID:            order.ID,
		OrderID:       order.OrderID,
		ClientOrderID: order.ClientOrderID,
		Symbol:        order.Symbol,
		Side:          string(order.Side),
		Type:          string(order.Type),
		Status:        string(order.Status),
		TimeInForce:   string(order.TimeInForce),
		Price:         order.Price,
		Quantity:      order.Quantity,
		ExecutedQty:   order.ExecutedQty,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}
}

// toDomain converts a database entity to a domain model
func (r *OrderRepository) toDomain(entity *OrderEntity) *model.Order {
	return &model.Order{
		ID:            entity.ID,
		OrderID:       entity.OrderID,
		ClientOrderID: entity.ClientOrderID,
		Symbol:        entity.Symbol,
		Side:          model.OrderSide(entity.Side),
		Type:          model.OrderType(entity.Type),
		Status:        model.OrderStatus(entity.Status),
		TimeInForce:   model.TimeInForce(entity.TimeInForce),
		Price:         entity.Price,
		Quantity:      entity.Quantity,
		ExecutedQty:   entity.ExecutedQty,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}
