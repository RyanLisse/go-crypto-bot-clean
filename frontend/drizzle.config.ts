import type { Config } from 'drizzle-kit';
import * as dotenv from 'dotenv';

dotenv.config();

export default {
  schema: './src/db/schema/*',
  out: './src/db/migrations',
  driver: 'turso',
  dbCredentials: {
    url: process.env.VITE_TURSO_URL as string,
    authToken: process.env.VITE_TURSO_AUTH_TOKEN as string,
  },
  verbose: true,
  strict: true,
} satisfies Config;
