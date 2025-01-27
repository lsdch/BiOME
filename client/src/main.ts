// Setup API client
import { client } from "@/api"
import { servers } from "../openapi.json"

// ❗ Use dynamic imports for anything relying on the client in the main script
// so that the updated config is used
client.setConfig({ baseUrl: `${servers[0].url}` })

import { createApp } from 'vue'
import App from './App.vue'

// Create app instance
const app = createApp(App)


// Setup TanStack VueQuery
import { QueryClient, VueQueryPlugin, VueQueryPluginOptions } from '@tanstack/vue-query'
const queryClient = new QueryClient()
const vueQueryPluginOptions: VueQueryPluginOptions = { queryClient }
app.use(VueQueryPlugin, vueQueryPluginOptions)


import { instanceSettingsOptions, listCountriesOptions } from "./api/gen/@tanstack/vue-query.gen"
// Prefetch instance settings
const settings = await queryClient.fetchQuery({
  ...instanceSettingsOptions(),
  gcTime: Infinity
}).catch((error) => {
  throw new Error("Failed to fetch instance settings", error)
})

// Prefetch countries
queryClient.prefetchQuery({
  ...listCountriesOptions(),
  gcTime: Infinity
})


// Setup router
import setupRouter from './router'
app.use(setupRouter(settings))

// Setup vuetify
import vuetify from "./vuetify"
app.use(vuetify)



// Setup pinia stores
import { createPinia, setActivePinia } from 'pinia'
const pinia = createPinia()
setActivePinia(pinia)
app.use(pinia)


// Setup authentication using refresh token
await import('./stores/user').then(async ({ useUserStore }) => {
  const { refresh, isAuthenticated, sessionExpired, getUser } = useUserStore()
  await getUser()
  client.interceptors.request.use(async (request) => {
    if (isAuthenticated && sessionExpired && !request.headers.has('noAuthRefresh')) {
      await refresh()
    }
    return request
  })
})



app.mount('#app')
