<template>
  <v-card :min-width="250">
    <v-list density="compact">
      <SiteListDialog :sites="data?.flatMap(({ data }) => data)" :max-width="1200">
        <template #activator="{ props }">
          <v-list-item
            v-bind="props"
            :title="pluralize(data?.length ?? 0, 'Site')"
            prepend-icon="mdi-map-marker-multiple"
          >
            <template #append>
              <v-badge inline :content="data?.length ?? 0" color="primary" />
            </template>
          </v-list-item>
        </template>
      </SiteListDialog>
      <SamplingListDialog
        with-site
        :samplings="
          data?.flatMap(({ data: { samplings, ...site } }) => {
            return samplings.flatMap((s) => ({ ...s, site }))
          })
        "
        :max-width="1200"
      >
        <template #activator="{ props }">
          <v-list-item
            :title="pluralize(samplingEventsCount, 'Sampling event')"
            prepend-icon="mdi-package-down"
            v-bind="props"
          >
            <template #append>
              <v-badge inline :content="samplingEventsCount" color="warning" />
            </template>
          </v-list-item>
        </template>
      </SamplingListDialog>
      <OccurrenceListDialog
        with-site
        :occurrences="
          data?.flatMap(({ data: { samplings, ...site } }) => {
            return samplings.flatMap(({ occurrences, date }) =>
              occurrences.map((o) => ({
                ...o,
                sampling_date: date,
                site
              }))
            )
          })
        "
        :max-width="1200"
      >
        <template #activator="{ props }">
          <v-list-item
            v-bind="props"
            :title="pluralize(occurrencesCount ?? 0, 'Occurrence')"
            prepend-icon="mdi-crosshairs-gps"
          >
            <template #append>
              <v-badge inline :content="occurrencesCount" color="success" />
            </template>
          </v-list-item>
        </template>
      </OccurrenceListDialog>
      <v-divider />
      <AreaSampledTaxa :data />
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import { HexPopupData } from '@/components/maps/BaseMap.vue'
import { pluralize } from '@/functions/text'
import { computed } from 'vue'
import OccurrenceListDialog from '../../components/occurrence/OccurrenceListDialog.vue'
import SamplingListDialog from '../../components/occurrence/SamplingListDialog.vue'
import SiteListDialog from '../../components/occurrence/SiteListDialog.vue'
import AreaSampledTaxa from './AreaSampledTaxa.vue'

const { data } = defineProps<{ data: HexPopupData<SiteWithOccurrences>[] | undefined }>()

const samplingEventsCount = computed(() => {
  return data?.flatMap(({ data }) => data.samplings).length ?? 0
})

const occurrencesCount = computed(() => {
  return data
    ?.flatMap(({ data }) => data.samplings.map((s) => s.occurrences.length))
    .reduce((a, b) => a + b, 0)
})
</script>

<style scoped lang="scss"></style>
