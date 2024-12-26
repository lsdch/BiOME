<template>
  <v-range-slider
    v-model="modelValue"
    :ticks="{ ...$TaxonRank.enum }"
    :max="$TaxonRank.enum.length - 1"
    :step="1"
    v-bind="$attrs"
  >
    <template #thumb-label="{ modelValue }">
      {{ $TaxonRank.enum[modelValue] }}
    </template>
  </v-range-slider>
</template>

<script setup lang="ts">
import { $TaxonRank, TaxonRank } from '@/api'
import { computed } from 'vue'
const model = defineModel<TaxonRank[]>()

const inverseEnumIndex = computed(() => {
  const index: Record<string, number> = {}
  for (const [i, v] of Object.entries($TaxonRank.enum)) {
    index[v] = Number(i)
  }
  return index
})

const modelValue = computed({
  get: () => model.value?.map((v) => inverseEnumIndex.value[v]),
  set: (value: number[]) => {
    model.value = value.map((v) => $TaxonRank.enum[v])
  }
})
</script>

<style scoped lang="scss"></style>
