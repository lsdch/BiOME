<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="SequencesService.listSeqDbs"
    entity-name="Seq. DB"
    :toolbar="{
      title: 'Sequence databases',
      icon: 'mdi-database-sync'
    }"
    append-actions
  >
    <template #expanded-row-inject="{ item }">
      <v-list>
        <v-list-item
          title="Description"
          :subtitle="item.description ?? 'No description provided'"
        ></v-list-item>
        <v-list-item title="Link template ">
          <template #subtitle>
            <code v-if="item.link_template">{{ item.link_template }}</code>
            <div v-else class="text-wrap overflow-auto">
              No template provided: direct link generation using accession numbers won't be
              available.
            </div>
          </template>
        </v-list-item>
      </v-list>
    </template>
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <SeqDBFormDialog
        :model-value="dialog"
        @close="onClose"
        @success="onSuccess"
        :edit="editItem"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SeqDb, SequencesService } from '@/api'
import SeqDBFormDialog from '@/components/sequences/SeqDBFormDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<SeqDb>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } },
  { key: 'label', title: 'Label' }
]
</script>

<style scoped lang="scss"></style>
