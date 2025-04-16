# Error Handling Standardization Strategy

## Overview

This document outlines our standardized approach to error handling across the application. The goal is to establish consistent, informative, and secure error handling patterns that improve developer experience, ease debugging, and provide clear error information to clients.

## Core Error Types

The foundation of our error handling is the `AppError` struct defined in `internal/apperror/errors.go`:

```go
type AppError struct {
	StatusCode int                    // HTTP status code
	Code       string                 // Error code (e.g., NOT_FOUND, UNAUTHORIZED)
	Message    string                 // User-friendly error message
	Details    interface{}            // Optional additional details (validation errors, etc.)
	Err        error                  // Original error that can be unwrapped
}
```

We have predefined common error types:
- `ErrInvalidInput` - 400 Bad Request
- `ErrNotFound` - 404 Not Found
- `ErrInternal` - 500 Internal Server Error
- `ErrUnauthorized` - 401 Unauthorized
- `ErrForbidden` - 403 Forbidden
- `ErrConflict` - 409 Conflict
- `ErrRateLimit` - 429 Too Many Requests
- `ErrExternalService` - 503 Service Unavailable
- `ErrValidation` - 400 Bad Request (for validation errors)

## Error Response Structure

The standardized error response format is:

```json
{
	"error": {
		"code": "NOT_FOUND",
		"message": "User with ID 123 not found",
		"details": { ... },
		"trace_id": "abcd1234-ef56-...",
		"field_errors": {
			"email": "Invalid email format",
			"password": "Password must be at least 8 characters"
		}
	}
}
```

## Standardized Error Handling

The new standardized error handling system in `internal/apperror/standardized.go` adds:

1. **Error Context**: Error handlers can be passed through request context
2. **Tracing**: All errors include a request ID for correlation
3. **Field Validation**: Field-specific validation errors can be included
4. **Wrapping & Type Checking**: Utilities for error wrapping and type checking
5. **Consistent Responses**: All errors are returned in a consistent format

## Unified Error Middleware

The `UnifiedErrorMiddleware` in `internal/adapter/http/middleware/unified_error_middleware.go` provides:

- Request ID generation and propagation
- Panic recovery with detailed logging
- Request and response logging with correlation IDs
- Context-based error handling
- HTTP response capturing

## Usage Guidelines

### HTTP Handler Layer

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		// Use the context-based error handling
		apperror.RespondWithError(w, r, err)
		return
	}
	
	// Normal response
	response.JSON(w, http.StatusOK, user)
}
```

For validation:

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperror.RespondWithError(w, r, apperror.NewInvalid("Invalid JSON", nil, err))
		return
	}
	
	// Validate input
	if input.Email == "" {
		fieldErrors := map[string]string{"email": "Email is required"}
		traceID := apperror.GetTraceID(r)
		apperror.WriteValidationErrors(w, fieldErrors, traceID)
		return
	}
	
	// Process valid input...
}
```

### Service Layer

```go
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		// Wrap repository errors with context
		return nil, apperror.WrapError(err, "Failed to get user")
	}
	
	return user, nil
}
```

### Repository Layer

```go
func (r *UserRepository) GetUser(ctx context.Context, id string) (*User, error) {
	var user User
	result := r.db.First(&user, "id = ?", id)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Return a standard not found error
			return nil, apperror.NewNotFound("user", id, result.Error)
		}
		// Return standard internal error
		return nil, apperror.NewInternal(result.Error)
	}
	
	return &user, nil
}
```

### External API Calls

```go
func (c *ExternalAPIClient) FetchData(ctx context.Context) ([]byte, error) {
	resp, err := c.httpClient.Get("https://api.example.com/data")
	if err != nil {
		return nil, apperror.NewExternalService("example-api", "Failed to connect", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return nil, apperror.NewExternalService(
			"example-api",
			fmt.Sprintf("API returned status %d", resp.StatusCode),
			fmt.Errorf("API error: %d", resp.StatusCode),
		)
	}
	
	// Process successful response...
}
```

## Logging Best Practices

1. **Include Request IDs**: Always include the request ID in log entries
2. **Log at Appropriate Levels**: Use Error for actual errors, Info for normal operations, Debug for details
3. **Structured Logging**: Use the structured logging provided by zerolog
4. **Don't Log Sensitive Data**: Avoid logging passwords, tokens, etc.

```go
// Example of good error logging
logger.Error().
	Str("request_id", requestID).
	Str("user_id", userID).
	Err(err).
	Str("operation", "update_user").
	Msg("Failed to update user profile")
```

## Error Handling Setup

To use the standardized error handling, update your application setup:

```go
// In your main.go or server setup
func setupMiddleware(router chi.Router) {
	// Create the unified error middleware
	errorMiddleware := middleware.NewUnifiedErrorMiddleware(logger)
	
	// Add it early in the middleware chain
	router.Use(errorMiddleware.Middleware())
	
	// Other middleware...
}
```

## Type Checking Errors

Use the provided helper functions to check error types:

```go
if apperror.IsNotFound(err) {
	// Handle not found case
}

if apperror.IsUnauthorized(err) {
	// Handle unauthorized case
}

// Get HTTP status code from any error
statusCode := apperror.GetStatusCode(err)
```

## Migration Plan

1. **Phase 1**: Implement the standardized error handling files and middleware
2. **Phase 2**: Update HTTP handlers to use the new error handling
3. **Phase 3**: Update services and repositories to use consistent error types
4. **Phase 4**: Update clients and tests to expect the new error format

## Conclusion

This standardized error handling approach provides several benefits:

1. **Consistency**: All errors follow the same structure
2. **Clarity**: Error messages are clear and context-rich
3. **Security**: Error details can be controlled to avoid leaking sensitive information
4. **Maintainability**: Error handling code is centralized and reusable
5. **Observability**: Tracing IDs and structured logs aid debugging
6. **Client Experience**: Errors are easy to understand and handle in clients

By adopting these patterns throughout the codebase, we create a more resilient and maintainable application. 