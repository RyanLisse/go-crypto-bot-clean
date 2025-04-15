package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// RiskService defines the interface for risk management operations
type RiskService interface {
	// AssessOrderRisk evaluates the risk of a new order
	AssessOrderRisk(ctx context.Context, userID string, orderRequest *model.OrderRequest) ([]*model.RiskAssessment, error)

	// AssessPositionRisk evaluates the risk of an existing or potential position
	AssessPositionRisk(ctx context.Context, userID string, position *model.Position) ([]*model.RiskAssessment, error)

	// AssessPortfolioRisk evaluates the risk of the entire portfolio
	AssessPortfolioRisk(ctx context.Context, userID string) ([]*model.RiskAssessment, error)

	// CalculateRiskMetrics calculates current risk metrics for a user
	CalculateRiskMetrics(ctx context.Context, userID string) (*model.RiskMetrics, error)

	// CheckConstraints checks if a proposed order violates any risk constraints
	CheckConstraints(ctx context.Context, userID string, orderRequest *model.OrderRequest) (bool, []*model.RiskConstraint, error)

	// GetUserRiskProfile retrieves the risk profile for a user
	GetUserRiskProfile(ctx context.Context, userID string) (*model.RiskProfile, error)

	// UpdateUserRiskProfile updates the risk profile for a user
	UpdateUserRiskProfile(ctx context.Context, profile *model.RiskProfile) error

	// GetActiveRisks retrieves all active risks for a user
	GetActiveRisks(ctx context.Context, userID string) ([]*model.RiskAssessment, error)

	// ResolveRisk marks a risk as resolved
	ResolveRisk(ctx context.Context, riskID string) error

	// IgnoreRisk marks a risk as ignored
	IgnoreRisk(ctx context.Context, riskID string) error

	// AddConstraint adds a new risk constraint
	AddConstraint(ctx context.Context, constraint *model.RiskConstraint) error

	// UpdateConstraint updates an existing risk constraint
	UpdateConstraint(ctx context.Context, constraint *model.RiskConstraint) error

	// DeleteConstraint removes a risk constraint
	DeleteConstraint(ctx context.Context, constraintID string) error

	// GetActiveConstraints retrieves all active constraints for a user
	GetActiveConstraints(ctx context.Context, userID string) ([]*model.RiskConstraint, error)
}
