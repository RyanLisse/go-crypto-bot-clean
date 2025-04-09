package report

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"go-crypto-bot-clean/backend/internal/repository/database"
	"go.uber.org/zap"
)

// ReportRepository is an implementation of the report repository using our database abstraction layer
type ReportRepository struct {
	db     database.Repository
	logger *zap.Logger
}

// NewReportRepository creates a new ReportRepository
func NewReportRepository(db database.Repository, logger *zap.Logger) repository.ReportRepository {
	return &ReportRepository{
		db:     db,
		logger: logger,
	}
}

// Initialize initializes the repository
func (r *ReportRepository) Initialize(ctx context.Context) error {
	// Create the reports table
	query := `
	CREATE TABLE IF NOT EXISTS performance_reports (
		id TEXT PRIMARY KEY,
		generated_at TIMESTAMP NOT NULL,
		period TEXT NOT NULL,
		analysis TEXT NOT NULL,
		metrics_json TEXT NOT NULL,
		insights_json TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_performance_reports_period ON performance_reports(period);
	CREATE INDEX IF NOT EXISTS idx_performance_reports_generated_at ON performance_reports(generated_at);
	`

	_, err := r.db.Execute(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create reports table: %w", err)
	}

	return nil
}

// SaveReport saves a report
func (r *ReportRepository) SaveReport(ctx context.Context, report *models.PerformanceReport) error {
	// Marshal metrics to JSON
	metricsJSON, err := json.Marshal(report.Metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// Marshal insights to JSON
	insightsJSON, err := json.Marshal(report.Insights)
	if err != nil {
		return fmt.Errorf("failed to marshal insights: %w", err)
	}

	// Insert the report
	query := `
	INSERT INTO performance_reports (
		id, generated_at, period, analysis, metrics_json, insights_json
	) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Execute(
		ctx,
		query,
		report.ID,
		report.GeneratedAt,
		report.Period,
		report.Analysis,
		string(metricsJSON),
		string(insightsJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert report: %w", err)
	}

	return nil
}

// GetReportByID gets a report by ID
func (r *ReportRepository) GetReportByID(ctx context.Context, id string) (*models.PerformanceReport, error) {
	query := `
	SELECT id, generated_at, period, analysis, metrics_json, insights_json
	FROM performance_reports
	WHERE id = ?
	`

	row := r.db.QueryRow(ctx, query, id)
	return r.scanReport(row)
}

// GetReportsByPeriod gets reports by period
func (r *ReportRepository) GetReportsByPeriod(ctx context.Context, period models.ReportPeriod, limit int) ([]*models.PerformanceReport, error) {
	query := `
	SELECT id, generated_at, period, analysis, metrics_json, insights_json
	FROM performance_reports
	WHERE period = ?
	ORDER BY generated_at DESC
	LIMIT ?
	`

	rows, err := r.db.Query(ctx, query, string(period), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query reports: %w", err)
	}
	defer rows.Close()

	reports := []*models.PerformanceReport{}
	for rows.Next() {
		report, err := r.scanReportFromRows(rows)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return reports, nil
}

// GetLatestReport gets the latest report
func (r *ReportRepository) GetLatestReport(ctx context.Context) (*models.PerformanceReport, error) {
	query := `
	SELECT id, generated_at, period, analysis, metrics_json, insights_json
	FROM performance_reports
	ORDER BY generated_at DESC
	LIMIT 1
	`

	row := r.db.QueryRow(ctx, query)
	return r.scanReport(row)
}

// scanReport scans a report from a row
func (r *ReportRepository) scanReport(row *sql.Row) (*models.PerformanceReport, error) {
	var (
		id          string
		generatedAt time.Time
		period      string
		analysis    string
		metricsJSON string
		insightsJSON string
	)

	err := row.Scan(&id, &generatedAt, &period, &analysis, &metricsJSON, &insightsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan report: %w", err)
	}

	// Unmarshal metrics
	var metrics models.PerformanceReportMetrics
	if err := json.Unmarshal([]byte(metricsJSON), &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	// Unmarshal insights
	var insights []string
	if err := json.Unmarshal([]byte(insightsJSON), &insights); err != nil {
		return nil, fmt.Errorf("failed to unmarshal insights: %w", err)
	}

	return &models.PerformanceReport{
		ID:          id,
		GeneratedAt: generatedAt,
		Period:      period,
		Analysis:    analysis,
		Metrics:     metrics,
		Insights:    insights,
	}, nil
}

// scanReportFromRows scans a report from rows
func (r *ReportRepository) scanReportFromRows(rows *sql.Rows) (*models.PerformanceReport, error) {
	var (
		id          string
		generatedAt time.Time
		period      string
		analysis    string
		metricsJSON string
		insightsJSON string
	)

	err := rows.Scan(&id, &generatedAt, &period, &analysis, &metricsJSON, &insightsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to scan report: %w", err)
	}

	// Unmarshal metrics
	var metrics models.PerformanceReportMetrics
	if err := json.Unmarshal([]byte(metricsJSON), &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	// Unmarshal insights
	var insights []string
	if err := json.Unmarshal([]byte(insightsJSON), &insights); err != nil {
		return nil, fmt.Errorf("failed to unmarshal insights: %w", err)
	}

	return &models.PerformanceReport{
		ID:          id,
		GeneratedAt: generatedAt,
		Period:      period,
		Analysis:    analysis,
		Metrics:     metrics,
		Insights:    insights,
	}, nil
}
