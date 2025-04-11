// Package risk provides risk management utilities such as position sizing calculations.
package risk

import (
	"errors"
)

// PositionSize calculates the appropriate position size based on account balance,
// risk percentage, and market volatility.
//
// Parameters:
//   - accountBalance: total account balance (must be > 0)
//   - riskPercentage: fraction of account balance to risk per trade (e.g., 0.01 for 1%, must be > 0)
//   - marketVolatility: a volatility measure such as ATR or standard deviation (must be > 0)
//
// Returns:
//   - positionSize: the calculated position size
//   - error: if any input is zero or negative
//
// Formula:
//
//	riskAmount = accountBalance * riskPercentage
//	positionSize = riskAmount / marketVolatility
func PositionSize(accountBalance, riskPercentage, marketVolatility float64) (float64, error) {
	if accountBalance <= 0 {
		return 0, errors.New("accountBalance must be positive")
	}
	if riskPercentage <= 0 {
		return 0, errors.New("riskPercentage must be positive")
	}
	if marketVolatility <= 0 {
		return 0, errors.New("marketVolatility must be positive")
	}

	riskAmount := accountBalance * riskPercentage
	positionSize := riskAmount / marketVolatility
	return positionSize, nil
}
