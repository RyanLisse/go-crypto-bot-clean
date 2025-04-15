package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// RiskProfileRepository defines the interface for risk profile persistence operations
type RiskProfileRepository interface {
	// Save creates or updates a risk profile
	Save(ctx context.Context, profile *model.RiskProfile) error
	// GetByUserID retrieves a risk profile for a specific user
	GetByUserID(ctx context.Context, userID string) (*model.RiskProfile, error)
	// Delete removes a risk profile
	Delete(ctx context.Context, id string) error
}

// RiskAssessmentRepository defines the interface for risk assessment persistence operations
type RiskAssessmentRepository interface {
	// Create adds a new risk assessment
	Create(ctx context.Context, assessment *model.RiskAssessment) error
	// Update updates an existing risk assessment
	Update(ctx context.Context, assessment *model.RiskAssessment) error
	// GetByID retrieves a risk assessment by its ID
	GetByID(ctx context.Context, id string) (*model.RiskAssessment, error)
	// GetByUserID retrieves risk assessments for a specific user with pagination
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.RiskAssessment, error)
	// GetActiveByUserID retrieves active risk assessments for a user
	GetActiveByUserID(ctx context.Context, userID string) ([]*model.RiskAssessment, error)
	// GetBySymbol retrieves risk assessments for a specific symbol
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.RiskAssessment, error)
	// GetByType retrieves risk assessments of a specific type
	GetByType(ctx context.Context, riskType model.RiskType, limit, offset int) ([]*model.RiskAssessment, error)
	// GetByLevel retrieves risk assessments of a specific level
	GetByLevel(ctx context.Context, level model.RiskLevel, limit, offset int) ([]*model.RiskAssessment, error)
	// GetByTimeRange retrieves risk assessments within a time range
	GetByTimeRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.RiskAssessment, error)
	// Count returns the total number of risk assessments matching the specified filters
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	// Delete removes a risk assessment
	Delete(ctx context.Context, id string) error
}

// RiskMetricsRepository defines the interface for risk metrics persistence operations
type RiskMetricsRepository interface {
	// Save creates or updates risk metrics
	Save(ctx context.Context, metrics *model.RiskMetrics) error
	// GetByUserID retrieves risk metrics for a specific user
	GetByUserID(ctx context.Context, userID string) (*model.RiskMetrics, error)
	// GetHistorical retrieves historical risk metrics for a user within a time range
	GetHistorical(ctx context.Context, userID string, from, to time.Time, interval string) ([]*model.RiskMetrics, error)
}

// RiskParameterRepository defines the interface for risk parameter persistence operations
type RiskParameterRepository interface {
	// GetParameters retrieves risk parameters for a user
	GetParameters(ctx context.Context, userID string) (*model.RiskParameters, error)
	// SaveParameters saves risk parameters for a user
	SaveParameters(ctx context.Context, params *model.RiskParameters) error
}

// RiskConstraintRepository defines the interface for risk constraint persistence operations
type RiskConstraintRepository interface {
	// Create adds a new risk constraint
	Create(ctx context.Context, constraint *model.RiskConstraint) error
	// Update updates an existing risk constraint
	Update(ctx context.Context, constraint *model.RiskConstraint) error
	// GetByID retrieves a risk constraint by its ID
	GetByID(ctx context.Context, id string) (*model.RiskConstraint, error)
	// GetByUserID retrieves risk constraints for a specific user
	GetByUserID(ctx context.Context, userID string) ([]*model.RiskConstraint, error)
	// GetActiveByUserID retrieves active risk constraints for a user
	GetActiveByUserID(ctx context.Context, userID string) ([]*model.RiskConstraint, error)
	// GetByType retrieves risk constraints of a specific type
	GetByType(ctx context.Context, userID string, riskType model.RiskType) ([]*model.RiskConstraint, error)
	// Delete removes a risk constraint
	Delete(ctx context.Context, id string) error
}
