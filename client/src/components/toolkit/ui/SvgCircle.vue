<template>
  <svg :width="size" :height="size" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
    <circle cx="50" cy="50" r="40" :fill="fillColor" :stroke="strokeColor" :stroke-width />
  </svg>
</template>

<script setup lang="ts">
const { size = 24, strokeWidth = 10 } = defineProps<{
  fillColor?: string
  strokeColor?: string
  size?: number | string
  strokeWidth?: number | string
}>()
</script>
<style scoped></style>
