package backtest

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// DataSet represents a collection of market data for backtesting
type DataSet struct {
	Symbol    string                                // Trading pair symbol
	Interval  string                                // Time interval (e.g., "1h", "1d")
	Klines    []*models.Kline                       // Candlestick data
	Tickers   []*models.Ticker                      // Ticker data (optional)
	OrderBook map[time.Time]*models.OrderBookUpdate // Order book snapshots (optional)
}

// DataLoaderOptions defines configuration options for the DataLoader
type DataLoaderOptions struct {
	FillMissingValues bool    // Whether to interpolate missing values
	DetectOutliers    bool    // Whether to detect and fix outliers
	OutlierThreshold  float64 // Number of standard deviations for outlier detection
	Resample          bool    // Whether to resample data to a different interval
	ResampleInterval  string  // Target interval for resampling
}

// DataLoader handles loading and preprocessing historical market data
type DataLoader struct {
	dataDir string
	options *DataLoaderOptions
}

// NewDataLoader creates a new DataLoader with default options
func NewDataLoader(dataDir string) *DataLoader {
	return &DataLoader{
		dataDir: dataDir,
		options: &DataLoaderOptions{
			FillMissingValues: false,
			DetectOutliers:    false,
			OutlierThreshold:  3.0,
			Resample:          false,
		},
	}
}

// NewDataLoaderWithOptions creates a new DataLoader with custom options
func NewDataLoaderWithOptions(dataDir string, options *DataLoaderOptions) *DataLoader {
	return &DataLoader{
		dataDir: dataDir,
		options: options,
	}
}

// LoadData loads historical market data for a symbol within a time range
func (l *DataLoader) LoadData(ctx context.Context, symbol, interval string, startTime, endTime time.Time) (*DataSet, error) {
	// Create a new dataset
	dataset := &DataSet{
		Symbol:    symbol,
		Interval:  interval,
		OrderBook: make(map[time.Time]*models.OrderBookUpdate),
	}

	// Load klines
	klines, err := l.loadKlinesFromCSV(symbol, interval, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to load klines: %w", err)
	}
	dataset.Klines = klines

	// Preprocess data if options are enabled
	if l.options.FillMissingValues {
		dataset.Klines = l.fillMissingValues(dataset.Klines, interval)
	}

	if l.options.DetectOutliers {
		err = l.detectAndFixOutliers(dataset)
		if err != nil {
			return nil, fmt.Errorf("failed to detect and fix outliers: %w", err)
		}
	}

	if l.options.Resample && l.options.ResampleInterval != "" {
		dataset.Klines = l.resampleData(dataset.Klines, interval, l.options.ResampleInterval)
		dataset.Interval = l.options.ResampleInterval
	}

	// Optionally load tickers and order book data
	// This is implemented separately as they're often not needed for basic backtesting

	return dataset, nil
}

// loadKlinesFromCSV loads candlestick data from a CSV file
func (l *DataLoader) loadKlinesFromCSV(symbol, interval string, startTime, endTime time.Time) ([]*models.Kline, error) {
	// Construct the file path
	filePath := filepath.Join(l.dataDir, fmt.Sprintf("%s_%s.csv", symbol, interval))

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Map header columns to indices
	columnMap := make(map[string]int)
	for i, column := range header {
		columnMap[column] = i
	}

	// Check required columns
	requiredColumns := []string{"timestamp", "open", "high", "low", "close", "volume"}
	for _, column := range requiredColumns {
		if _, ok := columnMap[column]; !ok {
			return nil, fmt.Errorf("CSV file missing required column: %s", column)
		}
	}

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	// Parse records into Kline objects
	klines := make([]*models.Kline, 0, len(records))
	for _, record := range records {
		// Parse timestamp
		timestampStr := record[columnMap["timestamp"]]
		var timestamp time.Time

		// Try different timestamp formats
		timestamp, err = time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			// Try Unix timestamp (milliseconds)
			if ms, err := strconv.ParseInt(timestampStr, 10, 64); err == nil {
				timestamp = time.Unix(ms/1000, (ms%1000)*1000000)
			} else {
				// Try other common formats
				formats := []string{
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05",
					"2006-01-02",
				}
				for _, format := range formats {
					if t, err := time.Parse(format, timestampStr); err == nil {
						timestamp = t
						break
					}
				}
			}
		}

		// Skip if timestamp is outside the requested range
		if timestamp.Before(startTime) || timestamp.After(endTime) {
			continue
		}

		// Parse numeric values
		open, _ := strconv.ParseFloat(record[columnMap["open"]], 64)
		high, _ := strconv.ParseFloat(record[columnMap["high"]], 64)
		low, _ := strconv.ParseFloat(record[columnMap["low"]], 64)
		close, _ := strconv.ParseFloat(record[columnMap["close"]], 64)
		volume, _ := strconv.ParseFloat(record[columnMap["volume"]], 64)

		// Create Kline object
		kline := &models.Kline{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  timestamp,
			CloseTime: timestamp.Add(parseInterval(interval)),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}

		klines = append(klines, kline)
	}

	// Sort klines by open time
	sort.Slice(klines, func(i, j int) bool {
		return klines[i].OpenTime.Before(klines[j].OpenTime)
	})

	return klines, nil
}

// fillMissingValues interpolates missing values in the kline data
func (l *DataLoader) fillMissingValues(klines []*models.Kline, interval string) []*models.Kline {
	if len(klines) == 0 {
		return klines
	}

	// Create a map of existing timestamps
	existingTimes := make(map[int64]bool)
	for _, kline := range klines {
		existingTimes[kline.OpenTime.Unix()] = true
	}

	// Sort klines by open time
	sort.Slice(klines, func(i, j int) bool {
		return klines[i].OpenTime.Before(klines[j].OpenTime)
	})

	// Calculate the expected interval in seconds
	intervalDuration := parseInterval(interval)
	intervalSeconds := int64(intervalDuration.Seconds())

	// Create a new slice for the result
	result := make([]*models.Kline, 0, len(klines))
	result = append(result, klines[0])

	// Iterate through the time range and fill gaps
	startTime := klines[0].OpenTime.Unix()
	endTime := klines[len(klines)-1].OpenTime.Unix()

	for t := startTime + intervalSeconds; t <= endTime; t += intervalSeconds {
		if existingTimes[t] {
			// Find the kline with this timestamp
			for _, kline := range klines {
				if kline.OpenTime.Unix() == t {
					result = append(result, kline)
					break
				}
			}
		} else {
			// Interpolate a new kline
			prevKline := result[len(result)-1]
			nextKlineIdx := -1

			// Find the next existing kline
			for i, kline := range klines {
				if kline.OpenTime.Unix() > t {
					nextKlineIdx = i
					break
				}
			}

			if nextKlineIdx != -1 {
				nextKline := klines[nextKlineIdx]

				// Calculate the weight for linear interpolation
				totalDiff := float64(nextKline.OpenTime.Unix() - prevKline.OpenTime.Unix())
				weight := float64(t-prevKline.OpenTime.Unix()) / totalDiff

				// Interpolate values
				interpolatedKline := &models.Kline{
					Symbol:    prevKline.Symbol,
					Interval:  interval,
					OpenTime:  time.Unix(t, 0),
					CloseTime: time.Unix(t, 0).Add(intervalDuration),
					Open:      prevKline.Close, // Use previous close as open
					High:      linearInterpolate(prevKline.High, nextKline.High, weight),
					Low:       linearInterpolate(prevKline.Low, nextKline.Low, weight),
					Close:     linearInterpolate(prevKline.Close, nextKline.Close, weight),
					Volume:    linearInterpolate(prevKline.Volume, nextKline.Volume, weight),
				}

				result = append(result, interpolatedKline)
			} else {
				// If there's no next kline, use the previous values
				interpolatedKline := &models.Kline{
					Symbol:    prevKline.Symbol,
					Interval:  interval,
					OpenTime:  time.Unix(t, 0),
					CloseTime: time.Unix(t, 0).Add(intervalDuration),
					Open:      prevKline.Close, // Use previous close as open
					High:      prevKline.High,
					Low:       prevKline.Low,
					Close:     prevKline.Close,
					Volume:    prevKline.Volume,
				}

				result = append(result, interpolatedKline)
			}
		}
	}

	return result
}

// detectAndFixOutliers identifies and corrects outliers in the kline data
func (l *DataLoader) detectAndFixOutliers(dataset *DataSet) error {
	if len(dataset.Klines) == 0 {
		return nil
	}

	// Calculate statistics for each price field
	openPrices := make([]float64, len(dataset.Klines))
	highPrices := make([]float64, len(dataset.Klines))
	lowPrices := make([]float64, len(dataset.Klines))
	closePrices := make([]float64, len(dataset.Klines))

	for i, kline := range dataset.Klines {
		openPrices[i] = kline.Open
		highPrices[i] = kline.High
		lowPrices[i] = kline.Low
		closePrices[i] = kline.Close
	}

	// Calculate statistics for each price type
	openMean, openStdDev := calculateStats(openPrices)
	highMean, highStdDev := calculateStats(highPrices)
	lowMean, lowStdDev := calculateStats(lowPrices)
	closeMean, closeStdDev := calculateStats(closePrices)

	// Fix outliers in each kline
	for _, kline := range dataset.Klines {
		kline.Open = fixOutlier(kline.Open, openMean, openStdDev, 3.0)
		kline.High = fixOutlier(kline.High, highMean, highStdDev, 3.0)
		kline.Low = fixOutlier(kline.Low, lowMean, lowStdDev, 3.0)
		kline.Close = fixOutlier(kline.Close, closeMean, closeStdDev, 3.0)

		// Ensure High is the highest price
		kline.High = math.Max(math.Max(math.Max(kline.Open, kline.High), kline.Low), kline.Close)
		// Ensure Low is the lowest price
		kline.Low = math.Min(math.Min(math.Min(kline.Open, kline.Low), kline.High), kline.Close)
	}

	return nil
}

// resampleData converts kline data to a different time interval
func (l *DataLoader) resampleData(klines []*models.Kline, sourceInterval, targetInterval string) []*models.Kline {
	if len(klines) == 0 {
		return klines
	}

	// Parse intervals to durations
	sourceDuration := parseInterval(sourceInterval)
	targetDuration := parseInterval(targetInterval)

	// If target interval is smaller than source, we can't resample
	if targetDuration < sourceDuration {
		return klines // Return original data
	}

	// Calculate how many source intervals fit in one target interval
	ratio := int(targetDuration / sourceDuration)

	// Create a new slice for the result
	result := make([]*models.Kline, 0, len(klines)/ratio+1)

	// Group klines by target interval
	for i := 0; i < len(klines); i += ratio {
		end := i + ratio
		if end > len(klines) {
			end = len(klines)
		}

		group := klines[i:end]
		if len(group) == 0 {
			continue
		}

		// Create a new resampled kline
		resampled := &models.Kline{
			Symbol:    group[0].Symbol,
			Interval:  targetInterval,
			OpenTime:  group[0].OpenTime,
			CloseTime: group[len(group)-1].CloseTime,
			Open:      group[0].Open,
			High:      group[0].High,
			Low:       group[0].Low,
			Close:     group[len(group)-1].Close,
			Volume:    0,
		}

		// Find highest high, lowest low, and sum volumes
		for _, kline := range group {
			resampled.High = math.Max(resampled.High, kline.High)
			resampled.Low = math.Min(resampled.Low, kline.Low)
			resampled.Volume += kline.Volume
		}

		result = append(result, resampled)
	}

	return result
}

// calculateStats returns the mean and standard deviation of a slice of values
func calculateStats(values []float64) (mean, stdDev float64) {
	if len(values) == 0 {
		return 0, 0
	}

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	stdDev = math.Sqrt(sumSquaredDiff / float64(len(values)))

	return mean, stdDev
}

// fixOutlier adjusts an outlier value to be within the acceptable range
func fixOutlier(value, mean, stdDev, threshold float64) float64 {
	if stdDev == 0 {
		return value // Can't fix outliers with zero standard deviation
	}

	zScore := (value - mean) / stdDev
	if math.Abs(zScore) > threshold {
		// Cap the value at mean ± (threshold * stdDev)
		if zScore > 0 {
			return mean + (threshold * stdDev)
		}
		return mean - (threshold * stdDev)
	}
	return value
}

// linearInterpolate performs linear interpolation between two values
func linearInterpolate(start, end, weight float64) float64 {
	return start + (end-start)*weight
}

// Note: Using parseInterval function from csv_data_provider.go
