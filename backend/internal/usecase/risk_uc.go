package usecase

import (
	"context"
	"errors"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// Common errors for risk management
var (
	ErrRiskProfileNotFound    = errors.New("risk profile not found")
	ErrRiskAssessmentNotFound = errors.New("risk assessment not found")
	ErrInvalidRiskData        = errors.New("invalid risk data")
)

// RiskUseCase defines the interface for risk management operations
type RiskUseCase interface {
	EvaluateOrderRisk(ctx context.Context, userID string, orderRequest model.OrderRequest) (bool, []*model.RiskAssessment, error)
	EvaluatePositionRisk(ctx context.Context, userID string, positionID string) ([]*model.RiskAssessment, error)
	EvaluatePortfolioRisk(ctx context.Context, userID string) ([]*model.RiskAssessment, error)
	GetRiskMetrics(ctx context.Context, userID string) (*model.RiskMetrics, error)
	GetHistoricalRiskMetrics(ctx context.Context, userID string, days int) ([]*model.RiskMetrics, error)
	GetActiveRisks(ctx context.Context, userID string) ([]*model.RiskAssessment, error)
	GetRiskAssessments(ctx context.Context, userID string, riskType *model.RiskType, level *model.RiskLevel, limit, offset int) ([]*model.RiskAssessment, error)
	ResolveRisk(ctx context.Context, riskID string) error
	IgnoreRisk(ctx context.Context, riskID string) error
	GetRiskProfile(ctx context.Context, userID string) (*model.RiskProfile, error)
	UpdateRiskProfile(ctx context.Context, profile *model.RiskProfile) error
	SaveRiskConstraint(ctx context.Context, constraint *model.RiskConstraint) error
	DeleteRiskConstraint(ctx context.Context, constraintID string) error
	GetActiveConstraints(ctx context.Context, userID string) ([]*model.RiskConstraint, error)
}
