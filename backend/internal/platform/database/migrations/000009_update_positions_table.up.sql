-- Add new fields to positions table
ALTER TABLE positions ADD COLUMN close_time TIMESTAMP;
ALTER TABLE positions ADD COLUMN entry_reason TEXT;
ALTER TABLE positions ADD COLUMN exit_reason TEXT;
ALTER TABLE positions ADD COLUMN strategy TEXT;
ALTER TABLE positions ADD COLUMN risk_reward_ratio REAL DEFAULT 0;
ALTER TABLE positions ADD COLUMN expected_profit REAL DEFAULT 0;
ALTER TABLE positions ADD COLUMN max_risk REAL DEFAULT 0;
ALTER TABLE positions ADD COLUMN take_profit_levels_json TEXT;
ALTER TABLE positions ADD COLUMN tags TEXT;
ALTER TABLE positions ADD COLUMN notes TEXT;

-- Add new fields to closed_positions table
ALTER TABLE closed_positions ADD COLUMN holding_time_ms INTEGER DEFAULT 0;
ALTER TABLE closed_positions ADD COLUMN entry_reason TEXT;
ALTER TABLE closed_positions ADD COLUMN strategy TEXT;
ALTER TABLE closed_positions ADD COLUMN initial_stop_loss REAL;
ALTER TABLE closed_positions ADD COLUMN initial_take_profit REAL;
ALTER TABLE closed_positions ADD COLUMN risk_reward_ratio REAL;
ALTER TABLE closed_positions ADD COLUMN actual_rr REAL;
ALTER TABLE closed_positions ADD COLUMN expected_value REAL;
ALTER TABLE closed_positions ADD COLUMN max_price REAL;
ALTER TABLE closed_positions ADD COLUMN min_price REAL;
ALTER TABLE closed_positions ADD COLUMN max_drawdown REAL;
ALTER TABLE closed_positions ADD COLUMN max_drawdown_percent REAL;
ALTER TABLE closed_positions ADD COLUMN max_profit REAL;
ALTER TABLE closed_positions ADD COLUMN max_profit_percent REAL;
ALTER TABLE closed_positions ADD COLUMN tags TEXT;
ALTER TABLE closed_positions ADD COLUMN notes TEXT;
ALTER TABLE closed_positions ADD COLUMN orders_json TEXT;

-- Rename reason to exit_reason in closed_positions table
ALTER TABLE closed_positions RENAME COLUMN reason TO exit_reason;
