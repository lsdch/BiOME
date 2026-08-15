<template>
  <ActivableCardDialog
    ref="sunburst"
    :title
    class="w-100 d-flex flex-column"
    v-model="fullscreen"
    fullscreen
    :min-height="300"
    :height
  >
    <template #append v-if="compact">
      <v-menu :close-on-content-click="false" location="bottom" origin="top center">
        <template #activator="{ props }">
          <v-btn icon="mdi-cog" variant="text" color="" v-bind="props" />
        </template>
        <v-list :width="300" max-width="100vw">
          <v-list-item title="Taxonomic scope">
            <TaxonRankSlider
              v-if="scope"
              v-model="scope"
              +
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
          v-if="scope"
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
import { AppError, OccurrenceOverviewItem, Taxon, TaxonRank } from '@/api/adapters'
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
import { SunburstSeriesOption } from 'echarts'

use([SVGRenderer, TitleComponent, SunburstChart, VisualMapComponent, DataZoomComponent])

const [fullscreen, toggleFullscreen] = useToggle(false)

const { height = 600, ...props } = defineProps<{
  items?: OccurrenceOverviewItem[]
  title?: string
  loading?: boolean
  error?: AppError | string | null
  compact?: boolean
  height?: number
}>()

const scope = defineModel<[TaxonRank.NoSubgenus, TaxonRank.NoSubgenus]>('scope')

const rankOccurrences = ref<Record<TaxonRank, number>>({
  kingdom: 0,
  phylum: 0,
  class: 0,
  order: 0,
  family: 0,
  genus: 0,
  species: 0,
  subgenus: 0,
  subspecies: 0
})

function defaultScope(
  rankOccurrences: Record<TaxonRank, number>
): [TaxonRank.NoSubgenus, TaxonRank.NoSubgenus] {
  const ranks = Object.entries(rankOccurrences)
    .filter(([_, count]) => count > 0)
    .map(([rank]) => rank as TaxonRank)

  if (ranks.length === 0) {
    return ['order', 'species'] as const
  }

  const minRank = TaxonRank.noSubgenus(ranks[ranks.length - 1])
  const maxRank = TaxonRank.noSubgenus(TaxonRank.parentRank(ranks[0]) ?? 'kingdom')

  return [minRank, maxRank] as const
}

const cumulate = defineModel<boolean>('cumulate', { default: false })

type SunburstData = {
  name: string
  full_name: string
  children?: SunburstData[]
  value: [number, number] // [byClade, byTaxon]
  rank: TaxonRank
}

type SunburstIndex = Record<UUID, SunburstData>

const data = ref<SunburstData[]>([])

const maxOccurrences = ref({ byClade: 0, byTaxon: 0 })

watch(
  () => [props.items, scope.value, cumulate.value],
  () => (data.value = props.items ? buildPlotData(props.items) : []),
  { immediate: true }
)

function buildPlotData(items: OccurrenceOverviewItem[]) {
  if (!items.length) return []

  const kingdoms: UUID[] = []

  rankOccurrences.value = {
    kingdom: 0,
    phylum: 0,
    class: 0,
    order: 0,
    family: 0,
    genus: 0,
    species: 0,
    subgenus: 0,
    subspecies: 0
  }

  const itemsByID = items.reduce<SunburstIndex>((acc, item) => {
    acc[item.id] = {
      name: Taxon.shortName(item.name),
      full_name: item.name,
      children: [],
      rank: item.rank,
      value: [item.occurrences, item.occurrences]
    }

    if (item.rank === 'kingdom') {
      kingdoms.push(item.id)
    }
    rankOccurrences.value[item.rank] += item.occurrences

    return acc
  }, {})

  items.forEach(({ id, parent_id }) => {
    if (!parent_id) return

    const parent = itemsByID[parent_id]
    const child = itemsByID[id]

    if (parent && child) {
      parent.children!.push(child)
    }
  })

  if (!scope.value) scope.value = defaultScope(rankOccurrences.value)

  maxOccurrences.value = { byClade: 0, byTaxon: 0 }

  kingdoms.forEach((id) => {
    computeTotalOccurrences(itemsByID[id], scope.value!)
  })

  return kingdoms.flatMap((id) => trim([itemsByID[id]], scope.value!))
}

// Trims sunburst data to show only the selected ranks
function trim(data: SunburstData[], [r1, r2]: [TaxonRank, TaxonRank]): SunburstData[] {
  const children_updated = data.map((d) => {
    if (!d.children) return d
    d.children = trim(d.children, [r1, r2])
    return d
  })
  const trimmed = children_updated.filter(({ rank }) => {
    return !(TaxonRank.isAscendant(rank, r2) || TaxonRank.isDescendant(rank, r1))
  })
  if (trimmed.length === 0) {
    return children_updated.flatMap(({ children }) => children ?? [])
  }
  return trimmed
}

function computeTotalOccurrences(d: SunburstData, [r1, r2]: [TaxonRank, TaxonRank]) {
  if (!d.children) return
  d.children.forEach((v) => computeTotalOccurrences(v, [r1, r2]))
  d.value[0] += d.children.reduce<number>((a, b) => a + b.value[0], 0) ?? 0
  if (
    d.rank === r1 ||
    d.rank === r2 ||
    (TaxonRank.isDescendant(d.rank, r2) && TaxonRank.isAscendant(d.rank, r1))
  ) {
    maxOccurrences.value = {
      byClade: Math.max(maxOccurrences.value.byClade ?? 0, d.value[0]),
      byTaxon: Math.max(maxOccurrences.value.byTaxon ?? 0, d.value[1])
    }
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

const option = computed<ECBasicOption>((): ECBasicOption => {
  const series: SunburstSeriesOption = {
    type: 'sunburst',
    data: data.value,
    radius: [0, '90%'],
    sort: undefined,
    emphasis: {
      focus: 'ancestor',
      label: {
        show: true,
        fontWeight: 'bold',
        color: '#ef0d00',
        // color: 'black',
        // textBorderType: 'solid',
        // textBorderColor: 'black',
        // textBorderWidth: 1,
        // shadowColor: 'white',
        // shadowBlur: 5,
        textShadowColor: 'white',
        textShadowBlur: 1,
        width: 150,
        minAngle: 0,
        // backgroundColor: 'white',
        // borderRadius: 5,

        formatter: (params) => {
          const { data } = params
          if (!data) return ''
          const d = data as SunburstData
          return d.full_name || d.name
        }
      }
    },
    itemStyle: {
      borderWidth: 0.3
    },
    label: {
      rotate: 'tangential',
      minAngle: 25,
      // position: 'top',
      // distance: 100,
      padding: 5,
      width: 80,
      overflow: 'break'
      // align: 'right'
    }
  }

  return {
    visualMap: visualMap.value,
    series
  }
})
</script>

<style scoped lang="scss"></style>
