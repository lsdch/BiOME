<template>
  <v-range-slider v-bind="props" :model-value @update:model-value="setValue" :max :step>
    <template #thumb-label="{ modelValue }">
      {{ $TaxonRank.enum[modelValue] }}
    </template>
  </v-range-slider>
</template>

<script setup lang="ts">
import { $Taxon, $TaxonRank, TaxonRank } from '@/api'
import { computed } from 'vue'
import { VRangeSlider } from 'vuetify/components'
const model = defineModel<[TaxonRank, TaxonRank]>({ required: true })

type RangeSliderProps = InstanceType<typeof VRangeSlider>['$props']
type TaxonRankSliderProps = /* @vue-ignore */ Omit<RangeSliderProps, 'modelValue'>

const {
  step = 1,
  max = $TaxonRank.enum.length - 1,
  ...props
} = defineProps<
  {
    step?: number
    max?: number
  } & TaxonRankSliderProps
>()

const inverseEnumIndex = computed(() => {
  const index: Record<string, number> = {}
  for (const [i, v] of Object.entries($TaxonRank.enum)) {
    index[v] = Number(i)
  }
  return index
})

const modelValue = computed(() => {
  const [start, end] = model.value
  return [inverseEnumIndex.value[start], inverseEnumIndex.value[end]]
})

function setValue([start, end]: [number, number]) {
  model.value = [$TaxonRank.enum[start], $TaxonRank.enum[end]]
}
</script>

<style scoped lang="scss"></style>
