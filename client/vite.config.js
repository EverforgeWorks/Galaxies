import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // Expose for Docker dev modes if needed
    port: 5173,
    proxy: {
      // Proxy API and WebSocket requests to Go backend during local dev
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ws': {
        target: 'http://localhost:8080',
        ws: true, // Vital for WebSocket proxying
        changeOrigin: true,
      }
    }
  }
})
