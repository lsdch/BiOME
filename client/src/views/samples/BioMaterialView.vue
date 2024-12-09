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
  >
    <template #search="">
      <v-inline-search-bar v-model="search.term" label="Search term" class="mx-1" />
      <v-select
        :items="$BioMaterialCategory.enum"
        v-model="search.category"
        label="Type"
        hide-details
        placeholder="Any"
        persistent-placeholder
        density="compact"
        class="mx-1"
        clearable
        persistent-clear
        :color="search.category ? 'primary' : undefined"
        :active="!!search.category"
        :prepend-inner-icon="
          search.category ? BioMaterialCategory.props[search.category].icon : undefined
        "
        :max-width="300"
      >
        <template #item="{ item, props }">
          <v-list-item
            v-bind="{
              ...props,
              ...BioMaterialCategory.props[item.raw]
            }"
            :class="`text-${BioMaterialCategory.props[item.raw].color}`"
          ></v-list-item>
        </template>
        <template #append>
          <v-btn
            :active="search.nomenclaturalType"
            active-color="primary"
            @click="toggleNomenclaturalType()"
            icon="mdi-star-four-points"
            size="small"
            rounded="sm"
            title="Show only nomenclatural type material"
          ></v-btn>
        </template>
      </v-select>
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
            v-bind="BioMaterialCategory.props[item.category]"
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
      <TaxonChip :taxon size="small" />
    </template>
    <template #item.identification.identified_by="{ value }: { value: PersonInner }">
      <PersonChip :person="value" size="small" />
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
          <v-card class="fill-height">
            <v-list density="compact">
              <v-list-subheader>Collection </v-list-subheader>
              <v-list-item :title="item.external.archive.collection">
                <template #append>
                  <v-btn
                    icon="mdi-link-variant"
                    :href="item.external.original_link"
                    size="small"
                    variant="tonal"
                  />
                </template>
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
          <v-card class="fill-height">
            <v-list density="compact">
              <v-list-subheader>Content</v-list-subheader>
              <v-list-item lines="one" :subtitle="item.external.content_description">
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
  $BioMaterial,
  $BioMaterialCategory,
  BioMaterial,
  PersonInner,
  SamplesService,
  SamplingInner,
  SiteInfo,
  Taxon
} from '@/api'
import { BioMaterialCategory, DateWithPrecision } from '@/api/adapters'
import SamplingCard from '@/components/events/SamplingCard.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import TaxonChip from '@/components/taxonomy/TaxonChip.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { computed, ref } from 'vue'

const focusSampling = {
  dialog: ref(false),
  sampling: ref<SamplingInner>()
}

type BiomatTableFilters = {
  term?: string
  category?: BioMaterialCategory
  nomenclaturalType?: boolean
}

const search = ref<BiomatTableFilters>({})
function toggleNomenclaturalType(value?: boolean) {
  search.value.nomenclaturalType = value ?? !search.value.nomenclaturalType
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
