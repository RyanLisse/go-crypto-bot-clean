package ports

// RepositoryFactory defines an interface for creating repository instances
type RepositoryFactory interface {
	// CreateOrderRepository creates a new order repository
	CreateOrderRepository() OrderRepository

	// CreatePositionRepository creates a new position repository
	CreatePositionRepository() PositionRepository

	// CreateTradeRepository creates a new trade repository
	CreateTradeRepository() TradeRepository
}
