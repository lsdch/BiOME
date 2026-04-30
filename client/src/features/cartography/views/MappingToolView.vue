<template>
  <div class="fill-height w-100 d-flex">
    <v-navigation-drawer
      :location="$vuetify.display.xs ? 'top' : 'left'"
      :width="500"
      v-model="drawer"
      :temporary="!drawerPinned"
    >
      <div class="fill-height d-flex flex-column">
        <v-tabs v-model="tab">
          <v-tab value="layers" prepend-icon="mdi-layers"> Layers </v-tab>
          <v-tab value="sites" prepend-icon="mdi-map-marker"> Sites </v-tab>
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
              v-model:hex-layer="hexLayerSpec"
              v-model:marker-layers="markerLayerSpecs"
            />
          </v-tabs-window-item>
          <v-tabs-window-item value="sites">
            <SiteSearchPanel v-model="siteMarkers" @focus-site="map?.fitViewToSite" />
          </v-tabs-window-item>
          <v-tabs-window-item value="config">
            <v-list>
              <ListItemInput
                label="Save map configuration"
                subtitle="Restore last map configuration on next visit"
              >
                <v-switch color="primary" hide-details />
              </ListItemInput>
              <CardDialog v-if="userStore.isGranted('Contributor')" title="Map presets">
                <template #append>
                  <v-switch
                    v-if="userStore.isGranted('Maintainer')"
                    v-model="showAllPresets"
                    label="Maintainer view"
                    hide-details
                    color="warning"
                    v-tooltip="'Display all registered presets'"
                  />
                </template>
                <template #activator="{ props }">
                  <v-list-item
                    title="Manage presets"
                    prepend-icon="mdi-folder-star-multiple"
                    v-bind="props"
                  />
                </template>
                <MapPresetManager :all="showAllPresets" />
              </CardDialog>
            </v-list>
          </v-tabs-window-item>
        </v-tabs-window>
      </div>
      <template #append>
        <v-divider />
        <div class="d-flex justify-space-between pa-2">
          <div>
            <MapPresetSaveDialog
              v-if="userStore.isAuthenticated"
              :specs="{
                hexgrid: hexLayerSpec,
                markers: markerLayerSpecs
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
                  hexLayerSpec = hexgrid
                  markerLayerSpecs.splice(0, markers.length, ...markers)
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
            </MapPresetLoadDialog>
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
        <v-tooltip content-class="bg-surface text-overline py-0" :height="48">
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
        </v-tooltip>
      </v-list>
    </v-navigation-drawer>
    <div class="fill-height w-100 d-flex flex-column">
      <v-progress-linear v-if="allPending" indeterminate color="warning" />
      <div :class="['fill-height w-100 position-relative']">
        <BaseMap
          ref="map"
          auto-fit
          :markers="siteMarkers"
          :hexgrid="hexgridLayer"
          :bounds="mapBounds"
          :marker-layers
        >
          <!-- <LControl v-if="isRefetching || isFetching" position="topleft">
            <v-progress-circular
              v-if="isPending || isRefetching"
              indeterminate
              color="warning"
              size="32"
              width="6"
            />
          </LControl> -->
          <!-- <LControl position="topright" v-if="sites">
            <MapStatsDialog :sites>
              <template #activator="{ props }">
                <v-btn v-bind="props" icon="mdi-poll" color="white" :width="45" :height="45" />
              </template>
            </MapStatsDialog>
          </LControl> -->
          <!-- <LControl position="topright" v-if="sites">
            <v-btn icon="mdi-shape-polygon-plus" @click="togglePolygonMode(true)"></v-btn>
          </LControl> -->
          <template #hex-popup="{ data }">
            <MapViewHexPopup :data />
          </template>
          <template #popup="{ item, popupOpen, zoom }">
            <KeepAlive>
              <MapViewSitePopup :item="item" :popupOpen="popupOpen" :zoom="zoom" :key="item.code" />
            </KeepAlive>
          </template>
        </BaseMap>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseMap from '@/features/cartography/components/BaseMap.vue'

import { SiteItem, SiteWithOccurrences } from '@/api'
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import { useScaleBinding } from '@/composables/occurrences'
import {
  HexgridLayer,
  HexLayerSpec,
  makeHexLayer,
  MarkerLayer,
  MarkerLayerSpec
} from '@/features/cartography/components/layers-manager/map-layers'
import MapPresetLoadDialog from '@/features/cartography/components/map-presets/MapPresetLoadDialog.vue'
import MapPresetManager from '@/features/cartography/components/map-presets/MapPresetManager.vue'
import MapPresetSaveDialog from '@/features/cartography/components/map-presets/MapPresetSaveDialog.vue'
import MapViewHexPopup from '@/features/cartography/components/popups/MapViewHexPopup.vue'
import MapViewSitePopup from '@/features/cartography/components/popups/MapViewSitePopup.vue'
import { palette } from '@/lib/color_brewer'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { useClipboard, useLocalStorage, useToggle } from '@vueuse/core'
import { LatLngExpression } from 'leaflet'
import { compressToEncodedURIComponent, decompressFromEncodedURIComponent } from 'lz-string'
import { computed, ref, useTemplateRef } from 'vue'
import { useRoute } from 'vue-router'
import { useLayerData } from '../components/layers-manager/layer-data'
import LayersManager from '../components/layers-manager/LayersManager.vue'
import SiteSearchPanel from '../components/SiteSearchPanel.vue'

const route = useRoute()

function getSpecQueryValue() {
  const value = route.query.h
  return Array.isArray(value) ? value[0] : value
}

const initialSpecs = computed(() => {
  const b64specs = getSpecQueryValue()
  if (!b64specs) return
  try {
    const decoded = decompressFromEncodedURIComponent(b64specs)
    return JSON.parse(decoded)
  } catch (e) {
    console.error('Failed to parse layer specs from URL', e)
  }
})

const siteMarkers = ref<SiteItem[]>([])
const mapBounds = ref<[LatLngExpression, LatLngExpression]>()

const showAllPresets = ref(false)

const drawerPinned = useLocalStorage('mapping-tool-drawer-pinned', false, {
  initOnMounted: true
})

const map = useTemplateRef('map')
const userStore = useUserStore()

const [drawer, toggleDrawer] = useToggle(false)

type MappingToolTab = 'layers' | 'sites' | 'config'

const tab = ref<MappingToolTab>('layers')

function toggleTab(newTab: MappingToolTab) {
  tab.value = newTab
  toggleDrawer(true)
}

const { feedback } = useFeedback()

const hexLayerSpec = ref<HexLayerSpec>(initialSpecs.value?.hexgrid ?? makeHexLayer())
const markerLayerSpecs = ref<MarkerLayerSpec[]>(initialSpecs.value?.markers ?? [])

const { allPending, anyLoading, layerData, data } = useLayerData()

const hexgridLayer = computed<HexgridLayer<SiteWithOccurrences>>(() => {
  const { id, name, active, config, bindings } = hexLayerSpec.value
  const remote = data.get(id)
  return {
    name,
    active,
    config: {
      ...config,
      colorRange: palette(config.colorRange ?? 'Viridis')
    },
    bindings: {
      radius: useScaleBinding(bindings.radius),
      color: useScaleBinding(bindings.color),
      opacity: useScaleBinding(bindings.opacity)
    },
    data: remote?.data.value
  }
})

const markerLayers = computed<MarkerLayer<SiteWithOccurrences>[]>(() => {
  return markerLayerSpecs.value.map((layer) => {
    const remote = data.get(layer.id)
    return {
      name: layer.name,
      config: layer.config,
      active: layer.active,
      clustered: false,
      data: remote?.data.value
    }
  })
})

const { copy } = useClipboard()

function encodeLayerSpecs() {
  const specs = {
    hexgrid: hexLayerSpec.value,
    markers: markerLayerSpecs.value
  }
  console.log(compressToEncodedURIComponent(JSON.stringify(specs)))
  return compressToEncodedURIComponent(JSON.stringify(specs))
}

function share() {
  const h = encodeLayerSpecs()
  const url = `${window.location.origin}${window.location.pathname}?h=${h}`
  copy(url)
  feedback({ message: 'Map configuration URL copied to clipboard', type: 'success' })
}
</script>

<style lang="scss">
@use 'vuetify';

.map-toolbar {
  height: 100%;
  width: 300px;
  background-color: rgb(var(--v-theme-surface));
}
</style>
