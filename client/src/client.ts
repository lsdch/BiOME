import { createClient } from '@hey-api/client-fetch'
import { servers } from "../openapi.json"

createClient({
  baseUrl: `${servers[0].url}`,
  global: true
})
