import { servers } from "./openapi.json"
import type { CreateClientConfig } from '@/api/gen/client/';

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseUrl: servers[0].url
})