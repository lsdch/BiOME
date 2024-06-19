import { UserConfig, defineConfig } from "@hey-api/openapi-ts"

const config: UserConfig = defineConfig({
  client: "@hey-api/client-fetch",
  input: "./openapi.json",
  output: {
    path: "src/api",
    format: "prettier",
    lint: "eslint"
  },
  services: {
    asClass: true,
  },
  types: {
    dates: true,
    name: "PascalCase",
  }
})

export default config