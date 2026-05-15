<template>
  <v-tabs v-model="tab" density="compact" class="ma-2" inset grow>
    <v-tab value="color">Color</v-tab>
    <v-tab value="radius">Radius</v-tab>
    <v-tab value="hover">Hover</v-tab>
  </v-tabs>
  <v-tabs-window v-model="tab">
    <v-tabs-window-item value="color">
      <v-list>
        <v-list-item>
          <ScaleBindingSelect
            v-model="layer.colorBinding"
            label="Color binding"
            density="compact"
            class="my-1"
          />
        </v-list-item>
        <v-list-item>
          <ColorPalettePicker v-model="layer.config.colorRange" label="Palette" class="my-1" />
        </v-list-item>
        <ListItemInput title="Opacity">
          <v-slider
            v-model="layer.config.opacity"
            :min="0"
            :max="1"
            :step="0.1"
            hide-details
            :width="250"
            thumb-label
          >
            <template #thumb-label="{ modelValue }"> {{ modelValue * 100 }}% </template>
          </v-slider>
        </ListItemInput>
      </v-list>
    </v-tabs-window-item>
    <v-tabs-window-item value="radius">
      <v-list>
        <v-list-item
          title="Radius"
          :subtitle="`${layer.config.radius} km (area ${Math.round(hexagonArea(layer.config.radius))} km²)`"
          lines="two"
        >
          <template #append>
            <v-slider
              v-model="layer.config.radius"
              density="compact"
              :min="5"
              :max="200"
              :step="1"
              :width="250"
              :ticks="[10, 50, 100, 150, 200]"
              show-ticks
              hide-details
            />
          </template>
        </v-list-item>
        <ListItemInput title="Coverage" :subtitle="`${(layer.config.coverage ?? 1) * 100}%`">
          <v-slider
            v-model="layer.config.coverage"
            :min="0.5"
            :max="1"
            :step="0.05"
            :width="250"
            :ticks="{ 0.5: '50%', 0.75: '75%', 1: '100%' }"
            show-ticks
            hide-details
          />
        </ListItemInput>
      </v-list>
    </v-tabs-window-item>
    <v-tabs-window-item value="hover">
      <v-list density="compact">
        <ListItemInput title="Show tooltip">
          <v-switch v-model="layer.config.hover.showTooltip" color="primary" hide-details />
        </ListItemInput>
        <ListItemInput title="Highlight on hover">
          <v-switch v-model="layer.config.hover.highlight" color="primary" hide-details />
        </ListItemInput>
      </v-list>
    </v-tabs-window-item>
  </v-tabs-window>
</template>

<script setup lang="ts">
import ColorPalettePicker from '@/components/toolkit/ui/ColorPalettePicker.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import { ref } from 'vue'
import { HexLayerSpec } from './map-layers'
import ScaleBindingSelect from './ScaleBindingSelect.vue'

const layer = defineModel<HexLayerSpec>({
  required: true
})

const tab = ref<'color' | 'radius' | 'hover'>('color')

function hexagonArea(radius: number) {
  return (3 * Math.sqrt(3) * radius ** 2) / 2
}
</script>

<style scoped lang="scss"></style>
