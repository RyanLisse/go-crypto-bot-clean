package repository

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// ReportRepository defines the interface for report persistence
type ReportRepository interface {
	// Initialize initializes the repository
	Initialize(ctx context.Context) error
	
	// SaveReport saves a report
	SaveReport(ctx context.Context, report *models.PerformanceReport) error
	
	// GetReportByID gets a report by ID
	GetReportByID(ctx context.Context, id string) (*models.PerformanceReport, error)
	
	// GetReportsByPeriod gets reports by period
	GetReportsByPeriod(ctx context.Context, period models.ReportPeriod, limit int) ([]*models.PerformanceReport, error)
	
	// GetLatestReport gets the latest report
	GetLatestReport(ctx context.Context) (*models.PerformanceReport, error)
}
