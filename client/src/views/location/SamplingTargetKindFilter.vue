<template>
  <v-select
    :items
    label="Sampling target"
    multiple
    chips
    closable-chips
    clearable
    placeholder="Any"
    item-props
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import { SamplingTargetKind } from '@/api'

const items: { title: SamplingTargetKind | 'Never'; subtitle: string }[] = [
  { title: 'Never', subtitle: 'Never sampled' },
  { title: 'Community', subtitle: 'Sampled for community' },
  { title: 'Taxa', subtitle: 'Sampled for specific taxonomic group' },
  { title: 'Unknown', subtitle: 'Unknown sampling target' }
]
</script>

<style scoped lang="scss"></style>
