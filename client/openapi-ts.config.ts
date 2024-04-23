import { defineConfig } from "@hey-api/openapi-ts"

export default defineConfig({
  input: "../server/docs/openapi.json",
  output: "src/api",
  client: "fetch",
  format: "prettier",
  types: {
    dates: true,
    name: "PascalCase",
  }
})