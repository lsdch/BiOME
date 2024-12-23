<template>
  <v-card title="Occurrence records" prepend-icon="mdi-chart-arc" variant="outlined">
    <template #append>
      <v-menu :close-on-content-click="false" location="left">
        <template #activator="{ props }">
          <v-btn icon="mdi-cog" variant="plain" color="" v-bind="props"></v-btn>
        </template>
        <v-list :width="300">
          <v-list-item>
            <TaxonRankPicker
              v-model="settings.maxRank"
              class="ma-1"
              label="Max rank"
              hide-details
              density="compact"
            />
          </v-list-item>
          <v-list-item>
            <v-number-input
              density="compact"
              class="ma-1"
              hide-details
              label="Max depth"
              v-model="settings.maxDepth"
              :max="$TaxonRank.enum.length"
              :min="2"
              control-variant="default"
              clearable
            ></v-number-input>
          </v-list-item>
        </v-list>
      </v-menu>
    </template>
    <v-card-text class="d-flex align-center justify-center pa-0" style="min-height: 400px">
      <div id="occurrence-sunburst">
        <v-progress-circular v-if="loading" indeterminate size="x-large"></v-progress-circular>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { $TaxonRank, OccurrencesService, TaxonRank } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'
import { Data, PlotData, newPlot } from 'plotly.js'
import { computed, onMounted, ref, watch } from 'vue'
import TaxonRankPicker from '../taxonomy/TaxonRankPicker.vue'
import { ranksUpTo } from '../taxonomy/rank'

const { items, fetch, loading } = useFetchItems(OccurrencesService.occurrenceOverview)

const settings = ref<{
  maxRank: TaxonRank
  maxDepth?: number
}>({ maxRank: 'Family' })

const pltData = computed(() => {
  return items.value.reduce<Data>(
    (acc, { name, parent_name, occurrences, rank }) => {
      if (!settings.value.maxRank || ranksUpTo(settings.value.maxRank).includes(rank)) {
        acc.labels!.push(name)
        acc.parents!.push(parent_name)
        acc.values!.push(occurrences)
      }
      return acc
    },
    {
      type: 'sunburst',
      maxdepth: settings.value.maxDepth,
      labels: [],
      parents: [],
      values: [],
      branchvalues: 'remainder',
      marker: { line: { width: 2 }, colorscale: 'Viridis' }
    }
  )
})

watch(pltData, (data) => {
  newPlot(
    'occurrence-sunburst',
    [data],
    {
      margin: { l: 0, r: 0, b: 0, t: 0 },
      paper_bgcolor: 'transparent',
      plot_bgcolor: 'transparent',
      height: 400,
      transition: { duration: 100, easing: 'elastic' }
    },
    { displaylogo: false, displayModeBar: true }
  )
})

onMounted(async () => {
  items.value = await fetch()
})
</script>

<style scoped lang="scss"></style>
