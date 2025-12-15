<template>
  <CRUDTableServer
    class="fill-height"
    entity-name="Occurrence"
    :headers
    :filters
    :toolbar="{ title: 'Occurrences', icon: 'mdi-package-variant' }"
    :fetch-items="listOccurrencesOptions"
    :delete="{
      mutation: deleteOccurrenceMutation,
      params: ({ code }: OccurrenceListItem) => ({ path: { code } })
    }"
    :mobile="xs"
    show-expand
    :sort-key-transform
    @clear-filters="filters = {}"
    @reload="invalidateQuery()"
  >
    <!-- Search and filters panel -->
    <template #menu>
      <v-row class="ma-0">
        <v-col cols="12" md="6">
          <v-list>
            <v-list-item prepend-icon="mdi-folder-table">
              <DatasetPicker
                v-model="filters.datasets"
                label="Datasets"
                class="mt-2"
                item-value="label"
                clearable
                multiple
                chips
                closable-chips
                density="compact"
                hide-details
              />
            </v-list-item>
            <v-list-item prepend-icon="mdi-star-four-points">
              <TypeStatusPicker
                v-model="filters.type_status"
                class="mt-2"
                label="Type status"
                multiple
                clearable
                chips
                closable-chips
                density="compact"
                hide-details
              />
            </v-list-item>
            <v-list-item prepend-icon="mdi-dna">
              <ClearableSwitch
                v-model="filters.has_sequences"
                class="pl-2"
                label="Sequences available"
                color-true="primary"
                color-false="red"
                hint="Show only bio material having registered sequences"
                persistent-hint
                density="compact"
              />
            </v-list-item>
            <v-list-item prepend-icon="mdi-calendar">
              <div class="d-flex align-center mt-2">
                <v-number-input
                  v-model="filters.year"
                  label="Year"
                  type="number"
                  :min="0"
                  density="compact"
                  hide-details
                  contentClass="rounded-e-0"
                />
                <v-btn
                  :active="useDateRange"
                  @click="toggleUseDateRange()"
                  color=""
                  active-color="primary"
                  variant="plain"
                  icon="mdi-calendar-expand-horizontal"
                  size="small"
                  class="rounded-s-0"
                  v-tooltip="`Toggle date range`"
                />
                <v-number-input
                  v-if="useDateRange"
                  v-model="filters.year_end"
                  label="End year"
                  type="number"
                  :min="filters.year || 0"
                  density="compact"
                  hide-details
                  contentClass="rounded-se-0"
                  :disabled="filters.year_end === null"
                >
                  <template #append-inner>
                    <v-btn
                      :active="filters.year_end === null"
                      @click="filters.year_end = filters.year_end === null ? undefined : null"
                      icon="mdi-calendar-arrow-right"
                    />
                  </template>
                </v-number-input>
              </div>
            </v-list-item>
            <v-list-item prepend-icon="mdi-calendar">
              <v-slider></v-slider>
            </v-list-item>
          </v-list>
        </v-col>
        <v-col cols="12" md="6">
          <v-list density="compact">
            <v-list-item prepend-icon="mdi-family-tree">
              <TaxonFilterPicker
                v-model:taxa="filters.taxa"
                v-model:whole-clade="filters.whole_clade"
                item-value="name"
                label="Assigned taxon"
                density="compact"
                class="mt-1"
                clearable
                multiple
                chips
                closable-chips
              />
              <div class="d-flex align-center ga-3 flex-wrap">
                <v-select
                  v-model="filters.rank"
                  :items="$TaxonRank.enum"
                  label="Rank"
                  density="compact"
                  hide-details
                  clearable
                  multiple
                  chips
                  closable-chips
                  :min-width="200"
                />
                <v-select
                  v-model="filters.status"
                  :items="$TaxonStatus.enum"
                  label="Status"
                  density="compact"
                  hide-details
                  clearable
                  :min-width="200"
                />
              </div>
              <ClearableSwitch
                v-model="filters.confer"
                class="pl-2"
                label="Confer (cf.)"
                color-true="primary"
                color-false="red"
                :hint="
                  filters.confer
                    ? 'Show only bio material with a confer identification'
                    : filters.confer !== undefined
                      ? 'Show only bio material without a confer identification'
                      : undefined
                "
                density="compact"
              />
            </v-list-item>
          </v-list>
        </v-col>
      </v-row>
    </template>

    <template #item.code="{ value, item }: { value: string; item: OccurrenceListItem }">
      <span class="d-flex justify-space-between align-center">
        <RouterLink :text="value" :to="{ name: 'occurrence-item', params: { code: value } }" />
        <span class="d-flex align-center ga-2 justify-end">
          <v-icon
            v-if="item.type_status"
            icon="mdi-star-four-points"
            size="small"
            v-tooltip="item.type_status"
            density="compact"
          />
          <v-icon
            v-if="item.has_sequences"
            size="small"
            icon="mdi-dna"
            v-tooltip="`Sequence(s) available`"
          />
        </span>
      </span>
    </template>
    <template #item.sampling.site="{ value: { code, name } }: { value: SiteItem }">
      <RouterLink :to="{ name: 'site-item', params: { code } }" :text="name || code" />
    </template>
    <template #item.sampling.performed_on="{ value }: { value?: DateWithPrecision }">
      <span :class="['font-monospace text-caption', { 'text-muted': !value }]">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>

    <template #item.identification="{ value: identification }: { value: Identification }">
      <IdentificationChip :identification size="small" short />
    </template>
    <template #item.identification.identified_by="{ value: person }">
      <PersonChip v-if="person" :person size="small" short />
      <span v-else class="text-muted text-caption">Unknown</span>
    </template>
    <template #item.identification.identified_on="{ value }: { value?: DateWithPrecision }">
      <span :class="['font-monospace text-caption', { 'text-muted': !value }]">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>

    <!-- <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <BioMaterialFormDialog
        :dialog
        :model-value="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template> -->
  </CRUDTableServer>
</template>

<script setup lang="ts">
import { $TaxonRank, $TaxonStatus, TaxonRank, TypeStatus } from '@/api'

import {
  BioMatSortKey,
  DateWithPrecision,
  Identification,
  OccurrenceListItem,
  SiteItem,
  TaxonStatus
} from '@/api'
import {
  deleteOccurrenceMutation,
  listOccurrencesOptions,
  listOccurrencesQueryKey
} from '@/api/gen/@tanstack/vue-query.gen'
// import BioMaterialFormDialog from '@/features/occurrences/components/BioMaterialFormDialog.vue'
import CRUDTableServer from '@/components/toolkit/tables/CRUDTableServer.vue'
import ClearableSwitch from '@/components/toolkit/ui/ClearableSwitch.vue'
import DatasetPicker from '@/features/datasets/components/DatasetPicker.vue'
import TypeStatusPicker from '@/features/occurrences/components/TypeStatusPicker'
import PersonChip from '@/features/people/components/PersonChip'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import TaxonFilterPicker from '@/features/taxonomy/components/TaxonFilterPicker.vue'
import { useQueryClient } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { onMounted, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs } = useDisplay()

const [useDateRange, toggleUseDateRange] = useToggle(false)

type BiomatTableFilters = {
  year?: number | null
  year_end?: number | null
  datasets?: string[]
  type_status?: TypeStatus[]
  has_sequences?: boolean
  confer?: boolean
  whole_clade?: boolean
  rank?: TaxonRank[]
  status?: TaxonStatus
  taxa?: string[]
}

const filters = ref<BiomatTableFilters>({})

const headers: CRUDTableHeader<OccurrenceListItem>[] = [
  {
    title: 'Occurrence',
    children: [{ key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } }]
  },
  {
    title: 'Sampling',
    align: 'center',
    headerProps: { class: 'border-s' },
    children: [
      {
        key: 'sampling.site',
        title: 'Site'
      },
      { key: 'sampling.performed_on', title: 'Date', align: 'end' }
    ]
  },
  {
    key: 'identification',
    title: 'Identification',
    align: 'center',
    sortable: false,
    headerProps: { class: 'border-s' },
    children: [
      Identification.tableHeader({ key: 'identification', sort: undefined }),
      {
        key: 'identification.identified_by',
        title: 'Done by',
        align: 'center'
      },
      {
        key: 'identification.identified_on',
        title: 'Date',
        align: 'end'
      }
    ]
  }
] as const

type SortableColumn = Exclude<
  Extract<
    Exclude<(typeof headers)[number]['children'], undefined>[number]['key'] | 'meta.last_updated',
    string
  >,
  `data-table-${string}`
>

const sortKeyMap: Record<SortableColumn, BioMatSortKey> = {
  'sampling.site': 'site',
  'sampling.performed_on': 'sampling_date',
  'identification.taxon': 'taxon',
  'identification.identified_by': 'identified_by',
  'identification.identified_on': 'identified_on',
  'meta.last_updated': 'last_updated',
  code: 'code'
}

function sortKeyTransform(key: string | undefined): BioMatSortKey | undefined {
  return key ? sortKeyMap[key as SortableColumn] : undefined
}

const queryClient = useQueryClient()
function invalidateQuery() {
  queryClient.invalidateQueries({ queryKey: listOccurrencesQueryKey() })
}
onMounted(invalidateQuery)
</script>

<style scoped lang="scss"></style>
