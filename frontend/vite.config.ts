import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "localhost",
    port: 5173,
    strictPort: false,
    hmr: {
      timeout: 10000
    }
  },
  plugins: [
    react()
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@playwright/test": path.resolve(__dirname, "src/__mocks__/@playwright/test.js"),
      "@google/generative-ai": path.resolve(__dirname, "src/__mocks__/@google/generative-ai.js"),
      "date-fns/differenceInCalendarISOWeekYears": "date-fns/differenceInCalendarYears"
    },
  },
  build: {
    sourcemap: process.env.NODE_ENV !== 'production',
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', 'react-router-dom'],
          ui: [
            '@radix-ui/react-dialog',
            '@radix-ui/react-popover',
            '@radix-ui/react-tabs',
            'lucide-react'
          ]
        }
      }
    }
  },
  optimizeDeps: {
    force: true
  }
});
