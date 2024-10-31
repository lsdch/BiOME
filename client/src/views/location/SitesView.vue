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
          <v-btn text="Import sites" :to="{ name: 'import-dataset' }"></v-btn>
        </template>
      </TableToolbar>
      <SitesMap ref="map" :items="sites" clustered>
        <template #marker="{ item }">
          <SitePopup :item :options="{ keepInView: false, autoPan: false }"></SitePopup>
        </template>
      </SitesMap>
    </div>
  </div>
</template>

<script setup lang="ts">
import SitesMap from '@/components/maps/SitesMap.vue'
import { ref } from 'vue'

import { LocationService } from '@/api'
import 'vue-leaflet-markercluster/dist/style.css'
import SitePopup from '@/components/sites/SitePopup.vue'
import DatasetPicker from '@/components/sites/DatasetPicker.vue'
import { useToggle } from '@vueuse/core'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { useDisplay } from 'vuetify'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'

const { xs } = useDisplay()

const [drawer, toggleDrawer] = useToggle(true)

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
