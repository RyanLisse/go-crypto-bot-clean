package advanced

import (
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// ParameterRange defines the range for a parameter
type ParameterRange struct {
	Min  float64
	Max  float64
	Step float64
	Type string // "float", "int", "bool"
	Name string
}

// ParameterSet represents a set of parameters
type ParameterSet map[string]interface{}

// OptimizationResult contains the result of parameter optimization
type OptimizationResult struct {
	BestParameters ParameterSet
	Performance    float64
	AllResults     []struct {
		Parameters  ParameterSet
		Performance float64
	}
	Iterations    int
	ExecutionTime time.Duration
}

// OptimizeParameters optimizes strategy parameters using walk-forward optimization
func OptimizeParameters(
	historicalData []*models.Candle,
	parameterRanges map[string]ParameterRange,
	evaluateFunc func(ParameterSet, []*models.Candle) (float64, error),
	iterations int,
) (*OptimizationResult, error) {
	if len(historicalData) < 100 {
		return nil, errors.New("not enough historical data for optimization")
	}

	startTime := time.Now()

	// Initialize with random parameters
	bestParams := generateRandomParameters(parameterRanges)
	bestPerformance, err := evaluateFunc(bestParams, historicalData)
	if err != nil {
		return nil, err
	}

	// Store all results
	allResults := []struct {
		Parameters  ParameterSet
		Performance float64
	}{
		{
			Parameters:  copyParameterSet(bestParams),
			Performance: bestPerformance,
		},
	}

	// Perform optimization iterations
	for i := 0; i < iterations; i++ {
		// Generate new parameters by mutating the best ones
		newParams := mutateParameters(bestParams, parameterRanges, float64(i)/float64(iterations))

		// Evaluate new parameters
		performance, err := evaluateFunc(newParams, historicalData)
		if err != nil {
			continue
		}

		// Store result
		allResults = append(allResults, struct {
			Parameters  ParameterSet
			Performance float64
		}{
			Parameters:  copyParameterSet(newParams),
			Performance: performance,
		})

		// Update best parameters if better
		if performance > bestPerformance {
			bestPerformance = performance
			bestParams = copyParameterSet(newParams)
		}
	}

	// Sort results by performance
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Performance > allResults[j].Performance
	})

	return &OptimizationResult{
		BestParameters: bestParams,
		Performance:    bestPerformance,
		AllResults:     allResults,
		Iterations:     iterations,
		ExecutionTime:  time.Since(startTime),
	}, nil
}

// generateRandomParameters generates random parameters within the specified ranges
func generateRandomParameters(ranges map[string]ParameterRange) ParameterSet {
	params := make(ParameterSet)
	for name, paramRange := range ranges {
		switch paramRange.Type {
		case "float":
			params[name] = paramRange.Min + rand.Float64()*(paramRange.Max-paramRange.Min)
		case "int":
			min := int(paramRange.Min)
			max := int(paramRange.Max)
			params[name] = min + rand.Intn(max-min+1)
		case "bool":
			params[name] = rand.Float64() < 0.5
		}
	}
	return params
}

// mutateParameters mutates parameters with decreasing mutation rate
func mutateParameters(params ParameterSet, ranges map[string]ParameterRange, progress float64) ParameterSet {
	newParams := copyParameterSet(params)

	// Calculate mutation rate (decreases as optimization progresses)
	mutationRate := 0.5 * (1.0 - progress)

	for name, paramRange := range ranges {
		// Decide whether to mutate this parameter
		if rand.Float64() < mutationRate {
			switch paramRange.Type {
			case "float":
				// Calculate mutation amount (smaller as we progress)
				mutationAmount := (paramRange.Max - paramRange.Min) * mutationRate * rand.Float64()

				// Get current value
				currentValue := params[name].(float64)

				// Apply mutation (50% chance of increasing or decreasing)
				if rand.Float64() < 0.5 {
					currentValue += mutationAmount
				} else {
					currentValue -= mutationAmount
				}

				// Ensure value stays within range
				currentValue = math.Max(paramRange.Min, math.Min(paramRange.Max, currentValue))

				// Round to step if needed
				if paramRange.Step > 0 {
					currentValue = math.Round(currentValue/paramRange.Step) * paramRange.Step
				}

				newParams[name] = currentValue

			case "int":
				min := int(paramRange.Min)
				max := int(paramRange.Max)
				step := int(paramRange.Step)
				if step < 1 {
					step = 1
				}

				// Get current value
				currentValue := params[name].(int)

				// Apply mutation (50% chance of increasing or decreasing)
				if rand.Float64() < 0.5 {
					currentValue += step
				} else {
					currentValue -= step
				}

				// Ensure value stays within range
				currentValue = int(math.Max(float64(min), math.Min(float64(max), float64(currentValue))))

				newParams[name] = currentValue

			case "bool":
				// 20% chance of flipping boolean
				if rand.Float64() < 0.2 {
					newParams[name] = !params[name].(bool)
				}
			}
		}
	}

	return newParams
}

// copyParameterSet creates a deep copy of a parameter set
func copyParameterSet(params ParameterSet) ParameterSet {
	newParams := make(ParameterSet)
	for k, v := range params {
		newParams[k] = v
	}
	return newParams
}

// AdaptParametersToRegime adapts strategy parameters based on market regime
func AdaptParametersToRegime(
	baseParams ParameterSet,
	regime MarketRegime,
) ParameterSet {
	newParams := copyParameterSet(baseParams)

	switch regime {
	case RegimeTrendingUp, RegimeTrendingDown:
		// In trending markets, increase trend-following parameters
		if period, ok := newParams["macdFastPeriod"].(int); ok {
			newParams["macdFastPeriod"] = int(float64(period) * 0.8) // Faster MACD
		}
		if period, ok := newParams["emaPeriod"].(int); ok {
			newParams["emaPeriod"] = int(float64(period) * 0.9) // Faster EMA
		}
		if threshold, ok := newParams["profitTarget"].(float64); ok {
			newParams["profitTarget"] = threshold * 1.2 // Higher profit target
		}
		if threshold, ok := newParams["stopLoss"].(float64); ok {
			newParams["stopLoss"] = threshold * 0.8 // Tighter stop loss
		}

	case RegimeRanging, RegimeChoppy:
		// In ranging markets, adjust for mean reversion
		if period, ok := newParams["rsiPeriod"].(int); ok {
			newParams["rsiPeriod"] = int(float64(period) * 0.8) // More responsive RSI
		}
		if threshold, ok := newParams["rsiOverbought"].(float64); ok {
			newParams["rsiOverbought"] = threshold * 0.95 // Lower overbought threshold
		}
		if threshold, ok := newParams["rsiOversold"].(float64); ok {
			newParams["rsiOversold"] = threshold * 1.05 // Higher oversold threshold
		}
		if threshold, ok := newParams["profitTarget"].(float64); ok {
			newParams["profitTarget"] = threshold * 0.8 // Lower profit target
		}

	case RegimeVolatile, RegimeBreakout, RegimeBreakdown:
		// In volatile markets, adjust for quick moves
		if period, ok := newParams["atrPeriod"].(int); ok {
			newParams["atrPeriod"] = int(float64(period) * 0.7) // More responsive ATR
		}
		if multiplier, ok := newParams["atrMultiplier"].(float64); ok {
			newParams["atrMultiplier"] = multiplier * 1.3 // Wider stops based on ATR
		}
		if threshold, ok := newParams["profitTarget"].(float64); ok {
			newParams["profitTarget"] = threshold * 1.3 // Higher profit target
		}
		if threshold, ok := newParams["stopLoss"].(float64); ok {
			newParams["stopLoss"] = threshold * 1.2 // Wider stop loss
		}
	}

	return newParams
}

// EvaluateStrategyPerformance evaluates strategy performance on historical data
func EvaluateStrategyPerformance(
	params ParameterSet,
	historicalData []*models.Candle,
	generateSignalFunc func(ParameterSet, *models.Candle) (*Signal, error),
) (float64, error) {
	if len(historicalData) < 10 {
		return 0, errors.New("not enough historical data for evaluation")
	}

	// Initialize performance metrics
	initialBalance := 1000.0
	balance := initialBalance
	position := 0.0
	entryPrice := 0.0
	trades := 0
	winningTrades := 0
	losingTrades := 0

	// Process each candle
	for i, candle := range historicalData {
		if i == 0 {
			continue // Skip first candle
		}

		// Generate signal for this candle
		signal, err := generateSignalFunc(params, candle)
		if err != nil {
			continue
		}

		// Process signal
		if signal != nil {
			switch signal.Type {
			case SignalBuy:
				if position == 0 {
					// Enter position
					positionSize := balance * signal.RecommendedSize
					if positionSize <= 0 {
						positionSize = balance * 0.1 // Default to 10% if not specified
					}
					position = positionSize / candle.ClosePrice
					entryPrice = candle.ClosePrice
					balance -= positionSize
				}

			case SignalSell, SignalClose:
				if position > 0 {
					// Exit position
					exitPrice := candle.ClosePrice
					positionValue := position * exitPrice
					pnl := positionValue - (position * entryPrice)
					balance += positionValue

					// Record trade result
					trades++
					if pnl > 0 {
						winningTrades++
					} else {
						losingTrades++
					}

					position = 0
					entryPrice = 0
				}
			}
		}

		// Apply stop loss and take profit if position is open
		if position > 0 {
			// Check for stop loss
			stopLossPercent := 0.05 // Default 5%
			if sl, ok := params["stopLoss"].(float64); ok {
				stopLossPercent = sl
			}
			stopLossPrice := entryPrice * (1 - stopLossPercent)

			// Check for take profit
			takeProfitPercent := 0.1 // Default 10%
			if tp, ok := params["profitTarget"].(float64); ok {
				takeProfitPercent = tp
			}
			takeProfitPrice := entryPrice * (1 + takeProfitPercent)

			// Check if stop loss or take profit was hit
			if candle.LowPrice <= stopLossPrice {
				// Stop loss hit
				positionValue := position * stopLossPrice
				// Calculate PnL for logging/metrics if needed
				_ = positionValue - (position * entryPrice)
				balance += positionValue

				trades++
				losingTrades++

				position = 0
				entryPrice = 0
			} else if candle.HighPrice >= takeProfitPrice {
				// Take profit hit
				positionValue := position * takeProfitPrice
				// Calculate PnL for logging/metrics if needed
				_ = positionValue - (position * entryPrice)
				balance += positionValue

				trades++
				winningTrades++

				position = 0
				entryPrice = 0
			}
		}
	}

	// Close any remaining position at the last price
	if position > 0 {
		lastPrice := historicalData[len(historicalData)-1].ClosePrice
		positionValue := position * lastPrice
		balance += positionValue
	}

	// Calculate performance metrics
	profitLoss := balance - initialBalance
	profitLossPercent := profitLoss / initialBalance * 100

	// Calculate win rate
	winRate := 0.0
	if trades > 0 {
		winRate = float64(winningTrades) / float64(trades)
	}

	// Calculate performance score (combination of profit and win rate)
	performanceScore := profitLossPercent * winRate

	return performanceScore, nil
}
