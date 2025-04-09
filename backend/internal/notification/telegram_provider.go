package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	telegramAPIBaseURL = "https://api.telegram.org/bot%s"
	telegramSendMessage = "sendMessage"
	telegramSendPhoto = "sendPhoto"
	telegramSendDocument = "sendDocument"
)

// TelegramProvider implements the NotificationProvider interface for Telegram
type TelegramProvider struct {
	*BaseProvider
	token    string
	chatIDs  []string
	client   *http.Client
	rateLimit int
	lastSent time.Time
	logger   *zap.Logger
}

// NewTelegramProvider creates a new Telegram provider
func NewTelegramProvider(logger *zap.Logger) *TelegramProvider {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &TelegramProvider{
		BaseProvider: NewBaseProvider("telegram"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		rateLimit: 20, // Default: 20 messages per minute
		logger:    logger,
	}
}

// Initialize initializes the Telegram provider
func (p *TelegramProvider) Initialize(config map[string]interface{}) error {
	p.SetConfig(config)

	// Get token
	token, ok := p.GetConfigString("token")
	if !ok || token == "" {
		p.SetAvailable(false)
		return fmt.Errorf("%w: missing or invalid token", ErrInvalidConfiguration)
	}
	p.token = token

	// Get chat IDs
	chatIDs, ok := p.GetConfigStringSlice("chat_ids")
	if !ok || len(chatIDs) == 0 {
		p.SetAvailable(false)
		return fmt.Errorf("%w: missing or invalid chat_ids", ErrInvalidConfiguration)
	}
	p.chatIDs = chatIDs

	// Get rate limit
	rateLimit, ok := p.GetConfigInt("rate_limit")
	if ok && rateLimit > 0 {
		p.rateLimit = rateLimit
	}

	// Check if enabled
	enabled, ok := p.GetConfigBool("enabled")
	if ok && !enabled {
		p.SetAvailable(false)
		p.logger.Info("Telegram provider is disabled")
		return nil
	}

	// Test connection
	if err := p.testConnection(); err != nil {
		p.SetAvailable(false)
		return fmt.Errorf("failed to connect to Telegram API: %w", err)
	}

	p.SetAvailable(true)
	p.logger.Info("Telegram provider initialized",
		zap.Int("chat_count", len(p.chatIDs)),
		zap.Int("rate_limit", p.rateLimit),
	)

	return nil
}

// Send sends a notification via Telegram
func (p *TelegramProvider) Send(ctx context.Context, notification *Notification) (*NotificationResult, error) {
	if !p.IsAvailable() {
		return nil, ErrProviderNotAvailable
	}

	// Check rate limit
	if !p.checkRateLimit() {
		return nil, ErrRateLimited
	}

	// Format message
	message := fmt.Sprintf("*%s*\n\n%s", notification.Title, notification.Message)

	// Send to all chat IDs
	var lastErr error
	var lastResult *NotificationResult

	for _, chatID := range p.chatIDs {
		// Check if we have attachments
		if len(notification.Attachments) > 0 {
			for _, attachment := range notification.Attachments {
				result, err := p.sendAttachment(ctx, chatID, message, &attachment)
				if err != nil {
					lastErr = err
					p.logger.Error("Failed to send attachment to Telegram",
						zap.String("notification_id", notification.ID),
						zap.String("chat_id", chatID),
						zap.Error(err),
					)
				} else {
					lastResult = result
				}
			}
		} else {
			// Send text message
			result, err := p.sendMessage(ctx, chatID, message)
			if err != nil {
				lastErr = err
				p.logger.Error("Failed to send message to Telegram",
					zap.String("notification_id", notification.ID),
					zap.String("chat_id", chatID),
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

// sendMessage sends a text message to a Telegram chat
func (p *TelegramProvider) sendMessage(ctx context.Context, chatID, message string) (*NotificationResult, error) {
	url := fmt.Sprintf(telegramAPIBaseURL+"/"+telegramSendMessage, p.token)

	// Prepare request body
	body := map[string]interface{}{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "Markdown",
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
		return nil, fmt.Errorf("telegram API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Check if successful
	ok, _ := result["ok"].(bool)
	if !ok {
		return nil, fmt.Errorf("telegram API returned error: %v", result["description"])
	}

	return &NotificationResult{
		NotificationID: "",  // Will be set by the service
		ProviderName:   p.GetName(),
		Success:        true,
		Timestamp:      time.Now(),
	}, nil
}

// sendAttachment sends an attachment to a Telegram chat
func (p *TelegramProvider) sendAttachment(ctx context.Context, chatID, caption string, attachment *Attachment) (*NotificationResult, error) {
	var url string
	var method string

	// Determine method based on attachment type
	switch attachment.Type {
	case "image":
		url = fmt.Sprintf(telegramAPIBaseURL+"/"+telegramSendPhoto, p.token)
		method = "photo"
	default:
		url = fmt.Sprintf(telegramAPIBaseURL+"/"+telegramSendDocument, p.token)
		method = "document"
	}

	// Create multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add chat ID
	if err := w.WriteField("chat_id", chatID); err != nil {
		return nil, err
	}

	// Add caption
	if caption != "" {
		if err := w.WriteField("caption", caption); err != nil {
			return nil, err
		}
		if err := w.WriteField("parse_mode", "Markdown"); err != nil {
			return nil, err
		}
	}

	// Add file
	var err error
	if attachment.URL != "" {
		// Use URL
		if err := w.WriteField(method+"_url", attachment.URL); err != nil {
			return nil, err
		}
	} else if len(attachment.Data) > 0 {
		// Use data
		filename := attachment.Filename
		if filename == "" {
			filename = "attachment"
		}

		var fw io.Writer
		fw, err = w.CreateFormFile(method, filename)
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
		return nil, fmt.Errorf("telegram API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Check if successful
	ok, _ := result["ok"].(bool)
	if !ok {
		return nil, fmt.Errorf("telegram API returned error: %v", result["description"])
	}

	return &NotificationResult{
		NotificationID: "",  // Will be set by the service
		ProviderName:   p.GetName(),
		Success:        true,
		Timestamp:      time.Now(),
	}, nil
}

// testConnection tests the connection to the Telegram API
func (p *TelegramProvider) testConnection() error {
	url := fmt.Sprintf(telegramAPIBaseURL+"/getMe", p.token)

	resp, err := p.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	ok, _ := result["ok"].(bool)
	if !ok {
		return fmt.Errorf("telegram API returned error: %v", result["description"])
	}

	return nil
}

// checkRateLimit checks if we're within the rate limit
func (p *TelegramProvider) checkRateLimit() bool {
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
