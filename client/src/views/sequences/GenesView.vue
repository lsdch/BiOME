<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Gene"
    :toolbar="{ title: 'Genes registry', icon: 'mdi-tag' }"
    :fetch-items="listGenesOptions()"
    append-actions
  >
    <template #item.is_MOTU_delimiter="{ value }">
      <v-icon v-if="value" icon="mdi-check" color="success" />
    </template>

    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <GeneFormDialogMutation
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
import { Gene } from '@/api'
import { listGenesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import GeneFormDialogMutation from '@/components/forms/GeneFormDialogMutation.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<Gene>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' },
  { key: 'is_MOTU_delimiter', title: 'MOTU', width: 0, align: 'center' }
]
</script>

<style scoped></style>
