<template>
  <DatasetItemView :slug :dataset="dataset">
    <template #map="{ isDialog, toggleMobileMap }">
      <BaseMap :hexgrid :closable="isDialog" @close="toggleMobileMap(false)" clustered> </BaseMap>
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
import { getSiteDatasetOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'
import BaseMap from '@/features/cartography/components/BaseMap.vue'
import { palette } from '@/lib/color_brewer'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import DatasetTabs from '../components/DatasetTabs.vue'

const { slug } = defineProps<{
  slug: string
}>()

const { data: dataset, error, isPending } = useQuery(getSiteDatasetOptions({ path: { slug } }))
const hexgrid = computed(() => {
  return {
    data: dataset.value?.sites,
    active: true,
    bindings: {},
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
