package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// FallbackStrategy represents a strategy for credential fallback
type FallbackStrategy string

const (
	// FallbackStrategyNone indicates no fallback
	FallbackStrategyNone FallbackStrategy = "none"

	// FallbackStrategyDefault indicates to use a default credential
	FallbackStrategyDefault FallbackStrategy = "default"

	// FallbackStrategyPool indicates to use a pool of credentials
	FallbackStrategyPool FallbackStrategy = "pool"

	// FallbackStrategyReadOnly indicates to use read-only credentials
	FallbackStrategyReadOnly FallbackStrategy = "read_only"
)

// FallbackConfig represents the configuration for credential fallback
type FallbackConfig struct {
	Strategy      FallbackStrategy
	DefaultCredID string
	PoolCredIDs   []string
	ReadOnlyMode  bool
}

// CredentialFallbackService handles fallback mechanisms for API credentials
type CredentialFallbackService struct {
	credentialRepo   port.APICredentialRepository
	fallbackConfigs  map[string]*FallbackConfig // Key: exchange
	configMutex      sync.RWMutex
	logger           *zerolog.Logger
	lastFailureTime  map[string]time.Time // Key: "userID:exchange"
	failureCounters  map[string]int       // Key: "userID:exchange"
	failureMutex     sync.RWMutex
	failureThreshold int
	cooldownPeriod   time.Duration
}

// NewCredentialFallbackService creates a new CredentialFallbackService
func NewCredentialFallbackService(
	credentialRepo port.APICredentialRepository,
	logger *zerolog.Logger,
	failureThreshold int,
	cooldownPeriod time.Duration,
) *CredentialFallbackService {
	service := &CredentialFallbackService{
		credentialRepo:   credentialRepo,
		fallbackConfigs:  make(map[string]*FallbackConfig),
		lastFailureTime:  make(map[string]time.Time),
		failureCounters:  make(map[string]int),
		logger:           logger,
		failureThreshold: failureThreshold,
		cooldownPeriod:   cooldownPeriod,
	}

	// Initialize fallback configurations for known exchanges
	service.initializeFallbackConfigs()

	// Start a background goroutine to clean up failure counters
	go service.startCleanupTask()

	return service
}

// GetCredentialWithFallback gets a credential with fallback mechanisms
func (s *CredentialFallbackService) GetCredentialWithFallback(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	// Generate key for failure tracking
	failureKey := fmt.Sprintf("%s:%s", userID, exchange)

	// Check if we should use fallback based on failure history
	useFallback := s.shouldUseFallback(failureKey)

	// If we should use fallback, get the fallback credential
	if useFallback {
		s.logger.Info().Str("userID", userID).Str("exchange", exchange).Msg("Using fallback credential due to failure history")
		return s.getFallbackCredential(ctx, exchange)
	}

	// Try to get the user's credential
	credential, err := s.credentialRepo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil {
		// Record the failure
		s.recordFailure(failureKey)

		// Try fallback if this failure puts us over the threshold
		if s.shouldUseFallback(failureKey) {
			s.logger.Info().Str("userID", userID).Str("exchange", exchange).Msg("Using fallback credential after failure")
			return s.getFallbackCredential(ctx, exchange)
		}

		return nil, err
	}

	// Reset failure counter on success
	s.resetFailureCounter(failureKey)

	return credential, nil
}

// SetFallbackConfig sets the fallback configuration for an exchange
func (s *CredentialFallbackService) SetFallbackConfig(exchange string, config *FallbackConfig) {
	s.configMutex.Lock()
	defer s.configMutex.Unlock()

	s.fallbackConfigs[exchange] = config
	s.logger.Info().Str("exchange", exchange).Str("strategy", string(config.Strategy)).Msg("Fallback configuration updated")
}

// GetFallbackConfig gets the fallback configuration for an exchange
func (s *CredentialFallbackService) GetFallbackConfig(exchange string) (*FallbackConfig, error) {
	s.configMutex.RLock()
	defer s.configMutex.RUnlock()

	config, found := s.fallbackConfigs[exchange]
	if !found {
		return nil, fmt.Errorf("no fallback configuration found for exchange: %s", exchange)
	}

	return config, nil
}

// ResetFailureCounters resets all failure counters
func (s *CredentialFallbackService) ResetFailureCounters() {
	s.failureMutex.Lock()
	defer s.failureMutex.Unlock()

	s.failureCounters = make(map[string]int)
	s.lastFailureTime = make(map[string]time.Time)
	s.logger.Info().Msg("All failure counters reset")
}

// initializeFallbackConfigs initializes fallback configurations for known exchanges
func (s *CredentialFallbackService) initializeFallbackConfigs() {
	// Initialize configurations for known exchanges
	exchanges := []string{"mexc", "binance", "coinbase", "kraken"}
	for _, exchange := range exchanges {
		config := &FallbackConfig{
			Strategy:      FallbackStrategyNone,
			DefaultCredID: "",
			PoolCredIDs:   []string{},
			ReadOnlyMode:  false,
		}
		s.fallbackConfigs[exchange] = config
	}
}

// shouldUseFallback checks if we should use fallback based on failure history
func (s *CredentialFallbackService) shouldUseFallback(key string) bool {
	s.failureMutex.RLock()
	defer s.failureMutex.RUnlock()

	// Check if we have a failure counter for this key
	count, found := s.failureCounters[key]
	if !found {
		return false
	}

	// Check if we're in the cooldown period
	lastFailure, found := s.lastFailureTime[key]
	if found && time.Since(lastFailure) > s.cooldownPeriod {
		// Reset counter if cooldown period has passed
		return false
	}

	// Use fallback if failure count exceeds threshold
	return count >= s.failureThreshold
}

// recordFailure records a failure for a key
func (s *CredentialFallbackService) recordFailure(key string) {
	s.failureMutex.Lock()
	defer s.failureMutex.Unlock()

	// Increment failure counter
	s.failureCounters[key]++
	s.lastFailureTime[key] = time.Now()

	s.logger.Debug().Str("key", key).Int("count", s.failureCounters[key]).Msg("Failure recorded")
}

// resetFailureCounter resets the failure counter for a key
func (s *CredentialFallbackService) resetFailureCounter(key string) {
	s.failureMutex.Lock()
	defer s.failureMutex.Unlock()

	delete(s.failureCounters, key)
	delete(s.lastFailureTime, key)

	s.logger.Debug().Str("key", key).Msg("Failure counter reset")
}

// startCleanupTask starts a background task to clean up old failure counters
func (s *CredentialFallbackService) startCleanupTask() {
	ticker := time.NewTicker(s.cooldownPeriod)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupOldFailures()
	}
}

// cleanupOldFailures removes old failure counters
func (s *CredentialFallbackService) cleanupOldFailures() {
	now := time.Now()
	s.failureMutex.Lock()
	defer s.failureMutex.Unlock()

	for key, lastFailure := range s.lastFailureTime {
		if now.Sub(lastFailure) > s.cooldownPeriod {
			delete(s.failureCounters, key)
			delete(s.lastFailureTime, key)
			s.logger.Debug().Str("key", key).Msg("Old failure counter removed")
		}
	}
}

// getFallbackCredential gets a fallback credential for an exchange
func (s *CredentialFallbackService) getFallbackCredential(ctx context.Context, exchange string) (*model.APICredential, error) {
	// Get the fallback configuration
	s.configMutex.RLock()
	config, found := s.fallbackConfigs[exchange]
	s.configMutex.RUnlock()

	if !found || config.Strategy == FallbackStrategyNone {
		return nil, fmt.Errorf("no fallback configuration found for exchange: %s", exchange)
	}

	// Get the fallback credential based on the strategy
	switch config.Strategy {
	case FallbackStrategyDefault:
		return s.getDefaultCredential(ctx, exchange, config)
	case FallbackStrategyPool:
		return s.getPoolCredential(ctx, exchange, config)
	case FallbackStrategyReadOnly:
		return s.getReadOnlyCredential(ctx, exchange, config)
	default:
		return nil, fmt.Errorf("unknown fallback strategy: %s", config.Strategy)
	}
}

// getDefaultCredential gets the default credential for an exchange
func (s *CredentialFallbackService) getDefaultCredential(ctx context.Context, exchange string, config *FallbackConfig) (*model.APICredential, error) {
	if config.DefaultCredID == "" {
		return nil, fmt.Errorf("no default credential ID configured for exchange: %s", exchange)
	}

	// Get the default credential
	credential, err := s.credentialRepo.GetByID(ctx, config.DefaultCredID)
	if err != nil {
		return nil, fmt.Errorf("failed to get default credential: %w", err)
	}

	// Mark the credential as a fallback
	if credential.Metadata == nil {
		credential.Metadata = &model.APICredentialMetadata{}
	}
	credential.Metadata.Custom = map[string]string{
		"fallback": "true",
		"strategy": string(FallbackStrategyDefault),
	}

	return credential, nil
}

// getPoolCredential gets a credential from the pool for an exchange
func (s *CredentialFallbackService) getPoolCredential(ctx context.Context, exchange string, config *FallbackConfig) (*model.APICredential, error) {
	if len(config.PoolCredIDs) == 0 {
		return nil, fmt.Errorf("no pool credential IDs configured for exchange: %s", exchange)
	}

	// Try each credential in the pool
	var lastErr error
	for _, credID := range config.PoolCredIDs {
		credential, err := s.credentialRepo.GetByID(ctx, credID)
		if err != nil {
			lastErr = err
			continue
		}

		// Mark the credential as a fallback
		if credential.Metadata == nil {
			credential.Metadata = &model.APICredentialMetadata{}
		}
		credential.Metadata.Custom = map[string]string{
			"fallback": "true",
			"strategy": string(FallbackStrategyPool),
		}

		return credential, nil
	}

	return nil, fmt.Errorf("failed to get any pool credential: %w", lastErr)
}

// getReadOnlyCredential gets a read-only credential for an exchange
func (s *CredentialFallbackService) getReadOnlyCredential(ctx context.Context, exchange string, config *FallbackConfig) (*model.APICredential, error) {
	// In a real implementation, we would have read-only credentials for each exchange
	// For now, we'll just return an error
	return nil, fmt.Errorf("read-only credentials not implemented for exchange: %s", exchange)
}
