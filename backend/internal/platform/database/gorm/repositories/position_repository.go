package repositories

import (
	"context"
	"errors"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/interfaces"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/logging"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GormPositionRepository implements the PositionRepository interface using GORM
type GormPositionRepository struct {
	db     *gorm.DB
	logger *logging.LoggerWrapper
}

// NewGormPositionRepository creates a new GormPositionRepository
func NewGormPositionRepository(db *gorm.DB, logger *logging.LoggerWrapper) *GormPositionRepository {
	return &GormPositionRepository{
		db:     db,
		logger: logger,
	}
}

// FindAll returns all positions matching the filter
func (r *GormPositionRepository) FindAll(ctx context.Context, filter interfaces.PositionFilter) ([]*models.Position, error) {
	var positions []*models.Position
	query := r.db.WithContext(ctx)

	if filter.Symbol != "" {
		query = query.Where("symbol = ?", filter.Symbol)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.MinPnL != nil {
		query = query.Where("pnl >= ?", *filter.MinPnL)
	}

	if filter.MaxPnL != nil {
		query = query.Where("pnl <= ?", *filter.MaxPnL)
	}

	if filter.FromDate != nil {
		query = query.Where("created_at >= ?", *filter.FromDate)
	}

	if filter.ToDate != nil {
		query = query.Where("created_at <= ?", *filter.ToDate)
	}

	result := query.Preload("Orders").Find(&positions)
	if result.Error != nil {
		r.logger.Error("Failed to find positions", zap.Error(result.Error))
		return nil, result.Error
	}

	r.logger.Info("Found positions", zap.Int("count", len(positions)))
	return positions, nil
}

// FindByID returns a specific position by ID
func (r *GormPositionRepository) FindByID(ctx context.Context, id string) (*models.Position, error) {
	var position models.Position
	result := r.db.WithContext(ctx).Preload("Orders").First(&position, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Info("Position not found", zap.String("id", id))
			return nil, nil // Return nil, nil when not found to match interface expectation
		}
		r.logger.Error("Failed to find position", zap.Error(result.Error), zap.String("id", id))
		return nil, result.Error
	}

	r.logger.Info("Found position", zap.String("id", id), zap.String("symbol", position.Symbol))
	return &position, nil
}

// FindBySymbol returns positions for a specific symbol
func (r *GormPositionRepository) FindBySymbol(ctx context.Context, symbol string) ([]*models.Position, error) {
	var positions []*models.Position
	result := r.db.WithContext(ctx).Preload("Orders").Where("symbol = ?", symbol).Find(&positions)
	if result.Error != nil {
		r.logger.Error("Failed to find positions by symbol", zap.Error(result.Error), zap.String("symbol", symbol))
		return nil, result.Error
	}

	r.logger.Info("Found positions by symbol", zap.String("symbol", symbol), zap.Int("count", len(positions)))
	return positions, nil
}

// Create adds a new position
func (r *GormPositionRepository) Create(ctx context.Context, position *models.Position) (string, error) {
	// Generate UUID if not provided
	if position.ID == "" {
		position.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	if position.CreatedAt.IsZero() {
		position.CreatedAt = now
	}
	position.UpdatedAt = now

	result := r.db.WithContext(ctx).Create(position)
	if result.Error != nil {
		r.logger.Error("Failed to create position",
			zap.Error(result.Error),
			zap.String("symbol", position.Symbol),
			zap.String("id", position.ID))
		return "", result.Error
	}

	r.logger.Info("Created position",
		zap.String("id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.Float64("entry_price", position.EntryPrice))
	return position.ID, nil
}

// Update modifies an existing position
func (r *GormPositionRepository) Update(ctx context.Context, position *models.Position) error {
	// Update timestamp
	position.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(position)
	if result.Error != nil {
		r.logger.Error("Failed to update position",
			zap.Error(result.Error),
			zap.String("id", position.ID))
		return result.Error
	}

	r.logger.Info("Updated position",
		zap.String("id", position.ID),
		zap.String("symbol", position.Symbol),
		zap.String("status", string(position.Status)))
	return nil
}

// Delete removes a position
func (r *GormPositionRepository) Delete(ctx context.Context, id string) error {
	// Use transaction to ensure we delete both position and its related orders
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete related orders first
		if err := tx.Where("position_id = ?", id).Delete(&models.Order{}).Error; err != nil {
			r.logger.Error("Failed to delete orders for position",
				zap.Error(err),
				zap.String("position_id", id))
			return err
		}

		// Then delete the position
		result := tx.Delete(&models.Position{}, "id = ?", id)
		if result.Error != nil {
			r.logger.Error("Failed to delete position",
				zap.Error(result.Error),
				zap.String("id", id))
			return result.Error
		}

		if result.RowsAffected == 0 {
			r.logger.Warn("No position found to delete", zap.String("id", id))
			return gorm.ErrRecordNotFound
		}

		r.logger.Info("Deleted position", zap.String("id", id))
		return nil
	})
}

// AddOrder adds an order to a position
func (r *GormPositionRepository) AddOrder(ctx context.Context, positionID string, order *models.Order) error {
	// Generate order ID if not provided
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Set position ID and timestamps
	order.PositionID = positionID
	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	order.UpdatedAt = now

	// Add the order
	result := r.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		r.logger.Error("Failed to add order to position",
			zap.Error(result.Error),
			zap.String("position_id", positionID),
			zap.String("order_id", order.ID))
		return result.Error
	}

	// Update position's UpdatedAt timestamp
	err := r.db.WithContext(ctx).Model(&models.Position{}).Where("id = ?", positionID).Update("updated_at", now).Error
	if err != nil {
		r.logger.Error("Failed to update position timestamp",
			zap.Error(err),
			zap.String("position_id", positionID))
		return err
	}

	r.logger.Info("Added order to position",
		zap.String("position_id", positionID),
		zap.String("order_id", order.ID),
		zap.String("order_type", string(order.Type)),
		zap.Float64("quantity", order.Quantity),
		zap.Float64("price", order.Price))
	return nil
}

// UpdateOrder updates an order in a position
func (r *GormPositionRepository) UpdateOrder(ctx context.Context, positionID string, order *models.Order) error {
	// Update timestamp
	order.UpdatedAt = time.Now()

	// Ensure the order belongs to the specified position
	if order.PositionID != positionID {
		r.logger.Error("Order does not belong to the specified position",
			zap.String("order_position_id", order.PositionID),
			zap.String("requested_position_id", positionID))
		return errors.New("order does not belong to the specified position")
	}

	// Use transaction to update both the order and position's timestamp
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update the order
		if err := tx.Save(order).Error; err != nil {
			r.logger.Error("Failed to update order",
				zap.Error(err),
				zap.String("order_id", order.ID))
			return err
		}

		// Update position's UpdatedAt timestamp
		err := tx.Model(&models.Position{}).Where("id = ?", positionID).Update("updated_at", time.Now()).Error
		if err != nil {
			r.logger.Error("Failed to update position timestamp",
				zap.Error(err),
				zap.String("position_id", positionID))
			return err
		}

		r.logger.Info("Updated order in position",
			zap.String("position_id", positionID),
			zap.String("order_id", order.ID),
			zap.String("status", string(order.Status)))
		return nil
	})
}
