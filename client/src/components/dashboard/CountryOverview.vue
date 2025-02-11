<template>
  <ActivableCardDialog
    ref="country-overview"
    title="Sampled sites by country"
    class="w-100 d-flex flex-column"
    v-model="fullscreen"
    fullscreen
    :min-height="300"
    :height="500"
    :max-height="500"
  >
    <template #append>
      <v-btn
        color=""
        variant="text"
        :icon="fullscreen ? 'mdi-close' : 'mdi-fullscreen'"
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
import { getSitesCountByCountryOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { TreemapChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, VisualMapComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { SVGRenderer } from 'echarts/renderers'
import { ECBasicOption, TreemapSeriesOption } from 'echarts/types/dist/shared'
import { computed, ref, watch } from 'vue'
import VChart from 'vue-echarts'
import ActivableCardDialog from '../toolkit/ui/ActivableCardDialog.vue'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'

use([SVGRenderer, TitleComponent, TreemapChart, VisualMapComponent, TooltipComponent])

const [fullscreen, toggleFullscreen] = useToggle(false)

const { data: items, error, isPending } = useQuery(getSitesCountByCountryOptions())

type TreeMapData = {
  name: string
  code: string
  value: number
}

const data = ref<TreeMapData[]>([])
watch(items, (items) => {
  data.value =
    items?.map(({ code, name, sites_count }) => ({
      code,
      name,
      value: sites_count
    })) ?? []
})

const treemapSeries = computed<TreemapSeriesOption>(() => ({
  type: 'treemap',
  data: data.value,
  label: {
    show: true,
    formatter({ data }: { data: any }) {
      return `${data?.code}`
    }
  }
}))

const option = computed<ECBasicOption>(
  (): ECBasicOption => ({
    series: treemapSeries.value,
    tooltip: {
      formatter: function (info: any) {
        return `${info.name}: ${info.value} sites`
      }
    }
  })
)
</script>

<style scoped lang="scss"></style>
