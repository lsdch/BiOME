<template>
  <v-autocomplete
    v-model="model"
    v-model:search="search"
    :label
    :loading
    :chips
    required
    :items="items"
    item-title="name"
    variant="outlined"
    clear-on-select
    auto-select-first
    placeholder="Enter search terms..."
    :multiple
    v-bind="$attrs"
    :error-messages="error?.detail"
  >
    <template #chip="{ item, props }" v-if="chips">
      <TaxonChip :taxon="item.raw" />
    </template>
    <template #item="{ props, item }">
      <v-list-item v-bind="props" :title="item.raw.name" class="fuzzy-search-item">
        <!-- <template #title>
          <v-list-item-title v-html="highlight(item.raw, 'name')" />
        </template> -->
        <template #subtitle>
          {{ item.raw.authorship }}
        </template>
        <template #append>
          <v-chip>
            {{ item.raw.rank }}
            <!-- <span v-html="highlight(item.raw, 'rank')"></span> -->
          </v-chip>
          <FTaxonStatusIndicator :status="item.raw.status" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="Multiple extends boolean = false">
import { Taxon, TaxonRank } from '@/api'
import { listTaxaOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import { FTaxonStatusIndicator } from './functionals'
import TaxonChip from './TaxonChip'

type Multiplable<T> = true extends Multiple ? T[] : T

const model = defineModel<Multiplable<Taxon>>()

const {
  label = 'Taxon',
  // threshold = 0.7,
  // limit = 10,
  multiple = false,
  ranks
} = defineProps<{
  label?: string
  ranks?: TaxonRank[]
  threshold?: number
  limit?: number
  chips?: boolean
  multiple?: Multiple
}>()

const search = ref('')
const { data: items, isPending: loading, error } = useQuery(listTaxaOptions({ query: { ranks } }))

// TODO: Fix fix fuzzy search + highlight
// https://github.com/vuetifyjs/vuetify/issues/4417

// const keys: KeysDeclaration<TaxonWithParentRef> = ['name', 'rank', 'status']

// const { highlight, filteredItems } = useFuzzyItemsFilter(keys, search, items, {
//   threshold,
//   limit
// })
</script>

<style scoped></style>
