<template>
  <BaseMap :hexgrid>
    <template #hex-popup="{ data }">
      <MapViewHexPopup :data />
    </template>
    <template #popup="{ item, popupOpen, zoom }">
      <MapViewSitePopup :item :popupOpen :zoom :key="item.code" />
    </template>
  </BaseMap>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import BaseMap from '@/components/maps/BaseMap.vue'
import { HexgridLayer } from '@/components/maps/map-layers'
import MapViewHexPopup from '@/components/occurrence/MapViewHexPopup.vue'
import MapViewSitePopup from '@/components/occurrence/MapViewSitePopup.vue'
import { palette } from '@/functions/color_brewer'
import { computed } from 'vue'

const { sites } = defineProps<{
  sites: SiteWithOccurrences[]
}>()

const hexgrid = computed<HexgridLayer<SiteWithOccurrences>>(() => {
  return {
    data: sites,
    active: true,
    bindings: {
      color: (d) =>
        d.reduce((a, b) => a + b.data.samplings.flatMap(({ occurrences }) => occurrences).length, 0)
    },
    config: {
      radius: 8,
      opacity: 1,
      colorRange: palette('Viridis'),
      hover: {
        fill: true,
        useScale: false,
        scale: 1
      }
    }
  }
})
</script>

<style scoped lang="scss"></style>
