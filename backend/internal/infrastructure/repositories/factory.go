package repositories

import (
	"go-crypto-bot-clean/backend/internal/domain/ports"

	"gorm.io/gorm"
)

// Factory provides methods to create repository instances
type Factory struct {
	db *gorm.DB
}

// NewFactory creates a new repository factory
func NewFactory(db *gorm.DB) ports.RepositoryFactory {
	return &Factory{
		db: db,
	}
}

// CreateOrderRepository creates a new order repository
func (f *Factory) CreateOrderRepository() ports.OrderRepository {
	return NewOrderRepository(f.db)
}

// CreatePositionRepository creates a new position repository
func (f *Factory) CreatePositionRepository() ports.PositionRepository {
	return NewPositionRepository(f.db)
}

// CreateTradeRepository creates a new trade repository
func (f *Factory) CreateTradeRepository() ports.TradeRepository {
	return NewTradeRepository(f.db)
}
