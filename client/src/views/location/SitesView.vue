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
      <TableToolbar title="Sites" icon="mdi-map-marker" :togglable-search="false">
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
          <SiteFormDialog v-model="createDialog" @success="(site) => sites.unshift(site)" />

          <v-btn
            text="Import sites"
            class="mx-1"
            prepend-icon="mdi-upload"
            variant="tonal"
            :to="{ name: 'import-dataset' }"
          />
        </template>
      </TableToolbar>
      <SitesMap ref="map" :items="sites" clustered :auto-fit="sites.length > 1">
        <template #popup="{ item }">
          <KeepAlive>
            <SitePopup v-if="item" :item :options="{ keepInView: false, autoPan: false }" />
          </KeepAlive>
        </template>
      </SitesMap>
    </div>
  </div>
</template>

<script setup lang="ts">
import SitesMap from '@/components/maps/SitesMap.vue'
import { ref } from 'vue'

import { LocationService } from '@/api'
import SitePopup from '@/components/sites/SitePopup.vue'
import DatasetPicker from '@/components/datasets/DatasetPicker.vue'
import { useToggle } from '@vueuse/core'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { useDisplay } from 'vuetify'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import SiteFormDialog from '@/components/sites/SiteFormDialog.vue'
import { LCircle } from '@vue-leaflet/vue-leaflet'

const { xs } = useDisplay()

const [drawer, toggleDrawer] = useToggle(true)
const [createDialog, toggleCreate] = useToggle(false)

const sites = ref(
  (await LocationService.listSites().then(({ data, error }) => {
    if (error) {
      console.error(error)
      return []
    }
    return data
  })) ?? []
)
</script>

<style lang="scss">
@use 'vuetify';

.map-toolbar {
  height: 100%;
  width: 300px;
  background-color: rgb(var(--v-theme-surface));
}
</style>
