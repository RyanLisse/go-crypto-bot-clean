import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'], // Point to our updated setup file
    include: ['src/**/*.test.{ts,tsx,js,jsx}'],
  },
}); 