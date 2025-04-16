# Repository Standardization Strategy

This document outlines our approach to standardizing repository implementations across the application, focusing on consistency, maintainability, and proper abstraction.

## Core Repository Pattern

### Repository Interfaces

All repository interfaces will follow these principles:

1. Define in `internal/domain/port` directory
2. Focus on domain operations, not persistence details
3. Use domain models as parameters and return values
4. Use context for cancellation and tracing

Example:

```go
// internal/domain/port/user_repository.go

package port

import (
	"context"
	
	"github.com/your-org/your-app/internal/domain/model"
	"github.com/your-org/your-app/internal/pkg/errors"
)

// UserRepository defines operations for persisting and retrieving users
type UserRepository interface {
	// Create persists a new user
	Create(ctx context.Context, user *model.User) (*model.User, error)
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*model.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *model.User) (*model.User, error)
	
	// Delete removes a user
	Delete(ctx context.Context, id string) error
	
	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*model.User, int, error)
}
```

### Base Repository Implementation

For each persistence mechanism (e.g., GORM, MongoDB), define a base repository with common functionality:

```go
// internal/adapter/repository/gorm/base_repository.go

package gorm

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	appErrors "github.com/your-org/your-app/internal/pkg/errors"
)

// BaseRepository provides common functionality for GORM repositories
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{db: db}
}

// DB returns the database connection
func (r *BaseRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// HandleError maps GORM errors to domain errors
func (r *BaseRepository) HandleError(err error) error {
	if err == nil {
		return nil
	}
	
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return appErrors.ErrNotFound
	default:
		return err
	}
}
```

### Specific Repository Implementations

Each repository implementation will:

1. Embed the base repository
2. Implement the domain interface
3. Focus on mapping between entity and domain models
4. Use consistent error handling

Example:

```go
// internal/adapter/repository/gorm/user_repository.go

package gorm

import (
	"context"
	"fmt"
	
	"github.com/your-org/your-app/internal/adapter/repository/gorm/entity"
	"github.com/your-org/your-app/internal/domain/model"
	"github.com/your-org/your-app/internal/domain/port"
	"github.com/your-org/your-app/internal/pkg/errors"
)

// UserRepository implements the user repository using GORM
type UserRepository struct {
	BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(baseRepo BaseRepository) port.UserRepository {
	return &UserRepository{BaseRepository: baseRepo}
}

// Create persists a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.NewBadRequestError("User cannot be nil")
	}
	
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	userEntity := entity.FromUserModel(user)
	
	if result := r.DB(ctx).Create(&userEntity); result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %w", result.Error)
	}
	
	return userEntity.ToModel(), nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.NewBadRequestError("User ID cannot be empty")
	}
	
	var userEntity entity.User
	if result := r.DB(ctx).First(&userEntity, "id = ?", id); result.Error != nil {
		return nil, r.HandleError(result.Error)
	}
	
	return userEntity.ToModel(), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, errors.NewBadRequestError("Email cannot be empty")
	}
	
	var userEntity entity.User
	if result := r.DB(ctx).First(&userEntity, "email = ?", email); result.Error != nil {
		return nil, r.HandleError(result.Error)
	}
	
	return userEntity.ToModel(), nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.NewBadRequestError("User cannot be nil")
	}
	
	if user.ID == "" {
		return nil, errors.NewBadRequestError("User ID cannot be empty")
	}
	
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	// First check if user exists
	var existingUser entity.User
	if result := r.DB(ctx).First(&existingUser, "id = ?", user.ID); result.Error != nil {
		return nil, r.HandleError(result.Error)
	}
	
	// Update user
	userEntity := entity.FromUserModel(user)
	if result := r.DB(ctx).Save(&userEntity); result.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", result.Error)
	}
	
	return userEntity.ToModel(), nil
}

// Delete removes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.NewBadRequestError("User ID cannot be empty")
	}
	
	result := r.DB(ctx).Delete(&entity.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("User with ID %s not found", id)
	}
	
	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*model.User, int, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	
	if offset < 0 {
		offset = 0
	}
	
	var userEntities []entity.User
	var total int64
	
	// Get total count
	if err := r.DB(ctx).Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	// Get paginated results
	if err := r.DB(ctx).Offset(offset).Limit(limit).Find(&userEntities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	
	// Map to domain models
	users := make([]*model.User, len(userEntities))
	for i, userEntity := range userEntities {
		users[i] = userEntity.ToModel()
	}
	
	return users, int(total), nil
}
```

## Entity-Model Mapping

### Entities

Entities represent database structures and should:

1. Be defined in `internal/adapter/repository/{orm}/entity` directory
2. Include all necessary ORM tags and hooks
3. Include methods to convert to/from domain models

Example:

```go
// internal/adapter/repository/gorm/entity/user.go

package entity

import (
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"github.com/your-org/your-app/internal/domain/model"
)

// User represents a user in the database
type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)"`
	Email     string         `gorm:"uniqueIndex;type:varchar(255)"`
	Name      string         `gorm:"type:varchar(255)"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate is called before creating a new record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// ToModel converts the entity to a domain model
func (u *User) ToModel() *model.User {
	return &model.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromUserModel creates an entity from a domain model
func FromUserModel(user *model.User) *User {
	entity := &User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
	}
	
	// Preserve these fields if they exist in the model
	if !user.CreatedAt.IsZero() {
		entity.CreatedAt = user.CreatedAt
	}
	
	if !user.UpdatedAt.IsZero() {
		entity.UpdatedAt = user.UpdatedAt
	}
	
	return entity
}
```

## Transaction Management

Implement a transaction manager to ensure consistency:

```go
// internal/adapter/repository/gorm/transaction.go

package gorm

import (
	"context"
	"fmt"
	
	"gorm.io/gorm"
	
	"github.com/your-org/your-app/internal/domain/port"
)

// TransactionManager implements the transaction manager using GORM
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) port.TransactionManager {
	return &TransactionManager{db: db}
}

// Transaction executes the given function in a transaction
func (tm *TransactionManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		txCtx := context.WithValue(ctx, "tx", tx)
		
		// Execute the function
		if err := fn(txCtx); err != nil {
			return err
		}
		
		return nil
	})
}

// DB returns the transaction from context if it exists, otherwise the regular DB
func (tm *TransactionManager) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return tm.db.WithContext(ctx)
}
```

## Repository Factory

Create a repository factory to simplify dependency injection:

```go
// internal/adapter/repository/gorm/factory.go

package gorm

import (
	"github.com/your-org/your-app/internal/domain/port"
	"gorm.io/gorm"
)

// RepositoryFactory creates GORM repositories
type RepositoryFactory struct {
	db *gorm.DB
	baseRepository BaseRepository
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	baseRepo := NewBaseRepository(db)
	return &RepositoryFactory{
		db:             db,
		baseRepository: baseRepo,
	}
}

// NewUserRepository creates a new user repository
func (f *RepositoryFactory) NewUserRepository() port.UserRepository {
	return NewUserRepository(f.baseRepository)
}

// NewOrderRepository creates a new order repository
func (f *RepositoryFactory) NewOrderRepository() port.OrderRepository {
	return NewOrderRepository(f.baseRepository)
}

// NewTransactionManager creates a new transaction manager
func (f *RepositoryFactory) NewTransactionManager() port.TransactionManager {
	return NewTransactionManager(f.db)
}
```

## Mock Repositories

For testing, create mock implementations in a separate package:

```go
// internal/adapter/repository/mock/user_repository.go

package mock

import (
	"context"
	"sync"
	
	"github.com/google/uuid"
	
	"github.com/your-org/your-app/internal/domain/model"
	"github.com/your-org/your-app/internal/domain/port"
	"github.com/your-org/your-app/internal/pkg/errors"
)

// UserRepository implements a mock user repository for testing
type UserRepository struct {
	users  map[string]*model.User
	mutex  sync.RWMutex
	mockBehavior MockBehavior // Controls error injection for testing
}

// NewUserRepository creates a new mock user repository
func NewUserRepository() port.UserRepository {
	return &UserRepository{
		users: make(map[string]*model.User),
	}
}

// Create persists a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if r.mockBehavior != nil && r.mockBehavior.ShouldError("Create") {
		return nil, r.mockBehavior.Error("Create")
	}
	
	if user == nil {
		return nil, errors.NewBadRequestError("User cannot be nil")
	}
	
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if email exists
	for _, u := range r.users {
		if u.Email == user.Email {
			return nil, errors.NewConflictError("User with email %s already exists", user.Email)
		}
	}
	
	// Create a copy of the user
	newUser := *user
	
	// Generate ID if not provided
	if newUser.ID == "" {
		newUser.ID = uuid.New().String()
	}
	
	r.users[newUser.ID] = &newUser
	
	// Return a copy to prevent modification
	result := newUser
	return &result, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	if r.mockBehavior != nil && r.mockBehavior.ShouldError("GetByID") {
		return nil, r.mockBehavior.Error("GetByID")
	}
	
	if id == "" {
		return nil, errors.NewBadRequestError("User ID cannot be empty")
	}
	
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, errors.NewNotFoundError("User with ID %s not found", id)
	}
	
	// Return a copy to prevent modification
	result := *user
	return &result, nil
}

// Implement other methods...
```

## Testing Support

Create a mock behavior controller for testing edge cases:

```go
// internal/adapter/repository/mock/mock_behavior.go

package mock

import (
	"sync"
	
	"github.com/your-org/your-app/internal/pkg/errors"
)

// MockBehavior controls repository behavior for testing
type MockBehavior struct {
	errorMap map[string]error
	mutex    sync.RWMutex
}

// NewMockBehavior creates a new mock behavior controller
func NewMockBehavior() *MockBehavior {
	return &MockBehavior{
		errorMap: make(map[string]error),
	}
}

// SetError sets an error for a specific method
func (mb *MockBehavior) SetError(method string, err error) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	mb.errorMap[method] = err
}

// SetNotFoundError sets a not found error for a specific method
func (mb *MockBehavior) SetNotFoundError(method string) {
	mb.SetError(method, errors.NewNotFoundError("Resource not found"))
}

// ClearError clears the error for a specific method
func (mb *MockBehavior) ClearError(method string) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	delete(mb.errorMap, method)
}

// ClearAllErrors clears all errors
func (mb *MockBehavior) ClearAllErrors() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()
	mb.errorMap = make(map[string]error)
}

// ShouldError returns true if an error is set for the method
func (mb *MockBehavior) ShouldError(method string) bool {
	mb.mutex.RLock()
	defer mb.mutex.RUnlock()
	_, exists := mb.errorMap[method]
	return exists
}

// Error returns the error for a method
func (mb *MockBehavior) Error(method string) error {
	mb.mutex.RLock()
	defer mb.mutex.RUnlock()
	return mb.errorMap[method]
}
```

## Repository Factory Interface

Define a factory interface in the domain layer:

```go
// internal/domain/port/repository_factory.go

package port

import (
	"context"
)

// RepositoryFactory creates repositories
type RepositoryFactory interface {
	// NewUserRepository creates a new user repository
	NewUserRepository() UserRepository
	
	// NewOrderRepository creates a new order repository
	NewOrderRepository() OrderRepository
	
	// NewTransactionManager creates a new transaction manager
	NewTransactionManager() TransactionManager
}

// TransactionManager manages database transactions
type TransactionManager interface {
	// Transaction executes the given function in a transaction
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

## Selection Mechanism for Tests

Create a mechanism to switch between real and mock repositories:

```go
// internal/adapter/repository/factory.go

package repository

import (
	"github.com/your-org/your-app/internal/adapter/repository/gorm"
	"github.com/your-org/your-app/internal/adapter/repository/mock"
	"github.com/your-org/your-app/internal/domain/port"
)

// FactoryType defines the type of repository factory
type FactoryType string

const (
	// GormFactoryType uses GORM for database access
	GormFactoryType FactoryType = "gorm"
	
	// MockFactoryType uses mock repositories for testing
	MockFactoryType FactoryType = "mock"
)

// FactoryConfig configures the repository factory
type FactoryConfig struct {
	// Type is the type of repository factory to create
	Type FactoryType
	
	// DB is the GORM database connection (required for GORM factory)
	DB any
	
	// IsProd indicates if this is a production environment
	IsProd bool
}

// Factory creates the appropriate repository factory
func Factory(config FactoryConfig) (port.RepositoryFactory, error) {
	if config.IsProd && config.Type == MockFactoryType {
		// Prevent using mock repositories in production
		panic("Mock repositories cannot be used in production")
	}
	
	switch config.Type {
	case GormFactoryType:
		db, ok := config.DB.(*gorm.DB)
		if !ok {
			return nil, errors.New("DB must be a *gorm.DB for GORM factory")
		}
		return gorm.NewRepositoryFactory(db), nil
		
	case MockFactoryType:
		return mock.NewRepositoryFactory(), nil
		
	default:
		return nil, errors.New("Unsupported repository factory type: " + string(config.Type))
	}
}
```

## Migration Plan

To transition from multiple repository implementations to this standardized approach:

1. Create the base repository structures and interfaces
2. Implement one repository at a time, starting with the most critical ones
3. Update service layer to use the new repositories
4. Write comprehensive tests for each repository
5. Gradually replace old repositories with new implementations
6. Remove deprecated repository code

## Conclusion

This standardized repository approach provides several benefits:

1. **Consistency**: Uniform repository implementation across the application
2. **Testability**: Easy mocking and testing of repositories
3. **Maintainability**: Clear separation of concerns between domain and persistence
4. **Flexibility**: Ability to switch persistence mechanisms without changing domain code
5. **Security**: Prevention of using mock repositories in production
6. **Performance**: Consistent transaction management and connection handling 