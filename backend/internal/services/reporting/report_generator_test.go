package reporting

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// MockGeminiClient is a mock implementation of the Gemini client
type MockGeminiClient struct {
	mock.Mock
}

func (m *MockGeminiClient) AnalyzeMetrics(ctx context.Context, metrics models.PerformanceReportMetrics) (string, error) {
	args := m.Called(ctx, metrics)
	return args.String(0), args.Error(1)
}

func (m *MockGeminiClient) ExtractInsights(ctx context.Context, analysis string) ([]string, error) {
	args := m.Called(ctx, analysis)
	return args.Get(0).([]string), args.Error(1)
}

// MockReportRepository is a mock implementation of the report repository
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) SaveReport(ctx context.Context, report *models.PerformanceReport) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *MockReportRepository) GetReportByID(ctx context.Context, id string) (*models.PerformanceReport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PerformanceReport), args.Error(1)
}

func (m *MockReportRepository) GetReportsByPeriod(ctx context.Context, period models.ReportPeriod, limit int) ([]*models.PerformanceReport, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*models.PerformanceReport), args.Error(1)
}

func (m *MockReportRepository) GetLatestReport(ctx context.Context) (*models.PerformanceReport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PerformanceReport), args.Error(1)
}

// MockMetricsCollector is a mock implementation of the metrics collector
type MockMetricsCollector struct {
	mock.Mock
}

func (m *MockMetricsCollector) CollectMetrics(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestNewReportGenerator(t *testing.T) {
	// Create dependencies
	geminiClient := &MockGeminiClient{}
	repository := &MockReportRepository{}
	logger := zaptest.NewLogger(t)
	interval := 15 * time.Minute

	// Create generator for testing
	generator := &ReportGenerator{
		metrics:    make(chan models.PerformanceReportMetrics, 100),
		geminiAPI:  geminiClient,
		interval:   interval,
		repository: repository,
		logger:     logger,
		startTime:  time.Now(),
	}

	// Assert
	assert.NotNil(t, generator)
	assert.Equal(t, geminiClient, generator.geminiAPI)
	assert.Equal(t, repository, generator.repository)
	assert.Equal(t, interval, generator.interval)
	assert.Equal(t, logger, generator.logger)
}

func TestReportGenerator_generateReport(t *testing.T) {
	// Create dependencies
	geminiClient := &MockGeminiClient{}
	repository := &MockReportRepository{}
	logger := zaptest.NewLogger(t)
	interval := 15 * time.Minute

	// Create generator for testing
	generator := &ReportGenerator{
		metrics:    make(chan models.PerformanceReportMetrics, 100),
		geminiAPI:  geminiClient,
		interval:   interval,
		repository: repository,
		logger:     logger,
		startTime:  time.Now(),
	}

	// Create test data
	ctx := context.Background()
	analysis := "Test analysis"
	insights := []string{"Insight 1", "Insight 2"}
	metrics := models.PerformanceReportMetrics{
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
	}

	// Set up expectations
	repository.On("SaveReport", ctx, mock.AnythingOfType("*models.PerformanceReport")).Return(nil)

	// Call the method
	err := generator.generateReport(ctx, analysis, insights, metrics)

	// Assert
	assert.NoError(t, err)
	repository.AssertExpectations(t)
}

func TestReportGenerator_GetLatestReport(t *testing.T) {
	// Create dependencies
	geminiClient := &MockGeminiClient{}
	repository := &MockReportRepository{}
	logger := zaptest.NewLogger(t)
	interval := 15 * time.Minute

	// Create generator for testing
	generator := &ReportGenerator{
		metrics:    make(chan models.PerformanceReportMetrics, 100),
		geminiAPI:  geminiClient,
		interval:   interval,
		repository: repository,
		logger:     logger,
		startTime:  time.Now(),
	}

	// Create test data
	ctx := context.Background()
	report := &models.PerformanceReport{
		ID:          "test-id",
		GeneratedAt: time.Now(),
		Period:      "hourly",
		Analysis:    "Test analysis",
		Insights:    []string{"Insight 1", "Insight 2"},
	}

	// Set up expectations
	repository.On("GetLatestReport", ctx).Return(report, nil)

	// Call the method
	result, err := generator.GetLatestReport(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, report, result)
	repository.AssertExpectations(t)
}

func TestReportGenerator_GetReportByID(t *testing.T) {
	// Create dependencies
	geminiClient := &MockGeminiClient{}
	repository := &MockReportRepository{}
	logger := zaptest.NewLogger(t)
	interval := 15 * time.Minute

	// Create generator for testing
	generator := &ReportGenerator{
		metrics:    make(chan models.PerformanceReportMetrics, 100),
		geminiAPI:  geminiClient,
		interval:   interval,
		repository: repository,
		logger:     logger,
		startTime:  time.Now(),
	}

	// Create test data
	ctx := context.Background()
	id := "test-id"
	report := &models.PerformanceReport{
		ID:          id,
		GeneratedAt: time.Now(),
		Period:      "hourly",
		Analysis:    "Test analysis",
		Insights:    []string{"Insight 1", "Insight 2"},
	}

	// Set up expectations
	repository.On("GetReportByID", ctx, id).Return(report, nil)

	// Call the method
	result, err := generator.GetReportByID(ctx, id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, report, result)
	repository.AssertExpectations(t)
}

func TestReportGenerator_GetReportsByPeriod(t *testing.T) {
	// Create dependencies
	geminiClient := &MockGeminiClient{}
	repository := &MockReportRepository{}
	logger := zaptest.NewLogger(t)
	interval := 15 * time.Minute

	// Create generator for testing
	generator := &ReportGenerator{
		metrics:    make(chan models.PerformanceReportMetrics, 100),
		geminiAPI:  geminiClient,
		interval:   interval,
		repository: repository,
		logger:     logger,
		startTime:  time.Now(),
	}

	// Create test data
	ctx := context.Background()
	period := models.ReportPeriodHourly
	limit := 10
	reports := []*models.PerformanceReport{
		{
			ID:          "test-id-1",
			GeneratedAt: time.Now(),
			Period:      string(period),
			Analysis:    "Test analysis 1",
			Insights:    []string{"Insight 1", "Insight 2"},
		},
		{
			ID:          "test-id-2",
			GeneratedAt: time.Now().Add(-time.Hour),
			Period:      string(period),
			Analysis:    "Test analysis 2",
			Insights:    []string{"Insight 3", "Insight 4"},
		},
	}

	// Set up expectations
	repository.On("GetReportsByPeriod", ctx, period, limit).Return(reports, nil)

	// Call the method
	result, err := generator.GetReportsByPeriod(ctx, period, limit)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, reports, result)
	repository.AssertExpectations(t)
}

func TestGetPeriodFromInterval(t *testing.T) {
	// Test hourly
	assert.Equal(t, models.ReportPeriodHourly, getPeriodFromInterval(30*time.Minute))

	// Test daily
	assert.Equal(t, models.ReportPeriodDaily, getPeriodFromInterval(2*time.Hour))

	// Test weekly
	assert.Equal(t, models.ReportPeriodWeekly, getPeriodFromInterval(2*24*time.Hour))

	// Test monthly
	assert.Equal(t, models.ReportPeriodMonthly, getPeriodFromInterval(8*24*time.Hour))
}
