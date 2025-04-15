package gorm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BaseRepository provides common functionality for GORM repositories
type BaseRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository(db *gorm.DB, logger *zerolog.Logger) BaseRepository {
	return BaseRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new entity into the database
func (r *BaseRepository) Create(ctx context.Context, entity interface{}) error {
	result := r.GetDB(ctx).Create(entity)
	if result.Error != nil {
		r.logError(result.Error, "Failed to create entity", entity)
		return result.Error
	}
	return nil
}

// CreateInBatches inserts multiple entities in batches
func (r *BaseRepository) CreateInBatches(ctx context.Context, entities interface{}, batchSize int) error {
	result := r.GetDB(ctx).CreateInBatches(entities, batchSize)
	if result.Error != nil {
		r.logError(result.Error, "Failed to create entities in batches", entities)
		return result.Error
	}
	return nil
}

// Save updates an entity or creates it if it doesn't exist
func (r *BaseRepository) Save(ctx context.Context, entity interface{}) error {
	result := r.GetDB(ctx).Save(entity)
	if result.Error != nil {
		r.logError(result.Error, "Failed to save entity", entity)
		return result.Error
	}
	return nil
}

// Update updates an entity with the given fields
func (r *BaseRepository) Update(ctx context.Context, entity interface{}, updates map[string]interface{}) error {
	result := r.GetDB(ctx).Model(entity).Updates(updates)
	if result.Error != nil {
		r.logError(result.Error, "Failed to update entity", entity)
		return result.Error
	}
	return nil
}

// Upsert inserts or updates an entity based on conflict columns
func (r *BaseRepository) Upsert(ctx context.Context, entity interface{}, conflictColumns []string, updateColumns []string) error {
	// Create clauses for the conflict resolution
	var columns []clause.Column
	for _, col := range conflictColumns {
		columns = append(columns, clause.Column{Name: col})
	}

	// Determine what to update on conflict
	var onConflict clause.OnConflict
	if len(updateColumns) == 0 {
		// Update all columns if none specified
		onConflict = clause.OnConflict{
			Columns:   columns,
			UpdateAll: true,
		}
	} else {
		// Update only specified columns
		onConflict = clause.OnConflict{
			Columns:   columns,
			DoUpdates: clause.AssignmentColumns(updateColumns),
		}
	}

	// Execute the upsert
	result := r.GetDB(ctx).Clauses(onConflict).Create(entity)
	if result.Error != nil {
		r.logError(result.Error, "Failed to upsert entity", entity)
		return result.Error
	}
	return nil
}

// Delete removes an entity from the database
func (r *BaseRepository) Delete(ctx context.Context, entity interface{}) error {
	result := r.GetDB(ctx).Delete(entity)
	if result.Error != nil {
		r.logError(result.Error, "Failed to delete entity", entity)
		return result.Error
	}
	return nil
}

// DeleteByID removes an entity by ID
func (r *BaseRepository) DeleteByID(ctx context.Context, model interface{}, id interface{}) error {
	result := r.GetDB(ctx).Delete(model, id)
	if result.Error != nil {
		r.logger.Error().
			Err(result.Error).
			Interface("id", id).
			Str("model", reflect.TypeOf(model).String()).
			Msg("Failed to delete entity by ID")
		return result.Error
	}
	return nil
}

// FindByID retrieves an entity by ID
func (r *BaseRepository) FindByID(ctx context.Context, entity interface{}, id interface{}) error {
	result := r.GetDB(ctx).First(entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // Return nil for not found to match interface expectations
		}
		r.logger.Error().
			Err(result.Error).
			Interface("id", id).
			Str("model", reflect.TypeOf(entity).String()).
			Msg("Failed to find entity by ID")
		return result.Error
	}
	return nil
}

// FindOne retrieves a single entity matching the given conditions
func (r *BaseRepository) FindOne(ctx context.Context, entity interface{}, conditions interface{}, args ...interface{}) error {
	result := r.GetDB(ctx).Where(conditions, args...).First(entity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // Return nil for not found to match interface expectations
		}
		r.logger.Error().
			Err(result.Error).
			Interface("conditions", conditions).
			Str("model", reflect.TypeOf(entity).String()).
			Msg("Failed to find entity")
		return result.Error
	}
	return nil
}

// FindAll retrieves all entities matching the given conditions
func (r *BaseRepository) FindAll(ctx context.Context, entities interface{}, conditions interface{}, args ...interface{}) error {
	result := r.GetDB(ctx).Where(conditions, args...).Find(entities)
	if result.Error != nil {
		r.logger.Error().
			Err(result.Error).
			Interface("conditions", conditions).
			Str("model", reflect.TypeOf(entities).String()).
			Msg("Failed to find entities")
		return result.Error
	}
	return nil
}

// FindAllWithPagination retrieves entities with pagination
func (r *BaseRepository) FindAllWithPagination(ctx context.Context, entities interface{}, page, limit int, conditions interface{}, args ...interface{}) error {
	offset := (page - 1) * limit
	result := r.GetDB(ctx).Where(conditions, args...).Offset(offset).Limit(limit).Find(entities)
	if result.Error != nil {
		r.logger.Error().
			Err(result.Error).
			Interface("conditions", conditions).
			Int("page", page).
			Int("limit", limit).
			Str("model", reflect.TypeOf(entities).String()).
			Msg("Failed to find entities with pagination")
		return result.Error
	}
	return nil
}

// Count returns the number of entities matching the given conditions
func (r *BaseRepository) Count(ctx context.Context, model interface{}, count *int64, conditions interface{}, args ...interface{}) error {
	result := r.GetDB(ctx).Model(model).Where(conditions, args...).Count(count)
	if result.Error != nil {
		r.logger.Error().
			Err(result.Error).
			Interface("conditions", conditions).
			Str("model", reflect.TypeOf(model).String()).
			Msg("Failed to count entities")
		return result.Error
	}
	return nil
}

// Transaction executes operations within a database transaction
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.GetDB(ctx).Transaction(fn)
}

// WithTransaction returns a new BaseRepository that uses the given transaction
func (r *BaseRepository) WithTransaction(tx *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db:     tx,
		logger: r.logger,
	}
}

// GetDB returns the database instance, using the transaction from context if available
func (r *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return r.db
}

// FindAllInTimeRange retrieves entities within a time range
func (r *BaseRepository) FindAllInTimeRange(ctx context.Context, entities interface{}, timeField string, from, to time.Time, conditions interface{}, args ...interface{}) error {
	// Create time range condition
	timeCondition := fmt.Sprintf("%s BETWEEN ? AND ?", timeField)
	timeArgs := []interface{}{from, to}

	// Combine with other conditions if provided
	var finalCondition string
	var finalArgs []interface{}

	if conditions != nil {
		if strCond, ok := conditions.(string); ok {
			finalCondition = fmt.Sprintf("(%s) AND (%s)", strCond, timeCondition)
			finalArgs = append(args, timeArgs...)
		} else {
			r.logger.Error().
				Interface("conditions", conditions).
				Msg("Invalid conditions type for FindAllInTimeRange")
			return fmt.Errorf("invalid conditions type for FindAllInTimeRange")
		}
	} else {
		finalCondition = timeCondition
		finalArgs = timeArgs
	}

	// Execute the query
	result := r.GetDB(ctx).Where(finalCondition, finalArgs...).Find(entities)
	if result.Error != nil {
		r.logger.Error().
			Err(result.Error).
			Str("timeField", timeField).
			Time("from", from).
			Time("to", to).
			Interface("conditions", conditions).
			Msg("Failed to find entities in time range")
		return result.Error
	}
	return nil
}

// logError logs an error with context
func (r *BaseRepository) logError(err error, message string, entity interface{}) {
	r.logger.Error().
		Err(err).
		Interface("entity", entity).
		Msg(message)
}
