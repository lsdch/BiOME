<template>
  <ActivableCardDialog
    ref="sunburst"
    title="Occurrences overview"
    class="w-100 d-flex flex-column"
    v-model="fullscreen"
    fullscreen
    :min-height="300"
    :height="500"
  >
    <template #append>
      <v-menu :close-on-content-click="false" location="bottom" origin="top center">
        <template #activator="{ props }">
          <v-btn icon="mdi-cog" variant="text" color="" v-bind="props" />
        </template>
        <v-list :width="300" max-width="100vw">
          <v-list-item title="Taxonomic scope">
            <TaxonRankSlider
              v-model="settings.scope"
              class="pt-8 px-7"
              density="compact"
              thumb-label="always"
            />
          </v-list-item>
          <v-list-item>
            <v-switch
              label="Use total clade occurrences"
              v-model="settings.totalByClade"
              color="primary"
            />
          </v-list-item>
        </v-list>
      </v-menu>
      <v-btn
        color=""
        variant="text"
        :icon="fullscreen ? 'mdi-close' : 'mdi-fullscreen'"
        @click="toggleFullscreen()"
      />
    </template>
    <CenteredSpinner v-if="isPending" :height="200" size="large" color="primary" />
    <v-card-text v-else-if="error">
      <v-alert color="error"> Failed to load occurrences </v-alert>
    </v-card-text>
    <VChart v-else class="chart" :option autoresize />
  </ActivableCardDialog>
</template>

<script setup lang="ts">
import { OccurrenceOverviewItem, Taxon, TaxonRank } from '@/api/adapters'
import { occurrenceOverviewOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { SunburstChart } from 'echarts/charts'
import { DataZoomComponent, TitleComponent, VisualMapComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { SVGRenderer } from 'echarts/renderers'
import { ECBasicOption, VisualMapComponentOption } from 'echarts/types/dist/shared'
import { computed, ref, watch } from 'vue'
import VChart from 'vue-echarts'
import ActivableCardDialog from '../toolkit/ui/ActivableCardDialog.vue'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'
import TaxonRankSlider from './TaxonRankSlider.vue'

use([SVGRenderer, TitleComponent, SunburstChart, VisualMapComponent, DataZoomComponent])

const [fullscreen, toggleFullscreen] = useToggle(false)

const { data: items, error, isPending } = useQuery(occurrenceOverviewOptions())

const settings = ref<{
  scope: TaxonRank[]
  totalByClade: boolean
}>({
  scope: ['Order', 'Species'],
  totalByClade: false
})

type SunburstData = {
  name: string
  children?: SunburstData[]
  value: [number, number]
  rank: TaxonRank
}

type SunburstIndex = Record<string, SunburstData>

const data = ref<SunburstData[]>([])
const maxOccurrences = ref([0, 0])

watch(items, (items) => (data.value = items ? buildPlotData(items) : []), { immediate: true })
watch(settings, () => (data.value = buildPlotData(items.value ?? [])), { deep: true })

function buildPlotData(items: OccurrenceOverviewItem[]) {
  const itemsByName = items.reduce<SunburstIndex>(
    (acc, { name, occurrences, parent_name, rank }) => {
      acc[name] = acc[name] ?? {
        name: Taxon.shortName(name),
        children: [],
        value: [occurrences, occurrences]
      }
      acc[name].rank = rank

      acc[parent_name] = acc[parent_name] ?? {
        name: Taxon.shortName(parent_name),
        children: [],
        value: [0, 0],
        rank: TaxonRank.parentRank(rank)
      }
      acc[parent_name].children!.push(acc[name])

      return acc
    },
    {}
  )
  maxOccurrences.value = [0, 0]
  computeTotalOccurrences(itemsByName.Animalia)
  return trim([itemsByName.Animalia], settings.value.scope)
}

// Trims sunburst data to show only the selected ranks
function trim(data: SunburstData[], [r1, r2]: TaxonRank[]) {
  const children_updated = data.map((d) => {
    if (!d.children) return d
    d.children = trim(d.children, [r1, r2])
    return d
  })
  const trimmed = children_updated.filter(({ rank }) => {
    return !(TaxonRank.isAscendant(rank, r1) || TaxonRank.isDescendant(rank, r2))
  })
  if (trimmed.length === 0) {
    return children_updated.flatMap(({ children }) => children ?? [])
  }
  return trimmed
}

function computeTotalOccurrences(d: SunburstData) {
  if (!d.children) return
  d.children.forEach((v) => computeTotalOccurrences(v))
  d.value[0] += d.children.reduce<number>((a, b) => a + b.value[0], 0) ?? 0
  maxOccurrences.value = [
    Math.max(maxOccurrences.value[0], d.value[0]),
    Math.max(maxOccurrences.value[1], d.value[1])
  ]
}

const visualMap = computed<VisualMapComponentOption>(() => ({
  min: 0,
  max: maxOccurrences.value[settings.value.totalByClade ? 0 : 1],
  text: [maxOccurrences.value[settings.value.totalByClade ? 0 : 1].toString(), '0'],
  dimension: settings.value.totalByClade ? 0 : 1,
  top: 'center',
  left: 0,
  textStyle: {
    color: 'rgb(var(--v-theme-on-surface))'
  },
  // Map the score column to color
  inRange: {
    color: ['#440154', '#3b528b', '#21918c', '#5ec962', '#fde725']
  }
}))

const option = computed<ECBasicOption>(
  (): ECBasicOption => ({
    // title: {
    //   text: 'Occurrences overview',
    //   // subtext: 'Source: https://worldcoffeeresearch.org/work/sensory-lexicon/',
    //   textStyle: {
    //     fontSize: 14,
    //     align: 'center'
    //   },
    //   subtextStyle: {
    //     align: 'center'
    //   }
    //   // sublink: 'https://worldcoffeeresearch.org/work/sensory-lexicon/'
    // },
    visualMap: visualMap.value,
    series: {
      type: 'sunburst',
      data: data.value,
      radius: [0, '85%'],
      sort: undefined,
      emphasis: {
        focus: 'ancestor'
      },
      itemStyle: {
        borderWidth: 0.3
      },
      label: {
        rotate: 'tangential',
        minAngle: 10
      }
    }
  })
)
</script>

<style scoped lang="scss"></style>
