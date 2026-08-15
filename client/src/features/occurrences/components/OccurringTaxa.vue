<template>
  <!-- <v-tabs v-model="tab">
    <v-tab value="chart">Chart</v-tab>
    <v-tab value="table">Table</v-tab>
  </v-tabs> -->
  <!-- <v-tabs-window v-model="tab" class="overflow-y-auto">
    <v-tabs-window-item value="chart" class="pb-10"> -->
  <slot name="chart">
    <OccurrenceSunburst :items="taxa" flat />
  </slot>
  <!-- </v-tabs-window-item>
    <v-tabs-window-item value="table">
      <slot name="table">
        <SampledTaxaTable :occurrences />
      </slot>
    </v-tabs-window-item>
  </v-tabs-window> -->
</template>

<script setup lang="ts">
import { OccurrenceOverviewItem } from '@/api'
import OccurrenceSunburst from '@/features/dashboard/components/OccurrenceSunburst.vue'
import { ref } from 'vue'

const { taxa } = defineProps<{
  taxa: OccurrenceOverviewItem[]
}>()

const tab = ref<'chart' | 'table'>('chart')

// type RankKeys = Lowercase<TaxonRank> //& keyof (typeof $TaxonWithLineageNames)['properties']

// const rank_keys = $TaxonRank.enum
// .map((rank) => rank.toLocaleLowerCase() as TaxonRank.LowerCase)
// .filter((rank) => rank in $TaxonWithLineageNames['properties']) as RankKeys[]

// const sunburstItems = computed<OccurrenceOverviewItem[]>(() =>
//   taxa
//     .reduce((acc, { name, rank, occurrences, parent_name }) => {
//       const existing = acc.get(name)
//       rank_keys.forEach((rank) => {
//         const ancestor = taxon[rank]
//         if (!ancestor) return
//         if (acc.has(ancestor)) return
//         const parent_rank = TaxonRank.parentRank(TaxonRank.fromLowerCase(rank))?.toLowerCase() as
//           | RankKeys
//           | undefined
//         acc.set(ancestor, {
//           name: ancestor,
//           rank: TaxonRank.fromLowerCase(rank),
//           occurrences: 0,
//           parent_name: parent_rank ? taxon[parent_rank] : undefined
//         })
//       })
//       if (existing) {
//         existing.occurrences++
//       } else {
//         acc.set(taxon.name, {
//           name: taxon.name,
//           rank: taxon.rank,
//           occurrences: 1,
//           parent_name: taxon.parent
//         })
//       }
//       return acc
//     }, new Map<string, OccurrenceOverviewItem>())
//     .values()
//     .toArray()
// )
</script>

<style scoped lang="scss"></style>
