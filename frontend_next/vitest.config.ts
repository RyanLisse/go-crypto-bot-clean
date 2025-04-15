import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test-setup.ts'], // Use our more comprehensive setup file
    include: ['src/**/*.test.{ts,tsx,js,jsx}'],
  },
}); 