import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';
import { accounts } from './accounts';

export const trades = sqliteTable('trades', {
  id: text('id').primaryKey(),
  accountId: text('account_id').notNull().references(() => accounts.id),
  symbol: text('symbol').notNull(),
  side: text('side', { enum: ['buy', 'sell'] }).notNull(),
  type: text('type', { enum: ['market', 'limit', 'stop_loss', 'take_profit'] }).notNull(),
  status: text('status', { enum: ['new', 'partially_filled', 'filled', 'canceled', 'rejected', 'expired'] }).notNull(),
  quantity: real('quantity').notNull(),
  price: real('price').notNull(),
  stopPrice: real('stop_price'),
  commission: real('commission'),
  commissionAsset: text('commission_asset'),
  executedAt: integer('executed_at', { mode: 'timestamp' }),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
