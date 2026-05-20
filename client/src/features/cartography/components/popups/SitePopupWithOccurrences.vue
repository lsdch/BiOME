<template>
  <SingleSitePopup :item="item" :key="item.code">
    <template #append-items="{ item }">
      <v-divider />
      <SiteOverviewDialog :item :with-site="false" :max-width="1200">
        <template #activator="{ props, open }">
          <v-list-item @click="open('occurrences')" title="Occurrences">
            <template #append>
              <v-badge
                inline
                :content="item.samplings.reduce((sum, s) => sum + s.occurrences.length, 0)"
                color="success"
              />
            </template>
          </v-list-item>
          <v-list-item
            title="Sampling events"
            :subtitle="`Last visit: ${lastVisited}`"
            @click="open('samplings')"
          >
            <template #append>
              <v-badge inline :content="item.samplings.length" color="warning" />
            </template>
          </v-list-item>
        </template>
      </SiteOverviewDialog>
    </template>
  </SingleSitePopup>
</template>

<script setup lang="ts">
import { DateWithPrecision, SiteWithOccurrences } from '@/api'
import { lastSamplingDate } from '@/lib/dates'
import { computed } from 'vue'
import SingleSitePopup from './SingleSitePopup.vue'
import SiteOverviewDialog from '@/features/occurrences/components/SiteOverviewDialog.vue'

const props = defineProps<{
  item: SiteWithOccurrences
}>()

const lastVisited = computed(() => {
  const lastDate = lastSamplingDate(props.item)
  return lastDate
    ? DateWithPrecision.format(lastDate)
    : props.item.samplings.length
      ? 'Unknown'
      : 'Never'
})
</script>

<style scoped lang="scss"></style>
