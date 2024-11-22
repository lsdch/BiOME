import { UserConfig, defineConfig } from "@hey-api/openapi-ts"

const config: UserConfig = defineConfig({
  client: "@hey-api/client-fetch",
  input: "./openapi.json",
  output: {
    path: "src/api",
    format: "prettier",
    lint: "eslint",
  },
  plugins: [
    {
      name: "@hey-api/schemas",
      nameBuilder(name, schema) {
        return `$${name}`
      },
    },
    {
      name: "@hey-api/typescript",
      style: "PascalCase"
    },
    {
      name: "@hey-api/sdk",
      asClass: true,
    },
    { name: "@hey-api/transformers", dates: true },
  ],
})

export default config