<template>
  <CRUDTableServer
    class="fill-height"
    entity-name="Bio-material"
    :headers
    :filters
    :toolbar="{ title: 'Occurrences', icon: 'mdi-package-variant' }"
    :fetch-items="listBioMaterialOptions"
    :delete="{
      mutation: deleteBioMaterialMutation,
      params: ({ code }: BioMaterialWithDetails) => ({ path: { code } })
    }"
    :mobile="xs"
    show-expand
    :sort-key-transform
  >
    <!-- Search and filters panel -->
    <template #menu>
      <v-row class="ma-0">
        <v-col cols="12" md="6">
          <v-list>
            <v-list-item prepend-icon="mdi-package-variant">
              <OccurrenceCategorySelect class="mt-1" v-model="filters.category" label="Category" />
            </v-list-item>
            <v-divider />
            <v-list-item prepend-icon="mdi-family-tree">
              <TaxonPicker
                v-model="filters.taxon"
                item-value="name"
                label="Assigned taxon"
                density="compact"
                class="mt-1"
                hide-details
                clearable
              />
              <v-switch
                label="Include whole clade"
                color="primary"
                v-model="filters.whole_clade"
                hide-details
              />
            </v-list-item>
          </v-list>
        </v-col>
        <v-col cols="12" md="6">
          <v-list density="compact">
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
      </v-row>
    </template>

    <template #item.code="{ value, item }: { value: string; item: BioMaterial }">
      <span class="d-flex justify-space-between align-center">
        <RouterLink :text="value" :to="{ name: 'biomat-item', params: { code: value } }" />
        <span class="text-right">
          <v-icon
            v-if="item.is_type"
            icon="mdi-star-four-points"
            size="small"
            title="This is a nomenclatural type material"
            density="compact"
            class="mx-1"
          />
          <v-icon
            v-if="item.has_sequences"
            size="small"
            icon="mdi-dna"
            title="Sequence(s) available"
            class="mx-1"
          />
          <v-icon
            v-bind="OccurrenceCategory.props[item.category]"
            :title="item.category"
            class="mx-1"
          />
        </span>
      </span>
    </template>
    <template #item.event.site="{ value: { code, name } }: { value: SiteItem }">
      <RouterLink :to="{ name: 'site-item', params: { code } }" :text="name || code" />
    </template>
    <template #item.event.performed_on="{ value }: { value: DateWithPrecision }">
      <span
        :class="['font-monospace text-caption', { 'text-muted': value.precision == 'Unknown' }]"
      >
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>

    <template
      #item.identification.taxon="{ value: taxon, item }: { value: Taxon; item: BioMaterial }"
    >
      <TaxonChip :taxon size="small" short />
    </template>
    <template #item.identification.identified_by="{ value }: { value: PersonInner | undefined }">
      <PersonChip v-if="value" :person="value" size="small" short />
      <span v-else class="text-muted text-caption">Unknown</span>
    </template>
    <template #item.identification.identified_on="{ value }: { value: DateWithPrecision }">
      <span
        :class="['font-monospace text-caption', { 'text-muted': value.precision == 'Unknown' }]"
      >
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>
    <template #expanded-row-inject="{ item }">
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
import { BioMaterial, PersonInner, Taxon } from '@/api'
import {
  BioMaterialWithDetails,
  BioMatSortKey,
  DateWithPrecision,
  OccurrenceCategory,
  SiteItem
} from '@/api/adapters'
import {
  deleteBioMaterialMutation,
  listBioMaterialOptions
} from '@/api/gen/@tanstack/vue-query.gen'
// import BioMaterialFormDialog from '@/components/occurrence/BioMaterialFormDialog.vue'
import PersonChip from '@/components/people/PersonChip'
import ArticleChip from '@/components/references/ArticleChip'
import TaxonChip from '@/components/taxonomy/TaxonChip'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import OccurrenceCategorySelect from '@/components/toolkit/OccurrenceCategorySelect.vue'
import CRUDTableServer from '@/components/toolkit/tables/CRUDTableServer.vue'
import ClearableSwitch from '@/components/toolkit/ui/ClearableSwitch.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs } = useDisplay()

type BiomatTableFilters = {
  category?: OccurrenceCategory
  is_type?: boolean
  has_sequences?: boolean
  taxon?: string
  whole_clade?: boolean
}

const filters = ref<BiomatTableFilters>({})

const headers = [
  {
    title: 'BioMaterial',
    children: [{ key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } }]
  },
  {
    title: 'Sampling',
    align: 'center',
    headerProps: { class: 'border-s' },
    children: [
      {
        key: 'event.site',
        title: 'Site'
      },
      { key: 'event.performed_on', title: 'Date', align: 'end' }
    ]
  },
  {
    key: 'identification',
    title: 'Identification',
    align: 'center',
    headerProps: { class: 'border-s' },
    children: [
      {
        key: 'identification.taxon',
        title: 'Taxon',
        align: 'center'
      },
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
] as const satisfies CRUDTableHeader<BioMaterialWithDetails>[]

type SortableColumn = Extract<
  Exclude<(typeof headers)[number]['children'], undefined>[number]['key'] | 'meta.last_updated',
  string
>
const sortKeyMap: Record<SortableColumn, BioMatSortKey> = {
  'event.site': 'site',
  'event.performed_on': 'sampling_date',
  'identification.taxon': 'taxon',
  'identification.identified_by': 'identified_by',
  'identification.identified_on': 'identified_on',
  'meta.last_updated': 'last_updated',
  code: 'code'
}

function sortKeyTransform(key: string | undefined): BioMatSortKey | undefined {
  return key ? sortKeyMap[key as SortableColumn] : undefined
}
</script>

<style scoped lang="scss"></style>
