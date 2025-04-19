package repo

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
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

// GetDB returns the database connection, potentially retrieving it from context if using transactions
func (r *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	// // Uncomment this section if using TransactionManager pattern from the old backend
	// // Check if there's a transaction in the context
	// if tx, ok := ctx.Value(port.TxContextKey).(*gorm.DB); ok && tx != nil {
	// 	return tx
	// }
	// Otherwise, return the regular DB connection with context
	return r.db.WithContext(ctx)
}

// Create inserts a new entity into the database
func (r *BaseRepository) Create(ctx context.Context, entity interface{}) error {
	result := r.GetDB(ctx).Create(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to create entity")
		return result.Error
	}
	return nil
}

// Save updates an entity or creates it if it doesn't exist (Upsert based on primary key)
func (r *BaseRepository) Save(ctx context.Context, entity interface{}) error {
	result := r.GetDB(ctx).Save(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to save entity")
		return result.Error
	}
	return nil
}

// FindByID retrieves an entity by its primary key
func (r *BaseRepository) FindByID(ctx context.Context, entity interface{}, id interface{}) error {
	result := r.GetDB(ctx).First(entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Consistent behavior: return a specific domain/repo error for not found
			// For now, returning gorm error, can be mapped later.
			return result.Error
		}
		r.logger.Error().Err(result.Error).Interface("id", id).Msg("Failed to find entity by ID")
		return result.Error
	}
	return nil
}

// Add other common methods like FindAll, Delete, Update etc. if needed
