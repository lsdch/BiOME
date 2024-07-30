<template>
  <div class="fill-height w-100 d-flex">
    <div class="map-toolbar">
      <DatasetPicker />
    </div>
    <div class="fill-height w-100 d-flex flex-column">
      <v-toolbar>
        <v-toolbar-title> Sites </v-toolbar-title>
        <v-spacer />
        <v-btn text="Import sites" :to="{ name: 'import-sites' }"></v-btn>
      </v-toolbar>
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
