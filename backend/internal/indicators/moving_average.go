package indicators

import (
	"errors"
	"math"
)

// SMA calculates the Simple Moving Average for the given period
func SMA(data []float64, period int) ([]float64, error) {
	if period <= 0 {
		return nil, errors.New("period must be positive")
	}
	if len(data) < period {
		return nil, errors.New("not enough data points for the specified period")
	}

	result := make([]float64, len(data)-period+1)
	for i := 0; i <= len(data)-period; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += data[i+j]
		}
		result[i] = sum / float64(period)
	}

	return result, nil
}

// EMA calculates the Exponential Moving Average for the given period
func EMA(data []float64, period int) ([]float64, error) {
	if period <= 0 {
		return nil, errors.New("period must be positive")
	}
	if len(data) < period {
		return nil, errors.New("not enough data points for the specified period")
	}

	// Calculate the multiplier
	multiplier := 2.0 / (float64(period) + 1.0)

	// Calculate the initial SMA
	sma, err := SMA(data[:period], period)
	if err != nil {
		return nil, err
	}

	// Initialize the result with the SMA as the first value
	result := make([]float64, len(data)-period+1)
	result[0] = sma[0]

	// Calculate EMA for the rest of the data
	for i := 1; i < len(result); i++ {
		result[i] = (data[i+period-1]-result[i-1])*multiplier + result[i-1]
	}

	return result, nil
}

// MACD calculates the Moving Average Convergence Divergence
func MACD(data []float64, fastPeriod, slowPeriod, signalPeriod int) ([]float64, []float64, []float64, error) {
	if fastPeriod <= 0 || slowPeriod <= 0 || signalPeriod <= 0 {
		return nil, nil, nil, errors.New("periods must be positive")
	}
	if fastPeriod >= slowPeriod {
		return nil, nil, nil, errors.New("fast period must be less than slow period")
	}
	if len(data) < slowPeriod {
		return nil, nil, nil, errors.New("not enough data points for the specified periods")
	}

	// Calculate fast EMA
	fastEMA, err := EMA(data, fastPeriod)
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate slow EMA
	slowEMA, err := EMA(data, slowPeriod)
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate MACD line (fast EMA - slow EMA)
	macdLine := make([]float64, len(slowEMA))
	for i := 0; i < len(slowEMA); i++ {
		fastIndex := i + len(fastEMA) - len(slowEMA)
		macdLine[i] = fastEMA[fastIndex] - slowEMA[i]
	}

	// Calculate signal line (EMA of MACD line)
	signalLine, err := EMA(macdLine, signalPeriod)
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate histogram (MACD line - signal line)
	histogram := make([]float64, len(signalLine))
	for i := 0; i < len(signalLine); i++ {
		macdIndex := i + len(macdLine) - len(signalLine)
		histogram[i] = macdLine[macdIndex] - signalLine[i]
	}

	return macdLine, signalLine, histogram, nil
}

// RSI calculates the Relative Strength Index
func RSI(data []float64, period int) ([]float64, error) {
	if period <= 0 {
		return nil, errors.New("period must be positive")
	}
	if len(data) <= period {
		return nil, errors.New("not enough data points for the specified period")
	}

	// Calculate price changes
	changes := make([]float64, len(data)-1)
	for i := 0; i < len(changes); i++ {
		changes[i] = data[i+1] - data[i]
	}

	// Calculate gains and losses
	gains := make([]float64, len(changes))
	losses := make([]float64, len(changes))
	for i := 0; i < len(changes); i++ {
		if changes[i] > 0 {
			gains[i] = changes[i]
		} else {
			losses[i] = -changes[i]
		}
	}

	// Calculate average gains and losses
	result := make([]float64, len(data)-period)
	var avgGain, avgLoss float64

	// First average gain and loss
	for i := 0; i < period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate first RSI
	rs := avgGain / avgLoss
	result[0] = 100 - (100 / (1 + rs))

	// Calculate remaining RSIs
	for i := 1; i < len(result); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i+period-1]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i+period-1]) / float64(period)
		rs = avgGain / avgLoss
		result[i] = 100 - (100 / (1 + rs))
	}

	return result, nil
}

// BollingerBands calculates the Bollinger Bands
func BollingerBands(data []float64, period int, deviations float64) ([]float64, []float64, []float64, error) {
	if period <= 0 {
		return nil, nil, nil, errors.New("period must be positive")
	}
	if len(data) < period {
		return nil, nil, nil, errors.New("not enough data points for the specified period")
	}

	// Calculate SMA (middle band)
	middleBand, err := SMA(data, period)
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate standard deviation
	upperBand := make([]float64, len(middleBand))
	lowerBand := make([]float64, len(middleBand))

	for i := 0; i < len(middleBand); i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += math.Pow(data[i+j]-middleBand[i], 2)
		}
		stdDev := math.Sqrt(sum / float64(period))
		upperBand[i] = middleBand[i] + deviations*stdDev
		lowerBand[i] = middleBand[i] - deviations*stdDev
	}

	return upperBand, middleBand, lowerBand, nil
}

// ATR calculates the Average True Range
func ATR(high, low, close []float64, period int) ([]float64, error) {
	if period <= 0 {
		return nil, errors.New("period must be positive")
	}
	if len(high) != len(low) || len(high) != len(close) {
		return nil, errors.New("high, low, and close arrays must have the same length")
	}
	if len(high) <= period {
		return nil, errors.New("not enough data points for the specified period")
	}

	// Calculate true range
	tr := make([]float64, len(high))
	tr[0] = high[0] - low[0] // First TR is simply the first day's range

	for i := 1; i < len(high); i++ {
		// True range is the greatest of:
		// 1. Current high - current low
		// 2. Absolute value of current high - previous close
		// 3. Absolute value of current low - previous close
		tr[i] = math.Max(high[i]-low[i], math.Max(
			math.Abs(high[i]-close[i-1]),
			math.Abs(low[i]-close[i-1]),
		))
	}

	// Calculate ATR
	result := make([]float64, len(high)-period+1)

	// First ATR is the simple average of the first 'period' true ranges
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += tr[i]
	}
	result[0] = sum / float64(period)

	// Subsequent ATRs use the smoothing formula: ATR = ((period-1) * previous ATR + current TR) / period
	for i := 1; i < len(result); i++ {
		result[i] = (result[i-1]*float64(period-1) + tr[i+period-1]) / float64(period)
	}

	return result, nil
}
