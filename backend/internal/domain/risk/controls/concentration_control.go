package controls

import (
	"context"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// ConcentrationControl evaluates if the portfolio is too concentrated in one asset
type ConcentrationControl struct {
	BaseRiskControl
	marketDataService port.MarketDataService
	positionRepo      port.PositionRepository
	walletRepo        port.WalletRepository
}

// NewConcentrationControl creates a new concentration risk control
func NewConcentrationControl(
	marketDataService port.MarketDataService,
	positionRepo port.PositionRepository,
	walletRepo port.WalletRepository,
) *ConcentrationControl {
	return &ConcentrationControl{
		BaseRiskControl:   NewBaseRiskControl(model.RiskTypeConcentration, "Portfolio Concentration"),
		marketDataService: marketDataService,
		positionRepo:      positionRepo,
		walletRepo:        walletRepo,
	}
}

// Evaluate checks if the portfolio is too concentrated in any one asset
func (c *ConcentrationControl) Evaluate(ctx context.Context, userID string, profile *model.RiskProfile) ([]*model.RiskAssessment, error) {
	// Get all active positions for the user
	positions, err := c.positionRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// Get wallet to calculate total portfolio value
	wallet, err := c.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Calculate total portfolio value and value by symbol
	totalValue := 0.0
	valueBySymbol := make(map[string]float64)

	// Add wallet balances to total value
	for _, balance := range wallet.Balances {
		// For simplicity, we're only checking concentrations in positions, not wallet balances
		totalValue += balance.Total // Use Total field which is Free + Locked
	}

	// Add position values
	for _, position := range positions {
		ticker, err := c.marketDataService.GetTicker(ctx, position.Symbol)
		if err != nil {
			continue // Skip if we can't get market data
		}

		positionValue := position.Quantity * ticker.LastPrice
		valueBySymbol[position.Symbol] += positionValue
		totalValue += positionValue
	}

	var assessments []*model.RiskAssessment

	// Check concentration for each symbol
	for symbol, value := range valueBySymbol {
		if totalValue == 0 {
			continue
		}

		concentration := value / totalValue
		if concentration > profile.MaxConcentration {
			level := model.RiskLevelMedium
			if concentration > profile.MaxConcentration*1.5 {
				level = model.RiskLevelHigh
			}

			assessment := model.NewRiskAssessment(
				userID,
				model.RiskTypeConcentration,
				level,
				fmt.Sprintf("Portfolio concentration in %s is %.2f%%, exceeding limit of %.2f%%",
					symbol, concentration*100, profile.MaxConcentration*100),
			)
			assessment.Symbol = symbol
			assessment.Recommendation = "Diversify portfolio by reducing position or adding other assets"
			assessments = append(assessments, assessment)
		}
	}

	return assessments, nil
}
