import { UserConfig, defineConfig, defaultPlugins } from "@hey-api/openapi-ts"

const config: Promise<UserConfig> = defineConfig({
  input: "./openapi.json",
  output: {
    path: "src/api/gen/",
    postProcess: ["eslint", "prettier"],
  },
  parser: {
    transforms: {

      readWrite: {
        enabled: false,
      }
    }
  },
  plugins: [
    ...defaultPlugins,
    {
      name: '@hey-api/client-fetch',
      runtimeConfigPath: '../openapi-client.cfg.ts'
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
    },
    {
      name: "@hey-api/sdk",
      operations: {
        strategy: 'byTags',
        containerName(name) {
          return `${name}Service`
        },
      },
      transformer: true,
    },
    { name: "@hey-api/transformers", dates: true, bigInt: false },
  ],
})

export default config