<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="DatasetsService.listDatasets"
    entity-name="Dataset"
    :toolbar="{ title: 'Datasets', icon: 'mdi-folder-table' }"
  >
    <template #item.label="{ item }: { item: Dataset }">
      <RouterLink :to="{ name: 'dataset-item', params: { slug: item.slug } }" :text="item.label" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Dataset, DatasetsService } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader[] = [
  { key: 'label', title: 'Label' },
  { key: 'description', title: 'Description' },
  {
    key: 'sites',
    title: 'Sites',
    width: 0,
    value(item, fallback) {
      return item.sites.length
    }
  }
]
</script>

<style scoped></style>
