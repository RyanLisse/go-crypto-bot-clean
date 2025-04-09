import { createClient } from '@libsql/client';
import { drizzle } from 'drizzle-orm/libsql';
import * as schema from './schema';

// Get Turso credentials from environment variables
const url = import.meta.env.VITE_TURSO_URL;
const authToken = import.meta.env.VITE_TURSO_AUTH_TOKEN;

if (!url || !authToken) {
  throw new Error('Turso credentials not found in environment variables');
}

// Create a Turso client
const tursoClient = createClient({
  url,
  authToken,
});

// Create a Drizzle ORM instance
export const db = drizzle(tursoClient, { schema });

export default db;
