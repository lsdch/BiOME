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
              <div>Zoom: {{ zoom.toFixed(4) }}</div>
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
      <LMarkerClusterGroup
        v-if="clustered"
        remove-outside-visible-bounds
        show-coverage-on-hover
        :maxClusterRadius="70"
      >
        <LCircleMarker
          v-for="(item, key) in items"
          :key
          :lat-lng="[item.coordinates.latitude, item.coordinates.longitude]"
          v-bind="marker"
          @click="selectSite(item)"
        >
        </LCircleMarker>
      </LMarkerClusterGroup>
      <LCircleMarker
        v-else
        v-for="(item, key) in items"
        :key
        :latLng="[item.coordinates.latitude, item.coordinates.longitude]"
        v-bind="marker"
        @click="selectSite(item)"
      >
      </LCircleMarker>
      <slot name="default" :map></slot>
      <LCircle
        v-if="selected && showRadius(selected.coordinates.precision)"
        :lat-lng="[selected.coordinates.latitude, selected.coordinates.longitude]"
        :radius="precisionRadius(selected.coordinates.precision)"
        :interactive="false"
      ></LCircle>

      <!-- Shared site popup -->
      <LLayerGroup ref="popup-layer">
        <KeepAlive>
          <slot name="popup" v-if="selected" :item="selected"></slot>
        </KeepAlive>
      </LLayerGroup>
    </l-map>
  </div>
</template>

<script setup lang="ts" generic="SiteItem extends { id: string } & Geocoordinates">
import {
  LCircle,
  LCircleMarker,
  LControl,
  LControlLayers,
  LControlScale,
  LLayerGroup,
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
  type LeafletMouseEvent,
  type Map
} from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { ref, useTemplateRef, watch } from 'vue'
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster'
import { Geocoordinates } from '.'

import { CoordinatesPrecision } from '@/api'
import { vElementVisibility } from '@vueuse/components'

const zoom = ref(1)
const map = ref<HTMLElement>()
const popupLayer = useTemplateRef<InstanceType<typeof LLayerGroup>>('popup-layer')

const cursorCoordinates = ref<LatLngLiteral>()

const selected = ref<SiteItem>()

function selectSite(item: SiteItem) {
  selected.value = item
  popupLayer.value?.leafletObject?.openPopup(Geocoordinates.LatLng(item))
}

function showRadius(precision?: CoordinatesPrecision): precision is CoordinatesPrecision {
  switch (precision) {
    case undefined:
    case '<100m':
    case 'Unknown':
      return false
    case '<1km':
      return zoom.value > 10
    default:
      return zoom.value > 6
  }
}

function precisionRadius(precision: CoordinatesPrecision): number {
  switch (precision) {
    case '10-100km':
      return 100_000
    case '<10km':
      return 10_000
    case '<1km':
      return 1000
    default:
      return 0
  }
}

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
    autoFit: true,
    marker: () => ({
      color: 'white',
      fill: true,
      fillColor: 'orangered',
      fillOpacity: 1,
      radius: 8
    })
  }
)

defineSlots<{
  default: (map?: HTMLElement) => any
  popup: (props: { item: SiteItem }) => any
}>()

const mapBounds = ref(L.latLngBounds(...props.bounds))

watch(
  () => props.items,
  (items) => {
    if (props.autoFit && items) fitBounds(items)
  }
)

function onReady(mapInstance: Map) {
  // nextTick(fitBounds)
  if (props.autoFit) setTimeout(fitBounds, 200)
}

function onVisible(visible: boolean) {
  if (visible && props.autoFit) fitBounds()
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
