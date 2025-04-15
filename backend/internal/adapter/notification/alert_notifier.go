package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
)

// AlertLevel defines the severity level of an alert
type AlertLevel string

const (
	// AlertLevelInfo is for informational alerts
	AlertLevelInfo AlertLevel = "info"
	// AlertLevelWarning is for warning alerts
	AlertLevelWarning AlertLevel = "warning"
	// AlertLevelError is for error alerts
	AlertLevelError AlertLevel = "error"
	// AlertLevelCritical is for critical alerts
	AlertLevelCritical AlertLevel = "critical"
)

// Alert represents a system alert
type Alert struct {
	ID        string     `json:"id"`
	Level     AlertLevel `json:"level"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Source    string     `json:"source"`
	Timestamp time.Time  `json:"timestamp"`
	Resolved  bool       `json:"resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

// AlertNotifier implements the StatusNotifier interface with alert generation
type AlertNotifier struct {
	logger      *zerolog.Logger
	alertChan   chan Alert
	enabled     bool
	alertStore  []Alert
	maxAlerts   int
	subscribers []AlertSubscriber
}

// AlertSubscriber defines the interface for alert subscribers
type AlertSubscriber interface {
	// HandleAlert processes an alert
	HandleAlert(alert Alert) error
	// GetName returns the name of the subscriber
	GetName() string
}

// NewAlertNotifier creates a new alert notifier
func NewAlertNotifier(logger *zerolog.Logger, maxAlerts int) *AlertNotifier {
	if maxAlerts <= 0 {
		maxAlerts = 100
	}

	notifier := &AlertNotifier{
		logger:      logger,
		alertChan:   make(chan Alert, 100),
		enabled:     true,
		alertStore:  make([]Alert, 0, maxAlerts),
		maxAlerts:   maxAlerts,
		subscribers: make([]AlertSubscriber, 0),
	}

	// Start the alert processor
	go notifier.processAlerts()

	return notifier
}

// NotifyStatusChange sends a notification about a status change
func (n *AlertNotifier) NotifyStatusChange(ctx context.Context, component string, oldStatus, newStatus status.Status, message string) error {
	if !n.enabled {
		return nil
	}

	// Determine alert level based on status change
	var level AlertLevel
	switch newStatus {
	case status.StatusError:
		level = AlertLevelError
	case status.StatusWarning:
		level = AlertLevelWarning
	case status.StatusStopped:
		if oldStatus == status.StatusRunning {
			level = AlertLevelWarning
		} else {
			level = AlertLevelInfo
		}
	case status.StatusRunning:
		if oldStatus == status.StatusError || oldStatus == status.StatusWarning {
			level = AlertLevelInfo
		} else {
			return nil // Don't alert for normal transitions to running
		}
	default:
		return nil // Don't alert for other status changes
	}

	// Create alert
	alert := Alert{
		ID:        fmt.Sprintf("comp-%s-%d", component, time.Now().UnixNano()),
		Level:     level,
		Title:     fmt.Sprintf("Component %s status changed", component),
		Message:   fmt.Sprintf("Status changed from %s to %s: %s", oldStatus, newStatus, message),
		Source:    component,
		Timestamp: time.Now(),
		Resolved:  newStatus == status.StatusRunning,
	}

	// Try to send to channel with timeout
	select {
	case n.alertChan <- alert:
		n.logger.Debug().
			Str("component", component).
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Str("level", string(level)).
			Msg("Component status change alert queued")
	case <-time.After(time.Second):
		n.logger.Warn().
			Str("component", component).
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Msg("Failed to queue component status change alert: channel full")
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// NotifySystemStatusChange sends a notification about a system status change
func (n *AlertNotifier) NotifySystemStatusChange(ctx context.Context, oldStatus, newStatus status.Status, message string) error {
	if !n.enabled {
		return nil
	}

	// Determine alert level based on status change
	var level AlertLevel
	switch newStatus {
	case status.StatusError:
		level = AlertLevelCritical
	case status.StatusWarning:
		level = AlertLevelWarning
	case status.StatusStopped:
		level = AlertLevelWarning
	case status.StatusRunning:
		if oldStatus == status.StatusError || oldStatus == status.StatusWarning {
			level = AlertLevelInfo
		} else {
			return nil // Don't alert for normal transitions to running
		}
	default:
		return nil // Don't alert for other status changes
	}

	// Create alert
	alert := Alert{
		ID:        fmt.Sprintf("system-%d", time.Now().UnixNano()),
		Level:     level,
		Title:     "System status changed",
		Message:   fmt.Sprintf("Status changed from %s to %s: %s", oldStatus, newStatus, message),
		Source:    "system",
		Timestamp: time.Now(),
		Resolved:  newStatus == status.StatusRunning,
	}

	// Try to send to channel with timeout
	select {
	case n.alertChan <- alert:
		n.logger.Debug().
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Str("level", string(level)).
			Msg("System status change alert queued")
	case <-time.After(time.Second):
		n.logger.Warn().
			Str("old_status", string(oldStatus)).
			Str("new_status", string(newStatus)).
			Msg("Failed to queue system status change alert: channel full")
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// CreateAlert creates a new alert
func (n *AlertNotifier) CreateAlert(ctx context.Context, level AlertLevel, title, message, source string) error {
	if !n.enabled {
		return nil
	}

	alert := Alert{
		ID:        fmt.Sprintf("%s-%d", source, time.Now().UnixNano()),
		Level:     level,
		Title:     title,
		Message:   message,
		Source:    source,
		Timestamp: time.Now(),
		Resolved:  false,
	}

	// Try to send to channel with timeout
	select {
	case n.alertChan <- alert:
		n.logger.Debug().
			Str("source", source).
			Str("level", string(level)).
			Str("title", title).
			Msg("Alert queued")
	case <-time.After(time.Second):
		n.logger.Warn().
			Str("source", source).
			Str("level", string(level)).
			Str("title", title).
			Msg("Failed to queue alert: channel full")
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// ResolveAlert marks an alert as resolved
func (n *AlertNotifier) ResolveAlert(ctx context.Context, alertID string) error {
	n.logger.Debug().Str("alert_id", alertID).Msg("Resolving alert")

	for i, alert := range n.alertStore {
		if alert.ID == alertID && !alert.Resolved {
			now := time.Now()
			n.alertStore[i].Resolved = true
			n.alertStore[i].ResolvedAt = &now

			// Notify subscribers
			for _, subscriber := range n.subscribers {
				if err := subscriber.HandleAlert(n.alertStore[i]); err != nil {
					n.logger.Error().
						Err(err).
						Str("subscriber", subscriber.GetName()).
						Str("alert_id", alertID).
						Msg("Failed to notify subscriber about resolved alert")
				}
			}

			return nil
		}
	}

	return fmt.Errorf("alert not found or already resolved: %s", alertID)
}

// GetAlerts returns all alerts
func (n *AlertNotifier) GetAlerts(ctx context.Context, onlyActive bool) []Alert {
	if onlyActive {
		activeAlerts := make([]Alert, 0)
		for _, alert := range n.alertStore {
			if !alert.Resolved {
				activeAlerts = append(activeAlerts, alert)
			}
		}
		return activeAlerts
	}
	return n.alertStore
}

// AddSubscriber adds a subscriber for alerts
func (n *AlertNotifier) AddSubscriber(subscriber AlertSubscriber) {
	n.subscribers = append(n.subscribers, subscriber)
	n.logger.Info().Str("subscriber", subscriber.GetName()).Msg("Added alert subscriber")
}

// Enable enables the notifier
func (n *AlertNotifier) Enable() {
	n.enabled = true
	n.logger.Info().Msg("Alert notifier enabled")
}

// Disable disables the notifier
func (n *AlertNotifier) Disable() {
	n.enabled = false
	n.logger.Info().Msg("Alert notifier disabled")
}

// processAlerts processes alerts in the background
func (n *AlertNotifier) processAlerts() {
	for alert := range n.alertChan {
		n.processAlert(alert)
	}
}

// processAlert processes a single alert
func (n *AlertNotifier) processAlert(alert Alert) {
	// Log the alert
	logEvent := n.logger.Info()
	if alert.Level == AlertLevelError || alert.Level == AlertLevelCritical {
		logEvent = n.logger.Error()
	} else if alert.Level == AlertLevelWarning {
		logEvent = n.logger.Warn()
	}

	logEvent.
		Str("alert_id", alert.ID).
		Str("level", string(alert.Level)).
		Str("source", alert.Source).
		Str("title", alert.Title).
		Str("message", alert.Message).
		Time("timestamp", alert.Timestamp).
		Bool("resolved", alert.Resolved).
		Msg("Alert generated")

	// Store the alert
	n.storeAlert(alert)

	// Notify subscribers
	for _, subscriber := range n.subscribers {
		if err := subscriber.HandleAlert(alert); err != nil {
			n.logger.Error().
				Err(err).
				Str("subscriber", subscriber.GetName()).
				Str("alert_id", alert.ID).
				Msg("Failed to notify subscriber about alert")
		}
	}
}

// storeAlert stores an alert in the alert store
func (n *AlertNotifier) storeAlert(alert Alert) {
	// Add to the store
	n.alertStore = append(n.alertStore, alert)

	// Trim if needed
	if len(n.alertStore) > n.maxAlerts {
		// Remove oldest alerts
		excess := len(n.alertStore) - n.maxAlerts
		n.alertStore = n.alertStore[excess:]
	}
}
