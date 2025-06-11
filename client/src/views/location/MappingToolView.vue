<template>
  <div class="fill-height w-100 d-flex">
    <v-navigation-drawer
      :location="$vuetify.display.xs ? 'top' : 'left'"
      :width="500"
      v-model="drawer"
      :temporary="!drawerPinned"
    >
      <template #append>
        <v-divider />
        <div class="d-flex justify-space-between pa-2">
          <div>
            <v-tooltip location="top">
              <template #activator="{ props }">
                <v-btn v-bind="props" variant="text" icon="mdi-content-save-plus"></v-btn>
              </template>
              Save map view settings
            </v-tooltip>
            <v-tooltip location="top">
              <template #activator="{ props }">
                <v-btn v-bind="props" variant="text" icon="mdi-file-download"></v-btn>
              </template>
              Load settings
            </v-tooltip>
          </div>
          <v-tooltip>
            <template #activator="{ props }">
              <v-btn
                v-bind="props"
                icon="mdi-pin"
                size="small"
                :variant="drawerPinned ? 'tonal' : 'plain'"
                @click="drawerPinned = !drawerPinned"
              />
            </template>
            Toggle permanent options menu
          </v-tooltip>
        </div>
      </template>
      <v-tabs v-model="tab">
        <v-tab value="feeds" prepend-icon="mdi-database-arrow-right">
          Data feeds
          <v-badge inline :content="registry.length" color="purple" />
        </v-tab>
        <!-- <v-tab text="Filters" value="filters" prepend-icon="mdi-filter-variant" /> -->
        <v-tab value="bindings" prepend-icon="mdi-layers">
          Layers
          <v-badge inline :content="markerLayerOptions.length + 1" color="success" />
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
        <!-- <v-tabs-window-item value="filters">
          <MappingToolFilters v-model="filters" />
        </v-tabs-window-item> -->
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
      </v-tabs-window>
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
        <!-- <v-overlay
          contained
          :model-value="!isRefetching && !!error"
          class="align-center justify-center"
        >
          <v-alert color="error" variant="elevated">Failed to load sampling sites</v-alert>
        </v-overlay> -->
        <SitesMap
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
        </SitesMap>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import SitesMap, {
  HexgridConfig,
  HexgridLayer,
  HexgridScaleBindings,
  MarkerLayer
} from '@/components/maps/SitesMap.vue'

import { SiteWithOccurrences } from '@/api'
import {
  HexgridLayerDefinition,
  MarkerLayerDefinition,
  SitesFilter
} from '@/components/maps/map-layers'
import DataFeedPicker from '@/components/occurrence/DataFeedPicker.vue'
import OccurrenceDataFeedManager from '@/components/occurrence/OccurrenceDataFeedManager.vue'
import { useDataFeeds } from '@/components/occurrence/data_feeds'
import ConfirmDialog from '@/components/toolkit/ui/ConfirmDialog.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import { palette, withOpacity } from '@/functions/color_brewer'
import { useLocalStorage, useToggle } from '@vueuse/core'
import { UUID } from 'crypto'
import { computed, reactive, ref } from 'vue'
import ColorPalettePicker from '../../components/toolkit/ui/ColorPalettePicker.vue'
import LayerOptionsCard from './LayerOptionsCard.vue'
import MapViewHexPopup from '../../components/occurrence/MapViewHexPopup.vue'
import MapViewSitePopup from '../../components/occurrence/MapViewSitePopup.vue'
import ScaleBindingSelect from './ScaleBindingSelect.vue'
import SiteSamplingStatusFilter from './SiteSamplingStatusFilter.vue'
import MarkerLayerCard from '@/components/maps/MarkerLayerCard.vue'
import HexgridLayerCard from '@/components/maps/HexgridLayerCard.vue'

const [polygonMode, togglePolygonMode] = useToggle(false)

const drawerPinned = useLocalStorage('mapping-tool-drawer-pinned', false, {
  initOnMounted: true
})

const [drawer, toggleDrawer] = useToggle(false)

const tab = ref<'bindings' | 'feeds'>('feeds')

function toggleTab(newTab: 'feeds' | 'bindings') {
  tab.value = newTab
  toggleDrawer(true)
}

const { remotes, registry, allPending, anyLoading } = useDataFeeds()

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

const hexgridLayerOptions = reactive<HexgridLayerDefinition>({
  name: 'Hexgrid',
  active: true,
  dataFeedID: registry.value[0].id,
  filterType: 'Occurrences',
  config: {
    radius: 10,
    radiusRange: [0, 10],
    hover: {
      fill: true,
      useScale: false,
      scale: 1.5
    },
    colorRange: palette('Viridis'),
    opacity: 0.8,
    opacityRange: [0, 1]
  },
  bindings: {}
})

const hexgridLayer = computed<HexgridLayer<SiteWithOccurrences>>(() => {
  const feedID = hexgridLayerOptions.dataFeedID
  const remote = feedID ? remotes.get(feedID) : undefined
  return {
    name: hexgridLayerOptions.name,
    active: hexgridLayerOptions.active,
    config: hexgridLayerOptions.config,
    bindings: hexgridLayerOptions.bindings,
    data: filterSites(remote?.data.value, hexgridLayerOptions.filterType)
  }
})

const markerLayerOptions = reactive<MarkerLayerDefinition[]>([])
const markerLayers = computed<MarkerLayer<SiteWithOccurrences>[]>(() => {
  return markerLayerOptions.map((layer) => {
    const remote = layer.dataFeedID ? remotes.get(layer.dataFeedID) : undefined
    return {
      name: layer.name,
      config: layer.config,
      active: layer.active,
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

function newMarkerLayer(index: number = markerLayerOptions.length): MarkerLayerDefinition {
  return {
    filterType: 'Occurrences',
    active: true,
    config: {
      clustered: false,
      radius: 4,
      color: withOpacity(markerColorPalette[index % markerColorPalette.length], 0.8),
      fillColor: withOpacity(markerColorPalette[index % markerColorPalette.length], 0.3),
      weight: 2
    }
  }
}

function addMarkerLayer(index: number = markerLayerOptions.length) {
  const layer = newMarkerLayer(index)
  markerLayerOptions.push(layer)
  return layer
}

function resetMarkerLayer(index: number) {
  if (index < 0 || index >= markerLayerOptions.length) return
  markerLayerOptions[index] = newMarkerLayer(index)
}

function resetLayers() {
  hexgridLayerOptions.config = {
    radius: 10,
    radiusRange: [0, 10],
    hover: {
      fill: true,
      useScale: false,
      scale: 1.5
    },
    colorRange: palette('Viridis'),
    opacity: 0.8,
    opacityRange: [0, 1]
  }
  hexgridLayerOptions.bindings = {}
  markerLayerOptions.length = 0
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
