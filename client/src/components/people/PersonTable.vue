<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :crud="{
      list: PeopleService.getPeoplePersons,
      delete: (person: Person) => PeopleService.deletePerson(person.id)
    }"
    :filter="filter"
    :search="filters.term"
    entityName="Person"
    :itemRepr="(p: Person) => p.full_name"
    :toolbar="{
      title: 'People',
      icon: 'mdi-account'
    }"
    filter-mode="some"
    show-actions
    @create-item="create"
    @edit-item="edit"
  >
    <template v-slot:search>
      <PersonFilters v-model="filters" />
    </template>
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
      <v-card v-if="item.comment" flat>
        <v-list-item prepend-icon="mdi-comment">
          {{ item.comment }}
        </v-list-item>
      </v-card>
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
import { Institution, PeopleService, Person } from '@/api'

import { UserRole } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { computed, ref } from 'vue'
import { useEntityTable } from '../toolkit/tables'
import type { PersonFilters as Filters } from './PersonFilters.vue'
import PersonFilters from './PersonFilters.vue'
import PersonForm from './PersonForm.vue'
import { kindIcon } from './institutionKind'
import { orderedUserRoles, roleIcon } from './userRole'

const filters = ref<Filters>({
  term: '',
  status: undefined
})

const filter = computed(() => {
  const { status } = filters.value
  switch (status) {
    case undefined:
      return () => true
    case 'Registered user':
      return (item: Person) => Boolean(item.role)
    case 'Unregistered':
      return (item: Person) => !item.role
    default:
      return (item: Person) => item.role === status
  }
})

const headers: CRUDTableHeader[] = [
  {
    title: 'Role',
    key: 'role',
    width: 0,
    align: 'end',
    sort(a: UserRole, b: UserRole) {
      if (a === b) return 0
      return orderedUserRoles.indexOf(a) - orderedUserRoles.indexOf(b)
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
    sortable: false,
    filter: (value: Institution[], query: string) => {
      return value.find((inst) => inst.code.includes(query)) !== undefined
    }
  }
]

const { create, edit, editItem, onFormSuccess } = useEntityTable<Person>()
</script>

<style scoped></style>
