CREATE TABLE purchase_decisions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    coin_symbol TEXT NOT NULL,
    decision_time TIMESTAMP NOT NULL,
    decision TEXT NOT NULL,
    reason TEXT NOT NULL,
    price REAL NOT NULL,
    volume REAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_purchase_decisions_coin_symbol ON purchase_decisions(coin_symbol);
CREATE INDEX idx_purchase_decisions_decision_time ON purchase_decisions(decision_time);