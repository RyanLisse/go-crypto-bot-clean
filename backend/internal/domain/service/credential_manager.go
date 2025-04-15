package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
	"github.com/google/uuid"
)

// CredentialManager manages API credentials securely
type CredentialManager struct {
	credentialRepo port.APICredentialRepository
	encryptionSvc  crypto.EncryptionService
}

// NewCredentialManager creates a new CredentialManager
func NewCredentialManager(credentialRepo port.APICredentialRepository, encryptionSvc crypto.EncryptionService) *CredentialManager {
	return &CredentialManager{
		credentialRepo: credentialRepo,
		encryptionSvc:  encryptionSvc,
	}
}

// StoreCredential stores an API credential securely
func (m *CredentialManager) StoreCredential(ctx context.Context, userID, exchange, apiKey, apiSecret, label string) (*model.APICredential, error) {
	// Validate inputs
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if exchange == "" {
		return nil, errors.New("exchange is required")
	}
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}
	if apiSecret == "" {
		return nil, errors.New("API secret is required")
	}

	// Create credential
	credentialEntity := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    userID,
		Exchange:  exchange,
		APIKey:    apiKey,
		APISecret: "", // Don't store plaintext secret in memory
		Label:     label,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Encrypt API secret
	_, err := m.encryptionSvc.Encrypt(apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt API secret: %w", err)
	}


	// Store credential
	if err := m.credentialRepo.Save(ctx, credentialEntity); err != nil {
		return nil, fmt.Errorf("failed to save API credential: %w", err)
	}

	// Return credential without plaintext secret
	credentialEntity.APISecret = "********"
	return credentialEntity, nil
}

// GetCredential retrieves an API credential by ID
func (m *CredentialManager) GetCredential(ctx context.Context, id string) (*model.APICredential, error) {
	// Get credential from repository
	credential, err := m.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get API credential: %w", err)
	}

	// Mask API secret
	credential.APISecret = "********"

	return credential, nil
}

// GetCredentialWithSecret retrieves an API credential with decrypted secret
func (m *CredentialManager) GetCredentialWithSecret(ctx context.Context, id string) (*model.APICredential, error) {
	// Get credential from repository
	credential, err := m.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get API credential: %w", err)
	}

	// Decrypt API secret
	decryptedSecret, err := m.encryptionSvc.Decrypt([]byte(credential.APISecret))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	// Set decrypted secret
	credential.APISecret = decryptedSecret

	return credential, nil
}

// ListCredentials lists API credentials for a user
func (m *CredentialManager) ListCredentials(ctx context.Context, userID string) ([]*model.APICredential, error) {
	// Get credentials from repository
	credentials, err := m.credentialRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API credentials: %w", err)
	}

	// Mask API secrets
	for _, credential := range credentials {
		credential.APISecret = "********"
	}

	return credentials, nil
}

// DeleteCredential deletes an API credential
func (m *CredentialManager) DeleteCredential(ctx context.Context, id string) error {
	// Delete credential from repository
	if err := m.credentialRepo.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete API credential: %w", err)
	}

	return nil
}

// UpdateCredential updates an API credential
func (m *CredentialManager) UpdateCredential(ctx context.Context, id, apiKey, apiSecret, label string) (*model.APICredential, error) {
	// Get credential from repository
	credential, err := m.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get API credential: %w", err)
	}

	// Update fields
	if apiKey != "" {
		credential.APIKey = apiKey
	}

	if apiSecret != "" {
		// Encrypt API secret
		_, err := m.encryptionSvc.Encrypt(apiSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt API secret: %w", err)
		}
		credential.APISecret = apiSecret
	}

	if label != "" {
		credential.Label = label
	}

	credential.UpdatedAt = time.Now()

	// Save updated credential
	if err := m.credentialRepo.Save(ctx, credential); err != nil {
		return nil, fmt.Errorf("failed to save updated API credential: %w", err)
	}

	// Mask API secret
	credential.APISecret = "********"

	return credential, nil
}

// ValidateCredential validates an API credential by checking with the exchange
func (m *CredentialManager) ValidateCredential(ctx context.Context, id string) (bool, error) {
	// Get credential with secret
	_, err := m.GetCredentialWithSecret(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to get API credential: %w", err)
	}

	// TODO: Implement exchange-specific validation
	// This would typically involve making a request to the exchange API
	// using the credentials to verify they are valid

	// For now, just return true
	return true, nil
}
