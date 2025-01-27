<template>
  <v-select
    v-model="model"
    :items
    :loading
    v-bind="$attrs"
    item-title="code"
    :error-messages="error?.detail"
  >
    <template #item="{ item, props }">
      <v-list-item :subtitle="item.raw.label" v-bind="props" />
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { listGenesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const model = defineModel<any>()

const { data: items, isPending: loading, error } = useQuery(listGenesOptions())
</script>

<style scoped lang="scss"></style>
