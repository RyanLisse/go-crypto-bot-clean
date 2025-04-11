// Package risk provides stop-loss and take-profit management utilities.
package risk

import (
	"errors"
	"math"
)

// RiskParams holds parameters for risk-based SL/TP calculation.
type RiskParams struct {
	RiskPercentage float64 // e.g., 1.0 for 1%
	ATR            float64 // Average True Range, optional, can be zero if unused
	UseATR         bool    // If true, use ATR-based calculation
	RRRatio        float64 // Risk-Reward ratio, e.g., 2.0 means TP is 2x SL distance
}

// SLTPLevel holds stop-loss and take-profit price levels.
type SLTPLevel struct {
	StopLoss   float64
	TakeProfit float64
}

// CalculateSLTP computes stop-loss and take-profit prices based on entry price and risk parameters.
// Returns error if parameters are invalid.
func CalculateSLTP(entryPrice float64, params RiskParams) (SLTPLevel, error) {
	if entryPrice <= 0 {
		return SLTPLevel{}, errors.New("entry price must be positive")
	}
	if params.RiskPercentage <= 0 || params.RiskPercentage >= 100 {
		return SLTPLevel{}, errors.New("risk percentage must be between 0 and 100")
	}
	if params.RRRatio <= 0 {
		return SLTPLevel{}, errors.New("risk-reward ratio must be positive")
	}
	var riskAmount float64
	if params.UseATR {
		if params.ATR <= 0 {
			return SLTPLevel{}, errors.New("ATR must be positive when UseATR is true")
		}
		riskAmount = params.ATR
	} else {
		riskAmount = entryPrice * params.RiskPercentage / 100.0
	}

	stopLoss := entryPrice - riskAmount
	takeProfit := entryPrice + (riskAmount * params.RRRatio)

	if stopLoss <= 0 {
		stopLoss = 0.00000001 // minimal positive price
	}

	return SLTPLevel{
		StopLoss:   stopLoss,
		TakeProfit: takeProfit,
	}, nil
}

// IsStopLossHit returns true if the current price is less than or equal to the stop-loss level.
func IsStopLossHit(currentPrice float64, sltp SLTPLevel) bool {
	return currentPrice <= sltp.StopLoss
}

// IsTakeProfitHit returns true if the current price is greater than or equal to the take-profit level.
func IsTakeProfitHit(currentPrice float64, sltp SLTPLevel) bool {
	return currentPrice >= sltp.TakeProfit
}

// UpdateSLTP recalculates SL/TP based on new risk parameters.
// Returns error if parameters are invalid.
func UpdateSLTP(entryPrice float64, params RiskParams) (SLTPLevel, error) {
	return CalculateSLTP(entryPrice, params)
}

// CancelSLTP returns a zeroed SLTPLevel indicating no active SL/TP.
func CancelSLTP() SLTPLevel {
	return SLTPLevel{
		StopLoss:   math.NaN(),
		TakeProfit: math.NaN(),
	}
}
