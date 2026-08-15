<template>
  <v-range-slider v-bind="props" :model-value @update:model-value="setValue" :max :step>
    <template #thumb-label="{ modelValue }">
      {{ TaxonRank.enumNoSubgenus[modelValue] }}
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
  max = TaxonRank.enumNoSubgenus.length - 1,
  ...props
} = defineProps<
  {
    step?: number
    max?: number
  } & TaxonRankSliderProps
>()

const inverseEnumIndex = computed(() => {
  const index: Record<string, number> = {}
  TaxonRank.enumNoSubgenus.forEach((rank, i) => {
    index[rank] = i
  })
  return index
})

const modelValue = computed(() => {
  const [start, end] = model.value
  return [inverseEnumIndex.value[start], inverseEnumIndex.value[end]]
})

function setValue([start, end]: [number, number]) {
  model.value = [TaxonRank.enumNoSubgenus[start], TaxonRank.enumNoSubgenus[end]]
}
</script>

<style scoped lang="scss"></style>
