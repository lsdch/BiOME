<template>
  <l-popup class="site-popup" :options>
    <v-list-item
      class="text-no-wrap"
      :title="item.name"
      :subtitle="item.code"
      :to="{ name: 'site-item', params: { code: item.code } }"
    />
    <v-divider />
    <v-card-text>
      <div class="d-flex align-center my-2">
        <v-icon icon="mdi-crosshairs-gps" size="small" class="mr-3" />
        <div class="coordinates">
          <span class="label"> Lat </span>
          {{ item.coordinates.latitude }}
          <span class="label"> Lng </span>
          {{ item.coordinates.longitude }}
        </div>
      </div>
      <div v-if="item.locality || item.country.code" class="d-flex align-center my-2">
        <v-icon icon="mdi-map-marker" size="small" class="mr-3" />
        <div class="justify-space-between">
          {{ item.locality }}
          <v-chip :text="item.country.code" size="small" class="ml-2" />
        </div>
      </div>
    </v-card-text>
  </l-popup>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'
import { LPopup } from '@vue-leaflet/vue-leaflet'
import { PopupOptions } from 'leaflet'

defineProps<{ item: SiteItem; options?: PopupOptions }>()
</script>

<style lang="scss">
@use 'vuetify';
.leaflet-popup {
  .v-list-item-subtitle {
    font-size: 0.8rem;
  }
  .leaflet-popup-content-wrapper,
  .leaflet-popup-tip {
    background-color: rgb(var(--v-theme-surface));
    color: rgb(var(--v-theme-on-surface));
  }
  .leaflet-popup-content-wrapper {
    border-radius: 0;
    .leaflet-popup-content {
      margin: 0px;

      .coordinates {
        display: grid;
        grid-template-columns: [label] 0fr [value] 1fr;
        grid-template-rows: 1fr 1fr;
        column-gap: 10px;

        .label {
          grid-column: label;
          font-family: monospace;
        }
      }
    }
  }
}
</style>
