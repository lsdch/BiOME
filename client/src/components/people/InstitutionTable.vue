<template>
  <CRUDTable
    title="Institutions"
    :headers="headers"
    density="compact"
    :list="PeopleService.getPeopleInstitutions"
    :create="PeopleService.createInstitution"
    :update="PeopleService.updateInstitution"
    :delete="(inst: Institution) => PeopleService.deleteInstitution(inst.acronym)"
    :form="InstitutionForm"
    show-actions
    :item-repr="(inst) => inst.acronym"
    icon="mdi-domain"
  >
    <template v-slot:[`item.acronym`]="{ item }">
      <code>{{ item.acronym }}</code>
    </template>
    <template v-slot:[`item.people`]="{ item }">
      <v-badge color="primary" :content="item.people?.length ?? 0" inline></v-badge>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Institution, PeopleService } from '@/api'
import InstitutionForm from './InstitutionForm.vue'

import CRUDTable from '@/components/toolkit/CRUDTable.vue'

const headers: ReadonlyHeaders = [
  { title: 'Short name', key: 'acronym' },
  { title: 'Name', key: 'name' },
  { title: 'Description', key: 'description' },
  { title: 'People', key: 'people' },
  { title: 'Actions', key: 'actions' }
]
</script>

<style scoped></style>
