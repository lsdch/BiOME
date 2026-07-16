import { createApp } from 'vue'
import AppComponent from './App.vue'

import { QueryClient, VueQueryPlugin, VueQueryPluginOptions } from '@tanstack/vue-query'
import { createPinia, setActivePinia } from 'pinia'
import { getInstanceSettingsOptions, listCountriesOptions } from './api/gen/@tanstack/vue-query.gen'
import { settingsKey } from './lib/injection.ts'
import setupRouter from './router'
import { useUserStore } from './stores/user'
import vuetify from './vuetify'

// Create app instance
const app = createApp(AppComponent)

// Setup vuetify
app.use(vuetify)

// Setup TanStack VueQuery
const queryClient = new QueryClient()
const vueQueryPluginOptions: VueQueryPluginOptions = { queryClient }
app.use(VueQueryPlugin, vueQueryPluginOptions)

// Setup pinia stores
const pinia = createPinia()
setActivePinia(pinia)
app.use(pinia)

// Prefetch countries
queryClient.prefetchQuery({ ...listCountriesOptions(), gcTime: Infinity })

// Bootstrap authentication and instance settings
const userStore = useUserStore()
await Promise.all([
  userStore.bootstrapAuth(queryClient),
  initSettings(queryClient).then((settings) => {
    app.provide(settingsKey, settings)
    app.use(setupRouter(settings))
  })
])

app.mount('#app')

async function initSettings(client: QueryClient) {
  return client
    .fetchQuery({
      ...getInstanceSettingsOptions()
    })
    .catch((err) => {
      throw new Error(`Failed to fetch instance settings: ${err}`)
    })
}
