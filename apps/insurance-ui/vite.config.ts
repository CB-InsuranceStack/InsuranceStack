import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: process.env.BASE_PATH || '/',
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/api/accounts': {
        target: 'http://api-accounts:8001',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/accounts/, ''),
      },
      '/api/transactions': {
        target: 'http://api-transactions:8002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/transactions/, ''),
      },
      '/api/insights': {
        target: 'http://api-insights:8003',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/insights/, ''),
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
    },
  },
});
