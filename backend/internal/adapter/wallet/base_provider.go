package wallet

import (
	"context"
	"errors"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// BaseProvider provides common functionality for wallet providers
type BaseProvider struct {
	name   string
	typ    model.WalletType
	logger *zerolog.Logger
}

// NewBaseProvider creates a new base wallet provider
func NewBaseProvider(name string, typ model.WalletType, logger *zerolog.Logger) *BaseProvider {
	return &BaseProvider{
		name:   name,
		typ:    typ,
		logger: logger,
	}
}

// GetName returns the name of the wallet provider
func (p *BaseProvider) GetName() string {
	return p.name
}

// GetType returns the type of wallet provider
func (p *BaseProvider) GetType() model.WalletType {
	return p.typ
}

// Connect connects to the wallet provider
func (p *BaseProvider) Connect(ctx context.Context, params map[string]interface{}) (*model.Wallet, error) {
	return nil, errors.New("not implemented")
}

// Disconnect disconnects from the wallet provider
func (p *BaseProvider) Disconnect(ctx context.Context, walletID string) error {
	return errors.New("not implemented")
}

// Verify verifies a wallet connection using a signature
func (p *BaseProvider) Verify(ctx context.Context, address string, message string, signature string) (bool, error) {
	return false, errors.New("not implemented")
}

// GetBalance gets the balance for a wallet
func (p *BaseProvider) GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	return nil, errors.New("not implemented")
}

// IsValidAddress checks if an address is valid for this provider
func (p *BaseProvider) IsValidAddress(ctx context.Context, address string) (bool, error) {
	return false, errors.New("not implemented")
}

// Ensure BaseProvider implements port.WalletProvider
var _ port.WalletProvider = (*BaseProvider)(nil)
