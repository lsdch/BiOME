<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="listOccurrenceDatasetsOptions"
    entity-name="Occurrence dataset"
    :toolbar="{ title: 'Occurrence datasets', icon: 'mdi-crosshairs-gps' }"
  >
    <template #item.label="{ item }: { item: OccurrenceDataset }">
      <div class="d-flex justify-space-between ga-2">
        <RouterLink
          class="text-no-wrap"
          :to="{ name: 'occurrence-dataset-item', params: { slug: item.slug } }"
          :text="item.label"
        />
        <v-icon
          :color="item.is_congruent ? 'success' : 'warning'"
          :icon="item.is_congruent ? 'mdi-check-circle' : 'mdi-alert-circle'"
          :title="item.is_congruent ? 'Congruent' : 'Not congruent'"
        ></v-icon>
      </div>
    </template>
    <template #item.description="{ value }">
      <LineClampedText :title="value" :text="value" :lines="3" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { OccurrenceDataset, OccurrenceDatasetListItem } from '@/api'
import { listOccurrenceDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { LineClampedText } from '@/components/toolkit/ui/LineClampedText'

const headers: CRUDTableHeader<OccurrenceDatasetListItem>[] = [
  { key: 'label', title: 'Label' },
  {
    key: 'description',
    title: 'Description',
    cellProps: { class: 'text-caption' }
  },
  {
    key: 'sites',
    title: 'Sites',
    align: 'end',
    cellProps: { class: 'font-monospace' }
  },
  {
    key: 'occurrences',
    title: 'Occurrences',
    align: 'end',
    cellProps: { class: 'font-monospace' }
  }
]
</script>

<style scoped></style>
