package gorm

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// DatabaseStatusProvider provides status information for the database
type DatabaseStatusProvider struct {
	db            *gorm.DB
	logger        *zerolog.Logger
	lastCheckTime time.Time
	isRunning     bool
	name          string
}

// NewDatabaseStatusProvider creates a new database status provider
func NewDatabaseStatusProvider(db *gorm.DB, logger *zerolog.Logger) *DatabaseStatusProvider {
	return &DatabaseStatusProvider{
		db:        db,
		logger:    logger,
		isRunning: true,
		name:      "database",
	}
}

// GetStatus returns the current status of the database
func (p *DatabaseStatusProvider) GetStatus(ctx context.Context) (*status.ComponentStatus, error) {
	componentStatus := status.NewComponentStatus(p.name, status.StatusUnknown)
	now := time.Now()
	p.lastCheckTime = now

	// Check if the database is available by pinging it
	sqlDB, err := p.db.DB()
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to get database connection")
		componentStatus.Status = status.StatusError
		componentStatus.Message = "Failed to get database connection"
		componentStatus.LastError = err.Error()
		p.isRunning = false
		return componentStatus, nil
	}

	// Ping the database with context
	if err := sqlDB.PingContext(ctx); err != nil {
		p.logger.Error().Err(err).Msg("Database ping failed")
		componentStatus.Status = status.StatusError
		componentStatus.Message = "Database connection failed"
		componentStatus.LastError = err.Error()
		p.isRunning = false
		return componentStatus, nil
	}

	// Get database stats
	stats := sqlDB.Stats()

	// Database is responsive
	p.isRunning = true
	componentStatus.Status = status.StatusRunning
	componentStatus.Message = "Database connection is healthy"

	// Add some metrics
	componentStatus.AddMetric("open_connections", stats.OpenConnections)
	componentStatus.AddMetric("in_use", stats.InUse)
	componentStatus.AddMetric("idle", stats.Idle)
	componentStatus.AddMetric("max_open_connections", stats.MaxOpenConnections)
	// MaxIdleConnections is not available in sql.DBStats, use db.DB().SetMaxIdleConns() value instead
	maxIdleConns := 0
	if sqlDB, err := p.db.DB(); err == nil {
		maxIdleConns = sqlDB.Stats().Idle
	}
	componentStatus.AddMetric("max_idle_connections", maxIdleConns)
	componentStatus.AddMetric("wait_count", stats.WaitCount)
	componentStatus.AddMetric("wait_duration", stats.WaitDuration.String())
	componentStatus.AddMetric("max_idle_closed", stats.MaxIdleClosed)
	componentStatus.AddMetric("max_lifetime_closed", stats.MaxLifetimeClosed)
	componentStatus.AddMetric("last_check_time", now.Format(time.RFC3339))

	// Check for warning conditions
	if stats.OpenConnections > int(float64(stats.MaxOpenConnections)*0.8) {
		componentStatus.Status = status.StatusWarning
		componentStatus.Message = "Database connection pool nearing capacity"
	}

	return componentStatus, nil
}

// GetName returns the name of the component
func (p *DatabaseStatusProvider) GetName() string {
	return p.name
}

// IsRunning returns true if the component is running
func (p *DatabaseStatusProvider) IsRunning() bool {
	return p.isRunning
}
