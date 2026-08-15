<template>
  <div class="fill-height w-100 d-flex">
    <v-navigation-drawer
      :location="$vuetify.display.xs ? 'top' : 'left'"
      :width="500"
      v-model="drawer"
      :temporary="!drawerPinned"
    >
      <div class="fill-height d-flex flex-column">
        <v-tabs v-model="tab" class="flex-shrink-0">
          <v-tab value="layers" prepend-icon="mdi-layers"> Layers </v-tab>
          <!-- <v-tab value="sites" prepend-icon="mdi-map-marker">
            Sites
            <v-badge
              v-if="layerSpecs.sites?.length"
              :content="layerSpecs.sites.length"
              color="primary"
              inline
              class="ml-1"
            />
          </v-tab> -->
          <v-tab value="config">
            <v-icon icon="mdi-cog" />
          </v-tab>
          <v-spacer />

          <v-btn
            variant="plain"
            icon="mdi-chevron-left"
            rounded="xl"
            color=""
            @click="toggleDrawer(false)"
          />
        </v-tabs>
        <v-divider />
        <v-tabs-window v-model="tab" class="flex-grow-1 bg-main overflow-y-auto">
          <v-tabs-window-item eager value="layers">
            <LayersManager
              v-model:hex-layer="layerSpecs.hexgrid"
              v-model:marker-layers="layerSpecs.markers"
              :zoom
              :global-opts="markerOptions"
            />
          </v-tabs-window-item>
          <!-- <v-tabs-window-item value="sites">
            <SiteSearchPanel
              v-model="layerSpecs.sites"
              @focus-site="
                (site) => {
                  const info = siteMarkers?.find((s) => s.data.code === site.code)
                  if (info) map?.select({ type: 'site', info })
                  map?.fitViewToSite(site)
                }
              "
            />
          </v-tabs-window-item> -->
          <v-tabs-window-item value="config">
            <MapViewConfig v-model:save-layers="saveLayers" v-model:markerOptions="markerOptions" />
          </v-tabs-window-item>
        </v-tabs-window>
      </div>

      <!-- DRAWER FOOTER-->
      <template #append>
        <v-divider />
        <div class="d-flex justify-space-between pa-2">
          <div>
            <!-- <MapPresetSaveDialog
              v-if="userStore.isAuthenticated"
              :specs="{
                hexgrid: layerSpecs.hexgrid,
                markers: layerSpecs.markers
              }"
            >
              <template #activator="{ props }">
                <v-btn
                  variant="text"
                  icon="mdi-content-save"
                  v-tooltip="'Save map preset'"
                  v-bind="props"
                ></v-btn>
              </template>
            </MapPresetSaveDialog>
            <MapPresetLoadDialog
              @apply="
                ({ spec: { hexgrid, markers }, name }) => {
                  layerSpecs.hexgrid = hexgrid
                  layerSpecs.markers.splice(0, markers.length, ...markers)
                  feedback({ message: `Loaded preset '${name}'`, type: 'success' })
                }
              "
            >
              <template #activator="{ props }">
                <v-btn
                  variant="text"
                  icon="mdi-file-star"
                  v-tooltip="'Load preset'"
                  v-bind="props"
                />
              </template>
            </MapPresetLoadDialog> -->
            <v-btn
              icon="mdi-share"
              variant="text"
              v-tooltip="'Share current map configuration'"
              @click="share"
            />
          </div>
          <v-btn
            icon="mdi-pin"
            size="small"
            :variant="drawerPinned ? 'tonal' : 'plain'"
            @click="drawerPinned = !drawerPinned"
            v-tooltip="'Toggle permanent drawer'"
          />
        </div>
      </template>
    </v-navigation-drawer>

    <v-navigation-drawer v-if="!drawerPinned || !drawer" rail location="left" class="bg-main">
      <v-list>
        <v-tooltip content-class="bg-surface text-overline py-0" :height="48">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-layers"
              @click="toggleTab('layers')"
              :active="tab === 'layers' && drawer"
              color="primary"
            />
          </template>
          <v-sheet :height="48" class="my-0 d-flex align-center"> Layers </v-sheet>
        </v-tooltip>
        <!-- <v-tooltip content-class="bg-surface text-overline py-0" :height="48">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-map-marker"
              title="Sites"
              @click="toggleTab('sites')"
              :active="tab === 'sites' && drawer"
              color="primary"
            />
          </template>
          <v-sheet :height="48" class="my-0 d-flex align-center"> Sites </v-sheet>
        </v-tooltip> -->
      </v-list>
    </v-navigation-drawer>

    <div class="fill-height w-100 d-flex flex-column">
      <v-progress-linear v-if="allPending" indeterminate color="warning" />
      <div :class="['fill-height w-100 position-relative']">
        <DeckGlMap
          ref="map"
          :auto-fit="!anyLoading"
          :hexgrid="hexgridLayer"
          @toggle-hexgrid="(v) => (layerSpecs.hexgrid.active = v)"
          @toggle-markers="(i, v) => (layerSpecs.markers[i].active = v)"
          :marker-layers
          :marker-options
          v-model:zoom="zoom"
        >
          <!-- :pinMarkers="siteMarkers" -->
          <!-- <template #cluster-popup="{ data, resolution, params }"> </template>
          <template #pin-popup="{ item }">
            <KeepAlive>
              <SitePopupWithOccurrences :item="item.data" />
            </KeepAlive>
          </template> -->
          <template #popup="{ selection }">
            <SiteClusterPopup
              v-if="selection.type === 'hexagon' && !!selection.info.object"
              :data="selection.info.object"
              :resolution="selection.resolution"
              :params="selection.params"
              :attach="map?.el"
            />

            <SiteClusterPopup
              v-else-if="selection.type === 'marker' && !!selection.info.object"
              :data="selection.info.object"
              :resolution="selection.resolution"
              :params="selection.params"
              :attach="map?.el"
            />
          </template>
        </DeckGlMap>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import DeckGlMap, { GlobalMarkerOptions } from '@/features/cartography/components/DeckGlMap.vue'

import {
  HexgridLayer,
  HexLayerSpec,
  makeHexLayer,
  MarkerLayer,
  markerLayerFromSpec,
  MarkerLayerSpec,
  PinMarker
} from '@/features/cartography/components/layers-manager/map-layers'
// import MapPresetLoadDialog from '@/features/cartography/components/map-presets/MapPresetLoadDialog.vue'
// import MapPresetSaveDialog from '@/features/cartography/components/map-presets/MapPresetSaveDialog.vue'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import {
  computedAsync,
  useClipboard,
  useLocalStorage,
  useSessionStorage,
  useToggle
} from '@vueuse/core'
import { compressToEncodedURIComponent, decompressFromEncodedURIComponent } from 'lz-string'
import { onMounted, ref, useTemplateRef, watch } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { useRoute } from 'vue-router'
import { CellMarkerData, useLayerData } from '../components/layers-manager/layer-data'
import LayersManager from '../components/layers-manager/LayersManager.vue'
import MapViewConfig from '../components/MapViewConfig.vue'
import SiteClusterPopup from '../components/popups/MultiSamplingsPopup.vue'
import { H3CellWithRichness } from '@/api/adapters.ts'
import { cellToLatLng } from 'h3-js'
import { hexgridLayerFromSpec } from '../composables/hexgrid-layer'

const route = useRoute()

const zoom = ref(2)

type LayerSpecs = {
  hexgrid: HexLayerSpec
  markers: MarkerLayerSpec[]
  sites: PinMarker<CellMarkerData>[]
}

const layerSpecs = useSessionStorage<LayerSpecs>(
  'mapping-tool-layer-specs',
  {
    hexgrid: makeHexLayer(),
    markers: [],
    sites: []
  },
  { mergeDefaults: true, deep: true }
)

const saveLayers = useLocalStorage('mapping-tool-save-layers', false, {})
const markerOptions = useLocalStorage<GlobalMarkerOptions>('mapping-tool-marker-options', {
  cluster: {
    radiusScaleFactor: 0.5,
    labelZoomThreshold: 8
  },
  tooltips: true
})

watch(
  () => [layerSpecs.value, saveLayers.value],
  ([newSpecs, saveLayersValue]) => {
    if (saveLayersValue) {
      localStorage.setItem('mapping-tool-layer-specs', JSON.stringify(newSpecs))
    }
  },
  { deep: true }
)

// const singleSites = useQuery(
//   computed(() => {
//     return {
//       enabled: layerSpecs.value.sites.length > 0,
//       ...occurrencesBySiteOptions({
//         body: {
//           site_codes: layerSpecs.value.sites.map((s) => s.data.code),
//           sampling_target: {}
//         }
//       })
//     }
//   })
// )
// const siteMarkers = computed(() => {
//   if (!layerSpecs.value.sites?.length) return []
//   const colorMap = new Map(layerSpecs.value.sites.map((s) => [s.data.code, s.options?.color]))
//   return singleSites.data.value?.map<PinMarker<SiteWithOccurrences>>((site) => ({
//     data: site,
//     coordinates: site.coordinates,
//     options: { color: colorMap.get(site.code) }
//   }))
// })

const drawerPinned = useLocalStorage('mapping-tool-drawer-pinned', false, {
  initOnMounted: true
})

type DeckGlMapExposed = ComponentExposed<typeof DeckGlMap>
const map = useTemplateRef<DeckGlMapExposed>('map')
const userStore = useUserStore()

const [drawer, toggleDrawer] = useToggle(false)

type MappingToolTab = 'layers' | 'config'

const tab = ref<MappingToolTab>('layers')

function toggleTab(newTab: MappingToolTab) {
  tab.value = newTab
  toggleDrawer(true)
}

const { feedback } = useFeedback()

const { allPending, anyLoading, layerData, data, H3CellToMarkerData } = useLayerData()

const hexgridLayer = computedAsync<HexgridLayer<H3CellWithRichness>>(async () => {
  console.debug('Recomputing hexgrid layer with spec', layerSpecs.value.hexgrid)
  const remote = data.get(layerSpecs.value.hexgrid.id)
  if (remote?.isFetching?.value) await remote?.suspense()
  return hexgridLayerFromSpec<H3CellWithRichness>(
    layerSpecs.value.hexgrid,
    remote?.data.value
    // applySiteFilter(remote?.data.value, layerSpecs.value.hexgrid.include_sites)
  )
})

const markerLayers = computedAsync<MarkerLayer<H3CellWithRichness>[]>(async () => {
  console.log('Recomputing marker layers')
  return Promise.all(
    layerSpecs.value.markers.map(async (layer) => {
      const remote = data.get(layer.id)
      if (remote?.isFetching?.value) await remote?.suspense()
      return markerLayerFromSpec(layer, remote?.data.value ?? [], {
        radius(item) {
          return item.occurrences_count
        },
        getPosition: (item) => {
          const [lat, lng] = cellToLatLng(item.h3_index)
          return [lng, lat]
        }
      })
      // return markerLayerFromSpec(layer, applySiteFilter(remote?.data.value, layer.include_sites))
    })
  )
})

const { copy } = useClipboard()

function encodeLayerSpecs() {
  return compressToEncodedURIComponent(JSON.stringify(layerSpecs.value))
}

function share() {
  const h = encodeLayerSpecs()
  const url = `${window.location.origin}${window.location.pathname}?h=${h}`
  copy(url)
  feedback({ message: 'Map configuration URL copied to clipboard', type: 'success' })
}

function getSpecQueryValue() {
  const value = route.query.h
  return Array.isArray(value) ? value[0] : value
}

onMounted(() => {
  const encodedSpecs = getSpecQueryValue()
  if (encodedSpecs) {
    try {
      const decoded = decompressFromEncodedURIComponent(encodedSpecs)
      layerSpecs.value = JSON.parse(decoded)
    } catch (e) {
      console.error('Failed to parse layer specs from URL', e)
    }
  } else if (saveLayers.value) {
    const savedSpecs = localStorage.getItem('mapping-tool-layer-specs')
    if (savedSpecs) {
      try {
        layerSpecs.value = JSON.parse(savedSpecs)
      } catch (e) {
        console.error('Failed to parse saved layer specs from localStorage', e)
      }
    }
  }
})
</script>

<style lang="scss">
@use 'vuetify';

.map-toolbar {
  height: 100%;
  width: 300px;
  background-color: rgb(var(--v-theme-surface));
}
</style>
