<template>
  <div :style="{ width: `${width}px` }">
    <v-sparkline
      :gradient="palette(gradient)"
      gradient-direction="left"
      :model-value="Array.from({ length: 10 }, (_, i) => 0)"
      :line-width="thickness"
    />
  </div>
</template>

<script setup lang="ts">
import { ColorBrewerPaletteKey, palette } from '@/functions/color_brewer'

const { width = 100, thickness = 50 } = defineProps<{
  gradient: ColorBrewerPaletteKey
  width?: number
  thickness?: number
}>()
</script>

<style scoped lang="scss"></style>
