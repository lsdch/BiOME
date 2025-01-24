import { UserConfig, defineConfig, defaultPlugins } from "@hey-api/openapi-ts"

const config: UserConfig = defineConfig({
  client: "@hey-api/client-fetch",
  input: "./openapi.json",
  output: {
    path: "src/api/gen/",
    format: "prettier",
    lint: "eslint",
  },
  plugins: [
    ...defaultPlugins,
    {
      name: '@tanstack/vue-query',
    },
    {
      name: "@hey-api/schemas",
      nameBuilder(name, schema) {
        return `$${name}`
      },
    },
    '@hey-api/transformers',
    {
      name: "@hey-api/typescript",
      style: "PascalCase"
    },
    {
      name: "@hey-api/sdk",
      asClass: true,
      transformer: true,
    },
    { name: "@hey-api/transformers", dates: true },
  ],
})

export default config