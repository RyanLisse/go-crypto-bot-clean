package websocket

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"go.uber.org/zap"
)

// AccountUpdateCallback is a function that is called when an account update is received
type AccountUpdateCallback func(*models.Wallet)

// AccountHandler manages WebSocket account update subscriptions
type AccountHandler struct {
	client           *Client
	logger           *zap.Logger
	accountCallbacks []AccountUpdateCallback
	mutex            sync.RWMutex
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(client *Client, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{
		client:           client,
		logger:           logger,
		accountCallbacks: make([]AccountUpdateCallback, 0),
	}
}

// SubscribeToAccountUpdates registers a callback for account updates
func (h *AccountHandler) SubscribeToAccountUpdates(ctx context.Context, callback AccountUpdateCallback) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.logger.Debug("Registering account update callback")
	h.accountCallbacks = append(h.accountCallbacks, callback)
}

// UnsubscribeFromAccountUpdates clears all registered callbacks
func (h *AccountHandler) UnsubscribeFromAccountUpdates() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.logger.Debug("Clearing all account update callbacks")
	h.accountCallbacks = make([]AccountUpdateCallback, 0)
}

// HandleAccountUpdate processes account update messages from WebSocket
func (h *AccountHandler) HandleAccountUpdate(msg map[string]interface{}) {
	h.logger.Debug("Processing account update message")

	// Extract data from the message
	data, ok := msg["data"].(map[string]interface{})
	if !ok {
		h.logger.Error("Invalid account update message format")
		return
	}

	// Create a new wallet with the updated balances
	wallet := &models.Wallet{
		Balances:  make(map[string]*models.AssetBalance),
		UpdatedAt: time.Now(),
	}

	// Extract balances from the message
	if balances, ok := data["balances"].([]interface{}); ok {
		h.logger.Debug("Processing account balances", zap.Int("count", len(balances)))
		for _, bal := range balances {
			balance, ok := bal.(map[string]interface{})
			if !ok {
				continue
			}

			asset, ok := balance["asset"].(string)
			if !ok {
				continue
			}

			free := parseFloat(balance["free"])
			locked := parseFloat(balance["locked"])

			wallet.Balances[asset] = &models.AssetBalance{
				Asset:  asset,
				Free:   free,
				Locked: locked,
				Total:  free + locked,
			}
		}
	} else {
		h.logger.Warn("No balances found in account update message")
	}

	h.logger.Debug("Account update processed",
		zap.Time("updated_at", wallet.UpdatedAt),
		zap.Int("asset_count", len(wallet.Balances)))

	// Notify callbacks
	h.notifyCallbacks(wallet)
}

// notifyCallbacks sends the wallet update to all registered callbacks
func (h *AccountHandler) notifyCallbacks(wallet *models.Wallet) {
	h.mutex.RLock()
	callbacks := make([]AccountUpdateCallback, len(h.accountCallbacks))
	copy(callbacks, h.accountCallbacks)
	callbackCount := len(callbacks)
	h.mutex.RUnlock()

	h.logger.Debug("Notifying account update callbacks", zap.Int("callback_count", callbackCount))

	for _, callback := range callbacks {
		go callback(wallet)
	}
}

// ProcessMessage checks if a message is an account update and handles it
func (h *AccountHandler) ProcessMessage(msg map[string]interface{}) bool {
	// Check if it's an account update message
	if channel, ok := msg["channel"].(string); ok {
		if strings.Contains(channel, "spot@private.account.v3.api") {
			h.logger.Debug("Received account update message", zap.String("channel", channel))
			h.HandleAccountUpdate(msg)
			return true
		}
	}

	return false
}
