<template>
  <CRUDTable
    class="fill-height"
    :headers
    :toolbar="{
      title: 'Sequences',
      icon: 'mdi-dna'
    }"
    entity-name="Sequence"
    :fetch-items="listSequencesOptions"
    :delete="{
      mutation: deleteSequenceMutation,
      params: ({ code }: Sequence) => ({ path: { code } })
    }"
    :mobile="xs"
    :filter
    :search="search"
    append-actions
  >
    <template #menu>
      <v-row class="ma-0">
        <v-col cols="12" md="6">
          <v-list>
            <v-list-item prepend-icon="mdi-circle-half-full">
              <OccurrenceCategorySelect v-model="search.category" label="Category" class="mt-1" />
            </v-list-item>
            <v-list-item prepend-icon="mdi-family-tree">
              <TaxonPicker
                v-model="search.taxon"
                item-value="name"
                label="Assigned taxon"
                class="mt-1"
                hide-details
                clearable
                density="compact"
              />
            </v-list-item>
            <v-list-item prepend-icon="mdi-tag">
              <GenePicker
                v-model="search.gene"
                label="Gene"
                class="mt-1"
                hide-details
                clearable
                density="compact"
                item-value="code"
              />
            </v-list-item>
          </v-list>
        </v-col>
        <v-col cols="12" md="6">
          <v-list>
            <v-list-item prepend-icon="mdi-package-variant">
              <ClearableSwitch
                v-model="search.hasBiomaterial"
                class="pl-2"
                label="Has bio-material parent sample"
                color-true="primary"
                color-false="red"
                hint="Internal sequences always have related bio-material"
                persistent-hint
                density="compact"
              />
            </v-list-item>
          </v-list>
        </v-col>
      </v-row>
    </template>
    <template #item.code="{ value, item }: { value: string; item: Sequence }">
      <span class="d-flex justify-space-between align-center">
        <!-- Using zero-width spaces for better line breaks -->
        <RouterLink
          :text="CodeIdentifier.textWrap(value)"
          :to="{ name: 'sequence', params: { code: value } }"
        />
        <span class="text-right">
          <v-icon
            v-if="item.external != undefined"
            :icon="ExtSeqOrigin.icon(item.external.origin)"
            :title="ExtSeqOrigin.description(item.external.origin)"
            size="small"
            class="mx-1"
          />
          <v-icon
            v-bind="OccurrenceCategory.props[item.category]"
            :title="item.category"
            class="mx-1"
            size="small"
          />
          <v-icon
            v-if="item.external?.source_sample"
            class="mx-1"
            icon="mdi-package-variant"
            title="Has related bio-material"
            size="small"
          ></v-icon>
        </span>
      </span>
    </template>
    <template #item.gene="{ value: gene }: { value: Gene }">
      <GeneChip :gene size="small" />
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

    <!-- ROW EXPANSION -->
    <template #expanded-row-inject="{ item }">
      <v-row class="ma-0">
        <v-col cols="12" md="6">
          <v-card flat class="flex-grow-1">
            <v-list>
              <v-list-item
                prepend-icon="mdi-package-variant"
                title="Related bio material"
                :subtitle="item.external?.source_sample?.code ?? 'None registered'"
                :disabled="!item.external?.source_sample"
                :to="
                  item.external?.source_sample?.code
                    ? {
                        name: 'biomat-item',
                        params: { code: item.external?.source_sample?.code }
                      }
                    : undefined
                "
              ></v-list-item>
              <v-list-item prepend-icon="mdi-database" title="Database references">
                <SeqRefChip v-for="seqRef in item.external?.referenced_in" :seq-ref size="small" />
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
        <v-divider :vertical="mdAndUp" class="flex-grow-0"></v-divider>
        <v-col cols="12" md="6">
          <v-card
            title="Comments"
            prepend-icon="mdi-comment"
            class="small-card-title flex-grow-1"
            flat
          >
            <v-card-text>
              {{ item.comments }}
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { BioMaterial, Gene, PersonInner, Sequence, SequencesService, SiteInfo, Taxon } from '@/api'
import { CodeIdentifier, DateWithPrecision, ExtSeqOrigin, OccurrenceCategory } from '@/api/adapters'
import { deleteSequenceMutation, listSequencesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import PersonChip from '@/components/people/PersonChip'
import GeneChip from '@/components/sequences/GeneChip'
import GenePicker from '@/components/sequences/GenePicker.vue'
import SeqRefChip from '@/components/sequences/SeqRefChip.vue'
import TaxonChip from '@/components/taxonomy/TaxonChip'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import OccurrenceCategorySelect from '@/components/toolkit/OccurrenceCategorySelect.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import ClearableSwitch from '@/components/toolkit/ui/ClearableSwitch.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs, mdAndUp } = useDisplay()

type SeqTableFilters = {
  term?: string
  category?: OccurrenceCategory
  hasBiomaterial?: boolean
  gene?: string
  taxon?: string
}

const search = ref<SeqTableFilters>({})

const filter = computed(() => {
  const { category } = search.value
  switch (category) {
    case undefined:
    case null:
      return undefined
    default:
      return (item: Sequence) => item.category === category
  }
})

const headers: CRUDTableHeader<Sequence>[] = [
  {
    children: [
      { key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } },
      {
        key: 'gene',
        title: 'Gene',
        width: 0,
        sort(a: Gene, b: Gene) {
          return a.code.localeCompare(b.code)
        }
      }
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
      {
        key: 'identification.taxon',
        title: 'Taxon',
        align: 'center',
        sort(a: { name: string }, b: { name: string }) {
          return a.name.localeCompare(b.name)
        }
      },
      { key: 'identification.identified_by', title: 'Done by', align: 'center', sortable: false },
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
