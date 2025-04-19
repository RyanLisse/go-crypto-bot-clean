package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// contextKey is an unexported type for context keys defined in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

// TxContextKey is the key for storing/retrieving a *gorm.DB transaction in a context.
const TxContextKey contextKey = "databaseTransaction"

// TransactionManager defines an interface for managing database transactions.
// This allows services to perform multiple repository operations within a single transaction.
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// TriggerService defines the interface for checking trigger conditions.
type TriggerService interface {
	CheckCondition(condition *model.TriggerCondition, currentPrice float64) bool
}

// TODO: Define interfaces (ports) for repositories, external services, etc.
// Example:
// type OrderRepository interface {
//    Save(ctx context.Context, order *model.Order) error
//    GetByID(ctx context.Context, id string) (*model.Order, error)
// }
//
// type MEXCClient interface {
//    GetTicker(ctx context.Context, symbol string) (*model.Ticker, error)
//    // ... other methods
// }
