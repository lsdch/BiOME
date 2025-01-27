<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Fixative"
    :toolbar="{ title: 'Fixatives', icon: 'mdi-snowflake' }"
    :fetch-items="listFixativesOptions"
    :delete="({ code }: Fixative) => SamplingService.deleteFixative({ path: { code } })"
    appendActions
  >
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <FixativeFormDialog
        :model-value="dialog"
        @close="onClose"
        @success="onSuccess"
        :edit="editItem"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Fixative, SamplingService } from '@/api'
import { listFixativesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import FixativeFormDialog from '@/components/samples/FixativeFormDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<Fixative>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' }
]
</script>

<style scoped></style>
