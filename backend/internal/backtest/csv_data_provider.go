package backtest

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// CSVDataProvider implements the DataProvider interface using CSV files
type CSVDataProvider struct {
	dataDir string
}

// NewCSVDataProvider creates a new CSVDataProvider
func NewCSVDataProvider(dataDir string) *CSVDataProvider {
	return &CSVDataProvider{
		dataDir: dataDir,
	}
}

// GetHistoricalData retrieves historical candlestick data from CSV files
// (Renamed from GetKlines to match DataProvider interface)
func (p *CSVDataProvider) GetHistoricalData(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error) {
	// Construct the file path
	filePath := filepath.Join(p.dataDir, fmt.Sprintf("%s_%s.csv", symbol, interval))

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
		columnMap[strings.TrimSpace(column)] = i // Trim spaces from header
	}

	// Check required columns
	requiredColumns := []string{"timestamp", "open", "high", "low", "close", "volume"}
	for _, column := range requiredColumns {
		if _, ok := columnMap[column]; !ok {
			return nil, fmt.Errorf("CSV file \"%s\" missing required column: %s", filePath, column)
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
				timestamp = time.UnixMilli(ms)
			} else {
				// Try other common formats
				formats := []string{
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05Z07:00", // Added timezone support
					"2006-01-02T15:04:05",
					"2006/01/02 15:04:05",
					"2006-01-02",
				}
				parsed := false
				for _, format := range formats {
					if t, err := time.Parse(format, timestampStr); err == nil {
						timestamp = t
						parsed = true
						break
					}
				}
				if !parsed {
					// Log or handle error if no format matches
					fmt.Printf("Warning: Could not parse timestamp '%s' in file %s\n", timestampStr, filePath)
					continue // Skip record if timestamp is unparseable
				}
			}
		}

		// Skip if timestamp is outside the requested range [startTime, endTime]
		if timestamp.Before(startTime) || timestamp.After(endTime) {
			continue
		}

		// Parse numeric values
		open, _ := strconv.ParseFloat(record[columnMap["open"]], 64)
		high, _ := strconv.ParseFloat(record[columnMap["high"]], 64)
		low, _ := strconv.ParseFloat(record[columnMap["low"]], 64)
		closeVal, _ := strconv.ParseFloat(record[columnMap["close"]], 64) // Renamed to avoid conflict with defer file.Close()
		volume, _ := strconv.ParseFloat(record[columnMap["volume"]], 64)

		// Create Kline object
		kline := &models.Kline{
			Symbol:    symbol,
			Interval:  interval,
			OpenTime:  timestamp,
			CloseTime: timestamp.Add(parseInterval(interval)), // Calculate CloseTime based on interval
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closeVal,
			Volume:    volume,
		}

		klines = append(klines, kline)
	}

	return klines, nil
}

// // GetTickers retrieves historical ticker data from CSV files
// // Commenting out as it's not part of the current DataProvider interface
// func (p *CSVDataProvider) GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error) {
// 	// ... (Implementation similar to GetKlines, parsing ticker-specific columns)
// 	return nil, fmt.Errorf("GetTickers not implemented for CSVDataProvider")
// }

// // GetOrderBook retrieves historical order book snapshots from CSV files
// // Commenting out as it's not part of the current DataProvider interface
// func (p *CSVDataProvider) GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error) {
// 	// ... (Implementation would require specific CSV format for order books)
// 	return nil, fmt.Errorf("GetOrderBook not implemented for CSVDataProvider")
// }

// Helper function to parse interval string to duration
func parseInterval(interval string) time.Duration {
	switch interval {
	case "1m":
		return time.Minute
	case "3m":
		return 3 * time.Minute
	case "5m":
		return 5 * time.Minute
	case "15m":
		return 15 * time.Minute
	case "30m":
		return 30 * time.Minute
	case "1h":
		return time.Hour
	case "2h":
		return 2 * time.Hour
	case "4h":
		return 4 * time.Hour
	case "6h":
		return 6 * time.Hour
	case "8h":
		return 8 * time.Hour
	case "12h":
		return 12 * time.Hour
	case "1d":
		return 24 * time.Hour
	case "3d":
		return 3 * 24 * time.Hour
	case "1w":
		return 7 * 24 * time.Hour
	case "1M":
		return 30 * 24 * time.Hour
	default:
		return time.Hour // Default to 1 hour
	}
}
