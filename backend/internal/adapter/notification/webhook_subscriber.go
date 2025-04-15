package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// WebhookConfig contains configuration for webhook notifications
type WebhookConfig struct {
	Enabled   bool
	URL       string
	Method    string
	Headers   map[string]string
	MinLevel  AlertLevel
	Timeout   time.Duration
	BatchSize int
}

// WebhookSubscriber implements the AlertSubscriber interface for webhook notifications
type WebhookSubscriber struct {
	config WebhookConfig
	logger *zerolog.Logger
	client *http.Client
	batch  []Alert
}

// NewWebhookSubscriber creates a new webhook subscriber
func NewWebhookSubscriber(config WebhookConfig, logger *zerolog.Logger) *WebhookSubscriber {
	if config.Method == "" {
		config.Method = http.MethodPost
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 1 // Default to sending alerts individually
	}

	return &WebhookSubscriber{
		config: config,
		logger: logger,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		batch: make([]Alert, 0, config.BatchSize),
	}
}

// HandleAlert processes an alert and sends a webhook notification if needed
func (s *WebhookSubscriber) HandleAlert(alert Alert) error {
	if !s.config.Enabled {
		return nil
	}

	// Check if alert level meets minimum threshold
	if !s.shouldSendAlert(alert) {
		return nil
	}

	// If batching is disabled, send immediately
	if s.config.BatchSize <= 1 {
		return s.sendWebhook([]Alert{alert})
	}

	// Add to batch
	s.batch = append(s.batch, alert)

	// If batch is full, send it
	if len(s.batch) >= s.config.BatchSize {
		batch := s.batch
		s.batch = make([]Alert, 0, s.config.BatchSize)
		return s.sendWebhook(batch)
	}

	return nil
}

// GetName returns the name of the subscriber
func (s *WebhookSubscriber) GetName() string {
	return "webhook"
}

// FlushBatch sends any pending alerts in the batch
func (s *WebhookSubscriber) FlushBatch() error {
	if len(s.batch) == 0 {
		return nil
	}

	batch := s.batch
	s.batch = make([]Alert, 0, s.config.BatchSize)
	return s.sendWebhook(batch)
}

// shouldSendAlert determines if an alert should be sent
func (s *WebhookSubscriber) shouldSendAlert(alert Alert) bool {
	// Check minimum level
	switch s.config.MinLevel {
	case AlertLevelCritical:
		return alert.Level == AlertLevelCritical
	case AlertLevelError:
		return alert.Level == AlertLevelCritical || alert.Level == AlertLevelError
	case AlertLevelWarning:
		return alert.Level == AlertLevelCritical || alert.Level == AlertLevelError || alert.Level == AlertLevelWarning
	default:
		return true
	}
}

// sendWebhook sends alerts to the webhook endpoint
func (s *WebhookSubscriber) sendWebhook(alerts []Alert) error {
	// Prepare payload
	payload := map[string]interface{}{
		"alerts": alerts,
		"count":  len(alerts),
		"time":   time.Now().Format(time.RFC3339),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to marshal webhook payload")
		return err
	}

	// Create request
	req, err := http.NewRequest(s.config.Method, s.config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create webhook request")
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range s.config.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to send webhook request")
		return err
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("webhook returned non-success status: %d", resp.StatusCode)
		s.logger.Error().Err(err).Int("status_code", resp.StatusCode).Msg("Webhook request failed")
		return err
	}

	s.logger.Info().
		Int("count", len(alerts)).
		Int("status_code", resp.StatusCode).
		Msg("Webhook alerts sent")

	return nil
}
