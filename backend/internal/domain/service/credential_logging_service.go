package service

import (
	"context"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// CredentialLogLevel represents the log level for credential operations
type CredentialLogLevel string

const (
	// CredentialLogLevelDebug represents the debug log level
	CredentialLogLevelDebug CredentialLogLevel = "debug"

	// CredentialLogLevelInfo represents the info log level
	CredentialLogLevelInfo CredentialLogLevel = "info"

	// CredentialLogLevelWarn represents the warn log level
	CredentialLogLevelWarn CredentialLogLevel = "warn"

	// CredentialLogLevelError represents the error log level
	CredentialLogLevelError CredentialLogLevel = "error"
)

// CredentialOperationType represents the type of credential operation
type CredentialOperationType string

const (
	// CredentialOperationTypeCreate represents a create operation
	CredentialOperationTypeCreate CredentialOperationType = "create"

	// CredentialOperationTypeRead represents a read operation
	CredentialOperationTypeRead CredentialOperationType = "read"

	// CredentialOperationTypeUpdate represents an update operation
	CredentialOperationTypeUpdate CredentialOperationType = "update"

	// CredentialOperationTypeDelete represents a delete operation
	CredentialOperationTypeDelete CredentialOperationType = "delete"

	// CredentialOperationTypeVerify represents a verify operation
	CredentialOperationTypeVerify CredentialOperationType = "verify"

	// CredentialOperationTypeEncrypt represents an encrypt operation
	CredentialOperationTypeEncrypt CredentialOperationType = "encrypt"

	// CredentialOperationTypeDecrypt represents a decrypt operation
	CredentialOperationTypeDecrypt CredentialOperationType = "decrypt"

	// CredentialOperationTypeRotate represents a rotate operation
	CredentialOperationTypeRotate CredentialOperationType = "rotate"

	// CredentialOperationTypeExpire represents an expire operation
	CredentialOperationTypeExpire CredentialOperationType = "expire"

	// CredentialOperationTypeRevoke represents a revoke operation
	CredentialOperationTypeRevoke CredentialOperationType = "revoke"

	// CredentialOperationTypeActivate represents an activate operation
	CredentialOperationTypeActivate CredentialOperationType = "activate"

	// CredentialOperationTypeDeactivate represents a deactivate operation
	CredentialOperationTypeDeactivate CredentialOperationType = "deactivate"
)

// CredentialLogEntry represents a log entry for a credential operation
type CredentialLogEntry struct {
	Level     CredentialLogLevel
	Operation CredentialOperationType
	CredID    string
	UserID    string
	Exchange  string
	Status    model.APICredentialStatus
	Message   string
	Timestamp time.Time
	Duration  time.Duration
	Error     error
	RequestID string
	ClientIP  string
	UserAgent string
	Metadata  map[string]string
}

// CredentialLoggingService handles logging for API credentials
type CredentialLoggingService struct {
	logger *zerolog.Logger
}

// NewCredentialLoggingService creates a new CredentialLoggingService
func NewCredentialLoggingService(logger *zerolog.Logger) *CredentialLoggingService {
	return &CredentialLoggingService{
		logger: logger,
	}
}

// LogOperation logs a credential operation
func (s *CredentialLoggingService) LogOperation(ctx context.Context, entry *CredentialLogEntry) {
	// Create a logger event with the appropriate level
	var event *zerolog.Event
	switch entry.Level {
	case CredentialLogLevelDebug:
		event = s.logger.Debug()
	case CredentialLogLevelInfo:
		event = s.logger.Info()
	case CredentialLogLevelWarn:
		event = s.logger.Warn()
	case CredentialLogLevelError:
		event = s.logger.Error()
	default:
		event = s.logger.Info()
	}

	// Add common fields
	event = event.
		Str("operation", string(entry.Operation)).
		Time("timestamp", entry.Timestamp).
		Dur("duration", entry.Duration)

	// Add credential details if available
	if entry.CredID != "" {
		event = event.Str("credential_id", entry.CredID)
	}
	if entry.UserID != "" {
		event = event.Str("user_id", entry.UserID)
	}
	if entry.Exchange != "" {
		event = event.Str("exchange", entry.Exchange)
	}
	if entry.Status != "" {
		event = event.Str("status", string(entry.Status))
	}

	// Add request details if available
	if entry.RequestID != "" {
		event = event.Str("request_id", entry.RequestID)
	}
	if entry.ClientIP != "" {
		event = event.Str("client_ip", s.maskIP(entry.ClientIP))
	}
	if entry.UserAgent != "" {
		event = event.Str("user_agent", entry.UserAgent)
	}

	// Add metadata if available
	if entry.Metadata != nil {
		for key, value := range entry.Metadata {
			// Skip sensitive fields
			if s.isSensitiveField(key) {
				event = event.Str(key, s.maskSensitiveValue(value))
			} else {
				event = event.Str(key, value)
			}
		}
	}

	// Add error if available
	if entry.Error != nil {
		event = event.Err(entry.Error)
	}

	// Log the message
	event.Msg(entry.Message)
}

// LogCredentialCreate logs a credential create operation
func (s *CredentialLoggingService) LogCredentialCreate(ctx context.Context, credential *model.APICredential, duration time.Duration, err error) {
	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelInfo,
		Operation: CredentialOperationTypeCreate,
		CredID:    credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		Status:    credential.Status,
		Message:   "API credential created",
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to create API credential"
	}

	s.LogOperation(ctx, entry)
}

// LogCredentialRead logs a credential read operation
func (s *CredentialLoggingService) LogCredentialRead(ctx context.Context, credential *model.APICredential, duration time.Duration, err error) {
	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelDebug,
		Operation: CredentialOperationTypeRead,
		CredID:    credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		Status:    credential.Status,
		Message:   "API credential read",
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to read API credential"
	}

	s.LogOperation(ctx, entry)
}

// LogCredentialUpdate logs a credential update operation
func (s *CredentialLoggingService) LogCredentialUpdate(ctx context.Context, credential *model.APICredential, duration time.Duration, err error) {
	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelInfo,
		Operation: CredentialOperationTypeUpdate,
		CredID:    credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		Status:    credential.Status,
		Message:   "API credential updated",
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to update API credential"
	}

	s.LogOperation(ctx, entry)
}

// LogCredentialDelete logs a credential delete operation
func (s *CredentialLoggingService) LogCredentialDelete(ctx context.Context, credID, userID, exchange string, duration time.Duration, err error) {
	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelInfo,
		Operation: CredentialOperationTypeDelete,
		CredID:    credID,
		UserID:    userID,
		Exchange:  exchange,
		Message:   "API credential deleted",
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to delete API credential"
	}

	s.LogOperation(ctx, entry)
}

// LogCredentialVerify logs a credential verify operation
func (s *CredentialLoggingService) LogCredentialVerify(ctx context.Context, credential *model.APICredential, duration time.Duration, err error) {
	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelInfo,
		Operation: CredentialOperationTypeVerify,
		CredID:    credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		Status:    credential.Status,
		Message:   "API credential verified",
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to verify API credential"
	}

	s.LogOperation(ctx, entry)
}

// LogCredentialEncrypt logs a credential encrypt operation
func (s *CredentialLoggingService) LogCredentialEncrypt(ctx context.Context, credential *model.APICredential, duration time.Duration, err error) {
	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelDebug,
		Operation: CredentialOperationTypeEncrypt,
		CredID:    credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		Status:    credential.Status,
		Message:   "API credential encrypted",
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to encrypt API credential"
	}

	s.LogOperation(ctx, entry)
}

// LogCredentialDecrypt logs a credential decrypt operation
func (s *CredentialLoggingService) LogCredentialDecrypt(ctx context.Context, credential *model.APICredential, duration time.Duration, err error) {
	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelDebug,
		Operation: CredentialOperationTypeDecrypt,
		CredID:    credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		Status:    credential.Status,
		Message:   "API credential decrypted",
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to decrypt API credential"
	}

	s.LogOperation(ctx, entry)
}

// LogCredentialStatusChange logs a credential status change operation
func (s *CredentialLoggingService) LogCredentialStatusChange(ctx context.Context, credential *model.APICredential, oldStatus model.APICredentialStatus, duration time.Duration, err error) {
	// Determine the operation type based on the new status
	var operation CredentialOperationType
	switch credential.Status {
	case model.APICredentialStatusActive:
		operation = CredentialOperationTypeActivate
	case model.APICredentialStatusInactive:
		operation = CredentialOperationTypeDeactivate
	case model.APICredentialStatusRevoked:
		operation = CredentialOperationTypeRevoke
	case model.APICredentialStatusExpired:
		operation = CredentialOperationTypeExpire
	default:
		operation = CredentialOperationTypeUpdate
	}

	entry := &CredentialLogEntry{
		Level:     CredentialLogLevelInfo,
		Operation: operation,
		CredID:    credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		Status:    credential.Status,
		Message:   "API credential status changed from " + string(oldStatus) + " to " + string(credential.Status),
		Timestamp: time.Now(),
		Duration:  duration,
		Error:     err,
		Metadata: map[string]string{
			"old_status": string(oldStatus),
			"new_status": string(credential.Status),
		},
	}

	// If there was an error, change the level to error
	if err != nil {
		entry.Level = CredentialLogLevelError
		entry.Message = "Failed to change API credential status"
	}

	s.LogOperation(ctx, entry)
}

// isSensitiveField checks if a field is sensitive
func (s *CredentialLoggingService) isSensitiveField(field string) bool {
	sensitiveFields := []string{
		"api_key",
		"apikey",
		"api_secret",
		"apisecret",
		"secret",
		"password",
		"token",
		"access_token",
		"refresh_token",
		"private_key",
		"privatekey",
		"passphrase",
	}

	for _, sensitiveField := range sensitiveFields {
		if field == sensitiveField {
			return true
		}
	}

	return false
}

// maskSensitiveValue masks a sensitive value
func (s *CredentialLoggingService) maskSensitiveValue(value string) string {
	if len(value) <= 8 {
		return "********"
	}

	// Show first 4 and last 4 characters
	return value[:4] + "********" + value[len(value)-4:]
}

// maskIP masks an IP address
func (s *CredentialLoggingService) maskIP(ip string) string {
	// For IPv4, mask the last octet
	// For IPv6, mask the last 4 segments
	// For simplicity, we'll just mask the last part of the IP
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		// IPv4
		return parts[0] + "." + parts[1] + "." + parts[2] + ".xxx"
	}

	parts = strings.Split(ip, ":")
	if len(parts) > 4 {
		// IPv6
		return strings.Join(parts[:4], ":") + ":xxxx:xxxx:xxxx:xxxx"
	}

	// Unknown format, mask the whole IP
	return "xxx.xxx.xxx.xxx"
}
