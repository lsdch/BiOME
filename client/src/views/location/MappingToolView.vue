<template>
  <div class="fill-height w-100 d-flex">
    <MappingToolFilters v-model="filters" />
    <div class="fill-height w-100 d-flex flex-column">
      <!-- <TableToolbar
        title="Sites"
        icon="mdi-map-marker"
        :togglable-search="false"
        @reload="refetch()"
      >
        <template #search>
          <v-btn icon="mdi-menu" variant="tonal" v-if="xs" @click="toggleDrawer(true)" />
        </template>
        <template #append>
          <v-btn
            text="Create site"
            class="mx-1"
            prepend-icon="mdi-plus"
            variant="tonal"
            @click="toggleCreate(true)"
          />
          <SiteFormDialog v-model:dialog="createDialog" @success="onCreated" />

          <v-btn
            text="Import sites"
            class="mx-1"
            prepend-icon="mdi-upload"
            variant="tonal"
            :to="{ name: 'import-dataset' }"
          />
        </template>
      </TableToolbar> -->
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
import SitesMap from '@/components/maps/SitesMap.vue'

import { occurrencesBySiteOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { MapLayerMode } from '@/components/maps/MarkerControl.vue'
import { useQuery } from '@tanstack/vue-query'
import { useLocalStorage } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import MapViewHexPopup from './MapViewHexPopup.vue'
import MapViewSitePopup from './MapViewSitePopup.vue'
import MappingToolFilters, { MappingFilters } from './MappingToolFilters.vue'
import { LControl } from '@vue-leaflet/vue-leaflet'

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
