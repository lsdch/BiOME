<template>
  <div>
    <svg
      ref="colorLegend"
      :width
      :height
      xmlns="http://www.w3.org/2000/svg"
      @mousemove="
        (e) => {
          updateTooltip(e)
        }
      "
      @mouseleave="marker?.setAttribute('visibility', 'hidden')"
    >
      <defs>
        <linearGradient id="colorGradient" x1="0%" y1="100%" x2="0%" y2="0%">
          <stop offset="0%" stop-color="#440154" />
          <stop offset="25%" stop-color="#3b528b" />
          <stop offset="50%" stop-color="#21918c" />
          <stop offset="75%" stop-color="#5ec962" />
          <stop offset="100%" stop-color="#fde725" />
        </linearGradient>
      </defs>
      <rect
        :x="0"
        :y="0"
        :width
        :height
        :rx="rounded"
        :ry="rounded"
        stroke="black"
        :stroke-width="0.5"
        fill="url(#colorGradient)"
      />
      <rect
        ref="marker"
        :x="0"
        :y="(mouseY ?? 0) - 3"
        :width
        :height="6"
        :rx="2"
        :ry="2"
        fill="white"
        stroke="black"
        visibility="hidden"
      />
    </svg>
    <v-tooltip
      location="top left"
      :offset="[-15 - (mouseY ?? 0), -15]"
      :close-on-content-click="false"
      activator="parent"
    >
      {{ mouseY ?? 'undef' }} |
      {{ text }}
    </v-tooltip>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { computed, useTemplateRef } from 'vue'

const {
  width = 10,
  height = 200,
  rounded = 5,
  range
} = defineProps<{
  width?: number
  height?: number
  rounded?: number
  range: [number, number]
}>()

const colorLegend = useTemplateRef<SVGElement>('colorLegend')
const marker = useTemplateRef<SVGCircleElement>('marker')

const boundingRect = computed(() => colorLegend.value?.getBoundingClientRect())

const mouseY = ref<number>()
const text = computed(() => {
  if (!mouseY.value || !boundingRect.value) return ''
  const y = height - mouseY.value
  const percent = y / height
  return `${(percent * 100).toFixed(1)} %`
})
function updateTooltip(e: MouseEvent) {
  if (!boundingRect.value) return
  mouseY.value = e.clientY - boundingRect.value.top
  marker.value?.setAttribute('visibility', 'visible')
}
</script>

<style scoped lang="scss"></style>
