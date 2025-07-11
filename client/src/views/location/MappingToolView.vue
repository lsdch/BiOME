<template>
  <div class="fill-height w-100 d-flex">
    <v-navigation-drawer
      :location="$vuetify.display.xs ? 'top' : 'left'"
      :width="500"
      v-model="drawer"
      :temporary="!drawerPinned"
    >
      <v-tabs v-model="tab">
        <v-tab value="feeds" prepend-icon="mdi-database-arrow-right">
          Data feeds
          <v-badge inline :content="feeds.length" color="purple" />
        </v-tab>
        <!-- <v-tab text="Filters" value="filters" prepend-icon="mdi-filter-variant" /> -->
        <v-tab value="bindings" prepend-icon="mdi-layers">
          Layers
          <v-badge inline :content="markerLayerOptions.length + 1" color="success" />
        </v-tab>
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
      <v-tabs-window v-model="tab">
        <v-tabs-window-item eager value="feeds">
          <OccurrenceDataFeedManager />
        </v-tabs-window-item>
        <v-tabs-window-item value="bindings">
          <HexgridLayerCard v-model="hexgridLayerOptions" />
          <v-divider />
          <MarkerLayerCard
            v-for="(markerLayer, i) in markerLayerOptions"
            v-model="markerLayerOptions[i]"
            @delete="markerLayerOptions.splice(i, 1)"
            @reset="resetMarkerLayer(i)"
          />
          <v-divider />
          <div class="d-flex">
            <v-btn
              stacked
              class="flex-grow-1"
              variant="text"
              size="small"
              :rounded="0"
              prepend-icon="mdi-plus"
              text="Add marker layer"
              @click="addMarkerLayer()"
            />
            <v-divider vertical />
            <ConfirmDialog
              title="Reset layers"
              message="Are you sure you want to reset all layers?"
              @confirm="resetLayers()"
            >
              <template #activator="{ props }">
                <v-btn
                  stacked
                  class="flex-grow-1"
                  variant="text"
                  size="small"
                  :rounded="0"
                  prepend-icon="mdi-restore"
                  text="Reset layers"
                  v-bind="props"
                />
              </template>
            </ConfirmDialog>
          </div>
          <v-divider />
        </v-tabs-window-item>
        <v-tabs-window-item value="config">
          <v-list>
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
      <template #append>
        <v-divider />
        <div class="d-flex justify-space-between pa-2">
          <div>
            <MapPresetSaveDialog
              v-if="userStore.isAuthenticated"
              :specs="{
                feeds: feeds,
                hexgrid: hexgridLayerOptions,
                markers: markerLayerOptions
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
                ({ spec: { feeds, hexgrid, markers }, name }) => {
                  feeds.splice(0, feeds.length, ...feeds)
                  hexgridLayerOptions = hexgrid
                  markerLayerOptions.splice(0, markers.length, ...markers)
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
              prepend-icon="mdi-database-arrow-right"
              title="Data feeds"
              @click="toggleTab('feeds')"
              :active="tab === 'feeds' && drawer"
              color="primary"
            />
          </template>
          <v-sheet :height="48" class="my-0 d-flex align-center"> Data feeds </v-sheet>
        </v-tooltip>
        <v-tooltip content-class="bg-surface text-overline py-0" :height="48">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-layers"
              @click="toggleTab('bindings')"
              :active="tab === 'bindings' && drawer"
              color="primary"
            />
          </template>
          <v-sheet :height="48" class="my-0 d-flex align-center"> Layers </v-sheet>
        </v-tooltip>
      </v-list>
    </v-navigation-drawer>
    <div class="fill-height w-100 d-flex flex-column">
      <v-progress-linear v-if="allPending" indeterminate color="warning" />
      <div :class="['fill-height w-100 position-relative']">
        <BaseMap
          ref="map"
          auto-fit
          :marker-layers
          :hexgrid="hexgridLayer"
          v-model:polygon-mode="polygonMode"
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
import BaseMap from '@/components/maps/BaseMap.vue'

import { SiteWithOccurrences } from '@/api'
import HexgridLayerCard from '@/components/maps/HexgridLayerCard.vue'
import MarkerLayerCard from '@/components/maps/MarkerLayerCard.vue'
import {
  HexgridLayer,
  HexgridLayerSpec,
  MarkerLayer,
  MarkerLayerDefinition,
  SitesFilter
} from '@/components/maps/map-layers'
import OccurrenceDataFeedManager from '@/components/occurrence/OccurrenceDataFeedManager.vue'
import { useDataFeeds } from '@/components/occurrence/data_feeds'
import ConfirmDialog from '@/components/toolkit/ui/ConfirmDialog.vue'
import { useScaleBinding } from '@/composables/occurrences'
import { palette, withOpacity } from '@/functions/color_brewer'
import { useLocalStorage, useToggle } from '@vueuse/core'
import { computed, ref } from 'vue'
import MapViewHexPopup from '../../components/occurrence/MapViewHexPopup.vue'
import MapViewSitePopup from '../../components/occurrence/MapViewSitePopup.vue'
import MapPresetSaveDialog from '@/components/maps/MapPresetSaveDialog.vue'
import MapPresetLoadDialog from '@/components/maps/MapPresetLoadDialog.vue'
import { useFeedback } from '@/stores/feedback'
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'
import MapPresetManager from '@/components/maps/MapPresetManager.vue'
import { useUserStore } from '@/stores/user'

const [polygonMode, togglePolygonMode] = useToggle(false)

const showAllPresets = ref(false)

const drawerPinned = useLocalStorage('mapping-tool-drawer-pinned', false, {
  initOnMounted: true
})

const userStore = useUserStore()

const [drawer, toggleDrawer] = useToggle(false)

type MappingToolTab = 'feeds' | 'bindings' | 'config'

const tab = ref<MappingToolTab>('feeds')

function toggleTab(newTab: MappingToolTab) {
  tab.value = newTab
  toggleDrawer(true)
}

const { feedback } = useFeedback()

const { data, feeds, allPending, anyLoading } = useDataFeeds()

function filterSites(
  sites: SiteWithOccurrences[] | undefined,
  filterType: SitesFilter
): SiteWithOccurrences[] | undefined {
  if (!sites) return undefined
  switch (filterType) {
    case 'Sampled':
      return sites.filter((site) => site.samplings.length > 0)
    case 'Occurrences':
      return sites.filter((site) =>
        site.samplings.some(({ occurrences }) => occurrences.length > 0)
      )
    default:
      return sites
  }
}

const hexgridLayerOptions = useLocalStorage<HexgridLayerSpec>(
  'map-tool-hexgrid-layer',
  {
    name: 'Hexgrid',
    active: true,
    dataFeedID: feeds.value[0].id,
    filterType: 'Occurrences',
    config: {
      radius: 10,
      radiusRange: [0, 10],
      hover: {
        fill: true,
        useScale: false,
        scale: 1.5
      },
      colorRange: 'Viridis',
      opacity: 0.8,
      opacityRange: [0, 1]
    },
    bindings: {
      color: { log: false, binding: 'sites' },
      opacity: { log: false },
      radius: { log: false }
    }
  },
  { deep: true }
)

const hexgridLayer = computed<HexgridLayer<SiteWithOccurrences>>(() => {
  const { dataFeedID, name, active, config, bindings, filterType } = hexgridLayerOptions.value
  const remote = dataFeedID ? data.get(dataFeedID) : undefined
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
    data: filterSites(remote?.data.value, filterType)
  }
})

const markerLayerOptions = useLocalStorage<MarkerLayerDefinition[]>('maptool-marker-layers', [], {
  deep: true
})
const markerLayers = computed<MarkerLayer<SiteWithOccurrences>[]>(() => {
  return markerLayerOptions.value.map((layer) => {
    const remote = layer.dataFeedID ? data.get(layer.dataFeedID) : undefined
    return {
      name: layer.name,
      config: layer.config,
      active: layer.active,
      clustered: false,
      data: filterSites(remote?.data.value, layer.filterType)
    }
  })
})

const markerColorPalette = [
  '#e41a1c',
  '#377eb8',
  '#4daf4a',
  '#984ea3',
  '#ff7f00',
  '#ffff33',
  '#a65628',
  '#f781bf'
]

function newMarkerLayer(index: number = markerLayerOptions.value.length): MarkerLayerDefinition {
  return {
    filterType: 'Occurrences',
    active: true,
    clustered: false,
    config: {
      radius: 4,
      color: withOpacity(markerColorPalette[index % markerColorPalette.length], 0.8),
      fillColor: withOpacity(markerColorPalette[index % markerColorPalette.length], 0.3),
      weight: 2
    }
  }
}

function addMarkerLayer(index: number = markerLayerOptions.value.length) {
  const layer = newMarkerLayer(index)
  markerLayerOptions.value.push(layer)
  return layer
}

function resetMarkerLayer(index: number) {
  if (index < 0 || index >= markerLayerOptions.value.length) return
  markerLayerOptions.value[index] = newMarkerLayer(index)
}

function resetLayers() {
  hexgridLayerOptions.value.config = {
    radius: 10,
    radiusRange: [0, 10],
    hover: {
      fill: true,
      useScale: false,
      scale: 1.5
    },
    colorRange: 'Viridis',
    opacity: 0.8,
    opacityRange: [0, 1]
  }
  hexgridLayerOptions.value.bindings = {}
  markerLayerOptions.value.length = 0
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
