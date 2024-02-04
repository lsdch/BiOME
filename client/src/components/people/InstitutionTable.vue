<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :crud="{
      list: PeopleService.getPeopleInstitutions,
      create: PeopleService.createInstitution,
      update: PeopleService.updateInstitution,
      delete: (inst: Institution) => PeopleService.deleteInstitution(inst.acronym)
    }"
    :toolbar-props="{
      title: 'Institutions',
      entityName: 'Institution',
      form: InstitutionForm,
      icon: 'mdi-domain',
      itemRepr: (inst) => inst.acronym
    }"
    show-actions
    icon="mdi-domain"
    :filter-keys="['acronym', 'name']"
  >
    <template v-slot:[`item.acronym`]="{ item }">
      <code>{{ item.acronym }}</code>
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
      <!-- <tr class="pb-5">
        <td :colspan="columns.length">
          <div class="d-flex">
            <div class="flex-grow-0 mr-5">
              <ItemMetaInfos v-if="item.meta" :meta="item.meta" />
            </div>
            <v-divider vertical></v-divider -->
      <v-container>
        <v-row>
          <v-col>
            <v-list v-if="item.people?.length" lines="one">
              <v-list-item
                v-for="person in item.people"
                :key="person.id"
                :title="person.full_name"
              />
            </v-list>
            <span v-else>No people registered in this institution.</span>
          </v-col>
        </v-row>
      </v-container>
      <!-- </div>
        </td>
      </tr> -->
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
  { title: 'Description', key: 'description', sortable: false },
  { title: 'People', key: 'people', align: 'center' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
]
</script>

<style scoped></style>
