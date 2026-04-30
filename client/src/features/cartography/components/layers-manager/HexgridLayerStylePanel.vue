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
            v-model="layer.bindings.color"
            label="Color binding"
            density="compact"
            class="my-1"
          />
        </v-list-item>
        <v-list-item>
          <ColorPalettePicker v-model="layer.config.colorRange" label="Palette" class="my-1" />
        </v-list-item>
        <v-list-item>
          <ScaleBindingSelect
            v-model="layer.bindings.opacity"
            label="Opacity binding"
            density="compact"
            placeholder="Constant"
            persistent-placeholder
            clearable
            hide-details
            class="my-1"
          />
        </v-list-item>
        <ListItemInput :title="layer.bindings.opacity ? 'Opacity range' : 'Opacity'">
          <v-range-slider
            v-if="layer.bindings.opacity?.binding"
            v-model="layer.config.opacityRange"
            :min="0"
            :max="1"
            :step="0.1"
            :width="250"
            hide-details
            color="warning"
            thumb-label
          >
            <template #thumb-label="{ modelValue }"> {{ modelValue * 100 }}% </template>
          </v-range-slider>
          <v-slider
            v-else
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
        <v-list-item title="Grid cell">
          <template #append>
            <v-slider
              v-model="layer.config.radius"
              :min="2"
              :max="20"
              :step="1"
              :width="250"
              hide-details
              thumb-label
            />
          </template>
        </v-list-item>

        <v-list-item>
          <ScaleBindingSelect
            v-model="layer.bindings.radius"
            label="Radius binding"
            density="compact"
            placeholder="Constant"
            persistent-placeholder
            clearable
            hide-details
            class="my-1"
          />
        </v-list-item>
        <ListItemInput title="Radius range" v-if="layer.bindings.radius?.binding">
          <v-range-slider
            v-model="layer.config.radiusRange"
            :ticks="[layer.config.radius]"
            show-ticks="always"
            :min="2"
            :max="20"
            :step="0.5"
            :width="250"
            thumb-label
            hide-details
            color="warning"
          />
        </ListItemInput>
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

const tab = ref<'color' | 'radius' | 'hover'>('color')
</script>

<style scoped lang="scss"></style>
