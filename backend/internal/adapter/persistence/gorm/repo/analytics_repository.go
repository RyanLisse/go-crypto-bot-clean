package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// AnalyticsMetricEntity represents an analytics metric in the database
type AnalyticsMetricEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	Type      string    `gorm:"index;type:varchar(50)"`
	Timestamp time.Time `gorm:"index"`
	Data      []byte    `gorm:"type:json"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// StrategyPerformanceEntity represents strategy performance metrics in the database
type StrategyPerformanceEntity struct {
	ID         string    `gorm:"primaryKey;type:varchar(50)"`
	StrategyID string    `gorm:"index;type:varchar(50)"`
	Timestamp  time.Time `gorm:"index"`
	Metrics    []byte    `gorm:"type:json"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// GormAnalyticsRepository implements port.AnalyticsRepository using GORM
type GormAnalyticsRepository struct {
	BaseRepository
}

// NewGormAnalyticsRepository creates a new GormAnalyticsRepository
func NewGormAnalyticsRepository(db *gorm.DB, logger *zerolog.Logger) *GormAnalyticsRepository {
	return &GormAnalyticsRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveMetrics saves analytics metrics
func (r *GormAnalyticsRepository) SaveMetrics(ctx context.Context, metrics map[string]interface{}) error {
	// Extract metric type
	metricType, ok := metrics["type"].(string)
	if !ok {
		r.logger.Error().Interface("metrics", metrics).Msg("Missing type in metrics")
		return fmt.Errorf("missing type in metrics")
	}

	// Convert metrics to JSON
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		r.logger.Error().Err(err).Interface("metrics", metrics).Msg("Failed to marshal metrics")
		return err
	}

	// Create entity
	entity := &AnalyticsMetricEntity{
		ID:        uuid.New().String(),
		Type:      metricType,
		Timestamp: time.Now(),
		Data:      metricsJSON,
	}

	// Save entity
	return r.Create(ctx, entity)
}

// GetMetrics retrieves analytics metrics within a time range
func (r *GormAnalyticsRepository) GetMetrics(ctx context.Context, from, to time.Time) ([]map[string]interface{}, error) {
	var entities []AnalyticsMetricEntity

	query := r.GetDB(ctx)

	// Add time range conditions
	if !from.IsZero() {
		query = query.Where("timestamp >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("timestamp <= ?", to)
	}

	// Execute query
	err := query.
		Order("timestamp DESC").
		Find(&entities).Error
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to get metrics")
		return nil, err
	}

	// Convert to maps
	metrics := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		// Parse data
		var data map[string]interface{}
		if len(entity.Data) > 0 {
			if err := json.Unmarshal(entity.Data, &data); err != nil {
				r.logger.Error().Err(err).Str("metric_id", entity.ID).Msg("Failed to unmarshal metric data")
				continue
			}
		}

		// Add metadata
		data["id"] = entity.ID
		data["timestamp"] = entity.Timestamp
		data["created_at"] = entity.CreatedAt

		metrics[i] = data
	}

	return metrics, nil
}

// GetPerformanceByStrategy retrieves performance metrics for a specific strategy
func (r *GormAnalyticsRepository) GetPerformanceByStrategy(ctx context.Context, strategyID string, from, to time.Time) (map[string]interface{}, error) {
	var entities []StrategyPerformanceEntity

	query := r.GetDB(ctx).
		Where("strategy_id = ?", strategyID)

	// Add time range conditions
	if !from.IsZero() {
		query = query.Where("timestamp >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("timestamp <= ?", to)
	}

	// Execute query
	err := query.
		Order("timestamp DESC").
		Find(&entities).Error
	if err != nil {
		r.logger.Error().Err(err).Str("strategy_id", strategyID).Msg("Failed to get strategy performance")
		return nil, err
	}

	// Aggregate metrics
	result := map[string]interface{}{
		"strategy_id": strategyID,
		"from":        from,
		"to":          to,
		"metrics":     make([]map[string]interface{}, 0, len(entities)),
	}

	// Process each performance record
	for _, entity := range entities {
		// Parse metrics
		var metrics map[string]interface{}
		if len(entity.Metrics) > 0 {
			if err := json.Unmarshal(entity.Metrics, &metrics); err != nil {
				r.logger.Error().Err(err).Str("performance_id", entity.ID).Msg("Failed to unmarshal performance metrics")
				continue
			}
		}

		// Add metadata
		metrics["timestamp"] = entity.Timestamp

		// Add to result
		result["metrics"] = append(result["metrics"].([]map[string]interface{}), metrics)
	}

	// Calculate aggregated metrics
	if metrics, ok := result["metrics"].([]map[string]interface{}); ok && len(metrics) > 0 {
		// Example aggregations
		var totalPnl float64
		var winCount, loseCount int

		for _, m := range metrics {
			if pnl, ok := m["pnl"].(float64); ok {
				totalPnl += pnl
				if pnl > 0 {
					winCount++
				} else if pnl < 0 {
					loseCount++
				}
			}
		}

		// Add aggregated metrics
		result["total_pnl"] = totalPnl
		result["win_count"] = winCount
		result["lose_count"] = loseCount
		result["total_trades"] = winCount + loseCount

		if winCount+loseCount > 0 {
			result["win_rate"] = float64(winCount) / float64(winCount+loseCount)
		}
	}

	return result, nil
}

// SaveStrategyPerformance saves performance metrics for a strategy
func (r *GormAnalyticsRepository) SaveStrategyPerformance(ctx context.Context, strategyID string, metrics map[string]interface{}) error {
	// Convert metrics to JSON
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		r.logger.Error().Err(err).Interface("metrics", metrics).Msg("Failed to marshal strategy performance metrics")
		return err
	}

	// Create entity
	entity := &StrategyPerformanceEntity{
		ID:         uuid.New().String(),
		StrategyID: strategyID,
		Timestamp:  time.Now(),
		Metrics:    metricsJSON,
	}

	// Save entity
	return r.Create(ctx, entity)
}

// Ensure GormAnalyticsRepository implements port.AnalyticsRepository
var _ port.AnalyticsRepository = (*GormAnalyticsRepository)(nil)
