import { sqliteTable, text, integer, blob } from 'drizzle-orm/sqlite-core';
import { sql } from 'drizzle-orm';
import { users } from './users';

export const strategies = sqliteTable('strategies', {
  id: text('id').primaryKey(),
  userId: text('user_id').notNull().references(() => users.id),
  name: text('name').notNull(),
  description: text('description'),
  parameters: blob('parameters', { mode: 'json' }).notNull().default(sql`'{}'`),
  isEnabled: integer('is_enabled', { mode: 'boolean' }).notNull().default(false),
  status: text('status').notNull().default('active'),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
});

export const strategyPerformance = sqliteTable('strategy_performance', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  strategyId: text('strategy_id').notNull().references(() => strategies.id),
  winRate: blob('win_rate', { mode: 'number' }).notNull(),
  profitFactor: blob('profit_factor', { mode: 'number' }).notNull(),
  sharpeRatio: blob('sharpe_ratio', { mode: 'number' }).notNull(),
  maxDrawdown: blob('max_drawdown', { mode: 'number' }).notNull(),
  totalTrades: integer('total_trades').notNull(),
  periodStart: integer('period_start', { mode: 'timestamp' }).notNull(),
  periodEnd: integer('period_end', { mode: 'timestamp' }).notNull(),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
});

export const strategyParameters = sqliteTable('strategy_parameters', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  strategyId: text('strategy_id').notNull().references(() => strategies.id),
  name: text('name').notNull(),
  type: text('type').notNull(),
  description: text('description'),
  defaultValue: text('default_value'),
  min: text('min'),
  max: text('max'),
  options: blob('options', { mode: 'json' }),
  required: integer('required', { mode: 'boolean' }).notNull().default(false),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(sql`CURRENT_TIMESTAMP`),
});
