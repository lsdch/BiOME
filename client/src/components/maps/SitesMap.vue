<template>
  <l-map
    ref="map"
    v-model:zoom="zoom"
    :center="[0, 0]"
    :use-global-leaflet="true"
    v-model:bounds="mapBounds"
    :max-bounds="latLngBounds(latLng(90, -360), latLng(-90, 360))"
    :max-bounds-viscosity="1.0"
    :min-zoom="2"
    :max-zoom="18"
    @ready="onReady"
    :options="{
      gestureHandling: true,
      worldCopyJump: true,
      wheelPxPerZoomLevel: 100,
      zoomSnap: 0.5
    }"
  >
    <LControlScale position="bottomleft" metric :imperial="false" />
    <LControl position="topright" class="ma-0 d-flex justify-end">
      <v-btn
        v-if="closable"
        title="Close"
        color="white"
        :rounded="false"
        icon="mdi-close"
        :width="35"
        :height="35"
        density="compact"
        @click="emit('close')"
      />
    </LControl>

    <LControl position="topleft">
      <div class="leaflet-bar">
        <v-btn
          v-if="items"
          title="Fit view"
          color="white"
          :rounded="false"
          icon="mdi-fit-to-screen"
          :width="30"
          density="compact"
          @click="fitBounds(items)"
        />
      </div>
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
      :visible="false"
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
</template>

<script setup lang="ts" generic="SiteItem extends Geocoordinates">
import {
  LCircleMarker,
  LControl,
  LControlLayers,
  LControlScale,
  LMap,
  LMarker,
  LTileLayer
} from '@vue-leaflet/vue-leaflet'
import L, {
  CircleMarkerOptions,
  latLng,
  latLngBounds,
  LatLngExpression,
  LatLngLiteral,
  type Map
} from 'leaflet'
import 'leaflet.fullscreen'
import 'leaflet.fullscreen/Control.FullScreen.css'
import 'leaflet/dist/leaflet.css'
import { ref, watch } from 'vue'
import { Geocoordinates } from '.'
import { onKeyStroke, useFullscreen } from '@vueuse/core'
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster'

const zoom = ref(1)
const map = ref<HTMLElement>()

const { isFullscreen, enter, exit } = useFullscreen(map)
onKeyStroke('f', enter)
onKeyStroke('Escape', exit)

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
  }>(),
  {
    bounds: () => [latLng(90, -360), latLng(-90, 360)]
  }
)

const mapBounds = ref(L.latLngBounds(...props.bounds))

defineExpose({ fitBounds })

watch(
  () => props.items,
  (items) => {
    if (props.autoFit && items) fitBounds(items)
  },
  { immediate: true }
)

function onReady(mapInstance: Map) {
  console.log('Map ready')
  L.control.fullscreen().addTo(mapInstance)
  fitBounds()
}

function fitBounds(items: SiteItem[] = props.items ?? []) {
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
    mapBounds.value = latLngBounds(minMaxCoords.sw, minMaxCoords.ne).pad(0.1)
  }
}
</script>

<style lang="scss">
@use 'vuetify';
.leaflet-container {
  background-color: rgb(var(--v-theme-surface));
}
</style>
