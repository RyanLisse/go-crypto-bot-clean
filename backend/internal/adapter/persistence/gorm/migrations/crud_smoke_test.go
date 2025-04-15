package migrations

import (
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCRUDSmoke_NewTables(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	logger := zerolog.Nop()
	migrator := NewMigrator(db, &logger)
	RegisterMigrations(migrator, &logger)
	if err := migrator.RunMigrations(); err != nil {
		t.Fatalf("migrations failed: %v", err)
	}

	// AccountEntity CRUD
	account := &entity.AccountEntity{ID: "acc1", UserID: "user1", Email: "user1@example.com"}
	if err := db.Create(account).Error; err != nil {
		t.Errorf("failed to create AccountEntity: %v", err)
	}
	var accountRead entity.AccountEntity
	if err := db.First(&accountRead, "id = ?", account.ID).Error; err != nil {
		t.Errorf("failed to read AccountEntity: %v", err)
	}
	if err := db.Model(&accountRead).Update("ID", "acc2").Error; err != nil {
		t.Errorf("failed to update AccountEntity: %v", err)
	}
	if err := db.Delete(&entity.AccountEntity{}, "id = ?", "acc2").Error; err != nil {
		t.Errorf("failed to delete AccountEntity: %v", err)
	}

	// WalletEntity CRUD
	wallet := &entity.WalletEntity{ID: "wal1", AccountID: "acc1", Exchange: "binance", TotalUSD: 100.0}
	if err := db.Create(wallet).Error; err != nil {
		t.Errorf("failed to create WalletEntity: %v", err)
	}
	var walletRead entity.WalletEntity
	if err := db.First(&walletRead, "id = ?", wallet.ID).Error; err != nil {
		t.Errorf("failed to read WalletEntity: %v", err)
	}
	if err := db.Model(&walletRead).Update("ID", "wal2").Error; err != nil {
		t.Errorf("failed to update WalletEntity: %v", err)
	}
	if err := db.Delete(&entity.WalletEntity{}, "id = ?", "wal2").Error; err != nil {
		t.Errorf("failed to delete WalletEntity: %v", err)
	}

	// OrderEntity CRUD
	order := &entity.OrderEntity{ID: "ord1", AccountID: "acc1", Symbol: "BTCUSDT", Side: "BUY", Type: "LIMIT", Quantity: 1.0, Price: 50000.0, Status: "NEW"}
	if err := db.Create(order).Error; err != nil {
		t.Errorf("failed to create OrderEntity: %v", err)
	}
	var orderRead entity.OrderEntity
	if err := db.First(&orderRead, "id = ?", order.ID).Error; err != nil {
		t.Errorf("failed to read OrderEntity: %v", err)
	}
	if err := db.Model(&orderRead).Update("ID", "ord2").Error; err != nil {
		t.Errorf("failed to update OrderEntity: %v", err)
	}
	if err := db.Delete(&entity.OrderEntity{}, "id = ?", "ord2").Error; err != nil {
		t.Errorf("failed to delete OrderEntity: %v", err)
	}

	// PositionEntity CRUD
	position := &entity.PositionEntity{ID: "pos1", AccountID: "acc1", Symbol: "BTCUSDT", Side: "LONG", Quantity: 1.0, EntryPrice: 50000.0, Status: "OPEN", OpenedAt: time.Now()}
	if err := db.Create(position).Error; err != nil {
		t.Errorf("failed to create PositionEntity: %v", err)
	}
	var positionRead entity.PositionEntity
	if err := db.First(&positionRead, "id = ?", position.ID).Error; err != nil {
		t.Errorf("failed to read PositionEntity: %v", err)
	}
	if err := db.Model(&positionRead).Update("ID", "pos2").Error; err != nil {
		t.Errorf("failed to update PositionEntity: %v", err)
	}
	if err := db.Delete(&entity.PositionEntity{}, "id = ?", "pos2").Error; err != nil {
		t.Errorf("failed to delete PositionEntity: %v", err)
	}

	// TransactionEntity CRUD
	transaction := &entity.TransactionEntity{ID: "txn1", AccountID: "acc1", Type: "DEPOSIT", Asset: "BTC", Amount: 0.1, Status: "COMPLETED", Timestamp: time.Now()}
	if err := db.Create(transaction).Error; err != nil {
		t.Errorf("failed to create TransactionEntity: %v", err)
	}
	var transactionRead entity.TransactionEntity
	if err := db.First(&transactionRead, "id = ?", transaction.ID).Error; err != nil {
		t.Errorf("failed to read TransactionEntity: %v", err)
	}
	if err := db.Model(&transactionRead).Update("ID", "txn2").Error; err != nil {
		t.Errorf("failed to update TransactionEntity: %v", err)
	}
	if err := db.Delete(&entity.TransactionEntity{}, "id = ?", "txn2").Error; err != nil {
		t.Errorf("failed to delete TransactionEntity: %v", err)
	}
}
