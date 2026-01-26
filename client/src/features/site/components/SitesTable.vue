<template>
  <v-text-field v-model="search.term" label="Search" class="ma-2" clearable />
  <CRUDTable :headers :items="sites" entityName="Site" density="compact" :search>
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
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import CountryChip from './CountryChip'
import { ref } from 'vue'

const { sites } = defineProps<{ sites: SiteItem[] }>()

const search = ref({ term: undefined })

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
