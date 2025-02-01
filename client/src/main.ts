
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
import { useUserStore } from './stores/user'
const pinia = createPinia()
setActivePinia(pinia)
app.use(pinia)
useUserStore().refreshSession()


app.mount('#app')
