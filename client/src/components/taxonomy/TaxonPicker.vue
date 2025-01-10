<template>
  <v-autocomplete
    v-model="model"
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
import { TaxonomyService, TaxonRank, TaxonWithParentRef } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'
import { KeysDeclaration, useFuzzyItemsFilter } from '@/composables/fuzzy_search'
import { ref } from 'vue'
import { FTaxonStatusIndicator } from './functionals'

const model = defineModel<any>()
const { threshold, limit, ranks } = withDefaults(
  defineProps<{
    label?: string
    ranks?: TaxonRank[]
    threshold?: number
    limit?: number
  }>(),
  {
    label: 'Taxon',
    limit: 10,
    threshold: 0.7
  }
)
const search = ref('')

const { items, loading, error } = useFetchItems(() =>
  // FIXME : not reactive wrt ranks prop
  TaxonomyService.listTaxa({ query: { ranks } })
)

const keys: KeysDeclaration<TaxonWithParentRef> = ['name', 'rank', 'status']

const { highlight, filteredItems } = useFuzzyItemsFilter(keys, search, items, {
  threshold,
  limit
})
</script>

<style scoped></style>
