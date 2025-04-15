package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// CredentialValidationError represents a validation error
type CredentialValidationError struct {
	Field   string
	Message string
}

func (e CredentialValidationError) Error() string {
	return fmt.Sprintf("validation error for field %s: %s", e.Field, e.Message)
}

// CredentialValidationService handles validation of API credentials
type CredentialValidationService struct {
	credentialRepo port.APICredentialRepository
	logger         *zerolog.Logger
}

// NewCredentialValidationService creates a new CredentialValidationService
func NewCredentialValidationService(
	credentialRepo port.APICredentialRepository,
	logger *zerolog.Logger,
) *CredentialValidationService {
	return &CredentialValidationService{
		credentialRepo: credentialRepo,
		logger:         logger,
	}
}

// ValidateCredential validates an API credential
func (s *CredentialValidationService) ValidateCredential(ctx context.Context, credential *model.APICredential) error {
	// Validate required fields
	if credential.UserID == "" {
		return &CredentialValidationError{Field: "UserID", Message: "user ID is required"}
	}

	if credential.Exchange == "" {
		return &CredentialValidationError{Field: "Exchange", Message: "exchange is required"}
	}

	if credential.APIKey == "" {
		return &CredentialValidationError{Field: "APIKey", Message: "API key is required"}
	}

	if credential.APISecret == "" {
		return &CredentialValidationError{Field: "APISecret", Message: "API secret is required"}
	}

	// Validate exchange-specific formats
	if err := s.validateExchangeSpecificFormat(credential); err != nil {
		return err
	}

	// Validate label uniqueness
	if credential.Label != "" {
		if err := s.validateLabelUniqueness(ctx, credential); err != nil {
			return err
		}
	}

	// Validate expiration date
	if credential.ExpiresAt != nil && credential.ExpiresAt.Before(time.Now()) {
		return &CredentialValidationError{Field: "ExpiresAt", Message: "expiration date cannot be in the past"}
	}

	// Validate rotation due date
	if credential.RotationDue != nil && credential.RotationDue.Before(time.Now()) {
		return &CredentialValidationError{Field: "RotationDue", Message: "rotation due date cannot be in the past"}
	}

	return nil
}

// validateExchangeSpecificFormat validates exchange-specific formats for API credentials
func (s *CredentialValidationService) validateExchangeSpecificFormat(credential *model.APICredential) error {
	switch strings.ToLower(credential.Exchange) {
	case "mexc":
		return s.validateMEXCCredential(credential)
	case "binance":
		return s.validateBinanceCredential(credential)
	case "coinbase":
		return s.validateCoinbaseCredential(credential)
	case "kraken":
		return s.validateKrakenCredential(credential)
	default:
		// For unknown exchanges, just do basic validation
		return nil
	}
}

// validateMEXCCredential validates MEXC API credentials
func (s *CredentialValidationService) validateMEXCCredential(credential *model.APICredential) error {
	// MEXC API keys are typically 32 characters
	if len(credential.APIKey) < 16 || len(credential.APIKey) > 64 {
		return &CredentialValidationError{Field: "APIKey", Message: "MEXC API key should be between 16 and 64 characters"}
	}

	// MEXC API secrets are typically 32 characters
	if len(credential.APISecret) < 16 || len(credential.APISecret) > 64 {
		return &CredentialValidationError{Field: "APISecret", Message: "MEXC API secret should be between 16 and 64 characters"}
	}

	// MEXC API keys and secrets are typically alphanumeric
	apiKeyPattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !apiKeyPattern.MatchString(credential.APIKey) {
		return &CredentialValidationError{Field: "APIKey", Message: "MEXC API key should be alphanumeric"}
	}

	apiSecretPattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !apiSecretPattern.MatchString(credential.APISecret) {
		return &CredentialValidationError{Field: "APISecret", Message: "MEXC API secret should be alphanumeric"}
	}

	return nil
}

// validateBinanceCredential validates Binance API credentials
func (s *CredentialValidationService) validateBinanceCredential(credential *model.APICredential) error {
	// Binance API keys are typically 64 characters
	if len(credential.APIKey) < 16 || len(credential.APIKey) > 64 {
		return &CredentialValidationError{Field: "APIKey", Message: "Binance API key should be between 16 and 64 characters"}
	}

	// Binance API secrets are typically 64 characters
	if len(credential.APISecret) < 16 || len(credential.APISecret) > 64 {
		return &CredentialValidationError{Field: "APISecret", Message: "Binance API secret should be between 16 and 64 characters"}
	}

	// Binance API keys and secrets are typically alphanumeric
	apiKeyPattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !apiKeyPattern.MatchString(credential.APIKey) {
		return &CredentialValidationError{Field: "APIKey", Message: "Binance API key should be alphanumeric"}
	}

	apiSecretPattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !apiSecretPattern.MatchString(credential.APISecret) {
		return &CredentialValidationError{Field: "APISecret", Message: "Binance API secret should be alphanumeric"}
	}

	return nil
}

// validateCoinbaseCredential validates Coinbase API credentials
func (s *CredentialValidationService) validateCoinbaseCredential(credential *model.APICredential) error {
	// Coinbase API keys typically start with a specific prefix
	if !strings.HasPrefix(credential.APIKey, "cb") {
		return &CredentialValidationError{Field: "APIKey", Message: "Coinbase API key should start with 'cb'"}
	}

	// Coinbase API secrets are typically 64 characters
	if len(credential.APISecret) < 16 || len(credential.APISecret) > 128 {
		return &CredentialValidationError{Field: "APISecret", Message: "Coinbase API secret should be between 16 and 128 characters"}
	}

	return nil
}

// validateKrakenCredential validates Kraken API credentials
func (s *CredentialValidationService) validateKrakenCredential(credential *model.APICredential) error {
	// Kraken API keys typically start with a specific prefix
	if !strings.HasPrefix(credential.APIKey, "K-") {
		return &CredentialValidationError{Field: "APIKey", Message: "Kraken API key should start with 'K-'"}
	}

	// Kraken API secrets are typically base64 encoded
	apiSecretPattern := regexp.MustCompile("^[a-zA-Z0-9+/=]+$")
	if !apiSecretPattern.MatchString(credential.APISecret) {
		return &CredentialValidationError{Field: "APISecret", Message: "Kraken API secret should be base64 encoded"}
	}

	return nil
}

// validateLabelUniqueness validates that the label is unique for the user and exchange
func (s *CredentialValidationService) validateLabelUniqueness(ctx context.Context, credential *model.APICredential) error {
	// Skip validation if the label is empty
	if credential.Label == "" {
		return nil
	}

	// Check if a credential with the same label already exists
	existingCredential, err := s.credentialRepo.GetByUserIDAndLabel(ctx, credential.UserID, credential.Exchange, credential.Label)
	if err != nil {
		// If the error is not "credential not found", return the error
		if err != model.ErrCredentialNotFound {
			s.logger.Error().Err(err).Str("userID", credential.UserID).Str("exchange", credential.Exchange).Str("label", credential.Label).Msg("Failed to check label uniqueness")
			return fmt.Errorf("failed to check label uniqueness: %w", err)
		}
		// If the error is "credential not found", the label is unique
		return nil
	}

	// If the credential exists and it's not the same credential, return an error
	if existingCredential != nil && existingCredential.ID != credential.ID {
		return &CredentialValidationError{Field: "Label", Message: "label must be unique for the user and exchange"}
	}

	return nil
}
