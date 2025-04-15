package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// SignatureVerificationService defines the interface for signature verification operations
type SignatureVerificationService interface {
	// GenerateChallenge generates a challenge message for a wallet to sign
	GenerateChallenge(ctx context.Context, walletID string) (string, error)

	// VerifySignature verifies a signature against a challenge
	VerifySignature(ctx context.Context, walletID, challenge, signature string) (bool, error)

	// GetWalletStatus gets the verification status of a wallet
	GetWalletStatus(ctx context.Context, walletID string) (model.WalletStatus, error)

	// SetWalletStatus sets the verification status of a wallet
	SetWalletStatus(ctx context.Context, walletID string, status model.WalletStatus) error
}

// Challenge represents a signature challenge
type Challenge struct {
	WalletID  string    `json:"wallet_id"`
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

// signatureVerificationService implements the SignatureVerificationService interface
type signatureVerificationService struct {
	providerRegistry *wallet.ProviderRegistry
	walletRepo       port.WalletRepository
	challenges       map[string]*Challenge // In-memory store of challenges (should be replaced with a persistent store in production)
	logger           *zerolog.Logger
}

// NewSignatureVerificationService creates a new signature verification service
func NewSignatureVerificationService(
	providerRegistry *wallet.ProviderRegistry,
	walletRepo port.WalletRepository,
	logger *zerolog.Logger,
) SignatureVerificationService {
	return &signatureVerificationService{
		providerRegistry: providerRegistry,
		walletRepo:       walletRepo,
		challenges:       make(map[string]*Challenge),
		logger:           logger,
	}
}

// GenerateChallenge generates a challenge message for a wallet to sign
func (s *signatureVerificationService) GenerateChallenge(ctx context.Context, walletID string) (string, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return "", err
	}
	if wallet == nil {
		return "", errors.New("wallet not found")
	}

	// Generate a random challenge
	challengeBytes := make([]byte, 32)
	if _, err := rand.Read(challengeBytes); err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate random challenge")
		return "", err
	}
	challengeStr := base64.StdEncoding.EncodeToString(challengeBytes)

	// Create a message to sign
	message := fmt.Sprintf("Sign this message to verify your wallet ownership: %s\nWallet ID: %s\nTimestamp: %d",
		challengeStr, walletID, time.Now().Unix())

	// Store the challenge
	s.challenges[walletID] = &Challenge{
		WalletID:  walletID,
		Message:   message,
		ExpiresAt: time.Now().Add(15 * time.Minute), // Challenge expires in 15 minutes
	}

	return message, nil
}

// VerifySignature verifies a signature against a challenge
func (s *signatureVerificationService) VerifySignature(ctx context.Context, walletID, challenge, signature string) (bool, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return false, err
	}
	if wallet == nil {
		return false, errors.New("wallet not found")
	}

	// Check if challenge exists and is valid
	storedChallenge, ok := s.challenges[walletID]
	if !ok {
		return false, errors.New("no challenge found for wallet")
	}
	if time.Now().After(storedChallenge.ExpiresAt) {
		delete(s.challenges, walletID)
		return false, errors.New("challenge expired")
	}
	if storedChallenge.Message != challenge {
		return false, errors.New("challenge mismatch")
	}

	// Get provider
	var providerName string
	var address string
	if wallet.Type == model.WalletTypeExchange {
		providerName = wallet.Exchange
		address = wallet.Exchange // For exchange wallets, address is not applicable
	} else if wallet.Type == model.WalletTypeWeb3 {
		providerName = wallet.Metadata.Network
		address = wallet.Metadata.Address
	} else {
		return false, errors.New("unsupported wallet type")
	}

	provider, err := s.providerRegistry.GetProvider(providerName)
	if err != nil {
		s.logger.Error().Err(err).Str("provider", providerName).Msg("Failed to get provider")
		return false, err
	}

	// Verify signature
	verified, err := provider.Verify(ctx, address, challenge, signature)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to verify signature")
		return false, err
	}

	// If verified, update wallet status
	if verified {
		wallet.Status = model.WalletStatusVerified
		wallet.LastUpdated = time.Now()
		if err := s.walletRepo.Save(ctx, wallet); err != nil {
			s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to update wallet verification status")
			return false, err
		}

		// Clean up challenge
		delete(s.challenges, walletID)
	}

	return verified, nil
}

// GetWalletStatus gets the verification status of a wallet
func (s *signatureVerificationService) GetWalletStatus(ctx context.Context, walletID string) (model.WalletStatus, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return "", err
	}
	if wallet == nil {
		return "", errors.New("wallet not found")
	}

	return wallet.Status, nil
}

// SetWalletStatus sets the verification status of a wallet
func (s *signatureVerificationService) SetWalletStatus(ctx context.Context, walletID string, status model.WalletStatus) error {
	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Update status
	wallet.Status = status
	wallet.LastUpdated = time.Now()

	// Save wallet
	return s.walletRepo.Save(ctx, wallet)
}
