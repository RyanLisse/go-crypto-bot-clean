package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewReportRepository(t *testing.T) {
	// Create dependencies
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Assert
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
	assert.Equal(t, logger, repo.logger)
}

func TestReportRepository_Initialize(t *testing.T) {
	// Create dependencies
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Set up expectations
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS performance_reports`).WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the method
	err = repo.Initialize(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_SaveReport(t *testing.T) {
	// Create dependencies
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Create test data
	report := &models.PerformanceReport{
		ID:          "test-id",
		GeneratedAt: time.Now(),
		Period:      "hourly",
		Analysis:    "Test analysis",
		Metrics: models.PerformanceReportMetrics{
			Timestamp: time.Now(),
			Metrics: map[string]interface{}{
				"test_metric": 123.45,
			},
			SystemState: models.SystemState{
				CPUUsage:    10.5,
				MemoryUsage: 256.0,
				Latency:     50,
				Goroutines:  10,
				Uptime:      "1h 30m",
			},
		},
		Insights: []string{"Insight 1", "Insight 2"},
	}

	// Set up expectations
	mock.ExpectExec(`INSERT INTO performance_reports`).
		WithArgs(report.ID, report.GeneratedAt, report.Period, report.Analysis, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = repo.SaveReport(context.Background(), report)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetReportByID(t *testing.T) {
	// Create dependencies
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Create test data
	id := "test-id"
	now := time.Now()
	metricsJSON := `{"timestamp":"2023-04-06T12:34:56Z","metrics":{"test_metric":123.45},"system_state":{"cpu_usage":10.5,"memory_usage":256,"latency_ms":50,"goroutines":10,"uptime":"1h 30m"}}`
	insightsJSON := `["Insight 1","Insight 2"]`

	// Set up expectations
	rows := sqlmock.NewRows([]string{"id", "generated_at", "period", "analysis", "metrics_json", "insights_json"}).
		AddRow(id, now, "hourly", "Test analysis", metricsJSON, insightsJSON)

	mock.ExpectQuery(`SELECT id, generated_at, period, analysis, metrics_json, insights_json FROM performance_reports WHERE id = \?`).
		WithArgs(id).
		WillReturnRows(rows)

	// Call the method
	report, err := repo.GetReportByID(context.Background(), id)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, id, report.ID)
	assert.Equal(t, "hourly", report.Period)
	assert.Equal(t, "Test analysis", report.Analysis)
	assert.Equal(t, []string{"Insight 1", "Insight 2"}, report.Insights)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetReportByID_NotFound(t *testing.T) {
	// Create dependencies
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Create test data
	id := "test-id"

	// Set up expectations
	mock.ExpectQuery(`SELECT id, generated_at, period, analysis, metrics_json, insights_json FROM performance_reports WHERE id = \?`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	// Call the method
	report, err := repo.GetReportByID(context.Background(), id)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetReportsByPeriod(t *testing.T) {
	// Create dependencies
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Create test data
	period := models.ReportPeriodHourly
	limit := 10
	now := time.Now()
	metricsJSON := `{"timestamp":"2023-04-06T12:34:56Z","metrics":{"test_metric":123.45},"system_state":{"cpu_usage":10.5,"memory_usage":256,"latency_ms":50,"goroutines":10,"uptime":"1h 30m"}}`
	insightsJSON := `["Insight 1","Insight 2"]`

	// Set up expectations
	rows := sqlmock.NewRows([]string{"id", "generated_at", "period", "analysis", "metrics_json", "insights_json"}).
		AddRow("test-id-1", now, string(period), "Test analysis 1", metricsJSON, insightsJSON).
		AddRow("test-id-2", now.Add(-time.Hour), string(period), "Test analysis 2", metricsJSON, insightsJSON)

	mock.ExpectQuery(`SELECT id, generated_at, period, analysis, metrics_json, insights_json FROM performance_reports WHERE period = \? ORDER BY generated_at DESC LIMIT \?`).
		WithArgs(string(period), limit).
		WillReturnRows(rows)

	// Call the method
	reports, err := repo.GetReportsByPeriod(context.Background(), period, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, reports)
	assert.Len(t, reports, 2)
	assert.Equal(t, "test-id-1", reports[0].ID)
	assert.Equal(t, "Test analysis 1", reports[0].Analysis)
	assert.Equal(t, "test-id-2", reports[1].ID)
	assert.Equal(t, "Test analysis 2", reports[1].Analysis)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetLatestReport(t *testing.T) {
	// Create dependencies
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Create test data
	now := time.Now()
	metricsJSON := `{"timestamp":"2023-04-06T12:34:56Z","metrics":{"test_metric":123.45},"system_state":{"cpu_usage":10.5,"memory_usage":256,"latency_ms":50,"goroutines":10,"uptime":"1h 30m"}}`
	insightsJSON := `["Insight 1","Insight 2"]`

	// Set up expectations
	rows := sqlmock.NewRows([]string{"id", "generated_at", "period", "analysis", "metrics_json", "insights_json"}).
		AddRow("test-id", now, "hourly", "Test analysis", metricsJSON, insightsJSON)

	mock.ExpectQuery(`SELECT id, generated_at, period, analysis, metrics_json, insights_json FROM performance_reports ORDER BY generated_at DESC LIMIT 1`).
		WillReturnRows(rows)

	// Call the method
	report, err := repo.GetLatestReport(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "test-id", report.ID)
	assert.Equal(t, "hourly", report.Period)
	assert.Equal(t, "Test analysis", report.Analysis)
	assert.Equal(t, []string{"Insight 1", "Insight 2"}, report.Insights)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_GetLatestReport_NotFound(t *testing.T) {
	// Create dependencies
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := zaptest.NewLogger(t)

	// Create repository
	repo := NewReportRepository(db, logger)

	// Set up expectations
	mock.ExpectQuery(`SELECT id, generated_at, period, analysis, metrics_json, insights_json FROM performance_reports ORDER BY generated_at DESC LIMIT 1`).
		WillReturnError(sql.ErrNoRows)

	// Call the method
	report, err := repo.GetLatestReport(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}
