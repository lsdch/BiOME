<template>
  <l-map
    ref="map"
    v-model:zoom="zoom"
    :center="[0, 0]"
    :use-global-leaflet="true"
    :bounds="latLngBounds(latLng(90, -360), latLng(-90, 360))"
    :max-bounds="latLngBounds(latLng(90, -360), latLng(-90, 360))"
    :max-bounds-viscosity="1.0"
    :min-zoom="2"
    :max-zoom="12"
    @ready="onReady"
    :options="{
      gestureHandling: true,
      worldCopyJump: true,
      wheelPxPerZoomLevel: 100,
      zoomSnap: 0.5
    }"
  >
    <l-tile-layer
      :subdomains="['server', 'services']"
      url="https://{s}.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}"
      attributionUrl="https://static.arcgis.com/attribution/World_Imagery"
      attribution="Powered by <a href='https://www.esri.com/'>Esri</a> &mdash; Source: Esri, Maxar, Earthstar Geographics, and the GIS User Community"
      layer-type="base"
      name="Base layer"
    />
  </l-map>
</template>

<script setup lang="ts">
import L from 'leaflet'
import 'leaflet.fullscreen'
import 'leaflet.fullscreen/Control.FullScreen.css'
import { latLngBounds, latLng, Control, type Map } from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { LMap, LTileLayer } from '@vue-leaflet/vue-leaflet'
import { ref } from 'vue'

const zoom = ref(1)

const map = ref<HTMLElement>()

function onReady(mapInstance: Map) {
  L.control.fullscreen().addTo(mapInstance)
}
</script>

<style lang="scss">
@use 'vuetify';
.leaflet-container {
  background-color: rgb(var(--v-theme-surface));
}
</style>
