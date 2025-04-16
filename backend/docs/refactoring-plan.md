# HTTP Layer Unification (Gin to Chi) - Implementation Plan

## Current State Analysis

1. **Router Definition**: 
   - The main entry point in `cmd/server/main.go` already uses Chi
   - A Gin router also exists in `internal/adapter/delivery/http/router.go`

2. **Handler Types**:
   - Some handlers use Chi (e.g., `position_handler_chi.go`, `market_data.go`)
   - Others use Gin (e.g., `position_handler.go`, `risk_handler.go`, `autobuy_handler.go`)

3. **Middleware**:
   - Chi middleware exists in `internal/adapter/http/middleware/middleware.go`
   - Gin middleware exists in various handler files and in the router.go

## Implementation Steps

### 1. Remove or Rename Gin-based Router

```bash
# Delete Gin router
rm internal/adapter/delivery/http/router.go
```

### 2. Convert Each Gin Handler to Chi

For each handler using Gin, follow this pattern:

#### Example Migration Pattern (converting risk_handler.go)

Original Gin-based handler:
```go
// RegisterRoutes registers risk-related routes with the Gin engine
func (h *RiskHandler) RegisterRoutes(router *gin.RouterGroup) {
    riskGroup := router.Group("/risk")
    {
        riskGroup.GET("/profile", h.GetRiskProfile)
        // ...other routes
    }
}

// GetRiskProfile handles profile requests
func (h *RiskHandler) GetRiskProfile(c *gin.Context) {
    userID := getUserIDFromContext(c.Request)
    // ...handler logic
    c.JSON(http.StatusOK, response.Success(profile))
}
```

New Chi-based handler:
```go
// RegisterRoutes registers risk-related routes with the Chi router
func (h *RiskHandler) RegisterRoutes(r chi.Router) {
    r.Route("/risk", func(r chi.Router) {
        r.Get("/profile", h.GetRiskProfile)
        // ...other routes
    })
}

// GetRiskProfile handles profile requests
func (h *RiskHandler) GetRiskProfile(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r)
    // ...handler logic
    response.WriteJSON(w, http.StatusOK, response.Success(profile))
}
```

### 3. Specific Handler Files to Convert

1. **Position Handler** (priority since there are two versions):
   - Remove `position_handler.go` (Gin version)
   - Rename `position_handler_chi.go` to `position_handler.go`

2. **Risk Handler**:
   - Convert `risk_handler.go` from Gin to Chi

3. **Autobuy Handler**:
   - Convert `autobuy_handler.go` from Gin to Chi

4. **Websocket Handler**:
   - Convert `websocket.go` from Gin to Chi

### 4. Update Handler Tests

Test files like `position_handler_test.go` need to be updated to use Chi instead of Gin.

Example update:
```go
// Before: Gin setup
router := gin.Default()
router.GET("/positions/:id", handler.GetPositionByID)
req, _ := http.NewRequest("GET", "/positions/123", nil)
recorder := httptest.NewRecorder()
router.ServeHTTP(recorder, req)

// After: Chi setup
router := chi.NewRouter()
router.Get("/positions/{id}", handler.GetPositionByID)
req, _ := http.NewRequest("GET", "/positions/123", nil)
recorder := httptest.NewRecorder()
router.ServeHTTP(recorder, req)
```

### 5. Update Dependencies

1. Remove Gin dependency from go.mod:
```bash
go mod edit -droprequire=github.com/gin-gonic/gin
go mod tidy
```

### 6. Update Main Server File

Update `cmd/server/main.go` to register all converted handlers.

### 7. Execute Incremental Testing

1. Convert one handler at a time
2. Run tests after each conversion
3. Verify API endpoints work correctly

## Migration Checklist

- [x] Remove Gin router
- [x] Rename `position_handler_chi.go` to `position_handler.go`
- [x] Convert `risk_handler.go` to use Chi
- [x] Convert `autobuy_handler.go` to use Chi
- [x] Convert `websocket.go` to use Chi
- [x] Update handler tests to use Chi
- [ ] Remove Gin dependency
- [ ] Verify all API endpoints work correctly 

# Refactoring Plan: Consolidation of Factory, Repository, and Middleware Components

## Overview

This document outlines a comprehensive plan to consolidate and standardize redundant components in the codebase, with a focus on:

1. Factory pattern implementations
2. Repository pattern implementations 
3. HTTP middleware components
4. Error handling standards
5. Migration strategy documentation

## 1. Factory Pattern Consolidation

### Current Issues
- Multiple factory implementations (MarketFactory, TradeFactory, WalletFactory, etc.)
- Redundant creation methods across factories
- Inconsistent initialization of dependencies
- Some factories create repositories that others also create

### Consolidation Plan

#### 1.1 Create a Unified AppFactory Structure
- Consolidate all factory functionality into a single `AppFactory` that manages all component creation
- Use a builder pattern for composable initialization

```go
// internal/factory/app_factory.go
package factory

import (
    // imports
)

// AppFactory is the single entry point for creating all application components
type AppFactory struct {
    config *config.Config
    logger *zerolog.Logger
    db     *gorm.DB
    
    // Core shared components
    mexcClient port.MEXCClient
    
    // Cached repositories
    repositories map[string]interface{}
    
    // Cached services
    services map[string]interface{}
}

// NewAppFactory creates a new unified app factory
func NewAppFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *AppFactory {
    return &AppFactory{
        config:       cfg,
        logger:       logger,
        db:           db,
        repositories: make(map[string]interface{}),
        services:     make(map[string]interface{}),
    }
}

// GetMEXCClient returns a shared MEXC client instance
func (f *AppFactory) GetMEXCClient() port.MEXCClient {
    if f.mexcClient == nil {
        f.mexcClient = mexc.NewClient(f.config.MEXC.APIKey, f.config.MEXC.APISecret, f.logger)
    }
    return f.mexcClient
}

// GetWalletRepository returns a wallet repository instance
func (f *AppFactory) GetWalletRepository() port.WalletRepository {
    repoKey := "wallet_repository"
    if repo, ok := f.repositories[repoKey]; ok {
        return repo.(port.WalletRepository)
    }
    
    repo := repo.NewConsolidatedWalletRepository(f.db, f.logger)
    f.repositories[repoKey] = repo
    return repo
}

// Additional methods for all other components
// ...
```

#### 1.2 Migrate All Factory Methods
- Systematically move all creation methods from existing factories to the unified AppFactory
- Maintain the same method naming convention for consistency
- Implement caching of instances to avoid duplicate creation

#### 1.3 Update Dependency Injection in main.go
- Refactor application bootstrapping to use the unified factory
- Simplify component initialization

## 2. Repository Pattern Standardization

### Current Issues
- Duplicate base repository implementations (`internal/adapter/persistence/gorm/repo/base_repository.go` and `internal/adapter/repository/gorm/base_repository.go`)
- Inconsistent repository methods across implementations
- Multiple repository factory methods

### Consolidation Plan

#### 2.1 Standardize on a Single BaseRepository
- Keep only the more comprehensive base repository implementation
- Add missing methods from the other implementation

```go
// internal/adapter/persistence/gorm/base_repository.go
package gorm

import (
    // imports
)

// BaseRepository provides common functionality for GORM repositories
type BaseRepository struct {
    db     *gorm.DB
    logger *zerolog.Logger
}

// Methods combined from both existing implementations
// ...
```

#### 2.2 Consolidate Repository Implementations
- Move all repositories to a consistent package structure
- Update all repositories to use the standardized BaseRepository
- Ensure consistent error handling and logging

## 3. HTTP Middleware Consolidation

### Current Issues
- Multiple authentication middleware implementations (ClerkMiddleware, EnhancedClerkMiddleware, SimpleAuthMiddleware, ConsolidatedAuthMiddleware)
- Inconsistent error handling across middleware
- No clear documentation on when to use which middleware
- Potential for test middleware to be enabled in production

### Consolidation Plan

#### 3.1 Standardize on ConsolidatedAuthMiddleware
- The `ConsolidatedAuthMiddleware` already exists and should be the single auth middleware
- Remove redundant middleware implementations:
  - ClerkMiddleware
  - EnhancedClerkMiddleware
  - SimpleAuthMiddleware

#### 3.2 Add Environment Flags for Test Middleware
- Add clear environment checks for test middleware
- Add warnings in logs when test middleware is activated
- Prevent test middleware from running in production

```go
// internal/adapter/http/middleware/test_auth_middleware.go

// NewTestAuthMiddleware creates a new test auth middleware
func NewTestAuthMiddleware(logger *zerolog.Logger, cfg *config.Config) (*TestAuthMiddleware, error) {
    // Check environment to prevent accidental usage in production
    if cfg.Environment == "production" {
        return nil, errors.New("test auth middleware cannot be used in production")
    }
    
    // Log warning about test middleware usage
    logger.Warn().Msg("TEST AUTH MIDDLEWARE ENABLED - NOT FOR PRODUCTION USE")
    
    return &TestAuthMiddleware{
        logger: logger,
    }, nil
}
```

#### 3.3 Document Middleware Usage Guidelines
Create a clear middleware documentation file that explains:
- When to use each middleware
- Configuration options
- Proper middleware ordering

## 4. Error Handling Standardization

### Current Issues
- Inconsistent error handling across the codebase
- Multiple ways to create and return errors
- No clear documentation on error handling best practices

### Consolidation Plan

#### 4.1 Standardize on AppError
- Use `apperror.AppError` as the standard error type throughout the application
- Ensure all errors returned to clients follow the standard format
- Document error codes and their meanings

#### 4.2 Ensure Consistent Error Middleware Usage
- Use the ErrorMiddleware consistently across all routes
- Standardize error response format

#### 4.3 Create Error Handling Guidelines
Document clear guidelines for error handling including:
- When to create new error types
- How to propagate errors
- How to log errors appropriately

## 5. Migration Strategy Documentation

### Current Issues
- Both SQL migrations and GORM AutoMigrate mentioned in the codebase
- Unclear which approach should be used

### Consolidation Plan

#### 5.1 Standardize on GORM AutoMigrate
- Document GORM's AutoMigrate as the official migration strategy
- Remove any SQL migration files or references to other migration tools
- Update the migration documentation to clearly state the chosen approach

#### 5.2 Create Migration Testing Guidelines
- Document how to test migrations in development and staging environments
- Create guidelines for handling breaking schema changes

## Implementation Timeline

1. **Phase 1: Documentation and Planning** (Week 1)
   - Create detailed documentation for the new consolidated patterns
   - Map all components that need to be migrated

2. **Phase 2: Base Components** (Week 2)
   - Implement the unified AppFactory
   - Standardize the BaseRepository
   - Create middleware usage guidelines

3. **Phase 3: Migration** (Weeks 3-4)
   - Systematically migrate components to the new patterns
   - Write tests for the consolidated components
   - Update application bootstrapping

4. **Phase 4: Validation and Cleanup** (Week 5)
   - Verify all functionality works with the new components
   - Remove deprecated components
   - Update all documentation

## Conclusion

This refactoring plan will significantly improve code maintainability by reducing duplication and standardizing patterns across the codebase. The unified factory pattern will make dependency injection more straightforward, while standardized repositories and middleware will ensure consistent behavior throughout the application. 