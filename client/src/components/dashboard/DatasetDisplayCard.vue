<template>
  <v-card
    :title="dataset.label"
    class="small-card-title"
    variant="tonal"
    density="compact"
    :to="{ name: targetRouteName, params: { slug: dataset.slug } }"
  >
    <template #subtitle>
      <v-card-subtitle class="text-caption">
        <v-icon :icon="Meta.icon(dataset.meta)" size="small"></v-icon>
        {{ Meta.toString(dataset.meta) }}
      </v-card-subtitle>
    </template>
    <template #prepend>
      <DatasetCategoryIcon :category="dataset.category" size="small" />
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { Meta, Dataset } from '@/api'
import { computed } from 'vue'
import DatasetCategoryIcon from '../datasets/DatasetCategoryIcon'

const { dataset } = defineProps<{
  dataset: Dataset
}>()

const targetRouteName = computed<string>(() => {
  switch (dataset.category) {
    case 'Occurrence':
      return 'occurrence-dataset-item'
    case 'Site':
      return 'site-dataset-item'
    case 'Seq':
      return 'seq-dataset-item'
  }
})
</script>

<style scoped lang="scss"></style>
