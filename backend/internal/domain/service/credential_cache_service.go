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

// CachedCredential represents a cached API credential
type CachedCredential struct {
	Credential  *model.APICredential
	ExpiresAt   time.Time
	LastUpdated time.Time
}

// CredentialCacheService handles caching of API credentials
type CredentialCacheService struct {
	credentialRepo port.APICredentialRepository
	cache          map[string]*CachedCredential // Key: "userID:exchange"
	cacheMutex     sync.RWMutex
	ttl            time.Duration
	logger         *zerolog.Logger
}

// NewCredentialCacheService creates a new CredentialCacheService
func NewCredentialCacheService(
	credentialRepo port.APICredentialRepository,
	ttl time.Duration,
	logger *zerolog.Logger,
) *CredentialCacheService {
	service := &CredentialCacheService{
		credentialRepo: credentialRepo,
		cache:          make(map[string]*CachedCredential),
		ttl:            ttl,
		logger:         logger,
	}

	// Start a background goroutine to clean up expired cache entries
	go service.startCleanupTask()

	return service
}

// GetCredential gets a credential from the cache or repository
func (s *CredentialCacheService) GetCredential(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	// Generate cache key
	cacheKey := fmt.Sprintf("%s:%s", userID, exchange)

	// Try to get from cache first
	s.cacheMutex.RLock()
	cachedCred, found := s.cache[cacheKey]
	s.cacheMutex.RUnlock()

	// If found in cache and not expired, return it
	if found && time.Now().Before(cachedCred.ExpiresAt) {
		s.logger.Debug().Str("userID", userID).Str("exchange", exchange).Msg("Credential cache hit")
		return cachedCred.Credential, nil
	}

	// Not found in cache or expired, get from repository
	s.logger.Debug().Str("userID", userID).Str("exchange", exchange).Msg("Credential cache miss")
	credential, err := s.credentialRepo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil {
		return nil, err
	}

	// Cache the credential
	s.cacheCredential(cacheKey, credential)

	return credential, nil
}

// GetCredentialByID gets a credential by ID from the cache or repository
func (s *CredentialCacheService) GetCredentialByID(ctx context.Context, id string) (*model.APICredential, error) {
	// Try to find in cache first
	s.cacheMutex.RLock()
	for _, cachedCred := range s.cache {
		if cachedCred.Credential.ID == id && time.Now().Before(cachedCred.ExpiresAt) {
			s.cacheMutex.RUnlock()
			s.logger.Debug().Str("id", id).Msg("Credential cache hit by ID")
			return cachedCred.Credential, nil
		}
	}
	s.cacheMutex.RUnlock()

	// Not found in cache or expired, get from repository
	s.logger.Debug().Str("id", id).Msg("Credential cache miss by ID")
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the credential
	cacheKey := fmt.Sprintf("%s:%s", credential.UserID, credential.Exchange)
	s.cacheCredential(cacheKey, credential)

	return credential, nil
}

// InvalidateCredential invalidates a cached credential
func (s *CredentialCacheService) InvalidateCredential(userID, exchange string) {
	cacheKey := fmt.Sprintf("%s:%s", userID, exchange)
	s.cacheMutex.Lock()
	delete(s.cache, cacheKey)
	s.cacheMutex.Unlock()
	s.logger.Debug().Str("userID", userID).Str("exchange", exchange).Msg("Credential cache invalidated")
}

// InvalidateCredentialByID invalidates a cached credential by ID
func (s *CredentialCacheService) InvalidateCredentialByID(id string) {
	s.cacheMutex.Lock()
	for key, cachedCred := range s.cache {
		if cachedCred.Credential.ID == id {
			delete(s.cache, key)
			s.logger.Debug().Str("id", id).Msg("Credential cache invalidated by ID")
			break
		}
	}
	s.cacheMutex.Unlock()
}

// InvalidateAllCredentials invalidates all cached credentials
func (s *CredentialCacheService) InvalidateAllCredentials() {
	s.cacheMutex.Lock()
	s.cache = make(map[string]*CachedCredential)
	s.cacheMutex.Unlock()
	s.logger.Debug().Msg("All credential cache invalidated")
}

// cacheCredential caches a credential
func (s *CredentialCacheService) cacheCredential(key string, credential *model.APICredential) {
	// Don't cache credentials with certain statuses
	if credential.Status == model.APICredentialStatusRevoked ||
		credential.Status == model.APICredentialStatusExpired ||
		credential.Status == model.APICredentialStatusFailed {
		return
	}

	// Calculate cache expiration time
	expiresAt := time.Now().Add(s.ttl)

	// If the credential has an expiration date, use the earlier of the two
	if credential.ExpiresAt != nil && credential.ExpiresAt.Before(expiresAt) {
		expiresAt = *credential.ExpiresAt
	}

	// Create cached credential
	cachedCred := &CachedCredential{
		Credential:  credential,
		ExpiresAt:   expiresAt,
		LastUpdated: time.Now(),
	}

	// Store in cache
	s.cacheMutex.Lock()
	s.cache[key] = cachedCred
	s.cacheMutex.Unlock()

	s.logger.Debug().Str("key", key).Time("expiresAt", expiresAt).Msg("Credential cached")
}

// startCleanupTask starts a background task to clean up expired cache entries
func (s *CredentialCacheService) startCleanupTask() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredEntries()
	}
}

// cleanupExpiredEntries removes expired entries from the cache
func (s *CredentialCacheService) cleanupExpiredEntries() {
	now := time.Now()
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for key, cachedCred := range s.cache {
		if now.After(cachedCred.ExpiresAt) {
			delete(s.cache, key)
			s.logger.Debug().Str("key", key).Msg("Expired credential removed from cache")
		}
	}
}

// GetCacheStats returns statistics about the cache
func (s *CredentialCacheService) GetCacheStats() map[string]interface{} {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	stats := map[string]interface{}{
		"size":            len(s.cache),
		"ttl_seconds":     s.ttl.Seconds(),
		"active_entries":  0,
		"expired_entries": 0,
	}

	now := time.Now()
	for _, cachedCred := range s.cache {
		if now.Before(cachedCred.ExpiresAt) {
			stats["active_entries"] = stats["active_entries"].(int) + 1
		} else {
			stats["expired_entries"] = stats["expired_entries"].(int) + 1
		}
	}

	return stats
}
