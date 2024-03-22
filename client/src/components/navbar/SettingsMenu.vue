<template>
  <v-menu
    location="bottom right"
    :close-on-content-click="false"
    target="#app-bar"
    :width="400"
    :min-width="100"
  >
    <template v-slot:activator="{ props }">
      <v-btn icon="mdi-cog" v-bind="props" />
    </template>
    <v-list>
      <v-list-subheader class="text-overline"> Settings </v-list-subheader>
      <v-divider />
      <v-list-item>
        <v-switch
          label="Dark theme"
          v-model="theme.global.current.value.dark"
          @click="toggleTheme"
          color="purple"
        />
      </v-list-item>
      <v-divider />
      <v-list-item prepend-icon="mdi-tools" :to="{ name: 'app-settings' }">
        Database settings
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { useTheme } from 'vuetify'
const theme = useTheme()

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}
</script>

<style scoped></style>
