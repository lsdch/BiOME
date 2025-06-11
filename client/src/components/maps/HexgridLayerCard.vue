<template>
  <LayerOptionsCard
    title="Hexgrid"
    :subtitle="layer.filterType"
    v-model="layer.active"
    prepend-icon="mdi-hexagon-multiple-outline"
    flat
  >
    <template #header>
      <DataFeedPicker
        v-model="layer.dataFeedID"
        class="mx-3 mb-2"
        density="compact"
        hide-details
        mandatory
      />
    </template>
    <div class="bg-main">
      <v-list class="bg-main">
        <v-list-item>
          <SiteSamplingStatusFilter
            class="my-1"
            density="compact"
            v-model="layer.filterType"
            hide-details
          />
        </v-list-item>
        <div class="d-flex align-center">
          <v-list-subheader>Color</v-list-subheader>
          <v-divider />
        </div>

        <v-list-item>
          <ScaleBindingSelect
            label="Color binding"
            density="compact"
            class="my-1"
            @update-fn="(f) => (layer.bindings.color = f)"
          />
        </v-list-item>
        <v-list-item>
          <ColorPalettePicker
            label="Palette"
            class="my-1"
            @update:model-value="(v) => (layer.config.colorRange = palette(v))"
          />
        </v-list-item>
        <v-list-item>
          <ScaleBindingSelect
            label="Opacity binding"
            density="compact"
            placeholder="Constant"
            persistent-placeholder
            clearable
            hide-details
            class="my-1"
            @update-fn="(f) => (layer.bindings.opacity = f)"
          />
        </v-list-item>
        <ListItemInput :title="layer.bindings.opacity ? 'Opacity range' : 'Opacity'">
          <v-range-slider
            v-if="layer.bindings.opacity"
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

        <div class="d-flex align-center">
          <v-list-subheader>Radius</v-list-subheader>
          <v-divider />
        </div>
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
            label="Radius binding"
            density="compact"
            placeholder="Constant"
            persistent-placeholder
            clearable
            hide-details
            class="my-1"
            @update-fn="(f) => (layer.bindings.radius = f)"
          />
        </v-list-item>
        <ListItemInput title="Radius range" v-if="layer.bindings.radius">
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

        <div class="d-flex align-center">
          <v-list-subheader>Hover</v-list-subheader>
          <v-divider />
        </div>
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
    </div>
  </LayerOptionsCard>
</template>

<script setup lang="ts">
import LayerOptionsCard from '@/views/location/LayerOptionsCard.vue'
import DataFeedPicker from '../occurrence/DataFeedPicker.vue'
import { HexgridLayerDefinition } from './map-layers'
import SiteSamplingStatusFilter from '@/views/location/SiteSamplingStatusFilter.vue'
import ScaleBindingSelect from '@/views/location/ScaleBindingSelect.vue'
import ColorPalettePicker from '../toolkit/ui/ColorPalettePicker.vue'
import { palette } from '@/functions/color_brewer'
import ListItemInput from '../toolkit/ui/ListItemInput.vue'

const layer = defineModel<HexgridLayerDefinition>({ required: true })
</script>

<style scoped lang="scss"></style>
