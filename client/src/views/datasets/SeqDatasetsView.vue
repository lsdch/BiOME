<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="DatasetsService.listSequenceDatasets"
    entity-name="Sequence dataset"
    :toolbar="{ title: 'Sequence datasets', icon: 'mdi-crosshairs-gps' }"
  >
    <template #item.label="{ item }: { item: SequenceDataset }">
      <RouterLink
        :to="{ name: 'seq-dataset-item', params: { slug: item.slug } }"
        :text="item.label"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { DatasetsService, SequenceDataset } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<SequenceDataset>[] = [
  { key: 'label', title: 'Label' },
  { key: 'description', title: 'Description' },
  {
    key: 'sequences',
    title: 'Sequences',
    width: 0,
    value(item: SequenceDataset, fallback) {
      return item.sequences.length
    }
  }
]
</script>

<style scoped></style>
