<template>
  <v-tabs v-model="tab" density="compact" class="ma-2" inset grow>
    <v-tab value="color">Color</v-tab>
    <v-tab value="radius">Radius</v-tab>
    <v-tab value="stroke">Stroke</v-tab>
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
        <v-list-item title="Radius" :subtitle="`${layer.config.radius} km`">
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
        <ListItemInput title="Coverage">
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
    <v-tabs-window-item value="stroke">
      <v-list>
        <v-list-item title="Width">
          <template #append>
            <v-slider
              v-model="layer.config.strokeWidth"
              :min="0"
              :max="5"
              :step="0.5"
              :width="250"
              hide-details
              thumb-label
            >
              <template #thumb-label="{ modelValue }"> {{ modelValue }}px </template>
            </v-slider>
          </template>
        </v-list-item>
        <v-list-item title="Opacity">
          <template #append>
            <v-slider
              v-model="layer.config.strokeOpacity"
              :min="0"
              :max="1"
              :step="0.1"
              :width="250"
              hide-details
              thumb-label
            >
              <template #thumb-label="{ modelValue }">
                {{ (modelValue * 100).toFixed(0) }}%
              </template>
            </v-slider>
          </template>
        </v-list-item>
      </v-list>
    </v-tabs-window-item>
    <v-tabs-window-item value="hover">
      <v-list>
        <v-list-item title="Fill cell">
          <template #prepend>
            <v-checkbox v-model="layer.config.hover.fill" hide-details />
          </template>
        </v-list-item>
        <v-list-item title="Upscale">
          <template #prepend>
            <v-checkbox v-model="layer.config.hover.useScale" hide-details />
          </template>
          <template #append>
            <v-slider
              v-model="layer.config.hover.scale"
              :disabled="!layer.config.hover.useScale"
              :min="1"
              :max="5"
              :step="0.2"
              :width="250"
              :ticks="Object.fromEntries(Array.from({ length: 5 }, (_, i) => [i + 1, `×${i + 1}`]))"
              show-ticks="always"
              hide-details
              thumb-label
            >
              <template #thumb-label="{ modelValue }"> ×{{ modelValue }} </template>
            </v-slider>
          </template>
        </v-list-item>
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

const tab = ref<'color' | 'radius' | 'stroke' | 'hover'>('color')
</script>

<style scoped lang="scss"></style>
