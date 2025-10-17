<template>
  <v-combobox
    :items
    :loading="isPending"
    chips
    closable-chips
    multiple
    :error-message="error?.detail"
  >
    <template #chip="{ item, props }">
      <v-chip
        v-bind="props"
        :text="item.value"
        :color="items?.includes(item.value) ? 'primary' : 'success'"
      />
    </template>
  </v-combobox>
</template>

<script setup lang="ts">
import { getAccessPointsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const { data: items, isPending, error } = useQuery(getAccessPointsOptions())
</script>

<style scoped lang="scss"></style>
