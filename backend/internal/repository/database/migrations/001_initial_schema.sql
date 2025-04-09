-- Initial schema for TursoDB migration
-- This ensures our schema is compatible with both SQLite and TursoDB

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Enable WAL mode for better concurrency
PRAGMA journal_mode = WAL;

-- Balance history table
CREATE TABLE IF NOT EXISTS balance_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME NOT NULL,
    balance REAL NOT NULL,
    equity REAL NOT NULL,
    free_balance REAL NOT NULL,
    free_margin REAL NOT NULL DEFAULT 0,
    locked_balance REAL NOT NULL,
    unrealized_pnl REAL NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on timestamp for faster queries
CREATE INDEX IF NOT EXISTS idx_balance_history_timestamp ON balance_history(timestamp);

-- Migration metadata table
CREATE TABLE IF NOT EXISTS migration_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version TEXT NOT NULL,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT NOT NULL
);

-- Insert initial migration record
INSERT INTO migration_history (version, description)
VALUES ('001', 'Initial schema for TursoDB migration');
