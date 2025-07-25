<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :fetch-items="listOrganisationsOptions()"
    :delete="{
      mutation: deleteOrganisationMutation,
      params: (inst: Organisation) => ({ path: { code: inst.code } })
    }"
    entityName="Organisation"
    :itemRepr="(inst) => inst.code"
    :toolbar="{
      title: 'Organisations',
      icon: 'mdi-domain'
    }"
    v-model:search="filters"
    :filter="filter"
    filter-mode="some"
    appendActions
    :filter-keys="['code', 'name', 'kind']"
  >
    <template #menu>
      <v-row class="ma-0">
        <v-col cols="12" md="6" class="pa-0">
          <v-list-item>
            <OrgKindPicker
              v-model="filters.kind"
              class="mt-1 mb-2"
              label="Kind"
              placeholder="Any"
              density="compact"
              hide-details
              clearable
            />
          </v-list-item>
        </v-col>
      </v-row>
    </template>
    <template #form="{ dialog, editItem, onSuccess, onClose }">
      <OrganisationFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :model-value="editItem"
        @success="onSuccess"
        @close="onClose"
      />
    </template>

    <!-- Header slots -->
    <template #[`header.people`]>
      <v-icon title="People">mdi-account-group </v-icon>
    </template>

    <!-- Item slots -->
    <template #[`item.code`]="{ item }">
      <code>{{ item.code }}</code>
    </template>
    <template #[`item.kind`]="{ item, value }">
      <OrgKindChip size="small" :kind="item.kind" :label="value" :hide-label="!mdAndUp" />
    </template>
    <template #[`item.people`]="{ value, toggleExpand, internalItem }">
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
    <template #expanded-row-inject="{ item }">
      <div class="w-100">
        <v-card density="compact" flat>
          <v-card-title class="text-body-2">Description</v-card-title>
          <v-card-text class="text-caption">{{ item.description ?? 'None provided' }}</v-card-text>
        </v-card>
        <v-divider />
        <v-card flat :min-width="300">
          <v-list lines="one" density="compact" prepend-icon="mdi-account">
            <v-list-subheader>
              {{ item.people?.length ? 'PEOPLE' : 'No people registered in this organisation.' }}
            </v-list-subheader>
            <PersonChip
              v-for="person in item.people"
              :key="person.id"
              :person="person"
              size="small"
              class="ma-1"
            />
          </v-list>
        </v-card>
      </div>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Organisation, OrgKind } from '@/api'
import {
  deleteOrganisationMutation,
  listOrganisationsOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'
import OrganisationFormDialogMutation from '../forms/people/OrganisationFormDialogMutation.vue'
import { enumAsString } from '../toolkit/enums'
import OrgKindChip from './OrgKindChip'
import OrgKindPicker from './OrgKindPicker.vue'
import PersonChip from './PersonChip'

const { mdAndUp } = useDisplay()

const filters = ref<{ term?: string; kind?: OrgKind }>({})

const filter = computed(() => {
  return filters.value.kind ? (item: Organisation) => item.kind === filters.value.kind : () => true
})

const headers = computed((): CRUDTableHeader<Organisation>[] => [
  { title: 'Short name', key: 'code' },
  { title: 'Name', key: 'name' },
  {
    title: 'Kind',
    key: 'kind',
    value: (item: Organisation) => enumAsString(item.kind)
  },
  { title: 'People', key: 'people', align: 'center' }
])
</script>

<style scoped>
.item-person {
  min-height: unset;
}
</style>
