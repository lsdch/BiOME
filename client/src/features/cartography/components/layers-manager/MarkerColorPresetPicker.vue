<template>
  <v-menu>
    <template #activator="{ props, isActive }">
      <v-btn variant="outlined" color="" rounded="md" :height="50" v-bind="props">
        <SvgCircle :stroke-color="layer.config.color" :fill-color="layer.config.fillColor" />
        <template #append>
          <v-icon :icon="isActive ? 'mdi-menu-up' : 'mdi-menu-down'" size="large" />
        </template>
      </v-btn>
    </template>
    <v-card :width="300">
      <v-card-text class="d-flex flex-wrap ga-1">
        <v-btn
          v-for="color in markerColorPalette"
          :key="color"
          width="20%"
          variant="text"
          rounded="sm"
          @click="usePreset(color)"
        >
          <SvgCircle :stroke-color="color" :fill-color="withOpacity(color, 0.5)" />
        </v-btn>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { withOpacity } from '@/lib/color_brewer'
import { markerColorPalette, MarkerLayerSpec } from './map-layers'
import SvgCircle from '@/components/toolkit/ui/SvgCircle.vue'

const layer = defineModel<MarkerLayerSpec>({
  required: true
})

function usePreset(color: string) {
  layer.value.config.color = color
  layer.value.config.fillColor = withOpacity(color, 0.5)
}
</script>

<style scoped lang="scss"></style>
