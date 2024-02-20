import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import mkcert from "vite-plugin-mkcert"

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      // redirects API calls made on the client host towards the API server
      '^/api': {
        target: "http://localhost:8080",
        changeOrigin: true
      }
    }
  },
  plugins: [
    vue(),
    vueJsx(),
    mkcert()
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
