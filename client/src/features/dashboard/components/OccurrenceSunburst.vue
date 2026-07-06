<template>
  <ActivableCardDialog
    ref="sunburst"
    :title
    class="w-100 d-flex flex-column"
    v-model="fullscreen"
    fullscreen
    :min-height="300"
    :height="500"
  >
    <template #append v-if="compact">
      <v-menu :close-on-content-click="false" location="bottom" origin="top center">
        <template #activator="{ props }">
          <v-btn icon="mdi-cog" variant="text" color="" v-bind="props" />
        </template>
        <v-list :width="300" max-width="100vw">
          <v-list-item title="Taxonomic scope">
            <TaxonRankSlider
              v-model="scope"
              class="pt-8 px-7"
              density="compact"
              thumb-label="always"
            />
          </v-list-item>
          <v-list-item>
            <v-switch label="Use cumulative clade occurrences" v-model="cumulate" color="primary" />
          </v-list-item>
        </v-list>
      </v-menu>
      <v-btn
        v-if="!fullscreen"
        color=""
        variant="text"
        icon="mdi-fullscreen"
        @click="toggleFullscreen()"
      />
    </template>

    <slot v-if="loading" name="loading">
      <CenteredSpinner v-if="loading" :height="200" size="large" color="primary" />
    </slot>
    <slot v-else-if="error" name="error">
      <v-card-text>
        <v-alert color="error"> Failed to load occurrences </v-alert>
      </v-card-text>
    </slot>
    <div v-else-if="items?.length" class="d-flex align-center fill-height">
      <VChart class="chart" :option autoresize />
      <div v-if="!compact" class="d-flex flex-column align-start flex-shrink-0 pe-4">
        <v-btn
          v-if="!fullscreen"
          color=""
          variant="text"
          icon="mdi-fullscreen"
          @click="toggleFullscreen()"
        />
        <TaxonRankSlider
          label="Scope"
          v-model="scope"
          density="compact"
          thumb-label="always"
          direction="vertical"
          reverse
          show-ticks="always"
        />
        <v-switch
          label="Cumulative"
          class="text-no-wrap"
          v-model="cumulate"
          color="primary"
        ></v-switch>
      </div>
    </div>
    <v-card-text v-else>
      <v-alert>No occurrences to display</v-alert>
    </v-card-text>
  </ActivableCardDialog>
</template>

<script setup lang="ts">
import { ErrorModel, OccurrenceOverviewItem, Taxon, TaxonRank } from '@/api/adapters'
import ActivableCardDialog from '@/components/toolkit/ui/ActivableCardDialog.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { useToggle } from '@vueuse/core'
import { SunburstChart } from 'echarts/charts'
import { DataZoomComponent, TitleComponent, VisualMapComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { SVGRenderer } from 'echarts/renderers'
import { ECBasicOption, VisualMapComponentOption } from 'echarts/types/dist/shared'
import { computed, ref, watch } from 'vue'
import VChart from 'vue-echarts'
import TaxonRankSlider from './TaxonRankSlider.vue'

use([SVGRenderer, TitleComponent, SunburstChart, VisualMapComponent, DataZoomComponent])

const [fullscreen, toggleFullscreen] = useToggle(false)

const props = defineProps<{
  items?: OccurrenceOverviewItem[]
  title?: string
  loading?: boolean
  error?: ErrorModel | string | null
  compact?: boolean
}>()

const scope = defineModel<[TaxonRank, TaxonRank]>('scope', {
  default: () => ['Order', 'Species']
})
const cumulate = defineModel<boolean>('cumulate', { default: false })

type SunburstData = {
  name: string
  children?: SunburstData[]
  value: [number, number] // [byClade, byTaxon]
  rank: TaxonRank
}

type SunburstIndex = Record<string, SunburstData>

const data = ref<SunburstData[]>([])

const maxOccurrences = ref({ byClade: 0, byTaxon: 0 })

watch(
  () => [props.items, scope.value, cumulate.value],
  () => (data.value = props.items ? buildPlotData(props.items) : []),
  { immediate: true }
)

function buildPlotData(items: OccurrenceOverviewItem[]) {
  const itemsByName = items.reduce<SunburstIndex>(
    (acc, { name, occurrences, parent_name, rank }) => {
      acc[name] = acc[name] ?? {
        name: Taxon.shortName(name),
        children: [],
        rank,
        value: [occurrences, occurrences]
      }
      acc[name].rank = rank

      if (parent_name) {
        acc[parent_name] = acc[parent_name] ?? {
          name: Taxon.shortName(parent_name),
          children: [],
          value: [0, 0],
          rank: TaxonRank.parentRank(rank)!
        }
        acc[parent_name].children!.push(acc[name])
      }
      return acc
    },
    {}
  )
  maxOccurrences.value = { byClade: 0, byTaxon: 0 }
  computeTotalOccurrences(itemsByName['Animalia']!)
  return trim([itemsByName['Animalia']!], scope.value)
}

// Trims sunburst data to show only the selected ranks
function trim(data: SunburstData[], [r1, r2]: [TaxonRank, TaxonRank]): SunburstData[] {
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
  maxOccurrences.value = {
    byClade: Math.max(maxOccurrences.value.byClade ?? 0, d.value[0]),
    byTaxon: Math.max(maxOccurrences.value.byTaxon ?? 0, d.value[1])
  }
}

const visualMap = computed<VisualMapComponentOption>(() => ({
  min: 0,
  max: maxOccurrences.value[cumulate.value ? 'byClade' : 'byTaxon'],
  text: [maxOccurrences.value[cumulate.value ? 'byClade' : 'byTaxon']?.toString() ?? '0', '0'],
  dimension: cumulate.value ? 0 : 1,
  top: 'center',
  left: 0,
  textStyle: {
    color: 'rgb(var(--v-theme-on-surface))',
    fontWeight: 'bold'
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
