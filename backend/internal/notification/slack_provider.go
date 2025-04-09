package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	slackAPIBaseURL = "https://slack.com/api"
	slackPostMessage = "chat.postMessage"
	slackUploadFile = "files.upload"
)

// SlackProvider implements the NotificationProvider interface for Slack
type SlackProvider struct {
	*BaseProvider
	token     string
	channels  []string
	client    *http.Client
	rateLimit int
	lastSent  time.Time
	logger    *zap.Logger
}

// NewSlackProvider creates a new Slack provider
func NewSlackProvider(logger *zap.Logger) *SlackProvider {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &SlackProvider{
		BaseProvider: NewBaseProvider("slack"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		rateLimit: 30, // Default: 30 messages per minute
		logger:    logger,
	}
}

// Initialize initializes the Slack provider
func (p *SlackProvider) Initialize(config map[string]interface{}) error {
	p.SetConfig(config)

	// Get token
	token, ok := p.GetConfigString("token")
	if !ok || token == "" {
		p.SetAvailable(false)
		return fmt.Errorf("%w: missing or invalid token", ErrInvalidConfiguration)
	}
	p.token = token

	// Get channels
	channels, ok := p.GetConfigStringSlice("channels")
	if !ok || len(channels) == 0 {
		p.SetAvailable(false)
		return fmt.Errorf("%w: missing or invalid channels", ErrInvalidConfiguration)
	}
	p.channels = channels

	// Get rate limit
	rateLimit, ok := p.GetConfigInt("rate_limit")
	if ok && rateLimit > 0 {
		p.rateLimit = rateLimit
	}

	// Check if enabled
	enabled, ok := p.GetConfigBool("enabled")
	if ok && !enabled {
		p.SetAvailable(false)
		p.logger.Info("Slack provider is disabled")
		return nil
	}

	// Test connection
	if err := p.testConnection(); err != nil {
		p.SetAvailable(false)
		return fmt.Errorf("failed to connect to Slack API: %w", err)
	}

	p.SetAvailable(true)
	p.logger.Info("Slack provider initialized",
		zap.Int("channel_count", len(p.channels)),
		zap.Int("rate_limit", p.rateLimit),
	)

	return nil
}

// Send sends a notification via Slack
func (p *SlackProvider) Send(ctx context.Context, notification *Notification) (*NotificationResult, error) {
	if !p.IsAvailable() {
		return nil, ErrProviderNotAvailable
	}

	// Check rate limit
	if !p.checkRateLimit() {
		return nil, ErrRateLimited
	}

	// Send to all channels
	var lastErr error
	var lastResult *NotificationResult

	for _, channel := range p.channels {
		// Check if we have attachments
		if len(notification.Attachments) > 0 {
			for _, attachment := range notification.Attachments {
				result, err := p.sendAttachment(ctx, channel, notification.Title, notification.Message, &attachment)
				if err != nil {
					lastErr = err
					p.logger.Error("Failed to send attachment to Slack",
						zap.String("notification_id", notification.ID),
						zap.String("channel", channel),
						zap.Error(err),
					)
				} else {
					lastResult = result
				}
			}
		} else {
			// Send text message
			result, err := p.sendMessage(ctx, channel, notification.Title, notification.Message, notification.Level)
			if err != nil {
				lastErr = err
				p.logger.Error("Failed to send message to Slack",
					zap.String("notification_id", notification.ID),
					zap.String("channel", channel),
					zap.Error(err),
				)
			} else {
				lastResult = result
			}
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, lastErr)
	}

	return lastResult, nil
}

// sendMessage sends a text message to a Slack channel
func (p *SlackProvider) sendMessage(ctx context.Context, channel, title, message string, level NotificationLevel) (*NotificationResult, error) {
	url := fmt.Sprintf("%s/%s", slackAPIBaseURL, slackPostMessage)

	// Determine color based on level
	color := p.getLevelColor(level)

	// Prepare request body
	body := map[string]interface{}{
		"channel": channel,
		"attachments": []map[string]interface{}{
			{
				"title":      title,
				"text":       message,
				"color":      color,
				"mrkdwn_in":  []string{"text"},
				"fallback":   fmt.Sprintf("%s: %s", title, message),
				"ts":         time.Now().Unix(),
			},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.token)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Update rate limit tracking
	p.lastSent = time.Now()

	// Check response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("slack API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Check if successful
	ok, _ := result["ok"].(bool)
	if !ok {
		return nil, fmt.Errorf("slack API returned error: %v", result["error"])
	}

	return &NotificationResult{
		NotificationID: "",  // Will be set by the service
		ProviderName:   p.GetName(),
		Success:        true,
		Timestamp:      time.Now(),
	}, nil
}

// sendAttachment sends an attachment to a Slack channel
func (p *SlackProvider) sendAttachment(ctx context.Context, channel, title, message string, attachment *Attachment) (*NotificationResult, error) {
	url := fmt.Sprintf("%s/%s", slackAPIBaseURL, slackUploadFile)

	// Create multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add channel
	if err := w.WriteField("channels", channel); err != nil {
		return nil, err
	}

	// Add title and message
	if title != "" {
		if err := w.WriteField("title", title); err != nil {
			return nil, err
		}
	}
	if message != "" {
		if err := w.WriteField("initial_comment", message); err != nil {
			return nil, err
		}
	}

	// Add file
	var err error
	if attachment.URL != "" {
		// For URL, we need to download the file first
		// This is a simplified version - in a real implementation, you'd want to stream the file
		resp, err := http.Get(attachment.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Get filename from URL
		filename := attachment.Filename
		if filename == "" {
			parts := strings.Split(attachment.URL, "/")
			filename = parts[len(parts)-1]
			if filename == "" {
				filename = "attachment"
			}
		}

		// Create form file
		var fw io.Writer
		fw, err = w.CreateFormFile("file", filename)
		if err != nil {
			return nil, err
		}

		// Copy file data
		_, err = io.Copy(fw, resp.Body)
		if err != nil {
			return nil, err
		}
	} else if len(attachment.Data) > 0 {
		// Use data
		filename := attachment.Filename
		if filename == "" {
			filename = "attachment"
		}

		var fw io.Writer
		fw, err = w.CreateFormFile("file", filename)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fw, bytes.NewReader(attachment.Data))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("attachment has no URL or data")
	}

	// Close multipart writer
	w.Close()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+p.token)

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Update rate limit tracking
	p.lastSent = time.Now()

	// Check response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("slack API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Check if successful
	ok, _ := result["ok"].(bool)
	if !ok {
		return nil, fmt.Errorf("slack API returned error: %v", result["error"])
	}

	return &NotificationResult{
		NotificationID: "",  // Will be set by the service
		ProviderName:   p.GetName(),
		Success:        true,
		Timestamp:      time.Now(),
	}, nil
}

// testConnection tests the connection to the Slack API
func (p *SlackProvider) testConnection() error {
	url := fmt.Sprintf("%s/auth.test", slackAPIBaseURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("slack API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	ok, _ := result["ok"].(bool)
	if !ok {
		return fmt.Errorf("slack API returned error: %v", result["error"])
	}

	return nil
}

// checkRateLimit checks if we're within the rate limit
func (p *SlackProvider) checkRateLimit() bool {
	// If we haven't sent anything yet, we're good
	if p.lastSent.IsZero() {
		return true
	}

	// Calculate time since last send
	elapsed := time.Since(p.lastSent)

	// Calculate minimum time between messages
	minTime := time.Minute / time.Duration(p.rateLimit)

	// If we've waited long enough, we're good
	return elapsed >= minTime
}

// getLevelColor returns a color for a notification level
func (p *SlackProvider) getLevelColor(level NotificationLevel) string {
	switch level {
	case LevelInfo:
		return "#2196F3" // Blue
	case LevelWarning:
		return "#FF9800" // Orange
	case LevelError:
		return "#F44336" // Red
	case LevelCritical:
		return "#9C27B0" // Purple
	case LevelTrade:
		return "#4CAF50" // Green
	default:
		return "#9E9E9E" // Grey
	}
}
