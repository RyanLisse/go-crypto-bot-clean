package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// CredentialErrorType represents the type of credential error
type CredentialErrorType string

const (
	// CredentialErrorTypeValidation represents a validation error
	CredentialErrorTypeValidation CredentialErrorType = "validation"

	// CredentialErrorTypeEncryption represents an encryption error
	CredentialErrorTypeEncryption CredentialErrorType = "encryption"

	// CredentialErrorTypeDecryption represents a decryption error
	CredentialErrorTypeDecryption CredentialErrorType = "decryption"

	// CredentialErrorTypeDatabase represents a database error
	CredentialErrorTypeDatabase CredentialErrorType = "database"

	// CredentialErrorTypeNotFound represents a not found error
	CredentialErrorTypeNotFound CredentialErrorType = "not_found"

	// CredentialErrorTypePermission represents a permission error
	CredentialErrorTypePermission CredentialErrorType = "permission"

	// CredentialErrorTypeExpired represents an expired credential error
	CredentialErrorTypeExpired CredentialErrorType = "expired"

	// CredentialErrorTypeRevoked represents a revoked credential error
	CredentialErrorTypeRevoked CredentialErrorType = "revoked"

	// CredentialErrorTypeInactive represents an inactive credential error
	CredentialErrorTypeInactive CredentialErrorType = "inactive"

	// CredentialErrorTypeUnknown represents an unknown error
	CredentialErrorTypeUnknown CredentialErrorType = "unknown"
)

// CredentialError represents an error related to API credentials
type CredentialError struct {
	Type      CredentialErrorType
	Message   string
	CredID    string
	UserID    string
	Exchange  string
	Timestamp time.Time
	Cause     error
}

// Error returns the error message
func (e *CredentialError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s error: %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *CredentialError) Unwrap() error {
	return e.Cause
}

// CredentialErrorService handles errors related to API credentials
type CredentialErrorService struct {
	credentialRepo port.APICredentialRepository
	logger         *zerolog.Logger
}

// NewCredentialErrorService creates a new CredentialErrorService
func NewCredentialErrorService(
	credentialRepo port.APICredentialRepository,
	logger *zerolog.Logger,
) *CredentialErrorService {
	return &CredentialErrorService{
		credentialRepo: credentialRepo,
		logger:         logger,
	}
}

// HandleError handles an error related to API credentials
func (s *CredentialErrorService) HandleError(ctx context.Context, err error, credID, userID, exchange string) error {
	// Create a new credential error
	credErr := &CredentialError{
		Type:      CredentialErrorTypeUnknown,
		Message:   "An unknown error occurred",
		CredID:    credID,
		UserID:    userID,
		Exchange:  exchange,
		Timestamp: time.Now(),
		Cause:     err,
	}

	// Determine the error type
	var validationErr *CredentialValidationError
	if errors.As(err, &validationErr) {
		credErr.Type = CredentialErrorTypeValidation
		credErr.Message = fmt.Sprintf("Validation error for field %s: %s", validationErr.Field, validationErr.Message)
	} else if errors.Is(err, model.ErrCredentialNotFound) {
		credErr.Type = CredentialErrorTypeNotFound
		credErr.Message = "Credential not found"
	} else if errors.Is(err, model.ErrInvalidUserID) {
		credErr.Type = CredentialErrorTypeValidation
		credErr.Message = "Invalid user ID"
	} else if errors.Is(err, model.ErrInvalidExchange) {
		credErr.Type = CredentialErrorTypeValidation
		credErr.Message = "Invalid exchange"
	} else if errors.Is(err, model.ErrInvalidAPIKey) {
		credErr.Type = CredentialErrorTypeValidation
		credErr.Message = "Invalid API key"
	} else if errors.Is(err, model.ErrInvalidAPISecret) {
		credErr.Type = CredentialErrorTypeValidation
		credErr.Message = "Invalid API secret"
	} else {
		// Try to determine the error type from the error message
		errMsg := err.Error()
		switch {
		case contains(errMsg, "encrypt", "encryption", "encrypting"):
			credErr.Type = CredentialErrorTypeEncryption
			credErr.Message = "Failed to encrypt credential"
		case contains(errMsg, "decrypt", "decryption", "decrypting"):
			credErr.Type = CredentialErrorTypeDecryption
			credErr.Message = "Failed to decrypt credential"
		case contains(errMsg, "database", "db", "sql", "query", "transaction"):
			credErr.Type = CredentialErrorTypeDatabase
			credErr.Message = "Database error"
		case contains(errMsg, "permission", "access", "unauthorized", "forbidden"):
			credErr.Type = CredentialErrorTypePermission
			credErr.Message = "Permission denied"
		case contains(errMsg, "expired", "expiration"):
			credErr.Type = CredentialErrorTypeExpired
			credErr.Message = "Credential expired"
		case contains(errMsg, "revoked", "revocation"):
			credErr.Type = CredentialErrorTypeRevoked
			credErr.Message = "Credential revoked"
		case contains(errMsg, "inactive", "disabled"):
			credErr.Type = CredentialErrorTypeInactive
			credErr.Message = "Credential inactive"
		}
	}

	// Log the error
	s.logError(credErr)

	// Update credential status if necessary
	if credID != "" {
		s.updateCredentialStatus(ctx, credErr)
	}

	return credErr
}

// logError logs a credential error
func (s *CredentialErrorService) logError(err *CredentialError) {
	// Create a logger event with the error details
	event := s.logger.Error().
		Str("error_type", string(err.Type)).
		Str("message", err.Message).
		Time("timestamp", err.Timestamp)

	// Add credential details if available
	if err.CredID != "" {
		event = event.Str("credential_id", err.CredID)
	}
	if err.UserID != "" {
		event = event.Str("user_id", err.UserID)
	}
	if err.Exchange != "" {
		event = event.Str("exchange", err.Exchange)
	}

	// Add the underlying error if available
	if err.Cause != nil {
		event = event.Err(err.Cause)
	}

	// Log the error
	event.Msg("Credential error occurred")
}

// updateCredentialStatus updates the status of a credential based on the error
func (s *CredentialErrorService) updateCredentialStatus(ctx context.Context, err *CredentialError) {
	// Only update status for certain error types
	var status model.APICredentialStatus
	switch err.Type {
	case CredentialErrorTypeExpired:
		status = model.APICredentialStatusExpired
	case CredentialErrorTypeRevoked:
		status = model.APICredentialStatusRevoked
	case CredentialErrorTypeInactive:
		status = model.APICredentialStatusInactive
	case CredentialErrorTypeDecryption:
		// Increment failure count for decryption errors
		if incrementErr := s.credentialRepo.IncrementFailureCount(ctx, err.CredID); incrementErr != nil {
			s.logger.Error().Err(incrementErr).Str("credential_id", err.CredID).Msg("Failed to increment failure count")
		}

		// Get the credential to check the failure count
		credential, getErr := s.credentialRepo.GetByID(ctx, err.CredID)
		if getErr != nil {
			s.logger.Error().Err(getErr).Str("credential_id", err.CredID).Msg("Failed to get credential")
			return
		}

		// Update status to failed if failure count exceeds threshold
		if credential.FailureCount >= 5 {
			status = model.APICredentialStatusFailed
		} else {
			// Don't update status if failure count is below threshold
			return
		}
	default:
		// Don't update status for other error types
		return
	}

	// Update the credential status
	if updateErr := s.credentialRepo.UpdateStatus(ctx, err.CredID, status); updateErr != nil {
		s.logger.Error().Err(updateErr).Str("credential_id", err.CredID).Str("status", string(status)).Msg("Failed to update credential status")
	}
}

// contains checks if any of the substrings are contained in the string
func contains(s string, substrings ...string) bool {
	for _, substring := range substrings {
		if strings.Contains(s, substring) {
			return true
		}
	}
	return false
}
