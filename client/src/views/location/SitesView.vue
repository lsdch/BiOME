<template>
  <div class="fill-height w-100 d-flex">
    <v-navigation-drawer v-model="drawer" :location="xs ? 'top' : 'left'">
      <v-list-subheader class="px-2"> Filters </v-list-subheader>
      <v-divider></v-divider>
      <div class="mt-3 px-2">
        <DatasetPicker density="compact" />
        <CountryPicker density="compact" />
        <TaxonPicker density="compact" />
        <v-select label="Sampled" clearable density="compact" />
      </div>
    </v-navigation-drawer>
    <div class="fill-height w-100 d-flex flex-column">
      <TableToolbar
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
      </TableToolbar>
      <v-progress-linear v-if="isPending || isRefetching" indeterminate color="warning" />
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
          <template #popup="{ item, popupOpen, zoom }">
            <KeepAlive>
              <SitePopup
                v-if="item"
                :item
                :options="{ keepInView: false, autoPan: false }"
                :showRadius="popupOpen"
                :zoom
              />
            </KeepAlive>
          </template>
        </SitesMap>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import SitesMap from '@/components/maps/SitesMap.vue'

import { Site } from '@/api'
import { listSitesOptions, listSitesQueryKey } from '@/api/gen/@tanstack/vue-query.gen'
import DatasetPicker from '@/components/datasets/DatasetPicker.vue'
import { MarkerLayer } from '@/components/maps/MarkerControl.vue'
import SiteFormDialog from '@/components/sites/SiteFormDialog.vue'
import SitePopup from '@/components/sites/SitePopup.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { useLocalStorage, useToggle } from '@vueuse/core'
import { useDisplay } from 'vuetify'
import { ref, watch } from 'vue'

const { xs } = useDisplay()

const [drawer, toggleDrawer] = useToggle(true)
const [createDialog, toggleCreate] = useToggle(false)

const { data: sites, error, isPending, isRefetching, refetch } = useQuery(listSitesOptions())

const queryClient = useQueryClient()
function onCreated(newSite: Site) {
  queryClient.setQueryData<Site[]>(listSitesQueryKey(), (sites) =>
    sites ? [newSite, ...sites] : [newSite]
  )
}

const markerMode = useLocalStorage<MarkerLayer>('site-view-marker-mode', 'cluster', {
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
