import path from 'node:path'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    // Increase warning limit for main chunk
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      output: {
        // Split chunks for better caching and load performance
        manualChunks: {
          // Vue ecosystem
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          // TanStack libraries
          'tanstack': ['@tanstack/vue-query', '@tanstack/vue-table'],
          // UI components
          'ui-vendor': ['reka-ui', 'lucide-vue-next'],
          // Utilities
          'utils': ['@vueuse/core', 'class-variance-authority', 'clsx', 'tailwind-merge'],
        },
      },
    },
  },
})
