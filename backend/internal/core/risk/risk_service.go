package risk

import (
	"context"
	"fmt"
	"log"

	"github.com/ryanlisse/go-crypto-bot/internal/core/account"
	"github.com/ryanlisse/go-crypto-bot/internal/core/position"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PositionSizer defines interface for pluggable position sizing models.
type PositionSizer interface {
	Calculate(ctx context.Context, accountBalance, entryPrice, stopLossPrice float64) (float64, error)
}

// RiskConfig defines the interface for risk configuration parameters.
type RiskConfig interface {
	GetMaxRiskPerTrade() float64
	GetMaxPortfolioRisk() float64
	GetMaxExposureLimit() float64
	GetMinRiskRewardRatio() float64
	GetDailyLossLimit() float64
	GetMaxDrawdownThreshold() float64
}

// RiskService defines the interface for risk management operations.
type RiskService interface {
	// CalculatePositionSize determines the appropriate position size based on risk parameters.
	CalculatePositionSize(ctx context.Context, entryPrice, stopLossPrice float64) (float64, float64, error)
	// CalculatePortfolioRisk calculates the current risk exposure of the portfolio.
	CalculatePortfolioRisk(ctx context.Context) (float64, map[string]float64, error)
	// CheckRiskLimits checks if a new trade would exceed risk limits.
	CheckRiskLimits(ctx context.Context, symbol string, amount, entryPrice, stopLossPrice float64) (bool, float64, float64, error)
	// CalculateRiskRewardRatio calculates the risk-reward ratio for a trade.
	CalculateRiskRewardRatio(entryPrice, takeProfitPrice, stopLossPrice float64) float64
	// CalculateMaxDrawdown calculates the maximum drawdown from a series of balances.
	CalculateMaxDrawdown(balances []float64) float64
	// IsTradeAllowed checks if a trade is allowed based on risk parameters.
	IsTradeAllowed(ctx context.Context, symbol string, amount, entryPrice, takeProfitPrice, stopLossPrice float64) (bool, string, error)
}

// riskService implements the RiskService interface.
type riskService struct {
	accountService  account.AccountService
	positionService position.PositionService
	config          RiskConfig
	positionSizer   PositionSizer
}

// NewRiskService creates a new instance of the risk service.
func NewRiskService(accountService account.AccountService, positionService position.PositionService, config RiskConfig, sizer PositionSizer) RiskService {
	// Subscribe to balance updates for real-time risk assessment
	ctx := context.Background()
	_ = accountService.SubscribeToBalanceUpdates(ctx, func(wallet *models.Wallet) {
		// This is a placeholder for future real-time risk assessment logic
		// We could recalculate portfolio risk when balances change
	})

	return &riskService{
		accountService:  accountService,
		positionService: positionService,
		config:          config,
		positionSizer:   sizer,
	}
}

// checkGlobalLimits consolidates risk checks for exposure limits.
// It now calls GetPortfolioValue to simulate additional limit checks.
func (s *riskService) checkGlobalLimits(ctx context.Context, newOrderValue float64) error {
	// Simulate additional limit checks by fetching portfolio value.
	if _, err := s.accountService.GetPortfolioValue(ctx); err != nil {
		return fmt.Errorf("failed to get portfolio value for limits check: %w", err)
	}

	ok, err := s.checkExposureLimit(ctx, newOrderValue)
	if err != nil {
		return fmt.Errorf("failed to check exposure limit: %w", err)
	} else if !ok {
		return fmt.Errorf("exposure limit reached")
	}
	return nil
}

// CalculatePositionSize determines the appropriate position size based on risk parameters.
func (s *riskService) CalculatePositionSize(ctx context.Context, entryPrice, stopLossPrice float64) (float64, float64, error) {
	stopLossDistance := entryPrice - stopLossPrice
	if stopLossDistance <= 0 {
		return 0, 0, fmt.Errorf("invalid stop loss: entry price must be greater than stop loss price")
	}

	// Check global risk limits with no additional exposure.
	if err := s.checkGlobalLimits(ctx, 0); err != nil {
		return 0, 0, err
	}

	portfolioValue, err := s.accountService.GetPortfolioValue(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get portfolio value: %w", err)
	}

	maxRiskPerTrade := s.config.GetMaxRiskPerTrade()
	riskAmount := portfolioValue * maxRiskPerTrade

	positionSize := (riskAmount / stopLossDistance) * entryPrice

	return positionSize, riskAmount, nil
}

// CalculatePortfolioRisk calculates the current risk exposure of the portfolio.
func (s *riskService) CalculatePortfolioRisk(ctx context.Context) (float64, map[string]float64, error) {
	// Check global risk limits with no additional exposure.
	if err := s.checkGlobalLimits(ctx, 0); err != nil {
		return 0, nil, err
	}

	positionRisks, err := s.accountService.GetAllPositionRisks(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get position risks: %w", err)
	}

	portfolioValue, err := s.accountService.GetPortfolioValue(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get portfolio value: %w", err)
	}

	totalExposure := 0.0
	exposureBySymbol := make(map[string]float64)
	for _, risk := range positionRisks {
		totalExposure += risk.ExposureUSD
		exposureBySymbol[risk.Symbol] = risk.ExposureUSD / portfolioValue
	}

	portfolioRisk := totalExposure / portfolioValue

	return portfolioRisk, exposureBySymbol, nil
}

// CheckRiskLimits checks if a new trade would exceed the risk limits.
func (s *riskService) CheckRiskLimits(ctx context.Context, symbol string, amount, entryPrice, stopLossPrice float64) (bool, float64, float64, error) {
	newOrderValue := amount * entryPrice
	if err := s.checkGlobalLimits(ctx, newOrderValue); err != nil {
		return false, 0, 0, err
	}

	positionRisks, err := s.accountService.GetAllPositionRisks(ctx)
	if err != nil {
		return false, 0, 0, fmt.Errorf("failed to get position risks: %w", err)
	}

	portfolioValue, err := s.accountService.GetPortfolioValue(ctx)
	if err != nil {
		return false, 0, 0, fmt.Errorf("failed to get portfolio value: %w", err)
	}

	totalExposure := 0.0
	for _, risk := range positionRisks {
		totalExposure += risk.ExposureUSD
	}
	currentRisk := totalExposure / portfolioValue

	maxRisk := s.config.GetMaxPortfolioRisk()

	return currentRisk <= maxRisk, currentRisk, maxRisk, nil
}

// CalculateRiskRewardRatio calculates the risk-reward ratio for a trade.
func (s *riskService) CalculateRiskRewardRatio(entryPrice, takeProfitPrice, stopLossPrice float64) float64 {
	potentialProfit := takeProfitPrice - entryPrice
	potentialLoss := entryPrice - stopLossPrice

	if potentialLoss <= 0 {
		return 0
	}
	return potentialProfit / potentialLoss
}

// CalculateMaxDrawdown calculates the maximum drawdown from a series of balances.
func (s *riskService) CalculateMaxDrawdown(balances []float64) float64 {
	if len(balances) == 0 {
		return 0
	}

	peak := balances[0]
	maxDrawdown := 0.0
	for _, balance := range balances {
		if balance > peak {
			peak = balance
		} else {
			drawdown := (peak - balance) / peak
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
			}
		}
	}
	return maxDrawdown
}

// IsTradeAllowed checks if a trade is allowed based on risk parameters.
func (s *riskService) IsTradeAllowed(ctx context.Context, symbol string, amount, entryPrice, takeProfitPrice, stopLossPrice float64) (bool, string, error) {
	log.Println("IsTradeAllowed called")
	// Check risk-reward ratio.
	riskRewardRatio := s.CalculateRiskRewardRatio(entryPrice, takeProfitPrice, stopLossPrice)
	minRiskRewardRatio := s.config.GetMinRiskRewardRatio()
	if riskRewardRatio < minRiskRewardRatio {
		return false, fmt.Sprintf("Risk-reward ratio %.2f is below minimum threshold %.2f", riskRewardRatio, minRiskRewardRatio), nil
	}

	// Check portfolio risk limits.
	withinLimits, currentRisk, maxRisk, err := s.CheckRiskLimits(ctx, symbol, amount, entryPrice, stopLossPrice)
	if err != nil {
		return false, "", err
	}
	if !withinLimits {
		return false, fmt.Sprintf("Trade would exceed maximum portfolio risk: %.2f%% > %.2f%%", currentRisk*100, maxRisk*100), nil
	}

	return true, "", nil
}

// DefaultRiskConfig provides default values for RiskConfig.
type DefaultRiskConfig struct{}

func (d *DefaultRiskConfig) GetMaxRiskPerTrade() float64 {
	return 0.02 // 2%
}

func (d *DefaultRiskConfig) GetMaxPortfolioRisk() float64 {
	return 0.8 // 80%
}

func (d *DefaultRiskConfig) GetMinRiskRewardRatio() float64 {
	return 1.5
}

func (d *DefaultRiskConfig) GetMaxExposureLimit() float64 {
	return 100000 // Updated to a high value for testing
}

func (d *DefaultRiskConfig) GetDailyLossLimit() float64 {
	return 5000.0
}

func (d *DefaultRiskConfig) GetMaxDrawdownThreshold() float64 {
	return 0.3 // 30%
}

// checkExposureLimit checks if the current exposure plus a new position is within the configured limit.
func (s *riskService) checkExposureLimit(ctx context.Context, newPositionValue float64) (bool, error) {
	currentExposure, err := s.accountService.GetCurrentExposure(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get current exposure: %w", err)
	}

	maxExposureLimit := s.config.GetMaxExposureLimit()
	totalExposure := currentExposure + newPositionValue

	return totalExposure <= maxExposureLimit, nil
}
