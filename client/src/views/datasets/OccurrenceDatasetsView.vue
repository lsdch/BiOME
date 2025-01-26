<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="DatasetsService.listOccurrenceDatasets"
    entity-name="Occurrence dataset"
    :toolbar="{ title: 'Occurrence datasets', icon: 'mdi-crosshairs-gps' }"
  >
    <template #item.label="{ item }: { item: OccurrenceDataset }">
      <RouterLink
        :to="{ name: 'occurrence-dataset-item', params: { slug: item.slug } }"
        :text="item.label"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { DatasetsService, OccurrenceDataset } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<OccurrenceDataset>[] = [
  { key: 'label', title: 'Label' },
  { key: 'description', title: 'Description' },
  {
    key: 'occurrences',
    title: 'Occurrences',
    width: 0,
    value(item: OccurrenceDataset, fallback) {
      return item.occurrences.length
    }
  },
  {
    key: 'is_congruent',
    title: 'Congruent'
  }
]
</script>

<style scoped></style>
