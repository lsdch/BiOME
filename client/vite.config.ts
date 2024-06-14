import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import mkcert from "vite-plugin-mkcert"
import { watchAndRun } from 'vite-plugin-watch-and-run'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  optimizeDeps: {
    exclude: ['leaflet.fullscreen'],
  },
  server: {
    proxy: {
      // redirects API calls made on the client host towards the API server
      '^/api': {
        target: "http://localhost:8080",
        changeOrigin: true
      },
      "^/assets": {
        target: "http://localhost:8080",
        changeOrigin: true
      }
    },
    watch: {
      awaitWriteFinish: true,
      ignored: "openapi.json"
    }
  },
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => ["elements-api"].includes(tag)
        }
      }
    }),
    vueJsx(),
    mkcert(),
    // Generate API client when OpenAPI spec changes
    watchAndRun([
      {
        name: 'gen',
        watchKind: ['add', 'change'],
        watch: path.resolve('openapi.json'),
        run: 'echo "✨ Generating API client" && pnpm run gen-client',
        delay: 0,
        logs: ['streamData', 'streamError']
      }
    ])
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
