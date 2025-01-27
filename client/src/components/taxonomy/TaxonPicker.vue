<template>
  <v-autocomplete
    v-model:search="search"
    no-filter
    :label
    :loading
    required
    :items="filteredItems"
    item-title="name"
    variant="outlined"
    clear-on-select
    auto-select-first
    placeholder="Enter search terms..."
    v-bind="$attrs"
    :error-messages="error?.detail"
  >
    <template #item="{ props, item }">
      <v-list-item v-bind="props" :title="item.raw.obj.name" class="fuzzy-search-item">
        <template #title>
          <v-list-item-title v-html="highlight(item.raw, 'name')" />
        </template>
        <template #subtitle>
          {{ item.raw.obj.authorship }}
        </template>
        <template #append>
          <v-chip>
            <span v-html="highlight(item.raw, 'rank')"></span>
          </v-chip>
          <FTaxonStatusIndicator :status="item.raw.obj.status" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { TaxonRank, TaxonWithParentRef } from '@/api'
import { listTaxaOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { KeysDeclaration, useFuzzyItemsFilter } from '@/composables/fuzzy_search'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import { FTaxonStatusIndicator } from './functionals'

const {
  label = 'Taxon',
  threshold = 0.7,
  limit = 10,
  ranks
} = defineProps<{
  label?: string
  ranks?: TaxonRank[]
  threshold?: number
  limit?: number
}>()

const search = ref('')
const { data: items, isPending: loading, error } = useQuery(listTaxaOptions({ query: { ranks } }))

const keys: KeysDeclaration<TaxonWithParentRef> = ['name', 'rank', 'status']

const { highlight, filteredItems } = useFuzzyItemsFilter(keys, search, items, {
  threshold,
  limit
})
</script>

<style scoped></style>
