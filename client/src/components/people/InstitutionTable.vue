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

    <!-- Header slots -->
    <template v-slot:[`header.people`]>
      <v-icon title="People">mdi-account-group </v-icon>
    </template>

    <!-- Item slots -->
    <template v-slot:[`item.code`]="{ item }">
      <code>{{ item.code }}</code>
    </template>
    <template v-slot:[`item.kind`]="{ item }">
      <v-chip
        label
        variant="outlined"
        color="primary"
        text-color="primary"
        v-bind="kindIcon(item.kind)"
        :title="enumAsString(item.kind)"
      >
        <template v-slot:prepend>
          <v-icon v-bind="kindIcon(item.kind)"></v-icon>
        </template>
        {{ mdAndUp ? enumAsString(item.kind) : '' }}
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
        <div v-if="!mdAndUp">
          <v-card density="compact" flat>
            <v-card-title class="text-body-2">Description</v-card-title>
            <v-card-text class="text-caption">{{ item.description }}</v-card-text>
          </v-card>
          <v-divider />
        </div>
        <v-card flat :min-width="300">
          <v-list lines="one" density="compact" prepend-icon="mdi-account">
            <v-list-subheader>
              {{ item.people?.length ? 'PEOPLE' : 'No people registered in this institution.' }}
            </v-list-subheader>
            <v-list-item v-for="person in item.people" :key="person.id" class="item-person">
              <v-list-item-title class="text-body-2">
                {{ person.full_name }}
              </v-list-item-title>
              <template v-slot:prepend>
                <v-icon v-bind="roleIcon(person.role)" size="small"></v-icon>
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
import { computed } from 'vue'
import { useDisplay } from 'vuetify'
import { enumAsString } from '../toolkit/enums'
import { useEntityTable } from '../toolkit/tables'
import { kindIcon } from './institutionKind'
import { roleIcon } from './userRole'

const { mdAndUp } = useDisplay()

const headers = computed((): ReadonlyHeaders => {
  const headers: DataTableHeader[] = [
    { title: 'Short name', key: 'code' },
    { title: 'Name', key: 'name' },
    { title: 'Kind', key: 'kind' }
  ]
  if (mdAndUp.value) {
    headers.push({ title: 'Description', key: 'description', sortable: false })
  }
  return headers.concat([
    { title: 'People', key: 'people', align: 'center' },
    { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
  ])
})

const { create, edit, editItem, onFormSuccess } = useEntityTable<Institution>()
</script>

<style scoped>
.item-person {
  min-height: unset;
}
</style>
