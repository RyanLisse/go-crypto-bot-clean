import { sqliteTable, text, integer, blob, real } from 'drizzle-orm/sqlite-core';
import { sql } from 'drizzle-orm';
import { users } from './users';
import { strategies } from './strategies';

export const backtests = sqliteTable('backtests', {
  id: text('id').primaryKey(),
  userId: text('user_id').notNull().references(() => users.id),
  strategyId: text('strategy_id').notNull().references(() => strategies.id),
  name: text('name').notNull(),
  description: text('description'),
  startDate: integer('start_date', { mode: 'timestamp' }).notNull(),
  endDate: integer('end_date', { mode: 'timestamp' }).notNull(),
  initialBalance: real('initial_balance').notNull(),
  finalBalance: real('final_balance').notNull(),
  totalTrades: integer('total_trades').notNull(),
  winningTrades: integer('winning_trades').notNull(),
  losingTrades: integer('losing_trades').notNull(),
  winRate: real('win_rate').notNull(),
  profitFactor: real('profit_factor').notNull(),
  sharpeRatio: real('sharpe_ratio').notNull(),
  maxDrawdown: real('max_drawdown').notNull(),
  parameters: blob('parameters', { mode: 'json' }).notNull().default(sql`'{}'`),
  status: text('status').notNull(),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
});

export const backtestTrades = sqliteTable('backtest_trades', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  backtestId: text('backtest_id').notNull().references(() => backtests.id),
  symbol: text('symbol').notNull(),
  entryTime: integer('entry_time', { mode: 'timestamp' }).notNull(),
  entryPrice: real('entry_price').notNull(),
  exitTime: integer('exit_time', { mode: 'timestamp' }),
  exitPrice: real('exit_price'),
  quantity: real('quantity').notNull(),
  direction: text('direction').notNull(),
  profitLoss: real('profit_loss'),
  profitLossPct: real('profit_loss_pct'),
  exitReason: text('exit_reason'),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
});

export const backtestEquity = sqliteTable('backtest_equity', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  backtestId: text('backtest_id').notNull().references(() => backtests.id),
  timestamp: integer('timestamp', { mode: 'timestamp' }).notNull(),
  equity: real('equity').notNull(),
  balance: real('balance').notNull(),
  drawdown: real('drawdown').notNull(),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
});
