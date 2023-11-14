
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Vuetify
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'

import { createVuetify, type ThemeDefinition } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import { md3 } from 'vuetify/blueprints'
import { createPinia, setActivePinia } from 'pinia'


const light: ThemeDefinition = {
  dark: false,
  colors: {
    primary: '#1071B0',
  },
}

const dark: ThemeDefinition = {
  dark: true,
  colors: {
    primary: '#057C9D',
  }
}


const vuetify = createVuetify({
  blueprint: md3,
  components: {
    ...components
  },
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi }
  },
  theme: {
    defaultTheme: "light",
    themes: {
      light, dark
    }
  },
  defaults: {
    VTextField: { variant: "outlined" }
  }
})

const app = createApp(App)

app.use(router)
app.use(vuetify)
const pinia = createPinia()
setActivePinia(pinia)
app.use(pinia)


app.mount('#app')
