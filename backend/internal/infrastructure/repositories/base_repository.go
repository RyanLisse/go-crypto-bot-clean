package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// BaseRepository provides common functionality for all repositories
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
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByID retrieves a record by its ID
func (r *BaseRepository) FindByID(ctx context.Context, model interface{}, id string) error {
	result := r.db.WithContext(ctx).First(model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return result.Error
	}
	return nil
}

// Update updates an existing record in the database
func (r *BaseRepository) Update(ctx context.Context, model interface{}) error {
	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete removes a record from the database
func (r *BaseRepository) Delete(ctx context.Context, model interface{}) error {
	result := r.db.WithContext(ctx).Delete(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindAll retrieves all records of a given model type with optional conditions
func (r *BaseRepository) FindAll(ctx context.Context, models interface{}, conditions ...interface{}) error {
	result := r.db.WithContext(ctx).Find(models, conditions...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Transaction executes operations within a database transaction
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// Count returns the number of records matching the given conditions
func (r *BaseRepository) Count(ctx context.Context, model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
