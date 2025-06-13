<template>
  <DatasetItemView :slug :dataset="dataset">
    <template #map="{ isDialog, toggleMobileMap }">
      <BaseMap :hexgrid :closable="isDialog" @close="toggleMobileMap(false)" clustered>
        <template #hex-popup="{ data }">
          <MapViewHexPopup :data />
        </template>
        <template #popup="{ item, popupOpen, zoom }">
          <KeepAlive>
            <MapViewSitePopup :item :popupOpen :zoom :key="item.code" />
          </KeepAlive>
        </template>
      </BaseMap>
    </template>
    <template #details>
      <CenteredSpinner v-if="isPending" :height="300" size="large" color="primary" />
      <PageErrors v-else-if="error" :error class="flex-grow-1" />
      <div v-else-if="dataset" class="flex-grow-1">
        <DatasetTabs :dataset flat />
      </div>
    </template>
  </DatasetItemView>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import { getOccurrenceDatasetOptions } from '@/api/gen/@tanstack/vue-query.gen'
import BaseMap, { HexgridLayer } from '@/components/maps/BaseMap.vue'
import MapViewHexPopup from '@/components/occurrence/MapViewHexPopup.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { palette } from '@/functions/color_brewer'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import DatasetItemView from './DatasetItemView.vue'
import DatasetTabs from './DatasetTabs.vue'
import MapViewSitePopup from '@/components/occurrence/MapViewSitePopup.vue'

const { slug } = defineProps<{
  slug: string
}>()

const {
  data: dataset,
  error,
  refetch,
  isPending
} = useQuery(getOccurrenceDatasetOptions({ path: { slug } }))

const hexgrid = computed<HexgridLayer<SiteWithOccurrences>>(() => {
  return {
    data: dataset.value?.sites,
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

<style lang="scss">
.v-list-item.empty .v-list-item-subtitle {
  font-style: italic;
}
</style>
