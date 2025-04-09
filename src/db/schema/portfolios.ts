import { sqliteTable, text, integer, real } from 'drizzle-orm/sqlite-core';
import { users } from './users';

export const portfolios = sqliteTable('portfolios', {
  id: text('id').primaryKey(),
  userId: text('user_id').notNull().references(() => users.id),
  name: text('name').notNull(),
  description: text('description'),
  totalValue: real('total_value').notNull().default(0),
  createdAt: integer('created_at', { mode: 'timestamp' }).notNull().default(Date.now()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).notNull().default(Date.now()),
});
