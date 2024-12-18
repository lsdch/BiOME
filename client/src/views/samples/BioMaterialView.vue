<template>
  <CRUDTable
    class="fill-height"
    entity-name="Bio material"
    :headers
    :toolbar="{ title: 'Bio material', icon: 'mdi-package-variant' }"
    :fetch-items="SamplesService.listBioMaterial"
    :delete="({ code }: BioMaterial) => SamplesService.deleteBioMaterial({ path: { code } })"
    append-actions
    :search="search.term"
    :filter
    :mobile="xs"
    :owned-filter="search.owned"
  >
    <template #search="{ toggleMenu, menuOpen }">
      <v-inline-search-bar v-if="mdAndUp" v-model="search.term" label="Search term" class="mx-1" />
      <v-btn
        icon="mdi-dots-vertical"
        color="primary"
        @click="toggleMenu(true)"
        :active="menuOpen"
        size="small"
      ></v-btn>
    </template>

    <template #menu="{ toggleMenu }">
      <v-card rounded="t-0">
        <v-card-text>
          <v-inline-search-bar v-model="search.term" label="Search term" class="mx-1" />
        </v-card-text>
        <v-divider></v-divider>
        <v-row class="ma-0">
          <v-col cols="12" md="6">
            <v-list>
              <v-list-item prepend-icon="mdi-package-variant">
                <OccurrenceCategorySelect class="mt-1" v-model="search.category" label="Category" />
              </v-list-item>
              <v-list-item prepend-icon="mdi-family-tree">
                <TaxonPicker
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
                  hint="Show only nomenclatural type material"
                  persistent-hint
                  density="compact"
                />
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
        <v-divider> </v-divider>
        <v-list-item>
          <template #title>
            <v-switch
              v-model="search.owned"
              label="Owned items"
              color="primary"
              hint="Restrict the list to elements you contributed"
              persistent-hint
              class="ml-2"
              density="compact"
            />
          </template>
        </v-list-item>
        <v-divider></v-divider>
        <v-card-actions>
          <v-btn color="primary" text="OK" @click="toggleMenu(false)"></v-btn>
          <v-spacer></v-spacer>
          <v-btn color="" text="Clear" @click="search = {}"></v-btn>
        </v-card-actions>
      </v-card>
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
      #item.event.site="{ value: { code, name }, item }: { value: SiteInfo; item: BioMaterial }"
    >
      <RouterLink :to="{ name: 'site-item', params: { code } }" :text="name" />
    </template>
    <template #item.event.performed_on="{ value }: { value: DateWithPrecision }">
      <span>{{ DateWithPrecision.format(value) }}</span>
    </template>

    <template
      #item.identification.taxon="{ value: taxon, item }: { value: Taxon; item: BioMaterial }"
    >
      <TaxonChip :taxon size="small" short />
    </template>
    <template #item.identification.identified_by="{ value }: { value: PersonInner }">
      <PersonChip :person="value" size="small" short />
    </template>
    <template #item.identification.identified_on="{ value }">
      {{ DateWithPrecision.format(value) }}
    </template>
    <template #expanded-row-inject="{ item }">
      <v-row v-if="item.external" class="ma-0">
        <v-col v-if="item.external.original_taxon != undefined">
          <v-card class="fill-height">
            <v-list>
              <v-list-subheader>Original identification</v-list-subheader>
              <v-chip :text="item.external.original_taxon"></v-chip>
            </v-list>
          </v-card>
        </v-col>
        <v-col>
          <v-card
            title="References"
            class="fill-height small-card-title muted-title"
            density="compact"
          >
            <template #append>
              <v-btn
                icon="mdi-link-variant"
                :href="item.external.original_link"
                size="x-small"
                variant="tonal"
              />
            </template>
            <v-card-text>
              <ArticleChip
                v-for="article in item.published_in"
                :article
                class="ma-1"
                size="small"
              />
            </v-card-text>
            <v-divider></v-divider>
            <v-list density="compact">
              <v-list-item title="In collection">
                <div>
                  {{ item.external.archive.collection }}
                </div>
                Vouchers:
                <v-chip
                  v-for="voucher in item.external.archive.vouchers"
                  size="small"
                  :text="voucher"
                  class="mx-1"
                />
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
        <v-col>
          <v-card title="Content" class="fill-height small-card-title muted-title">
            <v-list density="compact">
              <v-list-item
                lines="one"
                :subtitle="item.external.content_description ?? 'No further description'"
              >
                <template #title>
                  Quantity: <v-chip :text="item.external.quantity" size="small" />
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import {
  $OccurrenceCategory,
  BioMaterial,
  PersonInner,
  SamplesService,
  SiteInfo,
  Taxon
} from '@/api'
import { DateWithPrecision, OccurrenceCategory } from '@/api/adapters'
import PersonChip from '@/components/people/PersonChip.vue'
import ArticleChip from '@/components/references/ArticleChip.vue'
import TaxonChip from '@/components/taxonomy/TaxonChip.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import OccurrenceCategorySelect from '@/components/toolkit/OccurrenceCategorySelect.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import ClearableSwitch from '@/components/toolkit/ui/ClearableSwitch.vue'
import { useToggle } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs, mdAndUp } = useDisplay()

type BiomatTableFilters = {
  term?: string
  category?: OccurrenceCategory
  nomenclaturalType?: boolean
  hasSequences?: boolean
  owned?: boolean
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

const headers: CRUDTableHeader[] = [
  {
    children: [{ key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } }]
  },
  {
    title: 'Sampling',
    align: 'center',
    headerProps: { class: 'border-s' },
    children: [
      { key: 'event.site', title: 'Site' },
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
      { key: 'identification.identified_by', title: 'Done by', align: 'center' },
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
