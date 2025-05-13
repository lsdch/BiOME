<template>
  <v-menu location="right" origin="left" :close-on-content-click="false">
    <template #activator="{ props }">
      <v-list-item v-bind="props" append-icon="mdi-chevron-right">
        {{ occurringTaxaCount }}
        occurring {{ occurringTaxaCount === 1 ? 'taxon' : 'taxa' }}
      </v-list-item>
    </template>
    <OccurringTaxaList v-if="occurrences" :occurrences="occurrences" />
  </v-menu>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import { HexPopupData } from '@/components/maps/SitesMap.vue'
import { computed } from 'vue'
import OccurringTaxaList from './OccurringTaxaList.vue'
const { data } = defineProps<{
  data?: HexPopupData<SiteWithOccurrences>[]
}>()

const occurrences = computed(() => {
  return data?.reduce(
    (acc, { data }) => {
      return acc.concat(data.occurrences)
    },
    [] as SiteWithOccurrences['occurrences']
  )
})

const occurringTaxaCount = computed(() => {
  return occurrences.value?.reduce((acc, { taxon }) => {
    return acc.add(taxon.name)
  }, new Set()).size
})
</script>

<style scoped lang="scss"></style>
