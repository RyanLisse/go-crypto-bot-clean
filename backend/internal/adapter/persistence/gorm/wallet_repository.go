package gorm

import (
	"context"
	"strconv"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure WalletRepository implements the port.WalletRepository interface
var _ port.WalletRepository = (*WalletRepository)(nil)

// WalletEntity represents the database model for wallet
type WalletEntity struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        string    `gorm:"size:50;not null;index"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null"`
	LastUpdated   time.Time `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// BalanceEntity represents the database model for a balance
type BalanceEntity struct {
	ID        uint    `gorm:"primaryKey"`
	WalletID  uint    `gorm:"not null;index"`
	Asset     string  `gorm:"size:20;not null"`
	Free      float64 `gorm:"type:decimal(18,8);not null"`
	Locked    float64 `gorm:"type:decimal(18,8);not null"`
	Total     float64 `gorm:"type:decimal(18,8);not null"`
	USDValue  float64 `gorm:"type:decimal(18,8);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BalanceHistoryEntity represents the database model for balance history
type BalanceHistoryEntity struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"size:50;not null;index"`
	Asset     string    `gorm:"size:20;not null"`
	Free      float64   `gorm:"type:decimal(18,8);not null"`
	Locked    float64   `gorm:"type:decimal(18,8);not null"`
	Total     float64   `gorm:"type:decimal(18,8);not null"`
	USDValue  float64   `gorm:"type:decimal(18,8);not null"`
	Timestamp time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}

// WalletRepository implements the port.WalletRepository interface with GORM
type WalletRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewWalletRepository creates a new WalletRepository
func NewWalletRepository(db *gorm.DB, logger *zerolog.Logger) *WalletRepository {
	return &WalletRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists a wallet to the database
func (r *WalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	r.logger.Debug().
		Str("userID", wallet.UserID).
		Float64("totalUSDValue", wallet.TotalUSDValue).
		Msg("Saving wallet to database")

	// Begin transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find or create wallet entity
	var walletEntity WalletEntity
	result := tx.Where("user_id = ?", wallet.UserID).First(&walletEntity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new wallet entity
			walletEntity = WalletEntity{
				UserID:        wallet.UserID,
				TotalUSDValue: wallet.TotalUSDValue,
				LastUpdated:   wallet.LastUpdated,
			}
			if err := tx.Create(&walletEntity).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return result.Error
		}
	} else {
		// Update existing wallet entity
		walletEntity.TotalUSDValue = wallet.TotalUSDValue
		walletEntity.LastUpdated = wallet.LastUpdated
		if err := tx.Save(&walletEntity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete old balances
	if err := tx.Where("wallet_id = ?", walletEntity.ID).Delete(&BalanceEntity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new balances
	for asset, balance := range wallet.Balances {
		balanceEntity := BalanceEntity{
			WalletID: walletEntity.ID,
			Asset:    string(asset),
			Free:     balance.Free,
			Locked:   balance.Locked,
			Total:    balance.Total,
			USDValue: balance.USDValue,
		}
		if err := tx.Create(&balanceEntity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

// GetByUserID retrieves a wallet by user ID
func (r *WalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	r.logger.Debug().
		Str("userID", userID).
		Msg("Getting wallet by user ID")

	// Find wallet entity
	var walletEntity WalletEntity
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&walletEntity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return empty wallet if not found
			return model.NewWallet(userID), nil
		}
		return nil, result.Error
	}

	// Find balances for this wallet
	var balanceEntities []BalanceEntity
	if err := r.db.WithContext(ctx).Where("wallet_id = ?", walletEntity.ID).Find(&balanceEntities).Error; err != nil {
		return nil, err
	}

	// Create domain wallet
	wallet := &model.Wallet{
		UserID:        walletEntity.UserID,
		Balances:      make(map[model.Asset]*model.Balance),
		TotalUSDValue: walletEntity.TotalUSDValue,
		LastUpdated:   walletEntity.LastUpdated,
	}

	// Add balances
	for _, entity := range balanceEntities {
		asset := model.Asset(entity.Asset)
		wallet.Balances[asset] = &model.Balance{
			Asset:    asset,
			Free:     entity.Free,
			Locked:   entity.Locked,
			Total:    entity.Total,
			USDValue: entity.USDValue,
		}
	}

	return wallet, nil
}

// SaveBalanceHistory saves a balance history record
func (r *WalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	r.logger.Debug().
		Str("userID", history.UserID).
		Str("asset", string(history.Asset)).
		Time("timestamp", history.Timestamp).
		Msg("Saving balance history")

	entity := BalanceHistoryEntity{
		UserID:    history.UserID,
		Asset:     string(history.Asset),
		Free:      history.Free,
		Locked:    history.Locked,
		Total:     history.Total,
		USDValue:  history.USDValue,
		Timestamp: history.Timestamp,
	}

	return r.db.WithContext(ctx).Create(&entity).Error
}

// GetBalanceHistory retrieves balance history records
func (r *WalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	r.logger.Debug().
		Str("userID", userID).
		Str("asset", string(asset)).
		Time("from", from).
		Time("to", to).
		Msg("Getting balance history")

	var entities []BalanceHistoryEntity
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND asset = ? AND timestamp BETWEEN ? AND ?", userID, string(asset), from, to).
		Order("timestamp ASC").
		Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}

	history := make([]*model.BalanceHistory, len(entities))
	for i, entity := range entities {
		history[i] = &model.BalanceHistory{
			ID:        strconv.FormatUint(uint64(entity.ID), 10),
			UserID:    entity.UserID,
			Asset:     model.Asset(entity.Asset),
			Free:      entity.Free,
			Locked:    entity.Locked,
			Total:     entity.Total,
			USDValue:  entity.USDValue,
			Timestamp: entity.Timestamp,
		}
	}

	return history, nil
}
