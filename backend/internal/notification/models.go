package notification

import (
	"time"

	"github.com/google/uuid"
)

// NotificationLevel represents the severity level of a notification
type NotificationLevel string

const (
	// LevelInfo is for informational notifications
	LevelInfo NotificationLevel = "INFO"
	// LevelWarning is for warning notifications
	LevelWarning NotificationLevel = "WARNING"
	// LevelError is for error notifications
	LevelError NotificationLevel = "ERROR"
	// LevelCritical is for critical notifications
	LevelCritical NotificationLevel = "CRITICAL"
	// LevelTrade is for trade-related notifications
	LevelTrade NotificationLevel = "TRADE"
)

// Notification represents a notification to be sent
type Notification struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Level       NotificationLevel      `json:"level"`
	Source      string                 `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Attachments []Attachment           `json:"attachments,omitempty"`
	Providers   []string               `json:"providers,omitempty"`
	Priority    int                    `json:"priority"`
	Retries     int                    `json:"-"`
	MaxRetries  int                    `json:"-"`
}

// NewNotification creates a new notification
func NewNotification(title, message string, level NotificationLevel) *Notification {
	return &Notification{
		ID:         uuid.New().String(),
		Title:      title,
		Message:    message,
		Level:      level,
		Timestamp:  time.Now(),
		Priority:   getPriorityForLevel(level),
		MaxRetries: 3,
	}
}

// Attachment represents a file attachment for a notification
type Attachment struct {
	Type        string `json:"type"`
	URL         string `json:"url,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Data        []byte `json:"data,omitempty"`
	Filename    string `json:"filename,omitempty"`
}

// NotificationTemplate represents a template for notifications
type NotificationTemplate struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Level     NotificationLevel `json:"level"`
	Providers []string          `json:"providers,omitempty"`
	Priority  int               `json:"priority"`
}

// NotificationResult represents the result of sending a notification
type NotificationResult struct {
	NotificationID string    `json:"notification_id"`
	ProviderName   string    `json:"provider_name"`
	Success        bool      `json:"success"`
	Error          string    `json:"error,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}

// getPriorityForLevel returns a default priority based on the notification level
func getPriorityForLevel(level NotificationLevel) int {
	switch level {
	case LevelCritical:
		return 100
	case LevelError:
		return 80
	case LevelWarning:
		return 60
	case LevelTrade:
		return 40
	case LevelInfo:
		return 20
	default:
		return 0
	}
}
