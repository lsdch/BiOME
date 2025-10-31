<template>
  <v-menu location="right" origin="left" :close-on-content-click="false">
    <template #activator="{ props }">
      <v-list-item
        v-bind="props"
        append-icon="mdi-chevron-right"
        :title="
          pluralizeWithCount(occurringTaxaCount ?? 0, 'sampled taxon', { plural: 'sampled taxa' })
        "
      >
      </v-list-item>
    </template>
    <OccurringTaxaList v-if="occurrences" :occurrences="occurrences" />
  </v-menu>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import { HexPopupData } from '@/features/cartography/components/BaseMap.vue'
import { computed } from 'vue'
import OccurringTaxaList from '@/features/occurrences/components/OccurringTaxaList.vue'
import { pluralizeWithCount } from '@/lib/text'
const { data } = defineProps<{
  data?: HexPopupData<SiteWithOccurrences>[]
}>()

const occurrences = computed(() => {
  return data?.reduce(
    (acc, { data }) => {
      return acc.concat(data.samplings.flatMap((s) => s.occurrences))
    },
    [] as SiteWithOccurrences['samplings'][number]['occurrences']
  )
})

const occurringTaxaCount = computed(() => {
  return occurrences.value?.reduce((acc, { identification: { taxon } }) => {
    return acc.add(taxon.name)
  }, new Set()).size
})
</script>

<style scoped lang="scss"></style>
