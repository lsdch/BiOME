<template>
  <v-progress-linear v-if="loading" indeterminate />
  <TableToolbar title="Datasets"></TableToolbar>
  <v-container>
    <v-row>
      <v-col v-if="!datasets" class="text-center">There are currently no datasets</v-col>
      <v-col v-for="(dataset, key) in datasets" :key>
        <DatasetCard :dataset />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ErrorModel, LocationService } from '@/api'
import DatasetCard from '@/components/datasets/DatasetCard.vue'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { ref } from 'vue'

const loading = ref(true)
const datasets = ref(await fetch())
const error = ref<ErrorModel>()

async function fetch() {
  loading.value = true
  const { data, error: err } = await LocationService.listSiteDatasets()
  loading.value = false
  if (err !== undefined) {
    error.value = err
    return undefined
  }
  return data
}
</script>

<style scoped></style>
