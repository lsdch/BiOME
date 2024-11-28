<template>
  <v-autocomplete
    label="Programs"
    :items
    item-title="label"
    item-subtitle="description"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import { EventsService } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'

const { items } = useFetchItems(EventsService.listPrograms)
</script>

<style scoped lang="scss"></style>
