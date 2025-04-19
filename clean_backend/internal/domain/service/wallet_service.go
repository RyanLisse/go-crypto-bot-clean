package service

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var (
	// ErrWalletNotFound is returned when a wallet is not found for the given criteria.
	ErrWalletNotFound = errors.New("wallet not found")
	// ErrUnauthorizedAccess is returned when a user tries to access a wallet they don't own.
	ErrUnauthorizedAccess = errors.New("unauthorized access to wallet")
	// ErrInvalidWalletType is returned when an operation is attempted on an unsuitable wallet type.
	ErrInvalidWalletType = errors.New("invalid wallet type for this operation")
	// ErrVerificationFailed might be used by underlying verification mechanisms.
	ErrVerificationFailed = errors.New("wallet verification failed")
	// ErrPrimaryWalletExists is returned when trying to set a primary wallet and one already exists (if applicable).
	ErrPrimaryWalletExists = errors.New("primary wallet already exists for this user")
)

// WalletService defines the interface for wallet business operations
type WalletService interface {
	// Web3 Wallet Operations
	CreateWeb3Wallet(ctx context.Context, userID uuid.UUID, address, network string, isPrimary bool) (*model.Wallet, error)

	// Exchange Wallet Operations
	CreateExchangeWallet(ctx context.Context, userID uuid.UUID, exchangeID string, exchangeName string, isPrimary bool) (*model.Wallet, error)

	// General Wallet Operations
	GetWallet(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*model.Wallet, error)
	GetPrimaryWallet(ctx context.Context, userID uuid.UUID) (*model.Wallet, error)
	SetPrimaryWallet(ctx context.Context, walletID, userID uuid.UUID) error
	UpdateBalances(ctx context.Context, walletID uuid.UUID, balances []*model.Balance) error
	VerifyWallet(ctx context.Context, walletID uuid.UUID) error
	CheckVerificationStatus(ctx context.Context, walletID uuid.UUID) (model.VerificationStatus, error)
	DeleteWallet(ctx context.Context, walletID, userID uuid.UUID) error
}

// walletService implements the WalletService interface
type walletService struct {
	repo   port.WalletRepository
	logger *zerolog.Logger
}

// NewWalletService creates a new wallet service
func NewWalletService(repo port.WalletRepository, logger *zerolog.Logger) WalletService {
	return &walletService{
		repo:   repo,
		logger: logger,
	}
}

// CreateWeb3Wallet creates a new Web3 wallet
func (s *walletService) CreateWeb3Wallet(ctx context.Context, userID uuid.UUID, address, network string, isPrimary bool) (*model.Wallet, error) {
	s.logger.Info().Str("userID", userID.String()).Str("address", address).Str("network", network).Bool("isPrimary", isPrimary).Msg("Creating Web3 wallet")

	// Create a new wallet
	wallet := &model.Wallet{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      model.WalletTypeWeb3,
		Status:    model.WalletStatusPending,
		IsPrimary: isPrimary,
		Address:   address,
		Network:   network,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the wallet
	err := s.repo.Create(ctx, wallet)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create Web3 wallet")
		return nil, err
	}

	return wallet, nil
}

// CreateExchangeWallet creates a new exchange wallet
func (s *walletService) CreateExchangeWallet(ctx context.Context, userID uuid.UUID, exchangeID, exchangeName string, isPrimary bool) (*model.Wallet, error) {
	s.logger.Info().Str("userID", userID.String()).Str("exchangeID", exchangeID).Str("exchangeName", exchangeName).Bool("isPrimary", isPrimary).Msg("Creating exchange wallet")

	// Create a new wallet
	wallet := &model.Wallet{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         model.WalletTypeExchange,
		Status:       model.WalletStatusPending,
		IsPrimary:    isPrimary,
		ExchangeID:   exchangeID,
		ExchangeName: exchangeName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save the wallet
	err := s.repo.Create(ctx, wallet)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create exchange wallet")
		return nil, err
	}

	return wallet, nil
}

// GetWallet retrieves a wallet by ID
func (s *walletService) GetWallet(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	s.logger.Info().Str("walletID", id.String()).Msg("Getting wallet")
	return s.repo.GetByID(ctx, id)
}

// GetUserWallets retrieves all wallets for a user
func (s *walletService) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*model.Wallet, error) {
	s.logger.Info().Str("userID", userID.String()).Msg("Getting user wallets")
	return s.repo.GetByUserID(ctx, userID)
}

// GetPrimaryWallet retrieves the primary wallet for a user
func (s *walletService) GetPrimaryWallet(ctx context.Context, userID uuid.UUID) (*model.Wallet, error) {
	s.logger.Info().Str("userID", userID.String()).Msg("Getting primary wallet")
	// Since we don't have a specific wallet type, we'll just use Web3 as a placeholder
	return s.repo.FindPrimary(ctx, userID, model.WalletTypeWeb3)
}

// SetPrimaryWallet sets a wallet as primary
func (s *walletService) SetPrimaryWallet(ctx context.Context, walletID, userID uuid.UUID) error {
	s.logger.Info().Str("walletID", walletID.String()).Str("userID", userID.String()).Msg("Setting primary wallet")

	// Get the wallet to determine its type
	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get wallet")
		return err
	}

	// Check ownership
	if wallet.UserID != userID {
		s.logger.Error().Msg("Unauthorized access to wallet")
		return ErrUnauthorizedAccess
	}

	// Set as primary
	return s.repo.SetPrimary(ctx, userID, walletID, wallet.Type)
}

// UpdateBalances updates wallet balances
func (s *walletService) UpdateBalances(ctx context.Context, walletID uuid.UUID, balances []*model.Balance) error {
	s.logger.Info().Str("walletID", walletID.String()).Msg("Updating wallet balances")

	// Convert []*model.Balance to []model.Balance
	balanceSlice := make([]model.Balance, len(balances))
	for i, b := range balances {
		if b != nil {
			balanceSlice[i] = *b
		}
	}

	return s.repo.UpdateBalance(ctx, walletID, balanceSlice)
}

// VerifyWallet verifies a wallet
func (s *walletService) VerifyWallet(ctx context.Context, walletID uuid.UUID) error {
	s.logger.Info().Str("walletID", walletID.String()).Msg("Verifying wallet")
	return s.repo.UpdateStatus(ctx, walletID, model.WalletStatusVerified)
}

// CheckVerificationStatus checks the verification status of a wallet
func (s *walletService) CheckVerificationStatus(ctx context.Context, walletID uuid.UUID) (model.VerificationStatus, error) {
	s.logger.Info().Str("walletID", walletID.String()).Msg("Checking verification status")

	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get wallet")
		return "", err
	}

	// Convert WalletStatus to VerificationStatus
	switch wallet.Status {
	case model.WalletStatusVerified:
		return model.VerificationStatusVerified, nil
	case model.WalletStatusPending:
		return model.VerificationStatusPending, nil
	default:
		return model.VerificationStatusNone, nil
	}
}

// DeleteWallet deletes a wallet
func (s *walletService) DeleteWallet(ctx context.Context, walletID, userID uuid.UUID) error {
	s.logger.Info().Str("walletID", walletID.String()).Str("userID", userID.String()).Msg("Deleting wallet")

	// Get the wallet to check ownership
	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get wallet")
		return err
	}

	// Check ownership
	if wallet.UserID != userID {
		s.logger.Error().Msg("Unauthorized access to wallet")
		return ErrUnauthorizedAccess
	}

	return s.repo.Delete(ctx, walletID)
}
