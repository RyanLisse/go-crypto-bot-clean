package dto

import (
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/google/uuid"
)

// --- Input DTOs ---

// CreateWeb3WalletRequest defines the input for creating a Web3 wallet.
type CreateWeb3WalletRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Address   string    `json:"address" validate:"required,eth_addr"` // Example validation tag
	Network   string    `json:"network" validate:"required"`
	IsPrimary bool      `json:"is_primary"`
}

// CreateExchangeWalletRequest defines the input for creating an exchange wallet.
type CreateExchangeWalletRequest struct {
	UserID       uuid.UUID `json:"user_id" validate:"required"`
	ExchangeID   string    `json:"exchange_id" validate:"required"`
	ExchangeName string    `json:"exchange_name" validate:"required"`
	IsPrimary    bool      `json:"is_primary"`
	// APIKey    string    `json:"api_key" validate:"required"` // If handling keys directly
	// APISecret string    `json:"api_secret" validate:"required"`
}

// UpdateBalancesRequest defines the input for updating wallet balances.
type UpdateBalancesRequest struct {
	Balances []*model.Balance `json:"balances" validate:"required,dive"` // Validate nested elements
}

// --- Output DTOs ---

// WalletResponse defines the standard response structure for a wallet.
type WalletResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Type      model.WalletType   `json:"type"`
	Status    model.WalletStatus `json:"status"`
	IsPrimary bool               `json:"is_primary"`

	Address      string `json:"address,omitempty"`
	Network      string `json:"network,omitempty"`
	ExchangeID   string `json:"exchange_id,omitempty"`
	ExchangeName string `json:"exchange_name,omitempty"`

	Balances   []*BalanceResponse `json:"balances"`
	VerifiedAt *time.Time         `json:"verified_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// BalanceResponse defines the response structure for a wallet balance.
type BalanceResponse struct {
	Token  string  `json:"token"`
	Amount string  `json:"amount"`
	Value  float64 `json:"value"`
}

// --- Mapping Functions (Optional but Recommended) ---

// ToWalletResponse converts a domain Wallet model to a WalletResponse DTO.
func ToWalletResponse(w *model.Wallet) *WalletResponse {
	if w == nil {
		return nil
	}
	balancesDTO := make([]*BalanceResponse, len(w.Balances))
	for i, b := range w.Balances {
		balancesDTO[i] = ToBalanceResponse(b)
	}
	return &WalletResponse{
		ID:           w.ID,
		UserID:       w.UserID,
		Type:         w.Type,
		Status:       w.Status,
		IsPrimary:    w.IsPrimary,
		Address:      w.Address,
		Network:      w.Network,
		ExchangeID:   w.ExchangeID,
		ExchangeName: w.ExchangeName,
		Balances:     balancesDTO,
		VerifiedAt:   w.VerifiedAt,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}

// ToBalanceResponse converts a domain Balance model to a BalanceResponse DTO.
func ToBalanceResponse(b *model.Balance) *BalanceResponse {
	if b == nil {
		return nil
	}
	return &BalanceResponse{
		Token:  string(b.Asset),
		Amount: fmt.Sprintf("%.8f", b.Free),
		Value:  b.USDValue,
	}
}

// ToListWalletResponse converts a slice of domain Wallet models to WalletResponse DTOs.
func ToListWalletResponse(wallets []*model.Wallet) []*WalletResponse {
	list := make([]*WalletResponse, len(wallets))
	for i, w := range wallets {
		list[i] = ToWalletResponse(w)
	}
	return list
}
