package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// ContentValidationRule represents a rule for validating content
type ContentValidationRule struct {
	Name        string
	Description string
	Pattern     *regexp.Regexp
	Severity    string // "ERROR", "WARNING", "INFO"
}

// ContentValidator validates content for security and compliance
type ContentValidator struct {
	rules  []ContentValidationRule
	logger *zap.Logger
}

// ValidationResult represents the result of content validation
type ValidationResult struct {
	Valid       bool
	Violations  []ContentViolation
	SanitizedContent string
}

// ContentViolation represents a content validation violation
type ContentViolation struct {
	Rule        string
	Description string
	Severity    string
	Match       string
}

// NewContentValidator creates a new content validator
func NewContentValidator(logger *zap.Logger) *ContentValidator {
	validator := &ContentValidator{
		logger: logger,
	}

	// Add default rules
	validator.AddRule("SENSITIVE_DATA_CREDIT_CARD", "Credit card number", regexp.MustCompile(`\b(?:\d{4}[-\s]?){3}\d{4}\b`), "ERROR")
	validator.AddRule("SENSITIVE_DATA_SSN", "Social security number", regexp.MustCompile(`\b\d{3}[-\s]?\d{2}[-\s]?\d{4}\b`), "ERROR")
	validator.AddRule("SENSITIVE_DATA_API_KEY", "API key", regexp.MustCompile(`(?i)\b(api[-_]?key|apikey|access[-_]?key|auth[-_]?token|client[-_]?secret)[-_]?[:=]\s*["']?([a-zA-Z0-9]{16,})["']?`), "ERROR")
	validator.AddRule("SENSITIVE_DATA_PASSWORD", "Password", regexp.MustCompile(`(?i)"(password|passwd|pwd)":\s*"[^"]*"`), "ERROR")
	validator.AddRule("MALICIOUS_CODE_JS", "JavaScript code", regexp.MustCompile(`<script[\s\S]*?</script>`), "ERROR")
	validator.AddRule("MALICIOUS_CODE_HTML", "HTML tags", regexp.MustCompile(`<(?!code|pre|br|p|b|i|strong|em)([a-z][a-z0-9]*)\b[^>]*>(.*?)</\1>`), "WARNING")
	validator.AddRule("TRADING_ADVICE_DISCLAIMER", "Trading advice without disclaimer", regexp.MustCompile(`(?i)\b(buy|sell|invest|trade)\b`), "WARNING")
	validator.AddRule("FINANCIAL_ADVICE", "Financial advice", regexp.MustCompile(`(?i)\b(guarantee|guaranteed|promise|assured|certain)\b.*\b(return|profit|gain|income)\b`), "ERROR")
	validator.AddRule("MARKET_MANIPULATION", "Market manipulation", regexp.MustCompile(`(?i)\b(pump|dump|manipulate|scheme|scam)\b`), "ERROR")

	return validator
}

// AddRule adds a rule to the validator
func (v *ContentValidator) AddRule(name, description string, pattern *regexp.Regexp, severity string) {
	v.rules = append(v.rules, ContentValidationRule{
		Name:        name,
		Description: description,
		Pattern:     pattern,
		Severity:    severity,
	})
}

// Validate validates content against all rules
func (v *ContentValidator) Validate(ctx context.Context, content string) ValidationResult {
	result := ValidationResult{
		Valid:            true,
		Violations:       []ContentViolation{},
		SanitizedContent: content,
	}

	// Check each rule
	for _, rule := range v.rules {
		matches := rule.Pattern.FindAllString(content, -1)
		if len(matches) > 0 {
			// For ERROR severity, mark as invalid
			if rule.Severity == "ERROR" {
				result.Valid = false
			}

			// Add violation for each match
			for _, match := range matches {
				result.Violations = append(result.Violations, ContentViolation{
					Rule:        rule.Name,
					Description: rule.Description,
					Severity:    rule.Severity,
					Match:       match,
				})

				// Sanitize content
				if rule.Severity == "ERROR" {
					result.SanitizedContent = rule.Pattern.ReplaceAllString(result.SanitizedContent, "[REDACTED]")
				}
			}
		}
	}

	// Add disclaimer for trading advice if needed
	if containsTradingAdvice(content) && !containsDisclaimer(content) {
		disclaimer := "\n\nDISCLAIMER: This information is for educational purposes only and not financial advice. Trading cryptocurrencies involves significant risk. Always do your own research before making investment decisions."
		result.SanitizedContent += disclaimer
	}

	return result
}

// containsTradingAdvice checks if content contains trading advice
func containsTradingAdvice(content string) bool {
	tradingAdvicePattern := regexp.MustCompile(`(?i)\b(buy|sell|invest|trade|recommendation|suggest|advise)\b`)
	return tradingAdvicePattern.MatchString(content)
}

// containsDisclaimer checks if content contains a disclaimer
func containsDisclaimer(content string) bool {
	disclaimerPattern := regexp.MustCompile(`(?i)\b(disclaimer|not financial advice|not investment advice|educational purposes only|do your own research|dyor)\b`)
	return disclaimerPattern.MatchString(content)
}

// ValidateAndSanitize validates and sanitizes content
func (v *ContentValidator) ValidateAndSanitize(ctx context.Context, content string) (string, error) {
	result := v.Validate(ctx, content)

	// Log violations
	for _, violation := range result.Violations {
		if violation.Severity == "ERROR" {
			v.logger.Warn("Content validation error",
				zap.String("rule", violation.Rule),
				zap.String("description", violation.Description),
				zap.String("match", violation.Match),
			)
		} else {
			v.logger.Info("Content validation warning",
				zap.String("rule", violation.Rule),
				zap.String("description", violation.Description),
				zap.String("match", violation.Match),
			)
		}
	}

	// Return error if content is invalid
	if !result.Valid {
		return result.SanitizedContent, fmt.Errorf("content validation failed: %d violations", len(result.Violations))
	}

	return result.SanitizedContent, nil
}
