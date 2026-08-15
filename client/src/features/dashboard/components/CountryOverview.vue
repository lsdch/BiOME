<template>
  <ActivableCardDialog
    ref="country-overview"
    :title="`Sampled ${facet} by country`"
    class="w-100 d-flex flex-column"
    v-model="fullscreen"
    fullscreen
    :min-height="300"
    :height
  >
    <template #append>
      <v-chip-group v-model="facet" mandatory color="primary">
        <v-chip value="occurrences">Occurrences</v-chip>
        <v-chip value="samplings">Samplings</v-chip>
      </v-chip-group>
      <v-btn
        v-if="!fullscreen"
        color=""
        variant="text"
        icon="mdi-fullscreen"
        @click="toggleFullscreen()"
      ></v-btn>
    </template>
    <CenteredSpinner v-if="isPending" :height="200" size="large" color="primary" />
    <v-card-text v-else-if="error">
      <v-alert color="error"> Failed to load {{ facet }} </v-alert>
    </v-card-text>
    <VChart v-else class="chart" :option autoresize />
  </ActivableCardDialog>
</template>

<script setup lang="ts">
import { listCountriesSummaryOptions } from '@/api/gen/@tanstack/vue-query.gen'
import ActivableCardDialog from '@/components/toolkit/ui/ActivableCardDialog.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { TreemapChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, VisualMapComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { SVGRenderer } from 'echarts/renderers'
import { ECBasicOption, TreemapSeriesOption } from 'echarts/types/dist/shared'
import { computed, ref } from 'vue'
import VChart from 'vue-echarts'

use([SVGRenderer, TitleComponent, TreemapChart, VisualMapComponent, TooltipComponent])

const { height = 600 } = defineProps<{
  height?: number
}>()

const [fullscreen, toggleFullscreen] = useToggle(false)

const { data: items, error, isPending } = useQuery(listCountriesSummaryOptions())

type Facet = 'samplings' | 'occurrences'
const facet = ref<Facet>('occurrences')

type TreeMapData = {
  name: string
  code?: string
  value: number
  children?: TreeMapData[]
}

const data = computed<TreeMapData[]>(() => {
  const bySubcontinent = items.value?.reduce(
    (acc, curr) => {
      const { subcontinent, ...rest } = curr
      acc[subcontinent] = acc[subcontinent] || {
        name: subcontinent,
        value: 0,
        children: []
      }
      const value = facet.value === 'samplings' ? curr.sampling_count : curr.occurrence_count
      acc[subcontinent].value += value
      acc[subcontinent].children?.push({
        ...rest,
        value
      })
      return acc
    },
    {} as Record<string, TreeMapData>
  )

  if (!bySubcontinent) return []
  return Object.values(bySubcontinent)
})
//   return (
//     items.value?.map(({ code, name, subcontinent, sampling_count, occurrence_count }) => ({
//       code,
//       name,
//       value: facet.value === 'samplings' ? sampling_count : occurrence_count
//     })) ?? []
//   )
// })

const treemapSeries = computed<TreemapSeriesOption>(() => {
  return {
    type: 'treemap',
    data: data.value,
    label: {
      show: true,
      formatter({ data }: { data: any }) {
        return `${data?.code}`
      },
      itemStyle: {
        borderColor: '#fff'
      }
    },
    levels: getLevelOption()
  }
})

function getLevelOption() {
  return [
    {
      itemStyle: {
        borderWidth: 0,
        gapWidth: 2
      }
    },
    {
      colorSaturation: [0.35, 0.5],
      itemStyle: {
        gapWidth: 1,
        borderColorSaturation: 0.6
      }
    }
  ]
}

const option = computed<ECBasicOption>((): ECBasicOption => {
  return {
    series: treemapSeries.value,
    tooltip: {
      formatter: function (info: any) {
        return `${info.name}: ${info.value} ${facet.value}`
      }
    }
  }
})
</script>

<style scoped lang="scss"></style>
