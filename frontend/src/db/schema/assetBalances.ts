import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';
import { wallets } from './wallets';

export const assetBalances = sqliteTable('asset_balances', {
  id: text('id').primaryKey(),
  walletId: text('wallet_id').notNull().references(() => wallets.id),
  asset: text('asset').notNull(),
  free: real('free').notNull().default(0),
  locked: real('locked').notNull().default(0),
  total: real('total').notNull().default(0),
  price: real('price').notNull().default(0), // Current price in USDT
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
