package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// StatusRecord represents a status record in the database
type StatusRecord struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	Type          string `gorm:"index;not null"` // "system" or "component"
	ComponentName string `gorm:"index"`
	Status        string `gorm:"not null"`
	Message       string
	Data          []byte    `gorm:"type:blob"`
	CreatedAt     time.Time `gorm:"index;not null"`
}

// StatusRepository implements the SystemStatusRepository interface
type StatusRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewStatusRepository creates a new status repository
func NewStatusRepository(db *gorm.DB, logger *zerolog.Logger) *StatusRepository {
	return &StatusRepository{
		db:     db,
		logger: logger,
	}
}

// GetDB returns the database connection
func (r *StatusRepository) GetDB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// SaveSystemStatus saves the current system status
func (r *StatusRepository) SaveSystemStatus(ctx context.Context, systemStatus *status.SystemStatus) error {
	data, err := json.Marshal(systemStatus)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to marshal system status")
		return err
	}

	record := StatusRecord{
		Type:      "system",
		Status:    string(systemStatus.Status),
		Data:      data,
		CreatedAt: time.Now(),
	}

	result := r.GetDB(ctx).Create(&record)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to save system status")
		return result.Error
	}

	return nil
}

// GetSystemStatus retrieves the current system status
func (r *StatusRepository) GetSystemStatus(ctx context.Context) (*status.SystemStatus, error) {
	var record StatusRecord
	result := r.GetDB(ctx).
		Where("type = ?", "system").
		Order("created_at DESC").
		First(&record)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Debug().Msg("No system status found")
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Msg("Failed to get system status")
		return nil, result.Error
	}

	var systemStatus status.SystemStatus
	if err := json.Unmarshal(record.Data, &systemStatus); err != nil {
		r.logger.Error().Err(err).Msg("Failed to unmarshal system status")
		return nil, err
	}

	return &systemStatus, nil
}

// SaveComponentStatus saves a component status
func (r *StatusRepository) SaveComponentStatus(ctx context.Context, componentStatus *status.ComponentStatus) error {
	data, err := json.Marshal(componentStatus)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to marshal component status")
		return err
	}

	record := StatusRecord{
		Type:          "component",
		ComponentName: componentStatus.Name,
		Status:        string(componentStatus.Status),
		Message:       componentStatus.Message,
		Data:          data,
		CreatedAt:     time.Now(),
	}

	result := r.GetDB(ctx).Create(&record)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to save component status")
		return result.Error
	}

	return nil
}

// GetComponentStatus retrieves a component status by name
func (r *StatusRepository) GetComponentStatus(ctx context.Context, name string) (*status.ComponentStatus, error) {
	var record StatusRecord
	result := r.GetDB(ctx).
		Where("type = ? AND component_name = ?", "component", name).
		Order("created_at DESC").
		First(&record)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Debug().Str("component", name).Msg("No component status found")
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Str("component", name).Msg("Failed to get component status")
		return nil, result.Error
	}

	var componentStatus status.ComponentStatus
	if err := json.Unmarshal(record.Data, &componentStatus); err != nil {
		r.logger.Error().Err(err).Str("component", name).Msg("Failed to unmarshal component status")
		return nil, err
	}

	return &componentStatus, nil
}

// GetComponentHistory retrieves historical status for a component
func (r *StatusRepository) GetComponentHistory(ctx context.Context, name string, limit int) ([]*status.ComponentStatus, error) {
	if limit <= 0 {
		limit = 10
	}

	var records []StatusRecord
	result := r.GetDB(ctx).
		Where("type = ? AND component_name = ?", "component", name).
		Order("created_at DESC").
		Limit(limit).
		Find(&records)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("component", name).Msg("Failed to get component history")
		return nil, result.Error
	}

	history := make([]*status.ComponentStatus, 0, len(records))
	for _, record := range records {
		var componentStatus status.ComponentStatus
		if err := json.Unmarshal(record.Data, &componentStatus); err != nil {
			r.logger.Error().Err(err).Str("component", name).Msg("Failed to unmarshal component status")
			continue
		}
		history = append(history, &componentStatus)
	}

	return history, nil
}
