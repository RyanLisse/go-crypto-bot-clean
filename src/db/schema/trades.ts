import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';
import { portfolios } from './portfolios';

export const trades = sqliteTable('trades', {
  id: text('id').primaryKey(),
  portfolioId: text('portfolio_id').notNull().references(() => portfolios.id),
  symbol: text('symbol').notNull(),
  type: text('type', { enum: ['buy', 'sell'] }).notNull(),
  quantity: real('quantity').notNull(),
  price: real('price').notNull(),
  totalAmount: real('total_amount').notNull(),
  status: text('status', { enum: ['pending', 'completed', 'failed'] }).notNull().default('pending'),
  executedAt: integer('executed_at', { mode: 'timestamp' }),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
