<template>
  <div class="d-flex align-center ga-2">
    <v-chip
      v-if="target.kind === 'Unknown'"
      :text="target.kind"
      prepend-icon="mdi-help-circle"
      size="small"
    />
    <v-chip
      v-else-if="target.kind === 'Community'"
      :text="target.kind"
      prepend-icon="mdi-grid"
      size="small"
    />
    <TaxonChip
      v-else-if="target.kind === 'Taxa'"
      v-for="taxon in target.taxa"
      :taxon
      rounded
      size="small"
    />
  </div>
</template>

<script setup lang="ts">
import { SamplingTarget } from '@/api'
import TaxonChip from '../taxonomy/TaxonChip'

defineProps<{ target: SamplingTarget }>()
</script>

<style scoped lang="scss"></style>
