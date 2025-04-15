package service

import (
	"context"
	"fmt"
	"math"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// RiskService provides risk management functionality
type RiskService struct {
	riskProfileRepo    port.RiskProfileRepository
	riskAssessmentRepo port.RiskAssessmentRepository
	riskMetricsRepo    port.RiskMetricsRepository
	riskConstraintRepo port.RiskConstraintRepository
	positionRepo       port.PositionRepository
	orderRepo          port.OrderRepository
	walletRepo         port.WalletRepository
	marketDataService  port.MarketDataService
	logger             zerolog.Logger
}

// NewRiskService creates a new RiskService with all the required dependencies
func NewRiskService(
	profileRepo port.RiskProfileRepository,
	assessmentRepo port.RiskAssessmentRepository,
	metricsRepo port.RiskMetricsRepository,
	constraintRepo port.RiskConstraintRepository,
	positionRepo port.PositionRepository,
	orderRepo port.OrderRepository,
	walletRepo port.WalletRepository,
	marketDataService port.MarketDataService,
	logger zerolog.Logger,
) *RiskService {
	return &RiskService{
		riskProfileRepo:    profileRepo,
		riskAssessmentRepo: assessmentRepo,
		riskMetricsRepo:    metricsRepo,
		riskConstraintRepo: constraintRepo,
		positionRepo:       positionRepo,
		orderRepo:          orderRepo,
		walletRepo:         walletRepo,
		marketDataService:  marketDataService,
		logger:             logger,
	}
}

// AssessOrderRisk evaluates the risk of a new order
func (s *RiskService) AssessOrderRisk(ctx context.Context, userID string, orderRequest *model.OrderRequest) ([]*model.RiskAssessment, error) {
	s.logger.Info().
		Str("userId", userID).
		Str("symbol", orderRequest.Symbol).
		Str("side", string(orderRequest.Side)).
		Float64("quantity", orderRequest.Quantity).
		Msg("Assessing order risk")

	assessments := make([]*model.RiskAssessment, 0)

	// Get user's risk profile
	profile, err := s.GetUserRiskProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk profile: %w", err)
	}

	// Get market data for the symbol
	marketData, err := s.marketDataService.GetTicker(ctx, orderRequest.Symbol) // Corrected to assumed method
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	// Calculate order value in quote currency
	orderValue := orderRequest.Quantity * marketData.Price // Adjusted to use the correct field

	// Check position size risk
	if orderValue > profile.MaxPositionSize {
		assessment := model.NewRiskAssessment(
			userID,
			model.RiskTypePosition,
			model.RiskLevelHigh,
			fmt.Sprintf("Order value %.2f exceeds maximum position size %.2f", orderValue, profile.MaxPositionSize),
		)
		assessment.Symbol = orderRequest.Symbol
		assessment.Recommendation = fmt.Sprintf("Reduce order size to below %.2f", profile.MaxPositionSize)
		assessments = append(assessments, assessment)
	}

	// Check portfolio concentration risk
	totalValue, err := s.calculateTotalPortfolioValue(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userId", userID).Msg("Failed to calculate portfolio value")
	} else if totalValue > 0 {
		concentration := orderValue / totalValue
		if concentration > profile.MaxConcentration {
			assessment := model.NewRiskAssessment(
				userID,
				model.RiskTypeConcentration,
				model.RiskLevelMedium,
				fmt.Sprintf("Order would result in %.2f%% concentration in %s, exceeding limit of %.2f%%",
					concentration*100, orderRequest.Symbol, profile.MaxConcentration*100),
			)
			assessment.Symbol = orderRequest.Symbol
			assessment.Recommendation = "Diversify portfolio by reducing position size or adding positions in other assets"
			assessments = append(assessments, assessment)
		}
	}

	// Check liquidity risk
	if marketData.Volume*marketData.Price < profile.MinLiquidity {
		assessment := model.NewRiskAssessment(
			userID,
			model.RiskTypeLiquidity,
			model.RiskLevelMedium,
			fmt.Sprintf("Symbol %s has low liquidity (%.2f)", orderRequest.Symbol, marketData.Volume*marketData.Price),
		)
		assessment.Symbol = orderRequest.Symbol
		assessment.Recommendation = "Consider trading assets with higher liquidity"
		assessments = append(assessments, assessment)
	}

	// Check volatility risk
	volatility := math.Abs(marketData.PercentChange)
	if volatility > profile.VolatilityThreshold*100 { // Convert from decimal to percentage
		level := model.RiskLevelMedium
		if volatility > profile.VolatilityThreshold*200 {
			level = model.RiskLevelHigh
		}

		assessment := model.NewRiskAssessment(
			userID,
			model.RiskTypeVolatility,
			level,
			fmt.Sprintf("Symbol %s has high volatility (%.2f%%)", orderRequest.Symbol, volatility),
		)
		assessment.Symbol = orderRequest.Symbol
		assessment.Recommendation = "Consider reducing position size or using limit orders"
		assessments = append(assessments, assessment)
	}

	// Check total exposure risk
	exposure, err := s.calculateTotalExposure(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userId", userID).Msg("Failed to calculate total exposure")
	} else {
		newExposure := exposure + orderValue
		if newExposure > profile.MaxTotalExposure {
			assessment := model.NewRiskAssessment(
				userID,
				model.RiskTypeExposure,
				model.RiskLevelHigh,
				fmt.Sprintf("Order would increase total exposure to %.2f, exceeding limit of %.2f",
					newExposure, profile.MaxTotalExposure),
			)
			assessment.Recommendation = "Close existing positions or reduce order size"
			assessments = append(assessments, assessment)
		}
	}

	return assessments, nil
}

// AssessPositionRisk evaluates the risk of an existing or potential position
func (s *RiskService) AssessPositionRisk(ctx context.Context, userID string, position *model.Position) ([]*model.RiskAssessment, error) {
	s.logger.Info().
		Str("userId", userID).
		Str("positionId", position.ID).
		Str("symbol", position.Symbol).
		Msg("Assessing position risk")

	assessments := make([]*model.RiskAssessment, 0)

	// Get user's risk profile
	profile, err := s.GetUserRiskProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk profile: %w", err)
	}

	// Get market data for the symbol
	marketData, err := s.marketDataService.GetTicker(ctx, position.Symbol) // Corrected to assumed method
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	// Check drawdown risk
	drawdown := 0.0
	if position.EntryPrice > 0 {
		if position.Side == model.PositionSideLong {
			drawdown = 1 - (marketData.Price / position.EntryPrice)
		} else {
			drawdown = 1 - (position.EntryPrice / marketData.Price)
		}
	}

	if drawdown > profile.MaxDrawdown {
		level := model.RiskLevelHigh
		if drawdown > profile.MaxDrawdown*1.5 {
			level = model.RiskLevelCritical
		}

		assessment := model.NewRiskAssessment(
			userID,
			model.RiskTypePosition,
			level,
			fmt.Sprintf("Position drawdown %.2f%% exceeds maximum %.2f%%",
				drawdown*100, profile.MaxDrawdown*100),
		)
		assessment.PositionID = position.ID
		assessment.Symbol = position.Symbol
		assessment.Recommendation = "Consider setting a stop loss or closing the position"
		assessments = append(assessments, assessment)
	}

	// More risk assessments can be added here...

	return assessments, nil
}

// AssessPortfolioRisk evaluates the risk of the entire portfolio
func (s *RiskService) AssessPortfolioRisk(ctx context.Context, userID string) ([]*model.RiskAssessment, error) {
	s.logger.Info().Str("userId", userID).Msg("Assessing portfolio risk")

	assessments := make([]*model.RiskAssessment, 0)

	// Get user's risk profile
	profile, err := s.GetUserRiskProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk profile: %w", err)
	}

	// Get all open positions
	positions, err := s.positionRepo.GetOpenPositionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// Calculate total exposure
	totalExposure := 0.0
	for _, pos := range positions {
		totalExposure += pos.Quantity * pos.CurrentPrice
	}

	// Check total exposure risk
	if totalExposure > profile.MaxTotalExposure {
		assessment := model.NewRiskAssessment(
			userID,
			model.RiskTypeExposure,
			model.RiskLevelHigh,
			fmt.Sprintf("Total exposure %.2f exceeds maximum %.2f", totalExposure, profile.MaxTotalExposure),
		)
		assessment.Recommendation = "Close some positions to reduce overall exposure"
		assessments = append(assessments, assessment)
	}

	// Check portfolio concentration
	if len(positions) > 0 {
		// Group positions by symbol
		symbolMap := make(map[string]float64)
		for _, pos := range positions {
			symbolMap[pos.Symbol] += pos.Quantity * pos.CurrentPrice
		}

		// Find highest concentration
		highestSymbol := ""
		highestValue := 0.0
		for sym, value := range symbolMap {
			if value > highestValue {
				highestValue = value
				highestSymbol = sym
			}
		}

		concentration := highestValue / totalExposure
		if concentration > profile.MaxConcentration {
			assessment := model.NewRiskAssessment(
				userID,
				model.RiskTypeConcentration,
				model.RiskLevelMedium,
				fmt.Sprintf("%.2f%% of portfolio is concentrated in %s, exceeding limit of %.2f%%",
					concentration*100, highestSymbol, profile.MaxConcentration*100),
			)
			assessment.Symbol = highestSymbol
			assessment.Recommendation = "Diversify portfolio by reducing this position or adding positions in other assets"
			assessments = append(assessments, assessment)
		}
	}

	// More portfolio-level risk assessments can be added here...

	return assessments, nil
}

// CalculateRiskMetrics calculates current risk metrics for a user
func (s *RiskService) CalculateRiskMetrics(ctx context.Context, userID string) (*model.RiskMetrics, error) {
	metrics := model.NewRiskMetrics(userID)

	// Get all open positions
	positions, err := s.positionRepo.GetOpenPositionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// Calculate total exposure
	totalExposure := 0.0
	for _, pos := range positions {
		totalExposure += pos.Quantity * pos.CurrentPrice
	}
	metrics.TotalExposure = totalExposure

	// Calculate highest concentration
	if len(positions) > 0 && totalExposure > 0 {
		// Group positions by symbol
		symbolMap := make(map[string]float64)
		for _, pos := range positions {
			symbolMap[pos.Symbol] += pos.Quantity * pos.CurrentPrice
		}

		// Find highest concentration
		highestValue := 0.0
		for _, value := range symbolMap {
			concentration := value / totalExposure
			if concentration > highestValue {
				highestValue = concentration
			}
		}

		metrics.HighestConcentration = highestValue
	}

	// Get active risks
	activeRisks, err := s.riskAssessmentRepo.GetActiveByUserID(ctx, userID)
	if err == nil {
		metrics.ActiveRiskCount = len(activeRisks)

		// Count high/critical risks
		highCount := 0
		for _, risk := range activeRisks {
			if risk.Level == model.RiskLevelHigh || risk.Level == model.RiskLevelCritical {
				highCount++
			}
		}
		metrics.HighRiskCount = highCount
	}

	// Calculate daily P&L
	// This would require historical position data

	// Calculate portfolio volatility
	// This would require historical price data

	return metrics, nil
}

// CheckConstraints checks if a proposed order violates any risk constraints
func (s *RiskService) CheckConstraints(ctx context.Context, userID string, orderRequest *model.OrderRequest) (bool, []*model.RiskConstraint, error) {
	s.logger.Info().
		Str("userId", userID).
		Str("symbol", orderRequest.Symbol).
		Msg("Checking risk constraints")

	// Get active constraints for the user
	constraints, err := s.riskConstraintRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get constraints: %w", err)
	}

	violatedConstraints := make([]*model.RiskConstraint, 0)

	// Get market data for the symbol
	marketData, err := s.marketDataService.GetTicker(ctx, orderRequest.Symbol)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get market data: %w", err)
	}

	// Calculate order value in quote currency
	orderValue := orderRequest.Quantity * marketData.Price

	// Check each constraint
	for _, constraint := range constraints {
		violated := false

		switch constraint.Type {
		case model.RiskTypePosition:
			// Check position size constraints
			if constraint.Parameter == "max_position_size" && constraint.Operator == "LT" {
				if orderValue >= constraint.Value {
					violated = true
				}
			}

		case model.RiskTypeExposure:
			// Check total exposure constraints
			if constraint.Parameter == "max_total_exposure" && constraint.Operator == "LT" {
				exposure, err := s.calculateTotalExposure(ctx, userID)
				if err == nil && exposure+orderValue >= constraint.Value {
					violated = true
				}
			}

		case model.RiskTypeLiquidity:
			// Check liquidity constraints
			if constraint.Parameter == "min_volume" && constraint.Operator == "GT" {
				if marketData.Volume <= constraint.Value {
					violated = true
				}
			}

		case model.RiskTypeVolatility:
			// Check volatility constraints
			if constraint.Parameter == "max_volatility" && constraint.Operator == "LT" {
				volatility := math.Abs(marketData.PercentChange)
				if volatility >= constraint.Value {
					violated = true
				}
			}
		}

		if violated {
			violatedConstraints = append(violatedConstraints, constraint)

			// If action is "BLOCK", return immediately
			if constraint.Action == "BLOCK" {
				return false, violatedConstraints, nil
			}
		}
	}

	// If we got here, either no constraints were violated, or none had "BLOCK" action
	return len(violatedConstraints) == 0, violatedConstraints, nil
}

// GetUserRiskProfile retrieves a user's risk profile, creating one with defaults if it doesn't exist
func (s *RiskService) GetUserRiskProfile(ctx context.Context, userID string) (*model.RiskProfile, error) {
	profile, err := s.riskProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userId", userID).Msg("Error retrieving risk profile")

		// If the profile doesn't exist, create a new one with defaults
		if err.Error() == "record not found" {
			s.logger.Info().Str("userId", userID).Msg("Creating default risk profile")
			profile = model.NewRiskProfile(userID)
			if err := s.riskProfileRepo.Save(ctx, profile); err != nil {
				return nil, fmt.Errorf("failed to create default risk profile: %w", err)
			}
			return profile, nil
		}

		return nil, err
	}

	return profile, nil
}

// UpdateUserRiskProfile updates the risk profile for a user
func (s *RiskService) UpdateUserRiskProfile(ctx context.Context, profile *model.RiskProfile) error {
	return s.riskProfileRepo.Save(ctx, profile)
}

// GetActiveRisks retrieves all active risks for a user
func (s *RiskService) GetActiveRisks(ctx context.Context, userID string) ([]*model.RiskAssessment, error) {
	return s.riskAssessmentRepo.GetActiveByUserID(ctx, userID)
}

// ResolveRisk marks a risk as resolved
func (s *RiskService) ResolveRisk(ctx context.Context, riskID string) error {
	// Get the risk assessment
	risk, err := s.riskAssessmentRepo.GetByID(ctx, riskID)
	if err != nil {
		return err
	}
	if risk == nil {
		return fmt.Errorf("risk assessment not found")
	}

	// Mark as resolved
	risk.Resolve()

	// Update in repository
	return s.riskAssessmentRepo.Update(ctx, risk)
}

// IgnoreRisk marks a risk as ignored
func (s *RiskService) IgnoreRisk(ctx context.Context, riskID string) error {
	// Get the risk assessment
	risk, err := s.riskAssessmentRepo.GetByID(ctx, riskID)
	if err != nil {
		return err
	}
	if risk == nil {
		return fmt.Errorf("risk assessment not found")
	}

	// Mark as ignored
	risk.Ignore()

	// Update in repository
	return s.riskAssessmentRepo.Update(ctx, risk)
}

// AddConstraint adds a new risk constraint
func (s *RiskService) AddConstraint(ctx context.Context, constraint *model.RiskConstraint) error {
	return s.riskConstraintRepo.Create(ctx, constraint)
}

// UpdateConstraint updates an existing risk constraint
func (s *RiskService) UpdateConstraint(ctx context.Context, constraint *model.RiskConstraint) error {
	return s.riskConstraintRepo.Update(ctx, constraint)
}

// DeleteConstraint removes a risk constraint
func (s *RiskService) DeleteConstraint(ctx context.Context, constraintID string) error {
	return s.riskConstraintRepo.Delete(ctx, constraintID)
}

// GetActiveConstraints retrieves all active constraints for a user
func (s *RiskService) GetActiveConstraints(ctx context.Context, userID string) ([]*model.RiskConstraint, error) {
	return s.riskConstraintRepo.GetActiveByUserID(ctx, userID)
}

// Helper methods

// calculateTotalExposure calculates the total market exposure for a user
func (s *RiskService) calculateTotalExposure(ctx context.Context, userID string) (float64, error) {
	// Get all open positions for the user
	positions, err := s.positionRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get positions: %w", err)
	}

	var totalExposure float64

	// Calculate total exposure across all positions
	for _, position := range positions {
		ticker, err := s.marketDataService.GetTicker(ctx, position.Symbol)
		if err != nil {
			s.logger.Warn().Err(err).Str("symbol", position.Symbol).Msg("Failed to get price for position")
			continue
		}

		positionValue := position.Quantity * ticker.Price
		totalExposure += positionValue
	}

	return totalExposure, nil
}

// calculateTotalPortfolioValue calculates the total value of a user's portfolio
func (s *RiskService) calculateTotalPortfolioValue(ctx context.Context, userID string) (float64, error) {
	// Get wallet for the user
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	var totalValue float64

	// Add all balances converted to quote currency
	for asset, balance := range wallet.Balances {
		// Skip assets with zero balance
		if balance.Total <= 0 {
			continue
		}

		// For quote currency (e.g., USDT, USD), add directly
		if asset == model.Asset("USDT") || asset == model.Asset("USD") {
			totalValue += balance.Total
			continue
		}

		// For other assets, convert to quote currency using current market price
		symbol := string(asset) + "USDT" // Assuming USDT as quote currency
		ticker, err := s.marketDataService.GetTicker(ctx, symbol)
		if err != nil {
			s.logger.Warn().Err(err).Str("asset", string(asset)).Msg("Failed to get price for asset")
			continue
		}

		assetValue := balance.Total * ticker.Price
		totalValue += assetValue
	}

	// Get all open positions
	positions, err := s.positionRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return totalValue, fmt.Errorf("failed to get positions: %w", err)
	}

	// Add the current value of all open positions
	for _, position := range positions {
		ticker, err := s.marketDataService.GetTicker(ctx, position.Symbol)
		if err != nil {
			s.logger.Warn().Err(err).Str("symbol", position.Symbol).Msg("Failed to get price for position")
			continue
		}

		positionValue := position.Quantity * ticker.Price
		totalValue += positionValue
	}

	return totalValue, nil
}
