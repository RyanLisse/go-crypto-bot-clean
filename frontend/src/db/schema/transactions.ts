import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';
import { accounts } from './accounts';

export const transactions = sqliteTable('transactions', {
  id: text('id').primaryKey(),
  accountId: text('account_id').notNull().references(() => accounts.id),
  type: text('type', { enum: ['deposit', 'withdrawal', 'transfer'] }).notNull(),
  asset: text('asset').notNull(),
  amount: real('amount').notNull(),
  status: text('status', { enum: ['pending', 'completed', 'failed'] }).notNull().default('pending'),
  txId: text('tx_id'),
  address: text('address'),
  network: text('network'),
  fee: real('fee'),
  feeAsset: text('fee_asset'),
  completedAt: integer('completed_at', { mode: 'timestamp' }),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
