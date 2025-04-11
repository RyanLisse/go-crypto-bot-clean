// Package risk provides risk management utilities including risk limits and exposure tracking.
package risk

import (
	"errors"
	"fmt"
	"sync"
)

// RiskLimits defines the configurable risk parameters.
type RiskLimits struct {
	MaxRiskPerTradePercent float64 // e.g., 1.0 means 1% of account balance
	MaxTotalExposure       float64 // total absolute exposure allowed (e.g., in USD)
}

// Validate checks if the risk limits are valid (positive values).
func (rl *RiskLimits) Validate() error {
	if rl.MaxRiskPerTradePercent <= 0 {
		return errors.New("MaxRiskPerTradePercent must be positive")
	}
	if rl.MaxTotalExposure <= 0 {
		return errors.New("MaxTotalExposure must be positive")
	}
	return nil
}

// Position represents an open position.
type Position struct {
	Symbol string
	Size   float64 // number of units/contracts
	Value  float64 // absolute USD value of the position
}

// ExposureTracker tracks open positions and checks risk limits.
type ExposureTracker struct {
	mu             sync.Mutex
	positions      map[string]Position
	RiskLimits     RiskLimits
	AccountBalance float64 // current account balance in USD
}

// NewExposureTracker creates a new tracker with given limits and account balance.
func NewExposureTracker(limits RiskLimits, accountBalance float64) (*ExposureTracker, error) {
	if err := limits.Validate(); err != nil {
		return nil, err
	}
	if accountBalance <= 0 {
		return nil, errors.New("Account balance must be positive")
	}
	return &ExposureTracker{
		positions:      make(map[string]Position),
		RiskLimits:     limits,
		AccountBalance: accountBalance,
	}, nil
}

// AddOrUpdatePosition adds or updates an open position.
func (et *ExposureTracker) AddOrUpdatePosition(pos Position) {
	et.mu.Lock()
	defer et.mu.Unlock()
	et.positions[pos.Symbol] = pos
}

// RemovePosition removes a position by symbol.
func (et *ExposureTracker) RemovePosition(symbol string) {
	et.mu.Lock()
	defer et.mu.Unlock()
	delete(et.positions, symbol)
}

// TotalExposure calculates the sum of absolute values of all open positions.
func (et *ExposureTracker) TotalExposure() float64 {
	et.mu.Lock()
	defer et.mu.Unlock()
	var total float64
	for _, pos := range et.positions {
		total += pos.Value
	}
	return total
}

// CanOpenTrade checks if opening a new trade with given position size and value would breach limits.
// positionValue is the absolute USD value of the new position (from position sizing calculator).
func (et *ExposureTracker) CanOpenTrade(positionValue float64) (bool, error) {
	if positionValue <= 0 {
		return false, errors.New("Position value must be positive")
	}

	// Check per-trade risk limit
	maxRiskValue := et.AccountBalance * et.RiskLimits.MaxRiskPerTradePercent / 100.0
	if positionValue > maxRiskValue {
		return false, fmt.Errorf("Trade risk %.2f exceeds max per-trade risk %.2f", positionValue, maxRiskValue)
	}

	// Check total exposure limit
	totalExposure := et.TotalExposure()
	if totalExposure+positionValue > et.RiskLimits.MaxTotalExposure {
		return false, fmt.Errorf("Total exposure %.2f would exceed max total exposure %.2f", totalExposure+positionValue, et.RiskLimits.MaxTotalExposure)
	}

	return true, nil
}
