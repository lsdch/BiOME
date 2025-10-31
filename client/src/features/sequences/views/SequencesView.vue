<template>
  <CRUDTable
    class="fill-height"
    :headers
    :toolbar="{
      title: 'Sequences',
      icon: 'mdi-dna'
    }"
    entity-name="Sequence"
    :fetch-items="listSequencesOptions()"
    :delete="{
      mutation: deleteSequenceMutation,
      params: ({ code }: SequenceListItem) => ({ path: { code } })
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
      </v-row>
    </template>
    <template #item.code="{ value, item }: { value: string; item: Sequence }">
      <span class="d-flex justify-space-between align-center">
        <!-- Using zero-width spaces for better line breaks -->
        <RouterLink
          :text="CodeIdentifier.textWrap(value)"
          :to="{ name: 'sequence', params: { code: value } }"
        />
        <span class="d-flex align-center ga-2 justify-end">
          <!-- <v-icon
            v-if="item.external != undefined"
            :icon="ExtSeqOrigin.icon(item.external.origin)"
            v-tooltip="ExtSeqOrigin.description(item.external.origin)"
            size="small"
          /> -->
          <v-icon
            v-bind="OccurrenceCategory.props[item.category]"
            v-tooltip="item.category"
            size="small"
          />
        </span>
      </span>
    </template>
    <template #item.gene="{ value: gene }: { value: Gene }">
      <GeneChip :gene size="small" />
    </template>
    <template
      #item.occurrence.sampling.site="{
        value: { code, name },
        item
      }: {
        value: SiteItem
        item: Sequence
      }"
    >
      <RouterLink
        :class="{ 'font-monospace': !name }"
        :to="{ name: 'site-item', params: { code } }"
        :text="name || code"
      />
    </template>
    <template #item.occurrence.sampling.performed_on="{ value }: { value?: DateWithPrecision }">
      <span :class="['font-monospace', { 'text-muted': !value }]">{{
        DateWithPrecision.format(value)
      }}</span>
    </template>

    <template #item.identification="{ value: identification }">
      <IdentificationChip :identification size="small" short />
    </template>
    <template #item.identification.identified_by="{ value }: { value: Person }">
      <PersonChip :person="value" size="small" short />
    </template>
    <template #item.identification.identified_on="{ value }">
      <span :class="['font-monospace', { 'text-muted': !value }]">{{
        DateWithPrecision.format(value)
      }}</span>
    </template>

    <!-- ROW EXPANSION -->
    <template #expanded-row-footer-append="{ item }">
      <v-chip
        :text="item.occurrence.code"
        :to="{ name: 'occurrence-item', params: { code: item.occurrence.code } }"
        label
        :color="OccurrenceCategory.props[item.category].color"
        prepend-icon="mdi-package-variant"
        class="ma-2"
        v-tooltip="'Related occurrence'"
      >
      </v-chip>
    </template>
    <!-- <template #expanded-row-inject="{ item }">
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
                        name: occurrence-item,
                        params: { code: item.external?.source_sample?.code }
                      }
                    : undefined
                "
              ></v-list-item>
              <v-list-item prepend-icon="mdi-database" title="Database references">
                <SeqRefChip v-for="seqRef in item.external?.referenced_in" :seqRef size="small" />
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
    </template> -->
  </CRUDTable>
</template>

<script setup lang="ts">
import { Gene, Sequence, Taxon } from '@/api'
import {
  CodeIdentifier,
  DateWithPrecision,
  Identification,
  OccurrenceCategory,
  Person,
  SequenceListItem,
  SiteItem
} from '@/api/adapters'
import { deleteSequenceMutation, listSequencesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import OccurrenceCategorySelect from '@/features/occurrences/components/OccurrenceCategorySelect.vue'
import PersonChip from '@/features/people/components/PersonChip'
import GeneChip from '@/features/sequences/components/GeneChip'
import GenePicker from '@/features/sequences/components/GenePicker.vue'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import TaxonPicker from '@/features/taxonomy/components/TaxonPicker.vue'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'

const { xs, mdAndUp } = useDisplay()

type SeqTableFilters = {
  term?: string
  category?: OccurrenceCategory
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
      return (item: SequenceListItem) => item.category === category
  }
})

const headers: CRUDTableHeader<SequenceListItem>[] = [
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
      { key: 'occurrence.sampling.site', title: 'Site' },
      {
        key: 'occurrence.sampling.performed_on',
        title: 'Date',
        align: 'end',
        sort: DateWithPrecision.compare
      }
    ]
  },
  {
    key: 'identification',
    title: 'Identification',
    align: 'center',
    headerProps: { class: 'border-s' },
    children: [
      Identification.tableHeader({ key: 'identification' }),
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
