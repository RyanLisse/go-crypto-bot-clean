package repository

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// TradeDecisionFilter defines filters for querying trade decisions
type TradeDecisionFilter struct {
	Symbol      string
	Type        models.DecisionType
	Status      models.DecisionStatus
	Reason      models.DecisionReason
	Strategy    string
	StartTime   time.Time
	EndTime     time.Time
	PositionID  string
	OrderID     string
	Tags        []string
}

// TradeDecisionRepository defines the interface for trade decision persistence
type TradeDecisionRepository interface {
	// Create adds a new trade decision
	Create(ctx context.Context, decision *models.TradeDecision) (string, error)
	
	// Update modifies an existing trade decision
	Update(ctx context.Context, decision *models.TradeDecision) error
	
	// FindByID retrieves a trade decision by ID
	FindByID(ctx context.Context, id string) (*models.TradeDecision, error)
	
	// FindAll retrieves all trade decisions matching the filter
	FindAll(ctx context.Context, filter TradeDecisionFilter) ([]*models.TradeDecision, error)
	
	// FindByPositionID retrieves all trade decisions for a position
	FindByPositionID(ctx context.Context, positionID string) ([]*models.TradeDecision, error)
	
	// FindByOrderID retrieves all trade decisions for an order
	FindByOrderID(ctx context.Context, orderID string) ([]*models.TradeDecision, error)
	
	// FindBySymbol retrieves all trade decisions for a symbol
	FindBySymbol(ctx context.Context, symbol string) ([]*models.TradeDecision, error)
	
	// FindByTimeRange retrieves all trade decisions within a time range
	FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.TradeDecision, error)
	
	// CountByFilter counts trade decisions matching the filter
	CountByFilter(ctx context.Context, filter TradeDecisionFilter) (int, error)
	
	// GetSummary generates a summary of trade decisions for a time period
	GetSummary(ctx context.Context, startTime, endTime time.Time) (*models.TradeDecisionSummary, error)
	
	// Delete removes a trade decision
	Delete(ctx context.Context, id string) error
}
