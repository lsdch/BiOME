<template>
  <v-autocomplete
    v-model="model"
    :items
    :loading
    item-title="label"
    item-value="code"
    item-subtitle="description"
    prepend-inner-icon="mdi-snowflake"
    :error-messages="error?.detail"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import { Fixative } from '@/api'
import { listFixativesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const model = defineModel<string | string[] | Fixative | Fixative[] | null>()
const { data: items, isPending: loading, error } = useQuery(listFixativesOptions())
</script>

<style scoped lang="scss"></style>
