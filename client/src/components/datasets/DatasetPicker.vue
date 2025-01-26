<template>
  <v-select
    v-model="model"
    :label
    :items
    :loading="isPending"
    item-title="label"
    return-item
    :error-messages="error?.detail"
  >
    <template #item="{ props, item }">
      <v-list-item v-bind="props">
        <template #append>
          <v-chip rounded> {{ item.raw.sites.length }}</v-chip>
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { SiteDataset } from '@/api'
import { listSiteDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const model = defineModel<SiteDataset>()
withDefaults(defineProps<{ label?: string }>(), { label: 'Dataset' })

const { data: items, error, isPending } = useQuery(listSiteDatasetsOptions())
</script>

<style scoped></style>
