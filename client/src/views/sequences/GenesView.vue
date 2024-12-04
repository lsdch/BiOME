<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Gene"
    :toolbar="{ title: 'Genes registry', icon: 'mdi-tag' }"
    :fetch-items="SequencesService.listGenes"
    append-actions
  >
    <template #item.is_MOTU_delimiter="{ value }">
      <v-icon v-if="value" icon="mdi-check" color="success" />
    </template>

    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <GeneFormDialog
        :model-value="dialog"
        @close="onClose"
        @success="onSuccess"
        :edit="editItem"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SequencesService } from '@/api'
import GeneFormDialog from '@/components/sequences/GeneFormDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' },
  { key: 'is_MOTU_delimiter', title: 'MOTU', width: 0, align: 'center' }
]
</script>

<style scoped></style>
