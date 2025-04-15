package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// AutoBuyRuleRepository defines the interface for auto-buy rule persistence operations
type AutoBuyRuleRepository interface {
	// Create adds a new auto-buy rule
	Create(ctx context.Context, rule *model.AutoBuyRule) error

	// Update updates an existing auto-buy rule
	Update(ctx context.Context, rule *model.AutoBuyRule) error

	// GetByID retrieves an auto-buy rule by its ID
	GetByID(ctx context.Context, id string) (*model.AutoBuyRule, error)

	// GetByUserID retrieves auto-buy rules for a specific user
	GetByUserID(ctx context.Context, userID string) ([]*model.AutoBuyRule, error)

	// GetBySymbol retrieves auto-buy rules for a specific symbol
	GetBySymbol(ctx context.Context, symbol string) ([]*model.AutoBuyRule, error)

	// GetActive retrieves all active auto-buy rules
	GetActive(ctx context.Context) ([]*model.AutoBuyRule, error)

	// GetActiveByUserID retrieves active auto-buy rules for a specific user
	GetActiveByUserID(ctx context.Context, userID string) ([]*model.AutoBuyRule, error)

	// GetActiveBySymbol retrieves active auto-buy rules for a specific symbol
	GetActiveBySymbol(ctx context.Context, symbol string) ([]*model.AutoBuyRule, error)

	// GetByTriggerType retrieves auto-buy rules with a specific trigger type
	GetByTriggerType(ctx context.Context, triggerType model.TriggerType) ([]*model.AutoBuyRule, error)

	// Delete removes an auto-buy rule
	Delete(ctx context.Context, id string) error

	// Count returns the total number of auto-buy rules matching the specified filters
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
}

// AutoBuyExecutionRepository defines the interface for auto-buy execution persistence operations
type AutoBuyExecutionRepository interface {
	// Create adds a new auto-buy execution record
	Create(ctx context.Context, execution *model.AutoBuyExecution) error

	// GetByID retrieves an auto-buy execution by its ID
	GetByID(ctx context.Context, id string) (*model.AutoBuyExecution, error)

	// GetByRuleID retrieves execution records for a specific rule
	GetByRuleID(ctx context.Context, ruleID string, limit, offset int) ([]*model.AutoBuyExecution, error)

	// GetByUserID retrieves execution records for a specific user
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.AutoBuyExecution, error)

	// GetBySymbol retrieves execution records for a specific symbol
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.AutoBuyExecution, error)

	// GetByTimeRange retrieves execution records within a time range
	GetByTimeRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.AutoBuyExecution, error)

	// Count returns the total number of execution records matching the specified filters
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
}
