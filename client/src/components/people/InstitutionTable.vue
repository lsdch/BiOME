<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :crud="{
      list: PeopleService.listInstitutions,
      delete: (inst: Institution) => PeopleService.deleteInstitution(inst.code)
    }"
    entityName="Institution"
    :itemRepr="(inst) => inst.code"
    :toolbar="{
      title: 'Institutions',
      icon: 'mdi-domain'
    }"
    show-actions
    :filter-keys="['code', 'name']"
    @create-item="create"
    @edit-item="edit"
  >
    <template v-slot:form>
      <InstitutionForm :edit="editItem" @success="onFormSuccess" />
    </template>

    <template v-slot:[`item.code`]="{ item }">
      <code>{{ item.code }}</code>
    </template>
    <template v-slot:[`item.people`]="{ value, toggleExpand, internalItem }">
      <v-btn
        icon
        color="primary"
        variant="tonal"
        density="compact"
        @click="toggleExpand(internalItem)"
      >
        {{ value?.length ?? 0 }}
      </v-btn>
    </template>
    <template v-slot:expanded-row-inject="{ item }">
      <v-card flat :min-width="300">
        <v-list lines="one" density="compact" prepend-icon="mdi-account">
          <v-list-subheader>
            {{ item.people?.length ? 'PEOPLE' : 'No people registered in this institution.' }}
          </v-list-subheader>
          <v-list-item
            v-for="person in item.people"
            :key="person.id"
            :subtitle="person.full_name"
          />
        </v-list>
      </v-card>
      <v-divider vertical />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Institution, PeopleService } from '@/api'
import InstitutionForm from './InstitutionForm.vue'

import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { useEntityTable } from '../toolkit/tables'

const headers: ReadonlyHeaders = [
  { title: 'Short name', key: 'code' },
  { title: 'Name', key: 'name' },
  { title: 'Description', key: 'description', sortable: false },
  { title: 'People', key: 'people', align: 'center' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
]

const { create, edit, editItem, onFormSuccess } = useEntityTable<Institution>()
</script>

<style scoped></style>
