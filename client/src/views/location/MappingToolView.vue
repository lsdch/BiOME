<template>
  <div class="fill-height w-100 d-flex">
    <v-navigation-drawer :location="$vuetify.display.xs ? 'top' : 'left'" :width="500">
      <v-tabs v-model="tab">
        <v-tab text="Filters" value="filters" prepend-icon="mdi-filter-variant" />
        <v-tab text="Layers" value="bindings" prepend-icon="mdi-layers" />
      </v-tabs>
      <v-divider />
      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="filters">
          <MappingToolFilters v-model="filters" />
        </v-tabs-window-item>
        <v-tabs-window-item value="bindings">
          <LayerOptionsCard
            title="Hexgrid"
            v-model="hexgridConfig.active"
            prepend-icon="mdi-hexagon-multiple-outline"
            flat
          >
            <div class="bg-main">
              <v-list class="bg-main">
                <div class="d-flex align-center">
                  <v-list-subheader>Color</v-list-subheader>
                  <v-divider />
                </div>

                <v-list-item>
                  <ScaleBindingSelect
                    label="Color binding"
                    density="compact"
                    class="my-1"
                    @update-fn="(f) => (hexgridConfig.bindings.color = f)"
                  />
                </v-list-item>
                <v-list-item>
                  <ColorPalettePicker
                    label="Palette"
                    class="my-1"
                    @update:model-value="(v) => (hexgridConfig.colorRange = palette(v))"
                  />
                </v-list-item>
                <v-list-item>
                  <ScaleBindingSelect
                    label="Opacity binding"
                    density="compact"
                    placeholder="Constant"
                    persistent-placeholder
                    clearable
                    hide-details
                    class="my-1"
                    @update-fn="(f) => (hexgridConfig.bindings.opacity = f)"
                  />
                </v-list-item>
                <ListItemInput
                  :title="hexgridConfig.bindings.opacity ? 'Opacity range' : 'Opacity'"
                >
                  <v-range-slider
                    v-if="hexgridConfig.bindings.opacity"
                    v-model="hexgridConfig.opacityRange"
                    :min="0"
                    :max="1"
                    :step="0.1"
                    :width="250"
                    hide-details
                    color="warning"
                    thumb-label
                  >
                    <template #thumb-label="{ modelValue }"> {{ modelValue * 100 }}% </template>
                  </v-range-slider>
                  <v-slider
                    v-else
                    v-model="hexgridConfig.opacity"
                    :min="0"
                    :max="1"
                    :step="0.1"
                    hide-details
                    :width="250"
                    thumb-label
                  >
                    <template #thumb-label="{ modelValue }"> {{ modelValue * 100 }}% </template>
                  </v-slider>
                </ListItemInput>

                <div class="d-flex align-center">
                  <v-list-subheader>Radius</v-list-subheader>
                  <v-divider />
                </div>
                <v-list-item title="Grid cell">
                  <template #append>
                    <v-slider
                      v-model="hexgridConfig.radius"
                      :min="2"
                      :max="20"
                      :step="1"
                      :width="250"
                      hide-details
                      thumb-label
                    />
                  </template>
                </v-list-item>

                <v-list-item>
                  <ScaleBindingSelect
                    label="Radius binding"
                    density="compact"
                    placeholder="Constant"
                    persistent-placeholder
                    clearable
                    hide-details
                    class="my-1"
                    @update-fn="(f) => (hexgridConfig.bindings.radius = f)"
                  />
                </v-list-item>
                <ListItemInput title="Radius range" v-if="hexgridConfig.bindings.radius">
                  <v-range-slider
                    v-model="hexgridConfig.radiusRange"
                    :ticks="[hexgridConfig.radius]"
                    show-ticks="always"
                    :min="2"
                    :max="20"
                    :step="0.5"
                    :width="250"
                    thumb-label
                    hide-details
                    color="warning"
                  />
                </ListItemInput>

                <div class="d-flex align-center">
                  <v-list-subheader>Hover</v-list-subheader>
                  <v-divider />
                </div>
                <v-list-item title="Fill cell">
                  <template #prepend>
                    <v-checkbox v-model="hexgridConfig.hover.fill" hide-details />
                  </template>
                </v-list-item>
                <v-list-item title="Upscale">
                  <template #prepend>
                    <v-checkbox v-model="hexgridConfig.hover.useScale" hide-details />
                  </template>
                  <template #append>
                    <v-slider
                      v-model="hexgridConfig.hover.scale"
                      :disabled="!hexgridConfig.hover.useScale"
                      :min="1"
                      :max="5"
                      :step="0.2"
                      :width="250"
                      :ticks="
                        Object.fromEntries(
                          Array.from({ length: 5 }, (_, i) => [i + 1, `×${i + 1}`])
                        )
                      "
                      show-ticks="always"
                      hide-details
                      thumb-label
                    >
                      <template #thumb-label="{ modelValue }"> ×{{ modelValue }} </template>
                    </v-slider>
                  </template>
                </v-list-item>
              </v-list>
            </div>
          </LayerOptionsCard>

          <v-divider />

          <LayerOptionsCard
            title="Markers"
            v-model="markerConfig.active"
            prepend-icon="mdi-circle-multiple-outline"
          >
            <div class="bg-main">
              <v-list>
                <ListItemInput label="Clustered" subtitle="Aggregate marker clusters">
                  <v-switch v-model="markerConfig.clustered" hide-details />
                </ListItemInput>
                <ListItemInput label="Radius">
                  <v-slider
                    :min="1"
                    :max="20"
                    :step="0.5"
                    v-model="markerConfig.radius"
                    hide-details
                    :width="250"
                    thumb-label
                  />
                </ListItemInput>
                <ListItemInput label="Stroke color" subtitle="Hue and opacity">
                  <ColorPickerMenu v-model="markerConfig.color" hide-details />
                </ListItemInput>
                <ListItemInput label="Fill color" subtitle="Hue and opacity">
                  <ColorPickerMenu v-model="markerConfig.fillColor" hide-details />
                </ListItemInput>
                <ListItemInput label="Stroke width">
                  <v-slider
                    :min="1"
                    :max="5"
                    v-model="markerConfig.weight"
                    hide-details
                    :width="250"
                    thumb-label
                  />
                </ListItemInput>
              </v-list>
            </div>
          </LayerOptionsCard>
          <v-divider />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-navigation-drawer>
    <div class="fill-height w-100 d-flex flex-column">
      <v-progress-linear v-if="isPending && !initialFetchDone" indeterminate color="warning" />
      <div class="fill-height w-100 position-relative">
        <v-overlay
          contained
          :model-value="!isRefetching && !!error"
          class="align-center justify-center"
        >
          <v-alert color="error" variant="elevated">Failed to load sampling sites</v-alert>
        </v-overlay>
        <SitesMap
          ref="map"
          :items="sites"
          clustered
          :auto-fit="(sites?.length ?? 0) > 1"
          v-model:marker-mode="markerMode"
          :marker-config
          :hexgrid-config
        >
          <LControl v-if="isRefetching || isFetching" position="topleft">
            <v-progress-circular
              v-if="isPending || isRefetching"
              indeterminate
              color="warning"
              size="32"
              width="6"
            />
          </LControl>
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
import SitesMap, { HexgridConfig, MarkerConfig } from '@/components/maps/SitesMap.vue'

import { SiteWithOccurrences } from '@/api'
import { occurrencesBySiteOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { MapLayerMode } from '@/components/maps/MarkerControl.vue'
import ColorPickerMenu from '@/components/toolkit/ui/ColorPickerMenu.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import { palette } from '@/views/location/color_brewer'
import { useQuery } from '@tanstack/vue-query'
import { LControl } from '@vue-leaflet/vue-leaflet'
import { useLocalStorage } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import ColorPalettePicker from './ColorPalettePicker.vue'
import LayerOptionsCard from './LayerOptionsCard.vue'
import MapViewHexPopup from './MapViewHexPopup.vue'
import MapViewSitePopup from './MapViewSitePopup.vue'
import MappingToolFilters, { MappingFilters } from './MappingToolFilters.vue'
import ScaleBindingSelect from './ScaleBindingSelect.vue'

const tab = ref<'filters' | 'bindings'>('filters')

const hexgridConfig = ref<HexgridConfig<SiteWithOccurrences>>({
  active: true,
  radius: 10,
  radiusRange: [0, 10],
  hover: {
    fill: true,
    useScale: false,
    scale: 1.5
  },
  colorRange: palette('Viridis'),
  opacity: 0.8,
  opacityRange: [0, 1],
  bindings: {}
})

const markerConfig = ref<MarkerConfig>({
  active: false,
  clustered: false,
  radius: 4,
  color: '#FF000055',
  fillColor: '#FF000000',
  weight: 1
})

const filters = ref<MappingFilters>({})

const {
  data: sites,
  error,
  isPending,
  isFetching,
  isRefetching,
  refetch
} = useQuery(
  computed(() =>
    occurrencesBySiteOptions({
      query: {
        ...filters.value,
        habitats: filters.value.habitats?.map(({ label }) => label)
      }
    })
  )
)

const initialFetchDone = ref(false)
watch(isPending, (pending) => {
  if (!pending && !initialFetchDone.value) {
    initialFetchDone.value = true
  }
})

const markerMode = useLocalStorage<MapLayerMode>('site-view-marker-mode', 'markers', {
  initOnMounted: true
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
