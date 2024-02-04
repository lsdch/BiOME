<template>
  <CRUDTable
    :crud="{
      list: () => TaxonomyService.taxonomyList(),
      update: (item: TaxonWithRelatives) =>
        TaxonomyService.updateTaxon(item.code, item as TaxonInput),
      delete: (item: TaxonWithRelatives) => TaxonomyService.deleteTaxon(item.code)
    }"
    :toolbar-props="{
      title: 'Taxonomy',
      entityName: 'Taxon',
      itemRepr: (item: TaxonWithRelatives) => item.name,
      icon: 'mdi-graph',
      togglableSearch: true
    }"
    :headers="headers"
    :filter="filter"
    showActions
    density="compact"
    v-model:search="searchName"
    :filter-keys="['name']"
  >
    <template v-slot:search>
      <TaxaTableFilters v-model="filters" v-model:name="searchName" />
    </template>
    <template v-slot:[`item.code`]="{ item }">
      <code>{{ item.code }}</code>
    </template>
    <template v-slot:[`item.status`]="{ item }">
      <StatusIcon :status="item.status" size="small" />
      <LinkIconGBIF v-if="item.GBIF_ID" :GBIF_ID="item.GBIF_ID" variant="text" size="x-small" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { TaxonInput, TaxonStatus, TaxonWithRelatives, TaxonomyService } from '@/api'
import { computed } from 'vue'
import CRUDTable from '../toolkit/CRUDTable.vue'
import LinkIconGBIF from './LinkIconGBIF.vue'
import StatusIcon from './StatusIcon.vue'
import TaxaTableFilters from './TaxaTableFilters.vue'

const searchName = ref('')
const filters = ref({
  rank: undefined,
  status: TaxonStatus.Accepted
})

const filter = computed(() => {
  const { rank, status } = filters.value
  if (rank || status)
    return (item: TaxonWithRelatives) => {
      return (rank ? item.rank === rank : true) && (status ? item.status === status : true)
    }
  else return undefined
})

const headers: ReadonlyHeaders = [
  { title: 'Name', key: 'name' },
  { title: 'Code', key: 'code' },
  { title: 'Rank', key: 'rank' },
  { title: 'Status', key: 'status', width: 0, align: 'center' }
]
</script>

<style scoped></style>
