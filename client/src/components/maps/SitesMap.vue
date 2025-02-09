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
      <LControl position="bottomleft" v-if="clustered">
        <MarkerControl v-model="markerMode" />
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
      <slot name="default" :map :zoom></slot>
      <template v-if="clustered">
        <LHexbinLayer
          v-if="markerMode === 'hexgrid'"
          :data="items"
          :accessor="(item) => Geocoordinates.LatLng(item)"
          :radius="10"
          :radius-range="[9.5, 10]"
          :hover="{ fill: true }"
          :color-range="['#440154', '#3b528b', '#21918c', '#5ec962', '#fde725']"
          :opacity="[0.8, 0.9]"
          style="cursor: pointer"
        ></LHexbinLayer>

        <LMarkerClusterGroup
          v-else-if="markerMode === 'cluster'"
          remove-outside-visible-bounds
          show-coverage-on-hover
          :maxClusterRadius="70"
        >
          <LCircleMarker
            v-for="item in items"
            :key="item.id"
            :lat-lng="[item.coordinates.latitude, item.coordinates.longitude]"
            v-bind="marker"
            @click="selectSite(item)"
            @popupopen="console.log('open')"
          >
          </LCircleMarker>
        </LMarkerClusterGroup>
      </template>

      <LCircleMarker
        v-else
        v-for="item in items"
        :key="item.id"
        :latLng="[item.coordinates.latitude, item.coordinates.longitude]"
        v-bind="marker"
        @click="selectSite(item)"
      >
      </LCircleMarker>

      <!-- Shared site popup -->
      <LLayerGroup ref="popup-layer" @popupopen="popupOpen = true" @popupclose="popupOpen = false">
        <KeepAlive>
          <slot name="popup" v-if="selected" :item="selected" :popupOpen :zoom> </slot>
        </KeepAlive>
      </LLayerGroup>
    </l-map>
  </div>
</template>

<script setup lang="ts" generic="SiteItem extends { id: string } & Geocoordinates">
import 'leaflet/dist/leaflet.css'
import LHexbinLayer from 'vue-leaflet-hexbin'
import 'vue-leaflet-markercluster/dist/style.css'

import {
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

import { nextTick, ref, useTemplateRef, watch } from 'vue'
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster'
import { Geocoordinates } from '.'

import { vElementVisibility } from '@vueuse/components'
import MarkerControl, { MarkerLayer } from './MarkerControl.vue'

const markerMode = ref<MarkerLayer>('cluster')

const zoom = ref(1)
const map = ref<HTMLElement>()
const popupLayer = useTemplateRef<InstanceType<typeof LLayerGroup>>('popup-layer')

const popupOpen = ref(false)

const cursorCoordinates = ref<LatLngLiteral>()

const selected = ref<SiteItem>()

function selectSite(item: SiteItem) {
  selected.value = item
  nextTick(() => popupLayer.value?.leafletObject?.openPopup(Geocoordinates.LatLng(item)))
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
    autoFit?: boolean | number
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
  default: (props: { zoom: number; map?: HTMLElement }) => any
  popup: (props: { item: SiteItem; popupOpen: boolean; zoom: number }) => any
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

const fitBounds = useDebounceFn(
  (items: SiteItem[] = props.items ?? [], fitPad: number = props.fitPad) => {
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
      minMaxCoords.sw.lat -= fitPad
      minMaxCoords.sw.lng -= fitPad
      minMaxCoords.ne.lat += fitPad
      minMaxCoords.ne.lng += fitPad
      mapBounds.value = latLngBounds(minMaxCoords.sw, minMaxCoords.ne).pad(0.1)
    }
  },
  200
)

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

.hexbin-hexagon {
  stroke: white;
  stroke-opacity: 0.5;
  stroke-width: 1;
  cursor: pointer;
}
</style>
