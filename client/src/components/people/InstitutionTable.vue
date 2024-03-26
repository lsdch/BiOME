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
    :search="filters.term"
    :filter="filter"
    filter-mode="some"
    show-actions
    :filter-keys="['code', 'name', 'kind']"
    @create-item="create"
    @edit-item="edit"
  >
    <template v-slot:search>
      <InstitutionFilters v-model="filters" />
    </template>
    <template v-slot:form>
      <InstitutionForm :edit="editItem" @success="onFormSuccess" />
    </template>

    <!-- Header slots -->
    <template v-slot:[`header.people`]>
      <v-icon title="People">mdi-account-group </v-icon>
    </template>

    <!-- Item slots -->
    <template v-slot:[`item.code`]="{ item }">
      <code>{{ item.code }}</code>
    </template>
    <template v-slot:[`item.kind`]="{ item, value }">
      <v-chip
        label
        variant="outlined"
        color="primary"
        text-color="primary"
        v-bind="kindIcon(item.kind)"
      >
        <template v-slot:prepend>
          <v-icon v-bind="kindIcon(item.kind)"></v-icon>
        </template>
        {{ mdAndUp ? value : '' }}
      </v-chip>
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
      <div class="w-100">
        <v-card density="compact" flat>
          <v-card-title class="text-body-2">Description</v-card-title>
          <v-card-text class="text-caption">{{ item.description }}</v-card-text>
        </v-card>
        <v-divider />
        <v-card flat :min-width="300">
          <v-list lines="one" density="compact" prepend-icon="mdi-account">
            <v-list-subheader>
              {{ item.people?.length ? 'PEOPLE' : 'No people registered in this institution.' }}
            </v-list-subheader>
            <v-list-item v-for="person in item.people" :key="person.id" class="item-person py-0">
              <v-list-item-title class="text-body-2">
                {{ person.full_name }}
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ `@${person.alias}` }}
              </v-list-item-subtitle>
              <template v-slot:prepend>
                <v-icon v-bind="roleIcon(person.role)" size="small" />
              </template>
            </v-list-item>
          </v-list>
        </v-card>
      </div>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Institution, PeopleService } from '@/api'
import InstitutionForm from './InstitutionForm.vue'

import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'
import { enumAsString } from '../toolkit/enums'
import { useEntityTable } from '../toolkit/tables'
import InstitutionFilters from './InstitutionFilters.vue'
import { Filters } from './InstitutionFilters.vue'
import { kindIcon } from './institutionKind'
import { roleIcon } from './userRole'

const { mdAndUp } = useDisplay()

const filters = ref<Filters>({
  term: '',
  kind: undefined
})

const filter = computed(() => {
  return filters.value.kind ? (item: Institution) => item.kind === filters.value.kind : () => true
})

const headers = computed(
  (): CRUDTableHeaders => [
    { title: 'Short name', key: 'code' },
    { title: 'Name', key: 'name' },
    {
      title: 'Kind',
      key: 'kind',
      value(item: Institution) {
        return enumAsString(item.kind)
      }
    },
    { title: 'People', key: 'people', align: 'center' }
  ]
)

const { create, edit, editItem, onFormSuccess } = useEntityTable<Institution>()
</script>

<style scoped>
.item-person {
  min-height: unset;
}
</style>
