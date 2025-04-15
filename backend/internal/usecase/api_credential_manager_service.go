package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
	"github.com/rs/zerolog"
)

// APICredentialManagerService defines the interface for managing API credentials
type APICredentialManagerService interface {
	// CreateCredential creates a new API credential
	CreateCredential(ctx context.Context, userID, exchange, apiKey, apiSecret, label string) (*model.APICredential, error)

	// GetCredential gets an API credential by ID
	GetCredential(ctx context.Context, id string) (*model.APICredential, error)

	// GetCredentialByUserIDAndExchange gets an API credential by user ID and exchange
	GetCredentialByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error)

	// GetCredentialByUserIDAndLabel gets an API credential by user ID, exchange, and label
	GetCredentialByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error)

	// UpdateCredential updates an API credential
	UpdateCredential(ctx context.Context, id, apiKey, apiSecret, label string) (*model.APICredential, error)

	// DeleteCredential deletes an API credential
	DeleteCredential(ctx context.Context, id string) error

	// ListCredentialsByUserID lists API credentials by user ID
	ListCredentialsByUserID(ctx context.Context, userID string) ([]*model.APICredential, error)

	// VerifyCredential verifies an API credential with the exchange
	VerifyCredential(ctx context.Context, id string) (bool, error)

	// RotateCredential rotates an API credential
	RotateCredential(ctx context.Context, id string, newAPIKey, newAPISecret string) (*model.APICredential, error)

	// MarkCredentialAsUsed marks an API credential as used
	MarkCredentialAsUsed(ctx context.Context, id string) error

	// GetCredentialForExchange gets a valid API credential for an exchange
	GetCredentialForExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error)
}

// apiCredentialManagerService implements APICredentialManagerService
type apiCredentialManagerService struct {
	repo            port.APICredentialRepository
	encryptionSvc   crypto.EncryptionService
	exchangeClients map[string]port.ExchangeWalletProvider
	logger          *zerolog.Logger
}

// NewAPICredentialManagerService creates a new API credential manager service
func NewAPICredentialManagerService(
	repo port.APICredentialRepository,
	encryptionSvc crypto.EncryptionService,
	providerRegistry *wallet.ProviderRegistry,
	logger *zerolog.Logger,
) APICredentialManagerService {
	// Get all exchange providers from the registry
	exchangeProviders := providerRegistry.GetAllExchangeProviders()
	exchangeClients := make(map[string]port.ExchangeWalletProvider, len(exchangeProviders))

	// Map exchange providers by name
	for _, provider := range exchangeProviders {
		exchangeClients[provider.GetName()] = provider
	}

	return &apiCredentialManagerService{
		repo:            repo,
		encryptionSvc:   encryptionSvc,
		exchangeClients: exchangeClients,
		logger:          logger,
	}
}

// CreateCredential creates a new API credential
func (s *apiCredentialManagerService) CreateCredential(
	ctx context.Context,
	userID, exchange, apiKey, apiSecret, label string,
) (*model.APICredential, error) {
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

	// Check if exchange is supported
	if _, ok := s.exchangeClients[exchange]; !ok {
		return nil, fmt.Errorf("unsupported exchange: %s", exchange)
	}

	// Check if credential already exists for this user and exchange
	existingCred, err := s.repo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil && !errors.Is(err, model.ErrCredentialNotFound) {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to check for existing credential")
		return nil, err
	}

	if existingCred != nil && label == "" {
		return nil, errors.New("credential already exists for this user and exchange, please provide a label")
	}

	// Check if credential with the same label already exists
	if label != "" {
		existingLabelCred, err := s.repo.GetByUserIDAndLabel(ctx, userID, exchange, label)
		if err != nil && !errors.Is(err, model.ErrCredentialNotFound) {
			s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Str("label", label).Msg("Failed to check for existing credential with label")
			return nil, err
		}

		if existingLabelCred != nil {
			return nil, errors.New("credential with this label already exists")
		}
	}

	// Create new credential
	credential := model.NewAPICredential(userID, exchange, apiKey, apiSecret, label)

	// Verify credential with exchange
	if err := s.verifyWithExchange(ctx, credential); err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to verify credential with exchange")
		credential.Status = model.APICredentialStatusFailed
		credential.FailureCount = 1
	} else {
		credential.Status = model.APICredentialStatusActive
		now := time.Now()
		credential.LastVerified = &now
	}

	// Save credential
	if err := s.repo.Save(ctx, credential); err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to save credential")
		return nil, err
	}

	s.logger.Info().Str("userID", userID).Str("exchange", exchange).Str("id", credential.ID).Msg("Created API credential")
	return credential, nil
}

// GetCredential gets an API credential by ID
func (s *apiCredentialManagerService) GetCredential(ctx context.Context, id string) (*model.APICredential, error) {
	if id == "" {
		return nil, errors.New("credential ID is required")
	}

	credential, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get credential")
		return nil, err
	}

	return credential, nil
}

// GetCredentialByUserIDAndExchange gets an API credential by user ID and exchange
func (s *apiCredentialManagerService) GetCredentialByUserIDAndExchange(
	ctx context.Context,
	userID, exchange string,
) (*model.APICredential, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if exchange == "" {
		return nil, errors.New("exchange is required")
	}

	credential, err := s.repo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to get credential")
		return nil, err
	}

	return credential, nil
}

// GetCredentialByUserIDAndLabel gets an API credential by user ID, exchange, and label
func (s *apiCredentialManagerService) GetCredentialByUserIDAndLabel(
	ctx context.Context,
	userID, exchange, label string,
) (*model.APICredential, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if exchange == "" {
		return nil, errors.New("exchange is required")
	}
	if label == "" {
		return nil, errors.New("label is required")
	}

	credential, err := s.repo.GetByUserIDAndLabel(ctx, userID, exchange, label)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Str("label", label).Msg("Failed to get credential")
		return nil, err
	}

	return credential, nil
}

// UpdateCredential updates an API credential
func (s *apiCredentialManagerService) UpdateCredential(
	ctx context.Context,
	id, apiKey, apiSecret, label string,
) (*model.APICredential, error) {
	if id == "" {
		return nil, errors.New("credential ID is required")
	}

	// Get existing credential
	credential, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get credential")
		return nil, err
	}

	// Update fields
	updated := false
	if apiKey != "" && apiKey != credential.APIKey {
		credential.APIKey = apiKey
		updated = true
	}
	if apiSecret != "" && apiSecret != credential.APISecret {
		credential.APISecret = apiSecret
		updated = true
	}
	if label != "" && label != credential.Label {
		// Check if credential with the same label already exists
		existingLabelCred, err := s.repo.GetByUserIDAndLabel(ctx, credential.UserID, credential.Exchange, label)
		if err != nil && !errors.Is(err, model.ErrCredentialNotFound) {
			s.logger.Error().Err(err).Str("userID", credential.UserID).Str("exchange", credential.Exchange).Str("label", label).Msg("Failed to check for existing credential with label")
			return nil, err
		}

		if existingLabelCred != nil && existingLabelCred.ID != id {
			return nil, errors.New("credential with this label already exists")
		}

		credential.Label = label
		updated = true
	}

	if !updated {
		return credential, nil
	}

	// If API key or secret was updated, verify with exchange
	if apiKey != "" || apiSecret != "" {
		if err := s.verifyWithExchange(ctx, credential); err != nil {
			s.logger.Error().Err(err).Str("id", id).Msg("Failed to verify updated credential with exchange")
			credential.Status = model.APICredentialStatusFailed
			credential.FailureCount++
		} else {
			credential.Status = model.APICredentialStatusActive
			now := time.Now()
			credential.LastVerified = &now
			credential.FailureCount = 0
		}
	}

	// Update timestamp
	credential.UpdatedAt = time.Now()

	// Save updated credential
	if err := s.repo.Save(ctx, credential); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to save updated credential")
		return nil, err
	}

	s.logger.Info().Str("id", id).Msg("Updated API credential")
	return credential, nil
}

// DeleteCredential deletes an API credential
func (s *apiCredentialManagerService) DeleteCredential(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("credential ID is required")
	}

	// Get credential to check if it exists
	credential, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get credential")
		return err
	}

	// Delete credential
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to delete credential")
		return err
	}

	s.logger.Info().Str("id", id).Str("userID", credential.UserID).Str("exchange", credential.Exchange).Msg("Deleted API credential")
	return nil
}

// ListCredentialsByUserID lists API credentials by user ID
func (s *apiCredentialManagerService) ListCredentialsByUserID(ctx context.Context, userID string) ([]*model.APICredential, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	credentials, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to list credentials")
		return nil, err
	}

	return credentials, nil
}

// VerifyCredential verifies an API credential with the exchange
func (s *apiCredentialManagerService) VerifyCredential(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("credential ID is required")
	}

	// Get credential
	credential, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get credential")
		return false, err
	}

	// Verify with exchange
	err = s.verifyWithExchange(ctx, credential)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to verify credential with exchange")

		// Update status and failure count
		credential.Status = model.APICredentialStatusFailed
		credential.FailureCount++

		if err := s.repo.Save(ctx, credential); err != nil {
			s.logger.Error().Err(err).Str("id", id).Msg("Failed to update credential status after verification failure")
		}

		return false, err
	}

	// Update verification timestamp and status
	now := time.Now()
	credential.LastVerified = &now
	credential.Status = model.APICredentialStatusActive
	credential.FailureCount = 0

	if err := s.repo.Save(ctx, credential); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to update credential after successful verification")
	}

	return true, nil
}

// RotateCredential rotates an API credential
func (s *apiCredentialManagerService) RotateCredential(
	ctx context.Context,
	id string,
	newAPIKey, newAPISecret string,
) (*model.APICredential, error) {
	if id == "" {
		return nil, errors.New("credential ID is required")
	}
	if newAPIKey == "" {
		return nil, errors.New("new API key is required")
	}
	if newAPISecret == "" {
		return nil, errors.New("new API secret is required")
	}

	// Get existing credential
	credential, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get credential")
		return nil, err
	}

	// Store old values
	oldAPIKey := credential.APIKey
	oldAPISecret := credential.APISecret

	// Update with new values
	credential.APIKey = newAPIKey
	credential.APISecret = newAPISecret
	credential.UpdatedAt = time.Now()
	now := time.Now()
	credential.RotationDue = &now

	// Verify new credentials with exchange
	err = s.verifyWithExchange(ctx, credential)
	if err != nil {
		// Revert to old values
		credential.APIKey = oldAPIKey
		credential.APISecret = oldAPISecret

		s.logger.Error().Err(err).Str("id", id).Msg("Failed to verify new credentials with exchange")
		return nil, fmt.Errorf("failed to verify new credentials: %w", err)
	}

	// Update verification timestamp and status
	credential.LastVerified = &now
	credential.Status = model.APICredentialStatusActive
	credential.FailureCount = 0

	// Save updated credential
	if err := s.repo.Save(ctx, credential); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to save rotated credential")
		return nil, err
	}

	s.logger.Info().Str("id", id).Msg("Rotated API credential")
	return credential, nil
}

// MarkCredentialAsUsed marks an API credential as used
func (s *apiCredentialManagerService) MarkCredentialAsUsed(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("credential ID is required")
	}

	// Update last used timestamp
	now := time.Now()
	if err := s.repo.UpdateLastUsed(ctx, id, now); err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to update last used timestamp")
		return err
	}

	return nil
}

// GetCredentialForExchange gets a valid API credential for an exchange
func (s *apiCredentialManagerService) GetCredentialForExchange(
	ctx context.Context,
	userID, exchange string,
) (*model.APICredential, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if exchange == "" {
		return nil, errors.New("exchange is required")
	}

	// Get credential
	credential, err := s.repo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to get credential")
		return nil, err
	}

	// Check if credential is active
	if credential.Status != model.APICredentialStatusActive {
		return nil, fmt.Errorf("credential is not active (status: %s)", credential.Status)
	}

	// Mark as used
	now := time.Now()
	if err := s.repo.UpdateLastUsed(ctx, credential.ID, now); err != nil {
		s.logger.Error().Err(err).Str("id", credential.ID).Msg("Failed to update last used timestamp")
		// Continue anyway, this is not critical
	}

	return credential, nil
}

// verifyWithExchange verifies an API credential with the exchange
func (s *apiCredentialManagerService) verifyWithExchange(ctx context.Context, credential *model.APICredential) error {
	// Get exchange client
	exchangeClient, ok := s.exchangeClients[credential.Exchange]
	if !ok {
		return fmt.Errorf("unsupported exchange: %s", credential.Exchange)
	}

	// Set API credentials
	if err := exchangeClient.SetAPICredentials(ctx, credential.APIKey, credential.APISecret); err != nil {
		return fmt.Errorf("failed to set API credentials: %w", err)
	}

	// Verify by attempting to get account information
	_, err := exchangeClient.GetBalance(ctx, &model.Wallet{
		UserID:   credential.UserID,
		Exchange: credential.Exchange,
	})

	return err
}
