-- Drop indexes
DROP INDEX IF EXISTS idx_trade_performance_symbol;
DROP INDEX IF EXISTS idx_trade_performance_entry_time;
DROP INDEX IF EXISTS idx_trade_performance_exit_time;
DROP INDEX IF EXISTS idx_trade_performance_profit_loss;
DROP INDEX IF EXISTS idx_trade_performance_position_id;
DROP INDEX IF EXISTS idx_trade_performance_strategy;
DROP INDEX IF EXISTS idx_balance_history_timestamp;

-- Drop tables
DROP TABLE IF EXISTS trade_performance;
DROP TABLE IF EXISTS balance_history;
