package backtest

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataLoader(t *testing.T) {
	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "dataloader_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test CSV data
	csvContent := `timestamp,open,high,low,close,volume
2023-01-01T00:00:00Z,100.0,105.0,95.0,102.0,1000.0
2023-01-01T01:00:00Z,102.0,107.0,101.0,106.0,1200.0
2023-01-01T02:00:00Z,106.0,110.0,104.0,108.0,1500.0
`
	// Write test CSV file
	symbol := "BTCUSDT"
	interval := "1h"
	csvPath := filepath.Join(tempDir, symbol+"_"+interval+".csv")
	err = os.WriteFile(csvPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Create DataLoader
	loader := NewDataLoader(tempDir)

	// Test loading data
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 1, 3, 0, 0, 0, time.UTC)

	dataset, err := loader.LoadData(context.Background(), symbol, interval, startTime, endTime)
	require.NoError(t, err)

	// Verify dataset
	assert.Equal(t, symbol, dataset.Symbol)
	assert.Equal(t, interval, dataset.Interval)
	assert.Len(t, dataset.Klines, 3)

	// Verify first kline
	firstKline := dataset.Klines[0]
	assert.Equal(t, startTime, firstKline.OpenTime)
	assert.Equal(t, 100.0, firstKline.Open)
	assert.Equal(t, 105.0, firstKline.High)
	assert.Equal(t, 95.0, firstKline.Low)
	assert.Equal(t, 102.0, firstKline.Close)
	assert.Equal(t, 1000.0, firstKline.Volume)
}

func TestDataPreprocessing(t *testing.T) {
	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "dataloader_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test CSV data with missing values and outliers
	// Using more realistic price movements with an obvious outlier
	csvContent := `timestamp,open,high,low,close,volume
2023-01-01T00:00:00Z,100.0,102.0,98.0,101.0,1000.0
2023-01-01T01:00:00Z,101.0,103.0,99.0,102.0,1200.0
2023-01-01T02:00:00Z,102.0,104.0,100.0,103.0,1100.0
2023-01-01T03:00:00Z,103.0,150.0,101.0,104.0,1300.0
2023-01-01T05:00:00Z,105.0,107.0,103.0,106.0,1200.0
`
	// Write test CSV file
	symbol := "BTCUSDT"
	interval := "1h"
	csvPath := filepath.Join(tempDir, symbol+"_"+interval+".csv")
	err = os.WriteFile(csvPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Create DataLoader with preprocessing options
	options := &DataLoaderOptions{
		FillMissingValues: true,
		DetectOutliers:    true,
		OutlierThreshold:  3.0, // 3 standard deviations
	}
	loader := NewDataLoaderWithOptions(tempDir, options)

	// Test loading and preprocessing data
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 1, 6, 0, 0, 0, time.UTC)

	dataset, err := loader.LoadData(context.Background(), symbol, interval, startTime, endTime)
	require.NoError(t, err)

	// Verify dataset
	assert.Equal(t, symbol, dataset.Symbol)
	assert.Equal(t, interval, dataset.Interval)

	// Should have 6 klines (0h through 5h) with 4h being interpolated
	assert.Len(t, dataset.Klines, 6)

	// Check that the outlier at 3h was fixed
	outlierKline := dataset.Klines[3]
	assert.Less(t, outlierKline.High, 150.0, "Outlier high value should have been reduced")
	assert.Greater(t, outlierKline.High, 100.0, "Outlier high value should still be reasonable")
	assert.Equal(t, 103.0, outlierKline.Open, "Non-outlier values should remain unchanged")
	assert.Equal(t, 101.0, outlierKline.Low, "Non-outlier values should remain unchanged")
	assert.Equal(t, 104.0, outlierKline.Close, "Non-outlier values should remain unchanged")

	// Check that the missing value at 4h was interpolated
	missingKline := dataset.Klines[4]
	require.NotNil(t, missingKline, "Missing kline should have been interpolated")

	// Verify interpolated values are reasonable
	assert.InDelta(t, 104.5, missingKline.Open, 0.1, "Interpolated open price should be halfway between surrounding values")
	assert.InDelta(t, 105.5, missingKline.Close, 0.1, "Interpolated close price should be halfway between surrounding values")
	assert.True(t, missingKline.High > missingKline.Open && missingKline.High > missingKline.Close, "Interpolated high should be above open and close")
	assert.True(t, missingKline.Low < missingKline.Open && missingKline.Low < missingKline.Close, "Interpolated low should be below open and close")
	assert.InDelta(t, 1250.0, missingKline.Volume, 100.0, "Interpolated volume should be reasonable")
}
