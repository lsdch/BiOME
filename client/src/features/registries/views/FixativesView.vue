<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Fixative"
    :toolbar="{ title: 'Fixatives', icon: 'mdi-snowflake' }"
    :items
    :error
    :loading
    @reload="refetch()"
    appendActions
  >
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <FixativeFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :item="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template>

    <template #expanded-row-footer-append="{ item }">
      <DeleteBtn
        title="Delete fixative ?"
        message="Deleted fixative will be removed for all associated samplings"
        @confirm="deleteMutation({ path: { code: item.code } }).then(() => refetch())"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Fixative } from '@/api'
import { deleteFixativeMutation, listFixativesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import FixativeFormDialogMutation from '@/components/forms/FixativeFormDialogMutation.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import { useMutation, useQuery } from '@tanstack/vue-query'

const headers: CRUDTableHeader<Fixative>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'name', title: 'Name' }
]

const { data: items, refetch, error, isPending: loading } = useQuery(listFixativesOptions())
const { mutateAsync: deleteMutation } = useMutation(deleteFixativeMutation())
</script>

<style scoped></style>
