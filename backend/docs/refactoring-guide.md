# Backend Refactoring Guide

This document provides an overview of the refactoring work completed on the backend codebase, including the rationale behind the changes, the benefits gained, and guidelines for future development.

## Table of Contents

1. [Overview](#overview)
2. [Data Flow Improvements](#data-flow-improvements)
3. [Redundancy Consolidation](#redundancy-consolidation)
4. [Authentication Simplification](#authentication-simplification)
5. [Secure Key Handling](#secure-key-handling)
6. [Migration Standardization](#migration-standardization)
7. [Implementation Fixes](#implementation-fixes)
8. [Error Handling](#error-handling)
9. [Logging](#logging)
10. [Best Practices](#best-practices)

## Overview

The backend refactoring project addressed several critical issues in the codebase:

1. **Data Flow**: Fixed inconsistent data flow between layers
2. **Redundancy**: Consolidated duplicate code and components
3. **Authentication**: Simplified and standardized authentication
4. **Key Handling**: Improved security of sensitive data
5. **Migrations**: Standardized database migration approach
6. **Implementation**: Fixed various implementation issues

## Data Flow Improvements

### Changes Made

- Updated MEXC client to use live API calls instead of sample data
- Removed direct API calls from HTTP handlers
- Consolidated data fetching logic in appropriate layers
- Added proper repository support for market data

### Benefits

- Clear separation of concerns between layers
- Consistent data flow from API to database
- Proper use of the repository pattern
- Improved error handling and logging

### Example

Before:
```go
// HTTP handler making direct API calls
func (h *Handler) GetMarketData(w http.ResponseWriter, r *http.Request) {
    client := mexc.NewClient(h.config.APIKey, h.config.APISecret)
    data, err := client.GetMarketData(r.Context(), "BTCUSDT")
    // Handle response...
}
```

After:
```go
// HTTP handler using service layer
func (h *Handler) GetMarketData(w http.ResponseWriter, r *http.Request) {
    data, err := h.marketService.GetMarketData(r.Context(), "BTCUSDT")
    // Handle response...
}

// Service using repository
func (s *MarketService) GetMarketData(ctx context.Context, symbol string) (*model.MarketData, error) {
    // Try to get from cache/database first
    data, err := s.repository.GetMarketData(ctx, symbol)
    if err == nil && data != nil {
        return data, nil
    }
    
    // Fetch from external API if not found
    data, err = s.client.GetMarketData(ctx, symbol)
    if err != nil {
        return nil, err
    }
    
    // Save to database for future use
    if err := s.repository.SaveMarketData(ctx, data); err != nil {
        s.logger.Warn().Err(err).Msg("Failed to save market data")
    }
    
    return data, nil
}
```

## Redundancy Consolidation

### Changes Made

- Created unified ConsolidatedFactory in factory package
- Consolidated redundant entity definitions
- Created consolidated repository implementations
- Removed redundant files and code

### Benefits

- Reduced code duplication
- Improved maintainability
- Consistent patterns across the codebase
- Easier to find and modify code

### Example

Before:
```go
// Multiple factory implementations
type RepositoryFactory struct { /* ... */ }
type ServiceFactory struct { /* ... */ }
type HandlerFactory struct { /* ... */ }

// Multiple entity definitions
type APICredentialEntity struct { /* ... */ }
type APICredential struct { /* ... */ }
```

After:
```go
// Unified factory
type ConsolidatedFactory struct { /* ... */ }

// Single entity definition
type APICredential struct { /* ... */ }
```

## Authentication Simplification

### Changes Made

- Standardized on Clerk as the primary authentication strategy
- Created consolidated authentication middleware
- Updated auth factory to use the consolidated middleware
- Fixed context key types for improved type safety

### Benefits

- Consistent authentication checks
- Type-safe context keys
- Proper token validation
- Simplified authentication flow

### Example

Before:
```go
// Multiple authentication middlewares
type ClerkMiddleware struct { /* ... */ }
type EnhancedClerkMiddleware struct { /* ... */ }
type SimpleAuthMiddleware struct { /* ... */ }

// String-based context keys
const UserIDKey = "user_id"
const RoleKey = "role"
```

After:
```go
// Single authentication middleware
type ConsolidatedAuthMiddleware struct { /* ... */ }

// Type-safe context keys
type UserIDKey struct{}
type RoleKey struct{}
```

## Secure Key Handling

### Changes Made

- Made encryption key mandatory in production
- Improved key validation with better error messages
- Added secure fallback for development environments
- Enhanced error handling for key-related issues

### Benefits

- Prevents accidental use of default keys in production
- Ensures proper key length and format
- Provides clear error messages for misconfiguration
- Allows development without requiring key configuration

### Example

Before:
```go
// Using a default key if not provided
keyB64 := os.Getenv("ENCRYPTION_KEY")
if keyB64 == "" {
    // For development, use a default key
    keyB64 = "Wn3PvhLOYk0QpFdod9qUDRRik9cI8jD3noi0TgrTJ1M="
}
```

After:
```go
// Requiring a key in production
keyB64 := os.Getenv("ENCRYPTION_KEY")
if keyB64 == "" {
    // Check if we're in production
    if os.Getenv("ENV") == "production" || os.Getenv("GO_ENV") == "production" {
        return nil, errors.New("ENCRYPTION_KEY environment variable is required in production")
    }
    
    // For non-production, log a warning and use a temporary key
    fmt.Fprintf(os.Stderr, "WARNING: Using a temporary encryption key. This is insecure and should only be used for development.\n")
    
    // Generate a temporary key
    tempKey, err := GenerateEncryptionKey()
    if err != nil {
        return nil, fmt.Errorf("failed to generate temporary encryption key: %w", err)
    }
    keyB64 = tempKey
}
```

## Migration Standardization

### Changes Made

- Standardized on GORM AutoMigrate for all database migrations
- Created unified migration system in auto_migrate.go
- Updated dedicated migration command
- Removed redundant migration methods

### Benefits

- Single, consistent approach to database migrations
- No need to maintain separate SQL migration files
- Reduced context switching between Go code and SQL
- Schema changes are directly tied to entity definitions

### Example

Before:
```go
// SQL-based migrations
func RunMigrations(db *sql.DB) error {
    migrations := []string{
        `CREATE TABLE IF NOT EXISTS users (...)`,
        `ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255)`,
        // More SQL statements...
    }
    
    for _, migration := range migrations {
        if _, err := db.Exec(migration); err != nil {
            return err
        }
    }
    
    return nil
}
```

After:
```go
// GORM AutoMigrate
func RunMigrations(db *gorm.DB, logger *zerolog.Logger) error {
    logger.Info().Msg("Starting database migrations")
    
    // List of all entity models to migrate
    models := []interface{}{
        &entity.User{},
        &entity.APICredential{},
        // More models...
    }
    
    // Run AutoMigrate on all models
    for _, model := range models {
        logger.Debug().Str("model", getModelName(model)).Msg("Migrating model")
        if err := db.AutoMigrate(model); err != nil {
            logger.Error().Err(err).Str("model", getModelName(model)).Msg("Failed to migrate model")
            return err
        }
    }
    
    logger.Info().Msg("All models migrated successfully")
    return nil
}
```

## Implementation Fixes

### Changes Made

- Fixed MarketDataRepository implementation
- Created database infrastructure package
- Updated migration command to use the new database package
- Fixed API credential repository implementation

### Benefits

- Improved code organization
- Fixed field naming inconsistencies
- Enhanced error handling and logging
- Simplified database connection management

## Error Handling

### Approach

We've implemented a centralized error handling system with the following components:

1. **AppError Type**: A structured error type that includes:
   - Error type (validation, not found, unauthorized, etc.)
   - HTTP status code
   - Error message
   - Optional details
   - Original error (for logging)

2. **Error Middleware**: A middleware that:
   - Captures panics and converts them to structured errors
   - Adds request IDs to all responses
   - Ensures consistent error responses

3. **Helper Functions**: Functions to create specific error types:
   - `NewValidation`: For validation errors
   - `NewNotFound`: For not found errors
   - `NewUnauthorized`: For authentication errors
   - And more...

### Usage

```go
// Creating errors
if user == nil {
    return apperror.NewNotFound("User", userID, nil)
}

if err := validate(input); err != nil {
    return apperror.NewValidation("Invalid input", err.Error(), err)
}

// Handling errors in HTTP handlers
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.userService.GetUser(r.Context(), userID)
    if err != nil {
        var appErr *apperror.AppError
        if errors.As(err, &appErr) {
            apperror.WriteError(w, appErr)
        } else {
            apperror.WriteError(w, apperror.NewInternal(err))
        }
        return
    }
    
    // Handle success...
}
```

## Logging

### Approach

We've implemented a structured logging system using zerolog with the following components:

1. **Logger Configuration**: Configured zerolog for structured JSON logging in production and human-readable logging in development.

2. **Logging Middleware**: A middleware that logs HTTP requests with:
   - Request method and path
   - Remote address and user agent
   - Request ID for correlation
   - Response status code
   - Request duration

3. **Contextual Logging**: Added context to logs throughout the application:
   - Component name (e.g., "http", "service", "repository")
   - User ID when available
   - Request-specific information

### Usage

```go
// Service-level logging
func (s *UserService) GetUser(ctx context.Context, userID string) (*model.User, error) {
    s.logger.Debug().Str("user_id", userID).Msg("Getting user")
    
    user, err := s.repository.GetUser(ctx, userID)
    if err != nil {
        s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get user")
        return nil, err
    }
    
    return user, nil
}
```

## Best Practices

### Repository Pattern

- Use the repository pattern for all database access
- Implement interfaces in the domain layer
- Keep repository methods focused on CRUD operations
- Use meaningful method names that reflect the operation

### Error Handling

- Use structured errors with appropriate HTTP status codes
- Include enough context in error messages
- Log errors with stack traces at the appropriate level
- Don't expose internal errors to clients

### Dependency Injection

- Use the factory pattern for creating components
- Inject dependencies through constructors
- Avoid global state and singletons
- Use interfaces for testability

### Configuration

- Use environment variables for configuration
- Provide sensible defaults for development
- Validate configuration at startup
- Document required configuration

### Testing

- Write unit tests for business logic
- Use interfaces and mocks for dependencies
- Test error cases as well as happy paths
- Use table-driven tests for multiple scenarios
