<template>
  <div class="fill-height w-100 d-flex">
    <div class="map-toolbar">FILTERS</div>
    <div class="fill-height w-100 d-flex flex-column">
      <v-toolbar>
        <v-toolbar-title> Sites </v-toolbar-title>
        <v-spacer />
        <v-btn text="Import sites" :to="{ name: 'import-sites' }"></v-btn>
      </v-toolbar>
      <SitesMap />
    </div>
  </div>
</template>

<script setup lang="ts">
import SitesMap from '@/components/maps/SitesMap.vue'
</script>

<style lang="scss">
@use 'vuetify';

.map-toolbar {
  height: 100%;
  width: 300px;
  background-color: rgb(var(--v-theme-surface));
}
</style>
