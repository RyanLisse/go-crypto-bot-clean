package backtest

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseDataLoader implements data loading from a database
type DatabaseDataLoader struct {
	db      *gorm.DB
	options *DataLoaderOptions
}

// NewDatabaseDataLoader creates a new DatabaseDataLoader with default options
func NewDatabaseDataLoader(dbPath string) (*DatabaseDataLoader, error) {
	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DatabaseDataLoader{
		db: db,
		options: &DataLoaderOptions{
			FillMissingValues: false,
			DetectOutliers:    false,
			OutlierThreshold:  3.0,
			Resample:          false,
		},
	}, nil
}

// NewDatabaseDataLoaderWithOptions creates a new DatabaseDataLoader with custom options
func NewDatabaseDataLoaderWithOptions(dbPath string, options *DataLoaderOptions) (*DatabaseDataLoader, error) {
	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DatabaseDataLoader{
		db:      db,
		options: options,
	}, nil
}

// LoadData loads historical market data from the database
func (l *DatabaseDataLoader) LoadData(ctx context.Context, symbol, interval string, startTime, endTime time.Time) (*DataSet, error) {
	// Create a new dataset
	dataset := &DataSet{
		Symbol:    symbol,
		Interval:  interval,
		OrderBook: make(map[time.Time]*models.OrderBookUpdate),
	}

	// Load klines from the database
	var klines []*models.Kline
	result := l.db.WithContext(ctx).
		Where("symbol = ? AND interval = ? AND open_time >= ? AND open_time <= ?", 
			symbol, interval, startTime, endTime).
		Order("open_time ASC").
		Find(&klines)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to load klines from database: %w", result.Error)
	}

	dataset.Klines = klines

	// Preprocess data if options are enabled
	if l.options.FillMissingValues {
		dataset.Klines = l.fillMissingValues(dataset.Klines, interval)
	}

	if l.options.DetectOutliers {
		dataset.Klines = l.detectAndFixOutliers(dataset.Klines)
	}

	if l.options.Resample && l.options.ResampleInterval != "" {
		dataset.Klines = l.resampleData(dataset.Klines, interval, l.options.ResampleInterval)
		dataset.Interval = l.options.ResampleInterval
	}

	// Optionally load tickers
	if ctx.Value("loadTickers") != nil {
		var tickers []*models.Ticker
		result := l.db.WithContext(ctx).
			Where("symbol = ? AND timestamp >= ? AND timestamp <= ?", 
				symbol, startTime, endTime).
			Order("timestamp ASC").
			Find(&tickers)

		if result.Error != nil {
			return nil, fmt.Errorf("failed to load tickers from database: %w", result.Error)
		}

		dataset.Tickers = tickers
	}

	// Optionally load order book data
	if ctx.Value("loadOrderBook") != nil {
		var orderBooks []*models.OrderBook
		result := l.db.WithContext(ctx).
			Where("symbol = ? AND timestamp >= ? AND timestamp <= ?", 
				symbol, startTime, endTime).
			Order("timestamp ASC").
			Preload("Bids").
			Preload("Asks").
			Find(&orderBooks)

		if result.Error != nil {
			return nil, fmt.Errorf("failed to load order books from database: %w", result.Error)
		}

		// Convert to OrderBookUpdate format
		for _, ob := range orderBooks {
			bids := make([]models.OrderBookEntry, len(ob.Bids))
			asks := make([]models.OrderBookEntry, len(ob.Asks))

			for i, bid := range ob.Bids {
				bids[i] = models.OrderBookEntry{
					Price:    bid.Price,
					Quantity: bid.Quantity,
				}
			}

			for i, ask := range ob.Asks {
				asks[i] = models.OrderBookEntry{
					Price:    ask.Price,
					Quantity: ask.Quantity,
				}
			}

			dataset.OrderBook[ob.Timestamp] = &models.OrderBookUpdate{
				Symbol:    ob.Symbol,
				Timestamp: ob.Timestamp,
				Bids:      bids,
				Asks:      asks,
			}
		}
	}

	return dataset, nil
}

// fillMissingValues interpolates missing values in the kline data
func (l *DatabaseDataLoader) fillMissingValues(klines []*models.Kline, interval string) []*models.Kline {
	// Reuse the implementation from DataLoader
	dataLoader := &DataLoader{options: l.options}
	return dataLoader.fillMissingValues(klines, interval)
}

// detectAndFixOutliers identifies and corrects outliers in the kline data
func (l *DatabaseDataLoader) detectAndFixOutliers(klines []*models.Kline) []*models.Kline {
	// Reuse the implementation from DataLoader
	dataLoader := &DataLoader{options: l.options}
	return dataLoader.detectAndFixOutliers(klines)
}

// resampleData converts kline data to a different time interval
func (l *DatabaseDataLoader) resampleData(klines []*models.Kline, sourceInterval, targetInterval string) []*models.Kline {
	// Reuse the implementation from DataLoader
	dataLoader := &DataLoader{options: l.options}
	return dataLoader.resampleData(klines, sourceInterval, targetInterval)
}
