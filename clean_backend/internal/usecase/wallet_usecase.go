package usecase

import (
	"context"
	"errors" // Use domain errors or apperror package

	"github.com/rs/zerolog"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/google/uuid"
)

// Ensure walletUseCase implements the domain service interface
var _ port.WalletService = (*walletUseCase)(nil)

type walletUseCase struct {
	walletRepo port.WalletRepository
	logger     *zerolog.Logger
	// transactionMgr port.TransactionManager // Optional: if complex operations need transactions
	// Add other dependencies like notification service, verification gateway etc.
}

// NewWalletUseCase creates a new wallet use case/service instance.
func NewWalletUseCase(repo port.WalletRepository, logger *zerolog.Logger /*, txMgr port.TransactionManager*/) port.WalletService {
	return &walletUseCase{
		walletRepo: repo,
		logger:     logger,
		// transactionMgr: txMgr,
	}
}

// CreateWeb3Wallet handles the creation of a new Web3 wallet.
func (uc *walletUseCase) CreateWeb3Wallet(ctx context.Context, userID uuid.UUID, address, network string, isPrimary bool) (*model.Wallet, error) {
	// TODO: Add validation for address format based on network?

	// Check for existing primary wallet if needed (could be domain service logic)
	if isPrimary {
		primary, err := uc.walletRepo.FindPrimary(ctx, userID, model.WalletTypeWeb3)
		if err != nil && !errors.Is(err, model.ErrWalletNotFound) {
			// TODO: Wrap error (apperror?)
			return nil, err
		}
		if primary != nil {
			// TODO: Handle error - primary already exists (apperror?)
			return nil, errors.New("web3 primary wallet already exists for this user")
		}
	}

	newWallet := model.NewWeb3Wallet(userID, address, network)
	newWallet.SetPrimary(isPrimary)

	// Persist the new wallet
	err := uc.walletRepo.Create(ctx, newWallet)
	if err != nil {
		// TODO: Wrap/handle specific repo errors (e.g., duplicate address)
		return nil, err
	}

	// Optionally set primary status transactionally if needed
	if isPrimary {
		err = uc.walletRepo.SetPrimary(ctx, userID, newWallet.ID, newWallet.Type)
		if err != nil {
			// TODO: Log this error? Potential inconsistency if Create worked but SetPrimary failed.
			// Consider transactional approach using transactionMgr
			return nil, err
		}
	}

	// TODO: Trigger verification process if needed (async?)

	return newWallet, nil
}

// CreateExchangeWallet handles the creation of a new exchange wallet.
func (uc *walletUseCase) CreateExchangeWallet(ctx context.Context, userID uuid.UUID, exchangeID, exchangeName string, isPrimary bool) (*model.Wallet, error) {
	// TODO: Add validation if needed

	// Check for existing primary wallet if needed
	if isPrimary {
		primary, err := uc.walletRepo.FindPrimary(ctx, userID, model.WalletTypeExchange)
		if err != nil && !errors.Is(err, model.ErrWalletNotFound) {
			return nil, err // TODO: Wrap error
		}
		if primary != nil {
			return nil, errors.New("exchange primary wallet already exists for this user") // TODO: Wrap error
		}
	}

	newWallet := model.NewExchangeWallet(userID, exchangeID, exchangeName)
	newWallet.SetPrimary(isPrimary)

	err := uc.walletRepo.Create(ctx, newWallet)
	if err != nil {
		return nil, err // TODO: Wrap/handle error
	}

	if isPrimary {
		err = uc.walletRepo.SetPrimary(ctx, userID, newWallet.ID, newWallet.Type)
		if err != nil {
			return nil, err // TODO: Log, consider transaction
		}
	}

	// TODO: Trigger verification/sync process?

	return newWallet, nil
}

// GetWallet retrieves a single wallet by ID.
func (uc *walletUseCase) GetWallet(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	// Directly call repository, error handling might involve mapping repo errors to use case/app errors.
	w, err := uc.walletRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			return nil, model.ErrWalletNotFound // Or a specific use case error
		}
		return nil, err // TODO: Wrap error
	}
	return w, nil
}

// GetUserWallets retrieves all wallets for a given user.
func (uc *walletUseCase) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*model.Wallet, error) {
	wallets, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err // TODO: Wrap error
	}
	// Return empty slice, not nil, if no wallets are found
	if wallets == nil {
		return []*model.Wallet{}, nil
	}
	return wallets, nil
}

// GetPrimaryWallet retrieves the user's primary wallet (type needs clarification - Web3 or Exchange?).
// Assuming we need separate methods or a type parameter if both can be primary.
// This implementation fetches the first primary found, regardless of type.
func (uc *walletUseCase) GetPrimaryWallet(ctx context.Context, userID uuid.UUID) (*model.Wallet, error) {
	// Need clarification: Should this return primary Web3 or Exchange? Or just the first one found?
	// Let's assume we look for Web3 first, then Exchange if not found.
	w, err := uc.walletRepo.FindPrimary(ctx, userID, model.WalletTypeWeb3)
	if err == nil && w != nil {
		return w, nil // Found primary Web3
	}
	if err != nil && !errors.Is(err, model.ErrWalletNotFound) {
		return nil, err // Error checking Web3 primary
	}

	// If Web3 not found or error was ErrWalletNotFound, check for Exchange primary
	w, err = uc.walletRepo.FindPrimary(ctx, userID, model.WalletTypeExchange)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			return nil, model.ErrWalletNotFound // No primary wallet found at all
		}
		return nil, err // Error checking Exchange primary
	}
	return w, nil // Found primary Exchange
}

// SetPrimaryWallet sets a specific wallet as the primary for the user.
func (uc *walletUseCase) SetPrimaryWallet(ctx context.Context, walletID, userID uuid.UUID) error {
	// 1. Get the target wallet to know its type
	w, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			return model.ErrWalletNotFound
		}
		return err // TODO: Wrap
	}

	// 2. Check ownership (redundant if GetByID implicitly checks user? Depends on repo implementation)
	if w.UserID != userID {
		return errors.New("unauthorized to set primary status") // TODO: AppError
	}

	// 3. Use the repository's SetPrimary function which handles unsetting others of the same type.
	err = uc.walletRepo.SetPrimary(ctx, userID, walletID, w.Type)
	if err != nil {
		return err // TODO: Wrap
	}

	return nil
}

// UpdateBalances updates the balances for a specific wallet.
// Note: This use case assumes balances are provided externally.
// A real implementation might involve fetching from an exchange/node.
func (uc *walletUseCase) UpdateBalances(ctx context.Context, walletID uuid.UUID, balances []*model.Balance) error {
	// TODO: Validate input balances?

	// Check if wallet exists first (optional, repo update might handle not found)
	_, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			return model.ErrWalletNotFound
		}
		return err // TODO: Wrap
	}

	// Convert []*model.Balance to []model.Balance if repo expects value type
	valueBalances := make([]model.Balance, len(balances))
	for i, b := range balances {
		if b == nil {
			return errors.New("nil balance provided in update request") // TODO: AppError
		}
		valueBalances[i] = *b
	}

	err = uc.walletRepo.UpdateBalance(ctx, walletID, valueBalances)
	if err != nil {
		return err // TODO: Wrap
	}
	return nil
}

// VerifyWallet is a placeholder - implementation depends on verification method.
func (uc *walletUseCase) VerifyWallet(ctx context.Context, walletID uuid.UUID) error {
	// 1. Get wallet
	// 2. Determine type (Web3/Exchange)
	// 3. Call appropriate verification service/gateway
	// 4. Update status in repository
	_, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			uc.logger.Error().Str("walletID", walletID.String()).Msg("Wallet not found during verification")
			return model.ErrWalletNotFound
		}
		uc.logger.Error().Err(err).Str("walletID", walletID.String()).Msg("Error retrieving wallet for verification")
		return err
	}
	// Placeholder
	uc.logger.Info().Str("walletID", walletID.String()).Msg("VerifyWallet use case called")
	return errors.New("VerifyWallet not implemented")
}

// CheckVerificationStatus is a placeholder.
func (uc *walletUseCase) CheckVerificationStatus(ctx context.Context, walletID uuid.UUID) (model.VerificationStatus, error) {
	// 1. Get wallet status from repository
	// 2. Optionally query external verification service
	// 3. Update repo if needed
	w, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			uc.logger.Error().Str("walletID", walletID.String()).Msg("Wallet not found when checking verification status")
			return model.VerificationStatusNone, model.ErrWalletNotFound
		}
		uc.logger.Error().Err(err).Str("walletID", walletID.String()).Msg("Error retrieving wallet for verification status check")
		return "", err
	}

	uc.logger.Info().Str("walletID", walletID.String()).Msg("Checking wallet verification status")

	// Map wallet status to verification status
	var verificationStatus model.VerificationStatus
	switch w.Status {
	case model.WalletStatusVerified:
		verificationStatus = model.VerificationStatusVerified
	case model.WalletStatusPending:
		verificationStatus = model.VerificationStatusPending
	default:
		verificationStatus = model.VerificationStatusNone
	}

	return verificationStatus, nil
}

// DeleteWallet handles deleting a wallet, ensuring ownership.
func (uc *walletUseCase) DeleteWallet(ctx context.Context, walletID, userID uuid.UUID) error {
	// 1. Get wallet to verify ownership
	w, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			return model.ErrWalletNotFound // Or return nil if delete is idempotent?
		}
		return err // TODO: Wrap
	}

	// 2. Check ownership
	if w.UserID != userID {
		return errors.New("unauthorized to delete wallet") // TODO: AppError
	}

	// 3. Call repository delete
	err = uc.walletRepo.Delete(ctx, walletID)
	if err != nil {
		return err // TODO: Wrap
	}

	return nil
}
