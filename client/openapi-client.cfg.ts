import { servers } from "./openapi.json"
import type { CreateClientConfig } from '@hey-api/client-fetch';

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseUrl: servers[0].url
})