package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"go.uber.org/zap"
)

// EventType represents the type of audit event
type EventType string

const (
	// EventTypeAuth represents authentication events
	EventTypeAuth EventType = "AUTH"

	// EventTypeAI represents AI interaction events
	EventTypeAI EventType = "AI"

	// EventTypeAdmin represents administrative events
	EventTypeAdmin EventType = "ADMIN"

	// EventTypeTrading represents trading events
	EventTypeTrading EventType = "TRADING"

	// EventTypeSecurity represents security events
	EventTypeSecurity EventType = "SECURITY"
)

// EventSeverity represents the severity of an audit event
type EventSeverity string

const (
	// EventSeverityInfo represents informational events
	EventSeverityInfo EventSeverity = "INFO"

	// EventSeverityWarning represents warning events
	EventSeverityWarning EventSeverity = "WARNING"

	// EventSeverityError represents error events
	EventSeverityError EventSeverity = "ERROR"

	// EventSeverityCritical represents critical events
	EventSeverityCritical EventSeverity = "CRITICAL"
)

// AuditEvent represents an audit event
type AuditEvent struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      int            `json:"user_id" gorm:"index"`
	Type        EventType      `json:"type" gorm:"index"`
	Severity    EventSeverity  `json:"severity" gorm:"index"`
	Action      string         `json:"action" gorm:"index"`
	Description string         `json:"description"`
	Metadata    string         `json:"metadata" gorm:"type:text"`
	IP          string         `json:"ip"`
	UserAgent   string         `json:"user_agent"`
	RequestID   string         `json:"request_id" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at" gorm:"index"`
}

// TableName returns the table name for the audit event
func (AuditEvent) TableName() string {
	return "audit_events"
}

// Service defines the interface for the audit service
type Service interface {
	// LogEvent logs an audit event
	LogEvent(ctx context.Context, event *AuditEvent) error

	// GetEvents retrieves audit events with optional filtering
	GetEvents(ctx context.Context, userID int, eventType EventType, severity EventSeverity, startTime, endTime time.Time, limit, offset int) ([]*AuditEvent, error)

	// GetEventByID retrieves an audit event by ID
	GetEventByID(ctx context.Context, id uint) (*AuditEvent, error)

	// GetEventsByRequestID retrieves audit events by request ID
	GetEventsByRequestID(ctx context.Context, requestID string) ([]*AuditEvent, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewService creates a new audit service
func NewService(db *gorm.DB, logger *zap.Logger) (Service, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&AuditEvent{}); err != nil {
		return nil, fmt.Errorf("failed to migrate audit_events table: %w", err)
	}

	return &serviceImpl{
		db:     db,
		logger: logger,
	}, nil
}

// LogEvent logs an audit event
func (s *serviceImpl) LogEvent(ctx context.Context, event *AuditEvent) error {
	// Set created time if not set
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	// Log the event
	s.logger.Info("Audit event",
		zap.String("type", string(event.Type)),
		zap.String("severity", string(event.Severity)),
		zap.String("action", event.Action),
		zap.String("description", event.Description),
		zap.Int("user_id", event.UserID),
		zap.String("request_id", event.RequestID),
	)

	// Store the event in the database
	result := s.db.WithContext(ctx).Create(event)
	if result.Error != nil {
		return fmt.Errorf("failed to store audit event: %w", result.Error)
	}

	return nil
}

// GetEvents retrieves audit events with optional filtering
func (s *serviceImpl) GetEvents(
	ctx context.Context,
	userID int,
	eventType EventType,
	severity EventSeverity,
	startTime, endTime time.Time,
	limit, offset int,
) ([]*AuditEvent, error) {
	var events []*AuditEvent

	// Build query
	query := s.db.WithContext(ctx)

	// Apply filters
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if eventType != "" {
		query = query.Where("type = ?", eventType)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Execute query
	result := query.Order("created_at DESC").Find(&events)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve audit events: %w", result.Error)
	}

	return events, nil
}

// GetEventByID retrieves an audit event by ID
func (s *serviceImpl) GetEventByID(ctx context.Context, id uint) (*AuditEvent, error) {
	var event AuditEvent

	result := s.db.WithContext(ctx).First(&event, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve audit event: %w", result.Error)
	}

	return &event, nil
}

// GetEventsByRequestID retrieves audit events by request ID
func (s *serviceImpl) GetEventsByRequestID(ctx context.Context, requestID string) ([]*AuditEvent, error) {
	var events []*AuditEvent

	result := s.db.WithContext(ctx).Where("request_id = ?", requestID).Order("created_at ASC").Find(&events)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve audit events: %w", result.Error)
	}

	return events, nil
}

// CreateAuditEvent creates a new audit event
func CreateAuditEvent(
	userID int,
	eventType EventType,
	severity EventSeverity,
	action string,
	description string,
	metadata interface{},
	ip string,
	userAgent string,
	requestID string,
) (*AuditEvent, error) {
	// Convert metadata to JSON
	var metadataJSON string
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(metadataBytes)
	}

	// Create event
	event := &AuditEvent{
		UserID:      userID,
		Type:        eventType,
		Severity:    severity,
		Action:      action,
		Description: description,
		Metadata:    metadataJSON,
		IP:          ip,
		UserAgent:   userAgent,
		RequestID:   requestID,
		CreatedAt:   time.Now().UTC(),
	}

	return event, nil
}
