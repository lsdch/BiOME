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
      <OccurrenceListDialog :occurrences="item.occurrences" :with-site="false" :max-width="1200">
        <template #title>
          <div class="d-flex ga-2">
            <RouterLink
              class="font-monospace"
              :to="{
                name: 'site-item',
                params: { code: item.code }
              }"
            >
              {{ item.code }}
            </RouterLink>
            <CountryChip v-if="item.country" :country="item.country" />
          </div>
        </template>
        <template #subtitle>
          <span class="font-monospace">
            {{ item.coordinates.latitude }}, {{ item.coordinates.longitude }}
          </span>
        </template>
        <template #activator="{ props }">
          <v-list-item v-bind="props">
            {{ pluralizeWithCount(item.occurrences.length, 'occurrence') }}
          </v-list-item>
        </template>
      </OccurrenceListDialog>
      <v-list-item class="font-monospace">
        {{ item.last_visited ? DateWithPrecision.format(item.last_visited) : 'Never' }}
        <template #append>
          <span class="text-caption text-muted"> Last visit</span>
        </template>
      </v-list-item>
    </template>
  </SitePopup>
</template>

<script setup lang="ts">
import { DateWithPrecision, SiteWithOccurrences } from '@/api'
import SitePopup from '@/components/sites/SitePopup.vue'
import { pluralizeWithCount } from '@/functions/text'
import OccurrenceListDialog from './OccurrenceListDialog.vue'
import CountryChip from '@/components/sites/CountryChip'

defineProps<{
  item: SiteWithOccurrences
  popupOpen: boolean
  zoom: number
}>()
</script>

<style scoped lang="scss"></style>
