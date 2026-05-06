<template>
  <div
    ref="mapContainer"
    id="map-container"
    class="deck-map-container fill-height"
    @mouseleave="cursorCoordinates = undefined"
  >
    <div ref="mapHost" class="deck-map"></div>

    <div class="map-control top-left">
      <div class="d-flex flex-column ga-1">
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
        <layers-control
          :hexgrid
          :marker-layers
          :hasSiteMarkers="!!markers?.length"
          v-model:site-markers-visible="siteMarkersVisible"
          v-model:regions="regions"
          v-model:roads="roads"
          @toggleHexgrid="(v) => emit('toggleHexgrid', v)"
          @toggleMarkers="(index, v) => emit('toggleMarkers', index, v)"
        />
      </div>
    </div>

    <div v-if="closable" class="map-control top-right">
      <v-btn
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
    </div>

    <div class="d-flex top-right ga-2 map-control">
      <div v-if="cursorCoordinates" class="pointer-events-none">
        <v-card density="compact" class="pa-2 opacity-70" theme="light" rounded="0">
          <code class="text-label-small font-monospace d-block">
            <div class="d-flex justify-space-between ga-2">
              <span>Lat:</span> {{ cursorCoordinates.lat.toFixed(4) }}
            </div>
            <div class="d-flex justify-space-between ga-2">
              <span>Lng:</span> {{ cursorCoordinates.lng.toFixed(4) }}
            </div>
            <div class="d-flex justify-space-between ga-2">
              <span>Zoom:</span> {{ currentZoom.toFixed(2) }}
            </div>
          </code>
        </v-card>
      </div>

      <div v-if="hexgrid?.active && hexgridColorDomain">
        <color-scale-widget
          :min="hexgridColorDomain.min"
          :max="hexgridColorDomain.max"
          :color-range="hexgridColorRange"
          :binding-label="bindingLabels[hexgrid.colorBinding.binding]"
          :log="hexgrid.colorBinding.log"
          :hidden="false"
        />
      </div>
    </div>

    <div
      v-if="selected?.type === 'item' && selected.info.object"
      class="map-popup site-popup bottom-left"
    >
      <slot name="popup" :item="selected.info.object" :zoom="currentZoom" />
    </div>
    <div v-else-if="selected?.type === 'site'" class="map-popup site-popup bottom-left">
      <slot name="popup" :item="selected.info" :zoom="currentZoom" />
    </div>
    <div
      v-else-if="selected?.type === 'cluster' && selected?.info.object"
      class="map-popup hex-popup"
    >
      <slot name="cluster-popup" :data="selected.info.object.items" :type="selected.type" />
    </div>
    <div
      v-else-if="
        selected?.type === 'hexagon' &&
        selected?.info.object &&
        selected.info.object.points?.length === 1
      "
      class="map-popup site-popup bottom-left"
    >
      <slot name="popup" :item="selected.info.object.points[0]" :zoom="currentZoom" />
    </div>

    <div
      v-else-if="selected?.type === 'hexagon' && selected?.info.object"
      class="map-popup hex-popup"
    >
      <slot name="cluster-popup" :data="selected.info.object.points" :type="selected.type" />
    </div>

    <div
      v-if="hoverTooltip"
      class="map-popup hex-hover-tooltip pointer-events-none"
      :style="{ left: `${hoverTooltip.x}px`, top: `${hoverTooltip.y}px` }"
    >
      <div class="hex-hover-tooltip__content text-label-small font-monospace">
        {{ hoverTooltip.text }}
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'DeckGlMap'
}
</script>

<script setup lang="ts" generic="Marker extends SiteWithOccurrences & { color?: string }">
import { HexagonLayer } from '@deck.gl/aggregation-layers'
import type { Layer } from '@deck.gl/core'
import { MapboxOverlay } from '@deck.gl/mapbox'
import { onKeyStroke, useDebounceFn, useFullscreen, useThrottleFn } from '@vueuse/core'
import maplibregl, { LngLatBounds, type StyleSpecification } from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'
import {
  computed,
  nextTick,
  onMounted,
  onUnmounted,
  ref,
  shallowRef,
  watch,
  watchEffect
} from 'vue'

import { CoordinatesPrecision, type SiteWithOccurrences } from '@/api'
import type { Geocoordinates } from '@/features/cartography/coordinates'
import { paletteRGB, RGBToArray } from '@/lib/color_brewer'
import { bindingLabels, getBindingFn, hexagonLayerColorBinding } from '../bindings'
import { useMarkerLayers } from '../composables/useMarkerLayers'
import type { HexgridLayerDeck, MarkerLayer } from './layers-manager/map-layers'

import '@deck.gl/widgets/stylesheet.css'
import { circle } from '@turf/turf'
import { useMarkerSelection } from '../composables/marker-selection'
import {
  createRegionLayers,
  createRoadsLayers,
  regionsLayerStyle,
  updateRegionsLayerVisibility,
  updateRoadsLayerVisibility
} from '../regions-overlay'
import ColorScaleWidget from './controls/ColorScaleWidget.vue'
import LayersControl from './controls/LayersControl.vue'
import { useHexgridLayer } from '../composables/hexgrid-layer'

export type MarkerOptions = {
  cluster: {
    radiusScaleFactor: number
    labelZoomThreshold: number
  }
  tooltips?: boolean
}
type Item = SiteWithOccurrences
type MarkerWithColor = Marker

const mapContainer = ref<HTMLElement>()
const mapHost = ref<HTMLElement>()
const map = shallowRef<maplibregl.Map>()
const overlay = shallowRef<MapboxOverlay>()
const mapInitialized = ref(false)

const regions = defineModel<boolean>('regions', { default: true })
const roads = defineModel<boolean>('roads', { default: false })
watch(regions, (value) => updateRegionsLayerVisibility(map, value), { immediate: true })
watch(roads, (value) => updateRoadsLayerVisibility(map, value), { immediate: true })

const {
  hexgrid,
  minZoom = 1,
  maxZoom = 18,
  autoFit = true,
  center = [0, 0],
  zoom = 2,
  ...props
} = defineProps<{
  hexgrid?: HexgridLayerDeck<Item>
  markers?: MarkerWithColor[]
  markerLayers?: MarkerLayer<Item>[]
  autoFit?: boolean | number
  closable?: boolean
  regions?: boolean
  center?: [number, number]
  minZoom?: number
  maxZoom?: number
  zoom?: number
  markerOptions?: MarkerOptions
}>()

defineSlots<{
  popup: (props: { item: Item; zoom: number }) => any
  'cluster-popup': (props: { data?: Item[]; type: 'cluster' | 'hexagon' }) => any
}>()

const emit = defineEmits<{
  toggleHexgrid: [active: boolean]
  toggleMarkers: [layerIndex: number, active: boolean]
  close: []
}>()

const cursorCoordinates = ref<{ lat: number; lng: number }>()
const currentZoom = ref(zoom)
const hoverTooltip = ref<{ x: number; y: number; text: string }>()

const { isFullscreen, exit, toggle } = useFullscreen(mapContainer)
onKeyStroke('Escape', () => exit())
const toggleFullscreen = useThrottleFn(toggle, 150)

const mapStyle: StyleSpecification = {
  version: 8,
  glyphs: regionsLayerStyle.glyphs,
  sources: {
    arcgisImagery: {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
        'https://services.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'
      ],
      tileSize: 256,
      attribution:
        "Powered by <a href='https://www.esri.com/'>Esri</a> — Source: Esri, Maxar, Earthstar Geographics, and the GIS User Community"
    },
    openmaptiles: {
      type: 'vector',
      url: regionsLayerStyle.sources.openmaptiles.url
    }
  },
  layers: [
    {
      id: 'arcgis-imagery',
      type: 'raster',
      source: 'arcgisImagery'
    },
    ...createRegionLayers(regions.value),
    ...createRoadsLayers(false)
  ]
}

const { selected, select, highlightLayer, clear } = useMarkerSelection<Item>()

function displayRadius(site: Item) {
  const precision = site.coordinates.precision
  if (precision === 'Unknown') return
  const radius = CoordinatesPrecision.radius(precision)
  const c = circle([site.coordinates.longitude, site.coordinates.latitude], radius, {
    steps: 64,
    units: 'meters'
  })
  map.value?.addSource(`site-radius`, {
    type: 'geojson',
    data: c
  })
  map.value?.addLayer({
    id: `site-radius`,
    type: 'fill',
    source: `site-radius`,
    paint: {
      'fill-color': 'rgba(0, 123, 255, 0.2)'
    }
  })
  map.value?.addLayer({
    id: `site-radius-outline`,
    type: 'line',
    source: `site-radius`,
    paint: {
      'line-color': '#fb8c00',
      'line-width': 2,
      'line-dasharray': [2, 2]
    }
  })
}

function removeRadius(mapInstance: maplibregl.Map) {
  mapInstance.getLayer(`site-radius`) && mapInstance.removeLayer(`site-radius`)
  mapInstance.getLayer(`site-radius-outline`) && mapInstance.removeLayer(`site-radius-outline`)
  mapInstance.getSource(`site-radius`) && mapInstance.removeSource(`site-radius`)
}

watch(selected, () => {
  if (!map.value) return
  removeRadius(map.value)
  if (selected.value?.type === 'item') displayRadius(selected.value.info.object!)
})

const { markerDeckLayers } = useMarkerLayers(props, {
  currentZoom,
  hoverTooltip,
  selected,
  select
})

const { hexgridLayer, hexgridColorDomain, hexgridColorRange } = useHexgridLayer<Item>(
  { hexgrid: () => hexgrid },
  {
    selected,
    select,
    currentZoom,
    hoverTooltip
  }
)

const deckLayers = computed<Layer[]>(() => {
  console.debug('Recomputing deck layers')
  const layers = [hexgridLayer.value, ...markerDeckLayers.value, ...highlightLayer.value].filter(
    (layer) => layer !== undefined
  )
  return layers satisfies Layer[]
})

const fitSignature = computed(() => {
  const hexSignature = hexgrid
    ? [
        hexgrid.data?.length ?? 0,
        hexgrid.data?.[0]?.coordinates.latitude ?? '',
        hexgrid.data?.[0]?.coordinates.longitude ?? ''
      ].join(':')
    : 'none'

  const markerSignature = (props.markerLayers ?? [])
    .map((layer) =>
      [
        layer.data?.length ?? 0,
        layer.data?.[0]?.coordinates.latitude ?? '',
        layer.data?.[0]?.coordinates.longitude ?? ''
      ].join(':')
    )
    .join('|')

  const markerSignatureSingle = (props.markers ?? []).length

  return [hexSignature, markerSignature, markerSignatureSingle].join('::')
})

function computeAllPoints() {
  const hexPoints = hexgrid?.active ? (hexgrid.data ?? []) : []
  const markerLayerPoints = (props.markerLayers ?? [])
    .filter((layer) => layer.active)
    .flatMap((layer) => layer.data ?? [])
  const singleMarkers = props.markers ?? []

  return [...hexPoints, ...markerLayerPoints, ...singleMarkers].map(({ coordinates }) => [
    coordinates.longitude,
    coordinates.latitude
  ]) as [number, number][]
}

function fitMapView(radiusMeters = 0) {
  if (autoFit === false || !map.value || !mapInitialized.value) return

  const points = computeAllPoints()
  if (!points.length) return

  const bounds = points.reduce((acc, [lng, lat]) => acc.extend([lng, lat]), new LngLatBounds())

  map.value.fitBounds(bounds, {
    padding: Math.max(40, radiusMeters ? Number(radiusMeters) / 200 : 40),
    duration: 350,
    maxZoom: maxZoom
  })
}

function fitViewToSite({ coordinates: { latitude, longitude } }: Geocoordinates) {
  map.value?.easeTo({
    center: [longitude, latitude],
    zoom: Math.max(10, currentZoom.value),
    duration: 300
  })
}

watch(deckLayers, () => scheduleOverlayUpdate(), { flush: 'post' })
watch(
  fitSignature,
  () => {
    if (!mapInitialized.value) return
    fitMapView(typeof autoFit === 'number' ? autoFit : 0)
  },
  { flush: 'post' }
)

const mapMarkers = ref<maplibregl.Marker[]>()
const siteMarkersVisible = ref<boolean>(true)
watchEffect(() => {
  if (!map.value) return
  mapMarkers.value?.forEach((marker) => marker.remove())
  if (!siteMarkersVisible.value) return
  mapMarkers.value = props.markers?.map((marker) => {
    const m = new maplibregl.Marker({ color: marker.color })
      .setLngLat([marker.coordinates.longitude, marker.coordinates.latitude])
      .addTo(map.value!)

    m.on('click', () => {
      select({ type: 'site', info: marker })
    })
    return m
  })
})

function updateOverlayLayers() {
  overlay.value?.setProps({ layers: deckLayers.value })
}

const scheduleOverlayUpdate = useDebounceFn(async () => {
  if (!mapInitialized.value) return
  await nextTick()
  updateOverlayLayers()
}, 16)

onMounted(() => {
  void nextTick(() => {
    if (!mapHost.value || !mapContainer.value) return
    if (mapHost.value.clientWidth === 0 || mapHost.value.clientHeight === 0) return

    const [lat, lng] = center
    const mapInstance = new maplibregl.Map({
      container: mapHost.value,
      style: mapStyle,
      center: [lng, lat],
      zoom: zoom,
      minZoom: minZoom,
      maxZoom: maxZoom,
      attributionControl: {
        compact: true
      }
    })

    map.value = mapInstance

    // disable map rotation using right click + drag
    mapInstance.dragRotate.disable()
    // disable map rotation using keyboard
    mapInstance.keyboard.disable()
    // disable map rotation using touch rotation gesture
    mapInstance.touchZoomRotate.disableRotation()

    // track cursor position and zoom level
    mapInstance.on('mousemove', (event) => {
      cursorCoordinates.value = {
        lat: event.lngLat.lat,
        lng: event.lngLat.lng
      }
    })
    mapInstance.on('move', () => {
      currentZoom.value = mapInstance.getZoom()
    })

    mapInstance.on('click', () => {
      clear()
    })

    mapInstance.on('load', () => {
      const deckOverlay = new MapboxOverlay({
        getCursor({ isDragging, isHovering }) {
          if (isDragging) return 'grabbing'
          if (isHovering) return 'pointer'
          return 'grab'
        },
        interleaved: true,
        layers: deckLayers.value
      })

      mapInstance.addControl(deckOverlay)
      overlay.value = deckOverlay
      mapInitialized.value = true
      updateOverlayLayers()
      fitMapView(typeof autoFit === 'number' ? autoFit : 0)
    })
  })
})

onUnmounted(() => {
  overlay.value?.finalize()
  map.value?.remove()
  overlay.value = undefined
  map.value = undefined
  mapInitialized.value = false
})

defineExpose({ fitMapView, fitViewToSite, select, clear, el: mapContainer })
</script>

<style scoped lang="scss">
.deck-map-container {
  position: relative;
  width: 100%;
  height: 100%;
  // overflow: hidden;
}

.deck-map {
  width: 100%;
  height: 100%;
}

.map-control {
  position: absolute;
  z-index: 10;
}

.top-left {
  top: 12px;
  left: 12px;
}

.top-right {
  top: 12px;
  right: 12px;
}

.bottom-right {
  bottom: 12px;
  right: 12px;
}
.bottom-left {
  bottom: 12px;
  left: 12px;
}

.pointer-events-none {
  pointer-events: none;
}

.map-popup {
  position: absolute;
  z-index: 11;
  max-width: min(480px, calc(100% - 24px));
}

.hex-popup {
  left: 12px;
  bottom: 12px;
}

.hex-hover-tooltip {
  transform: translate(12px, 12px);
  max-width: min(240px, calc(100% - 24px));
}

.hex-hover-tooltip__content {
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(0, 0, 0, 0.12);
  box-shadow: 0 4px 18px rgba(0, 0, 0, 0.16);
  padding: 8px 12px;
}
</style>
