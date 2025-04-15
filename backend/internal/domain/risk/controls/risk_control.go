package controls

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// RiskControl defines an interface for evaluating different types of trading risks
type RiskControl interface {
	// Evaluate analyzes a specific risk factor and returns assessments if risk thresholds are exceeded
	Evaluate(ctx context.Context, userID string, profile *model.RiskProfile) ([]*model.RiskAssessment, error)

	// GetRiskType returns the type of risk this control manages
	GetRiskType() model.RiskType

	// GetName returns the human-readable name of this risk control
	GetName() string
}

// BaseRiskControl provides common functionality for risk controls
type BaseRiskControl struct {
	riskType model.RiskType
	name     string
}

// NewBaseRiskControl creates a new base risk control
func NewBaseRiskControl(riskType model.RiskType, name string) BaseRiskControl {
	return BaseRiskControl{
		riskType: riskType,
		name:     name,
	}
}

// GetRiskType returns the type of risk this control manages
func (b BaseRiskControl) GetRiskType() model.RiskType {
	return b.riskType
}

// GetName returns the human-readable name of this risk control
func (b BaseRiskControl) GetName() string {
	return b.name
}

// createRiskAssessment creates a new risk assessment - helper function since there is no model.NewRiskAssessment
func createRiskAssessment(userID string, riskType model.RiskType, level model.RiskLevel, message string) *model.RiskAssessment {
	now := time.Now()
	// Calculate a risk score based on level
	var score float64
	switch level {
	case model.RiskLevelLow:
		score = 25.0
	case model.RiskLevelMedium:
		score = 50.0
	case model.RiskLevelHigh:
		score = 75.0
	case model.RiskLevelCritical:
		score = 100.0
	default:
		score = 0.0
	}

	return &model.RiskAssessment{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      riskType,
		Level:     level,
		Status:    model.RiskStatusActive,
		Score:     score,
		Message:   message,
		CreatedAt: now,
	}
}

// RiskEvaluator manages a collection of risk controls
type RiskEvaluator struct {
	controls []RiskControl
}

// NewRiskEvaluator creates a new risk evaluator with all standard controls
func NewRiskEvaluator(marketDataService port.MarketDataService,
	positionRepo port.PositionRepository,
	orderRepo port.OrderRepository,
	walletRepo port.WalletRepository) *RiskEvaluator {

	// Create all risk controls
	controls := []RiskControl{
		NewPositionSizeControl(marketDataService, orderRepo),
		NewConcentrationControl(marketDataService, positionRepo, walletRepo),
		NewLiquidityControl(marketDataService),
		NewVolatilityControl(marketDataService),
		NewExposureControl(marketDataService, positionRepo),
		NewDrawdownControl(marketDataService, positionRepo),
	}

	return &RiskEvaluator{
		controls: controls,
	}
}

// EvaluateAllRisks runs all risk controls and aggregates the results
func (e *RiskEvaluator) EvaluateAllRisks(ctx context.Context, userID string, profile *model.RiskProfile) ([]*model.RiskAssessment, error) {
	var allAssessments []*model.RiskAssessment

	for _, control := range e.controls {
		assessments, err := control.Evaluate(ctx, userID, profile)
		if err != nil {
			return nil, fmt.Errorf("error evaluating %s risk: %w", control.GetName(), err)
		}

		allAssessments = append(allAssessments, assessments...)
	}

	return allAssessments, nil
}

// AddControl adds a new risk control to the evaluator
func (e *RiskEvaluator) AddControl(control RiskControl) {
	e.controls = append(e.controls, control)
}

// GetControlByType retrieves a risk control by its type
func (e *RiskEvaluator) GetControlByType(riskType model.RiskType) (RiskControl, bool) {
	for _, control := range e.controls {
		if control.GetRiskType() == riskType {
			return control, true
		}
	}

	return nil, false
}
