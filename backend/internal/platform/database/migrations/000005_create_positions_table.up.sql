CREATE TABLE IF NOT EXISTS positions (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    quantity REAL NOT NULL,
    entry_price REAL NOT NULL,
    current_price REAL NOT NULL,
    open_time TIMESTAMP NOT NULL,
    stop_loss REAL,
    take_profit REAL,
    trailing_stop REAL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    pnl REAL DEFAULT 0,
    pnl_percentage REAL DEFAULT 0,
    status TEXT NOT NULL,
    orders_json TEXT
);

CREATE INDEX IF NOT EXISTS idx_positions_symbol ON positions(symbol);
CREATE INDEX IF NOT EXISTS idx_positions_status ON positions(status);
CREATE INDEX IF NOT EXISTS idx_positions_open_time ON positions(open_time);
