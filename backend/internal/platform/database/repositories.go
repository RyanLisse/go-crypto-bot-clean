package database

import (
	"github.com/jmoiron/sqlx"
	"go-crypto-bot-clean/backend/internal/platform/database/repositories"
	"go-crypto-bot-clean/backend/internal/domain/repository"
)

// NewSQLiteBoughtCoinRepository creates a new SQLite implementation of BoughtCoinRepository
func NewSQLiteBoughtCoinRepository(db *sqlx.DB) repository.BoughtCoinRepository {
	return repositories.NewSQLiteBoughtCoinRepository(db)
}

// NewSQLiteNewCoinRepository creates a new SQLite implementation of NewCoinRepository
func NewSQLiteNewCoinRepository(db *sqlx.DB) repository.NewCoinRepository {
	return repositories.NewSQLiteNewCoinRepository(db)
}
