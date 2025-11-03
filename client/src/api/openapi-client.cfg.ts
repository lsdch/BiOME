import { servers } from "../../openapi.json"
import type { CreateClientConfig } from '@/api/gen/client/';

if (servers.length === 0) {
  throw new Error("No servers defined in OpenAPI specification.")
}

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseUrl: servers[0]!.url
})