package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// EnhancedWalletRepository implements the WalletRepository interface using GORM
type EnhancedWalletRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewEnhancedWalletRepository creates a new EnhancedWalletRepository
func NewEnhancedWalletRepository(db *gorm.DB, logger *zerolog.Logger) port.WalletRepository {
	return &EnhancedWalletRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists a wallet to the database
func (r *EnhancedWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	r.logger.Debug().
		Str("id", wallet.ID).
		Str("userID", wallet.UserID).
		Str("type", string(wallet.Type)).
		Str("status", string(wallet.Status)).
		Float64("totalUSDValue", wallet.TotalUSDValue).
		Msg("Saving enhanced wallet to database")

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

	// Create wallet metadata entity
	metadata := entity.WalletMetadataEntity{}
	if wallet.Metadata != nil {
		metadata = entity.WalletMetadataEntity{
			Name:        wallet.Metadata.Name,
			Description: wallet.Metadata.Description,
			Tags:        wallet.Metadata.Tags,
			IsPrimary:   wallet.Metadata.IsPrimary,
			Network:     wallet.Metadata.Network,
			Address:     wallet.Metadata.Address,
			Custom:      wallet.Metadata.Custom,
		}
	}

	// Create wallet entity
	walletEntity := entity.EnhancedWalletEntity{
		ID:            wallet.ID,
		UserID:        wallet.UserID,
		Exchange:      wallet.Exchange,
		Type:          string(wallet.Type),
		Status:        string(wallet.Status),
		TotalUSDValue: wallet.TotalUSDValue,
		Metadata:      metadata,
		LastUpdated:   wallet.LastUpdated,
		LastSyncAt:    &wallet.LastSyncAt,
		CreatedAt:     wallet.CreatedAt,
		UpdatedAt:     wallet.UpdatedAt,
	}

	// Save wallet entity
	if err := tx.Save(&walletEntity).Error; err != nil {
		tx.Rollback()
		r.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to save wallet entity")
		return err
	}

	// Delete existing balances for this wallet
	if err := tx.Where("wallet_id = ?", wallet.ID).Delete(&entity.EnhancedWalletBalanceEntity{}).Error; err != nil {
		tx.Rollback()
		r.logger.Error().Err(err).Str("walletID", wallet.ID).Msg("Failed to delete existing balances")
		return err
	}

	// Save balances
	for asset, balance := range wallet.Balances {
		balanceEntity := entity.EnhancedWalletBalanceEntity{
			WalletID: wallet.ID,
			Asset:    string(asset),
			Free:     balance.Free,
			Locked:   balance.Locked,
			Total:    balance.Total,
			USDValue: balance.USDValue,
		}

		if err := tx.Create(&balanceEntity).Error; err != nil {
			tx.Rollback()
			r.logger.Error().Err(err).Str("walletID", wallet.ID).Str("asset", string(asset)).Msg("Failed to save balance entity")
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		r.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to commit transaction")
		return err
	}

	return nil
}

// GetByID retrieves a wallet by its ID
func (r *EnhancedWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	r.logger.Debug().Str("id", id).Msg("Getting wallet by ID")

	// Get wallet entity
	var walletEntity entity.EnhancedWalletEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&walletEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to get wallet entity")
		return nil, err
	}

	// Get balances
	var balanceEntities []entity.EnhancedWalletBalanceEntity
	if err := r.db.WithContext(ctx).Where("wallet_id = ?", id).Find(&balanceEntities).Error; err != nil {
		r.logger.Error().Err(err).Str("walletID", id).Msg("Failed to get balance entities")
		return nil, err
	}

	// Create domain wallet
	wallet := &model.Wallet{
		ID:            walletEntity.ID,
		UserID:        walletEntity.UserID,
		Exchange:      walletEntity.Exchange,
		Type:          model.WalletType(walletEntity.Type),
		Status:        model.WalletStatus(walletEntity.Status),
		Balances:      make(map[model.Asset]*model.Balance),
		TotalUSDValue: walletEntity.TotalUSDValue,
		LastUpdated:   walletEntity.LastUpdated,
		CreatedAt:     walletEntity.CreatedAt,
		UpdatedAt:     walletEntity.UpdatedAt,
	}

	// Set LastSyncAt if not nil
	if walletEntity.LastSyncAt != nil {
		wallet.LastSyncAt = *walletEntity.LastSyncAt
	}

	// Set metadata
	wallet.Metadata = &model.WalletMetadata{
		Name:        walletEntity.Metadata.Name,
		Description: walletEntity.Metadata.Description,
		Tags:        walletEntity.Metadata.Tags,
		IsPrimary:   walletEntity.Metadata.IsPrimary,
		Network:     walletEntity.Metadata.Network,
		Address:     walletEntity.Metadata.Address,
		Custom:      walletEntity.Metadata.Custom,
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

	return wallet, nil
}

// GetByUserID retrieves a wallet by user ID
func (r *EnhancedWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	r.logger.Debug().Str("userID", userID).Msg("Getting wallet by user ID")

	// Get wallet entity
	var walletEntity entity.EnhancedWalletEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&walletEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallet entity")
		return nil, err
	}

	// Get the wallet by ID
	return r.GetByID(ctx, walletEntity.ID)
}

// GetWalletsByUserID retrieves all wallets for a user
func (r *EnhancedWalletRepository) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	r.logger.Debug().Str("userID", userID).Msg("Getting all wallets for user")

	// Get wallet entities
	var walletEntities []entity.EnhancedWalletEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&walletEntities).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallet entities")
		return nil, err
	}

	// If no wallets found, return empty slice
	if len(walletEntities) == 0 {
		return []*model.Wallet{}, nil
	}

	// Get all wallets by ID
	wallets := make([]*model.Wallet, 0, len(walletEntities))
	for _, walletEntity := range walletEntities {
		wallet, err := r.GetByID(ctx, walletEntity.ID)
		if err != nil {
			r.logger.Error().Err(err).Str("id", walletEntity.ID).Msg("Failed to get wallet by ID")
			continue
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// SaveBalanceHistory saves a balance history record
func (r *EnhancedWalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	r.logger.Debug().
		Str("id", history.ID).
		Str("userID", history.UserID).
		Str("walletID", history.WalletID).
		Time("timestamp", history.Timestamp).
		Msg("Saving balance history")

	// Create balance history entity with balances as JSON
	balancesJSON, err := json.Marshal(history.Balances)
	if err != nil {
		r.logger.Error().Err(err).Str("userID", history.UserID).Msg("Failed to marshal balances")
		return err
	}

	// Create balance history entity
	historyEntity := entity.EnhancedWalletBalanceHistoryEntity{
		ID:            history.ID,
		UserID:        history.UserID,
		WalletID:      history.WalletID,
		BalancesJSON:  balancesJSON,
		TotalUSDValue: history.TotalUSDValue,
		Timestamp:     history.Timestamp,
	}

	// Save balance history entity
	if err := r.db.WithContext(ctx).Create(&historyEntity).Error; err != nil {
		r.logger.Error().Err(err).Str("id", history.ID).Msg("Failed to save balance history entity")
		return err
	}

	return nil
}

// GetBalanceHistory retrieves balance history for a user and asset
func (r *EnhancedWalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	r.logger.Debug().
		Str("userID", userID).
		Str("asset", string(asset)).
		Time("from", from).
		Time("to", to).
		Msg("Getting balance history")

	// Get balance history entities
	var historyEntities []entity.EnhancedWalletBalanceHistoryEntity
	query := r.db.WithContext(ctx).Where("user_id = ? AND timestamp BETWEEN ? AND ?", userID, from, to)
	if err := query.Order("timestamp ASC").Find(&historyEntities).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get balance history entities")
		return nil, err
	}

	// Convert to domain model
	history := make([]*model.BalanceHistory, len(historyEntities))
	for i, historyEntity := range historyEntities {
		// Unmarshal balances JSON
		var balancesMap map[string]float64
		if err := json.Unmarshal(historyEntity.BalancesJSON, &balancesMap); err != nil {
			r.logger.Error().Err(err).Str("id", historyEntity.ID).Msg("Failed to unmarshal balances JSON")
			balancesMap = make(map[string]float64)
		}

		// Convert to map[model.Asset]*model.Balance
		balances := make(map[model.Asset]*model.Balance)
		for assetStr, value := range balancesMap {
			assetObj := model.Asset(assetStr)
			balances[assetObj] = &model.Balance{
				Asset: assetObj,
				Free:  value,
				Total: value,
			}
		}

		history[i] = &model.BalanceHistory{
			ID:            historyEntity.ID,
			UserID:        historyEntity.UserID,
			WalletID:      historyEntity.WalletID,
			Balances:      balances,
			TotalUSDValue: historyEntity.TotalUSDValue,
			Timestamp:     historyEntity.Timestamp,
		}
	}

	return history, nil
}

// DeleteWallet deletes a wallet by ID
func (r *EnhancedWalletRepository) DeleteWallet(ctx context.Context, id string) error {
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
	if err := tx.Where("wallet_id = ?", id).Delete(&entity.EnhancedWalletBalanceEntity{}).Error; err != nil {
		tx.Rollback()
		r.logger.Error().Err(err).Str("walletID", id).Msg("Failed to delete balances")
		return err
	}

	// Delete wallet
	if err := tx.Where("id = ?", id).Delete(&entity.EnhancedWalletEntity{}).Error; err != nil {
		tx.Rollback()
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to delete wallet")
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to commit transaction")
		return err
	}

	return nil
}
