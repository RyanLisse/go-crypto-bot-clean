package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/database/repositories"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
)

// NewSQLiteBoughtCoinRepository creates a new SQLite implementation of BoughtCoinRepository
func NewSQLiteBoughtCoinRepository(db *sqlx.DB) repository.BoughtCoinRepository {
	return repositories.NewSQLiteBoughtCoinRepository(db)
}

// NewSQLiteNewCoinRepository creates a new SQLite implementation of NewCoinRepository
func NewSQLiteNewCoinRepository(db *sqlx.DB) repository.NewCoinRepository {
	return repositories.NewSQLiteNewCoinRepository(db)
}
