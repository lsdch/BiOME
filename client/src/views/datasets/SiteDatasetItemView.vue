<template>
  <DatasetItemView :slug :dataset="dataset">
    <template #map="{ isDialog, toggleMobileMap }">
      <SitesMap :hexgrid :closable="isDialog" @close="toggleMobileMap(false)" clustered> </SitesMap>
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
import SitesMap from '@/components/maps/SitesMap.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { palette } from '@/functions/color_brewer'
import { useUserStore } from '@/stores/user'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import DatasetTabs from './DatasetTabs.vue'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'

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
