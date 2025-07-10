<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Abiotic parameter"
    :toolbar="{ title: 'Abiotic parameters', icon: 'mdi-gauge' }"
    :fetch-items="listAbioticParametersOptions()"
    :delete="{
      mutation: deleteAbioticParameterMutation,
      params: ({ code }: AbioticParameter) => ({ path: { code } })
    }"
    appendActions
  >
    <template #[`item.unit`]="{ value }">
      <code>{{ value }}</code>
    </template>
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <AbioticParameterFormDialog
        :dialog
        :model-value="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { AbioticParameter } from '@/api'
import {
  deleteAbioticParameterMutation,
  listAbioticParametersOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import AbioticParameterFormDialog from '@/components/forms/AbioticParamFormDialogMutation.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<AbioticParameter>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' },
  { key: 'unit', title: 'Unit' }
]
</script>

<style scoped></style>
