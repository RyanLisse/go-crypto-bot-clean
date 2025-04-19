package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// RiskAssessmentRepository defines the interface for risk assessment persistence
type RiskAssessmentRepository interface {
	Create(ctx context.Context, assessment *model.RiskAssessment) error
	Update(ctx context.Context, assessment *model.RiskAssessment) error
	GetByID(ctx context.Context, id string) (*model.RiskAssessment, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.RiskAssessment, error)
	GetActiveByUserID(ctx context.Context, userID string) ([]*model.RiskAssessment, error)
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.RiskAssessment, error)
	GetByType(ctx context.Context, riskType model.RiskType, limit, offset int) ([]*model.RiskAssessment, error)
	GetByLevel(ctx context.Context, level model.RiskLevel, limit, offset int) ([]*model.RiskAssessment, error)
	GetByTimeRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.RiskAssessment, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) error
}

// RiskMetricsRepository defines the interface for risk metrics persistence
type RiskMetricsRepository interface {
	Save(ctx context.Context, metrics *model.RiskMetrics) error
	GetByID(ctx context.Context, id string) (*model.RiskMetrics, error)
	GetLatestByUserID(ctx context.Context, userID string) (*model.RiskMetrics, error)
	GetByUserID(ctx context.Context, userID string) (*model.RiskMetrics, error)
	GetByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*model.RiskMetrics, error)
	GetByUserIDAndPeriod(ctx context.Context, userID string, period string, limit int) ([]*model.RiskMetrics, error)
	Delete(ctx context.Context, id string) error
	DeleteOlderThan(ctx context.Context, date time.Time) error
	GetHistorical(ctx context.Context, userID string, from, to time.Time, interval string) ([]*model.RiskMetrics, error)
}

// RiskParameterRepository defines the interface for risk parameter persistence
type RiskParameterRepository interface {
	GetParameters(ctx context.Context, userID string) (*model.RiskParameters, error)
	SaveParameters(ctx context.Context, params *model.RiskParameters) error
}

// RiskProfileRepository defines the interface for risk profile persistence
type RiskProfileRepository interface {
	Save(ctx context.Context, profile *model.RiskProfile) error
	GetByUserID(ctx context.Context, userID string) (*model.RiskProfile, error)
	Delete(ctx context.Context, id string) error
}
