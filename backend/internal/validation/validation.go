package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a validation error with details
type ValidationError struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string, details interface{}) *ValidationError {
	return &ValidationError{
		Message: message,
		Details: details,
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct using validator tags
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return FormatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

// FormatValidationErrors formats validation errors into user-friendly messages
func FormatValidationErrors(errs validator.ValidationErrors) error {
	var errMsgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid email address", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be at least %s", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must not exceed %s", err.Field(), err.Param()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid URL", err.Field()))
		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s failed validation: %s", err.Field(), err.Tag()))
		}
	}
	return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, "; "))
}

// ValidatePositiveFloat validates that a float64 value is positive
func ValidatePositiveFloat(value float64) error {
	if value <= 0 {
		return fmt.Errorf("value must be positive, got %v", value)
	}
	return nil
}

// ValidatePercentage validates that a float64 value is between 0 and 100
func ValidatePercentage(value float64) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("percentage must be between 0 and 100, got %v", value)
	}
	return nil
}

// ValidateSymbol validates that a trading symbol is not empty
func ValidateSymbol(symbol string) error {
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	return nil
}

// ValidateOrderType validates that the order type is either 'market' or 'limit'
func ValidateOrderType(orderType string) error {
	orderType = strings.ToLower(orderType)
	if orderType != "market" && orderType != "limit" {
		return fmt.Errorf("invalid order type: %s (must be 'market' or 'limit')", orderType)
	}
	return nil
}

// ValidateTradeRequest validates a trade request based on the provided parameters
func ValidateTradeRequest(symbol string, amount float64, orderType string, price *float64) error {
	if err := ValidateSymbol(symbol); err != nil {
		return err
	}

	if err := ValidatePositiveFloat(amount); err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if err := ValidateOrderType(orderType); err != nil {
		return err
	}

	orderType = strings.ToLower(orderType)
	if orderType == "limit" && price == nil {
		return fmt.Errorf("price is required for limit orders")
	}

	if price != nil {
		if err := ValidatePositiveFloat(*price); err != nil {
			return fmt.Errorf("invalid price: %w", err)
		}
	}

	return nil
}

// ValidateRiskParameters validates risk management parameters
func ValidateRiskParameters(maxDrawdown, riskPerTrade, maxExposure, dailyLossLimit, minBalance float64) error {
	// Validate percentages
	for _, value := range []struct {
		name  string
		value float64
	}{
		{"maxDrawdown", maxDrawdown},
		{"riskPerTrade", riskPerTrade},
		{"maxExposure", maxExposure},
		{"dailyLossLimit", dailyLossLimit},
	} {
		if err := ValidatePercentage(value.value); err != nil {
			return fmt.Errorf("invalid %s: %w", value.name, err)
		}
	}

	// Validate minBalance
	if err := ValidatePositiveFloat(minBalance); err != nil {
		return fmt.Errorf("invalid minBalance: %w", err)
	}

	return nil
}
