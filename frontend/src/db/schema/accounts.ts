import { sqliteTable, text, integer } from 'drizzle-orm/sqlite-core';
import { users } from './users';

export const accounts = sqliteTable('accounts', {
  id: text('id').primaryKey(),
  userId: text('user_id').notNull().references(() => users.id),
  name: text('name').notNull(),
  type: text('type').notNull(), // spot, margin, futures
  status: text('status').notNull().default('active'),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
