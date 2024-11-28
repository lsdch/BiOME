<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Gene"
    :toolbar="{ title: 'Genes registry', icon: 'mdi-tag' }"
    :fetch-items="SequencesService.listGenes"
    appendActions
  />
</template>

<script setup lang="ts">
import { SequencesService } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' }
]
</script>

<style scoped></style>
