package port

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// WalletCache defines the interface for wallet data caching
// Used to store/retrieve wallet data by unique key (userID+walletID or similar)
type WalletCache interface {
	Set(key string, wallet *model.Wallet, ttl ...time.Duration)
	Get(key string) (*model.Wallet, bool)
	Delete(key string)
	Clear()
	Keys() []string
	IsExpired(key string) bool
}
