<template>
  <v-select
    v-model="model"
    label="Include sites"
    :items
    item-props
    :hint="items.find(({ value }) => value == model)?.subtitle"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import { SiteSamplingStatus } from '@/api'

const model = defineModel<SiteSamplingStatus>({ default: 'Occurrences' })

const items = [
  { title: 'All', subtitle: 'All sites, including those never sampled', value: 'All' },
  { title: 'Sampled', subtitle: 'Only sites that have been sampled', value: 'Sampled' },
  {
    title: 'With occurrences',
    subtitle: 'Only sites that have been sampled and have occurrences',
    value: 'Occurrences'
  }
]
</script>

<style scoped lang="scss"></style>
