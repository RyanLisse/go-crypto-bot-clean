package gorm

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure OrderRepository implements the port.OrderRepository interface
var _ port.OrderRepository = (*OrderRepository)(nil)

// OrderEntity is defined in entity.go

// OrderRepository implements the port.OrderRepository interface using GORM
type OrderRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *gorm.DB, logger *zerolog.Logger) port.OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger,
	}
}

// toDomain converts a GORM entity to a domain model
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

// toEntity converts a domain model to a GORM entity
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

// Create adds a new order to the database
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	entity := r.toEntity(order)
	result := r.db.WithContext(ctx).Create(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).
			Str("orderId", order.ID).
			Str("symbol", order.Symbol).
			Msg("Failed to create order in database")
		return result.Error
	}
	return nil
}

// Update updates an existing order in the database
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	entity := r.toEntity(order)
	result := r.db.WithContext(ctx).Save(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).
			Str("orderId", order.ID).
			Str("symbol", order.Symbol).
			Msg("Failed to update order in database")
		return result.Error
	}
	return nil
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var entity OrderEntity
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil for not found to match interface expectation
		}
		r.logger.Error().Err(result.Error).
			Str("orderId", id).
			Msg("Failed to get order by ID from database")
		return nil, result.Error
	}
	return r.toDomain(&entity), nil
}

// GetByOrderID retrieves an order by its exchange-specific order ID
func (r *OrderRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Order, error) {
	var entity OrderEntity
	result := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(result.Error).
			Str("orderID", orderID).
			Msg("Failed to get order by exchange order ID from database")
		return nil, result.Error
	}
	return r.toDomain(&entity), nil
}

// GetBySymbol retrieves orders for a symbol with pagination
func (r *OrderRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	var entities []OrderEntity
	query := r.db.WithContext(ctx).Where("symbol = ?", symbol)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Order("created_at DESC").Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).
			Str("symbol", symbol).
			Msg("Failed to get orders by symbol from database")
		return nil, result.Error
	}

	orders := make([]*model.Order, len(entities))
	for i, entity := range entities {
		orders[i] = r.toDomain(&entity)
	}

	return orders, nil
}

// GetByStatus retrieves orders with a specific status
func (r *OrderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error) {
	var entities []OrderEntity
	query := r.db.WithContext(ctx).Where("status = ?", status)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Order("created_at DESC").Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).
			Str("status", string(status)).
			Msg("Failed to get orders by status from database")
		return nil, result.Error
	}

	orders := make([]*model.Order, len(entities))
	for i, entity := range entities {
		orders[i] = r.toDomain(&entity)
	}

	return orders, nil
}

// Delete removes an order from the database
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&OrderEntity{}, "id = ?", id)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).
			Str("orderId", id).
			Msg("Failed to delete order from database")
		return result.Error
	}
	return nil
}

// Count returns the total number of orders matching the specified filters
func (r *OrderRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&OrderEntity{})

	// Apply all filters in the map
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	result := query.Count(&count)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).
			Interface("filters", filters).
			Msg("Failed to count orders from database")
		return 0, result.Error
	}

	return count, nil
}

// GetByClientOrderID retrieves an order by its client order ID
func (r *OrderRepository) GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error) {
	var entity OrderEntity
	result := r.db.WithContext(ctx).Where("client_order_id = ?", clientOrderID).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(result.Error).
			Str("clientOrderID", clientOrderID).
			Msg("Failed to get order by client order ID from database")
		return nil, result.Error
	}
	return r.toDomain(&entity), nil
}

// GetByUserID retrieves orders for a specific user with pagination
func (r *OrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error) {
	// Since the order model currently doesn't have a user association,
	// this is a placeholder implementation to satisfy the interface
	// Can be updated once user-order relationships are implemented
	r.logger.Warn().
		Str("userID", userID).
		Msg("GetByUserID called but user-order association not implemented yet")

	return []*model.Order{}, nil
}
