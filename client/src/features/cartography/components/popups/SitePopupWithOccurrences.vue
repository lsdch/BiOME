<template>
  <SingleSitePopup :item="item" :key="item.code">
    <template #append-items="{ item }">
      <v-divider />
      <OccurrenceListDialog
        :occurrences="
          item.samplings.flatMap((s) => s.occurrences.map((o) => ({ ...o, sampling_date: s.date })))
        "
        :with-site="false"
        :max-width="1200"
      >
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
      <SamplingTableDialog :samplings="item.samplings" :with-site="false" :max-width="1200">
        <template #activator="{ props }">
          <v-list-item
            title="Sampling events"
            :subtitle="`Last visit: ${item.last_visited ? DateWithPrecision.format(item.last_visited) : item.samplings.length ? 'Unknown' : 'Never'}`"
            v-bind="props"
          >
            <template #append>
              <v-badge inline :content="item.samplings.length" color="warning" />
            </template>
          </v-list-item>
        </template>
      </SamplingTableDialog>
    </template>
  </SingleSitePopup>
</template>

<script setup lang="ts">
import { DateWithPrecision, SiteWithOccurrences } from '@/api'
import SingleSitePopup from './SingleSitePopup.vue'
import OccurrenceListDialog from '@/features/occurrences/components/OccurrenceListDialog.vue'
import SamplingTableDialog from '@/features/occurrences/components/SamplingTableDialog.vue'

const props = defineProps<{
  item: SiteWithOccurrences
}>()
</script>

<style scoped lang="scss"></style>
