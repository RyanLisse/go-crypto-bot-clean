# Error Handling in the API

This document describes the standardized error handling approach used in the API.

## Error Response Format

All API errors are returned in a consistent format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "trace_id": "unique-request-id",
    "details": { /* Optional additional error details */ }
  }
}
```

## Error Codes

The API uses the following standardized error codes:

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `BAD_REQUEST` | 400 | The request was invalid or cannot be served. |
| `UNAUTHORIZED` | 401 | Authentication is required and has failed or has not been provided. |
| `FORBIDDEN` | 403 | The request is understood, but it has been refused or access is not allowed. |
| `NOT_FOUND` | 404 | The requested resource does not exist. |
| `METHOD_NOT_ALLOWED` | 405 | The request method is not supported for the requested resource. |
| `CONFLICT` | 409 | The request conflicts with the current state of the server. |
| `VALIDATION_ERROR` | 422 | The request was well-formed but was unable to be followed due to semantic errors. |
| `RATE_LIMIT_EXCEEDED` | 429 | The user has sent too many requests in a given amount of time. |
| `INTERNAL_ERROR` | 500 | An unexpected condition was encountered on the server. |
| `SERVICE_UNAVAILABLE` | 503 | The server is currently unavailable (because it is overloaded or down for maintenance). |
| `BAD_GATEWAY` | 502 | The server received an invalid response from an upstream server. |
| `GATEWAY_TIMEOUT` | 504 | The server was acting as a gateway or proxy and did not receive a timely response from the upstream server. |
| `DATABASE_ERROR` | 500 | An error occurred while interacting with the database. |
| `EXTERNAL_SERVICE_ERROR` | 500 | An error occurred while interacting with an external service. |

## Error Details

Some errors may include additional details to help diagnose the issue. For example:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "trace_id": "abc-123-xyz",
    "details": {
      "fields": {
        "email": "Invalid email format",
        "password": "Password must be at least 8 characters"
      }
    }
  }
}
```

## Trace IDs

Each request is assigned a unique trace ID, which is included in the error response. This ID can be used to correlate the error with server logs for troubleshooting.

## Using Error Handling in Handlers

Handlers should use the `apperror` package to create and return errors:

```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
    id := chi.URLParam(r, "id")
    
    user, err := h.userService.GetUser(r.Context(), id)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return apperror.NewNotFound("User", id, err)
        }
        return apperror.NewInternal(err)
    }
    
    return apperror.RespondWithOK(w, r, user)
}
```

Then use the `HandleError` middleware to handle these errors:

```go
router.Get("/users/{id}", apperror.HandleError(handler.GetUser, logger))
```

## Panic Recovery

The `UnifiedErrorMiddleware` automatically recovers from panics and converts them to appropriate error responses. This ensures that the API remains stable even in the face of unexpected errors.
