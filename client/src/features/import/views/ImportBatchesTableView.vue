<template>
  <CRUDTable
    class="fill-height"
    :items="data"
    :headers
    entity-name="Import batch"
    :toolbar="{
      title: 'Import batches',
      icon: 'mdi-file-table-outline'
    }"
  >
    <template #item.label="{ value, item }">
      <router-link :to="{ name: 'import-batch-item', params: { uuid: item.id } }">
        {{ value }}
      </router-link>
    </template>
    <template #item.created_at="{ value }">
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
      <DeleteBtn
        title="Confirm deletion"
        @confirm="deleteBatch({ path: { id: item.id } }).then(() => refetch())"
      >
        <template #message>
          <v-card-text>
            <b> Are you sure you want to delete this import batch ? </b>
            <ul>
              <li>
                All {{ item.occurrence_count }} occurrences created by this import batch will be
                deleted.
              </li>
              <li>
                All samplings create by this import batch, for which there are no remaining
                occurrences, will also be deleted.
              </li>
              <li>Taxa created by this import batch will not be deleted.</li>
            </ul>
            <span class="text-error">This action cannot be undone.</span>
          </v-card-text>
        </template>
      </DeleteBtn>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { ImportBatchWithContent } from '@/api'
import {
  deleteImportBatchMutation,
  listImportBatchesWithContentOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import UserChip from '@/features/users/components/UserChip'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { DateTime } from 'luxon'

const { data, isPending, error, refetch } = useQuery(listImportBatchesWithContentOptions())

const { mutateAsync: deleteBatch } = useMutation(deleteImportBatchMutation())

const headers: CRUDTableHeader<ImportBatchWithContent>[] = [
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
]
</script>

<style scoped lang="scss"></style>
