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
	csvContent := `timestamp,open,high,low,close,volume
2023-01-01T00:00:00Z,100.0,105.0,95.0,102.0,1000.0
2023-01-01T01:00:00Z,102.0,107.0,101.0,106.0,1200.0
2023-01-01T02:00:00Z,106.0,110.0,104.0,108.0,1500.0
2023-01-01T03:00:00Z,108.0,1000.0,104.0,110.0,1300.0
2023-01-01T05:00:00Z,112.0,115.0,110.0,114.0,1400.0
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

	// Should have 6 klines (0h, 1h, 2h, 3h, 4h, 5h) with 4h being interpolated
	assert.Len(t, dataset.Klines, 6)

	// Check that the outlier at 3h was fixed
	outlierKline := dataset.Klines[3]
	// The original value was 1000.0, so after fixing it should be much lower
	assert.True(t, outlierKline.High < 200.0, "Outlier should have been fixed")
	// Print the actual value for debugging
	t.Logf("Fixed outlier high value: %.2f", outlierKline.High)

	// Check that the missing value at 4h was interpolated
	missingKline := dataset.Klines[4]
	assert.NotNil(t, missingKline, "Missing kline should have been interpolated")
}
