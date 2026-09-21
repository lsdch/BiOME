<template>
  <CRUDTable :headers="headers" :items entityName="Worklow">
    <template #item.batch.label="{ item }">
      <RouterLink :to="{ name: 'import-batch-item', params: { uuid: item.batch.id } }">
        {{ item.batch.label }}
      </RouterLink>
    </template>
    <template #expanded-row-footer-append="{ item }">
      <DeleteBtn
        title="Delete batch workflow ?"
        @confirm="deleteBatch({ path: { id: item.batch.id } }).then(() => refetch())"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { BatchSnapshot } from '@/api'
import {
  deleteBatchWorkflowMutation,
  listImportsForCurrentUserOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import { useMutation, useQuery } from '@tanstack/vue-query'

const { data: items, refetch } = useQuery(listImportsForCurrentUserOptions())

const { mutateAsync: deleteBatch } = useMutation(deleteBatchWorkflowMutation())

const headers: CRUDTableHeader<BatchSnapshot>[] = [
  { key: 'batch.label', title: 'Label' },
  { key: 'status', title: 'Status' },
  { key: 'batch.assembled_by', title: 'Assembled by' },
  { key: 'batch.created_at', title: 'Created at' },
  { key: 'batch.updated_at', title: 'Updated at' }
]
</script>

<style scoped lang="scss"></style>
