<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Abiotic parameter"
    :toolbar="{ title: 'Abiotic parameters', icon: 'mdi-gauge' }"
    :items
    :error
    :loading
    @reload="refetch()"
    appendActions
  >
    <!-- :delete="{
      mutation: deleteAbioticParameterMutation,
      params: ({ code }: AbioticParameter) => ({ path: { code } })
    }" -->
    <template #[`item.unit`]="{ value }">
      <code>{{ value }}</code>
    </template>
    <!-- <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <AbioticParameterFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :item="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template> -->
  </CRUDTable>
</template>

<script setup lang="ts">
import { AbioticParam } from '@/api'
import {
  // deleteAbioticParameterMutation,
  listAbioticParametersOptions
} from '@/api/gen/@tanstack/vue-query.gen'
// import AbioticParameterFormDialogMutation from '@/components/forms/AbioticParamFormDialogMutation.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { useQuery } from '@tanstack/vue-query'

const { data: items, refetch, isPending: loading, error } = useQuery(listAbioticParametersOptions())

const headers: CRUDTableHeader<AbioticParam>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' },
  { key: 'unit', title: 'Unit' }
]
</script>

<style scoped></style>
