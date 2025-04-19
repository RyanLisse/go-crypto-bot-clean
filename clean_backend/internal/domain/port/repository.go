package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model" // Import model
)

// BaseRepository defines common methods for all repositories
type BaseRepository interface {
	// Create inserts a new entity into the database
	Create(ctx context.Context, entity interface{}) error

	// Save updates an entity or creates it if it doesn't exist
	Save(ctx context.Context, entity interface{}) error

	// FindByID retrieves an entity by ID
	FindByID(ctx context.Context, entity interface{}, id interface{}) error

	// FindAll retrieves all entities matching the given conditions
	FindAll(ctx context.Context, entities interface{}, conditions interface{}, args ...interface{}) error

	// FindAllWithPagination retrieves entities with pagination
	FindAllWithPagination(ctx context.Context, entities interface{}, page, limit int, conditions interface{}, args ...interface{}) error

	// Count returns the number of entities matching the given conditions
	Count(ctx context.Context, model interface{}, count *int64, conditions interface{}, args ...interface{}) error

	// Delete removes an entity from the database
	Delete(ctx context.Context, entity interface{}) error

	// Transaction executes operations within a database transaction
	Transaction(ctx context.Context, fn func(tx interface{}) error) error

	// FindOne retrieves a single entity matching the given conditions
	FindOne(ctx context.Context, entity interface{}, conditions interface{}, args ...interface{}) error

	// DeleteByID removes an entity by ID
	DeleteByID(ctx context.Context, model interface{}, id interface{}) error

	// Update updates an entity with the given fields
	Update(ctx context.Context, entity interface{}, updates map[string]interface{}) error
}

// RepositoryFactory defines the interface for creating repository instances
type RepositoryFactory interface {
	WithEncryptionService(encryptionSvc EncryptionService) RepositoryFactory // Allow chaining
	CreateMarketRepository() MarketRepository
	CreateSymbolRepository() SymbolRepository
	CreateUserRepository() UserRepository
	CreateWalletRepository() WalletRepository
	CreateOrderRepository() OrderRepository
	CreateAPICredentialRepository() APICredentialRepository
	CreateNewCoinRepository() NewCoinRepository // Added
	CreateEventRepository() EventRepository     // Added
	// CreateRiskRepository() RiskRepository // Keep commented if not implemented
}

// NewCoinRepository defines methods for interacting with new coin data
type NewCoinRepository interface {
	Save(ctx context.Context, coin *model.NewCoin) error
	GetBySymbol(ctx context.Context, symbol string) (*model.NewCoin, error)
	GetByID(ctx context.Context, id string) (*model.NewCoin, error)
	Update(ctx context.Context, coin *model.NewCoin) error
	Delete(ctx context.Context, id string) error
	GetByStatus(ctx context.Context, status model.CoinStatus) ([]*model.NewCoin, error)
	GetRecent(ctx context.Context, limit int) ([]*model.NewCoin, error)
}

// EventRepository defines methods for storing and retrieving domain events
type EventRepository interface {
	SaveEvent(ctx context.Context, event *model.NewCoinEvent) error // Example, adjust model as needed
	// Add other methods like GetEventsByCoinID, GetEventsByType, etc.
}
