<template>
  <v-autocomplete
    label="Programs"
    :items
    item-title="label"
    item-subtitle="description"
    v-bind="$attrs"
    :loading="isPending"
    :error-messages="error?.detail"
  />
</template>

<script setup lang="ts">
import { listProgramsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const { data: items, error, isPending } = useQuery(listProgramsOptions())
</script>

<style scoped lang="scss"></style>
