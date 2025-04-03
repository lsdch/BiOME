<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :fetch-items="listPersonsOptions"
    :delete="{
      mutation: deletePersonMutation,
      params: ({ id }) => ({ path: { id } })
    }"
    :filter
    v-model:search="filters"
    entityName="Person"
    :itemRepr="(p: Person) => p.full_name"
    :toolbar="{
      title: 'People',
      icon: 'mdi-account'
    }"
    filter-mode="some"
    appendActions
  >
    <template #menu>
      <PersonFilters v-model="filters" />
    </template>
    <template #form="{ dialog, onClose, onSuccess, editItem }">
      <PersonFormDialogMutation
        @success="onSuccess"
        @close="onClose"
        :dialog
        :model-value="editItem"
      />
    </template>

    <!-- User Role column -->
    <template #[`header.role`]="slotProps">
      <IconTableHeader v-bind="slotProps" icon="mdi-account-badge" />
    </template>
    <template #[`header.organisations`]="slotProps">
      <IconTableHeader v-bind="slotProps" icon="mdi-domain" :expanded="smAndUp" />
    </template>

    <template #[`item.role`]="{ value }">
      <UserRole.Icon :role="value" size="x-small" />
    </template>

    <template #[`item.alias`]="{ value }">
      <span class="font-weight-light"> {{ `@${value}` }}</span>
    </template>
    <template #[`item.organisations`]="{ value }">
      <OrgKindChip
        v-for="inst in value"
        :key="inst.code"
        :kind="inst.kind"
        :label="inst.code"
        :hide-label="xs"
        size="x-small"
      />
    </template>

    <template #[`expanded-row-inject`]="{ item }">
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
import { $UserRole, Organisation, Person } from '@/api'

import { UserRole } from '@/api'
import { deletePersonMutation, listPersonsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import IconTableHeader from '@/components/toolkit/tables/IconTableHeader.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'
import OrgKindChip from './OrgKindChip'
import type { AccountStatus, PersonFilters as Filters } from './PersonFilters.vue'
import PersonFilters from './PersonFilters.vue'
import PersonFormDialogMutation from '../forms/people/PersonFormDialogMutation.vue'

const { xs, smAndUp } = useDisplay()

const filters = ref<Filters>({})

function filterStatus(item: Person, status: AccountStatus) {
  return status === 'Registered user' ? Boolean(item.role) : item.role === status
}

const filter = computed(() => {
  const { status, organisations } = filters.value
  return (item: Person) =>
    (!status || filterStatus(item, status)) &&
    (!organisations || item.organisations.some(({ code }) => organisations.includes(code)))
})

const headers: CRUDTableHeader<Person>[] = [
  {
    title: 'Role',
    key: 'role',
    width: 0,
    align: 'end',
    sort(a: UserRole, b: UserRole) {
      if (a === b) return 0
      return $UserRole.enum.indexOf(a) - $UserRole.enum.indexOf(b)
    }
  },
  { title: 'Name', key: 'full_name' },
  {
    title: 'Alias',
    key: 'alias'
  },
  {
    title: 'Organisations',
    key: 'organisations',
    sortable: false,
    filter: (value: Organisation[], query: string) => {
      return value.find((inst) => inst.code.includes(query)) !== undefined
    }
  }
]
</script>

<style scoped></style>
