<template>
  <v-card :max-width="300" v-if="sampledTaxa">
    <v-expansion-panels
      variant="accordion"
      multiple
      :model-value="
        Object.entries(sampledTaxa).reduce<TaxonRank[]>((acc, [rank, taxa]) => {
          if (Object.values(taxa).length <= 3) {
            acc.push(rank as TaxonRank)
          }
          return acc
        }, [])
      "
    >
      <template v-for="(taxa, rank) in sampledTaxa">
        <v-expansion-panel
          v-if="Object.values(taxa).length"
          :title="`${Object.values(taxa).length} ${rank}`"
          :value="rank"
        >
          <template #text>
            <TaxonChip v-for="taxon in taxa" :taxon size="small" class="ma-1" />
          </template>
        </v-expansion-panel>
      </template>
    </v-expansion-panels>
  </v-card>
</template>

<script setup lang="ts">
import { OccurrenceAtSite, TaxonRank } from '@/api'
import TaxonChip from '@/components/taxonomy/TaxonChip'
import { occurringTaxa, SampledTaxa } from '@/functions/occurrences'
import { computed } from 'vue'

const { occurrences } = defineProps<{
  occurrences: OccurrenceAtSite[]
}>()
const sampledTaxa = computed<SampledTaxa | undefined>(() => {
  return occurringTaxa(occurrences)
})
</script>

<style scoped lang="scss"></style>
