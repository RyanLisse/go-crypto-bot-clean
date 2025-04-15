package validation

import (
	"regexp"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
)

// SupportedExchanges is a list of supported exchanges
var SupportedExchanges = []string{
	"mexc",
	"binance",
	"coinbase",
	"kraken",
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// CredentialValidator provides validation for API credentials
type CredentialValidator struct {
	errors []ValidationError
}

// NewCredentialValidator creates a new CredentialValidator
func NewCredentialValidator() *CredentialValidator {
	return &CredentialValidator{
		errors: make([]ValidationError, 0),
	}
}

// ValidateExchange validates the exchange field
func (v *CredentialValidator) ValidateExchange(exchange string) *CredentialValidator {
	if exchange == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   "exchange",
			Message: "Exchange is required",
		})
		return v
	}

	// Convert to lowercase for case-insensitive comparison
	exchangeLower := strings.ToLower(exchange)
	valid := false
	for _, supportedExchange := range SupportedExchanges {
		if exchangeLower == supportedExchange {
			valid = true
			break
		}
	}

	if !valid {
		v.errors = append(v.errors, ValidationError{
			Field:   "exchange",
			Message: "Unsupported exchange. Supported exchanges: " + strings.Join(SupportedExchanges, ", "),
		})
	}

	return v
}

// ValidateAPIKey validates the API key field
func (v *CredentialValidator) ValidateAPIKey(apiKey, exchange string) *CredentialValidator {
	if apiKey == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   "apiKey",
			Message: "API key is required",
		})
		return v
	}

	// Exchange-specific validation
	switch strings.ToLower(exchange) {
	case "mexc":
		// MEXC API keys are typically alphanumeric and between 16-64 characters
		if len(apiKey) < 16 || len(apiKey) > 64 {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiKey",
				Message: "MEXC API key should be between 16 and 64 characters",
			})
		}

		apiKeyPattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
		if !apiKeyPattern.MatchString(apiKey) {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiKey",
				Message: "MEXC API key should be alphanumeric",
			})
		}
	case "binance":
		// Binance API keys are typically alphanumeric and 64 characters
		if len(apiKey) != 64 {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiKey",
				Message: "Binance API key should be 64 characters",
			})
		}
	case "coinbase":
		// Coinbase API keys typically start with a specific prefix
		if !strings.HasPrefix(apiKey, "CB") {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiKey",
				Message: "Coinbase API key should start with 'CB'",
			})
		}
	case "kraken":
		// Kraken API keys typically start with a specific prefix
		if !strings.HasPrefix(apiKey, "K-") {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiKey",
				Message: "Kraken API key should start with 'K-'",
			})
		}
	}

	return v
}

// ValidateAPISecret validates the API secret field
func (v *CredentialValidator) ValidateAPISecret(apiSecret, exchange string) *CredentialValidator {
	if apiSecret == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   "apiSecret",
			Message: "API secret is required",
		})
		return v
	}

	// Exchange-specific validation
	switch strings.ToLower(exchange) {
	case "mexc":
		// MEXC API secrets are typically alphanumeric and between 16-64 characters
		if len(apiSecret) < 16 || len(apiSecret) > 64 {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiSecret",
				Message: "MEXC API secret should be between 16 and 64 characters",
			})
		}

		apiSecretPattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
		if !apiSecretPattern.MatchString(apiSecret) {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiSecret",
				Message: "MEXC API secret should be alphanumeric",
			})
		}
	case "kraken":
		// Kraken API secrets are typically base64 encoded
		apiSecretPattern := regexp.MustCompile("^[a-zA-Z0-9+/=]+$")
		if !apiSecretPattern.MatchString(apiSecret) {
			v.errors = append(v.errors, ValidationError{
				Field:   "apiSecret",
				Message: "Kraken API secret should be base64 encoded",
			})
		}
	}

	return v
}

// ValidateLabel validates the label field
func (v *CredentialValidator) ValidateLabel(label string) *CredentialValidator {
	// Label is optional, but if provided, it should not be too long
	if label != "" && len(label) > 50 {
		v.errors = append(v.errors, ValidationError{
			Field:   "label",
			Message: "Label should not exceed 50 characters",
		})
	}

	return v
}

// GetErrors returns all validation errors
func (v *CredentialValidator) GetErrors() []ValidationError {
	return v.errors
}

// HasErrors checks if there are any validation errors
func (v *CredentialValidator) HasErrors() bool {
	return len(v.errors) > 0
}

// ToAppError converts validation errors to an AppError
func (v *CredentialValidator) ToAppError() *apperror.AppError {
	if !v.HasErrors() {
		return nil
	}

	// Create a map of field errors
	fieldErrors := make(map[string]string)
	for _, err := range v.errors {
		fieldErrors[err.Field] = err.Message
	}

	return apperror.NewValidation("Validation error", fieldErrors, nil)
}
