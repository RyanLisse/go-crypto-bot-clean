-- Create balance entities table
CREATE TABLE balance_entities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wallet_id INTEGER NOT NULL,
    asset VARCHAR(20) NOT NULL,
    free DECIMAL(18,8) NOT NULL,
    locked DECIMAL(18,8) NOT NULL,
    total DECIMAL(18,8) NOT NULL,
    usd_value DECIMAL(18,8) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES wallet_entities(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_balance_entities_wallet_id ON balance_entities(wallet_id);
CREATE INDEX idx_balance_entities_asset ON balance_entities(asset);
CREATE UNIQUE INDEX idx_balance_entities_wallet_id_asset ON balance_entities(wallet_id, asset);
