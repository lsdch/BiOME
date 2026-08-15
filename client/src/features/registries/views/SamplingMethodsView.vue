<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Sampling method"
    :toolbar="{ title: 'Sampling methods', icon: 'mdi-hook' }"
    :items
    :error
    :loading
    @reload="refetch()"
    appendActions
  >
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <SamplingMethodFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :item="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template>
    <template #expanded-row-footer-append="{ item }">
      <DeleteBtn
        title="Confirm deletion"
        @confirm="deleteSamplingMethod({ path: { code: item.code } })"
      >
        <template #message>
          <v-card-text>
            <b> Are you sure you want to delete this sampling method ? </b>
            <p>Samplings referencing this method will be preserved.</p>
          </v-card-text>
        </template>
      </DeleteBtn>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SamplingMethod } from '@/api'
import {
  deleteSamplingMethodMutation,
  listSamplingMethodsOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import SamplingMethodFormDialogMutation from '@/components/forms/SamplingMethodFormDialogMutation.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import { useMutation, useQuery } from '@tanstack/vue-query'

const { data: items, refetch, error, isPending: loading } = useQuery(listSamplingMethodsOptions())

const { mutateAsync: deleteMutation } = useMutation(deleteSamplingMethodMutation())
async function deleteSamplingMethod(params: { path: { code: string } }) {
  await deleteMutation(params, { onSuccess: () => refetch() })
}

const headers: CRUDTableHeader<SamplingMethod>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'name', title: 'Name' }
]
</script>

<style scoped></style>
