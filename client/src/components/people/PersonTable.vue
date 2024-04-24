<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :crud="{
      list: PeopleService.listPersons,
      delete: (person: Person) => PeopleService.deletePerson({ id: person.id })
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
    <!-- User Role column -->
    <template #[`header.role`]="slotProps">
      <IconTableHeader v-bind="slotProps" icon="mdi-account-badge" />
    </template>
    <template #[`item.role`]="{ value }">
      <v-icon v-bind="roleIcon(value)" size="x-small" :title="value" />
    </template>

    <template v-slot:search>
      <PersonFilters v-model="filters" />
    </template>
    <template v-slot:form>
      <PersonForm :edit="editItem" @success="onFormSuccess"></PersonForm>
    </template>

    <template v-slot:[`item.alias`]="{ value }">
      <span class="font-weight-light"> {{ `@${value}` }}</span>
    </template>
    <template v-slot:[`item.institutions`]="{ value }">
      <InstitutionKindChip
        v-for="inst in value"
        :key="inst.code"
        :kind="inst.kind"
        :label="inst.code"
        :hide-label="xs"
        size="x-small"
      />
    </template>

    <template v-slot:[`expanded-row-inject`]="{ item }">
      <v-card v-if="item.comment" flat>
        <v-list-item prepend-icon="mdi-comment">
          {{ item.comment }}
        </v-list-item>
      </v-card>
      <v-card flat>
        <v-btn
          v-if="item.contact"
          variant="plain"
          prepend-icon="mdi-at"
          :href="`mailto:${item.contact}`"
          :text="item.contact"
          size="small"
          class="mx-1"
        />
        <v-btn
          color="primary"
          variant="plain"
          prepend-icon="mdi-account-box-plus-outline"
          text="Invite"
        />
      </v-card>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Institution, PeopleService, Person } from '@/api'

import { UserRole } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'
import { useEntityTable } from '../toolkit/tables'
import InstitutionKindChip from './InstitutionKindChip.vue'
import type { PersonFilters as Filters } from './PersonFilters.vue'
import PersonFilters from './PersonFilters.vue'
import PersonForm from './PersonForm.vue'
import { orderedUserRoles, roleIcon } from './userRole'
import IconTableHeader from '@/components/toolkit/tables/IconTableHeader.vue'

const { xs } = useDisplay()

const filters = ref<Filters>({
  term: '',
  status: undefined
})

const filter = computed(() => {
  const { status } = filters.value
  switch (status) {
    case undefined:
    case null:
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
