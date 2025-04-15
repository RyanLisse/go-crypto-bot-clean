package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
	"github.com/rs/zerolog"
)

// CredentialEncryptionService handles encryption and decryption of API credentials
type CredentialEncryptionService struct {
	encryptionService crypto.EncryptionService
	credentialRepo    port.APICredentialRepository
	logger            *zerolog.Logger
}

// NewCredentialEncryptionService creates a new CredentialEncryptionService
func NewCredentialEncryptionService(
	encryptionService crypto.EncryptionService,
	credentialRepo port.APICredentialRepository,
	logger *zerolog.Logger,
) *CredentialEncryptionService {
	return &CredentialEncryptionService{
		encryptionService: encryptionService,
		credentialRepo:    credentialRepo,
		logger:            logger,
	}
}

// EncryptAndSaveCredential encrypts and saves an API credential
func (s *CredentialEncryptionService) EncryptAndSaveCredential(ctx context.Context, credential *model.APICredential) error {
	// Encrypt the API secret
	encryptedSecret, err := s.encryptionService.Encrypt(credential.APISecret)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", credential.UserID).Str("exchange", credential.Exchange).Msg("Failed to encrypt API secret")
		return fmt.Errorf("failed to encrypt API secret: %w", err)
	}

	// Create a copy of the credential with the encrypted secret
	encryptedCredential := &model.APICredential{
		ID:           credential.ID,
		UserID:       credential.UserID,
		Exchange:     credential.Exchange,
		APIKey:       credential.APIKey,
		APISecret:    string(encryptedSecret), // Store the encrypted secret
		Label:        credential.Label,
		Status:       credential.Status,
		LastUsed:     credential.LastUsed,
		LastVerified: credential.LastVerified,
		ExpiresAt:    credential.ExpiresAt,
		RotationDue:  credential.RotationDue,
		FailureCount: credential.FailureCount,
		Metadata:     credential.Metadata,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	}

	// Save the encrypted credential
	if err := s.credentialRepo.Save(ctx, encryptedCredential); err != nil {
		s.logger.Error().Err(err).Str("userID", credential.UserID).Str("exchange", credential.Exchange).Msg("Failed to save encrypted credential")
		return fmt.Errorf("failed to save encrypted credential: %w", err)
	}

	return nil
}

// GetDecryptedCredential retrieves and decrypts an API credential
func (s *CredentialEncryptionService) GetDecryptedCredential(ctx context.Context, id string) (*model.APICredential, error) {
	// Get the encrypted credential
	encryptedCredential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get credential")
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	// Decrypt the API secret
	decryptedSecret, err := s.decryptAPISecret(encryptedCredential.APISecret)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to decrypt API secret")
		return nil, fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	// Create a copy of the credential with the decrypted secret
	decryptedCredential := &model.APICredential{
		ID:           encryptedCredential.ID,
		UserID:       encryptedCredential.UserID,
		Exchange:     encryptedCredential.Exchange,
		APIKey:       encryptedCredential.APIKey,
		APISecret:    decryptedSecret, // Use the decrypted secret
		Label:        encryptedCredential.Label,
		Status:       encryptedCredential.Status,
		LastUsed:     encryptedCredential.LastUsed,
		LastVerified: encryptedCredential.LastVerified,
		ExpiresAt:    encryptedCredential.ExpiresAt,
		RotationDue:  encryptedCredential.RotationDue,
		FailureCount: encryptedCredential.FailureCount,
		Metadata:     encryptedCredential.Metadata,
		CreatedAt:    encryptedCredential.CreatedAt,
		UpdatedAt:    encryptedCredential.UpdatedAt,
	}

	// Update last used timestamp
	go s.updateLastUsed(context.Background(), id)

	return decryptedCredential, nil
}

// GetDecryptedCredentialByUserAndExchange retrieves and decrypts an API credential by user ID and exchange
func (s *CredentialEncryptionService) GetDecryptedCredentialByUserAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	// Get the encrypted credential
	encryptedCredential, err := s.credentialRepo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to get credential by user and exchange")
		return nil, fmt.Errorf("failed to get credential by user and exchange: %w", err)
	}

	// Decrypt the API secret
	decryptedSecret, err := s.decryptAPISecret(encryptedCredential.APISecret)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to decrypt API secret")
		return nil, fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	// Create a copy of the credential with the decrypted secret
	decryptedCredential := &model.APICredential{
		ID:           encryptedCredential.ID,
		UserID:       encryptedCredential.UserID,
		Exchange:     encryptedCredential.Exchange,
		APIKey:       encryptedCredential.APIKey,
		APISecret:    decryptedSecret, // Use the decrypted secret
		Label:        encryptedCredential.Label,
		Status:       encryptedCredential.Status,
		LastUsed:     encryptedCredential.LastUsed,
		LastVerified: encryptedCredential.LastVerified,
		ExpiresAt:    encryptedCredential.ExpiresAt,
		RotationDue:  encryptedCredential.RotationDue,
		FailureCount: encryptedCredential.FailureCount,
		Metadata:     encryptedCredential.Metadata,
		CreatedAt:    encryptedCredential.CreatedAt,
		UpdatedAt:    encryptedCredential.UpdatedAt,
	}

	// Update last used timestamp
	go s.updateLastUsed(context.Background(), encryptedCredential.ID)

	return decryptedCredential, nil
}

// GetDecryptedCredentialByUserAndLabel retrieves and decrypts an API credential by user ID, exchange, and label
func (s *CredentialEncryptionService) GetDecryptedCredentialByUserAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error) {
	// Get the encrypted credential
	encryptedCredential, err := s.credentialRepo.GetByUserIDAndLabel(ctx, userID, exchange, label)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Str("label", label).Msg("Failed to get credential by user, exchange, and label")
		return nil, fmt.Errorf("failed to get credential by user, exchange, and label: %w", err)
	}

	// Decrypt the API secret
	decryptedSecret, err := s.decryptAPISecret(encryptedCredential.APISecret)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Str("label", label).Msg("Failed to decrypt API secret")
		return nil, fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	// Create a copy of the credential with the decrypted secret
	decryptedCredential := &model.APICredential{
		ID:           encryptedCredential.ID,
		UserID:       encryptedCredential.UserID,
		Exchange:     encryptedCredential.Exchange,
		APIKey:       encryptedCredential.APIKey,
		APISecret:    decryptedSecret, // Use the decrypted secret
		Label:        encryptedCredential.Label,
		Status:       encryptedCredential.Status,
		LastUsed:     encryptedCredential.LastUsed,
		LastVerified: encryptedCredential.LastVerified,
		ExpiresAt:    encryptedCredential.ExpiresAt,
		RotationDue:  encryptedCredential.RotationDue,
		FailureCount: encryptedCredential.FailureCount,
		Metadata:     encryptedCredential.Metadata,
		CreatedAt:    encryptedCredential.CreatedAt,
		UpdatedAt:    encryptedCredential.UpdatedAt,
	}

	// Update last used timestamp
	go s.updateLastUsed(context.Background(), encryptedCredential.ID)

	return decryptedCredential, nil
}

// ListDecryptedCredentials retrieves and decrypts all API credentials for a user
func (s *CredentialEncryptionService) ListDecryptedCredentials(ctx context.Context, userID string) ([]*model.APICredential, error) {
	// Get the encrypted credentials
	encryptedCredentials, err := s.credentialRepo.ListByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to list credentials")
		return nil, fmt.Errorf("failed to list credentials: %w", err)
	}

	// Decrypt each credential
	decryptedCredentials := make([]*model.APICredential, 0, len(encryptedCredentials))
	for _, encryptedCredential := range encryptedCredentials {
		// Decrypt the API secret
		decryptedSecret, err := s.decryptAPISecret(encryptedCredential.APISecret)
		if err != nil {
			s.logger.Error().Err(err).Str("id", encryptedCredential.ID).Msg("Failed to decrypt API secret")
			// Skip this credential and continue with the others
			continue
		}

		// Create a copy of the credential with the decrypted secret
		decryptedCredential := &model.APICredential{
			ID:           encryptedCredential.ID,
			UserID:       encryptedCredential.UserID,
			Exchange:     encryptedCredential.Exchange,
			APIKey:       encryptedCredential.APIKey,
			APISecret:    decryptedSecret, // Use the decrypted secret
			Label:        encryptedCredential.Label,
			Status:       encryptedCredential.Status,
			LastUsed:     encryptedCredential.LastUsed,
			LastVerified: encryptedCredential.LastVerified,
			ExpiresAt:    encryptedCredential.ExpiresAt,
			RotationDue:  encryptedCredential.RotationDue,
			FailureCount: encryptedCredential.FailureCount,
			Metadata:     encryptedCredential.Metadata,
			CreatedAt:    encryptedCredential.CreatedAt,
			UpdatedAt:    encryptedCredential.UpdatedAt,
		}

		decryptedCredentials = append(decryptedCredentials, decryptedCredential)
	}

	return decryptedCredentials, nil
}

// DeleteCredential deletes an API credential
func (s *CredentialEncryptionService) DeleteCredential(ctx context.Context, id string) error {
	// Delete the credential
	if err := s.credentialRepo.DeleteByID(ctx, id); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to delete credential")
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	return nil
}

// VerifyCredential verifies an API credential by checking if it can be decrypted
func (s *CredentialEncryptionService) VerifyCredential(ctx context.Context, id string) error {
	// Get the encrypted credential
	encryptedCredential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get credential")
		return fmt.Errorf("failed to get credential: %w", err)
	}

	// Try to decrypt the API secret
	_, err = s.decryptAPISecret(encryptedCredential.APISecret)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to decrypt API secret")
		
		// Increment failure count
		if err := s.credentialRepo.IncrementFailureCount(ctx, id); err != nil {
			s.logger.Error().Err(err).Str("id", id).Msg("Failed to increment failure count")
		}
		
		// Update status to failed if failure count exceeds threshold
		if encryptedCredential.FailureCount >= 5 {
			if err := s.credentialRepo.UpdateStatus(ctx, id, model.APICredentialStatusFailed); err != nil {
				s.logger.Error().Err(err).Str("id", id).Msg("Failed to update credential status")
			}
		}
		
		return fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	// Update last verified timestamp and reset failure count
	now := time.Now()
	if err := s.credentialRepo.UpdateLastVerified(ctx, id, now); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to update last verified timestamp")
	}
	
	if err := s.credentialRepo.ResetFailureCount(ctx, id); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to reset failure count")
	}

	return nil
}

// UpdateCredentialStatus updates the status of an API credential
func (s *CredentialEncryptionService) UpdateCredentialStatus(ctx context.Context, id string, status model.APICredentialStatus) error {
	// Update the status
	if err := s.credentialRepo.UpdateStatus(ctx, id, status); err != nil {
		s.logger.Error().Err(err).Str("id", id).Str("status", string(status)).Msg("Failed to update credential status")
		return fmt.Errorf("failed to update credential status: %w", err)
	}

	return nil
}

// decryptAPISecret decrypts an API secret
func (s *CredentialEncryptionService) decryptAPISecret(encryptedSecret string) (string, error) {
	// Try to decrypt as bytes first
	decryptedSecret, err := s.encryptionService.Decrypt([]byte(encryptedSecret))
	if err != nil {
		// If that fails, try to decrypt as a string (for backward compatibility)
		return "", fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	return decryptedSecret, nil
}

// updateLastUsed updates the last used timestamp of an API credential
func (s *CredentialEncryptionService) updateLastUsed(ctx context.Context, id string) {
	now := time.Now()
	if err := s.credentialRepo.UpdateLastUsed(ctx, id, now); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to update last used timestamp")
	}
}
