<template>
  <v-card v-bind="$attrs">
    <template #title>
      {{ dataset.label }}
    </template>
    <v-divider></v-divider>
    <template #text>
      {{ dataset.description }}
      <v-chip prepend-icon="mdi-map-marker">{{ dataset.sites?.length || 'No' }} sites</v-chip>
    </template>
    <template #actions>
      <ItemDateChip :date="dataset.meta.last_updated" icon="mdi-update"></ItemDateChip>
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { Dataset } from '@/api'
import ItemDateChip from '../toolkit/ItemDateChip.vue'

defineProps<{ dataset: Dataset }>()
</script>

<style scoped></style>
