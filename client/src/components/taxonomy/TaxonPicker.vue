<template>
  <v-autocomplete
    v-model="model"
    no-filter
    :label
    :loading
    cache-items
    required
    :items
    return-object
    item-title="name"
    variant="outlined"
    @update:search="fetch"
    clear-on-select
    auto-select-first
    placeholder="Enter search terms..."
    v-bind="$attrs"
  >
    <template #item="{ props, item }">
      <v-list-item v-bind="props">
        <template #append>
          <v-chip>{{ item.raw.rank }}</v-chip>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { Taxon, TaxonomyService, TaxonRank } from '@/api'
import { useDebounceFn } from '@vueuse/core'
import { onMounted, ref } from 'vue'

const model = defineModel<Taxon>()

const props = withDefaults(
  defineProps<{
    label?: string
    ranks?: TaxonRank[]
  }>(),
  {
    label: 'Taxon'
  }
)

const loading = ref(false)

const items = ref<Taxon[]>()

const error = ref<string>()

async function _fetch(pattern?: string) {
  loading.value = true
  const { data, error: err } = await TaxonomyService.listTaxa({
    query: { limit: 10, pattern }
  }).finally(() => {
    loading.value = false
  })
  if (err !== undefined) {
    error.value = 'Failed to retrieve taxa'
    return
  }
  error.value = undefined
  items.value = props.ranks ? data.filter(({ rank }) => props.ranks?.includes(rank)) : data
}

const fetch = useDebounceFn(_fetch, 200)
onMounted(fetch)
</script>

<style scoped></style>
