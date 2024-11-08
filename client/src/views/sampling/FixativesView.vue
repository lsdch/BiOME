<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Fixative"
    :toolbar="{ title: 'Fixatives', icon: 'mdi-snowflake' }"
    :fetch-items="SamplingService.listFixatives"
    show-actions
  />
</template>

<script setup lang="ts">
import { SamplingService } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' }
]
</script>

<style scoped></style>
