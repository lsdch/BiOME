<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Abiotic parameter"
    :toolbar="{ title: 'Abiotic parameters', icon: 'mdi-gauge' }"
    :fetch-items="listAbioticParametersOptions"
    :delete="
      ({ code }: AbioticParameter) => SamplingService.deleteAbioticParameter({ path: { code } })
    "
    appendActions
  >
    <template #[`item.unit`]="{ value }">
      <code>{{ value }}</code>
    </template>
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <AbioticParameterFormDialog
        :model-value="dialog"
        @close="onClose"
        @success="onSuccess"
        :edit="editItem"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { AbioticParameter, SamplingService } from '@/api'
import { listAbioticParametersOptions } from '@/api/gen/@tanstack/vue-query.gen'
import AbioticParameterFormDialog from '@/components/events/AbioticParameterFormDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<AbioticParameter>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' },
  { key: 'unit', title: 'Unit' }
]
</script>

<style scoped></style>
