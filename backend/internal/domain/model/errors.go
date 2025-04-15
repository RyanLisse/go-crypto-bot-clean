package model

import "errors"

// Domain model errors
var (
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidExchange   = errors.New("invalid exchange")
	ErrInvalidAPIKey     = errors.New("invalid API key")
	ErrInvalidAPISecret  = errors.New("invalid API secret")
	ErrInvalidWalletID   = errors.New("invalid wallet ID")
	ErrInvalidAsset      = errors.New("invalid asset")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrCredentialNotFound = errors.New("credential not found")
)
