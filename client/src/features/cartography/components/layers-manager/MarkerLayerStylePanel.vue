<template>
  <v-list>
    <slot name="prepend-item"></slot>
    <!-- <ListItemInput label="Clustered" subtitle="Aggregate marker clusters">
      <v-switch v-model="layer.clustered" color="primary" hide-details />
    </ListItemInput> -->
    <ResolutionInputGroup v-model="layer" />
    <ListItemInput label="Radius">
      <v-slider
        :min="0.5"
        :max="10"
        :step="0.5"
        v-model="layer.config.baseRadius"
        hide-details
        :width="250"
        thumb-label
      />
    </ListItemInput>
    <ListItemInput
      label="Radius scale factor"
      :subtitle="layer.config.radiusScaleFactor ? layer.config.radiusScaleFactor : 'Disabled'"
    >
      <div class="d-flex align-center ga-2" :style="{ width: '250px' }">
        <v-slider
          class="flex-grow-1"
          :color="layer.config.radiusScaleFactor ? 'primary' : 'grey'"
          :min="0"
          :max="2"
          :step="0.1"
          v-model="layer.config.radiusScaleFactor"
          hide-details
          thumb-label
        />
        <InlineHelp class="flex-shrink-0">
          The amount of scaling applied to the radius based on the number of occurrences in a marker
          cluster.
        </InlineHelp>
      </div>
    </ListItemInput>
    <ListItemInput label="Stroke width" :subtitle="`${layer.config.weight} pixel`">
      <v-slider
        :min="1"
        :max="5"
        :step="0.5"
        v-model="layer.config.weight"
        hide-details
        :width="250"
        thumb-label
      />
    </ListItemInput>
    <ListItemInput label="Color">
      <MarkerColorPresetPicker v-model="layer" />
    </ListItemInput>
    <ListItemInput label="Stroke" subtitle="Hue and opacity">
      <ColorPickerMenu v-model="layer.config.color" hide-details show-swatches />
    </ListItemInput>
    <ListItemInput label="Fill" subtitle="Hue and opacity">
      <ColorPickerMenu v-model="layer.config.fillColor" hide-details show-swatches />
    </ListItemInput>
    <slot name="append-item"></slot>
  </v-list>
</template>

<script setup lang="ts">
import ColorPickerMenu from '@/components/toolkit/ui/ColorPickerMenu.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import { MarkerLayerSpec } from './map-layers'
import MarkerColorPresetPicker from './MarkerColorPresetPicker.vue'
import ResolutionInputGroup from './ResolutionInputGroup.vue'
import InlineHelp from '@/components/toolkit/ui/InlineHelp.vue'

const layer = defineModel<MarkerLayerSpec>({
  required: true
})
</script>

<style scoped lang="scss"></style>
