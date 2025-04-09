CREATE TABLE bought_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL UNIQUE,
    purchase_price REAL NOT NULL,
    quantity REAL NOT NULL,
    bought_at TIMESTAMP NOT NULL,
    stop_loss REAL NOT NULL,
    take_profit REAL NOT NULL,
    current_price REAL NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bought_coins_symbol ON bought_coins(symbol);