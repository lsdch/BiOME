<template>
  {{ polyline }}
  <div
    ref="map"
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
            v-if="items || marker"
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

      <LControl
        position="topleft"
        :options="{}"
        v-if="markerMode === 'hexgrid' && hexgridColorRange?.length"
      >
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
      <slot name="default" :map :zoom></slot>

      <!-- Hexagon layer  -->
      <LHexbinLayer
        v-if="hexgridConfig.active && items"
        pane="overlayPane"
        :data="items"
        :accessor="(item) => Geocoordinates.LatLng(item)"
        :radius="hexgridConfig.radius"
        :radius-range="
          hexgridConfig.bindings?.radius
            ? hexgridConfig.radiusRange
            : [hexgridConfig.radius, hexgridConfig.radius]
        "
        :opacity="
          hexgridConfig.bindings?.opacity ? hexgridConfig.opacityRange : hexgridConfig.opacity
        "
        :hover-fill="hexgridConfig.hover.fill"
        :hover-scale="hexgridConfig.hover.useScale ? hexgridConfig.hover.scale : undefined"
        :color-range="hexgridConfig.colorRange"
        style="cursor: pointer"
        :color-binding="hexgridConfig.bindings?.color"
        :opacity-binding="hexgridConfig.bindings?.opacity"
        :radius-binding="hexgridConfig.bindings?.radius"
        @update:color-scale-extent="(range) => console.log(range)"
        @click="
          (e) => {
            if (e.length === 1) selectSite(e[0].data)
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
      <LLayerGroup v-if="markerConfig.active && items" pane="markerPane">
        <LMarkerClusterGroup
          v-if="markerConfig.clustered"
          remove-outside-visible-bounds
          show-coverage-on-hover
          :maxClusterRadius="70"
        >
          <LCircleMarker
            v-for="item in items"
            pane="markerPane"
            :key="item.id"
            :lat-lng="[item.coordinates.latitude, item.coordinates.longitude]"
            v-bind="markerConfig"
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
          v-for="item in items"
          :key="item.id"
          :latLng="[item.coordinates.latitude, item.coordinates.longitude]"
          v-bind="markerConfig"
          :opacity="1"
          :fill-opacity="1"
          :options="{
            zIndexOffset: 10
          }"
          @click="selectSite(item)"
        >
        </LCircleMarker>
      </LLayerGroup>

      <slot
        v-if="marker"
        name="marker"
        :lat-lng="marker ? Geocoordinates.LatLng(marker) : undefined"
      >
        <LMarker :lat-lng="Geocoordinates.LatLng(marker)" />
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
          (ev) => {
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

<script setup lang="ts" generic="SiteItem extends { id: string } & Geocoordinates">
import 'leaflet/dist/leaflet.css'
import LHexbinLayer, { type ScaleBinding } from 'vue-leaflet-hexbin'
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
import {
  onKeyPressed,
  onKeyStroke,
  useDebounceFn,
  useFullscreen,
  useThrottleFn
} from '@vueuse/core'
import L, {
  CircleMarkerOptions,
  latLng,
  latLngBounds,
  LatLngExpression,
  LatLngLiteral,
  PointExpression,
  polygon,
  type LeafletMouseEvent,
  type Map
} from 'leaflet'

import { nextTick, reactive, ref, UnwrapRef, useTemplateRef, watch } from 'vue'
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster'
import { Geocoordinates } from '.'

import MapColorLegend from '@/views/location/MapColorLegend.vue'
import { vElementVisibility } from '@vueuse/components'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { MapLayerMode } from './MarkerControl.vue'
import { useKeyPress } from '@vue-flow/core'

export type HexPopupData<SiteItem> = {
  data: SiteItem
  coord: L.LatLngExpression
}

export type HexgridScaleBindings<SiteItem> = {
  color?: ScaleBinding<SiteItem>
  radius?: ScaleBinding<SiteItem>
  opacity?: ScaleBinding<SiteItem>
}

// Opacity can be controlled directly by 'color' and 'fill' properties
export type MarkerConfig = {
  active: boolean
  clustered: boolean
} & Overwrite<
  Omit<CircleMarkerOptions, 'opacity' | 'fillOpacity' | 'renderer'>,
  { dashArray?: string | undefined }
>

export type HexgridConfig<SiteItem> = {
  active: boolean
  radius: number
  radiusRange?: [number, number]
  colorRange?: string[]
  hover: {
    fill: boolean
    useScale: boolean
    scale: number
  }
  opacity: number
  opacityRange?: [number, number]
  bindings: HexgridScaleBindings<SiteItem>
}

const markerMode = defineModel<MapLayerMode>('marker-mode', { default: 'markers' })

const markerConfig = defineModel<MarkerConfig>('marker-config', {
  default: () =>
    reactive({
      active: false,
      radius: 8,
      clustered: false
    })
})

const hexgridConfig = defineModel<HexgridConfig<SiteItem>>('hexgrid-config', {
  default: () => ({
    active: true,
    radius: 10,
    radiusRange: [10, 10],
    opacity: 0.8,
    asRange: false,
    hover: {
      fill: true,
      scale: 1,
      useScale: false
    },
    bindings: {}
  })
})

const hexgridColorRange = ref<[number, number]>()

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
    marker?: Geocoordinates
    // markerOptions?: Omit<CircleMarkerOptions, 'dashArray'>
    bounds?: [LatLngExpression, LatLngExpression]
    autoFit?: boolean | number
    closable?: boolean
    regions?: boolean
    center?: PointExpression
    minZoom?: number
    maxZoom?: number
    hideMarkerControl?: boolean
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
  default: (props: { zoom: number; map?: HTMLElement }) => any
  popup: (props: { item: SiteItem; popupOpen: boolean; zoom: number }) => any
  marker: (props: { latLng?: LatLngExpression }) => any
  'hex-popup': (props: { data?: HexPopupData<UnwrapRef<SiteItem>>[] }) => any
}>()

const mapBounds = ref(L.latLngBounds(...props.bounds))

watch(
  () => [props.items, props.marker, props.autoFit],
  () => fitMapView()
)

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
      fitRadius(props.autoFit)
    } else {
      fitBounds(props.items)
    }
  }
}

const fitRadius = useDebounceFn((radius: number) => {
  if (props.marker) {
    let r = Math.max(100, radius)
    const { latitude, longitude } = props.marker.coordinates
    mapBounds.value = L.latLng(latitude, longitude)
      .toBounds(r + 100)
      .pad(0.5)
  }
}, 200)

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
    mapBounds.value = latLngBounds(minMaxCoords.sw, minMaxCoords.ne).pad(0.1)
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

defineExpose({ fitBounds })
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
</style>
