package reporting

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/services/gemini"
	"go.uber.org/zap"
)

// GeminiClient defines the interface for the Gemini API client
type GeminiClient interface {
	AnalyzeMetrics(ctx context.Context, metrics models.PerformanceReportMetrics) (string, error)
	ExtractInsights(ctx context.Context, analysis string) ([]string, error)
}

// ReportGenerator is a service that generates performance reports
type ReportGenerator struct {
	metrics    chan models.PerformanceReportMetrics
	geminiAPI  GeminiClient
	interval   time.Duration
	repository ReportRepository
	logger     *zap.Logger
	startTime  time.Time
	mu         sync.Mutex
}

// ReportRepository defines the interface for report persistence
type ReportRepository interface {
	SaveReport(ctx context.Context, report *models.PerformanceReport) error
	GetReportByID(ctx context.Context, id string) (*models.PerformanceReport, error)
	GetReportsByPeriod(ctx context.Context, period models.ReportPeriod, limit int) ([]*models.PerformanceReport, error)
	GetLatestReport(ctx context.Context) (*models.PerformanceReport, error)
}

// MetricsCollector defines the interface for metrics collection
type MetricsCollector interface {
	CollectMetrics(ctx context.Context, timeRanges ...time.Time) (map[string]interface{}, error)
}

// NewReportGenerator creates a new ReportGenerator
func NewReportGenerator(
	geminiAPI *gemini.GeminiClient,
	repository ReportRepository,
	interval time.Duration,
	logger *zap.Logger,
) *ReportGenerator {
	return &ReportGenerator{
		metrics:    make(chan models.PerformanceReportMetrics, 100),
		geminiAPI:  geminiAPI,
		interval:   interval,
		repository: repository,
		logger:     logger,
		startTime:  time.Now(),
	}
}

// StartCollection starts collecting metrics at regular intervals
func (r *ReportGenerator) StartCollection(ctx context.Context, collector MetricsCollector) {
	r.logger.Info("Starting metrics collection", zap.Duration("interval", r.interval))
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("Stopping metrics collection")
			return
		case <-ticker.C:
			r.logger.Debug("Collecting metrics")
			metrics, err := collector.CollectMetrics(ctx)
			if err != nil {
				r.logger.Error("Failed to collect metrics", zap.Error(err))
				continue
			}

			// Create performance metrics
			perfMetrics := models.PerformanceReportMetrics{
				Timestamp: time.Now(),
				Metrics:   metrics,
				SystemState: models.SystemState{
					CPUUsage:    getCPUUsage(),
					MemoryUsage: getMemoryUsage(),
					Latency:     getLatency(),
					Goroutines:  runtime.NumGoroutine(),
					Uptime:      getUptime(r.startTime),
				},
			}

			// Send metrics for processing
			r.metrics <- perfMetrics
		}
	}
}

// ProcessMetrics processes collected metrics
func (r *ReportGenerator) ProcessMetrics(ctx context.Context) {
	r.logger.Info("Starting metrics processing")
	for {
		select {
		case <-ctx.Done():
			r.logger.Info("Stopping metrics processing")
			return
		case metrics := <-r.metrics:
			r.logger.Debug("Processing metrics", zap.Time("timestamp", metrics.Timestamp))

			// Analyze metrics using Gemini
			analysis, err := r.geminiAPI.AnalyzeMetrics(ctx, metrics)
			if err != nil {
				r.logger.Error("Failed to analyze metrics", zap.Error(err))
				continue
			}

			// Extract insights
			insights, err := r.geminiAPI.ExtractInsights(ctx, analysis)
			if err != nil {
				r.logger.Error("Failed to extract insights", zap.Error(err))
				// Continue with empty insights
				insights = []string{}
			}

			// Generate report
			if err := r.generateReport(ctx, analysis, insights, metrics); err != nil {
				r.logger.Error("Failed to generate report", zap.Error(err))
			}
		}
	}
}

// generateReport generates a report from analyzed metrics
func (r *ReportGenerator) generateReport(
	ctx context.Context,
	analysis string,
	insights []string,
	metrics models.PerformanceReportMetrics,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Determine report period based on interval
	period := getPeriodFromInterval(r.interval)

	// Create report
	report := &models.PerformanceReport{
		ID:          uuid.New().String(),
		GeneratedAt: time.Now(),
		Period:      string(period),
		Analysis:    analysis,
		Metrics:     metrics,
		Insights:    insights,
	}

	// Save report
	if err := r.repository.SaveReport(ctx, report); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	r.logger.Info("Generated report",
		zap.String("id", report.ID),
		zap.String("period", string(period)),
		zap.Time("timestamp", metrics.Timestamp),
		zap.Int("insights", len(insights)),
	)

	return nil
}

// GetLatestReport gets the latest report
func (r *ReportGenerator) GetLatestReport(ctx context.Context) (*models.PerformanceReport, error) {
	return r.repository.GetLatestReport(ctx)
}

// GetReportByID gets a report by ID
func (r *ReportGenerator) GetReportByID(ctx context.Context, id string) (*models.PerformanceReport, error) {
	return r.repository.GetReportByID(ctx, id)
}

// GetReportsByPeriod gets reports by period
func (r *ReportGenerator) GetReportsByPeriod(ctx context.Context, period models.ReportPeriod, limit int) ([]*models.PerformanceReport, error) {
	return r.repository.GetReportsByPeriod(ctx, period, limit)
}

// Helper functions

// getCPUUsage gets the current CPU usage
func getCPUUsage() float64 {
	// This is a simplified implementation
	// In a real implementation, you would use a library like gopsutil
	return 0.0
}

// getMemoryUsage gets the current memory usage
func getMemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024 // MB
}

// getLatency gets the current latency
func getLatency() int64 {
	// This is a simplified implementation
	// In a real implementation, you would measure actual latency
	return 0
}

// getUptime gets the current uptime
func getUptime(startTime time.Time) string {
	uptime := time.Since(startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// getPeriodFromInterval determines the report period based on the interval
func getPeriodFromInterval(interval time.Duration) models.ReportPeriod {
	hours := interval.Hours()

	if hours < 1 {
		return models.ReportPeriodHourly
	} else if hours < 24 {
		return models.ReportPeriodDaily
	} else if hours < 24*7 {
		return models.ReportPeriodWeekly
	} else {
		return models.ReportPeriodMonthly
	}
}
