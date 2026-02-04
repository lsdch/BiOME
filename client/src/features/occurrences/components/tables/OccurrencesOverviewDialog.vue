<template>
  <CardDialog v-bind="props" v-model="dialog">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <template #title>
      <v-tabs v-model="tab">
        <v-tab value="sites">
          Sites
          <template #append>
            <v-badge :content="data.length" :color="data.length ? 'primary' : ''" inline />
          </template>
        </v-tab>
        <v-tab value="samplings">
          Sampling Events
          <template #append>
            <v-badge
              :content="samplings.length"
              :color="samplings.length ? 'warning' : ''"
              inline
            />
          </template>
        </v-tab>
        <v-tab value="occurrences">
          Occurrences
          <template #append>
            <v-badge
              :content="occurrences.length"
              :color="occurrences.length ? 'success' : ''"
              inline
            />
          </template>
        </v-tab>
        <v-tab value="sampled_taxa"> Sampled taxa </v-tab>
        <v-tab value="sunburst">
          <v-icon size="x-large">mdi-chart-donut-variant</v-icon>
        </v-tab>
      </v-tabs>
    </template>
    <slot name="prepend-body" />
    <v-tabs-window v-model="tab" class="overflow-y-auto">
      <v-tabs-window-item value="sites">
        <slot name="sites-table" :sites="data">
          <SiteWithOccurrencesTable :sites="data" />
        </slot>
      </v-tabs-window-item>
      <v-tabs-window-item value="samplings">
        <slot name="samplings-table" :samplings>
          <SamplingWithOccurrencesTable with-site :samplings />
        </slot>
      </v-tabs-window-item>
      <v-tabs-window-item value="occurrences">
        <slot name="occurrences-table" :occurrences>
          <OccurrencesTable with-site :occurrences />
        </slot>
      </v-tabs-window-item>
      <v-tabs-window-item value="sampled_taxa">
        <slot name="sampled-taxa-table" :occurrences>
          <SampledTaxaTable :occurrences />
        </slot>
      </v-tabs-window-item>
      <v-tabs-window-item value="sunburst">
        <OccurrenceSunburst :items="sunburstData" />
      </v-tabs-window-item>
    </v-tabs-window>
  </CardDialog>
</template>

<script setup lang="ts" generic="Data extends SiteWithOccurrences">
import { OccurrenceAtSite, OccurrenceOverviewItem, SiteWithOccurrences, TaxonRank } from '@/api'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import OccurrenceSunburst from '@/features/dashboard/components/OccurrenceSunburst.vue'
import OccurrencesTable, {
  type OccurrenceTableItem
} from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SampledTaxaTable from '@/features/occurrences/components/tables/SampledTaxaTable.vue'
import SamplingWithOccurrencesTable, {
  SamplingTableItem
} from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import SiteWithOccurrencesTable from '@/features/occurrences/components/tables/SiteWithOccurrencesTable.vue'
import { computed, useSlots } from 'vue'

const dialog = defineModel<boolean>({ default: false })

const { data, ...props } = defineProps<
  {
    data: Data[]
  } & CardDialogProps
>()

type Tab = 'sites' | 'samplings' | 'occurrences' | 'sampled_taxa' | 'sunburst'
const tab = defineModel<Tab>('tab', { default: 'sites' })

type Sampling = SamplingTableItem<Data['samplings'][number], true, Omit<Data, 'samplings'>>
const samplings = computed<Array<Sampling>>(() =>
  data.flatMap(({ samplings, ...site }) => samplings.map((s) => ({ ...s, site })))
)
type Occurrence = OccurrenceTableItem<OccurrenceAtSite, true, Omit<Data, 'samplings'>>
const occurrences = computed<Array<Occurrence>>(() =>
  data.flatMap(({ samplings, ...site }) =>
    samplings.flatMap(({ occurrences, date }) =>
      occurrences.map((o) => ({ sampling_date: date, site, ...o }))
    )
  )
)

// const sunburstData = computed(() => {
//   const taxonMap: Record<
//     string,
//     OccurrenceOverviewItem
//   > = {}

//   occurrences.value.forEach((occ) => {
//     const taxon = occ.identification.taxon
//     if (!taxonMap[taxon.name]) {
//       taxonMap[taxon.name] = { name: taxon.name, rank: taxon.rank, occurrences: 0, parent_name: taxon. }
//     }
//     taxonMap[taxon.name]!.occurrences += 1
//   })

//   return Object.values(taxonMap)
// })

function open(target: Tab = 'sites') {
  tab.value = target
  dialog.value = true
}

defineExpose({ open })

useSlots()
</script>

<style scoped lang="scss"></style>
