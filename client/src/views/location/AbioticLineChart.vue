<template>
  <div class="w-100">
    <VChart class="chart" :option autoresize :init-options="{ height: 400 }" />
  </div>
</template>

<script setup lang="ts">
import { AbioticParameter } from '@/api/'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TitleComponent,
  TooltipComponent,
  VisualMapComponent
} from 'echarts/components'
import { use } from 'echarts/core'
import { SVGRenderer } from 'echarts/renderers'
import { ECBasicOption } from 'echarts/types/dist/shared'
import { DateTime } from 'luxon'
import { computed } from 'vue'
import VChart from 'vue-echarts'

use([SVGRenderer, TitleComponent, LineChart, VisualMapComponent, GridComponent, TooltipComponent])

export type AbioticDataPoint = {
  date: Date
  y: number
}

export type AbioticData = {
  param: AbioticParameter
  points: AbioticDataPoint[]
}

const props = defineProps<{ data: AbioticData }>()

const option = computed<ECBasicOption>(() => {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    xAxis: [
      {
        type: 'time',
        position: 'bottom',
        splitLine: {
          show: false
        },
        axisLabel: {
          formatter(dateMs: number) {
            return DateTime.fromMillis(dateMs).toFormat('dd\nLLL\nyyyy')
          }
        }
      }
    ],
    yAxis: {
      type: 'value',
      splitLine: {
        show: false
      }
    },
    series: [
      {
        type: 'line',
        showSymbol: true,
        data: props.data.points.map(({ date, y }) => [date, y])
      }
    ]
  }
})
</script>

<style lang="scss">
@use 'vuetify';

.abiotic-chart {
  height: 50vh;
  text {
    fill: rgb(var(--v-theme-on-surface)) !important;
  }
}
</style>
