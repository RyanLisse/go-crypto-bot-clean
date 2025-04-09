package report

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/repository/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestReportRepository(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create a test database
	dbConfig := database.Config{
		DatabasePath:    ":memory:", // Use in-memory SQLite for testing
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Create SQLite repository
	db, err := database.NewSQLiteRepository(dbConfig)
	require.NoError(t, err)
	defer db.Close()

	// Create report repository
	repo := NewReportRepository(db, logger)

	// Initialize repository
	ctx := context.Background()
	err = repo.Initialize(ctx)
	require.NoError(t, err)

	// Test saving and retrieving a report
	t.Run("SaveAndGetReport", func(t *testing.T) {
		// Create a test report
		report := &models.PerformanceReport{
			ID:          uuid.New().String(),
			GeneratedAt: time.Now().UTC(),
			Period:      string(models.ReportPeriodDaily),
			Analysis:    "Test analysis",
			Metrics: models.PerformanceReportMetrics{
				Timestamp: time.Now().UTC(),
				Metrics: map[string]interface{}{
					"test_metric": 42.0,
				},
				SystemState: models.SystemState{
					CPUUsage:    10.5,
					MemoryUsage: 256.0,
					Latency:     15,
					Goroutines:  10,
					Uptime:      "1h",
				},
			},
			Insights: []string{"Test insight 1", "Test insight 2"},
		}

		// Save the report
		err := repo.SaveReport(ctx, report)
		require.NoError(t, err)

		// Get the report by ID
		retrieved, err := repo.GetReportByID(ctx, report.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// Verify the report
		assert.Equal(t, report.ID, retrieved.ID)
		assert.Equal(t, report.Period, retrieved.Period)
		assert.Equal(t, report.Analysis, retrieved.Analysis)
		assert.Equal(t, report.Insights, retrieved.Insights)
		assert.Equal(t, report.Metrics.SystemState.CPUUsage, retrieved.Metrics.SystemState.CPUUsage)
		assert.Equal(t, report.Metrics.SystemState.MemoryUsage, retrieved.Metrics.SystemState.MemoryUsage)
		assert.Equal(t, report.Metrics.SystemState.Latency, retrieved.Metrics.SystemState.Latency)
		assert.Equal(t, report.Metrics.SystemState.Goroutines, retrieved.Metrics.SystemState.Goroutines)
		assert.Equal(t, report.Metrics.SystemState.Uptime, retrieved.Metrics.SystemState.Uptime)
		assert.Equal(t, 42.0, retrieved.Metrics.Metrics["test_metric"])
	})

	// Test getting reports by period
	t.Run("GetReportsByPeriod", func(t *testing.T) {
		// Create multiple test reports
		for i := 0; i < 3; i++ {
			report := &models.PerformanceReport{
				ID:          uuid.New().String(),
				GeneratedAt: time.Now().UTC().Add(time.Duration(-i) * time.Hour),
				Period:      string(models.ReportPeriodHourly),
				Analysis:    fmt.Sprintf("Test analysis %d", i),
				Metrics: models.PerformanceReportMetrics{
					Timestamp: time.Now().UTC(),
					Metrics: map[string]interface{}{
						"test_metric": float64(i),
					},
					SystemState: models.SystemState{
						CPUUsage:    10.5,
						MemoryUsage: 256.0,
						Latency:     15,
						Goroutines:  10,
						Uptime:      "1h",
					},
				},
				Insights: []string{fmt.Sprintf("Test insight %d", i)},
			}

			err := repo.SaveReport(ctx, report)
			require.NoError(t, err)
		}

		// Get reports by period
		reports, err := repo.GetReportsByPeriod(ctx, models.ReportPeriodHourly, 10)
		require.NoError(t, err)
		assert.Len(t, reports, 3)

		// Verify reports are ordered by generated_at DESC
		for i := 0; i < len(reports)-1; i++ {
			assert.True(t, reports[i].GeneratedAt.After(reports[i+1].GeneratedAt))
		}
	})

	// Test getting the latest report
	t.Run("GetLatestReport", func(t *testing.T) {
		// Get the latest report
		latest, err := repo.GetLatestReport(ctx)
		require.NoError(t, err)
		require.NotNil(t, latest)

		// Verify it's the most recent one
		reports, err := repo.GetReportsByPeriod(ctx, models.ReportPeriod(latest.Period), 10)
		require.NoError(t, err)
		assert.Equal(t, latest.ID, reports[0].ID)
	})
}
