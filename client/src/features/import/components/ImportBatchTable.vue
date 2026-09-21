<template>
  <v-data-table :items="data" :headers="headers" v-bind="$attrs">
    <template #item.label="{ value, item }">
      <router-link :to="{ name: 'import-batch-item', params: { uuid: item.id } }">
        {{ value }}
      </router-link>
    </template>
    <template #item.completed_at="{ value }">
      <span class="text-muted">
        {{ DateTime.fromJSDate(value).toLocaleString(DateTime.DATETIME_SHORT) }}
      </span>
    </template>
    <template #item.created_by_user="{ value }">
      <UserChip :user="value" size="small" />
    </template>
    <template #item.completed_by_user="{ value }">
      <UserChip :user="value" size="small" />
    </template>
  </v-data-table>
</template>

<script setup lang="ts">
import { ImportBatchListItem, ImportBatchWithDetails } from '@/api'
import UserChip from '@/features/users/components/UserChip'
import { DateTime } from 'luxon'

const { data } = defineProps<{
  data?: ImportBatchListItem[]
}>()

const headers: DataTableHeader[] = [
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
    key: 'created_by_user',
    title: 'Submitted by'
  },
  {
    key: 'completed_at',
    title: 'Imported At'
  },
  {
    key: 'completed_by_user',
    title: 'Created by'
  }
] satisfies CRUDTableHeader<ImportBatchWithDetails>[]
</script>

<style scoped lang="scss"></style>
