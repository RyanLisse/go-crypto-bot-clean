package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MEXCProvider implements the ExchangeWalletProvider interface for MEXC
type MEXCProvider struct {
	*BaseProvider
	mexcClient port.MEXCClient
	apiKey     string
	apiSecret  string
}

// NewMEXCProvider creates a new MEXC wallet provider
func NewMEXCProvider(mexcClient port.MEXCClient, logger *zerolog.Logger) port.ExchangeWalletProvider {
	return &MEXCProvider{
		BaseProvider: NewBaseProvider("MEXC", model.WalletTypeExchange, logger),
		mexcClient:   mexcClient,
	}
}

// SetAPICredentials sets the API credentials for the exchange
func (p *MEXCProvider) SetAPICredentials(ctx context.Context, apiKey, apiSecret string) error {
	p.apiKey = apiKey
	p.apiSecret = apiSecret
	return nil
}

// Connect connects to the MEXC exchange
func (p *MEXCProvider) Connect(ctx context.Context, params map[string]interface{}) (*model.Wallet, error) {
	// Extract parameters
	userID, ok := params["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("user_id is required")
	}

	apiKey, ok := params["api_key"].(string)
	if !ok || apiKey == "" {
		return nil, errors.New("api_key is required")
	}

	apiSecret, ok := params["api_secret"].(string)
	if !ok || apiSecret == "" {
		return nil, errors.New("api_secret is required")
	}

	// Set API credentials
	if err := p.SetAPICredentials(ctx, apiKey, apiSecret); err != nil {
		return nil, err
	}

	// Get account from MEXC
	account, err := p.mexcClient.GetAccount(ctx)
	if err != nil {
		p.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get account from MEXC")
		return nil, err
	}

	// Create wallet
	wallet := model.NewExchangeWallet(userID, "MEXC")
	wallet.Balances = account.Balances
	wallet.TotalUSDValue = account.TotalUSDValue
	wallet.LastSyncAt = time.Now()
	wallet.LastUpdated = time.Now()

	// Set metadata
	wallet.SetMetadata("MEXC Exchange Wallet", "Connected via API", []string{"exchange", "mexc"})

	return wallet, nil
}

// GetBalance gets the balance for a wallet
func (p *MEXCProvider) GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	// Check if wallet is a MEXC wallet
	if wallet.Exchange != "MEXC" {
		return nil, errors.New("not a MEXC wallet")
	}

	// Get account from MEXC
	account, err := p.mexcClient.GetAccount(ctx)
	if err != nil {
		p.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to get account from MEXC")
		return nil, err
	}

	// Update wallet balances
	wallet.Balances = account.Balances
	wallet.TotalUSDValue = account.TotalUSDValue
	wallet.LastSyncAt = time.Now()
	wallet.LastUpdated = time.Now()

	return wallet, nil
}

// Verify verifies a wallet connection using API credentials
func (p *MEXCProvider) Verify(ctx context.Context, address string, message string, signature string) (bool, error) {
	// For exchange wallets, we verify by checking if we can get the account
	_, err := p.mexcClient.GetAccount(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to verify MEXC wallet")
		return false, err
	}

	return true, nil
}

// Disconnect disconnects from the MEXC exchange
func (p *MEXCProvider) Disconnect(ctx context.Context, walletID string) error {
	// For MEXC, we just clear the API credentials
	p.apiKey = ""
	p.apiSecret = ""
	return nil
}

// IsValidAddress checks if an address is valid for this provider
// For exchange wallets, this is not applicable
func (p *MEXCProvider) IsValidAddress(ctx context.Context, address string) (bool, error) {
	return false, errors.New("not applicable for exchange wallets")
}

// Ensure MEXCProvider implements port.ExchangeWalletProvider
var _ port.ExchangeWalletProvider = (*MEXCProvider)(nil)
