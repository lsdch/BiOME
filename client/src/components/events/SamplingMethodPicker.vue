<template>
  <v-autocomplete
    :items
    :loading
    item-title="label"
    item-subtitle="description"
    v-bind="$attrs"
    prepend-inner-icon="mdi-hook"
  >
  </v-autocomplete>
</template>

<script setup lang="ts">
import { SamplingService } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'

const { items, loading } = useFetchItems(SamplingService.listSamplingMethods)
</script>

<style scoped lang="scss"></style>
