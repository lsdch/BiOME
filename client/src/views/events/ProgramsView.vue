<template>
  <CRUDTable
    class="fill-height"
    entity-name="Program"
    :headers
    :toolbar="{
      title: 'Programs',
      icon: 'mdi-notebook'
    }"
    :fetch-items="listProgramsOptions"
    :delete="{
      mutation: deleteProgramMutation,
      params: ({ code }: Program) => ({ path: { code } })
    }"
    appendActions
  >
    <template #[`item.funding_agencies`]="{ value }">
      <OrgKindChip
        v-for="inst in value"
        :key="inst.code"
        :kind="inst.kind"
        :label="inst.code"
        size="small"
      />
    </template>
    <template #expanded-row-inject="{ item }">
      <div class="pa-3">
        <p class="mb-2">{{ item.description }}</p>
        <p>
          Managers: <PersonChip v-for="person in item.managers" :person size="small" class="ma-1" />
        </p>
      </div>
    </template>

    <!-- Form dialog -->
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <ProgramFormDialog
        :dialog
        @close="onClose"
        @success="onSuccess"
        :mode
        :fullscreen="xs"
        :model-value="editItem"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Program } from '@/api'
import { deleteProgramMutation, listProgramsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import ProgramFormDialog from '@/components/events/ProgramFormDialog.vue'
import OrgKindChip from '@/components/people/OrgKindChip.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { useDisplay } from 'vuetify'

const { xs } = useDisplay()

const headers: CRUDTableHeader<Program>[] = [
  { key: 'label', title: 'Label' },
  { key: 'code', title: 'Code', cellProps: { class: 'text-overline' } },
  { key: 'funding_agencies', title: 'Funding agencies' },
  { key: 'start_year', title: 'Start', width: 0 },
  { key: 'end_year', title: 'End', width: 0 }
]
</script>

<style scoped></style>
