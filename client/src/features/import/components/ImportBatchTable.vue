<template>
  <CRUDTable class="fill-height" :items="batches" :headers entity-name="Import batch">
    <template #item.label="{ value, item }">
      <router-link :to="{ name: 'import-batch-item', params: { uuid: item.id } }">
        {{ value }}
      </router-link>
    </template>
    <template #item.completed_at="{ value }">
      <span class="text-muted">{{
        DateTime.fromJSDate(value).toLocaleString(DateTime.DATETIME_SHORT)
      }}</span>
    </template>
    <template #item.created_by_user="{ value }">
      <UserChip :user="value" size="small" />
    </template>
    <template #item.completed_by_user="{ value }">
      <UserChip :user="value" size="small" />
    </template>

    <template #expanded-row-footer-append="{ item }">
      <ImportBatchCreationChip :batch="item" />
      <slot name="expanded-row-footer-append" :item="item" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { ImportBatchListItem, ImportBatchWithDetails } from '@/api'
import {
  deleteImportBatchMutation,
  listImportBatchesWithContentOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import UserChip from '@/features/users/components/UserChip'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { DateTime } from 'luxon'
import ImportBatchCreationChip from '../components/ImportBatchCreationChip.vue'

const { batches } = defineProps<{
  batches: ImportBatchListItem[]
}>()

const { mutateAsync: deleteBatch } = useMutation(deleteImportBatchMutation())

const headers: CRUDTableHeader<ImportBatchListItem>[] = [
  {
    key: 'label',
    title: 'Label'
  },
  {
    key: 'occurrence_count',
    title: 'Occurrences',
    width: 0,
    align: 'end',
    cellProps: { class: 'text-overline' }
  },
  {
    key: 'sampling_count',
    title: 'Samplings',
    width: 0,
    align: 'end',
    cellProps: { class: 'text-overline' }
  },
  {
    key: 'completed_at',
    title: 'Imported At'
  }
]
</script>

<style scoped lang="scss"></style>
