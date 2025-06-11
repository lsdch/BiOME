<template>
  <SitePopup
    v-if="item"
    :item
    :options="{ keepInView: false, autoPan: false }"
    :showRadius="popupOpen"
    :zoom
  >
    <template #append-items="{ item }">
      <v-divider />
      <OccurrenceListDialog
        :occurrences="
          item.samplings.flatMap((s) => s.occurrences.map((o) => ({ ...o, sampling_date: s.date })))
        "
        :with-site="false"
        :max-width="1200"
      >
        <template #title>
          <DialogTitle :item />
        </template>
        <template #subtitle>
          <DialogSubtitle :item />
        </template>
        <template #activator="{ props }">
          <v-list-item v-bind="props" title="Occurrences">
            <template #append>
              <v-badge
                inline
                :content="item.samplings.reduce((sum, s) => sum + s.occurrences.length, 0)"
                color="success"
              />
            </template>
          </v-list-item>
        </template>
      </OccurrenceListDialog>
      <SamplingListDialog :samplings="item.samplings" :with-site="false" :max-width="1200">
        <template #title>
          <DialogTitle :item />
        </template>
        <template #subtitle>
          <DialogSubtitle :item />
        </template>
        <template #activator="{ props }">
          <v-list-item
            title="Sampling events"
            :subtitle="`Last visit: ${item.last_visited ? DateWithPrecision.format(item.last_visited) : 'Never'}`"
            v-bind="props"
          >
            <template #append>
              <v-badge inline :content="item.samplings.length" color="warning" />
            </template>
          </v-list-item>
        </template>
      </SamplingListDialog>
    </template>
  </SitePopup>
</template>

<script setup lang="tsx">
import { DateWithPrecision, SiteWithOccurrences } from '@/api'
import CountryChip from '@/components/sites/CountryChip'
import SitePopup from '@/components/sites/SitePopup.vue'
import { RouterLink } from 'vue-router'
import OccurrenceListDialog from '../../components/occurrence/OccurrenceListDialog.vue'
import SamplingListDialog from '../../components/occurrence/SamplingListDialog.vue'

const { item } = defineProps<{
  item: SiteWithOccurrences
  popupOpen: boolean
  zoom: number
}>()

const DialogTitle = ({ item }: { item: SiteWithOccurrences }) => (
  <div class="d-flex ga-2">
    <RouterLink
      class="font-monospace"
      to={{
        name: 'site-item',
        params: { code: item.code }
      }}
    >
      {item.code}
    </RouterLink>
    {item.country ? <CountryChip country={item.country} /> : undefined}
  </div>
)

const DialogSubtitle = ({ item }: { item: SiteWithOccurrences }) => (
  <span class="font-monospace">
    {item.coordinates.latitude}, {item.coordinates.longitude}
  </span>
)
</script>

<style scoped lang="scss"></style>
