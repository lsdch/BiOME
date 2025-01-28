<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="listSequenceDatasetsOptions"
    entity-name="Sequence dataset"
    :toolbar="{ title: 'Sequence datasets', icon: 'mdi-crosshairs-gps' }"
  >
    <template #item.label="{ item }: { item: SequenceDataset }">
      <RouterLink
        :to="{ name: 'seq-dataset-item', params: { slug: item.slug } }"
        :text="item.label"
      />
    </template>
    <template #item.description="{ value }">
      <LineClampedText :title="value" :text="value" :lines="3" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SequenceDataset } from '@/api'
import { listSequenceDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { LineClampedText } from '@/components/toolkit/ui/LineClampedText'

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
