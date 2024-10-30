<template>
  <div
    ref="map"
    class="fill-height"
    v-element-visibility="onVisible"
    @mouseleave="cursorCoordinates = undefined"
  >
    <l-map
      v-model:zoom="zoom"
      v-bind="$attrs"
      :center="[0, 0]"
      :use-global-leaflet="true"
      v-model:bounds="mapBounds"
      :max-bounds="latLngBounds(latLng(90, -360), latLng(-90, 360))"
      :max-bounds-viscosity="1.0"
      :min-zoom="2"
      :max-zoom="18"
      @ready="onReady"
      @mousemove="
        ({ latlng }: LeafletMouseEvent) => {
          cursorCoordinates = latlng
        }
      "
      :options="{
        gestureHandling: true,
        worldCopyJump: true,
        wheelPxPerZoomLevel: 100,
        zoomSnap: 0.5
      }"
    >
      <LControlScale position="bottomright" metric :imperial="false" />
      <LControl position="bottomright">
        <v-card v-if="cursorCoordinates" class="coordinates-control">
          <template #text>
            <code>
              <div>Lat: {{ cursorCoordinates.lat.toFixed(5) }}</div>
              <div>Lng: {{ cursorCoordinates.lng.toFixed(5) }}</div>
            </code>
          </template>
        </v-card>
      </LControl>

      <LControl position="topright" class="ma-0 d-flex justify-end">
        <v-btn
          v-if="closable"
          title="Close"
          color="white"
          class="bg-white"
          :rounded="false"
          icon="mdi-close"
          :width="35"
          :height="35"
          density="compact"
          @click="emit('close')"
        />
      </LControl>

      <LControl position="topleft">
        <div class="leaflet-bar d-flex flex-column">
          <v-btn
            v-if="items"
            title="Fit view"
            class="bg-white"
            color="white"
            :rounded="false"
            icon="mdi-fit-to-screen"
            :width="30"
            density="compact"
            @click="fitBounds(items)"
          />
          <v-btn
            title="Toggle fullscreen"
            color="white"
            class="bg-white"
            :rounded="false"
            :icon="isFullscreen ? 'mdi-fullscreen-exit' : 'mdi-fullscreen'"
            :width="30"
            density="compact"
            @click="toggleFullscreen"
          />
        </div>
      </LControl>
      <LControl position="topleft">
        <div class="leaflet-bar"></div>
      </LControl>

      <LControlLayers hide-single-base />

      <l-tile-layer
        :subdomains="['server', 'services']"
        url="https://{s}.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}"
        attributionUrl="https://static.arcgis.com/attribution/World_Imagery"
        attribution="Powered by <a href='https://www.esri.com/'>Esri</a> &mdash; Source: Esri, Maxar, Earthstar Geographics, and the GIS User Community"
        layer-type="base"
        name="Base layer"
      />
      <l-tile-layer
        :subdomains="['server', 'services']"
        url="https://{s}.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile/{z}/{y}/{x}"
        layer-type="overlay"
        name="Regions"
        :opacity="0.75"
        :visible="regions"
      />
      <LMarkerClusterGroup v-if="clustered" remove-outside-visible-bounds show-coverage-on-hover>
        <LCircleMarker
          v-for="(item, key) in items"
          :key
          :lat-lng="[item.coordinates.latitude, item.coordinates.longitude]"
          v-bind="marker"
        >
          <slot name="marker" :item />
        </LCircleMarker>
      </LMarkerClusterGroup>
      <LCircleMarker
        v-else
        v-for="(item, key) in items"
        :key
        :latLng="[item.coordinates.latitude, item.coordinates.longitude]"
        v-bind="marker"
      >
        <slot name="marker" :item />
      </LCircleMarker>
      <slot :map></slot>
    </l-map>
  </div>
</template>

<script setup lang="ts" generic="SiteItem extends Geocoordinates">
import {
  LCircleMarker,
  LControl,
  LControlLayers,
  LControlScale,
  LMap,
  LTileLayer
} from '@vue-leaflet/vue-leaflet'
import { onKeyStroke, useDebounceFn, useFullscreen, useThrottleFn } from '@vueuse/core'
import L, {
  CircleMarkerOptions,
  latLng,
  latLngBounds,
  LatLngExpression,
  LatLngLiteral,
  type Map,
  type LeafletMouseEvent
} from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { ref, watch } from 'vue'
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster'
import { Geocoordinates } from '.'

import { vElementVisibility } from '@vueuse/components'

const zoom = ref(1)
const map = ref<HTMLElement>()

const cursorCoordinates = ref<LatLngLiteral>()

const { isFullscreen, enter, exit, toggle } = useFullscreen(map, {})
onKeyStroke('Escape', exit)
const toggleFullscreen = useThrottleFn(toggle)

const emit = defineEmits<{
  close: []
}>()

const props = withDefaults(
  defineProps<{
    items?: SiteItem[]
    marker?: Omit<CircleMarkerOptions, 'dashArray'>
    bounds?: [LatLngExpression, LatLngExpression]
    autoFit?: boolean
    closable?: boolean
    clustered?: boolean
    regions?: boolean
    fitPad?: number
  }>(),
  {
    bounds: () => [latLng(90, -360), latLng(-90, 360)],
    fitPad: 0,
    marker: () => ({
      color: 'white',
      fill: true,
      fillColor: 'orangered',
      fillOpacity: 1,
      radius: 8
    })
  }
)

const mapBounds = ref(L.latLngBounds(...props.bounds))

watch(
  () => props.items,
  (items) => {
    if (props.autoFit && items) fitBounds(items)
  }
)

function onReady(mapInstance: Map) {
  // nextTick(fitBounds)
  setTimeout(fitBounds, 200)
}

function onVisible(visible: boolean) {
  if (visible) fitBounds()
}

const fitBounds = useDebounceFn((items: SiteItem[] = props.items ?? []) => {
  console.log('[Map] Fit bounds')
  const minMaxCoords = items.reduce(
    (
      acc: { sw: LatLngLiteral; ne: LatLngLiteral } | null,
      { coordinates: { latitude, longitude } }: SiteItem
    ): { sw: LatLngLiteral; ne: LatLngLiteral } | null => {
      return acc === null
        ? {
            sw: { lat: latitude, lng: longitude },
            ne: { lat: latitude, lng: longitude }
          }
        : {
            sw: {
              lat: Math.min(acc.sw.lat, latitude),
              lng: Math.min(acc.sw.lng, longitude)
            },
            ne: {
              lat: Math.max(acc.ne.lat, latitude),
              lng: Math.max(acc.ne.lng, longitude)
            }
          }
    },
    null
  )

  if (minMaxCoords) {
    minMaxCoords.sw.lat -= props.fitPad
    minMaxCoords.sw.lng -= props.fitPad
    minMaxCoords.ne.lat += props.fitPad
    minMaxCoords.ne.lng += props.fitPad
    mapBounds.value = latLngBounds(minMaxCoords.sw, minMaxCoords.ne).pad(0.1)
  }
}, 200)

defineExpose({ fitBounds })
</script>

<style lang="scss">
@use 'vuetify';
.leaflet-container {
  background-color: rgb(var(--v-theme-surface));
}

.coordinates-control {
  background-color: rgb(var(--v-theme-surface), 0.5);
  code * {
    opacity: 1;
  }
}
</style>
