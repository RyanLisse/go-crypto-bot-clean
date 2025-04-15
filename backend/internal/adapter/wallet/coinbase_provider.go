package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// CoinbaseProvider implements the ExchangeWalletProvider interface for Coinbase
// This is a scaffold for Coinbase integration

type CoinbaseProvider struct {
	*BaseProvider
	coinbaseClient port.CoinbaseClient
	apiKey     string
	apiSecret  string
}

func NewCoinbaseProvider(coinbaseClient port.CoinbaseClient, logger *zerolog.Logger) port.ExchangeWalletProvider {
	return &CoinbaseProvider{
		BaseProvider:   NewBaseProvider("Coinbase", model.WalletTypeExchange, logger),
		coinbaseClient: coinbaseClient,
	}
}

func (p *CoinbaseProvider) SetAPICredentials(ctx context.Context, apiKey, apiSecret string) error {
	p.apiKey = apiKey
	p.apiSecret = apiSecret
	return nil
}

func (p *CoinbaseProvider) Connect(ctx context.Context, params map[string]interface{}) (*model.Wallet, error) {
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
	if err := p.SetAPICredentials(ctx, apiKey, apiSecret); err != nil {
		return nil, err
	}
	account, err := p.coinbaseClient.GetAccount(ctx)
	if err != nil {
		p.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get account from Coinbase")
		return nil, err
	}
	wallet := model.NewExchangeWallet(userID, "Coinbase")
	wallet.Balances = account.Balances
	wallet.TotalUSDValue = account.TotalUSDValue
	wallet.LastSyncAt = time.Now()
	wallet.LastUpdated = time.Now()
	wallet.SetMetadata("Coinbase Exchange Wallet", "Connected via API", []string{"exchange", "coinbase"})
	return wallet, nil
}

func (p *CoinbaseProvider) Disconnect(ctx context.Context, walletID string) error {
	p.apiKey = ""
	p.apiSecret = ""
	return nil
}

func (p *CoinbaseProvider) Verify(ctx context.Context, address string, message string, signature string) (bool, error) {
	_, err := p.coinbaseClient.GetAccount(ctx)
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to verify Coinbase wallet")
		return false, err
	}
	return true, nil
}

func (p *CoinbaseProvider) GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	if wallet.Exchange != "Coinbase" {
		return nil, errors.New("not a Coinbase wallet")
	}
	account, err := p.coinbaseClient.GetAccount(ctx)
	if err != nil {
		p.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to get account from Coinbase")
		return nil, err
	}
	wallet.Balances = account.Balances
	wallet.TotalUSDValue = account.TotalUSDValue
	wallet.LastSyncAt = time.Now()
	wallet.LastUpdated = time.Now()
	return wallet, nil
}

func (p *CoinbaseProvider) GetName() string {
	return p.name
}

func (p *CoinbaseProvider) GetType() model.WalletType {
	return p.typ
}

func (p *CoinbaseProvider) IsValidAddress(ctx context.Context, address string) (bool, error) {
	return false, errors.New("not applicable for exchange wallets")
}
