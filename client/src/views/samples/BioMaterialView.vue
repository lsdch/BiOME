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
  >
    <template #search="">
      <v-inline-search-bar v-model="search.term" label="Search term" class="mx-1" />
      <OccurrenceCategorySelect v-model="search.category">
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
      </OccurrenceCategorySelect>
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
import OccurrenceCategorySelect from '@/components/toolkit/OccurrenceCategorySelect.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs } = useDisplay()

type BiomatTableFilters = {
  term?: string
  category?: OccurrenceCategory
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
