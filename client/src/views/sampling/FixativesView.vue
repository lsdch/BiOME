<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Fixative"
    :toolbar="{ title: 'Fixatives', icon: 'mdi-snowflake' }"
    :fetch-items="listFixativesOptions"
    :delete="{
      mutation: deleteFixativeMutation,
      params: ({ code }: Fixative) => ({ path: { code } })
    }"
    appendActions
  >
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <FixativeFormDialog :dialog :model-value="editItem" @close="onClose" @success="onSuccess" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Fixative, SamplingService } from '@/api'
import { deleteFixativeMutation, listFixativesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import FixativeFormDialog from '@/components/samples/FixativeFormDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<Fixative>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'label', title: 'Label' }
]
</script>

<style scoped></style>
