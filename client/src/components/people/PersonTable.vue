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
    <template v-slot:[`item.role`]="{ value }">
      <v-icon v-bind="roleIcon(value)"></v-icon>
    </template>
    <template v-slot:[`item.alias`]="{ value }">
      <span class="font-weight-light"> {{ `@${value}` }}</span>
    </template>
    <template v-slot:[`item.institutions`]="{ value }">
      <v-chip
        label
        v-for="inst in value"
        :key="inst.code"
        class="text-overline mx-1"
        variant="outlined"
        v-bind="kindIcon(inst.kind)"
      >
        <template v-slot:prepend>
          <v-icon v-bind="kindIcon(inst.kind)"></v-icon>
        </template>
        {{ inst.code }}
      </v-chip>
    </template>

    <template v-slot:[`expanded-row-inject`]="{ item }">
      <v-btn
        v-if="item.contact"
        variant="plain"
        prepend-icon="mdi-at"
        :href="`mailto:${item.contact}`"
        :text="item.contact"
        size="small"
        class="mx-1"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { PeopleService, Person } from '@/api'

import { UserRole } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { useEntityTable } from '../toolkit/tables'
import PersonForm from './PersonForm.vue'
import { roleIcon } from './userRole'
import { kindIcon } from './institutionKind'

const role_order: UserRole[] = ['Guest', 'Contributor', 'ProjectMember', 'Admin']

const headers: CRUDTableHeader[] = [
  {
    title: 'Role',
    key: 'role',
    width: 0,
    align: 'end',
    sort(a, b) {
      if (a === b) return 0
      return role_order.indexOf(a) - role_order.indexOf(b)
    }
  },
  { title: 'Name', key: 'full_name' },
  {
    title: 'Alias',
    key: 'alias'
  },
  {
    title: 'Institutions',
    key: 'institutions',
    sortable: false
  }
]

const { create, edit, editItem, onFormSuccess } = useEntityTable<Person>()
</script>

<style scoped></style>
