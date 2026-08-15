<template>
  <OccurrencesOverviewDialog
    ref="dialog"
    :max-width="1200"
    :data="cellContent ?? []"
    :taxa="occurringTaxa"
    v-bind="$attrs"
  />
  <v-card :min-width="250">
    <v-list density="compact">
      <slot name="prepend-item" />
      <!-- <v-list-item
        :title="pluralize(data?.length ?? 0, 'Site')"
        prepend-icon="mdi-map-marker-multiple"
        @click="() => dialog?.open?.('sites')"
      >
        <template #append>
          <v-badge inline :content="data?.length ?? 0" color="primary" />
        </template>
      </v-list-item> -->
      <v-list-item
        v-if="data.samplings_count > 1"
        :title="pluralize(data?.samplings_count, 'Sampling event')"
        prepend-icon="mdi-package-down"
        @click="() => !loading && dialog?.open?.('samplings')"
      >
        <template #append>
          <v-badge inline :content="data?.samplings_count" color="info" />
        </template>
      </v-list-item>
      <template v-else-if="cellContent?.[0]">
        <v-list-item :title="cellContent[0].site.name" :subtitle="cellContent[0].site.locality">
          <template #append v-if="cellContent[0].site.country">
            <CountryChip :country="cellContent[0].site.country" size="small" class="ml-2" />
          </template>
        </v-list-item>
        <v-list-item prepend-icon="mdi-crosshairs-gps">
          <span class="text-overline">
            {{ cellContent[0].coordinates.latitude }}°N,
            {{ cellContent[0].coordinates.longitude }}°E
          </span>
        </v-list-item>
        <v-list-item
          v-if="cellContent[0].performed_on"
          prepend-icon="mdi-calendar"
          :title="DateWithPrecision.format(cellContent[0].performed_on, undefined, 'Unknown date')"
        >
        </v-list-item>
        <v-divider></v-divider>
      </template>

      <v-list-item
        v-if="
          data?.occurrences_count === 1 &&
          !!cellContent &&
          cellContent?.[0]?.occurrences?.length == 1
        "
        :title="cellContent?.[0]?.occurrences?.[0]?.code"
        :to="{
          name: 'occurrence-item',
          params: {
            id: cellContent[0].occurrences[0].id,
            code: cellContent[0].occurrences[0].code
          }
        }"
        target="_blank"
        append-icon="mdi-open-in-new"
        subtitle="Occurrence"
      >
        <template #title="{ title }">
          <span class="text-overline text-label-medium">{{ title }}</span>
        </template>
        <!-- <template #subtitle>
          <IdentificationChip
            size="small"
            :identification="cellContent?.[0]?.occurrences?.[0]?.identification"
          />
        </template> -->
      </v-list-item>
      <v-list-item
        v-else
        :title="pluralize(data?.occurrences_count ?? 0, 'Occurrence')"
        prepend-icon="mdi-crosshairs-gps"
        @click="() => !loading && dialog?.open?.('occurrences')"
      >
        <template #append>
          <v-badge inline :content="data?.occurrences_count" color="success" />
        </template>
      </v-list-item>

      <v-list-item
        v-if="occurringTaxaCount > 1"
        :title="pluralize(occurringTaxaCount, 'Sampled taxon', 'Sampled taxa')"
        prepend-icon="mdi-family-tree"
        @click="() => !loading && dialog?.open?.('sampled_taxa')"
      >
        <template #append>
          <v-badge inline :content="occurringTaxaCount" color="" />
        </template>
      </v-list-item>
      <v-list-item v-else-if="cellContent?.[0]?.occurrences?.[0]?.identification?.taxon">
        <IdentificationChip
          v-if="data.occurrences_count === 1"
          :identification="cellContent[0].occurrences[0].identification"
        />
        <TaxonChip v-else :taxon="cellContent[0].occurrences[0].identification.taxon" />
        <template #prepend>
          <span class="text-muted text-label"> Sampled taxon: </span>
        </template>
      </v-list-item>
      <slot name="append-item" />
    </v-list>
    <template #actions v-if="$slots['actions']">
      <slot name="actions"></slot>
    </template>
  </v-card>
</template>

<script setup lang="ts">
// import { SiteWithOccurrences } from '@/api'
import {
  DateWithPrecision,
  H3CellWithRichness,
  ListSamplingsWithOccurrencesAtCellData
} from '@/api'
import {
  listOccurringTaxaAtCellOptions,
  listSamplingsWithOccurrencesAtCellOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import OccurrencesOverviewDialog from '@/features/occurrences/components/tables/OccurrencesOverviewDialog.vue'
import CountryChip from '@/features/site/components/CountryChip'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import TaxonChip from '@/features/taxonomy/components/TaxonChip'
import { pluralize } from '@/lib/text'
import { useQuery } from '@tanstack/vue-query'
import { computed, useTemplateRef } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'

type ListOccurrencesParams = ListSamplingsWithOccurrencesAtCellData['query']

const { data, params, resolution } = defineProps<{
  data: H3CellWithRichness
  params: ListOccurrencesParams
  resolution: number
}>()

const occurringTaxaCount = computed(() => {
  return data ? data.species_richness + data.genus_richness + data.family_richness : 0
})

type DialogType = ComponentExposed<typeof OccurrencesOverviewDialog>
const dialog = useTemplateRef<DialogType>('dialog')

const { data: cellContent, isPending: loading } = useQuery(
  computed(() => ({
    initialData: [],
    ...listSamplingsWithOccurrencesAtCellOptions({
      query: params,
      path: { resolution, cell: data.h3_index }
    })
  }))
)

const { data: occurringTaxa } = useQuery(
  computed(() => ({
    initialData: [],
    ...listOccurringTaxaAtCellOptions({
      query: params,
      path: { resolution, cell: data.h3_index }
    })
  }))
)
</script>

<style scoped lang="scss"></style>
