<template>
  <ImportBatchTable
    :batches="data ?? []"
    :toolbar="{
      title: 'Import batches',
      icon: 'mdi-file-table-outline'
    }"
  >
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
  </ImportBatchTable>
</template>

<script setup lang="ts">
import {
  deleteImportBatchMutation,
  listImportBatchesWithContentOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import { useMutation, useQuery } from '@tanstack/vue-query'
import ImportBatchTable from '../components/ImportBatchTable.vue'

const { data, isPending, error, refetch } = useQuery(listImportBatchesWithContentOptions())

const { mutateAsync: deleteBatch } = useMutation(deleteImportBatchMutation())
</script>

<style scoped lang="scss"></style>
