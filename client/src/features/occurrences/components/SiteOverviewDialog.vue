<template>
  <CardDialog v-bind="props" v-model="dialog">
    <template #activator="props">
      <slot name="activator" v-bind="{ ...props, open }" />
    </template>
    <template #prepend>
      <v-card variant="tonal">
        <v-list-item :title="item.code" :subtitle="item.name"></v-list-item>
      </v-card>
    </template>
    <template #title>
      <v-tabs v-model="tab" class="text-body-large">
        <v-tab value="samplings">
          Sampling events
          <template #append>
            <v-badge
              :content="item.samplings.length"
              :color="item.samplings.length ? 'warning' : ''"
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
      </v-tabs>
    </template>

    <v-tabs-window class="overflow-y-auto" v-model="tab">
      <v-tabs-window-item value="samplings">
        <slot name="samplings-table" :samplings="item.samplings">
          <SamplingWithOccurrencesTable :samplings="item.samplings" />
        </slot>
      </v-tabs-window-item>
      <v-tabs-window-item value="occurrences">
        <slot name="occurrences-table" :occurrences="occurrences">
          <OccurrencesTable :occurrences />
        </slot>
      </v-tabs-window-item>
    </v-tabs-window>
  </CardDialog>
</template>

<script setup lang="ts" generic="WithSite extends boolean">
import { SiteWithOccurrences } from '@/api'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import OccurrencesTable from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import { computed, ref, useSlots } from 'vue'
import SamplingWithOccurrencesTable from './tables/SamplingWithOccurrencesTable.vue'

const dialog = defineModel<boolean>({ default: false })

const {
  title = 'Occurrences',
  item,
  ...props
} = defineProps<
  {
    item: SiteWithOccurrences
  } & CardDialogProps
>()

type Tab = 'samplings' | 'occurrences'
const tab = ref<Tab>('samplings')

const occurrences = computed(() => {
  return item.samplings.flatMap((s) => s.occurrences.map((o) => ({ ...o, sampling_date: s.date })))
})

useSlots()
function open(target: Tab = 'occurrences') {
  tab.value = target
  dialog.value = true
}

defineExpose({ open })
</script>

<style scoped lang="scss"></style>
