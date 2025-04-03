<template>
  <l-popup class="site-popup" :options @remove="console.log('popupclose')">
    <SiteRadius v-if="showRadius" :site="item" :zoom />
    <v-list-item
      class="text-no-wrap"
      :title="item.name"
      :subtitle="item.code"
      :to="{ name: 'site-item', params: { code: item.code } }"
    >
      <template #title="{ title }">
        <span class="text-wrap">{{ title }}</span>
      </template>
      <template #subtitle="{ subtitle }">
        <span class="font-monospace">{{ subtitle }}</span>
      </template>
    </v-list-item>
    <v-divider />
    <v-list>
      <v-list-item
        prepend-icon="mdi-crosshairs-gps"
        width="fit-content"
        @click.stop="copyCoordinates"
      >
        <div class="coordinates font-monospace cursor-pointer">
          <span class="label"> Lat </span>
          {{ item.coordinates.latitude }}
          <span class="label"> Lng </span>
          {{ item.coordinates.longitude }}
        </div>
        <v-overlay
          v-model="hasCopied"
          class="align-center justify-center"
          contained
          content-class="w-100"
        >
          <v-alert
            icon="mdi-content-copy"
            color="success"
            variant="elevated"
            density="compact"
            text="Copied"
            width="100%"
          />
        </v-overlay>
        <template #append>
          <CoordPrecisionChip :precision="item.coordinates.precision" size="small" />
        </template>
      </v-list-item>
      <v-list-item prepend-icon="mdi-map-marker">
        {{ item.locality }}
        <template #append v-if="item.country">
          <CountryChip :country="item.country" size="small" class="ml-2" />
        </template>
      </v-list-item>
    </v-list>
  </l-popup>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'
import { LPopup } from '@vue-leaflet/vue-leaflet'
import { useClipboard, useTimeoutFn, useToggle } from '@vueuse/core'
import { PopupOptions } from 'leaflet'
import CoordPrecisionChip from './CoordPrecisionChip'
import SiteRadius from './SiteRadius'
import CountryChip from './CountryChip'

const { zoom = 1, item } = defineProps<{
  item: SiteItem
  options?: PopupOptions
  showRadius?: boolean
  zoom?: number
}>()

const [hasCopied, toggleHasCopied] = useToggle(false)
const hasCopiedTimeout = useTimeoutFn(() => toggleHasCopied(false), 2000)

const { copy } = useClipboard()
function copyCoordinates() {
  hasCopiedTimeout.stop()
  copy(`${item.coordinates.latitude}, ${item.coordinates.longitude}`)
  toggleHasCopied(true)
  hasCopiedTimeout.start()
}
</script>

<style lang="scss"></style>
