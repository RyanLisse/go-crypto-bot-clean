// Package repositories contains the data access layer
package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go-crypto-bot-clean/backend/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MinimalRepository handles data access for the minimal API
type MinimalRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewMinimalRepository creates a new minimal repository
func NewMinimalRepository(db *gorm.DB, logger *zap.Logger) *MinimalRepository {
	return &MinimalRepository{
		db:     db,
		logger: logger,
	}
}

// SaveSystemInfo saves system information to the database
func (r *MinimalRepository) SaveSystemInfo(info *models.SystemInfo) error {
	if err := r.db.Create(info).Error; err != nil {
		r.logger.Error("Failed to save system info", zap.Error(err))
		return fmt.Errorf("failed to save system info: %w", err)
	}
	return nil
}

// GetSystemInfo retrieves the latest system information
func (r *MinimalRepository) GetSystemInfo() (*models.SystemInfo, error) {
	var info models.SystemInfo
	if err := r.db.Order("created_at DESC").First(&info).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get system info", zap.Error(err))
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}
	return &info, nil
}

// UpdateSystemInfo updates the system information
func (r *MinimalRepository) UpdateSystemInfo(info *models.SystemInfo) error {
	if err := r.db.Save(info).Error; err != nil {
		r.logger.Error("Failed to update system info", zap.Error(err))
		return fmt.Errorf("failed to update system info: %w", err)
	}
	return nil
}

// SaveHealthCheck saves a health check record
func (r *MinimalRepository) SaveHealthCheck(check *models.HealthCheck) error {
	if err := r.db.Create(check).Error; err != nil {
		r.logger.Error("Failed to save health check", zap.Error(err))
		return fmt.Errorf("failed to save health check: %w", err)
	}
	return nil
}

// GetHealthChecks retrieves health check records
func (r *MinimalRepository) GetHealthChecks(limit int) ([]models.HealthCheck, error) {
	var checks []models.HealthCheck
	if err := r.db.Order("timestamp DESC").Limit(limit).Find(&checks).Error; err != nil {
		r.logger.Error("Failed to get health checks", zap.Error(err))
		return nil, fmt.Errorf("failed to get health checks: %w", err)
	}
	return checks, nil
}

// GetHealthCheckByID retrieves a health check record by ID
func (r *MinimalRepository) GetHealthCheckByID(id uuid.UUID) (*models.HealthCheck, error) {
	var check models.HealthCheck
	if err := r.db.Where("id = ?", id).First(&check).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get health check by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get health check by ID: %w", err)
	}
	return &check, nil
}

// SaveLogEntry saves a log entry
func (r *MinimalRepository) SaveLogEntry(entry *models.LogEntry) error {
	if err := r.db.Create(entry).Error; err != nil {
		r.logger.Error("Failed to save log entry", zap.Error(err))
		return fmt.Errorf("failed to save log entry: %w", err)
	}
	return nil
}

// GetLogEntries retrieves log entries
func (r *MinimalRepository) GetLogEntries(limit int, level string) ([]models.LogEntry, error) {
	var entries []models.LogEntry
	query := r.db.Order("timestamp DESC").Limit(limit)
	
	if level != "" {
		query = query.Where("level = ?", level)
	}
	
	if err := query.Find(&entries).Error; err != nil {
		r.logger.Error("Failed to get log entries", zap.Error(err))
		return nil, fmt.Errorf("failed to get log entries: %w", err)
	}
	return entries, nil
}

// GetLogEntriesByTimeRange retrieves log entries within a time range
func (r *MinimalRepository) GetLogEntriesByTimeRange(start, end time.Time, level string) ([]models.LogEntry, error) {
	var entries []models.LogEntry
	query := r.db.Where("timestamp BETWEEN ? AND ?", start, end).Order("timestamp DESC")
	
	if level != "" {
		query = query.Where("level = ?", level)
	}
	
	if err := query.Find(&entries).Error; err != nil {
		r.logger.Error("Failed to get log entries by time range", zap.Error(err))
		return nil, fmt.Errorf("failed to get log entries by time range: %w", err)
	}
	return entries, nil
}

// CleanupOldLogs removes log entries older than the specified duration
func (r *MinimalRepository) CleanupOldLogs(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	result := r.db.Where("timestamp < ?", cutoff).Delete(&models.LogEntry{})
	if result.Error != nil {
		r.logger.Error("Failed to cleanup old logs", zap.Error(result.Error))
		return 0, fmt.Errorf("failed to cleanup old logs: %w", result.Error)
	}
	return result.RowsAffected, nil
}
