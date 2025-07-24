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
          class="px-3"
          label="Dark theme"
          :model-value="theme.name.value"
          @update:model-value="(v) => theme.change(v!)"
          false-value="light"
          true-value="dark"
          color="purple"
          :indeterminate="false"
          hide-details
        />
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { useTheme } from 'vuetify'
const theme = useTheme()

watch(theme.name, () => {
  localStorage.setItem('app-theme', theme.name.value)
})
</script>

<style scoped></style>
