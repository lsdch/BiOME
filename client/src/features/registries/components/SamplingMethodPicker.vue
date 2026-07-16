<template>
  <v-autocomplete v-bind="$attrs" :items :loading="isPending" item-title="name" item-value="code"
    item-subtitle="description" prepend-inner-icon="mdi-hook" :error-messages="error?.detail">
  </v-autocomplete>
</template>

<script setup lang="ts">
import { listSamplingMethodsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const { data: items, isPending, error } = useQuery(listSamplingMethodsOptions())
</script>

<style scoped lang="scss"></style>
