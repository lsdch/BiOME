<template>
  <v-select
    :items
    label="Sampling target"
    multiple
    chips
    closable-chips
    clearable
    placeholder="Any"
    persistent-placeholder
    item-props
    v-bind="$attrs"
  >
    <template #chip="{ item, props }" v-if="$slots.chip">
      <slot name="chip" :props :item> </slot>
    </template>
    <template #item="{ item, props }" v-if="$slots.item">
      <slot name="item" :props :item> </slot>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { SamplingTargetKind } from '@/api'

const items: { title: SamplingTargetKind; subtitle: string }[] = [
  { title: 'Community', subtitle: 'Sampled for community' },
  { title: 'Taxa', subtitle: 'Sampled for specific taxonomic group' },
  { title: 'Unknown', subtitle: 'Unknown sampling target' }
]
</script>

<style scoped lang="scss"></style>
