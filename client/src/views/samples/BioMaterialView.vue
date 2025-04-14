<template>
  <CRUDTableServer
    class="fill-height"
    entity-name="Bio-material"
    :headers
    :toolbar="{ title: 'Occurrences', icon: 'mdi-package-variant' }"
    :fetch-items="listBioMaterialOptions"
    :delete="{
      mutation: deleteBioMaterialMutation,
      params: ({ code }: BioMaterialWithDetails) => ({ path: { code } })
    }"
    :mobile="xs"
    show-expand
  >
    <!-- Search and filters panel -->
    <template #menu>
      <v-row class="ma-0">
        <v-col cols="12" md="6">
          <v-list>
            <v-list-item prepend-icon="mdi-package-variant">
              <OccurrenceCategorySelect class="mt-1" v-model="search.category" label="Category" />
            </v-list-item>
            <v-list-item prepend-icon="mdi-family-tree">
              <TaxonPicker
                v-model="search.taxon"
                item-value="name"
                label="Assigned taxon"
                density="compact"
                class="mt-1"
                hide-details
                clearable
              />
            </v-list-item>
          </v-list>
        </v-col>
        <v-col cols="12" md="6">
          <v-list density="compact">
            <v-list-item prepend-icon="mdi-star-four-points">
              <ClearableSwitch
                v-model="search.nomenclaturalType"
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
                v-model="search.hasSequences"
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
    <template
      #item.event.site="{ value: { code, name }, item }: { value: SiteItem; item: BioMaterial }"
    >
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
      <!-- <v-card flat class="fill-height small-card-title muted-title" density="compact"> -->
      <!-- <template #append>
              <v-btn
                icon="mdi-link-variant"
                :to="item.external.original_link"
                size="x-small"
                variant="tonal"
              />
            </template> -->
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
      <!-- </v-card> -->
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
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import CRUDTableServer from '@/components/toolkit/tables/CRUDTableServer.vue'
import ClearableSwitch from '@/components/toolkit/ui/ClearableSwitch.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs, mdAndUp } = useDisplay()

type BiomatTableFilters = {
  term?: string
  category?: OccurrenceCategory
  nomenclaturalType?: boolean
  hasSequences?: boolean
  taxon?: string
}

const search = ref<BiomatTableFilters>({})
function toggleNomenclaturalType(value?: boolean) {
  search.value.nomenclaturalType = value ?? !search.value.nomenclaturalType
}

function toggleSequenceFilter(value?: boolean) {
  search.value.hasSequences = value ?? !search.value.hasSequences
}

function nomenclaturalTypeFilter({ is_type }: BioMaterial) {
  return is_type
}
const filter = computed(() => {
  const { category, nomenclaturalType } = search.value
  switch (category) {
    case undefined:
    case null:
      return nomenclaturalType ? nomenclaturalTypeFilter : undefined
    default:
      return (item: BioMaterial) =>
        item.category === category && (nomenclaturalType ? nomenclaturalTypeFilter(item) : true)
  }
})

const headers: CRUDTableHeader<BioMaterialWithDetails>[] = [
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
        title: 'Site',
        sort(a, b) {
          return (a.name || a.code).localeCompare(b.name || b.code)
        }
      },
      { key: 'event.performed_on', title: 'Date', align: 'end', sort: DateWithPrecision.compare }
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
        align: 'center',
        sort(a: { name: string }, b: { name: string }) {
          return a.name.localeCompare(b.name)
        }
      },
      {
        key: 'identification.identified_by',
        title: 'Done by',
        align: 'center',
        sort(a, b) {
          return (a?.last_name || 'ZZZZZZZ').localeCompare(b?.last_name || 'ZZZZZZZ')
        }
      },
      {
        key: 'identification.identified_on',
        title: 'Date',
        align: 'end',
        sort: DateWithPrecision.compare
      }
    ]
  }
]
</script>

<style scoped lang="scss"></style>
