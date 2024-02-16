<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :crud="{
      list: PeopleService.getPeoplePersons,
      delete: (person: Person) => PeopleService.deletePerson(person.id)
    }"
    entityName="Person"
    :itemRepr="(p: Person) => p.full_name"
    :toolbar="{
      title: 'People',
      icon: 'mdi-account'
    }"
    show-actions
    @create-item="create"
    @edit-item="edit"
  >
    <template v-slot:form>
      <PersonForm :edit="editItem" @success="onFormSuccess"></PersonForm>
    </template>
    <template v-slot:[`item.institutions`]="{ value }">
      <v-chip v-for="inst in value" :key="inst.code" class="text-overline mx-1" rounded="xl">
        {{ inst.code }}
      </v-chip>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Person, PeopleService } from '@/api'

import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import PersonForm from './PersonForm.vue'
import { useEntityTable } from '../toolkit/tables'

const headers: ReadonlyHeaders = [
  { title: 'Name', key: 'full_name' },
  { title: 'Institutions', key: 'institutions' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
]

const { create, edit, editItem, onFormSuccess } = useEntityTable<Person>()
</script>

<style scoped></style>
