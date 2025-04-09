-- Rename exit_reason back to reason in closed_positions table
ALTER TABLE closed_positions RENAME COLUMN exit_reason TO reason;

-- Remove new fields from closed_positions table
ALTER TABLE closed_positions DROP COLUMN holding_time_ms;
ALTER TABLE closed_positions DROP COLUMN entry_reason;
ALTER TABLE closed_positions DROP COLUMN strategy;
ALTER TABLE closed_positions DROP COLUMN initial_stop_loss;
ALTER TABLE closed_positions DROP COLUMN initial_take_profit;
ALTER TABLE closed_positions DROP COLUMN risk_reward_ratio;
ALTER TABLE closed_positions DROP COLUMN actual_rr;
ALTER TABLE closed_positions DROP COLUMN expected_value;
ALTER TABLE closed_positions DROP COLUMN max_price;
ALTER TABLE closed_positions DROP COLUMN min_price;
ALTER TABLE closed_positions DROP COLUMN max_drawdown;
ALTER TABLE closed_positions DROP COLUMN max_drawdown_percent;
ALTER TABLE closed_positions DROP COLUMN max_profit;
ALTER TABLE closed_positions DROP COLUMN max_profit_percent;
ALTER TABLE closed_positions DROP COLUMN tags;
ALTER TABLE closed_positions DROP COLUMN notes;
ALTER TABLE closed_positions DROP COLUMN orders_json;

-- Remove new fields from positions table
ALTER TABLE positions DROP COLUMN close_time;
ALTER TABLE positions DROP COLUMN entry_reason;
ALTER TABLE positions DROP COLUMN exit_reason;
ALTER TABLE positions DROP COLUMN strategy;
ALTER TABLE positions DROP COLUMN risk_reward_ratio;
ALTER TABLE positions DROP COLUMN expected_profit;
ALTER TABLE positions DROP COLUMN max_risk;
ALTER TABLE positions DROP COLUMN take_profit_levels_json;
ALTER TABLE positions DROP COLUMN tags;
ALTER TABLE positions DROP COLUMN notes;
