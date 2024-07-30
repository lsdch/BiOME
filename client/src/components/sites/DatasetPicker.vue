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
import { LocationService, SiteDataset } from '@/api'
import { handleErrors } from '@/api/responses'
import { onMounted, ref } from 'vue'

const model = defineModel<SiteDataset>()
withDefaults(defineProps<{ label?: string }>(), { label: 'Dataset' })

const loading = ref<boolean>()

const items = ref<SiteDataset[]>()

onMounted(async () => (items.value = (await fetch()) ?? []))

async function fetch() {
  loading.value = true
  return LocationService.listSiteDatasets()
    .then(handleErrors((err) => console.error(err)))
    .finally(() => (loading.value = false))
}
</script>

<style scoped></style>
