<template>
  <v-slider
    v-model="model"
    label="Show nearby sites"
    hide-details
    :min="0"
    :max="5"
    :step="1"
    glow
    :color="model > 0 ? 'primary' : ''"
    :thumb-size="15"
    @update:model-value="emit('update:radius', proximityRadius.value)"
  >
    <template #append>
      <span class="text-caption text-muted">
        {{ proximityRadius.label }}
      </span>
    </template>
  </v-slider>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const model = defineModel<number>({ default: 0 })

const radiusOptions = [
  { value: 0, label: 'Disabled' },
  { value: 100, label: '100m' },
  { value: 1_000, label: '1km' },
  { value: 10_000, label: '10km' },
  { value: 50_000, label: '50km' },
  { value: 100_000, label: '100km' }
]

const emit = defineEmits<{
  'update:radius': [radius: number]
}>()

const proximityRadius = computed(() => {
  return radiusOptions[model.value]
})
</script>

<style lang="scss">
.v-input__control > .v-slider__container {
  width: 100%;
  margin-right: 15px;
}
</style>
