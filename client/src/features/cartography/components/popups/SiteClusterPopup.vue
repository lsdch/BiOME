<template>
  <OccurrencesOverviewDialog ref="dialog" :max-width="1200" :data="data ?? []" />
  <v-card :min-width="250">
    <v-list density="compact">
      <v-list-item
        :title="pluralize(data?.length ?? 0, 'Site')"
        prepend-icon="mdi-map-marker-multiple"
        @click="() => dialog?.open?.('sites')"
      >
        <template #append>
          <v-badge inline :content="data?.length ?? 0" color="primary" />
        </template>
      </v-list-item>
      <v-list-item
        :title="pluralize(samplingEventsCount, 'Sampling event')"
        prepend-icon="mdi-package-down"
        @click="() => dialog?.open?.('samplings')"
      >
        <template #append>
          <v-badge inline :content="samplingEventsCount" color="warning" />
        </template>
      </v-list-item>
      <v-list-item
        :title="pluralize(occurrencesCount ?? 0, 'Occurrence')"
        prepend-icon="mdi-crosshairs-gps"
        @click="() => dialog?.open?.('occurrences')"
      >
        <template #append>
          <v-badge inline :content="occurrencesCount" color="success" />
        </template>
      </v-list-item>
      <v-list-item
        :title="pluralize(occurringTaxaCount, 'Sampled taxon', 'Sampled taxa')"
        prepend-icon="mdi-family-tree"
        @click="() => dialog?.open?.('sampled_taxa')"
      >
        <template #append>
          <v-badge inline :content="occurringTaxaCount" color="" />
        </template>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import OccurrencesOverviewDialog from '@/features/occurrences/components/tables/OccurrencesOverviewDialog.vue'
import { pluralize } from '@/lib/text'
import { computed, useTemplateRef } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'

const { data } = defineProps<{ data: SiteWithOccurrences[] | undefined }>()

const samplingEventsCount = computed(() => {
  return data?.reduce((acc, site) => acc + site.samplings.length, 0) ?? 0
})

const occurrencesCount = computed(() => {
  return (
    data?.reduce(
      (acc, site) => acc + site.samplings.reduce((a, s) => a + s.occurrences.length, 0),
      0
    ) ?? 0
  )
})

const occurringTaxaCount = computed(() => {
  return data
    ?.reduce(
      (acc, { samplings }) => {
        return acc.concat(samplings.flatMap((s) => s.occurrences))
      },
      [] as SiteWithOccurrences['samplings'][number]['occurrences']
    )
    ?.reduce((acc, { identification: { taxon } }) => {
      return acc.add(taxon.name)
    }, new Set()).size
})

type DialogType = ComponentExposed<typeof OccurrencesOverviewDialog>
const dialog = useTemplateRef<DialogType>('dialog')
</script>

<style scoped lang="scss"></style>
