CREATE TABLE new_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL UNIQUE,
    found_at TIMESTAMP NOT NULL,
    base_volume REAL NOT NULL,
    quote_volume REAL NOT NULL,
    is_processed BOOLEAN NOT NULL DEFAULT 0,
    is_deleted BOOLEAN NOT NULL DEFAULT 0
);

CREATE INDEX idx_new_coins_symbol ON new_coins(symbol);
CREATE INDEX idx_new_coins_processed ON new_coins(is_processed);