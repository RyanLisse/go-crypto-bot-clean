package turso

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// BaseRepository provides common functionality for TursoDB repositories
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db: db,
	}
}

// Create inserts a new record into the database
func (r *BaseRepository) Create(ctx context.Context, model interface{}) error {
	result := r.db.WithContext(ctx).Create(model)
	return result.Error
}

// FindByID retrieves a record by its ID
func (r *BaseRepository) FindByID(ctx context.Context, model interface{}, id string) error {
	result := r.db.WithContext(ctx).First(model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return result.Error
	}
	return nil
}

// Update modifies an existing record
func (r *BaseRepository) Update(ctx context.Context, model interface{}) error {
	result := r.db.WithContext(ctx).Save(model)
	return result.Error
}

// Delete removes a record from the database
func (r *BaseRepository) Delete(ctx context.Context, model interface{}) error {
	result := r.db.WithContext(ctx).Delete(model)
	return result.Error
}

// Transaction executes operations within a database transaction
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
