<template>
  <CRUDTable
    entityName="Taxon"
    :fetch-items="() => TaxonomyService.listTaxa()"
    :delete="(item: Taxon) => TaxonomyService.deleteTaxon({ path: { code: item.code } })"
    :itemRepr="(item: Taxon) => item.name"
    :toolbar="{
      title: 'Taxonomy',
      icon: 'mdi-graph',
      togglableSearch: true
    }"
    :headers="headers"
    :filter="filter"
    showActions
    density="compact"
    v-model:search="searchName"
    :filter-keys="['name']"
    reload-on-delete
    fixed-header
    fixed-footer
    height="100"
    :items-per-page="15"
  >
    <template v-slot:form="{ dialog, onClose, onSuccess, editItem }">
      <TaxonFormDialog
        :model-value="dialog"
        :edit="editItem"
        @success="onSuccess"
        @onClose="onClose"
      />
    </template>
    <template v-slot:search>
      <TaxaTableFilters v-model="filters" v-model:name="searchName" />
    </template>
    <template v-slot:[`item.code`]="{ item }">
      <code>{{ item.code }}</code>
    </template>
    <template v-slot:[`item.status`]="{ item }">
      <StatusIcon :status="item.status" size="small" />
      <LinkIconGBIF v-if="item.GBIF_ID" :GBIF_ID="item.GBIF_ID" variant="text" size="x-small" />
    </template>
    <template v-slot:[`expanded-row-inject`]="{ item }">
      <v-card flat v-if="item.authorship">
        <v-card-title class="text-body-2">
          <v-icon size="small">mdi-newspaper-variant-outline</v-icon>
          {{ item.authorship }}
        </v-card-title>
      </v-card>
      <v-card flat v-if="item.comment">
        <v-card-title class="text-body-2">
          <v-icon size="small">mdi-comment-processing</v-icon>
          {{ item.comment }}
        </v-card-title>
      </v-card>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { Taxon, TaxonRank, TaxonStatus, TaxonomyService } from '@/api'
import { Ref, computed } from 'vue'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'
import LinkIconGBIF from './LinkIconGBIF.vue'
import StatusIcon from './StatusIcon.vue'
import TaxaTableFilters from './TaxaTableFilters.vue'
import TaxonFormDialog from './TaxonFormDialog.vue'

const searchName = ref('')
const filters: Ref<{ rank?: TaxonRank; status: TaxonStatus }> = ref({
  rank: undefined,
  status: 'Accepted'
})

const filter = computed(() => {
  const { rank, status } = filters.value
  if (rank || status)
    return (item: Taxon) => {
      return (rank ? item.rank === rank : true) && (status ? item.status === status : true)
    }
  else return undefined
})

const headers: CRUDTableHeaders = [
  { title: 'Name', key: 'name' },
  { title: 'Code', key: 'code' },
  { title: 'Rank', key: 'rank' },
  { title: 'Status', key: 'status', width: 0, align: 'center' }
]
</script>

<style scoped></style>
