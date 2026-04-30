<template>
  <!-- {{ polyline }} -->
  <div
    ref="mapContainer"
    :class="['fill-height', { 'polygon-mode': polygonMode }]"
    v-element-visibility="onVisible"
    @mouseleave="cursorCoordinates = undefined"
  >
    <l-map
      class="map"
      v-model:zoom="zoom"
      v-bind="$attrs"
      :use-global-leaflet="true"
      v-model:bounds="mapBounds"
      :max-bounds="latLngBounds(latLng(90, -360), latLng(-90, 360))"
      :max-bounds-viscosity="1.0"
      :min-zoom
      :max-zoom
      @ready="onReady"
      @mousemove="
        ({ latlng }: LeafletMouseEvent) => {
          cursorCoordinates = latlng
        }
      "
      :center
      :options="{
        gestureHandling: true,
        worldCopyJump: true,
        wheelPxPerZoomLevel: 100,
        zoomSnap: 0.5
      }"
      @click="
        (e: LeafletMouseEvent) => {
          if (polygonMode) {
            addPolylinePoint(e.latlng)
          }
        }
      "
    >
      <LControlScale position="bottomright" metric :imperial="false" />
      <LControl position="bottomright" class="coordinates-control">
        <v-card v-if="cursorCoordinates" density="compact" class="pa-2">
          <code class="text-caption font-monospace">
            <div class="d-flex justify-space-between ga-2">
              <span>Lat:</span> {{ cursorCoordinates.lat.toFixed(4) }}
            </div>
            <div class="d-flex justify-space-between ga-2">
              <span>Lng:</span> {{ cursorCoordinates.lng.toFixed(4) }}
            </div>
            <div class="d-flex justify-space-between ga-2">
              <span>Zoom:</span> {{ zoom.toFixed(2) }}
            </div>
          </code>
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
            title="Fit view"
            class="bg-white"
            color="white"
            :rounded="false"
            icon="mdi-fit-to-screen"
            :width="30"
            density="compact"
            @click="fitMapView()"
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
      <!-- <LControl position="topleft">
        <div class="leaflet-bar">
          <MarkerControl
            v-if="!hideMarkerControl"
            v-model="markerMode"
            v-model:hexgrid="hexgridConfig"
            v-model:marker="markerConfig"
          />
        </div>
      </LControl> -->

      <LControlLayers hide-single-base />

      <LControl position="topleft" :options="{}" v-if="hexgrid.active && hexgridColorRange?.length">
        <MapColorLegend class="mx-3" :range="hexgridColorRange" />
      </LControl>

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
      <slot name="default" :mapContainer :zoom></slot>

      <!-- Hexagon layer  -->
      <LHexbinLayer
        v-if="hexgrid.active"
        pane="overlayPane"
        :data="hexgrid.data"
        :accessor="(item) => Geocoordinates.LatLng(item)"
        :radius="hexgrid.config.radius"
        :radius-range="
          hexgrid.bindings?.radius
            ? hexgrid.config.radiusRange
            : [hexgrid.config.radius, hexgrid.config.radius]
        "
        :opacity="hexgrid.bindings?.opacity ? hexgrid.config.opacityRange : hexgrid.config.opacity"
        :hover-fill="hexgrid.config.hover.fill"
        :hover-scale="hexgrid.config.hover.useScale ? hexgrid.config.hover.scale : undefined"
        :color-range="hexgrid.config.colorRange"
        style="cursor: pointer"
        :color-binding="hexgrid.bindings?.color"
        :opacity-binding="hexgrid.bindings?.opacity"
        :radius-binding="hexgrid.bindings?.radius"
        @update:color-scale-extent="(range) => console.log('color scale extent', range)"
        @click="
          (e) => {
            if (e.length === 1) selectSite(e[0]!.data)
          }
        "
      >
        <template #popup="{ data }">
          <LPopup v-show="(data?.length ?? 0) > 1" :options="{ closeButton: false }">
            <slot name="hex-popup" :data>
              <v-card-text class="text-center">
                <code class="font-weight-bold">{{ data?.length }}</code>
              </v-card-text>
            </slot>
          </LPopup>
        </template>
      </LHexbinLayer>
      <LLayerGroup v-for="layer in markerLayers" pane="markerPane">
        <template v-if="layer.active">
          <LMarkerClusterGroup
            v-if="layer.clustered"
            remove-outside-visible-bounds
            show-coverage-on-hover
            :maxClusterRadius="layer.maxClusterRadius ?? 80"
          >
            <LCircleMarker
              v-for="item in unref(layer.data)"
              pane="markerPane"
              :key="item.id"
              :lat-lng="[item.coordinates.latitude, item.coordinates.longitude]"
              v-bind="layer.config"
              :opacity="1"
              :fill-opacity="1"
              @click="selectSite(item)"
              @popupopen="console.log('open')"
            >
            </LCircleMarker>
          </LMarkerClusterGroup>

          <!-- Marker layers -->
          <LCircleMarker
            v-else
            pane="markerPane"
            v-for="item in unref(layer.data)"
            :latLng="[item.coordinates.latitude, item.coordinates.longitude]"
            v-bind="layer.config"
            :opacity="1"
            :fill-opacity="1"
            :options="{
              zIndexOffset: 10
            }"
            @click="selectSite(item)"
          >
          </LCircleMarker>
        </template>
      </LLayerGroup>

      <slot v-if="markers" name="markers" :markers>
        <LMarker
          v-for="marker in markers"
          :lat-lng="Geocoordinates.LatLng(marker)"
          :icon="
            L.divIcon({
              className: '',
              iconSize: [36, 33],
              iconAnchor: [18, 33],
              shadowSize: [5, 5],
              html: `<i class='mdi-map-marker mdi marker-icon v-icon v-icon--size-x-large' style='color: ${marker.color ?? 'orangered'}; text-shadow: black 1px 0 5px;'></i>`
            }) as L.Icon
          "
          v-bind="markerProps"
        />
      </slot>

      <!-- Shared site popup -->
      <LLayerGroup ref="popup-layer" @popupopen="popupOpen = true" @popupclose="popupOpen = false">
        <KeepAlive>
          <slot name="popup" v-if="selected" :item="selected" :popupOpen :zoom> </slot>
        </KeepAlive>
      </LLayerGroup>

      <LPolygon
        v-if="!polygonMode && polyline.length > 2"
        :lat-lngs="polyline"
        color="orangered"
        :weight="2"
        fill
        no-clip
        fill-rule=""
        :fill-opacity="0.3"
      />
      <LPolyline
        v-if="polygonMode && polyline.length"
        :lat-lngs="[...polyline, ...(cursorCoordinates ? [cursorCoordinates] : [])]"
        color="orangered"
        :weight="2"
        fill
        no-clip
        fill-rule=""
        :fill-opacity="0.3"
        :interactive="false"
      />
      <LCircleMarker
        v-if="polygonMode || polyline.length > 1"
        v-for="(latLng, i) in polyline"
        interactive
        :lat-lng
        :radius="i === 0 || i === polyline.length - 1 ? 6 : 3"
        fill
        :fill-opacity="1"
        :fillColor="i === 0 ? 'green' : 'orangered'"
        :color="i === 0 ? 'green' : 'orangered'"
        @click="
          () => {
            console.log('click', i)
            if (polygonMode && i === 0) {
              if (polyline.length == 2) {
                clearPolyline()
              }
              polygonMode = false
            }
          }
        "
      />
    </l-map>
  </div>
</template>

<script lang="ts">
/**
 * Base map component for displaying geolocated items on a map.
 * Supports hexbin layers, marker layers, and custom popups.
 */
export default {
  name: 'BaseMap'
}
</script>

<script setup lang="ts" generic="Item extends { id: UUID } & Geocoordinates">
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
  LMarker,
  LPolygon,
  LPolyline,
  LPopup,
  LTileLayer
} from '@vue-leaflet/vue-leaflet'
import { onKeyStroke, useDebounceFn, useFullscreen, useThrottleFn } from '@vueuse/core'
import L, {
  latLng,
  latLngBounds,
  LatLngExpression,
  LatLngLiteral,
  PointExpression,
  type LeafletMouseEvent,
  type Map
} from 'leaflet'

import { nextTick, ref, unref, UnwrapRef, useTemplateRef, watch } from 'vue'
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster'
import { Geocoordinates } from '../coordinates'

import { Coordinates, CoordinatesPrecision } from '@/api'
import MapColorLegend from '@/features/cartography/components/MapColorLegend.vue'
import { vElementVisibility } from '@vueuse/components'
import { HexgridLayer, MarkerLayer } from './layers-manager/map-layers'

export type HexPopupData<Item> = {
  data: Item
  coord: L.LatLngExpression
}

const markerLayers = defineModel<MarkerLayer<Item>[]>('marker-layers')

const hexgrid = defineModel<HexgridLayer<Item>>('hexgrid', {
  default: (): HexgridLayer<Item> => ({
    active: true,
    config: {
      radius: 10,
      radiusRange: [10, 10],
      opacity: 0.8,
      hover: {
        fill: true,
        scale: 1,
        useScale: false
      }
    },
    bindings: {}
  })
})

const hexgridColorRange = ref<[number, number]>()

const zoom = ref(1)
const mapContainer = useTemplateRef<HTMLElement, 'mapContainer'>('mapContainer')
const popupLayer = useTemplateRef<InstanceType<typeof LLayerGroup>>('popup-layer')

const popupOpen = ref(false)

const cursorCoordinates = ref<LatLngLiteral>()

const selected = ref<Item>()

function selectSite(item: Item) {
  selected.value = item
  nextTick(() => popupLayer.value?.leafletObject?.openPopup(Geocoordinates.LatLng(item)))
}

const { isFullscreen, enter, exit, toggle } = useFullscreen(mapContainer, {})
onKeyStroke('Escape', exit)
const toggleFullscreen = useThrottleFn(toggle)

const emit = defineEmits<{
  close: []
}>()

const props = withDefaults(
  defineProps<{
    /**
     * Place a single marker at the given coordinates.
     */
    markers?: (Geocoordinates & { color?: string })[]
    markerProps?: Partial<InstanceType<typeof LMarker>['$props']> & {
      'onUpdate:latLng'?: (latLng: LatLngLiteral) => void
    }
    bounds?: [LatLngExpression, LatLngExpression]
    autoFit?: boolean | number
    closable?: boolean
    regions?: boolean
    center?: PointExpression
    minZoom?: number
    maxZoom?: number
  }>(),
  {
    bounds: () => [latLng(90, -360), latLng(-90, 360)],
    minZoom: 2,
    maxZoom: 18,
    autoFit: true,
    center: () => [0, 0]
  }
)

const slots = defineSlots<{
  default: (props: { zoom: number; mapContainer: HTMLElement | null }) => any
  popup: (props: { item: Item; popupOpen: boolean; zoom: number }) => any
  markers: (props: { markers: Geocoordinates[] }) => any
  'hex-popup': (props: { data?: HexPopupData<UnwrapRef<Item>>[] }) => any
}>()

const mapBounds = ref(L.latLngBounds(...props.bounds))

watch(
  () => [props.markers, props.autoFit, hexgrid.value.data, markerLayers.value?.map((l) => l.data)],
  () => fitMapView()
)

watch(
  () => markerLayers.value?.filter((l) => l.active).length,
  (active, previouslyActive) => {
    if ((active ?? 0) > (previouslyActive ?? 0)) {
      fitMapView()
    }
  }
)

function fitViewToSite({ coordinates: { latitude, longitude } }: Geocoordinates) {
  mapBounds.value = L.latLng(latitude, longitude).toBounds(500).pad(0.5)
}

function onReady(mapInstance: Map) {
  // nextTick(fitBounds)
  if (props.autoFit) setTimeout(fitMapView, 200)
}

function onVisible(visible: boolean) {
  if (!visible) return
  fitMapView()
}

function fitMapView() {
  if (props.autoFit !== false) {
    if (typeof props.autoFit == 'number') {
      fitBounds(props.autoFit)
    } else {
      fitBounds()
    }
  }
}

const fitRadius = useDebounceFn((coords: Coordinates, radius: number) => {
  if (props.markers?.length === 1) {
    let r = Math.max(100, radius)
    const { latitude, longitude } = coords
    mapBounds.value = L.latLng(latitude, longitude)
      .toBounds(r + 100)
      .pad(0.5)
  }
}, 200)

function computeBounds(items: Geocoordinates[]) {
  const minMaxCoords = items.reduce(
    (
      acc: { sw: LatLngLiteral; ne: LatLngLiteral } | null,
      { coordinates: { latitude, longitude } }: Geocoordinates
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
  return minMaxCoords ? latLngBounds(minMaxCoords.sw, minMaxCoords.ne) : undefined
}

const fitBounds = useDebounceFn((radius: number = 0) => {
  console.log('[Map] Fit bounds')
  let bounds = hexgrid.value.active ? computeBounds(unref(hexgrid.value.data) ?? []) : undefined
  markerLayers.value?.forEach((layer) => {
    if (!layer.active || !layer.data?.length) return
    const b = computeBounds(layer.data ?? [])
    if (bounds) {
      bounds.extend(b ?? [])
      return
    } else {
      bounds = b
    }
  })
  props.markers?.forEach(({ coordinates: { latitude, longitude, precision } }) => {
    const ll = latLng(latitude, longitude)
    if (precision) {
      const p = Math.max(CoordinatesPrecision.radius(precision), 100)
      const pointBounds = ll.toBounds(p + radius)
      if (bounds) {
        bounds.extend(pointBounds)
      } else {
        bounds = pointBounds
      }
      return
    }
    if (bounds) {
      bounds.extend(ll)
    } else {
      bounds = L.latLngBounds(ll, ll)
    }
  })

  if (bounds) {
    mapBounds.value = bounds.pad(0.1)
  }
}, 200)

const polygonMode = defineModel<boolean>('polygon-mode', { default: false })
const polyline = ref<LatLngExpression[]>([])
function addPolylinePoint(latlng: LatLngExpression) {
  polyline.value = [...polyline.value, latlng]
}
function clearPolyline() {
  polyline.value = []
}

onKeyStroke('Escape', () => {
  console.log('Escape pressed')
  if (polygonMode.value) {
    clearPolyline()
    polygonMode.value = false
  }
})

defineExpose({ fitBounds, fitViewToSite, el: mapContainer })
</script>

<style lang="scss">
@use 'vuetify';
.leaflet-container {
  background-color: rgb(var(--v-theme-surface));
}

.coordinates-control {
  pointer-events: none;
  * {
    pointer-events: none;
  }
  .v-card {
    background-color: rgb(var(--v-theme-surface), 0.5);
    code * {
      opacity: 1;
    }
  }
}

.polygon-mode .map {
  cursor: crosshair;
}

.hexbin-hexagon {
  stroke: white;
  stroke-opacity: 0.5;
  stroke-width: 1;
  cursor: pointer;
  &.hover {
    stroke-width: 2;
    stroke: orangered;
  }
}

.marker-icon {
  // margin-left: -1.5em !important;
  // width: 25px !important;
  // height: auto !important;
  // z-index: 509 !important;
  // position: absolute;
  // left: 0;
  // top: 0;
  font-size: 3em;
}
</style>
