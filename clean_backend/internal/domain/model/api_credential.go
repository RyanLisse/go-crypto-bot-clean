package model

import (
	"time"

	"github.com/google/uuid"
)

// APICredentialStatus represents the status of an API credential
type APICredentialStatus string

const (
	// APICredentialStatusActive indicates the credential is active and can be used
	APICredentialStatusActive APICredentialStatus = "active"

	// APICredentialStatusInactive indicates the credential is inactive and should not be used
	APICredentialStatusInactive APICredentialStatus = "inactive"

	// APICredentialStatusRevoked indicates the credential has been revoked and cannot be used
	APICredentialStatusRevoked APICredentialStatus = "revoked"

	// APICredentialStatusExpired indicates the credential has expired and cannot be used
	APICredentialStatusExpired APICredentialStatus = "expired"

	// APICredentialStatusFailed indicates the credential has failed authentication and should be verified
	APICredentialStatusFailed APICredentialStatus = "failed"
)

// APICredentialMetadata contains additional metadata for an API credential
type APICredentialMetadata struct {
	Permissions  []string          `json:"permissions,omitempty"`   // Permissions granted to this credential
	IPWhitelist  []string          `json:"ip_whitelist,omitempty"`  // IP addresses allowed to use this credential
	RateLimits   map[string]int    `json:"rate_limits,omitempty"`   // Rate limits for this credential
	Tags         []string          `json:"tags,omitempty"`          // Tags for categorizing credentials
	Description  string            `json:"description,omitempty"`   // Description of the credential
	ContactEmail string            `json:"contact_email,omitempty"` // Contact email for notifications
	Custom       map[string]string `json:"custom,omitempty"`        // Custom metadata
}

// APICredential represents a user's API credentials for a cryptocurrency exchange
type APICredential struct {
	ID                  string
	UserID              string
	Exchange            string
	APIKey              string
	APISecret           string // Encrypted
	APISecretKeyVersion string // Version of key used for encryption
	Label               string
	Status              APICredentialStatus
	LastUsed            *time.Time
	LastVerified        *time.Time
	ExpiresAt           *time.Time
	RotationDue         *time.Time
	FailureCount        int
	Metadata            *APICredentialMetadata
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewAPICredential creates a new APICredential
func NewAPICredential(userID, exchange, apiKey, apiSecret, label string) *APICredential {
	now := time.Now()
	return &APICredential{
		ID:           uuid.New().String(),
		UserID:       userID,
		Exchange:     exchange,
		APIKey:       apiKey,
		APISecret:    apiSecret,
		Label:        label,
		Status:       APICredentialStatusActive,
		FailureCount: 0,
		Metadata:     &APICredentialMetadata{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Validate validates the APICredential
func (c *APICredential) Validate() error {
	if c.UserID == "" {
		return ErrInvalidUserID
	}

	if c.Exchange == "" {
		return ErrInvalidExchange
	}

	if c.APIKey == "" {
		return ErrInvalidAPIKey
	}

	if c.APISecret == "" {
		return ErrInvalidAPISecret
	}

	return nil
}

// Update updates the APICredential
func (c *APICredential) Update(apiKey, apiSecret, label string) {
	if apiKey != "" {
		c.APIKey = apiKey
	}

	if apiSecret != "" {
		c.APISecret = apiSecret
	}

	if label != "" {
		c.Label = label
	}

	c.UpdatedAt = time.Now()
}
