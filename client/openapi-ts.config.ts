import { UserConfig, defineConfig, defaultPlugins } from "@hey-api/openapi-ts"

const config: Promise<UserConfig> = defineConfig({
  input: "./openapi.json",
  output: {
    path: "src/api/gen/",
    format: "prettier",
    lint: "eslint",
  },
  plugins: [
    ...defaultPlugins,
    {
      name: '@hey-api/client-fetch',
      runtimeConfigPath: './openapi-client.cfg.ts'
    },
    {
      name: '@tanstack/vue-query',
    },
    {
      name: "@hey-api/schemas",
      nameBuilder(name, schema) {
        return `$${name}`
      },
    },
    {
      name: "@hey-api/typescript",
      style: "PascalCase",
      readOnlyWriteOnlyBehavior: 'off',
    },
    {
      name: "@hey-api/sdk",
      asClass: true,
      transformer: true,
    },
    { name: "@hey-api/transformers", dates: true, bigInt: false },
  ],
})

export default config