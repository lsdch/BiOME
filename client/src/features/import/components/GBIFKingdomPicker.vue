<template>
  <v-select
    :items="items"
    item-title="canonicalName"
    v-bind="$attrs"
    :loading="isPending"
    :error-messages="error?.title"
  >
    <template #item="{ item, props }">
      <v-list-item :title="item.canonicalName" :subtitle="item.authorship" v-bind="props">
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { getGbifKingdomsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const { data: items, error, isPending, refetch } = useQuery(getGbifKingdomsOptions())
</script>

<style scoped lang="scss"></style>
