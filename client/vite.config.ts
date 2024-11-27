import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'
import vueJsx from '@vitejs/plugin-vue-jsx'
import mkcert from "vite-plugin-mkcert"
import { watchAndRun } from 'vite-plugin-watch-and-run'
import path from 'path'
import { exec } from 'node:child_process'

// https://vitejs.dev/config/
export default defineConfig({
  css: {
    preprocessorOptions: {
      scss: {
        api: 'modern-compiler'
      }
    }
  },
  optimizeDeps: {
    exclude: ['leaflet.fullscreen'],
    esbuildOptions: {
      define: {
        global: 'globalThis'
      }
    }
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
      // ignored: path.resolve("openapi.json")
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
    vuetify(),
    vueJsx(),
    mkcert(),
    {
      name: "openAPI-gen",
      handleHotUpdate({ file, server }) {
        if (file.endsWith('openapi.json')) {
          console.log("📃 OpenAPI spec written")
          exec('pnpm run --silent gen-client', (error, stdout, stderr) => {
            console.log(stdout)
            console.error(stderr)
          })
          server.ws.send({
            type: 'custom',
            event: 'openAPI-change',
            data: {}
          })
          return []
        }
      },
    },
    // Generate API client when OpenAPI spec changes
    // watchAndRun([
    //   {
    //     name: 'gen',
    //     watchKind: ['add', 'change'],
    //     watch: path.resolve('openapi.json'),
    //     run: 'echo "✨ Generating API client" && pnpm run gen-client',
    //     delay: 0,
    //     logs: ['streamData', 'streamError'],
    //   }
    // ])
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      stream: 'stream-browserify'
    }
  }
})
