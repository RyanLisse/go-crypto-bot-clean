-- Drop indexes
DROP INDEX IF EXISTS idx_trade_decisions_symbol;
DROP INDEX IF EXISTS idx_trade_decisions_type;
DROP INDEX IF EXISTS idx_trade_decisions_status;
DROP INDEX IF EXISTS idx_trade_decisions_reason;
DROP INDEX IF EXISTS idx_trade_decisions_created_at;
DROP INDEX IF EXISTS idx_trade_decisions_executed_at;
DROP INDEX IF EXISTS idx_trade_decisions_position_id;
DROP INDEX IF EXISTS idx_trade_decisions_order_id;

-- Drop table
DROP TABLE IF EXISTS trade_decisions;
