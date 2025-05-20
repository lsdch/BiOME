<template>
  <v-menu
    transition="scale-transition"
    location="right center"
    origin="left center"
    :offset="10"
    :close-on-content-click="false"
    attach=".map"
  >
    <template #activator="{ props }">
      <v-card theme="light" flat rounded="0" elevation="5" :width="30" :height="30">
        <v-btn
          icon="mdi-cog"
          variant="text"
          color=""
          :rounded="0"
          v-bind="props"
          size="small"
          :width="30"
          :height="30"
        />
      </v-card>
    </template>
    <v-card
      theme="light"
      :min-width="400"
      title="Marker settings"
      class="small-card-title opacity-80"
      :max-height="300"
      @click.stop
      @mousemove.stop
      @wheel.stop
      @touchmove.stop
      @scroll.stop
      :ripple="false"
    >
      <template #append>
        <v-btn-toggle
          mandatory
          :rounded="10"
          v-model="marker"
          density="compact"
          border="sm"
          color="success"
        >
          <v-btn text="Markers" size="x-small" density="compact" value="cluster"></v-btn>
          <v-btn text="Hex grid" size="x-small" density="compact" value="hexgrid"></v-btn>
        </v-btn-toggle>
      </template>
      <v-tabs-window v-model="marker">
        <v-tabs-window-item value="cluster">
          <v-card-text>
            <v-switch label="Clustered" v-model="markerConfig.clustered" />
            <ColorPickerMenu v-model="markerConfig.color" label="Stroke color" />
            <ColorPickerMenu v-model="markerConfig.fillColor" label="Fill color" />
            <v-slider label="Stroke width" :min="1" :max="5" v-model="markerConfig.weight" />
            <v-slider label="Radius" :min="1" :max="20" v-model="markerConfig.radius" />
          </v-card-text>
        </v-tabs-window-item>
        <v-tabs-window-item value="hexgrid">
          <v-divider />
          <HexgridConfigPanel v-model="hexgridConfig" />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { CircleMarkerOptions } from 'leaflet'
import HexgridConfigPanel, { HexgridConfig } from './HexgridConfigPanel.vue'
import ColorPickerMenu from '../toolkit/ui/ColorPickerMenu.vue'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'

export type MapLayerMode = 'markers' | 'hexgrid'

// Opacity can be controlled directly by 'color' and 'fill' properties
export type MarkerConfig = { clustered: boolean } & Overwrite<
  Omit<CircleMarkerOptions, 'opacity' | 'fillOpacity' | 'renderer'>,
  { dashArray?: string | undefined }
>

const marker = defineModel<MapLayerMode>({ required: true })
const hexgridConfig = defineModel<HexgridConfig>('hexgrid', { required: true })
const markerConfig = defineModel<MarkerConfig>('marker', { required: true })
</script>

<style scoped lang="scss"></style>
