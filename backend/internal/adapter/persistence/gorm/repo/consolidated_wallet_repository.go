package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// EnhancedWalletEntity represents a wallet in the database
type EnhancedWalletEntity struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	UserID        string    `gorm:"index;type:varchar(50);not null"`
	Exchange      string    `gorm:"index;type:varchar(50)"`
	Type          string    `gorm:"index;type:varchar(20);not null"`
	Status        string    `gorm:"index;type:varchar(20);not null"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Metadata      []byte    `gorm:"type:json"`
	LastUpdated   time.Time `gorm:"not null"`
	LastSyncAt    time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for EnhancedWalletEntity
func (EnhancedWalletEntity) TableName() string {
	return "enhanced_wallets"
}

// EnhancedWalletBalanceEntity represents a balance in the database
type EnhancedWalletBalanceEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	WalletID  string    `gorm:"index;type:varchar(50);not null"`
	Asset     string    `gorm:"index;type:varchar(20);not null"`
	Free      float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Locked    float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Total     float64   `gorm:"type:decimal(18,8);not null;default:0"`
	USDValue  float64   `gorm:"type:decimal(18,8);not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for EnhancedWalletBalanceEntity
func (EnhancedWalletBalanceEntity) TableName() string {
	return "enhanced_wallet_balances"
}

// EnhancedWalletBalanceHistoryEntity represents a balance history record in the database
type EnhancedWalletBalanceHistoryEntity struct {
	ID            string    `gorm:"primaryKey;type:varchar(50)"`
	UserID        string    `gorm:"index;type:varchar(50);not null"`
	WalletID      string    `gorm:"index;type:varchar(50);not null"`
	BalancesJSON  []byte    `gorm:"type:json"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null;default:0"`
	Timestamp     time.Time `gorm:"index;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// TableName sets the table name for EnhancedWalletBalanceHistoryEntity
func (EnhancedWalletBalanceHistoryEntity) TableName() string {
	return "enhanced_wallet_balance_history"
}

// ConsolidatedWalletRepository implements port.WalletRepository using GORM
type ConsolidatedWalletRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewConsolidatedWalletRepository creates a new ConsolidatedWalletRepository
func NewConsolidatedWalletRepository(db *gorm.DB, logger *zerolog.Logger) port.WalletRepository {
	return &ConsolidatedWalletRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists a wallet to the database
func (r *ConsolidatedWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	r.logger.Debug().
		Str("userID", wallet.UserID).
		Str("id", wallet.ID).
		Str("exchange", wallet.Exchange).
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

	// Convert metadata to JSON
	var metadataJSON []byte
	var err error
	if wallet.Metadata != nil {
		metadataJSON, err = json.Marshal(wallet.Metadata)
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to marshal wallet metadata")
			tx.Rollback()
			return err
		}
	}

	// Create or update wallet entity
	walletEntity := EnhancedWalletEntity{
		ID:            wallet.ID,
		UserID:        wallet.UserID,
		Exchange:      wallet.Exchange,
		Type:          string(wallet.Type),
		Status:        string(wallet.Status),
		TotalUSDValue: wallet.TotalUSDValue,
		Metadata:      metadataJSON,
		LastUpdated:   wallet.LastUpdated,
		LastSyncAt:    wallet.LastSyncAt,
	}

	// Check if wallet exists
	var existingWallet EnhancedWalletEntity
	result := tx.Where("id = ?", wallet.ID).First(&existingWallet)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		r.logger.Error().Err(result.Error).Str("id", wallet.ID).Msg("Error checking if wallet exists")
		tx.Rollback()
		return result.Error
	}

	// Create or update wallet
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new wallet
		if err := tx.Create(&walletEntity).Error; err != nil {
			r.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to create wallet")
			tx.Rollback()
			return err
		}
	} else {
		// Update existing wallet
		if err := tx.Model(&walletEntity).Updates(map[string]interface{}{
			"exchange":        wallet.Exchange,
			"type":            string(wallet.Type),
			"status":          string(wallet.Status),
			"total_usd_value": wallet.TotalUSDValue,
			"metadata":        metadataJSON,
			"last_updated":    wallet.LastUpdated,
			"last_sync_at":    wallet.LastSyncAt,
		}).Error; err != nil {
			r.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to update wallet")
			tx.Rollback()
			return err
		}
	}

	// Delete existing balances
	if err := tx.Where("wallet_id = ?", wallet.ID).Delete(&EnhancedWalletBalanceEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to delete existing balances")
		tx.Rollback()
		return err
	}

	// Create new balances
	for asset, balance := range wallet.Balances {
		balanceEntity := EnhancedWalletBalanceEntity{
			WalletID: wallet.ID,
			Asset:    string(asset),
			Free:     balance.Free,
			Locked:   balance.Locked,
			Total:    balance.Total,
			USDValue: balance.USDValue,
		}

		if err := tx.Create(&balanceEntity).Error; err != nil {
			r.logger.Error().Err(err).Str("id", wallet.ID).Str("asset", string(asset)).Msg("Failed to create balance")
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

// GetByID retrieves a wallet by ID
func (r *ConsolidatedWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	r.logger.Debug().Str("id", id).Msg("Getting wallet by ID")

	// Get wallet entity
	var walletEntity EnhancedWalletEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&walletEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found
		}
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to get wallet")
		return nil, err
	}

	// Get balances
	var balanceEntities []EnhancedWalletBalanceEntity
	if err := r.db.WithContext(ctx).Where("wallet_id = ?", id).Find(&balanceEntities).Error; err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to get balances")
		return nil, err
	}

	// Convert to domain model
	return r.toDomain(&walletEntity, balanceEntities), nil
}

// GetByUserID retrieves a wallet by user ID
func (r *ConsolidatedWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	r.logger.Debug().Str("userID", userID).Msg("Getting wallet by user ID")

	// Get wallet entity
	var walletEntity EnhancedWalletEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&walletEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found
		}
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallet")
		return nil, err
	}

	// Get balances
	var balanceEntities []EnhancedWalletBalanceEntity
	if err := r.db.WithContext(ctx).Where("wallet_id = ?", walletEntity.ID).Find(&balanceEntities).Error; err != nil {
		r.logger.Error().Err(err).Str("id", walletEntity.ID).Msg("Failed to get balances")
		return nil, err
	}

	// Convert to domain model
	return r.toDomain(&walletEntity, balanceEntities), nil
}

// GetWalletsByUserID retrieves all wallets for a user
func (r *ConsolidatedWalletRepository) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	r.logger.Debug().Str("userID", userID).Msg("Getting wallets by user ID")

	// Get wallet entities
	var walletEntities []EnhancedWalletEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&walletEntities).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallets")
		return nil, err
	}

	// Convert to domain models
	wallets := make([]*model.Wallet, len(walletEntities))
	for i, walletEntity := range walletEntities {
		// Get balances for this wallet
		var balanceEntities []EnhancedWalletBalanceEntity
		if err := r.db.WithContext(ctx).Where("wallet_id = ?", walletEntity.ID).Find(&balanceEntities).Error; err != nil {
			r.logger.Error().Err(err).Str("id", walletEntity.ID).Msg("Failed to get balances")
			return nil, err
		}

		// Convert to domain model
		wallets[i] = r.toDomain(&walletEntity, balanceEntities)
	}

	return wallets, nil
}

// DeleteWallet deletes a wallet
func (r *ConsolidatedWalletRepository) DeleteWallet(ctx context.Context, id string) error {
	r.logger.Debug().Str("id", id).Msg("Deleting wallet")

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

	// Delete balances
	if err := tx.Where("wallet_id = ?", id).Delete(&EnhancedWalletBalanceEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to delete balances")
		tx.Rollback()
		return err
	}

	// Delete wallet
	if err := tx.Where("id = ?", id).Delete(&EnhancedWalletEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to delete wallet")
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// SaveBalanceHistory saves a balance history record
func (r *ConsolidatedWalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	r.logger.Debug().
		Str("userID", history.UserID).
		Str("walletID", history.WalletID).
		Time("timestamp", history.Timestamp).
		Msg("Saving balance history")

	// Create history entity with balances as JSON
	balancesJSON, err := json.Marshal(history.Balances)
	if err != nil {
		r.logger.Error().Err(err).Str("userID", history.UserID).Msg("Failed to marshal balances")
		return err
	}

	// Create history entity
	historyEntity := EnhancedWalletBalanceHistoryEntity{
		ID:            uuid.New().String(),
		UserID:        history.UserID,
		WalletID:      history.WalletID,
		BalancesJSON:  balancesJSON,
		TotalUSDValue: history.TotalUSDValue,
		Timestamp:     history.Timestamp,
	}

	// Save history entity
	if err := r.db.WithContext(ctx).Create(&historyEntity).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", history.UserID).Msg("Failed to save balance history")
		return err
	}

	return nil
}

// GetBalanceHistory retrieves balance history for a user and asset within a time range
func (r *ConsolidatedWalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	r.logger.Debug().
		Str("userID", userID).
		Str("asset", string(asset)).
		Time("from", from).
		Time("to", to).
		Msg("Getting balance history")

	// Build query
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Add time range filters
	if !from.IsZero() {
		query = query.Where("timestamp >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("timestamp <= ?", to)
	}

	// Execute query
	var historyEntities []EnhancedWalletBalanceHistoryEntity
	if err := query.Order("timestamp ASC").Find(&historyEntities).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get balance history")
		return nil, err
	}

	// Convert to domain models
	history := make([]*model.BalanceHistory, len(historyEntities))
	for i, entity := range historyEntities {
		// Unmarshal balances JSON
		var balancesMap map[string]map[string]interface{}
		if err := json.Unmarshal(entity.BalancesJSON, &balancesMap); err != nil {
			r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to unmarshal balances JSON")
			balancesMap = make(map[string]map[string]interface{})
		}

		// Convert to map[model.Asset]*model.Balance
		balances := make(map[model.Asset]*model.Balance)
		for assetStr, balanceData := range balancesMap {
			asset := model.Asset(assetStr)
			balance := &model.Balance{
				Asset: asset,
			}

			// Extract values from the map
			if free, ok := balanceData["Free"].(float64); ok {
				balance.Free = free
			}
			if locked, ok := balanceData["Locked"].(float64); ok {
				balance.Locked = locked
			}
			if total, ok := balanceData["Total"].(float64); ok {
				balance.Total = total
			}
			if usdValue, ok := balanceData["USDValue"].(float64); ok {
				balance.USDValue = usdValue
			}

			balances[asset] = balance
		}

		history[i] = &model.BalanceHistory{
			ID:            entity.ID,
			UserID:        entity.UserID,
			WalletID:      entity.WalletID,
			Balances:      balances,
			TotalUSDValue: entity.TotalUSDValue,
			Timestamp:     entity.Timestamp,
		}
	}

	return history, nil
}

// toDomain converts database entities to a domain wallet
func (r *ConsolidatedWalletRepository) toDomain(walletEntity *EnhancedWalletEntity, balanceEntities []EnhancedWalletBalanceEntity) *model.Wallet {
	// Create wallet
	wallet := &model.Wallet{
		ID:            walletEntity.ID,
		UserID:        walletEntity.UserID,
		Exchange:      walletEntity.Exchange,
		Type:          model.WalletType(walletEntity.Type),
		Status:        model.WalletStatus(walletEntity.Status),
		TotalUSDValue: walletEntity.TotalUSDValue,
		LastUpdated:   walletEntity.LastUpdated,
		LastSyncAt:    walletEntity.LastSyncAt,
		CreatedAt:     walletEntity.CreatedAt,
		UpdatedAt:     walletEntity.UpdatedAt,
		Balances:      make(map[model.Asset]*model.Balance),
	}

	// Parse metadata
	if len(walletEntity.Metadata) > 0 {
		var metadata model.WalletMetadata
		if err := json.Unmarshal(walletEntity.Metadata, &metadata); err != nil {
			r.logger.Error().Err(err).Str("id", walletEntity.ID).Msg("Failed to unmarshal wallet metadata")
			wallet.Metadata = &model.WalletMetadata{}
		} else {
			wallet.Metadata = &metadata
		}
	} else {
		wallet.Metadata = &model.WalletMetadata{}
	}

	// Add balances
	for _, balanceEntity := range balanceEntities {
		asset := model.Asset(balanceEntity.Asset)
		wallet.Balances[asset] = &model.Balance{
			Asset:    asset,
			Free:     balanceEntity.Free,
			Locked:   balanceEntity.Locked,
			Total:    balanceEntity.Total,
			USDValue: balanceEntity.USDValue,
		}
	}

	return wallet
}

// Ensure ConsolidatedWalletRepository implements port.WalletRepository
var _ port.WalletRepository = (*ConsolidatedWalletRepository)(nil)
