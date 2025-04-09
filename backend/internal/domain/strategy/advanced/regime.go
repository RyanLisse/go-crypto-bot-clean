package advanced

import (
	"errors"
	"math"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// MarketRegime represents the current market condition
type MarketRegime string

const (
	RegimeTrendingUp   MarketRegime = "TRENDING_UP"
	RegimeTrendingDown MarketRegime = "TRENDING_DOWN"
	RegimeRanging      MarketRegime = "RANGING"
	RegimeVolatile     MarketRegime = "VOLATILE"
	RegimeChoppy       MarketRegime = "CHOPPY"
	RegimeBreakout     MarketRegime = "BREAKOUT"
	RegimeBreakdown    MarketRegime = "BREAKDOWN"
	RegimeUndetermined MarketRegime = "UNDETERMINED"
)

// RegimeDetectionResult contains the detected market regime and confidence
type RegimeDetectionResult struct {
	Regime     MarketRegime
	Confidence float64
	Metadata   map[string]interface{}
}

// DetectMarketRegime detects the current market regime using multiple methods
func DetectMarketRegime(candles []*models.Candle) (*RegimeDetectionResult, error) {
	if len(candles) < 50 {
		return nil, errors.New("not enough data for regime detection")
	}

	// Use multiple methods to detect regime
	trendResult, err := detectTrendUsingADX(candles)
	if err != nil {
		return nil, err
	}

	volatilityResult, err := detectVolatilityUsingATR(candles)
	if err != nil {
		return nil, err
	}

	rangeResult, err := detectRangeUsingBollingerBands(candles)
	if err != nil {
		return nil, err
	}

	breakoutResult, err := detectBreakoutUsingVolume(candles)
	if err != nil {
		return nil, err
	}

	// Combine results to determine overall regime
	result := combineRegimeResults(trendResult, volatilityResult, rangeResult, breakoutResult)

	return result, nil
}

// detectTrendUsingADX detects trend using Average Directional Index
func detectTrendUsingADX(candles []*models.Candle) (*RegimeDetectionResult, error) {
	// Calculate +DI and -DI
	period := 14
	if len(candles) < period+1 {
		return nil, errors.New("not enough data for ADX calculation")
	}

	// Calculate True Range (TR)
	tr := make([]float64, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		high := candles[i].HighPrice
		low := candles[i].LowPrice
		prevClose := candles[i-1].ClosePrice

		tr[i-1] = math.Max(high-low, math.Max(math.Abs(high-prevClose), math.Abs(low-prevClose)))
	}

	// Calculate +DM and -DM
	plusDM := make([]float64, len(candles)-1)
	minusDM := make([]float64, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		upMove := candles[i].HighPrice - candles[i-1].HighPrice
		downMove := candles[i-1].LowPrice - candles[i].LowPrice

		if upMove > downMove && upMove > 0 {
			plusDM[i-1] = upMove
		} else {
			plusDM[i-1] = 0
		}

		if downMove > upMove && downMove > 0 {
			minusDM[i-1] = downMove
		} else {
			minusDM[i-1] = 0
		}
	}

	// Calculate smoothed TR, +DM, and -DM
	smoothedTR := 0.0
	smoothedPlusDM := 0.0
	smoothedMinusDM := 0.0
	for i := 0; i < period; i++ {
		smoothedTR += tr[i]
		smoothedPlusDM += plusDM[i]
		smoothedMinusDM += minusDM[i]
	}

	// Calculate +DI and -DI
	plusDI := 100 * (smoothedPlusDM / smoothedTR)
	minusDI := 100 * (smoothedMinusDM / smoothedTR)

	// Calculate DX
	dx := 100 * (math.Abs(plusDI-minusDI) / (plusDI + minusDI))

	// Calculate ADX (smoothed DX)
	adx := dx
	for i := period; i < len(tr); i++ {
		// Update smoothed values
		smoothedTR = smoothedTR - (smoothedTR / float64(period)) + tr[i]
		smoothedPlusDM = smoothedPlusDM - (smoothedPlusDM / float64(period)) + plusDM[i]
		smoothedMinusDM = smoothedMinusDM - (smoothedMinusDM / float64(period)) + minusDM[i]

		// Calculate new +DI and -DI
		newPlusDI := 100 * (smoothedPlusDM / smoothedTR)
		newMinusDI := 100 * (smoothedMinusDM / smoothedTR)

		// Calculate new DX
		newDX := 100 * (math.Abs(newPlusDI-newMinusDI) / (newPlusDI + newMinusDI))

		// Smooth ADX
		adx = ((adx * float64(period-1)) + newDX) / float64(period)

		// Update +DI and -DI
		plusDI = newPlusDI
		minusDI = newMinusDI
	}

	// Determine trend based on ADX, +DI, and -DI
	var regime MarketRegime
	var confidence float64

	if adx > 25 {
		if plusDI > minusDI {
			regime = RegimeTrendingUp
			confidence = math.Min(adx/100, 0.95)
		} else {
			regime = RegimeTrendingDown
			confidence = math.Min(adx/100, 0.95)
		}
	} else if adx < 20 {
		regime = RegimeRanging
		confidence = math.Min((20-adx)/20, 0.9)
	} else {
		regime = RegimeUndetermined
		confidence = 0.5
	}

	return &RegimeDetectionResult{
		Regime:     regime,
		Confidence: confidence,
		Metadata: map[string]interface{}{
			"adx":     adx,
			"plusDI":  plusDI,
			"minusDI": minusDI,
		},
	}, nil
}

// detectVolatilityUsingATR detects volatility using Average True Range
func detectVolatilityUsingATR(candles []*models.Candle) (*RegimeDetectionResult, error) {
	period := 14
	if len(candles) < period+1 {
		return nil, errors.New("not enough data for ATR calculation")
	}

	// Calculate True Range (TR)
	tr := make([]float64, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		high := candles[i].HighPrice
		low := candles[i].LowPrice
		prevClose := candles[i-1].ClosePrice

		tr[i-1] = math.Max(high-low, math.Max(math.Abs(high-prevClose), math.Abs(low-prevClose)))
	}

	// Calculate ATR
	atr := 0.0
	for i := 0; i < period; i++ {
		atr += tr[i]
	}
	atr /= float64(period)

	// Smooth ATR
	for i := period; i < len(tr); i++ {
		atr = ((atr * float64(period-1)) + tr[i]) / float64(period)
	}

	// Calculate ATR percentage (ATR/Price)
	atrPercent := atr / candles[len(candles)-1].ClosePrice * 100

	// Determine volatility based on ATR percentage
	var regime MarketRegime
	var confidence float64

	if atrPercent > 5 {
		regime = RegimeVolatile
		confidence = math.Min(atrPercent/10, 0.95)
	} else if atrPercent < 2 {
		regime = RegimeRanging
		confidence = math.Min((2-atrPercent)/2, 0.9)
	} else {
		regime = RegimeUndetermined
		confidence = 0.5
	}

	return &RegimeDetectionResult{
		Regime:     regime,
		Confidence: confidence,
		Metadata: map[string]interface{}{
			"atr":        atr,
			"atrPercent": atrPercent,
		},
	}, nil
}

// detectRangeUsingBollingerBands detects ranging market using Bollinger Bands
func detectRangeUsingBollingerBands(candles []*models.Candle) (*RegimeDetectionResult, error) {
	// Calculate Bollinger Bands
	bb, err := CalculateBollingerBands(candles, 20, 2.0)
	if err != nil {
		return nil, err
	}

	// Get the latest values
	middle := bb.Values["middle"]
	upper := bb.Values["upper"]
	lower := bb.Values["lower"]
	width := bb.Values["width"]

	// Calculate percentage width
	percentWidth := (upper - lower) / middle * 100

	// Check if price is bouncing between bands
	bounceCount := 0
	touchUpperCount := 0
	touchLowerCount := 0
	for i := len(candles) - 20; i < len(candles); i++ {
		if candles[i].HighPrice >= upper*0.98 {
			touchUpperCount++
		}
		if candles[i].LowPrice <= lower*1.02 {
			touchLowerCount++
		}
	}

	if touchUpperCount > 0 && touchLowerCount > 0 {
		bounceCount = int(math.Min(float64(touchUpperCount), float64(touchLowerCount)))
	}

	// Determine regime based on Bollinger Bands
	var regime MarketRegime
	var confidence float64

	if percentWidth < 3 {
		// Tight bands indicate potential breakout
		regime = RegimeChoppy
		confidence = math.Min((3-percentWidth)/3, 0.9)
	} else if percentWidth > 6 && bounceCount >= 2 {
		// Wide bands with bouncing indicates ranging
		regime = RegimeRanging
		confidence = math.Min(float64(bounceCount)/5, 0.9)
	} else if width < 0.03 {
		// Narrow bands indicate low volatility
		regime = RegimeChoppy
		confidence = math.Min((0.03-width)/0.03, 0.9)
	} else {
		regime = RegimeUndetermined
		confidence = 0.5
	}

	return &RegimeDetectionResult{
		Regime:     regime,
		Confidence: confidence,
		Metadata: map[string]interface{}{
			"bbWidth":      width,
			"percentWidth": percentWidth,
			"bounceCount":  bounceCount,
		},
	}, nil
}

// detectBreakoutUsingVolume detects breakouts using volume and price action
func detectBreakoutUsingVolume(candles []*models.Candle) (*RegimeDetectionResult, error) {
	if len(candles) < 20 {
		return nil, errors.New("not enough data for breakout detection")
	}

	// Calculate average volume
	avgVolume := 0.0
	for i := len(candles) - 20; i < len(candles)-1; i++ {
		avgVolume += candles[i].Volume
	}
	avgVolume /= 19 // Last 19 candles (excluding the current one)

	// Get current volume
	currentVolume := candles[len(candles)-1].Volume

	// Calculate volume ratio
	volumeRatio := currentVolume / avgVolume

	// Calculate price change
	priceChange := (candles[len(candles)-1].ClosePrice - candles[len(candles)-2].ClosePrice) / candles[len(candles)-2].ClosePrice * 100

	// Determine if there's a breakout or breakdown
	var regime MarketRegime
	var confidence float64

	if volumeRatio > 2 && priceChange > 3 {
		// High volume with significant price increase
		regime = RegimeBreakout
		confidence = math.Min(volumeRatio/5, 0.95)
	} else if volumeRatio > 2 && priceChange < -3 {
		// High volume with significant price decrease
		regime = RegimeBreakdown
		confidence = math.Min(volumeRatio/5, 0.95)
	} else {
		regime = RegimeUndetermined
		confidence = 0.5
	}

	return &RegimeDetectionResult{
		Regime:     regime,
		Confidence: confidence,
		Metadata: map[string]interface{}{
			"volumeRatio": volumeRatio,
			"priceChange": priceChange,
		},
	}, nil
}

// combineRegimeResults combines multiple regime detection results
func combineRegimeResults(results ...*RegimeDetectionResult) *RegimeDetectionResult {
	// Initialize counters for each regime
	regimeScores := make(map[MarketRegime]float64)
	regimeCount := make(map[MarketRegime]int)

	// Combine results
	for _, result := range results {
		if result != nil {
			regimeScores[result.Regime] += result.Confidence
			regimeCount[result.Regime]++
		}
	}

	// Find the regime with the highest score
	var bestRegime MarketRegime
	var bestScore float64
	for regime, score := range regimeScores {
		if score > bestScore {
			bestRegime = regime
			bestScore = score
		}
	}

	// Calculate confidence
	confidence := 0.0
	if regimeCount[bestRegime] > 0 {
		confidence = bestScore / float64(regimeCount[bestRegime])
	}

	// If confidence is too low, return undetermined
	if confidence < 0.6 {
		bestRegime = RegimeUndetermined
		confidence = 0.5
	}

	return &RegimeDetectionResult{
		Regime:     bestRegime,
		Confidence: confidence,
		Metadata: map[string]interface{}{
			"regimeScores": regimeScores,
			"regimeCount":  regimeCount,
		},
	}
}
