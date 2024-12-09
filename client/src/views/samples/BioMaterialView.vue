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
        :items="$BioMaterialType.enum"
        v-model="search.type"
        label="Type"
        hide-details
        placeholder="Any"
        persistent-placeholder
        density="compact"
        class="mx-1"
        clearable
        persistent-clear
        :color="search.type ? 'primary' : undefined"
        :active="!!search.type"
      />
    </template>

    <template #item.code="{ value }">
      <RouterLink :text="value" :to="{ name: 'biomat-item', params: { code: value } }" />
    </template>
    <template #item.event.site="{ value: { code, name } }: { value: SiteInfo }">
      <RouterLink :to="{ name: 'site-item', params: { code } }" :text="name" />
    </template>
    <template #item.event.performed_on="{ value }: { value: DateWithPrecision }">
      <span>{{ DateWithPrecision.format(value) }}</span>
    </template>
    <template #item.type="{ value }: { value: BioMaterialType }">
      <v-icon
        v-bind="
          {
            Internal: { icon: 'mdi-cube-scan', color: 'primary' },
            External: { icon: 'mdi-open-in-new', color: 'warning' }
          }[value]
        "
        :title="value"
      />
    </template>
    <template #item.is_type="{ value }: { value: boolean }">
      <v-icon
        v-if="value"
        icon="mdi-star-four-points"
        size="small"
        title="This is a nomenclatural type"
        density="compact"
      />
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
  $BioMaterialType,
  BioMaterial,
  BioMaterialType,
  PersonInner,
  SamplesService,
  SamplingInner,
  SiteInfo,
  Taxon
} from '@/api'
import { DateWithPrecision } from '@/api/adapters'
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
  type?: BioMaterialType
}

const search = ref<BiomatTableFilters>({})
const filter = computed(() => {
  const { type } = search.value
  switch (type) {
    case undefined:
    case null:
      return () => true
    default:
      return (item: BioMaterial) => item.type === type
  }
})

const headers: CRUDTableHeader[] = [
  {
    children: [
      { key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } },
      { key: 'type', title: 'Category', width: 0, align: 'center' }
    ]
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
      { key: 'is_type', title: 'Nom. type', width: 0, align: 'center' },
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
