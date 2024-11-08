<template>
  <div id="abiotic-chart" class="abiotic-chart w-100"></div>
</template>

<script setup lang="ts">
import { AbioticParameter } from '@/api'
import Plotly, { Config, Data, Datum, Layout, PlotData } from 'plotly.js'
import { onMounted, ref, watch } from 'vue'

export type AbioticDataPoint = {
  date: Date
  y: number
}

export type AbioticData = {
  param: AbioticParameter
  points: AbioticDataPoint[]
}

const props = defineProps<{ data: AbioticData }>()

const config: Partial<Config> = {
  responsive: true,
  displayModeBar: false,
  displaylogo: false,

  modeBarButtonsToRemove: [
    'sendDataToCloud',
    'lasso2d',
    'select2d',
    'pan2d',
    'zoom2d',
    'zoomIn2d',
    'zoomOut2d',
    'autoScale2d',
    'resetScale2d'
  ]
}

function makePlotData(data: AbioticData): { data: Partial<Data>; layout: Partial<Layout> } {
  return {
    data: data.points.reduce<Partial<PlotData>>(
      (acc, { date, y }) => {
        ;(acc.x as Datum[]).push(date)
        ;(acc.y as Datum[]).push(y)

        return acc
      },
      {
        x: [],
        y: [],
        type: 'scatter'
      }
    ),
    layout: {
      title: data.param.label,
      xaxis: {
        showline: true,
        gridcolor: 'grey',
        tickformat: '%d %b %Y'
      },
      yaxis: {
        title: {
          text: data.param.unit,
          font: { weight: 800 }
        },
        gridcolor: 'grey',
        showline: true
        // ticksuffix: ` ${data.param.unit}`
      },
      paper_bgcolor: 'rgba(0,0,0,0)',
      plot_bgcolor: 'rgba(0,0,0,0)',
      margin: { b: 40, l: 40, r: 10, t: 30 }
    }
  }
}

onMounted(() => {
  const { data, layout } = makePlotData(props.data)
  Plotly.newPlot('abiotic-chart', [data], layout, config)
})

watch(
  () => props.data,
  (d) => {
    console.log('TOGGLE')
    const { data, layout } = makePlotData(d)
    Plotly.react('abiotic-chart', [data], layout)
  },
  { deep: true }
)

// const options: _DeepPartialObject<
//   CoreChartOptions<'line'> &
//     ElementChartOptions<'line'> &
//     PluginChartOptions<'line'> &
//     DatasetChartOptions<'line'> &
//     ScaleChartOptions<'line'> &
//     LineControllerChartOptions
// > = {
//   borderColor: '#08A1C3',
//   backgroundColor: '#08A1C3',
//   layout: {
//     padding: 20
//   },
//   scales: {
//     y: {
//       // beginAtZero: true,
//       ticks: {
//         callback(value) {
//           return `${value} °C`
//         }
//       }
//     }
//   }
// }
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
