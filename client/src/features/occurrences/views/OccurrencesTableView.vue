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
            <v-list-item prepend-icon="mdi-package-variant">
              <OccurrenceCategorySelect class="mt-1" v-model="filters.category" label="Category" />
            </v-list-item>
            <v-list-item prepend-icon="mdi-star-four-points">
              <ClearableSwitch
                v-model="filters.is_type"
                class="pl-2"
                label="Nomenclatural type"
                color-true="primary"
                color-false="red"
                hint="Show only <a href='https://en.wikipedia.org/wiki/Type_(biology)' target='_blank'>nomenclatural type</a> material"
                persistent-hint
                density="compact"
              >
                <template #message="{ message }">
                  <span v-html="message" />
                </template>
              </ClearableSwitch>
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
              <v-select
                v-model="filters.status"
                :items="$TaxonStatus.enum"
                label="Taxonomic status"
                density="compact"
                hide-details
              />
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
            v-if="item.is_type"
            icon="mdi-star-four-points"
            size="small"
            v-tooltip="`This is a nomenclatural type material`"
            density="compact"
          />
          <v-icon
            v-if="item.has_sequences"
            size="small"
            icon="mdi-dna"
            v-tooltip="`Sequence(s) available`"
          />
          <v-icon
            v-bind="OccurrenceCategory.props[item.category]"
            v-tooltip="item.category"
            size="small"
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
    <!-- <template #expanded-row-inject="{ item }">
      <v-list v-if="item.external">
        <v-list-item prepend-icon="mdi-newspaper-variant">
          <template #append>
            <span class="text-muted text-caption">Publications</span>
          </template>
          <ArticleChip v-for="article in item.published_in" :article class="ma-1" size="small" />
        </v-list-item>
        <v-list-item
          lines="one"
          :subtitle="item.external.content_description ?? 'No further description'"
          prepend-icon="mdi-hexagon-multiple"
        >
          <template #append>
            <span class="text-muted text-caption">Content</span>
          </template>
          <template #title>
            <v-chip :text="item.external.quantity" size="small" />
          </template>
        </v-list-item>
      </v-list>
      <v-divider v-if="item.external" />
    </template> -->
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
import { $TaxonStatus, OccurrencesService } from '@/api'

import {
  Identification,
  BioMatSortKey,
  DateWithPrecision,
  OccurrenceCategory,
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
import OccurrenceCategorySelect from '@/features/occurrences/components/OccurrenceCategorySelect.vue'
import PersonChip from '@/features/people/components/PersonChip'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import TaxonFilterPicker from '@/features/taxonomy/components/TaxonFilterPicker.vue'
import TaxonPicker from '@/features/taxonomy/components/TaxonPicker.vue'
import { useInfiniteQuery, useQueryClient } from '@tanstack/vue-query'
import { onMounted, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs } = useDisplay()

type BiomatTableFilters = {
  category?: OccurrenceCategory
  is_type?: boolean
  has_sequences?: boolean
  confer?: boolean
  whole_clade?: boolean
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
      Identification.tableHeader({ key: 'identification' }),
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
