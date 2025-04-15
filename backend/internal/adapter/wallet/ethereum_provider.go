package wallet

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
)

// EthereumProvider implements the Web3WalletProvider interface for Ethereum
type EthereumProvider struct {
	*BaseProvider
	chainID  int64
	network  string
	rpcURL   string
	explorer string
}

// NewEthereumProvider creates a new Ethereum wallet provider
func NewEthereumProvider(chainID int64, network, rpcURL, explorer string, logger *zerolog.Logger) port.Web3WalletProvider {
	return &EthereumProvider{
		BaseProvider: NewBaseProvider("Ethereum", model.WalletTypeWeb3, logger),
		chainID:      chainID,
		network:      network,
		rpcURL:       rpcURL,
		explorer:     explorer,
	}
}

// GetChainID returns the chain ID for the provider
func (p *EthereumProvider) GetChainID() int64 {
	return p.chainID
}

// GetNetwork returns the network for the provider
func (p *EthereumProvider) GetNetwork() string {
	return p.network
}

// Connect connects to an Ethereum wallet
func (p *EthereumProvider) Connect(ctx context.Context, params map[string]interface{}) (*model.Wallet, error) {
	// Extract parameters
	userID, ok := params["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("user_id is required")
	}

	address, ok := params["address"].(string)
	if !ok || address == "" {
		return nil, errors.New("address is required")
	}

	// Validate address
	valid, err := p.IsValidAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid Ethereum address")
	}

	// Create wallet
	wallet := model.NewWeb3Wallet(userID, p.network, address)
	wallet.LastSyncAt = time.Now()
	wallet.LastUpdated = time.Now()

	// Set metadata
	wallet.SetMetadata("Ethereum Wallet", "Connected via Web3", []string{"web3", "ethereum"})
	wallet.Metadata.Network = p.network
	wallet.Metadata.Address = address
	wallet.Metadata.ChainID = p.chainID
	wallet.Metadata.Explorer = p.explorer
	wallet.Metadata.Custom = make(map[string]string)
	wallet.Metadata.Custom["rpc_url"] = p.rpcURL

	return wallet, nil
}

// SignMessage signs a message with the wallet's private key
// Note: This is a placeholder. In a real implementation, this would be done client-side
func (p *EthereumProvider) SignMessage(ctx context.Context, message string) (string, error) {
	return "", errors.New("signing must be done client-side")
}

// Verify verifies a wallet connection using a signature
func (p *EthereumProvider) Verify(ctx context.Context, address string, message string, signature string) (bool, error) {
	// Validate address
	valid, err := p.IsValidAddress(ctx, address)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("invalid Ethereum address")
	}

	// Verify signature
	// The message is prefixed with "\x19Ethereum Signed Message:\n" + len(message) to prevent
	// malicious DApps from using the signature to perform contract calls
	prefixedMessage := "\x19Ethereum Signed Message:\n" + fmt.Sprintf("%d", len(message)) + message
	messageHash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// Convert signature to bytes
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	// The signature should be 65 bytes: R (32 bytes) + S (32 bytes) + V (1 byte)
	if len(signatureBytes) != 65 {
		return false, errors.New("invalid signature length")
	}

	// Adjust V value (last byte) if needed
	if signatureBytes[64] < 27 {
		signatureBytes[64] += 27
	}

	// Recover the public key from the signature
	sigPublicKey, err := crypto.Ecrecover(messageHash.Bytes(), signatureBytes)
	if err != nil {
		return false, err
	}

	// Convert the public key to an Ethereum address
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		return false, err
	}
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	// Compare the recovered address with the provided address
	return strings.EqualFold(recoveredAddress, address), nil
}

// GetBalance gets the balance for a wallet
// Note: This is a placeholder. In a real implementation, this would query the Ethereum blockchain
func (p *EthereumProvider) GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	// Check if wallet is an Ethereum wallet
	if wallet.Type != model.WalletTypeWeb3 || wallet.Metadata.Network != p.network {
		return nil, errors.New("not an Ethereum wallet")
	}

	// In a real implementation, we would query the Ethereum blockchain for the balance
	// For now, we'll just return the wallet as is
	wallet.LastSyncAt = time.Now()
	wallet.LastUpdated = time.Now()

	return wallet, nil
}

// Disconnect disconnects from the Ethereum wallet
func (p *EthereumProvider) Disconnect(ctx context.Context, walletID string) error {
	// For Ethereum wallets, there's no real "disconnection" - we just remove the wallet from our database
	return nil
}

// IsValidAddress checks if an address is valid for this provider
func (p *EthereumProvider) IsValidAddress(ctx context.Context, address string) (bool, error) {
	// Check if the address is a valid Ethereum address
	// Ethereum addresses are 42 characters long, starting with "0x" followed by 40 hexadecimal characters
	if !common.IsHexAddress(address) {
		return false, nil
	}

	// Additional validation: check if the address is checksummed
	// This is optional but recommended
	match, _ := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address)
	return match, nil
}

// Ensure EthereumProvider implements port.Web3WalletProvider
var _ port.Web3WalletProvider = (*EthereumProvider)(nil)
