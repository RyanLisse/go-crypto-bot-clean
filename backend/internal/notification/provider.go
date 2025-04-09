package notification

import (
	"context"
	"errors"
	"sync"
)

// Common errors
var (
	ErrProviderNotInitialized = errors.New("notification provider not initialized")
	ErrProviderNotAvailable   = errors.New("notification provider not available")
	ErrInvalidConfiguration   = errors.New("invalid provider configuration")
	ErrSendFailed             = errors.New("failed to send notification")
	ErrRateLimited            = errors.New("rate limit exceeded")
)

// NotificationProvider defines the interface for notification providers
type NotificationProvider interface {
	// Initialize sets up the provider with configuration
	Initialize(config map[string]interface{}) error

	// Send sends a notification
	Send(ctx context.Context, notification *Notification) (*NotificationResult, error)

	// GetName returns the provider name
	GetName() string

	// IsAvailable checks if the provider is available
	IsAvailable() bool
}

// ProviderRegistry manages notification providers
type ProviderRegistry struct {
	providers map[string]NotificationProvider
	mu        sync.RWMutex
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]NotificationProvider),
	}
}

// Register registers a notification provider
func (r *ProviderRegistry) Register(provider NotificationProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[provider.GetName()] = provider
}

// Get returns a notification provider by name
func (r *ProviderRegistry) Get(name string) (NotificationProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[name]
	return provider, ok
}

// GetAll returns all registered notification providers
func (r *ProviderRegistry) GetAll() []NotificationProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	providers := make([]NotificationProvider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}

// GetAvailable returns all available notification providers
func (r *ProviderRegistry) GetAvailable() []NotificationProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	providers := make([]NotificationProvider, 0, len(r.providers))
	for _, provider := range r.providers {
		if provider.IsAvailable() {
			providers = append(providers, provider)
		}
	}
	return providers
}

// BaseProvider provides common functionality for notification providers
type BaseProvider struct {
	name      string
	available bool
	config    map[string]interface{}
	mu        sync.RWMutex
}

// NewBaseProvider creates a new base provider
func NewBaseProvider(name string) *BaseProvider {
	return &BaseProvider{
		name:      name,
		available: false,
	}
}

// GetName returns the provider name
func (p *BaseProvider) GetName() string {
	return p.name
}

// IsAvailable checks if the provider is available
func (p *BaseProvider) IsAvailable() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.available
}

// SetAvailable sets the availability of the provider
func (p *BaseProvider) SetAvailable(available bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.available = available
}

// GetConfig returns the provider configuration
func (p *BaseProvider) GetConfig() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.config
}

// SetConfig sets the provider configuration
func (p *BaseProvider) SetConfig(config map[string]interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = config
}

// GetConfigString gets a string value from the configuration
func (p *BaseProvider) GetConfigString(key string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.config == nil {
		return "", false
	}
	value, ok := p.config[key]
	if !ok {
		return "", false
	}
	strValue, ok := value.(string)
	return strValue, ok
}

// GetConfigInt gets an int value from the configuration
func (p *BaseProvider) GetConfigInt(key string) (int, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.config == nil {
		return 0, false
	}
	value, ok := p.config[key]
	if !ok {
		return 0, false
	}
	intValue, ok := value.(int)
	return intValue, ok
}

// GetConfigBool gets a bool value from the configuration
func (p *BaseProvider) GetConfigBool(key string) (bool, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.config == nil {
		return false, false
	}
	value, ok := p.config[key]
	if !ok {
		return false, false
	}
	boolValue, ok := value.(bool)
	return boolValue, ok
}

// GetConfigStringSlice gets a string slice from the configuration
func (p *BaseProvider) GetConfigStringSlice(key string) ([]string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.config == nil {
		return nil, false
	}
	value, ok := p.config[key]
	if !ok {
		return nil, false
	}
	
	// Try to convert from []interface{} to []string
	if interfaceSlice, ok := value.([]interface{}); ok {
		strSlice := make([]string, len(interfaceSlice))
		for i, v := range interfaceSlice {
			if strValue, ok := v.(string); ok {
				strSlice[i] = strValue
			} else {
				return nil, false
			}
		}
		return strSlice, true
	}
	
	// Direct []string
	strSlice, ok := value.([]string)
	return strSlice, ok
}
