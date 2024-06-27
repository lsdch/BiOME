<template>
  <v-select
    v-bind="$attrs"
    v-model="model"
    :items="items"
    item-title="value"
    item-value="value"
    :label="label"
  >
    <template #item="{ item, props }">
      <v-list-item v-bind="props" :title="item.title" :subtitle="item.raw.description" />
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { CoordinatesPrecision } from '@/api'
import { Union } from 'ts-toolbelt'
import { ref } from 'vue'

const model = ref<CoordinatesPrecision>()

type PrecisionItem<T extends CoordinatesPrecision> = {
  value: T
  description: string
}

type PrecisionItems<P = Union.ListOf<CoordinatesPrecision>> = {
  [K in keyof P]: PrecisionItem<P[K] extends CoordinatesPrecision ? P[K] : never>
}

const items: PrecisionItems = [
  { value: '<100m', description: 'Coordinates of the site location' },
  { value: '<1km', description: 'Coordinates of the nearest landmark or populated place' },
  { value: '<10km', description: 'Coordinates of the nearest populated place' },
  { value: '10-100km', description: 'Coordinates of the region centroid' },
  { value: 'Unknown', description: 'Coordinates referential is unknown' }
]

withDefaults(defineProps<{ label?: string }>(), { label: 'Precision' })
</script>

<style scoped></style>
