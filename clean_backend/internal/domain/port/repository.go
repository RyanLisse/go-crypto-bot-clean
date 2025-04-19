package port

import (
	"context"
)

// BaseRepository defines common methods for all repositories
type BaseRepository interface {
	// Create inserts a new entity into the database
	Create(ctx context.Context, entity interface{}) error

	// Save updates an entity or creates it if it doesn't exist
	Save(ctx context.Context, entity interface{}) error

	// FindByID retrieves an entity by ID
	FindByID(ctx context.Context, entity interface{}, id interface{}) error

	// FindAll retrieves all entities matching the given conditions
	FindAll(ctx context.Context, entities interface{}, conditions interface{}, args ...interface{}) error

	// FindAllWithPagination retrieves entities with pagination
	FindAllWithPagination(ctx context.Context, entities interface{}, page, limit int, conditions interface{}, args ...interface{}) error

	// Count returns the number of entities matching the given conditions
	Count(ctx context.Context, model interface{}, count *int64, conditions interface{}, args ...interface{}) error

	// Delete removes an entity from the database
	Delete(ctx context.Context, entity interface{}) error

	// Transaction executes operations within a database transaction
	Transaction(ctx context.Context, fn func(tx interface{}) error) error

	// FindOne retrieves a single entity matching the given conditions
	FindOne(ctx context.Context, entity interface{}, conditions interface{}, args ...interface{}) error

	// DeleteByID removes an entity by ID
	DeleteByID(ctx context.Context, model interface{}, id interface{}) error

	// Update updates an entity with the given fields
	Update(ctx context.Context, entity interface{}, updates map[string]interface{}) error
}
