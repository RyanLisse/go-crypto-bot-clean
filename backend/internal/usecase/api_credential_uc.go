package usecase

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// APICredentialUseCase defines the API credential use case interface
type APICredentialUseCase interface {
	CreateCredential(ctx context.Context, credential *model.APICredential) error
	GetCredential(ctx context.Context, id string) (*model.APICredential, error)
	UpdateCredential(ctx context.Context, credential *model.APICredential) error
	DeleteCredential(ctx context.Context, id string) error
	ListCredentials(ctx context.Context, userID string) ([]*model.APICredential, error)
	GetCredentialByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error)
}

// apiCredentialUsecase implements APICredentialUseCase
type apiCredentialUsecase struct {
	repo   port.APICredentialRepository
	logger *zerolog.Logger
}

// NewAPICredentialUseCase creates a new APICredentialUseCase
func NewAPICredentialUseCase(repo port.APICredentialRepository, logger *zerolog.Logger) APICredentialUseCase {
	return &apiCredentialUsecase{
		repo:   repo,
		logger: logger,
	}
}

// CreateCredential creates a new API credential
func (uc *apiCredentialUsecase) CreateCredential(ctx context.Context, credential *model.APICredential) error {
	// Validate credential
	if err := credential.Validate(); err != nil {
		uc.logger.Error().Err(err).Str("userID", credential.UserID).Msg("Invalid API credential")
		return err
	}

	// Save credential
	if err := uc.repo.Save(ctx, credential); err != nil {
		uc.logger.Error().Err(err).Str("userID", credential.UserID).Msg("Failed to save API credential")
		return err
	}

	return nil
}

// GetCredential gets an API credential by ID
func (uc *apiCredentialUsecase) GetCredential(ctx context.Context, id string) (*model.APICredential, error) {
	// Get credential
	credential, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get API credential")
		return nil, err
	}

	return credential, nil
}

// UpdateCredential updates an API credential
func (uc *apiCredentialUsecase) UpdateCredential(ctx context.Context, credential *model.APICredential) error {
	// Validate credential
	if err := credential.Validate(); err != nil {
		uc.logger.Error().Err(err).Str("id", credential.ID).Msg("Invalid API credential")
		return err
	}

	// Save credential
	if err := uc.repo.Save(ctx, credential); err != nil {
		uc.logger.Error().Err(err).Str("id", credential.ID).Msg("Failed to update API credential")
		return err
	}

	return nil
}

// DeleteCredential deletes an API credential
func (uc *apiCredentialUsecase) DeleteCredential(ctx context.Context, id string) error {
	// Delete credential
	if err := uc.repo.DeleteByID(ctx, id); err != nil {
		uc.logger.Error().Err(err).Str("id", id).Msg("Failed to delete API credential")
		return err
	}

	return nil
}

// ListCredentials lists API credentials for a user
func (uc *apiCredentialUsecase) ListCredentials(ctx context.Context, userID string) ([]*model.APICredential, error) {
	// List credentials
	credentials, err := uc.repo.ListByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to list API credentials")
		return nil, err
	}

	return credentials, nil
}

// GetCredentialByUserIDAndExchange gets an API credential by user ID and exchange
func (uc *apiCredentialUsecase) GetCredentialByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	// Get credential
	credential, err := uc.repo.GetByUserIDAndExchange(ctx, userID, exchange)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to get API credential")
		return nil, err
	}

	return credential, nil
}
