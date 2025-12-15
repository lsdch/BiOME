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
      </v-tabs>
    </template>
    <slot name="prepend-body" />
    <v-tabs-window v-model="tab" class="overflow-y-auto">
      <v-tabs-window-item value="sites">
        <SiteWithOccurrencesTable :sites="data" />
      </v-tabs-window-item>
      <v-tabs-window-item value="samplings">
        <SamplingWithOccurrencesTable with-site :samplings />
      </v-tabs-window-item>
      <v-tabs-window-item value="occurrences">
        <OccurrencesTable with-site :occurrences />
      </v-tabs-window-item>
      <v-tabs-window-item value="sampled_taxa">
        <SampledTaxaTable :occurrences />
      </v-tabs-window-item>
    </v-tabs-window>
  </CardDialog>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import OccurrencesTable from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SampledTaxaTable from '@/features/occurrences/components/tables/SampledTaxaTable.vue'
import SamplingWithOccurrencesTable from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import SiteWithOccurrencesTable from '@/features/occurrences/components/tables/SiteWithOccurrencesTable.vue'
import { computed } from 'vue'

const dialog = defineModel<boolean>({ default: false })

const { data, ...props } = defineProps<
  {
    data: SiteWithOccurrences[]
  } & CardDialogProps
>()

type Tab = 'sites' | 'samplings' | 'occurrences' | 'sampled_taxa'
const tab = defineModel<Tab>('tab', { default: 'sites' })

const samplings = computed(() =>
  data.flatMap(({ samplings, ...site }) => samplings.map((s) => ({ ...s, site })))
)

const occurrences = computed(() =>
  data.flatMap(({ samplings, ...site }) =>
    samplings.flatMap(({ occurrences, date }) =>
      occurrences.map((o) => ({ sampling_date: date, site, ...o }))
    )
  )
)

function open(target: Tab = 'sites') {
  tab.value = target
  dialog.value = true
}

defineExpose({ open })
</script>

<style scoped lang="scss"></style>
