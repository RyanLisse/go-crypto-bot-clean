package handler

import "github.com/neo/crypto-bot/internal/adapter/http/response"

// Mock error codes for testing
var (
	ErrorCodeBadRequest    = response.ErrorCodeBadRequest
	ErrorCodeNotFound      = response.ErrorCodeNotFound
	ErrorCodeInternalError = response.ErrorCodeInternalError
)
