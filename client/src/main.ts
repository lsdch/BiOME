
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Vuetify
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css' // Ensure you are using css-loader

import { createVuetify, type ThemeDefinition } from 'vuetify'
import * as components from 'vuetify/components'
import { VSkeletonLoader } from 'vuetify/labs/VSkeletonLoader'
import {
  VDataTable,
  VDataTableServer,
  VDataTableVirtual,
} from "vuetify/labs/VDataTable";
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import { md3 } from 'vuetify/blueprints'

const lightTheme: ThemeDefinition = {
  dark: false,
  colors: {
    // background: '#FFFFFF',
    // surface: '#FFFFFF',
    primary: '#1071B0',
    // 'primary-darken-1': '#3700B3',
    // secondary: '#03DAC6',
    // 'secondary-darken-1': '#018786',
    // error: '#B00020',
    // info: '#2196F3',
    // success: '#4CAF50',
    // warning: '#FB8C00',
  },
}


const vuetify = createVuetify({
  blueprint: md3,
  components: {
    VSkeletonLoader,
    VDataTable,
    VDataTableServer,
    VDataTableVirtual,
    ...components
  },
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi }
  },
  theme: {
    defaultTheme: "lightTheme",
    themes: {
      lightTheme
    }
  },
  defaults: {
    VTextField: { variant: "outlined" }
  }
})

const app = createApp(App)

app.use(router)
app.use(vuetify)

app.mount('#app')
