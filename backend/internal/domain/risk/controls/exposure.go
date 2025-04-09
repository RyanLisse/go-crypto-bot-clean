package controls

import (
	"context"
	"fmt"
)

// ExposureMonitor tracks and limits total exposure
type ExposureMonitor struct {
	positionRepo PositionRepository
	accountSvc   AccountService
	logger       Logger
}

// GetAccountBalance returns the current account balance
func (em *ExposureMonitor) GetAccountBalance(ctx context.Context) (float64, error) {
	return em.accountSvc.GetBalance(ctx)
}

// PositionRepository defines the interface for accessing position data
type PositionRepository interface {
	GetOpenPositions(ctx context.Context) ([]Position, error)
}

// Position represents a trading position
type Position struct {
	Symbol     string
	Quantity   float64
	EntryPrice float64
}

// AccountService defines the interface for accessing account information
type AccountService interface {
	GetBalance(ctx context.Context) (float64, error)
}

// NewExposureMonitor creates a new ExposureMonitor
func NewExposureMonitor(
	positionRepo PositionRepository,
	accountSvc AccountService,
	logger Logger,
) *ExposureMonitor {
	return &ExposureMonitor{
		positionRepo: positionRepo,
		accountSvc:   accountSvc,
		logger:       logger,
	}
}

// CalculateTotalExposure calculates the total value of all open positions
func (em *ExposureMonitor) CalculateTotalExposure(ctx context.Context) (float64, error) {
	// Get open positions
	positions, err := em.positionRepo.GetOpenPositions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get open positions: %w", err)
	}

	// Calculate total exposure
	var totalExposure float64
	for _, pos := range positions {
		totalExposure += pos.Quantity * pos.EntryPrice
	}

	em.logger.Info("Calculated total exposure",
		"total_exposure", totalExposure,
		"position_count", len(positions))

	return totalExposure, nil
}

// CalculateExposurePercent calculates the percentage of account balance in open positions
func (em *ExposureMonitor) CalculateExposurePercent(ctx context.Context) (float64, error) {
	// Get total exposure
	totalExposure, err := em.CalculateTotalExposure(ctx)
	if err != nil {
		return 0, err
	}

	// Get account balance
	accountBalance, err := em.accountSvc.GetBalance(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get account balance: %w", err)
	}

	if accountBalance <= 0 {
		return 0, fmt.Errorf("invalid account balance: %f", accountBalance)
	}

	// Calculate exposure percentage
	exposurePercent := (totalExposure / accountBalance) * 100

	em.logger.Info("Calculated exposure percentage",
		"total_exposure", totalExposure,
		"account_balance", accountBalance,
		"exposure_percent", exposurePercent)

	return exposurePercent, nil
}

// CheckExposureLimit verifies if a new order would exceed exposure limits
func (em *ExposureMonitor) CheckExposureLimit(
	ctx context.Context,
	newOrderValue float64,
	maxExposurePercent float64,
) (bool, error) {
	// Get current account balance
	accountBalance, err := em.accountSvc.GetBalance(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get account balance: %w", err)
	}

	// Get current total exposure
	totalExposure, err := em.CalculateTotalExposure(ctx)
	if err != nil {
		return false, err
	}

	// Add new order value
	potentialExposure := totalExposure + newOrderValue

	// Calculate maximum allowed exposure
	maxExposure := accountBalance * (maxExposurePercent / 100)

	// Check if new total exposure exceeds limit
	allowed := potentialExposure <= maxExposure

	if !allowed {
		em.logger.Warn("Order rejected due to exposure limit",
			"current_exposure", totalExposure,
			"new_order_value", newOrderValue,
			"potential_exposure", potentialExposure,
			"max_allowed", maxExposure)
	}

	return allowed, nil
}
