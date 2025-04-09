-- Create wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id INTEGER PRIMARY KEY,
    updated_at TIMESTAMP NOT NULL
);

-- Create asset_balances table
CREATE TABLE IF NOT EXISTS asset_balances (
    asset TEXT NOT NULL,
    free REAL NOT NULL DEFAULT 0,
    locked REAL NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (asset)
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    balance REAL NOT NULL,
    reason TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

-- Create index on timestamp for faster queries
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions (timestamp);
