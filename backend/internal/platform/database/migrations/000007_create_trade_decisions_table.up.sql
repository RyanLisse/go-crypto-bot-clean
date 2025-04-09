-- Create trade_decisions table
CREATE TABLE IF NOT EXISTS trade_decisions (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    reason TEXT NOT NULL,
    detailed_reason TEXT,
    price REAL NOT NULL,
    quantity REAL NOT NULL,
    total_value REAL NOT NULL,
    confidence REAL DEFAULT 0,
    strategy TEXT,
    strategy_params TEXT,
    created_at TIMESTAMP NOT NULL,
    executed_at TIMESTAMP,
    position_id TEXT,
    order_id TEXT,
    stop_loss REAL,
    take_profit REAL,
    trailing_stop REAL,
    risk_reward_ratio REAL,
    expected_profit REAL,
    max_risk REAL,
    tags TEXT,
    metadata_json TEXT,
    
    FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_trade_decisions_symbol ON trade_decisions(symbol);
CREATE INDEX IF NOT EXISTS idx_trade_decisions_type ON trade_decisions(type);
CREATE INDEX IF NOT EXISTS idx_trade_decisions_status ON trade_decisions(status);
CREATE INDEX IF NOT EXISTS idx_trade_decisions_reason ON trade_decisions(reason);
CREATE INDEX IF NOT EXISTS idx_trade_decisions_created_at ON trade_decisions(created_at);
CREATE INDEX IF NOT EXISTS idx_trade_decisions_executed_at ON trade_decisions(executed_at);
CREATE INDEX IF NOT EXISTS idx_trade_decisions_position_id ON trade_decisions(position_id);
CREATE INDEX IF NOT EXISTS idx_trade_decisions_order_id ON trade_decisions(order_id);
