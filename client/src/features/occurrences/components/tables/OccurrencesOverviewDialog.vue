<template>
  <CardDialog v-bind="props" v-model="dialog">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <template #title>
      <v-tabs v-model="tab" class="text-body-large">
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
        <v-tab value="sampled_taxa" :disabled="!taxa"> Sampled taxa </v-tab>
      </v-tabs>
    </template>
    <slot name="prepend-body" />
    <v-tabs-window v-model="tab" class="overflow-y-auto">
      <!-- <v-tabs-window-item value="sites">
        <slot name="sites-table" :sites="data">
          <SiteWithOccurrencesTable :sites="data" />
        </slot>
      </v-tabs-window-item> -->
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
      <v-tabs-window-item value="sampled_taxa" v-if="taxa">
        <slot name="sampled-taxa-table" :occurrences>
          <OccurringTaxa :taxa />
          <!-- <SampledTaxaTable :occurrences /> -->
        </slot>
      </v-tabs-window-item>
    </v-tabs-window>
  </CardDialog>
</template>

<script setup lang="ts" generic="Data extends Sampling & { occurrences?: BaseOccurrence[] }">
import { BaseOccurrence, Occurrence, OccurrenceOverviewItem, Sampling } from '@/api'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import OccurrencesTable from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SamplingWithOccurrencesTable from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import { computed, useSlots } from 'vue'
import OccurringTaxa from '../OccurringTaxa.vue'

const dialog = defineModel<boolean>({ default: false })

const { data, ...props } = defineProps<
  {
    data: Data[]
    taxa?: OccurrenceOverviewItem[]
  } & CardDialogProps
>()

type Tab = 'samplings' | 'occurrences' | 'sampled_taxa'
const tab = defineModel<Tab>('tab', { default: 'samplings' })

const samplings = computed<Array<Data>>(() => data)

// type Occurrence = OccurrenceTableItem<Omit<Data, 'samplings'>>
const occurrences = computed<Array<Occurrence>>(() =>
  data.flatMap<Occurrence>(
    ({ occurrences, ...sampling }) => occurrences?.flatMap((o) => ({ ...o, sampling })) ?? []
  )
)

function open(target: Tab = 'samplings') {
  tab.value = target
  dialog.value = true
}

defineExpose({ open })

useSlots()
</script>

<style scoped lang="scss"></style>
