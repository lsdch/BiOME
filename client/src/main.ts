
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Vuetify
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

import { createPinia, setActivePinia } from 'pinia'
import { createVuetify, type ThemeDefinition } from 'vuetify'
import { md3 } from 'vuetify/blueprints'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'


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

import { VNumberInput } from 'vuetify/labs/VNumberInput'
import { VTextField } from 'vuetify/components'
import { VTreeview } from 'vuetify/labs/VTreeview'

import { useUserStore } from './stores/user'
const vuetify = createVuetify({
  blueprint: md3,
  components: {
    ...components,
    VNumberInput,
    VTreeview
  },
  aliases: {
    VInlineSearchBar: VTextField,
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
    VTextField: { variant: "outlined" },
    VSelect: { variant: "outlined" },
    VAutocomplete: { variant: "outlined" },
    VTextArea: { variant: "outlined" },
    VNumberInput: { variant: "outlined" },
    VInlineSearchBar: {
      density: "compact",
      clearable: true,
      hideDetails: true,
      color: "primary",
      variant: "outlined",
      prependInnerIcon: "mdi-magnify"
    }
  }
})

const app = createApp(App)

// Must be
const pinia = createPinia()
setActivePinia(pinia)
app.use(pinia)

await useUserStore().getUser()

app.use(router())
app.use(vuetify)


app.mount('#app')
