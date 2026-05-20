<template>
  <CardDialog v-bind="props" v-model="dialog">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <template #title>
      <v-tabs v-model="tab" class="text-body-large">
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
          <OccurringTaxa :occurrences />
          <!-- <SampledTaxaTable :occurrences /> -->
        </slot>
      </v-tabs-window-item>
    </v-tabs-window>
  </CardDialog>
</template>

<script setup lang="ts" generic="Data extends SiteWithOccurrences">
import { OccurrenceAtSite, SamplingDateWithOccurrences, SiteItem, SiteWithOccurrences } from '@/api'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import OccurrencesTable, {
  OccurrenceTableItem
} from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SamplingWithOccurrencesTable from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import SiteWithOccurrencesTable from '@/features/occurrences/components/tables/SiteWithOccurrencesTable.vue'
import { computed, useSlots } from 'vue'
import OccurringTaxa from '../OccurringTaxa.vue'

const dialog = defineModel<boolean>({ default: false })

const { data, ...props } = defineProps<
  {
    data: Data[]
  } & CardDialogProps
>()

type Tab = 'sites' | 'samplings' | 'occurrences' | 'sampled_taxa'
const tab = defineModel<Tab>('tab', { default: 'sites' })

type Sampling = SamplingDateWithOccurrences & { site: Omit<Data, 'samplings'> }
const samplings = computed<Array<Sampling>>(() =>
  data.flatMap(({ samplings, ...site }) => samplings.map((s) => ({ ...s, site })))
)

type Occurrence = OccurrenceTableItem<Omit<Data, 'samplings'>>
const occurrences = computed<Array<Occurrence>>(() =>
  data.flatMap<Occurrence>(({ samplings, ...site }) =>
    samplings.flatMap(({ occurrences, date }) =>
      occurrences.map<Occurrence>((o) => ({ sampling_date: date, site, ...o }))
    )
  )
)

function open(target: Tab = 'sites') {
  tab.value = target
  dialog.value = true
}

defineExpose({ open })

useSlots()
</script>

<style scoped lang="scss"></style>
