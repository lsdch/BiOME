<template>
  <CRUDTable :headers :items="sites" entityName="Site" density="compact">
    <template #[`item.name`]="{ item }: { item: SiteItem }">
      <RouterLink :to="{ name: 'site-item', params: { code: item.code } }">
        {{ item.name || item.code }}
      </RouterLink>
    </template>
    <template #item.country="{ value }">
      <CountryChip :country="value" label size="small" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'
import CountryChip from './CountryChip'

const { sites } = defineProps<{ sites: SiteItem[] }>()

const headers: CRUDTableHeader<SiteItem>[] = [
  { key: 'name', title: 'Name' },
  { key: 'coordinates.latitude', title: 'Latitude' },
  { key: 'coordinates.longitude', title: 'Longitude' },
  {
    key: 'country',
    title: 'Country',
    width: 0
  }
]
</script>

<style scoped lang="scss"></style>
