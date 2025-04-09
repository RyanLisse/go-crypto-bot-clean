package advanced

import (
	"errors"
	"sort"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// TimeframeMode represents a combination of timeframes to use
type TimeframeMode string

const (
	TimeframeModeShort   TimeframeMode = "SHORT"   // Short-term (1m, 5m, 15m)
	TimeframeModeMedium  TimeframeMode = "MEDIUM"  // Medium-term (15m, 1h, 4h)
	TimeframeModeLong    TimeframeMode = "LONG"    // Long-term (4h, 1d, 1w)
	TimeframeModeAll     TimeframeMode = "ALL"     // All timeframes
	TimeframeModeAdaptive TimeframeMode = "ADAPTIVE" // Adaptive based on market conditions
)

// TimeframeConfig defines the configuration for a timeframe mode
type TimeframeConfig struct {
	Mode       TimeframeMode
	Timeframes []string
	Weights    map[string]float64
}

// GetTimeframeConfig returns the configuration for a timeframe mode
func GetTimeframeConfig(mode TimeframeMode) TimeframeConfig {
	switch mode {
	case TimeframeModeShort:
		return TimeframeConfig{
			Mode:       TimeframeModeShort,
			Timeframes: []string{"1m", "5m", "15m"},
			Weights: map[string]float64{
				"1m":  0.2,
				"5m":  0.3,
				"15m": 0.5,
			},
		}
	case TimeframeModeMedium:
		return TimeframeConfig{
			Mode:       TimeframeModeMedium,
			Timeframes: []string{"15m", "1h", "4h"},
			Weights: map[string]float64{
				"15m": 0.2,
				"1h":  0.5,
				"4h":  0.3,
			},
		}
	case TimeframeModeLong:
		return TimeframeConfig{
			Mode:       TimeframeModeLong,
			Timeframes: []string{"4h", "1d", "1w"},
			Weights: map[string]float64{
				"4h": 0.3,
				"1d": 0.5,
				"1w": 0.2,
			},
		}
	case TimeframeModeAll:
		return TimeframeConfig{
			Mode:       TimeframeModeAll,
			Timeframes: []string{"1m", "5m", "15m", "1h", "4h", "1d", "1w"},
			Weights: map[string]float64{
				"1m":  0.05,
				"5m":  0.1,
				"15m": 0.15,
				"1h":  0.2,
				"4h":  0.25,
				"1d":  0.15,
				"1w":  0.1,
			},
		}
	case TimeframeModeAdaptive:
		return TimeframeConfig{
			Mode:       TimeframeModeAdaptive,
			Timeframes: []string{"5m", "15m", "1h", "4h"},
			Weights: map[string]float64{
				"5m":  0.25,
				"15m": 0.25,
				"1h":  0.25,
				"4h":  0.25,
			},
		}
	default:
		return TimeframeConfig{
			Mode:       TimeframeModeMedium,
			Timeframes: []string{"15m", "1h", "4h"},
			Weights: map[string]float64{
				"15m": 0.2,
				"1h":  0.5,
				"4h":  0.3,
			},
		}
	}
}

// TimeframeAnalysisResult contains the result of multi-timeframe analysis
type TimeframeAnalysisResult struct {
	Mode            TimeframeMode
	Signals         map[string]*Signal
	CombinedSignal  *Signal
	AlignmentScore  float64
	Timeframes      []string
	Weights         map[string]float64
}

// PerformMultiTimeframeAnalysis analyzes multiple timeframes and combines the results
func PerformMultiTimeframeAnalysis(
	candlesByTimeframe map[string][]*models.Candle,
	mode TimeframeMode,
	analyzeFunc func([]*models.Candle) (*Signal, error),
) (*TimeframeAnalysisResult, error) {
	// Get timeframe configuration
	config := GetTimeframeConfig(mode)

	// Check if we have data for all required timeframes
	for _, tf := range config.Timeframes {
		if _, ok := candlesByTimeframe[tf]; !ok {
			return nil, errors.New("missing data for timeframe: " + tf)
		}
	}

	// Analyze each timeframe
	signals := make(map[string]*Signal)
	for _, tf := range config.Timeframes {
		signal, err := analyzeFunc(candlesByTimeframe[tf])
		if err != nil {
			return nil, err
		}
		signal.Timeframe = tf
		signals[tf] = signal
	}

	// Combine signals
	combinedSignal, alignmentScore := combineTimeframeSignals(signals, config.Weights)

	return &TimeframeAnalysisResult{
		Mode:            mode,
		Signals:         signals,
		CombinedSignal:  combinedSignal,
		AlignmentScore:  alignmentScore,
		Timeframes:      config.Timeframes,
		Weights:         config.Weights,
	}, nil
}

// combineTimeframeSignals combines signals from multiple timeframes
func combineTimeframeSignals(signals map[string]*Signal, weights map[string]float64) (*Signal, float64) {
	if len(signals) == 0 {
		return nil, 0
	}

	// Get a sample signal to copy basic properties
	var sampleSignal *Signal
	for _, signal := range signals {
		sampleSignal = signal
		break
	}

	// Count signal types
	signalTypeCounts := make(map[SignalType]float64)
	signalTypeWeights := make(map[SignalType]float64)
	totalWeight := 0.0

	for tf, signal := range signals {
		weight := weights[tf]
		signalTypeCounts[signal.Type] += weight
		signalTypeWeights[signal.Type] += weight * signal.Confidence
		totalWeight += weight
	}

	// Normalize weights
	if totalWeight > 0 {
		for signalType := range signalTypeCounts {
			signalTypeCounts[signalType] /= totalWeight
			signalTypeWeights[signalType] /= totalWeight
		}
	}

	// Determine the dominant signal type
	var dominantType SignalType
	var maxCount float64
	for signalType, count := range signalTypeCounts {
		if count > maxCount {
			maxCount = count
			dominantType = signalType
		}
	}

	// Calculate alignment score (how aligned are the signals across timeframes)
	alignmentScore := maxCount

	// Create combined signal
	combinedSignal := &Signal{
		Symbol:          sampleSignal.Symbol,
		Type:            dominantType,
		Confidence:      signalTypeWeights[dominantType],
		Price:           sampleSignal.Price,
		Timestamp:       time.Now(),
		ExpirationTime:  time.Now().Add(1 * time.Hour),
		RecommendedSize: 0.0,
		Metadata:        make(map[string]interface{}),
	}

	// Calculate average target, stop loss, and take profit
	targetSum := 0.0
	stopLossSum := 0.0
	takeProfitSum := 0.0
	targetCount := 0
	stopLossCount := 0
	takeProfitCount := 0

	for _, signal := range signals {
		if signal.Type == dominantType {
			if signal.TargetPrice > 0 {
				targetSum += signal.TargetPrice
				targetCount++
			}
			if signal.StopLoss > 0 {
				stopLossSum += signal.StopLoss
				stopLossCount++
			}
			if signal.TakeProfit > 0 {
				takeProfitSum += signal.TakeProfit
				takeProfitCount++
			}
			// Adjust recommended size based on confidence
			combinedSignal.RecommendedSize += signal.RecommendedSize * signal.Confidence
		}
	}

	// Set average values
	if targetCount > 0 {
		combinedSignal.TargetPrice = targetSum / float64(targetCount)
	}
	if stopLossCount > 0 {
		combinedSignal.StopLoss = stopLossSum / float64(stopLossCount)
	}
	if takeProfitCount > 0 {
		combinedSignal.TakeProfit = takeProfitSum / float64(takeProfitCount)
	}

	// Normalize recommended size
	if combinedSignal.RecommendedSize > 1.0 {
		combinedSignal.RecommendedSize = 1.0
	}

	// Add metadata
	combinedSignal.Metadata["timeframes"] = getTimeframesSortedByWeight(weights)
	combinedSignal.Metadata["alignmentScore"] = alignmentScore
	combinedSignal.Metadata["signalCounts"] = signalTypeCounts

	return combinedSignal, alignmentScore
}

// getTimeframesSortedByWeight returns timeframes sorted by weight
func getTimeframesSortedByWeight(weights map[string]float64) []string {
	type timeframeWeight struct {
		Timeframe string
		Weight    float64
	}

	// Create slice of timeframe weights
	timeframeWeights := make([]timeframeWeight, 0, len(weights))
	for tf, weight := range weights {
		timeframeWeights = append(timeframeWeights, timeframeWeight{
			Timeframe: tf,
			Weight:    weight,
		})
	}

	// Sort by weight (descending)
	sort.Slice(timeframeWeights, func(i, j int) bool {
		return timeframeWeights[i].Weight > timeframeWeights[j].Weight
	})

	// Extract timeframes
	timeframes := make([]string, len(timeframeWeights))
	for i, tw := range timeframeWeights {
		timeframes[i] = tw.Timeframe
	}

	return timeframes
}

// AdaptTimeframeWeights adapts timeframe weights based on market conditions
func AdaptTimeframeWeights(
	regime MarketRegime,
	config TimeframeConfig,
) map[string]float64 {
	weights := make(map[string]float64)
	for tf, weight := range config.Weights {
		weights[tf] = weight
	}

	switch regime {
	case RegimeTrendingUp, RegimeTrendingDown:
		// In trending markets, give more weight to higher timeframes
		for tf, weight := range weights {
			switch tf {
			case "1m", "5m":
				weights[tf] = weight * 0.5
			case "15m", "1h":
				weights[tf] = weight * 1.0
			case "4h", "1d", "1w":
				weights[tf] = weight * 1.5
			}
		}
	case RegimeRanging, RegimeChoppy:
		// In ranging markets, give more weight to lower timeframes
		for tf, weight := range weights {
			switch tf {
			case "1m", "5m", "15m":
				weights[tf] = weight * 1.5
			case "1h":
				weights[tf] = weight * 1.0
			case "4h", "1d", "1w":
				weights[tf] = weight * 0.5
			}
		}
	case RegimeVolatile, RegimeBreakout, RegimeBreakdown:
		// In volatile markets, balance between timeframes
		for tf, weight := range weights {
			switch tf {
			case "1m":
				weights[tf] = weight * 0.7
			case "5m", "15m", "1h", "4h":
				weights[tf] = weight * 1.2
			case "1d", "1w":
				weights[tf] = weight * 0.8
			}
		}
	}

	// Normalize weights
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}
	if totalWeight > 0 {
		for tf := range weights {
			weights[tf] /= totalWeight
		}
	}

	return weights
}
