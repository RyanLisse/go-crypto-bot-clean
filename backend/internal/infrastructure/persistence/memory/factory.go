package memory

import (
	"go-crypto-bot-clean/backend/internal/domain/ports"
)

// Factory provides methods to create in-memory repository instances
type Factory struct{}

// NewFactory creates a new in-memory repository factory
func NewFactory() ports.RepositoryFactory {
	return &Factory{}
}

// CreateOrderRepository creates a new in-memory order repository
func (f *Factory) CreateOrderRepository() ports.OrderRepository {
	return NewOrderRepository()
}

// CreatePositionRepository creates a new in-memory position repository
func (f *Factory) CreatePositionRepository() ports.PositionRepository {
	return NewPositionRepository()
}

// CreateTradeRepository creates a new in-memory trade repository
func (f *Factory) CreateTradeRepository() ports.TradeRepository {
	return NewTradeRepository()
}
