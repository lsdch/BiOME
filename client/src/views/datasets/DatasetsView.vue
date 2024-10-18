<template>
  <v-progress-linear v-if="loading" indeterminate />
  <TableToolbar title="Datasets" icon="mdi-folder-table">
    <template #search>
      <v-inline-search-bar label="Search" v-model="search" />
    </template>
  </TableToolbar>
  <v-container fluid>
    <v-row>
      <v-col v-if="!filteredDatasets?.length" class="text-center">No datasets to display</v-col>
      <v-col v-for="(dataset, key) in filteredDatasets" :key sm="6">
        <DatasetCard :dataset :to="{ name: 'dataset-item', params: { slug: dataset.slug } }" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ErrorModel, LocationService } from '@/api'
import DatasetCard from '@/components/datasets/DatasetCard.vue'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { computed, ref } from 'vue'

const loading = ref(true)
const datasets = ref(await fetch())
const error = ref<ErrorModel>()

const search = ref<string>()

const filteredDatasets = computed(() => {
  return datasets.value?.filter((d) => d.label.toLowerCase().includes(search.value ?? ''))
})

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
