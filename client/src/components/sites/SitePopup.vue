<template>
  <l-popup class="site-popup" :options>
    <div class="text-subtitle-1 text-no-wrap">{{ item.name }}</div>
    <div class="text-overline text-secondary">
      <v-btn
        block
        color="primary"
        prepend-icon="mdi-identifier"
        variant="text"
        :text="item.code"
        class="px-0"
        :rounded="0"
        :to="{ name: 'site-item', params: { code: item.code } }"
      />
    </div>
    <v-divider />
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
  .leaflet-popup-content-wrapper,
  .leaflet-popup-tip {
    background-color: rgb(var(--v-theme-surface));
    color: rgb(var(--v-theme-on-surface));
  }
  .leaflet-popup-content-wrapper {
    border-radius: 0;
    .leaflet-popup-content {
      margin-top: 0px;
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
