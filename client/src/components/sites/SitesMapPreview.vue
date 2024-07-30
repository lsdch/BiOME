<template>
  <v-dialog v-model="isOpen" fullscreen>
    <SitesMap
      ref="map"
      closable
      auto-fit
      :items="mapItems"
      :marker="{
        color: 'white',
        fill: true,
        fillColor: 'orangered',
        fillOpacity: 1,
        radius: 8
      }"
      @close="isOpen = false"
    >
      <template #marker="{ item }">
        <l-popup class="site-popup">
          <div class="text-subtitle-1 text-no-wrap">{{ item.name }}</div>
          <div class="text-overline text-secondary">
            {{ item.code }}
          </div>
          <v-divider></v-divider>
          <div class="d-flex align-center my-2">
            <v-icon icon="mdi-crosshairs-gps" size="small" class="mr-3" />
            <div class="coordinates">
              <span class="label"> Lat </span>
              {{ item.coordinates.latitude }}
              <span class="label"> Lng </span>
              {{ item.coordinates.longitude }}
            </div>
          </div>
          <div v-if="item.locality || item.country_code" class="d-flex align-center my-2">
            <v-icon icon="mdi-map-marker" size="small" class="mr-3" />
            <div class="justify-space-between">
              {{ item.locality }}
              <v-chip :text="item.country_code" size="small" class="ml-2" />
            </div>
          </div>
        </l-popup>
      </template>
    </SitesMap>
  </v-dialog>
  <slot name="activator" :open :close></slot>
</template>

<script setup lang="ts">
import { LPopup } from '@vue-leaflet/vue-leaflet'
import { computed } from 'vue'
import { ImportItem } from '.'
import { Geocoordinates } from '../maps'
import SitesMap from '../maps/SitesMap.vue'
import { RecordElement } from './SiteTabularImport.vue'

const isOpen = defineModel<boolean>('open', { default: false })

const props = defineProps<{ sites: RecordElement[] }>()

type MapItem = Omit<RecordElement, 'coordinates'> & Geocoordinates

const mapItems = computed<MapItem[]>(() => {
  return props.sites.filter(
    (s) => s.coordinates && s.coordinates.latitude && s.coordinates.longitude
  ) as MapItem[]
})

function open() {
  isOpen.value = true
}

function close() {
  isOpen.value = false
}
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
