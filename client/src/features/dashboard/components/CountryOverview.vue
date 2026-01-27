<template>
  <ActivableCardDialog
    ref="country-overview"
    :title="`Sampled ${facet} by country`"
    class="w-100 d-flex flex-column"
    v-model="fullscreen"
    fullscreen
    :min-height="300"
    :height="500"
    :max-height="500"
  >
    <template #append>
      <v-chip-group v-model="facet" mandatory color="primary">
        <v-chip value="occurrences">Occurrences</v-chip>
        <v-chip value="sites">Sites</v-chip>
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
      <v-alert color="error"> Failed to load occurrences </v-alert>
    </v-card-text>
    <VChart v-else class="chart" :option autoresize />
  </ActivableCardDialog>
</template>

<script setup lang="ts">
import { getCountriesSummaryOptions } from '@/api/gen/@tanstack/vue-query.gen'
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

const [fullscreen, toggleFullscreen] = useToggle(false)

const { data: items, error, isPending } = useQuery(getCountriesSummaryOptions())

type Facet = 'sites' | 'occurrences'
const facet = ref<Facet>('occurrences')

type TreeMapData = {
  name: string
  code: string
  value: number
}

const data = computed<TreeMapData[]>(() => {
  return (
    items.value?.map(({ code, name, sites_count, occurrences_count }) => ({
      code,
      name,
      value: facet.value === 'sites' ? sites_count : occurrences_count
    })) ?? []
  )
})

const treemapSeries = computed<TreemapSeriesOption>(() => {
  return {
    type: 'treemap',
    data: data.value,
    label: {
      show: true,
      formatter({ data }: { data: any }) {
        return `${data?.code}`
      }
    }
  }
})

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
