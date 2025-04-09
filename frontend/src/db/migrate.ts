import { createClient } from '@libsql/client';
import { drizzle } from 'drizzle-orm/libsql';
import { migrate } from 'drizzle-orm/libsql/migrator';
import * as dotenv from 'dotenv';

dotenv.config();

// Get Turso credentials from environment variables
const url = process.env.VITE_TURSO_URL;
const authToken = process.env.VITE_TURSO_AUTH_TOKEN;

if (!url || !authToken) {
  throw new Error('Turso credentials not found in environment variables');
}

// Create a Turso client
const tursoClient = createClient({
  url,
  authToken,
});

// Create a Drizzle ORM instance
const db = drizzle(tursoClient);

// Run migrations
async function runMigrations() {
  console.log('Running migrations...');
  
  try {
    await migrate(db, { migrationsFolder: './src/db/migrations' });
    console.log('Migrations completed successfully');
  } catch (error) {
    console.error('Error running migrations:', error);
    process.exit(1);
  }
  
  process.exit(0);
}

runMigrations();
