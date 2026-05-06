<template>
  <div v-if="!hidden" class="color-scale-widget">
    <v-card density="compact" class="pa-2 opacity-70" theme="light" rounded="0">
      <div class="text-label-small mb-2">{{ bindingLabel }}</div>

      <div class="d-flex justify-space-between ga-2 text-label-small font-monospace">
        <span>{{ formatValue(min) }}</span>
        <div style="position: relative">
          <svg
            ref="colorLegend"
            :width="gradientWidth"
            :height="gradientHeight"
            xmlns="http://www.w3.org/2000/svg"
            @mousemove="handleMouseMove"
            @mouseleave="hoveredValue = null"
            style="cursor: crosshair"
          >
            <defs>
              <linearGradient id="colorScaleGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop
                  v-for="(color, index) in colorStops"
                  :key="index"
                  :offset="`${color.offset}%`"
                  :stop-color="`rgb(${color.rgb.join(',')})`"
                />
              </linearGradient>
            </defs>
            <rect
              :x="0"
              :y="0"
              :width="gradientWidth"
              :height="gradientHeight"
              stroke="#FFFFFF"
              stroke-width="1"
              :rx="10"
              :ry="10"
              fill="url(#colorScaleGradient)"
            />
          </svg>
          <div
            v-if="hoveredValue !== null"
            class="hover-tooltip"
            :style="{ left: tooltipX + 'px' }"
          >
            {{ formatValue(hoveredValue) }}
          </div>
        </div>
        <span>{{ formatValue(max) }}</span>
      </div>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, useTemplateRef } from 'vue'

interface ColorStop {
  offset: number
  rgb: [number, number, number]
}

const {
  colorRange = [],
  bindingLabel = 'Value',
  hidden = false,
  gradientWidth = 150,
  gradientHeight = 15,
  log = false,
  min,
  max
} = defineProps<{
  colorRange?: Array<[number, number, number] | { r: number; g: number; b: number }>
  min: number
  max: number
  bindingLabel?: string
  hidden?: boolean
  gradientWidth?: number
  gradientHeight?: number
  log?: boolean
}>()

const colorLegend = useTemplateRef<SVGElement>('colorLegend')

const hoveredValue = ref<number | null>(null)
const tooltipX = ref(0)

const normalizedColorRange = computed<[number, number, number][]>(() => {
  return colorRange.map((color) => {
    if (Array.isArray(color)) {
      return color as [number, number, number]
    }
    return [color.r, color.g, color.b]
  })
})

const colorStops = computed<ColorStop[]>(() => {
  const colors = normalizedColorRange.value
  if (colors.length === 0) return []

  return colors.map((rgb, index) => ({
    offset: (index / (colors.length - 1)) * 100,
    rgb
  }))
})

function formatValue(value: number): string {
  if (Number.isInteger(value)) {
    return value.toString()
  }
  return value.toFixed(2)
}

function handleMouseMove(event: MouseEvent) {
  if (!colorLegend.value) return

  const rect = colorLegend.value.getBoundingClientRect()
  const x = event.clientX - rect.left
  const percentage = Math.max(0, Math.min(1, x / gradientWidth))

  // Calculate the value based on position
  let value: number
  if (log) {
    // Logarithmic scale
    const logMin = Math.log(Math.max(min, 0.001))
    const logMax = Math.log(max)
    value = Math.exp(logMin + percentage * (logMax - logMin))
  } else {
    // Linear scale
    value = min + percentage * (max - min)
  }

  hoveredValue.value = Math.round(value)
  tooltipX.value = Math.round(Math.max(0, Math.min(x - 15, gradientWidth - 30)))
}
</script>

<style scoped lang="scss">
.color-scale-widget {
  svg {
    display: block;
  }
}

.hover-tooltip {
  position: absolute;
  top: -25px;
  background-color: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  white-space: nowrap;
  pointer-events: none;
  z-index: 10;
  transform: translateX(-50%);
  left: 0;
}
</style>
