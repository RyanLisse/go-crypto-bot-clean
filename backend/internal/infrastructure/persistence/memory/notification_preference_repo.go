package memory

import (
	"context"
	"sync"

	"go-crypto-bot-clean/backend/internal/domain/notification"
	"go-crypto-bot-clean/backend/internal/domain/notification/ports"
)

// InMemoryNotificationPreferenceRepository implements ports.NotificationPreferenceRepository using an in-memory map.
// Note: This is not thread-safe for concurrent writes without the mutex.
// It also implements SaveUserPreference for testing purposes.
type InMemoryNotificationPreferenceRepository struct {
	mu    sync.RWMutex
	prefs map[string][]notification.Preference // Map userID to list of preferences
}

// NewInMemoryNotificationPreferenceRepository creates a new in-memory repository.
func NewInMemoryNotificationPreferenceRepository() *InMemoryNotificationPreferenceRepository {
	return &InMemoryNotificationPreferenceRepository{
		prefs: make(map[string][]notification.Preference),
	}
}

// GetUserPreferences retrieves all active notification preferences for a given user.
func (r *InMemoryNotificationPreferenceRepository) GetUserPreferences(ctx context.Context, userID string) ([]notification.Preference, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userPrefs, exists := r.prefs[userID]
	if !exists {
		return []notification.Preference{}, nil // Return empty slice if user not found
	}

	// Filter for enabled preferences
	enabledPrefs := make([]notification.Preference, 0, len(userPrefs))
	for _, pref := range userPrefs {
		if pref.Enabled {
			enabledPrefs = append(enabledPrefs, pref)
		}
	}

	return enabledPrefs, nil
}

// SaveUserPreference saves or updates a user preference. // NOTE: Added for test seeding!
// This method is NOT part of the official ports.NotificationPreferenceRepository interface yet.
// It assumes a simple overwrite/append logic for demonstration.
func (r *InMemoryNotificationPreferenceRepository) SaveUserPreference(ctx context.Context, pref notification.Preference) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userPrefs := r.prefs[pref.UserID]

	// Simple logic: remove existing pref for the same channel, then append.
	// A real implementation might update in place or have more complex logic.
	updatedPrefs := make([]notification.Preference, 0, len(userPrefs)+1)
	found := false
	for _, existingPref := range userPrefs {
		if existingPref.Channel == pref.Channel {
			updatedPrefs = append(updatedPrefs, pref) // Replace with new one
			found = true
		} else {
			updatedPrefs = append(updatedPrefs, existingPref)
		}
	}
	if !found {
		updatedPrefs = append(updatedPrefs, pref) // Add if new channel
	}

	r.prefs[pref.UserID] = updatedPrefs
	return nil
}

// Ensure implementation satisfies the interface (compile-time check)
var _ ports.NotificationPreferenceRepository = (*InMemoryNotificationPreferenceRepository)(nil)

// We also satisfy the interface used for seeding in tests:
var _ interface {
	SaveUserPreference(context.Context, notification.Preference) error
	GetUserPreferences(context.Context, string) ([]notification.Preference, error)
} = (*InMemoryNotificationPreferenceRepository)(nil)
