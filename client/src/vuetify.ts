import "./styles/main.scss"

// Vuetify
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

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

export default createVuetify({
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
  display: {
    mobileBreakpoint: 'sm',
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
    VTextarea: { variant: "outlined" },
    VAlert: { variant: "tonal" },
    VNumberInput: {
      variant: "outlined",
      controlVariant: "stacked",
      VBtn: {
        color: undefined,
        rounded: 0
      }
    },
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
