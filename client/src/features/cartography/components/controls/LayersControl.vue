<template>
  <v-menu
    location="top left"
    origin="top left"
    :close-on-content-click="false"
    attach="#map-container"
  >
    <template #activator="{ props }">
      <v-btn
        icon="mdi-layers"
        v-bind="props"
        color="white"
        class="bg-white"
        :rounded="false"
      ></v-btn>
    </template>
    <v-list density="compact" theme="light">
      <v-list-subheader title="Layers"> </v-list-subheader>
      <v-list-item title="Regions" prepend-icon="mdi-map">
        <template #append>
          <v-checkbox-btn color="primary" v-model="regions" hide-details />
        </template>
      </v-list-item>
      <v-list-item title="Roads" prepend-icon="mdi-road">
        <template #append>
          <v-checkbox-btn color="primary" v-model="roadsVisible" hide-details />
        </template>
      </v-list-item>
      <v-list-item
        v-if="hexgrid"
        :title="hexgrid.name ?? 'Cells'"
        prepend-icon="mdi-hexagon-multiple"
      >
        <template #append>
          <v-checkbox-btn
            color="primary"
            :model-value="hexgrid.active"
            @update:model-value="(v) => emit('toggleHexgrid', v)"
            hide-details
          />
        </template>
      </v-list-item>
      <v-list-item
        v-for="(layer, index) in markerLayers"
        :key="index"
        :title="layer.name ?? `Markers #${index + 1}`"
        prepend-icon="mdi-circle-multiple-outline"
      >
        <template #prepend>
          <svg-circle
            :fill-color="layer.config.fillColor"
            :stroke-color="layer.config.color"
          ></svg-circle>
        </template>
        <template #append>
          <v-checkbox-btn
            color="primary"
            :model-value="layer.active"
            @update:model-value="(v) => emit('toggleMarkers', index, v)"
            hide-details
          />
        </template>
      </v-list-item>
      <v-list-item
        v-if="hasSiteMarkers"
        title="Site markers"
        prepend-icon="mdi-map-marker-multiple-outline"
      >
        <template #append>
          <v-checkbox-btn color="primary" v-model="siteMarkersVisible" hide-details />
        </template>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts" generic="HexData extends H3Cell, MarkerData extends H3Cell">
import SvgCircle from '@/components/toolkit/ui/SvgCircle.vue'
import { H3Cell, HexgridLayer, MarkerLayer } from '../layers-manager/map-layers'

const regions = defineModel<boolean>('regions', {
  required: true
})
const roadsVisible = defineModel<boolean>('roads', {
  required: true
})

const siteMarkersVisible = defineModel<boolean>('siteMarkersVisible', {
  default: true
})

const { markerLayers, hexgrid } = defineProps<{
  markerLayers?: MarkerLayer<MarkerData>[]
  hexgrid?: HexgridLayer<HexData>
  hasSiteMarkers?: boolean
}>()

const emit = defineEmits<{
  toggleHexgrid: [visible: boolean]
  toggleMarkers: [index: number, visible: boolean]
}>()
</script>

<style scoped lang="scss"></style>
