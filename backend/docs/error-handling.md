# Error Handling Strategy

This document defines the standardized approach to error handling throughout the application.

## Core Principles

1. **Consistency**: Use a uniform approach to error creation, propagation, and handling
2. **Context Preservation**: Include enough context to understand the error's origin and cause
3. **User Experience**: Provide appropriate error messages to end users without exposing internal details
4. **Traceability**: Ensure errors can be traced through logs for debugging
5. **Type Safety**: Use typed errors where appropriate for programmatic handling

## Error Types

### Domain Errors

Domain errors represent expected error conditions within the application's business logic. They should be defined in `internal/domain/error/errors.go`:

```go
package error

import "fmt"

// ErrorType represents the type of error
type ErrorType string

const (
    // Common error types
    NotFound           ErrorType = "NOT_FOUND"
    ValidationFailed   ErrorType = "VALIDATION_FAILED"
    Unauthorized       ErrorType = "UNAUTHORIZED"
    Forbidden          ErrorType = "FORBIDDEN"
    InternalError      ErrorType = "INTERNAL_ERROR"
    ExternalServiceError ErrorType = "EXTERNAL_SERVICE_ERROR"
    RateLimitExceeded  ErrorType = "RATE_LIMIT_EXCEEDED"
    
    // Domain-specific error types
    InsufficientFunds  ErrorType = "INSUFFICIENT_FUNDS"
    // Add other domain-specific error types here
)

// DomainError represents an error that occurs within the domain logic
type DomainError struct {
    Type        ErrorType
    Message     string
    InternalErr error
    Code        int
    Details     map[string]interface{}
}

// Error returns the error message
func (e *DomainError) Error() string {
    if e.InternalErr != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.InternalErr)
    }
    return e.Message
}

// Unwrap returns the wrapped error
func (e *DomainError) Unwrap() error {
    return e.InternalErr
}

// New creates a new DomainError
func New(errType ErrorType, message string) *DomainError {
    return &DomainError{
        Type:    errType,
        Message: message,
        Code:    errorTypeToCode(errType),
    }
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errType ErrorType, message string) *DomainError {
    return &DomainError{
        Type:        errType,
        Message:     message,
        InternalErr: err,
        Code:        errorTypeToCode(errType),
    }
}

// WithDetails adds context details to the error
func (e *DomainError) WithDetails(details map[string]interface{}) *DomainError {
    e.Details = details
    return e
}

// IsDomainError checks if an error is a DomainError of a specific type
func IsDomainError(err error, errType ErrorType) bool {
    var domainErr *DomainError
    if errors.As(err, &domainErr) {
        return domainErr.Type == errType
    }
    return false
}

// errorTypeToCode maps error types to HTTP status codes
func errorTypeToCode(errType ErrorType) int {
    switch errType {
    case NotFound:
        return 404
    case ValidationFailed:
        return 400
    case Unauthorized:
        return 401
    case Forbidden:
        return 403
    case RateLimitExceeded:
        return 429
    case ExternalServiceError:
        return 502
    default:
        return 500
    }
}
```

### Usage Example

```go
// In your service code
func (s *service) GetUser(id string) (*model.User, error) {
    if id == "" {
        return nil, error.New(error.ValidationFailed, "user ID cannot be empty")
    }

    user, err := s.repository.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, error.New(error.NotFound, "user not found")
        }
        return nil, error.Wrap(err, error.InternalError, "failed to fetch user")
    }

    return user, nil
}
```

## HTTP Error Handling Middleware

Define a centralized error handling middleware in `internal/adapter/http/middleware/error_handler.go`:

```go
package middleware

import (
    "encoding/json"
    "net/http"
    
    domainerror "github.com/yourusername/yourproject/internal/domain/error"
    "github.com/yourusername/yourproject/internal/domain/service"
    "go.uber.org/zap"
)

// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
    Error       string                 `json:"error"`
    ErrorType   string                 `json:"error_type,omitempty"`
    Details     map[string]interface{} `json:"details,omitempty"`
    RequestID   string                 `json:"request_id,omitempty"`
}

// ErrorHandlerMiddleware creates middleware that handles errors from handlers
func ErrorHandlerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Create a response writer that can capture the response
            rw := NewResponseWriter(w)
            
            // Call the next handler
            next.ServeHTTP(rw, r)
            
            // Check if we have a domainerror.DomainError in context
            if err := service.GetErrorFromContext(r.Context()); err != nil {
                handleError(rw, r, err, logger)
                return
            }
            
            // If we've already written a response, do nothing
            if rw.Written() {
                return
            }
        })
    }
}

// handleError processes the error and writes an appropriate response
func handleError(w http.ResponseWriter, r *http.Request, err error, logger *zap.Logger) {
    // Get request ID from context if available
    requestID := middleware.GetRequestIDFromContext(r.Context())
    
    // Log the error
    logger.Error("Request error",
        zap.String("request_id", requestID),
        zap.String("method", r.Method),
        zap.String("path", r.URL.Path),
        zap.Error(err),
    )
    
    // Determine the appropriate status code and response
    code := http.StatusInternalServerError
    response := ErrorResponse{
        Error:     "An unexpected error occurred",
        RequestID: requestID,
    }
    
    // Check if it's a domain error
    var domainErr *domainerror.DomainError
    if errors.As(err, &domainErr) {
        code = domainErr.Code
        response.Error = domainErr.Message
        response.ErrorType = string(domainErr.Type)
        response.Details = domainErr.Details
        
        // Don't expose internal error details to the client
        if domainErr.Type == domainerror.InternalError || domainErr.Type == domainerror.ExternalServiceError {
            response.Error = "An internal error occurred"
        }
    }
    
    // Write the response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(response)
}
```

## Error Propagation Guidelines

1. **Service Layer**:
   - Return domain errors using the error package
   - Don't log errors at the service level (except debug logs)
   - Add context to errors when propagating them up

2. **Repository Layer**:
   - Wrap database or external service errors with domain errors
   - Include enough context for debugging

3. **Handler Layer**:
   - Set errors in the request context for the error middleware to handle
   - Don't handle errors directly unless specific custom handling is required

Example for handlers:

```go
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    user, err := h.service.GetUser(id)
    if err != nil {
        service.SetErrorInContext(r, err)
        return
    }
    
    response.JSON(w, http.StatusOK, user)
}
```

## Logging Guidelines

1. **Use structured logging (zap or zerolog)**
2. **Log all errors at the entry points (HTTP, message consumers, etc.)**
3. **Include request IDs, user IDs, and other contextual information**
4. **Don't log sensitive data (passwords, tokens, etc.)**

Example:

```go
// In HTTP middleware
logger.Error("Request error",
    zap.String("request_id", requestID),
    zap.String("method", r.Method),
    zap.String("path", r.URL.Path),
    zap.String("user_id", userID),
    zap.Error(err),
)
```

## Error Context Helper

Create a helper package in `internal/domain/service/error_context.go`:

```go
package service

import (
    "context"
    "net/http"
)

type contextKey string

const errorContextKey contextKey = "error"

// SetErrorInContext saves an error in the request context
func SetErrorInContext(r *http.Request, err error) {
    *r = *r.WithContext(context.WithValue(r.Context(), errorContextKey, err))
}

// GetErrorFromContext retrieves an error from the context
func GetErrorFromContext(ctx context.Context) error {
    if err, ok := ctx.Value(errorContextKey).(error); ok {
        return err
    }
    return nil
}
```

## Testing Errors

1. **Verify error types**: Test that the correct error types are returned
2. **Test error propagation**: Ensure errors are properly wrapped and context is preserved
3. **Check error handling middleware**: Test that the middleware properly translates errors to HTTP responses

Example:

```go
func TestGetUserNotFound(t *testing.T) {
    // Setup
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mock.NewMockUserRepository(ctrl)
    mockRepo.EXPECT().FindByID("nonexistent").Return(nil, error.New(error.NotFound, "user not found"))
    
    service := NewUserService(mockRepo)
    
    // Act
    user, err := service.GetUser("nonexistent")
    
    // Assert
    assert.Nil(t, user)
    assert.True(t, error.IsDomainError(err, error.NotFound))
}
```

## Conclusion

This standardized approach to error handling ensures consistency, maintainability, and a good user experience. By following these guidelines, we create a system where errors are properly created, propagated, and handled throughout the application. 