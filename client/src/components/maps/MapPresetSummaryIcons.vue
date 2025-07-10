<template>
  <v-badge :content="spec.feeds.length" color="purple" :offset-x="-3" :offset-y="-3">
    <v-icon icon="mdi-database-arrow-right"></v-icon>
  </v-badge>
  <v-icon icon="mdi-hexagon-multiple" :disabled="!spec.hexgrid.active"></v-icon>
  <v-badge :content="spec.feeds.length" color="success" :offset-x="-3" :offset-y="-3">
    <v-icon icon="mdi-circle-multiple"></v-icon>
  </v-badge>
</template>

<script setup lang="ts">
import { MapPresetSpec } from './map-presets'

defineProps<{
  spec: MapPresetSpec
}>()
</script>

<style scoped lang="scss"></style>
