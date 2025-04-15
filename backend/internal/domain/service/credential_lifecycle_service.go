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

// CredentialLifecycleService handles the lifecycle of API credentials
type CredentialLifecycleService struct {
	credentialRepo    port.APICredentialRepository
	encryptionService crypto.EncryptionService
	validationService *CredentialValidationService
	errorService      *CredentialErrorService
	loggingService    *CredentialLoggingService
	logger            *zerolog.Logger
}

// NewCredentialLifecycleService creates a new CredentialLifecycleService
func NewCredentialLifecycleService(
	credentialRepo port.APICredentialRepository,
	encryptionService crypto.EncryptionService,
	validationService *CredentialValidationService,
	errorService *CredentialErrorService,
	loggingService *CredentialLoggingService,
	logger *zerolog.Logger,
) *CredentialLifecycleService {
	return &CredentialLifecycleService{
		credentialRepo:    credentialRepo,
		encryptionService: encryptionService,
		validationService: validationService,
		errorService:      errorService,
		loggingService:    loggingService,
		logger:            logger,
	}
}

// CreateCredential creates a new API credential
func (s *CredentialLifecycleService) CreateCredential(ctx context.Context, userID, exchange, apiKey, apiSecret, label string, expiresIn *time.Duration) (*model.APICredential, error) {
	startTime := time.Now()

	// Create a new credential
	credential := model.NewAPICredential(userID, exchange, apiKey, apiSecret, label)

	// Set expiration date if provided
	if expiresIn != nil {
		expiresAt := time.Now().Add(*expiresIn)
		credential.ExpiresAt = &expiresAt

		// Set rotation due date to 80% of the expiration time
		rotationDue := time.Now().Add(time.Duration(float64(*expiresIn) * 0.8))
		credential.RotationDue = &rotationDue
	}

	// Validate the credential
	if err := s.validationService.ValidateCredential(ctx, credential); err != nil {
		s.loggingService.LogCredentialCreate(ctx, credential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, credential.ID, credential.UserID, credential.Exchange)
	}

	// Encrypt the API secret
	encryptedSecret, err := s.encryptionService.Encrypt(credential.APISecret)
	if err != nil {
		s.loggingService.LogCredentialEncrypt(ctx, credential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, credential.ID, credential.UserID, credential.Exchange)
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
		s.loggingService.LogCredentialCreate(ctx, credential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, credential.ID, credential.UserID, credential.Exchange)
	}

	// Log the operation
	s.loggingService.LogCredentialCreate(ctx, credential, time.Since(startTime), nil)

	// Return the credential with masked secret
	credential.APISecret = "********"
	return credential, nil
}

// GetCredential retrieves an API credential by ID
func (s *CredentialLifecycleService) GetCredential(ctx context.Context, id string) (*model.APICredential, error) {
	startTime := time.Now()

	// Get the encrypted credential
	encryptedCredential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Decrypt the API secret
	decryptedSecret, err := s.decryptAPISecret(encryptedCredential.APISecret)
	if err != nil {
		s.loggingService.LogCredentialDecrypt(ctx, encryptedCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, encryptedCredential.ID, encryptedCredential.UserID, encryptedCredential.Exchange)
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

	// Check if the credential is expired
	if decryptedCredential.ExpiresAt != nil && decryptedCredential.ExpiresAt.Before(time.Now()) {
		// Update the status to expired
		oldStatus := decryptedCredential.Status
		decryptedCredential.Status = model.APICredentialStatusExpired
		if err := s.credentialRepo.UpdateStatus(ctx, decryptedCredential.ID, model.APICredentialStatusExpired); err != nil {
			s.logger.Error().Err(err).Str("id", decryptedCredential.ID).Msg("Failed to update credential status to expired")
		} else {
			s.loggingService.LogCredentialStatusChange(ctx, decryptedCredential, oldStatus, time.Since(startTime), nil)
		}
	}

	// Check if the credential needs rotation
	if decryptedCredential.RotationDue != nil && decryptedCredential.RotationDue.Before(time.Now()) {
		// Log a warning
		s.logger.Warn().Str("id", decryptedCredential.ID).Str("userID", decryptedCredential.UserID).Str("exchange", decryptedCredential.Exchange).Msg("API credential rotation is due")
	}

	// Update last used timestamp
	go s.updateLastUsed(context.Background(), decryptedCredential.ID)

	// Log the operation
	s.loggingService.LogCredentialRead(ctx, decryptedCredential, time.Since(startTime), nil)

	return decryptedCredential, nil
}

// GetCredentialByUserAndExchange retrieves an API credential by user ID and exchange
func (s *CredentialLifecycleService) GetCredentialByUserAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	startTime := time.Now()

	// Get the encrypted credential
	encryptedCredential, err := s.credentialRepo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil {
		return nil, s.errorService.HandleError(ctx, err, "", userID, exchange)
	}

	// Decrypt the API secret
	decryptedSecret, err := s.decryptAPISecret(encryptedCredential.APISecret)
	if err != nil {
		s.loggingService.LogCredentialDecrypt(ctx, encryptedCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, encryptedCredential.ID, encryptedCredential.UserID, encryptedCredential.Exchange)
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

	// Check if the credential is expired
	if decryptedCredential.ExpiresAt != nil && decryptedCredential.ExpiresAt.Before(time.Now()) {
		// Update the status to expired
		oldStatus := decryptedCredential.Status
		decryptedCredential.Status = model.APICredentialStatusExpired
		if err := s.credentialRepo.UpdateStatus(ctx, decryptedCredential.ID, model.APICredentialStatusExpired); err != nil {
			s.logger.Error().Err(err).Str("id", decryptedCredential.ID).Msg("Failed to update credential status to expired")
		} else {
			s.loggingService.LogCredentialStatusChange(ctx, decryptedCredential, oldStatus, time.Since(startTime), nil)
		}
	}

	// Check if the credential needs rotation
	if decryptedCredential.RotationDue != nil && decryptedCredential.RotationDue.Before(time.Now()) {
		// Log a warning
		s.logger.Warn().Str("id", decryptedCredential.ID).Str("userID", decryptedCredential.UserID).Str("exchange", decryptedCredential.Exchange).Msg("API credential rotation is due")
	}

	// Update last used timestamp
	go s.updateLastUsed(context.Background(), decryptedCredential.ID)

	// Log the operation
	s.loggingService.LogCredentialRead(ctx, decryptedCredential, time.Since(startTime), nil)

	return decryptedCredential, nil
}

// UpdateCredential updates an API credential
func (s *CredentialLifecycleService) UpdateCredential(ctx context.Context, id, apiKey, apiSecret, label string) (*model.APICredential, error) {
	startTime := time.Now()

	// Get the existing credential
	existingCredential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Create a copy of the existing credential
	updatedCredential := &model.APICredential{
		ID:           existingCredential.ID,
		UserID:       existingCredential.UserID,
		Exchange:     existingCredential.Exchange,
		APIKey:       existingCredential.APIKey,
		APISecret:    existingCredential.APISecret,
		Label:        existingCredential.Label,
		Status:       existingCredential.Status,
		LastUsed:     existingCredential.LastUsed,
		LastVerified: existingCredential.LastVerified,
		ExpiresAt:    existingCredential.ExpiresAt,
		RotationDue:  existingCredential.RotationDue,
		FailureCount: existingCredential.FailureCount,
		Metadata:     existingCredential.Metadata,
		CreatedAt:    existingCredential.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	// Update the fields if provided
	if apiKey != "" {
		updatedCredential.APIKey = apiKey
	}

	if apiSecret != "" {
		// Encrypt the new API secret
		encryptedSecret, err := s.encryptionService.Encrypt(apiSecret)
		if err != nil {
			s.loggingService.LogCredentialEncrypt(ctx, updatedCredential, time.Since(startTime), err)
			return nil, s.errorService.HandleError(ctx, err, updatedCredential.ID, updatedCredential.UserID, updatedCredential.Exchange)
		}
		updatedCredential.APISecret = string(encryptedSecret)

		// Reset the last verified timestamp and failure count
		now := time.Now()
		updatedCredential.LastVerified = &now
		updatedCredential.FailureCount = 0
	}

	if label != "" {
		updatedCredential.Label = label
	}

	// Validate the updated credential
	if err := s.validationService.ValidateCredential(ctx, updatedCredential); err != nil {
		s.loggingService.LogCredentialUpdate(ctx, updatedCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, updatedCredential.ID, updatedCredential.UserID, updatedCredential.Exchange)
	}

	// Save the updated credential
	if err := s.credentialRepo.Save(ctx, updatedCredential); err != nil {
		s.loggingService.LogCredentialUpdate(ctx, updatedCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, updatedCredential.ID, updatedCredential.UserID, updatedCredential.Exchange)
	}

	// Log the operation
	s.loggingService.LogCredentialUpdate(ctx, updatedCredential, time.Since(startTime), nil)

	// Return the updated credential with masked secret
	updatedCredential.APISecret = "********"
	return updatedCredential, nil
}

// DeleteCredential deletes an API credential
func (s *CredentialLifecycleService) DeleteCredential(ctx context.Context, id string) error {
	startTime := time.Now()

	// Get the credential to log the operation
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Delete the credential
	if err := s.credentialRepo.DeleteByID(ctx, id); err != nil {
		s.loggingService.LogCredentialDelete(ctx, id, credential.UserID, credential.Exchange, time.Since(startTime), err)
		return s.errorService.HandleError(ctx, err, id, credential.UserID, credential.Exchange)
	}

	// Log the operation
	s.loggingService.LogCredentialDelete(ctx, id, credential.UserID, credential.Exchange, time.Since(startTime), nil)

	return nil
}

// VerifyCredential verifies an API credential
func (s *CredentialLifecycleService) VerifyCredential(ctx context.Context, id string) error {
	startTime := time.Now()

	// Get the credential
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Try to decrypt the API secret
	_, err = s.decryptAPISecret(credential.APISecret)
	if err != nil {
		s.loggingService.LogCredentialVerify(ctx, credential, time.Since(startTime), err)

		// Increment failure count
		if err := s.credentialRepo.IncrementFailureCount(ctx, id); err != nil {
			s.logger.Error().Err(err).Str("id", id).Msg("Failed to increment failure count")
		}

		// Update status to failed if failure count exceeds threshold
		if credential.FailureCount >= 5 {
			oldStatus := credential.Status
			credential.Status = model.APICredentialStatusFailed
			if err := s.credentialRepo.UpdateStatus(ctx, id, model.APICredentialStatusFailed); err != nil {
				s.logger.Error().Err(err).Str("id", id).Msg("Failed to update credential status to failed")
			} else {
				s.loggingService.LogCredentialStatusChange(ctx, credential, oldStatus, time.Since(startTime), nil)
			}
		}

		return s.errorService.HandleError(ctx, err, id, credential.UserID, credential.Exchange)
	}

	// Update last verified timestamp and reset failure count
	now := time.Now()
	if err := s.credentialRepo.UpdateLastVerified(ctx, id, now); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to update last verified timestamp")
	}

	if err := s.credentialRepo.ResetFailureCount(ctx, id); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to reset failure count")
	}

	// Log the operation
	s.loggingService.LogCredentialVerify(ctx, credential, time.Since(startTime), nil)

	return nil
}

// RotateCredential rotates an API credential
func (s *CredentialLifecycleService) RotateCredential(ctx context.Context, id, newAPIKey, newAPISecret string) (*model.APICredential, error) {
	startTime := time.Now()

	// Get the existing credential
	existingCredential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Create a new credential with the same user ID, exchange, and label
	newCredential := model.NewAPICredential(
		existingCredential.UserID,
		existingCredential.Exchange,
		newAPIKey,
		newAPISecret,
		existingCredential.Label+" (rotated)",
	)

	// Copy expiration and rotation due dates
	newCredential.ExpiresAt = existingCredential.ExpiresAt
	newCredential.RotationDue = existingCredential.RotationDue

	// Validate the new credential
	if err := s.validationService.ValidateCredential(ctx, newCredential); err != nil {
		s.loggingService.LogCredentialCreate(ctx, newCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, newCredential.ID, newCredential.UserID, newCredential.Exchange)
	}

	// Encrypt the API secret
	encryptedSecret, err := s.encryptionService.Encrypt(newCredential.APISecret)
	if err != nil {
		s.loggingService.LogCredentialEncrypt(ctx, newCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, newCredential.ID, newCredential.UserID, newCredential.Exchange)
	}

	// Create a copy of the credential with the encrypted secret
	encryptedCredential := &model.APICredential{
		ID:           newCredential.ID,
		UserID:       newCredential.UserID,
		Exchange:     newCredential.Exchange,
		APIKey:       newCredential.APIKey,
		APISecret:    string(encryptedSecret), // Store the encrypted secret
		Label:        newCredential.Label,
		Status:       newCredential.Status,
		LastUsed:     newCredential.LastUsed,
		LastVerified: newCredential.LastVerified,
		ExpiresAt:    newCredential.ExpiresAt,
		RotationDue:  newCredential.RotationDue,
		FailureCount: newCredential.FailureCount,
		Metadata:     newCredential.Metadata,
		CreatedAt:    newCredential.CreatedAt,
		UpdatedAt:    newCredential.UpdatedAt,
	}

	// Save the new credential
	if err := s.credentialRepo.Save(ctx, encryptedCredential); err != nil {
		s.loggingService.LogCredentialCreate(ctx, newCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, newCredential.ID, newCredential.UserID, newCredential.Exchange)
	}

	// Update the status of the old credential to revoked
	oldStatus := existingCredential.Status
	existingCredential.Status = model.APICredentialStatusRevoked
	if err := s.credentialRepo.UpdateStatus(ctx, existingCredential.ID, model.APICredentialStatusRevoked); err != nil {
		s.logger.Error().Err(err).Str("id", existingCredential.ID).Msg("Failed to update old credential status to revoked")
	} else {
		s.loggingService.LogCredentialStatusChange(ctx, existingCredential, oldStatus, time.Since(startTime), nil)
	}

	// Log the operation
	s.loggingService.LogCredentialCreate(ctx, newCredential, time.Since(startTime), nil)

	// Return the new credential with masked secret
	newCredential.APISecret = "********"
	return newCredential, nil
}

// ExtendCredential extends the expiration date of an API credential
func (s *CredentialLifecycleService) ExtendCredential(ctx context.Context, id string, extension time.Duration) (*model.APICredential, error) {
	startTime := time.Now()

	// Get the existing credential
	existingCredential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Create a copy of the existing credential
	updatedCredential := &model.APICredential{
		ID:           existingCredential.ID,
		UserID:       existingCredential.UserID,
		Exchange:     existingCredential.Exchange,
		APIKey:       existingCredential.APIKey,
		APISecret:    existingCredential.APISecret,
		Label:        existingCredential.Label,
		Status:       existingCredential.Status,
		LastUsed:     existingCredential.LastUsed,
		LastVerified: existingCredential.LastVerified,
		FailureCount: existingCredential.FailureCount,
		Metadata:     existingCredential.Metadata,
		CreatedAt:    existingCredential.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	// Calculate new expiration date
	var newExpiresAt time.Time
	if existingCredential.ExpiresAt != nil {
		newExpiresAt = existingCredential.ExpiresAt.Add(extension)
	} else {
		newExpiresAt = time.Now().Add(extension)
	}
	updatedCredential.ExpiresAt = &newExpiresAt

	// Calculate new rotation due date (80% of the time to expiration)
	rotationDue := time.Now().Add(time.Duration(float64(extension) * 0.8))
	updatedCredential.RotationDue = &rotationDue

	// If the credential was expired, reactivate it
	if existingCredential.Status == model.APICredentialStatusExpired {
		oldStatus := existingCredential.Status
		updatedCredential.Status = model.APICredentialStatusActive
		s.loggingService.LogCredentialStatusChange(ctx, updatedCredential, oldStatus, time.Since(startTime), nil)
	}

	// Save the updated credential
	if err := s.credentialRepo.Save(ctx, updatedCredential); err != nil {
		s.loggingService.LogCredentialUpdate(ctx, updatedCredential, time.Since(startTime), err)
		return nil, s.errorService.HandleError(ctx, err, updatedCredential.ID, updatedCredential.UserID, updatedCredential.Exchange)
	}

	// Log the operation
	s.loggingService.LogCredentialUpdate(ctx, updatedCredential, time.Since(startTime), nil)

	// Return the updated credential with masked secret
	updatedCredential.APISecret = "********"
	return updatedCredential, nil
}

// ActivateCredential activates an API credential
func (s *CredentialLifecycleService) ActivateCredential(ctx context.Context, id string) error {
	startTime := time.Now()

	// Get the credential
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Check if the credential is already active
	if credential.Status == model.APICredentialStatusActive {
		return nil
	}

	// Check if the credential is expired
	if credential.ExpiresAt != nil && credential.ExpiresAt.Before(time.Now()) {
		return s.errorService.HandleError(ctx, fmt.Errorf("cannot activate expired credential"), id, credential.UserID, credential.Exchange)
	}

	// Update the status to active
	oldStatus := credential.Status
	credential.Status = model.APICredentialStatusActive
	if err := s.credentialRepo.UpdateStatus(ctx, id, model.APICredentialStatusActive); err != nil {
		s.loggingService.LogCredentialStatusChange(ctx, credential, oldStatus, time.Since(startTime), err)
		return s.errorService.HandleError(ctx, err, id, credential.UserID, credential.Exchange)
	}

	// Log the operation
	s.loggingService.LogCredentialStatusChange(ctx, credential, oldStatus, time.Since(startTime), nil)

	return nil
}

// DeactivateCredential deactivates an API credential
func (s *CredentialLifecycleService) DeactivateCredential(ctx context.Context, id string) error {
	startTime := time.Now()

	// Get the credential
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Check if the credential is already inactive
	if credential.Status == model.APICredentialStatusInactive {
		return nil
	}

	// Update the status to inactive
	oldStatus := credential.Status
	credential.Status = model.APICredentialStatusInactive
	if err := s.credentialRepo.UpdateStatus(ctx, id, model.APICredentialStatusInactive); err != nil {
		s.loggingService.LogCredentialStatusChange(ctx, credential, oldStatus, time.Since(startTime), err)
		return s.errorService.HandleError(ctx, err, id, credential.UserID, credential.Exchange)
	}

	// Log the operation
	s.loggingService.LogCredentialStatusChange(ctx, credential, oldStatus, time.Since(startTime), nil)

	return nil
}

// RevokeCredential revokes an API credential
func (s *CredentialLifecycleService) RevokeCredential(ctx context.Context, id string) error {
	startTime := time.Now()

	// Get the credential
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return s.errorService.HandleError(ctx, err, id, "", "")
	}

	// Check if the credential is already revoked
	if credential.Status == model.APICredentialStatusRevoked {
		return nil
	}

	// Update the status to revoked
	oldStatus := credential.Status
	credential.Status = model.APICredentialStatusRevoked
	if err := s.credentialRepo.UpdateStatus(ctx, id, model.APICredentialStatusRevoked); err != nil {
		s.loggingService.LogCredentialStatusChange(ctx, credential, oldStatus, time.Since(startTime), err)
		return s.errorService.HandleError(ctx, err, id, credential.UserID, credential.Exchange)
	}

	// Log the operation
	s.loggingService.LogCredentialStatusChange(ctx, credential, oldStatus, time.Since(startTime), nil)

	return nil
}

// ListCredentials lists all API credentials for a user
func (s *CredentialLifecycleService) ListCredentials(ctx context.Context, userID string) ([]*model.APICredential, error) {
	startTime := time.Now()

	// Get all credentials for the user
	credentials, err := s.credentialRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, s.errorService.HandleError(ctx, err, "", userID, "")
	}

	// Mask the API secrets
	for _, credential := range credentials {
		credential.APISecret = "********"
	}

	// Log the operation
	s.logger.Debug().Str("userID", userID).Int("count", len(credentials)).Dur("duration", time.Since(startTime)).Msg("Listed API credentials")

	return credentials, nil
}

// decryptAPISecret decrypts an API secret
func (s *CredentialLifecycleService) decryptAPISecret(encryptedSecret string) (string, error) {
	// Try to decrypt as bytes first
	decryptedSecret, err := s.encryptionService.Decrypt([]byte(encryptedSecret))
	if err != nil {
		// If that fails, try to decrypt as a string (for backward compatibility)
		return "", fmt.Errorf("failed to decrypt API secret: %w", err)
	}

	return decryptedSecret, nil
}

// updateLastUsed updates the last used timestamp of an API credential
func (s *CredentialLifecycleService) updateLastUsed(ctx context.Context, id string) {
	now := time.Now()
	if err := s.credentialRepo.UpdateLastUsed(ctx, id, now); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to update last used timestamp")
	}
}
