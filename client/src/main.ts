// Setup API client
import { client } from "@/api/services.gen"
import { servers } from "../openapi.json"

// ❗ Use dynamic imports for anything relying on the client in the main script
// so that the updated config is used
client.setConfig({ baseUrl: `${servers[0].url}` })

import { createApp } from 'vue'
import App from './App.vue'

// Create app instance
const settings = await import('./components/settings/').then(s => s.initInstanceSettings())
const app = createApp(App, { settings })

// Setup pinia stores
import { createPinia, setActivePinia } from 'pinia'
const pinia = createPinia()
setActivePinia(pinia)
app.use(pinia)

await import('./stores/user').then(s => s.useUserStore().getUser())
await import('./stores/countries').then(s => s.useCountries().fetch())

// Setup router
import setupRouter from './router'
app.use(setupRouter(settings))

// Setup vuetify
import vuetify from "./vuetify"
app.use(vuetify)


app.mount('#app')
