import { createApp } from 'vue'
import App from './App.vue'

// Create app instance
const app = createApp(App)

// Setup TanStack VueQuery
import { QueryClient, VueQueryPlugin, VueQueryPluginOptions } from '@tanstack/vue-query'
const queryClient = new QueryClient()
const vueQueryPluginOptions: VueQueryPluginOptions = { queryClient }
app.use(VueQueryPlugin, vueQueryPluginOptions)

// Prefetch countries
import { listCountriesOptions } from './api/gen/@tanstack/vue-query.gen'
queryClient.prefetchQuery({
  ...listCountriesOptions(),
  gcTime: Infinity
})

// Setup router
import setupRouter from './router'
app.use(setupRouter())

// Setup vuetify
import vuetify from './vuetify'
app.use(vuetify)

// Setup pinia stores
import { createPinia, setActivePinia } from 'pinia'
import { useUserStore } from './stores/user'
const pinia = createPinia()
setActivePinia(pinia)
app.use(pinia)
useUserStore().refreshSession()

app.mount('#app')
