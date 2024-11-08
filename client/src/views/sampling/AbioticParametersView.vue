<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Abiotic parameter"
    :toolbar="{ title: 'Abiotic parameters', icon: 'mdi-gauge' }"
    :fetch-items="SamplingService.listAbioticParameters"
    show-actions
  >
    <template #[`item.unit`]="{ value }">
      <code>{{ value }}</code>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SamplingService } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' },
  { key: 'unit', title: 'Unit' }
]
</script>

<style scoped></style>
