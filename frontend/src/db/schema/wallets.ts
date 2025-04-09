import { sqliteTable, text, integer } from 'drizzle-orm/sqlite-core';
import { accounts } from './accounts';

export const wallets = sqliteTable('wallets', {
  id: text('id').primaryKey(),
  accountId: text('account_id').notNull().references(() => accounts.id),
  type: text('type').notNull(), // spot, margin, futures
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
