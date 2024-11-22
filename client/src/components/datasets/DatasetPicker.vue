<template>
  <v-select v-model="model" :label :items :loading item-title="label" return-item>
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
import { Dataset, DatasetsService } from '@/api'
import { handleErrors } from '@/api/responses'
import { onMounted, ref } from 'vue'

const model = defineModel<Dataset>()
withDefaults(defineProps<{ label?: string }>(), { label: 'Dataset' })

const loading = ref<boolean>()

const items = ref<Dataset[]>()

onMounted(async () => (items.value = (await fetch()) ?? []))

async function fetch() {
  loading.value = true
  return DatasetsService.listDatasets()
    .then(handleErrors((err) => console.error(err)))
    .finally(() => (loading.value = false))
}
</script>

<style scoped></style>
