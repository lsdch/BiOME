<template>
  <v-list>
    <ListItemInput label="Clustered" subtitle="Aggregate marker clusters">
      <v-switch v-model="layer.clustered" color="primary" hide-details />
    </ListItemInput>
    <ListItemInput label="Radius">
      <v-slider
        :min="1"
        :max="20"
        :step="0.5"
        v-model="layer.config.radius"
        hide-details
        :width="250"
        thumb-label
      />
    </ListItemInput>
    <ListItemInput label="Stroke color" subtitle="Hue and opacity">
      <ColorPickerMenu v-model="layer.config.color" hide-details show-swatches />
    </ListItemInput>
    <ListItemInput label="Fill color" subtitle="Hue and opacity">
      <ColorPickerMenu v-model="layer.config.fillColor" hide-details show-swatches />
    </ListItemInput>
    <ListItemInput label="Stroke width">
      <v-slider
        :min="1"
        :max="5"
        v-model="layer.config.weight"
        hide-details
        :width="250"
        thumb-label
      />
    </ListItemInput>
  </v-list>
</template>

<script setup lang="ts">
import ColorPickerMenu from '@/components/toolkit/ui/ColorPickerMenu.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import { LayerSpec } from './map-layers'

const layer = defineModel<Extract<LayerSpec, { type: 'markers' }>>({
  required: true
})
</script>

<style scoped lang="scss"></style>
