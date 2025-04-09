package notification

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Common errors
var (
	ErrServiceNotInitialized = errors.New("notification service not initialized")
	ErrTemplateNotFound      = errors.New("notification template not found")
	ErrNoProvidersAvailable  = errors.New("no notification providers available")
	ErrQueueFull             = errors.New("notification queue is full")
)

// NotificationService manages sending notifications through various providers
type NotificationService struct {
	registry         *ProviderRegistry
	templates        map[string]*NotificationTemplate
	defaultProviders []string
	queue            chan *Notification
	results          chan *NotificationResult
	workerCount      int
	workerWg         sync.WaitGroup
	ctx              context.Context
	cancel           context.CancelFunc
	initialized      bool
	logger           *zap.Logger
	mu               sync.RWMutex
}

// NotificationServiceConfig contains configuration for the notification service
type NotificationServiceConfig struct {
	DefaultProviders []string               `json:"default_providers"`
	QueueCapacity    int                    `json:"queue_capacity"`
	WorkerCount      int                    `json:"worker_count"`
	Templates        map[string]interface{} `json:"templates"`
	Providers        map[string]interface{} `json:"providers"`
}

// NewNotificationService creates a new notification service
func NewNotificationService(logger *zap.Logger) *NotificationService {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &NotificationService{
		registry:    NewProviderRegistry(),
		templates:   make(map[string]*NotificationTemplate),
		queue:       nil, // Will be initialized in Initialize
		results:     nil, // Will be initialized in Initialize
		ctx:         ctx,
		cancel:      cancel,
		initialized: false,
		logger:      logger,
	}
}

// Initialize initializes the notification service
func (s *NotificationService) Initialize(config *NotificationServiceConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if config == nil {
		return errors.New("notification service config is nil")
	}

	// Set default providers
	s.defaultProviders = config.DefaultProviders

	// Set queue capacity and worker count
	queueCapacity := 1000
	if config.QueueCapacity > 0 {
		queueCapacity = config.QueueCapacity
	}

	workerCount := 5
	if config.WorkerCount > 0 {
		workerCount = config.WorkerCount
	}
	s.workerCount = workerCount

	// Initialize queue and results channel
	s.queue = make(chan *Notification, queueCapacity)
	s.results = make(chan *NotificationResult, queueCapacity)

	// Initialize templates
	if config.Templates != nil {
		for id, templateConfig := range config.Templates {
			if templateMap, ok := templateConfig.(map[string]interface{}); ok {
				template := &NotificationTemplate{
					ID: id,
				}

				if title, ok := templateMap["title"].(string); ok {
					template.Title = title
				}

				if message, ok := templateMap["message"].(string); ok {
					template.Message = message
				}

				if level, ok := templateMap["level"].(string); ok {
					template.Level = NotificationLevel(level)
				}

				if providers, ok := templateMap["providers"].([]interface{}); ok {
					providerStrings := make([]string, 0, len(providers))
					for _, p := range providers {
						if providerStr, ok := p.(string); ok {
							providerStrings = append(providerStrings, providerStr)
						}
					}
					template.Providers = providerStrings
				}

				if priority, ok := templateMap["priority"].(int); ok {
					template.Priority = priority
				} else {
					template.Priority = getPriorityForLevel(template.Level)
				}

				s.templates[id] = template
			}
		}
	}

	// Start workers
	for i := 0; i < workerCount; i++ {
		s.workerWg.Add(1)
		go s.worker(i)
	}

	s.initialized = true
	s.logger.Info("Notification service initialized",
		zap.Int("worker_count", workerCount),
		zap.Int("queue_capacity", queueCapacity),
		zap.Strings("default_providers", s.defaultProviders),
		zap.Int("template_count", len(s.templates)),
	)

	return nil
}

// RegisterProvider registers a notification provider
func (s *NotificationService) RegisterProvider(provider NotificationProvider) {
	s.registry.Register(provider)
	s.logger.Info("Registered notification provider", zap.String("provider", provider.GetName()))
}

// Send sends a notification
func (s *NotificationService) Send(notification *Notification) error {
	if !s.isInitialized() {
		return ErrServiceNotInitialized
	}

	// Generate ID if not set
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}

	// Set timestamp if not set
	if notification.Timestamp.IsZero() {
		notification.Timestamp = time.Now()
	}

	// Set providers if not set
	if len(notification.Providers) == 0 {
		notification.Providers = s.defaultProviders
	}

	// Set priority if not set
	if notification.Priority == 0 {
		notification.Priority = getPriorityForLevel(notification.Level)
	}

	// Send to queue
	select {
	case s.queue <- notification:
		s.logger.Debug("Notification queued",
			zap.String("id", notification.ID),
			zap.String("title", notification.Title),
			zap.String("level", string(notification.Level)),
		)
		return nil
	default:
		s.logger.Warn("Notification queue is full, dropping notification",
			zap.String("id", notification.ID),
			zap.String("title", notification.Title),
		)
		return ErrQueueFull
	}
}

// SendWithTemplate sends a notification using a template
func (s *NotificationService) SendWithTemplate(templateID string, data map[string]interface{}) error {
	if !s.isInitialized() {
		return ErrServiceNotInitialized
	}

	// Get template
	template, ok := s.getTemplate(templateID)
	if !ok {
		return ErrTemplateNotFound
	}

	// Create notification from template
	notification := &Notification{
		ID:        uuid.New().String(),
		Level:     template.Level,
		Providers: template.Providers,
		Priority:  template.Priority,
		Timestamp: time.Now(),
		Data:      data,
	}

	// Parse title template
	titleTmpl, err := s.parseTemplate(template.Title)
	if err != nil {
		s.logger.Error("Failed to parse title template",
			zap.String("template_id", templateID),
			zap.Error(err),
		)
		return err
	}

	// Execute title template
	var titleBuf bytes.Buffer
	if err := titleTmpl.Execute(&titleBuf, data); err != nil {
		s.logger.Error("Failed to execute title template",
			zap.String("template_id", templateID),
			zap.Error(err),
		)
		return err
	}
	notification.Title = titleBuf.String()

	// Parse message template
	messageTmpl, err := s.parseTemplate(template.Message)
	if err != nil {
		s.logger.Error("Failed to parse message template",
			zap.String("template_id", templateID),
			zap.Error(err),
		)
		return err
	}

	// Execute message template
	var messageBuf bytes.Buffer
	if err := messageTmpl.Execute(&messageBuf, data); err != nil {
		s.logger.Error("Failed to execute message template",
			zap.String("template_id", templateID),
			zap.Error(err),
		)
		return err
	}
	notification.Message = messageBuf.String()

	// Send notification
	return s.Send(notification)
}

// GetResults returns a channel for notification results
func (s *NotificationService) GetResults() <-chan *NotificationResult {
	return s.results
}

// Shutdown shuts down the notification service
func (s *NotificationService) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initialized {
		return
	}

	s.logger.Info("Shutting down notification service")
	s.cancel()
	s.workerWg.Wait()
	close(s.queue)
	close(s.results)
	s.initialized = false
	s.logger.Info("Notification service shut down")
}

// worker processes notifications from the queue
func (s *NotificationService) worker(id int) {
	defer s.workerWg.Done()

	s.logger.Debug("Starting notification worker", zap.Int("worker_id", id))

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Debug("Stopping notification worker", zap.Int("worker_id", id))
			return
		case notification, ok := <-s.queue:
			if !ok {
				return
			}

			s.processNotification(notification)
		}
	}
}

// processNotification processes a notification
func (s *NotificationService) processNotification(notification *Notification) {
	// Get available providers
	availableProviders := s.registry.GetAvailable()
	if len(availableProviders) == 0 {
		s.logger.Warn("No notification providers available",
			zap.String("notification_id", notification.ID),
		)
		s.sendResult(&NotificationResult{
			NotificationID: notification.ID,
			ProviderName:   "none",
			Success:        false,
			Error:          ErrNoProvidersAvailable.Error(),
			Timestamp:      time.Now(),
		})
		return
	}

	// Filter providers based on notification.Providers
	var providers []NotificationProvider
	if len(notification.Providers) > 0 {
		for _, providerName := range notification.Providers {
			for _, provider := range availableProviders {
				if provider.GetName() == providerName {
					providers = append(providers, provider)
					break
				}
			}
		}
	} else {
		providers = availableProviders
	}

	if len(providers) == 0 {
		s.logger.Warn("No matching notification providers available",
			zap.String("notification_id", notification.ID),
			zap.Strings("requested_providers", notification.Providers),
		)
		s.sendResult(&NotificationResult{
			NotificationID: notification.ID,
			ProviderName:   "none",
			Success:        false,
			Error:          "no matching providers available",
			Timestamp:      time.Now(),
		})
		return
	}

	// Send notification to each provider
	for _, provider := range providers {
		result, err := provider.Send(s.ctx, notification)
		if err != nil {
			s.logger.Error("Failed to send notification",
				zap.String("notification_id", notification.ID),
				zap.String("provider", provider.GetName()),
				zap.Error(err),
			)
			s.sendResult(&NotificationResult{
				NotificationID: notification.ID,
				ProviderName:   provider.GetName(),
				Success:        false,
				Error:          err.Error(),
				Timestamp:      time.Now(),
			})

			// Handle retries
			if notification.Retries < notification.MaxRetries {
				notification.Retries++
				s.logger.Info("Retrying notification",
					zap.String("notification_id", notification.ID),
					zap.String("provider", provider.GetName()),
					zap.Int("retry", notification.Retries),
					zap.Int("max_retries", notification.MaxRetries),
				)
				// Wait before retrying
				time.Sleep(time.Duration(notification.Retries) * time.Second)
				s.queue <- notification
			}
		} else {
			s.sendResult(result)
		}
	}
}

// sendResult sends a notification result
func (s *NotificationService) sendResult(result *NotificationResult) {
	select {
	case s.results <- result:
	default:
		s.logger.Warn("Results channel is full, dropping result",
			zap.String("notification_id", result.NotificationID),
			zap.String("provider", result.ProviderName),
		)
	}
}

// isInitialized checks if the service is initialized
func (s *NotificationService) isInitialized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.initialized
}

// getTemplate gets a template by ID
func (s *NotificationService) getTemplate(id string) (*NotificationTemplate, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	template, ok := s.templates[id]
	return template, ok
}

// parseTemplate parses a template string
func (s *NotificationService) parseTemplate(templateStr string) (*template.Template, error) {
	return template.New("notification").Parse(templateStr)
}
