<template>
  <v-card-text class="d-flex flex-column">
    <v-slider
      v-model="model.radius"
      :max="20"
      :min="4"
      label="Hex cell radius"
      dense
      thumb-label
      :step="1"
      hide-details
    />
    <div class="d-flex w-100 justify-space-between">
      <v-range-slider
        :model-value="model.useRadiusRange ? model.radiusRange : [model.radius, model.radius]"
        @update:model-value="(v) => (model.radiusRange = v)"
        :min="2"
        :max="20"
        :step="0.5"
        label="Radius range"
        thumb-label
        :disabled="!model.useRadiusRange"
        hide-details
        color="primary"
      >
      </v-range-slider>
      <v-checkbox-btn v-model="model.useRadiusRange" class="flex-grow-0 flex-start" />
    </div>
  </v-card-text>
  <v-divider />
  <v-card-text class="d-flex flex-column">
    <v-range-slider
      v-if="model.asRange"
      v-model="model.opacity"
      :min="0"
      :max="1"
      :step="0.1"
      label="Opacity"
      hide-details
      thumb-label
    />
    <v-slider
      v-else
      v-model="model.opacity"
      :min="0"
      :max="1"
      :step="0.1"
      label="Opacity"
      hide-details
      thumb-label
    />
    <v-switch label="As range" v-model="model.asRange" hide-details />
  </v-card-text>
  <v-divider />
  <v-card class="small-card-title" title="Hover" flat density="compact">
    <template #append>
      <v-switch label="Grow on hover" v-model="model.hover.fill" hide-details color="primary" />
    </template>
    <v-card-text class="d-flex flex-column">
      <div class="d-flex justify-space-between w-100">
        <v-checkbox-btn v-model="model.hover.useScale" class="flex-grow-0" />
        <v-slider
          label="Scale on hover"
          v-model="model.hover.scale"
          hide-details
          dense
          clearable
          :min="0"
          :max="3"
          :step="0.1"
          thumb-label
          :disabled="!model.hover.useScale"
        />
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { watch } from 'vue'

export type HexgridConfig = {
  radius: number
  useRadiusRange: boolean
  radiusRange: [number, number]
  hover: {
    fill: boolean
    useScale: boolean
    scale: number
  }
} & (
  | {
      opacity: number
      asRange: false
    }
  | {
      opacity: [number, number]
      asRange: true
    }
)

const model = defineModel<HexgridConfig>({ required: true })

watch(
  () => model.value.asRange,
  (asRange) => {
    model.value.opacity = asRange
      ? [model.value.opacity as number, model.value.opacity as number]
      : 0.8
  }
)
</script>

<style lang="scss">
.v-input__control > .v-slider__container {
  width: 200px;
}
</style>
