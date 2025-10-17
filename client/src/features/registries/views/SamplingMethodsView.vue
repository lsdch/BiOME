<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Sampling method"
    :toolbar="{ title: 'Sampling methods', icon: 'mdi-hook' }"
    :fetch-items="listSamplingMethodsOptions()"
    appendActions
  >
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <SamplingMethodFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :model-value="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SamplingMethod } from '@/api'
import { listSamplingMethodsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SamplingMethodFormDialogMutation from '@/components/forms/SamplingMethodFormDialogMutation.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<SamplingMethod>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' }
]
</script>

<style scoped></style>
