<template>
  <CRUDTable
    :headers="headers"
    density="compact"
    :fetch-items="PeopleService.listInstitutions"
    :delete="(inst: Institution) => PeopleService.deleteInstitution({ path: { code: inst.code } })"
    entityName="Institution"
    :itemRepr="(inst) => inst.code"
    :toolbar="{
      title: 'Institutions',
      icon: 'mdi-domain'
    }"
    :search="filters.term"
    :filter="filter"
    filter-mode="some"
    appendActions
    :filter-keys="['code', 'name', 'kind']"
  >
    <template #search>
      <InstitutionFilters v-model="filters" />
    </template>
    <template #form="{ dialog, editItem, onSuccess, onClose }">
      <InstitutionFormDialog
        :model-value="dialog"
        :edit="editItem"
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
      <InstitutionKindChip size="small" :kind="item.kind" :label="value" :hide-label="!mdAndUp" />
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
              <template #prepend>
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
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'
import { enumAsString } from '../toolkit/enums'
import InstitutionFilters, { Filters } from './InstitutionFilters.vue'
import InstitutionKindChip from './InstitutionKindChip.vue'
import { roleIcon } from './userRole'
import InstitutionFormDialog from './InstitutionFormDialog.vue'

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
      value: (item: Institution) => enumAsString(item.kind)
    },
    { title: 'People', key: 'people', align: 'center' }
  ]
)
</script>

<style scoped>
.item-person {
  min-height: unset;
}
</style>
