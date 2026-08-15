<template>
  <CRUDTable
    class="fill-height"
    :headers
    :items
    @reload="refetch()"
    :error
    :loading
    entity-name="Occurrence dataset"
    :toolbar="{ title: 'Occurrence datasets', icon: 'mdi-crosshairs-gps' }"
  >
    <template #item.label="{ item }: { item: Dataset }">
      <div class="d-flex justify-space-between ga-2">
        <RouterLink
          class="text-no-wrap"
          :to="{ name: 'occurrence-dataset-item', params: { slug: item.slug } }"
          :text="item.label"
        />
      </div>
    </template>
    <template #item.description="{ value }">
      <LineClampedText :title="value" :text="value" :lines="3" />
    </template>
    <!-- <template #expanded-row-inject="{ item }">
      <div class="d-flex ga-2 align-center ma-2">
        <span class="text-muted"> Maintainers: </span>
        <PersonChip v-for="person in item.maintainers" :person size="small" />
      </div>
    </template> -->
  </CRUDTable>
</template>

<script setup lang="ts">
import { Dataset } from '@/api'
import { listDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { LineClampedText } from '@/components/toolkit/ui/LineClampedText'
import { useQuery } from '@tanstack/vue-query'

const { data: items, refetch, error, isPending: loading } = useQuery(listDatasetsOptions())

const headers: CRUDTableHeader<Dataset>[] = [
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
