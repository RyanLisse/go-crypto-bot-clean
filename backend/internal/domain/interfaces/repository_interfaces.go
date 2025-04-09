package interfaces

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PositionRepository defines operations for managing positions
type PositionRepository interface {
	// FindAll returns all positions matching the filter
	FindAll(ctx context.Context, filter PositionFilter) ([]*models.Position, error)

	// FindByID returns a specific position by ID
	FindByID(ctx context.Context, id string) (*models.Position, error)

	// FindBySymbol returns positions for a specific symbol
	FindBySymbol(ctx context.Context, symbol string) ([]*models.Position, error)

	// Create adds a new position
	Create(ctx context.Context, position *models.Position) (string, error)

	// Update modifies an existing position
	Update(ctx context.Context, position *models.Position) error

	// Delete removes a position
	Delete(ctx context.Context, id string) error

	// AddOrder adds an order to a position
	AddOrder(ctx context.Context, positionID string, order *models.Order) error

	// UpdateOrder updates an order in a position
	UpdateOrder(ctx context.Context, positionID string, order *models.Order) error
}

// PositionFilter used for filtering positions in queries
type PositionFilter struct {
	Symbol   string
	Status   string
	MinPnL   *float64
	MaxPnL   *float64
	FromDate *time.Time
	ToDate   *time.Time
}

// BoughtCoinRepository defines operations for managing bought coins
type BoughtCoinRepository interface {
	// Create adds a new bought coin
	Create(ctx context.Context, coin *models.BoughtCoin) (int64, error)

	// FindByID returns a specific bought coin by ID
	FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error)

	// FindAll returns all bought coins
	FindAll(ctx context.Context) ([]models.BoughtCoin, error)

	// Update modifies an existing bought coin
	Update(ctx context.Context, coin *models.BoughtCoin) error

	// Delete removes a bought coin
	Delete(ctx context.Context, id int64) error
}

// PortfolioRepository defines the interface for portfolio data access
type PortfolioRepository interface {
	// GetPortfolio retrieves the user's portfolio
	GetPortfolio(ctx context.Context) (*models.Portfolio, error)
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// Create creates a new transaction record
	Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	
	// FindByID retrieves a transaction by its ID
	FindByID(ctx context.Context, id int64) (*models.Transaction, error)
	
	// FindByTimeRange retrieves transactions within a time range
	FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error)
	
	// FindAll retrieves all transactions
	FindAll(ctx context.Context) ([]*models.Transaction, error)
}
