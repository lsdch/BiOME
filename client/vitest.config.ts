import { fileURLToPath } from 'node:url'
import { configDefaults, defineConfig, mergeConfig } from 'vitest/config'
import viteConfig from './vite.config'

export default defineConfig((env) => {
  const resolvedViteConfig = typeof viteConfig === 'function' ? viteConfig(env) : viteConfig

  return mergeConfig(resolvedViteConfig, {
    test: {
      environment: 'jsdom',
      exclude: [...configDefaults.exclude, 'e2e/*'],
      root: fileURLToPath(new URL('./', import.meta.url)),
      // transformMode: {
      //   web: [/\.[jt]sx$/]
      // }
    }
  })
})
