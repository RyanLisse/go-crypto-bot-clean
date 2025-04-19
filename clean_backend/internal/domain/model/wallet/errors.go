package wallet

import "errors"

var (
	// ErrMissingUserID is returned when a wallet is created without a user ID
	ErrMissingUserID = errors.New("user ID is required")

	// ErrMissingAddress is returned when a Web3 wallet is created without an address
	ErrMissingAddress = errors.New("address is required for Web3 wallets")

	// ErrMissingNetwork is returned when a Web3 wallet is created without a network
	ErrMissingNetwork = errors.New("network is required for Web3 wallets")

	// ErrMissingExchangeID is returned when an exchange wallet is created without an exchange ID
	ErrMissingExchangeID = errors.New("exchange ID is required for exchange wallets")

	// ErrInvalidWalletType is returned when an invalid wallet type is specified
	ErrInvalidWalletType = errors.New("invalid wallet type")

	// ErrWalletNotFound is returned when a wallet cannot be found
	ErrWalletNotFound = errors.New("wallet not found")

	// ErrDuplicateWallet is returned when attempting to create a duplicate wallet
	ErrDuplicateWallet = errors.New("wallet already exists")

	// ErrInvalidAddress is returned when an invalid wallet address is provided
	ErrInvalidAddress = errors.New("invalid wallet address")

	// ErrInvalidNetwork is returned when an invalid network is specified
	ErrInvalidNetwork = errors.New("invalid network")

	// ErrInvalidExchangeID is returned when an invalid exchange ID is provided
	ErrInvalidExchangeID = errors.New("invalid exchange ID")

	// ErrVerificationFailed is returned when wallet ownership verification fails
	ErrVerificationFailed = errors.New("wallet ownership verification failed")

	// ErrInvalidSignature is returned when an invalid signature is provided for verification
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrInvalidChallenge is returned when an invalid or expired challenge is provided
	ErrInvalidChallenge = errors.New("invalid or expired challenge")
)
