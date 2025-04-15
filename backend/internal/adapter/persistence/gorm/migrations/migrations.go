package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// RegisterMigrations registers all available migrations with the migrator
func RegisterMigrations(migrator *Migrator, logger *zerolog.Logger) {
	// Register migrations in order
	migrations := []struct {
		name     string
		function func(*gorm.DB) error
	}{
		{"001_create_status_table", CreateStatusTable},
		{"002_create_ticker_table", CreateTickerTable},
		{"003_create_symbols_table", CreateSymbolsTable},
		{"004_create_wallet_table", CreateWalletTable},
		{"005_create_position_table", CreatePositionTable},
		{"006_add_symbol_columns", AddSymbolColumns},
		{"007_add_symbol_status", AddSymbolStatus},
		{"008_extend_symbol_metadata", ExtendSymbolMetadata},
		{"009_add_position_columns", AddPositionColumns},
		{"010_add_order_table", CreateOrderTable},
		{"011_add_position_order_id", AddPositionOrderId},
		{"012_add_order_extra_fields", AddOrderExtraFields},
		{"013_add_order_symbol_index", AddOrderSymbolIndex},
		{"014_add_ticker_symbols_index", AddTickerSymbolsIndex},
		{"015_create_auto_buy_rules_table", CreateAutoBuyRulesTable},
		{"016_create_auto_buy_executions_table", CreateAutoBuyExecutionsTable},
		{"017_add_position_order_price", AddPositionOrderPrice},
		{"018_add_symbols_usdt_index", AddSymbolsUsdtIndex},
		{"019_add_symbols_metadata_indexes", AddSymbolsMetadataIndexes},
		{"020_create_mexc_market_data_tables", CreateMexcMarketDataTables},
		{"021_create_account_table", CreateAccountTable},
		{"022_create_wallets_table", CreateWalletsTable},
		{"023_create_orders_table", CreateOrdersTable},
		{"024_create_positions_table", CreatePositionsTable},
		{"025_create_transactions_table", CreateTransactionsTable},
		{"026_create_mexc_api_credentials_table", CreateMexcApiCredentialsTable},
		{"027_create_users_table", CreateUsersTable},
		{"028_create_api_credentials_table", CreateAPICredentialsTable},
		{"029_create_wallet_entities_table", CreateWalletEntitiesTable},
		{"030_create_balance_entities_table", CreateBalanceEntitiesTable},
	}

	// Register object-based migrations
	migrator.AddObjectMigration(NewCreateEnhancedWalletsTable(logger))
	migrator.AddObjectMigration(NewCreateEnhancedWalletBalanceHistoryTable(logger))

	for _, migration := range migrations {
		migrator.AddMigration(migration.name, migration.function)
		logger.Debug().Str("migration", migration.name).Msg("Registered migration")
	}
}
