-- Create trade_performance table to track individual trade performance
CREATE TABLE IF NOT EXISTS trade_performance (
    id TEXT PRIMARY KEY,
    trade_id TEXT NOT NULL,
    position_id TEXT,
    symbol TEXT NOT NULL,
    entry_time TIMESTAMP NOT NULL,
    exit_time TIMESTAMP NOT NULL,
    entry_price REAL NOT NULL,
    exit_price REAL NOT NULL,
    quantity REAL NOT NULL,
    profit_loss REAL NOT NULL,
    profit_loss_percent REAL NOT NULL,
    holding_time_ms INTEGER NOT NULL,
    entry_reason TEXT,
    exit_reason TEXT,
    strategy TEXT,
    stop_loss REAL,
    take_profit REAL,
    risk_reward_ratio REAL,
    expected_value REAL,
    actual_rr REAL,
    tags TEXT,
    metadata_json TEXT,
    
    FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL
);

-- Create balance_history table to track account balance over time
CREATE TABLE IF NOT EXISTS balance_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp TIMESTAMP NOT NULL,
    balance REAL NOT NULL,
    equity REAL NOT NULL,
    free_balance REAL NOT NULL,
    locked_balance REAL NOT NULL,
    unrealized_pnl REAL NOT NULL
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_trade_performance_symbol ON trade_performance(symbol);
CREATE INDEX IF NOT EXISTS idx_trade_performance_entry_time ON trade_performance(entry_time);
CREATE INDEX IF NOT EXISTS idx_trade_performance_exit_time ON trade_performance(exit_time);
CREATE INDEX IF NOT EXISTS idx_trade_performance_profit_loss ON trade_performance(profit_loss);
CREATE INDEX IF NOT EXISTS idx_trade_performance_position_id ON trade_performance(position_id);
CREATE INDEX IF NOT EXISTS idx_trade_performance_strategy ON trade_performance(strategy);
CREATE INDEX IF NOT EXISTS idx_balance_history_timestamp ON balance_history(timestamp);
