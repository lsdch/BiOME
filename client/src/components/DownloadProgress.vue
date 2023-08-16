<template>
  <span v-if="finished"> Download complete : {{ printMB(completed) }} MB </span>
  <span v-else class="text-center w-100">
    Downloading {{ printMB(completed) }} / {{ printMB(total) }} MB
  </span>
  <span v-if="rate"> [ {{ printMB(rate) }} MB/s ]</span>
  <v-progress-linear
    v-if="progress"
    rounded
    :model-value="progress.progress"
    :color="finished ? 'green' : 'blue'"
  ></v-progress-linear>
</template>

<script lang="ts">
function MBytes(bytes: number) {
  return bytes / Math.pow(1024, 2)
}
export function printMB(bytes: number, digits = 2) {
  return MBytes(bytes).toFixed(digits)
}
</script>

<script setup lang="ts">
export type Progression = {
  completed: number
  total: number
  progress: number
  finished: boolean
  error?: string
  rate?: number
}

const progress = defineProps<Progression>()
</script>

<style scoped></style>
