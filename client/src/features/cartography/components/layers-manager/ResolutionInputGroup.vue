<template>
  <v-list-item
    :title="`Resolution ${model.resolution}`"
    lines="two"
    :subtitle="`≈ ${displayHexagonArea(model.resolution)} | ≈ ${displayHexagonRadius(model.resolution)} radius`"
  >
    <template #append>
      <v-chip-group v-model="model.resolutionMode" mandatory color="primary">
        <v-chip text="Auto" value="auto"></v-chip>
        <v-chip text="Manual" value="manual"></v-chip>
      </v-chip-group>
      <InlineHelp>
        Occurrences are clustered in high resolution hexagon cells. <br />
        The resolution can be automatically adjusted based on the zoom level of the map, or manually
        set to a specific resolution.
      </InlineHelp>
    </template>
  </v-list-item>
  <template v-if="model.resolutionMode === 'manual'">
    <v-list-item>
      <!-- <span class="text-muted">
            {{ `~ ${displayHexagonArea(model.resolution)}` }}
          </span> -->
      <v-slider
        class="pt-1 pb-3 px-3"
        v-model="model.resolution"
        :track-size="2"
        density="compact"
        :min="2"
        :max="12"
        :step="1"
        :ticks="[12, 8, 5, 2]"
        show-ticks
        hide-details
      />
    </v-list-item>
    <v-divider></v-divider>
  </template>
</template>

<script setup lang="ts">
import InlineHelp from '@/components/toolkit/ui/InlineHelp.vue'
import { getHexagonAreaAvg, getHexagonEdgeLengthAvg, UNITS } from 'h3-js'

interface ResolutionModel {
  resolution: number
  resolutionMode: 'auto' | 'manual'
}

const model = defineModel<ResolutionModel>({ required: true })

function hexagonArea(resolution: number) {
  // return (3 * Math.sqrt(3) * radius ** 2) / 2
  return getHexagonAreaAvg(resolution, resolution > 8 ? UNITS.m2 : UNITS.km2)
}

function displayHexagonArea(resolution: number) {
  const area = hexagonArea(resolution)
  return resolution > 8 ? `${Math.round(area)} m²` : `${area.toFixed(1)} km²`
}

function displayHexagonRadius(resolution: number) {
  const edgeLength = getHexagonEdgeLengthAvg(resolution, resolution > 7 ? UNITS.m : UNITS.km)
  return resolution > 7 ? `${Math.round(edgeLength)} m` : `${edgeLength.toFixed(1)} km`
}
</script>

<style scoped lang="scss"></style>
