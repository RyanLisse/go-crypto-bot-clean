package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go-crypto-bot-clean/backend/internal/domain/audit"
	"go-crypto-bot-clean/backend/internal/domain/security"

	"go.uber.org/zap"
)

// SecurityConfig contains configuration for AI security
type SecurityConfig struct {
	// EnableContentValidation enables content validation
	EnableContentValidation bool

	// EnableInputSanitization enables input sanitization
	EnableInputSanitization bool

	// EnableEncryption enables encryption of sensitive data
	EnableEncryption bool

	// Logger is the logger to use
	Logger *zap.Logger

	// AuditService is the audit service to use
	AuditService audit.Service

	// EncryptionService is the encryption service to use
	EncryptionService security.EncryptionService
}

// DefaultSecurityConfig returns the default security configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		EnableContentValidation: true,
		EnableInputSanitization: true,
		EnableEncryption:        true,
		Logger:                  zap.NewNop(),
	}
}

// AISecurityService provides security features for AI services
type AISecurityService struct {
	config          SecurityConfig
	contentValidator *ContentValidator
}

// NewAISecurityService creates a new AI security service
func NewAISecurityService(config SecurityConfig) *AISecurityService {
	return &AISecurityService{
		config:          config,
		contentValidator: NewContentValidator(config.Logger),
	}
}

// SanitizeInput sanitizes user input
func (s *AISecurityService) SanitizeInput(ctx context.Context, input string) (string, error) {
	if !s.config.EnableInputSanitization {
		return input, nil
	}

	// Remove potentially harmful patterns
	sanitized := input

	// Remove script tags
	scriptPattern := regexp.MustCompile(`<script[\s\S]*?</script>`)
	sanitized = scriptPattern.ReplaceAllString(sanitized, "")

	// Remove HTML tags
	htmlPattern := regexp.MustCompile(`<(?!code|pre|br|p|b|i|strong|em)([a-z][a-z0-9]*)\b[^>]*>(.*?)</\1>`)
	sanitized = htmlPattern.ReplaceAllString(sanitized, "$2")

	// Remove SQL injection patterns
	sqlPattern := regexp.MustCompile(`(?i)(union\s+select|select\s+.*\s+from|insert\s+into|update\s+.*\s+set|delete\s+from|drop\s+table|exec\s+xp_|exec\s+sp_|exec\s+master|declare\s+@|;--)`)
	sanitized = sqlPattern.ReplaceAllString(sanitized, "")

	// Log if input was sanitized
	if sanitized != input {
		s.config.Logger.Info("Input sanitized",
			zap.String("original_length", fmt.Sprintf("%d", len(input))),
			zap.String("sanitized_length", fmt.Sprintf("%d", len(sanitized))),
		)

		// Create audit event
		if s.config.AuditService != nil {
			event, err := audit.CreateAuditEvent(
				0, // User ID will be added by middleware
				audit.EventTypeSecurity,
				audit.EventSeverityWarning,
				"INPUT_SANITIZED",
				"Potentially harmful content removed from user input",
				map[string]interface{}{
					"original_length":  len(input),
					"sanitized_length": len(sanitized),
				},
				"", // IP will be added by middleware
				"", // User agent will be added by middleware
				"", // Request ID will be added by middleware
			)
			if err == nil {
				s.config.AuditService.LogEvent(ctx, event)
			}
		}
	}

	return sanitized, nil
}

// ValidateOutput validates and sanitizes AI output
func (s *AISecurityService) ValidateOutput(ctx context.Context, output string) (string, error) {
	if !s.config.EnableContentValidation {
		return output, nil
	}

	// Validate and sanitize content
	sanitized, err := s.contentValidator.ValidateAndSanitize(ctx, output)
	
	// Log if output was sanitized
	if sanitized != output {
		s.config.Logger.Info("Output sanitized",
			zap.String("original_length", fmt.Sprintf("%d", len(output))),
			zap.String("sanitized_length", fmt.Sprintf("%d", len(sanitized))),
			zap.Error(err),
		)

		// Create audit event
		if s.config.AuditService != nil {
			event, err := audit.CreateAuditEvent(
				0, // User ID will be added by middleware
				audit.EventTypeSecurity,
				audit.EventSeverityWarning,
				"OUTPUT_SANITIZED",
				"Potentially harmful content removed from AI output",
				map[string]interface{}{
					"original_length":  len(output),
					"sanitized_length": len(sanitized),
				},
				"", // IP will be added by middleware
				"", // User agent will be added by middleware
				"", // Request ID will be added by middleware
			)
			if err == nil {
				s.config.AuditService.LogEvent(ctx, event)
			}
		}
	}

	return sanitized, nil
}

// EncryptSensitiveData encrypts sensitive data
func (s *AISecurityService) EncryptSensitiveData(ctx context.Context, data string) (string, error) {
	if !s.config.EnableEncryption || s.config.EncryptionService == nil {
		return data, nil
	}

	// Encrypt data
	encrypted, err := s.config.EncryptionService.EncryptString(data)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	return encrypted, nil
}

// DecryptSensitiveData decrypts sensitive data
func (s *AISecurityService) DecryptSensitiveData(ctx context.Context, encryptedData string) (string, error) {
	if !s.config.EnableEncryption || s.config.EncryptionService == nil {
		return encryptedData, nil
	}

	// Decrypt data
	decrypted, err := s.config.EncryptionService.DecryptString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decrypted, nil
}
