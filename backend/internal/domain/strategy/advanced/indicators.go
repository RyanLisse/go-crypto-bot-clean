package advanced

import (
	"errors"
	"math"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// IndicatorType represents the type of technical indicator
type IndicatorType string

const (
	IndicatorSMA  IndicatorType = "SMA"  // Simple Moving Average
	IndicatorEMA  IndicatorType = "EMA"  // Exponential Moving Average
	IndicatorRSI  IndicatorType = "RSI"  // Relative Strength Index
	IndicatorMACD IndicatorType = "MACD" // Moving Average Convergence Divergence
	IndicatorBB   IndicatorType = "BB"   // Bollinger Bands
)

// IndicatorResult represents the result of a technical indicator calculation
type IndicatorResult struct {
	Type   IndicatorType         // Type of indicator
	Values map[string]float64    // Main indicator values
	Series map[string][]float64  // Historical series if applicable
	Metadata map[string]interface{} // Additional metadata
}

// CalculateSMA calculates the Simple Moving Average
func CalculateSMA(candles []*models.Candle, period int) (*IndicatorResult, error) {
	if len(candles) < period {
		return nil, errors.New("not enough data for SMA calculation")
	}

	// Calculate SMA
	sum := 0.0
	for i := len(candles) - period; i < len(candles); i++ {
		sum += candles[i].ClosePrice
	}
	sma := sum / float64(period)

	// Create result
	result := &IndicatorResult{
		Type: IndicatorSMA,
		Values: map[string]float64{
			"sma": sma,
		},
		Series: map[string][]float64{},
		Metadata: map[string]interface{}{
			"period": period,
		},
	}

	// Calculate historical series if we have enough data
	if len(candles) > period {
		smaSeries := make([]float64, len(candles)-period+1)
		for i := 0; i <= len(candles)-period; i++ {
			sum := 0.0
			for j := i; j < i+period; j++ {
				sum += candles[j].ClosePrice
			}
			smaSeries[i] = sum / float64(period)
		}
		result.Series["sma"] = smaSeries
	}

	return result, nil
}

// CalculateEMA calculates the Exponential Moving Average
func CalculateEMA(candles []*models.Candle, period int) (*IndicatorResult, error) {
	if len(candles) < period {
		return nil, errors.New("not enough data for EMA calculation")
	}

	// Calculate multiplier
	multiplier := 2.0 / (float64(period) + 1.0)

	// Calculate initial SMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += candles[i].ClosePrice
	}
	ema := sum / float64(period)

	// Calculate EMA for remaining periods
	emaSeries := make([]float64, len(candles)-period+1)
	emaSeries[0] = ema

	for i := period; i < len(candles); i++ {
		ema = (candles[i].ClosePrice-ema)*multiplier + ema
		emaSeries[i-period+1] = ema
	}

	// Create result
	result := &IndicatorResult{
		Type: IndicatorEMA,
		Values: map[string]float64{
			"ema": ema,
		},
		Series: map[string][]float64{
			"ema": emaSeries,
		},
		Metadata: map[string]interface{}{
			"period": period,
		},
	}

	return result, nil
}

// CalculateRSI calculates the Relative Strength Index
func CalculateRSI(candles []*models.Candle, period int) (*IndicatorResult, error) {
	if len(candles) < period+1 {
		return nil, errors.New("not enough data for RSI calculation")
	}

	// Calculate price changes
	changes := make([]float64, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		changes[i-1] = candles[i].ClosePrice - candles[i-1].ClosePrice
	}

	// Calculate gains and losses
	gains := make([]float64, len(changes))
	losses := make([]float64, len(changes))
	for i, change := range changes {
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = math.Abs(change)
		}
	}

	// Calculate average gains and losses
	avgGain := 0.0
	avgLoss := 0.0
	for i := 0; i < period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI
	rs := 0.0
	if avgLoss > 0 {
		rs = avgGain / avgLoss
	}
	rsi := 100.0 - (100.0 / (1.0 + rs))

	// Calculate RSI series
	rsiSeries := make([]float64, len(changes)-period+1)
	rsiSeries[0] = rsi

	for i := period; i < len(changes); i++ {
		// Smooth average gain and loss
		avgGain = ((avgGain * float64(period-1)) + gains[i]) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + losses[i]) / float64(period)

		// Calculate RS and RSI
		if avgLoss > 0 {
			rs = avgGain / avgLoss
		} else {
			rs = 100.0 // Avoid division by zero
		}
		rsi = 100.0 - (100.0 / (1.0 + rs))
		rsiSeries[i-period+1] = rsi
	}

	// Create result
	result := &IndicatorResult{
		Type: IndicatorRSI,
		Values: map[string]float64{
			"rsi": rsi,
		},
		Series: map[string][]float64{
			"rsi": rsiSeries,
		},
		Metadata: map[string]interface{}{
			"period": period,
		},
	}

	return result, nil
}

// CalculateMACD calculates the Moving Average Convergence Divergence
func CalculateMACD(candles []*models.Candle, fastPeriod, slowPeriod, signalPeriod int) (*IndicatorResult, error) {
	if len(candles) < slowPeriod+signalPeriod {
		return nil, errors.New("not enough data for MACD calculation")
	}

	// Calculate fast EMA
	fastEMA, err := CalculateEMA(candles, fastPeriod)
	if err != nil {
		return nil, err
	}

	// Calculate slow EMA
	slowEMA, err := CalculateEMA(candles, slowPeriod)
	if err != nil {
		return nil, err
	}

	// Calculate MACD line
	macdLine := make([]float64, len(fastEMA.Series["ema"]))
	for i := 0; i < len(macdLine); i++ {
		if i < len(slowEMA.Series["ema"]) {
			macdLine[i] = fastEMA.Series["ema"][i] - slowEMA.Series["ema"][i]
		}
	}

	// Calculate signal line (EMA of MACD line)
	signalLine := make([]float64, len(macdLine)-signalPeriod+1)
	
	// Calculate initial SMA for signal line
	sum := 0.0
	for i := 0; i < signalPeriod; i++ {
		sum += macdLine[i]
	}
	signal := sum / float64(signalPeriod)
	signalLine[0] = signal

	// Calculate EMA for signal line
	multiplier := 2.0 / (float64(signalPeriod) + 1.0)
	for i := signalPeriod; i < len(macdLine); i++ {
		signal = (macdLine[i]-signal)*multiplier + signal
		signalLine[i-signalPeriod+1] = signal
	}

	// Calculate histogram
	histogram := make([]float64, len(signalLine))
	for i := 0; i < len(histogram); i++ {
		if i < len(macdLine) {
			histogram[i] = macdLine[len(macdLine)-len(signalLine)+i] - signalLine[i]
		}
	}

	// Create result
	result := &IndicatorResult{
		Type: IndicatorMACD,
		Values: map[string]float64{
			"macd":      macdLine[len(macdLine)-1],
			"signal":    signalLine[len(signalLine)-1],
			"histogram": histogram[len(histogram)-1],
		},
		Series: map[string][]float64{
			"macd":      macdLine,
			"signal":    signalLine,
			"histogram": histogram,
		},
		Metadata: map[string]interface{}{
			"fastPeriod":   fastPeriod,
			"slowPeriod":   slowPeriod,
			"signalPeriod": signalPeriod,
		},
	}

	return result, nil
}

// CalculateBollingerBands calculates Bollinger Bands
func CalculateBollingerBands(candles []*models.Candle, period int, deviations float64) (*IndicatorResult, error) {
	if len(candles) < period {
		return nil, errors.New("not enough data for Bollinger Bands calculation")
	}

	// Calculate SMA
	sma, err := CalculateSMA(candles, period)
	if err != nil {
		return nil, err
	}

	// Calculate standard deviation
	stdDev := 0.0
	for i := len(candles) - period; i < len(candles); i++ {
		stdDev += math.Pow(candles[i].ClosePrice-sma.Values["sma"], 2)
	}
	stdDev = math.Sqrt(stdDev / float64(period))

	// Calculate upper and lower bands
	upperBand := sma.Values["sma"] + (deviations * stdDev)
	lowerBand := sma.Values["sma"] - (deviations * stdDev)

	// Calculate historical bands
	upperBands := make([]float64, len(sma.Series["sma"]))
	lowerBands := make([]float64, len(sma.Series["sma"]))
	stdDevs := make([]float64, len(sma.Series["sma"]))

	for i := 0; i < len(sma.Series["sma"]); i++ {
		// Calculate standard deviation for this period
		stdDev := 0.0
		startIdx := len(candles) - len(sma.Series["sma"]) + i - period + 1
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx := startIdx + period
		if endIdx > len(candles) {
			endIdx = len(candles)
		}

		for j := startIdx; j < endIdx; j++ {
			stdDev += math.Pow(candles[j].ClosePrice-sma.Series["sma"][i], 2)
		}
		stdDev = math.Sqrt(stdDev / float64(endIdx-startIdx))
		stdDevs[i] = stdDev

		// Calculate bands
		upperBands[i] = sma.Series["sma"][i] + (deviations * stdDev)
		lowerBands[i] = sma.Series["sma"][i] - (deviations * stdDev)
	}

	// Create result
	result := &IndicatorResult{
		Type: IndicatorBB,
		Values: map[string]float64{
			"middle": sma.Values["sma"],
			"upper":  upperBand,
			"lower":  lowerBand,
			"width":  (upperBand - lowerBand) / sma.Values["sma"], // Normalized width
		},
		Series: map[string][]float64{
			"middle": sma.Series["sma"],
			"upper":  upperBands,
			"lower":  lowerBands,
		},
		Metadata: map[string]interface{}{
			"period":     period,
			"deviations": deviations,
		},
	}

	return result, nil
}
