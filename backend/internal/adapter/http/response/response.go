package response

// APIResponse is a standardized response format for all API endpoints
type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *APIError   `json:"error,omitempty"`
}

// APIError represents an error in the API response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success creates a successful API response
func Success(data interface{}) APIResponse {
	return APIResponse{
		Status: "success",
		Data:   data,
	}
}

// Error creates an error API response
func Error(code string, message string) APIResponse {
	return APIResponse{
		Status: "error",
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

// ErrorCodes defines standard error codes for the API
const (
	ErrorCodeBadRequest       = "bad_request"
	ErrorCodeNotFound         = "not_found"
	ErrorCodeInternalError    = "internal_error"
	ErrorCodeRateLimitExceeded = "rate_limit_exceeded"
	ErrorCodeUnauthorized     = "unauthorized"
	ErrorCodeForbidden        = "forbidden"
)
