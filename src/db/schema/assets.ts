import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';
import { portfolios } from './portfolios';

export const assets = sqliteTable('assets', {
  id: text('id').primaryKey(),
  portfolioId: text('portfolio_id').notNull().references(() => portfolios.id),
  symbol: text('symbol').notNull(),
  name: text('name').notNull(),
  quantity: real('quantity').notNull().default(0),
  averageBuyPrice: real('average_buy_price').notNull().default(0),
  currentPrice: real('current_price').default(0),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
