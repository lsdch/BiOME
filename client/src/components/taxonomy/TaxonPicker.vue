<template>
  <v-autocomplete
    v-model="model"
    v-model:search="search"
    :label
    :loading
    :chips
    :items="items as Taxon[] | undefined"
    :item-value
    :multiple
    :return-object
    item-title="name"
    no-filter
    variant="outlined"
    clear-on-select
    placeholder="Enter search term..."
    :error-messages="error?.detail"
    class="taxon-picker"
    :list-props="{ class: 'position-relative' }"
    v-bind="$attrs"
  >
    <template #prepend-inner="props">
      <slot name="prepend-inner" v-bind="props" />
    </template>
    <template #prepend-item>
      <slot name="prepend-item" />
    </template>
    <template #chip="{ item, props }" v-if="chips">
      <TaxonChip :taxon="item.raw" v-bind="props" />
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

<script
  setup
  lang="ts"
  generic="Multiple extends boolean, ReturnObject extends boolean, ItemValue extends keyof Taxon"
>
import { Taxon, TaxonRank, TaxonWithParentRef } from '@/api'
import { listTaxaOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'
import { FTaxonStatusIndicator } from './functionals'
import TaxonChip from './TaxonChip'

export type TaxonPickerProps<
  ItemValue extends keyof Taxon = 'name',
  Multiple extends boolean = false,
  ReturnObject extends boolean = false
> = {
  label?: string
  ranks?: TaxonRank[] | TaxonRank
  threshold?: number
  limit?: number
  chips?: boolean
  multiple?: Multiple
  itemValue?: ItemValue
  sampledOnly?: boolean
  returnObject?: ReturnObject
}
const model = defineModel<Value<Taxon, ReturnObject, Multiple>>()

const {
  label = 'Taxon',
  // threshold = 0.7,
  // limit = 10,
  multiple = false,
  ranks,
  sampledOnly
} = defineProps<TaxonPickerProps<ItemValue, Multiple, ReturnObject>>()

const search = ref('')
const {
  data: items,
  isPending: loading,
  error
} = useQuery(
  computed(() => ({
    initialData: Array<TaxonWithParentRef>(),
    ...listTaxaOptions({
      query: {
        pattern: search.value,
        limit: 20,
        ranks: Array.isArray(ranks) || !ranks ? ranks : [ranks],
        sampled_only: sampledOnly
      }
    })
  }))
)

// TODO: Fix fix fuzzy search + highlight
// https://github.com/vuetifyjs/vuetify/issues/4417

// const keys: KeysDeclaration<TaxonWithParentRef> = ['name', 'rank', 'status']

// const { highlight, filteredItems } = useFuzzyItemsFilter(keys, search, items, {
//   threshold,
//   limit
// })
</script>

<style lang="scss">
.taxon-picker .v-list {
  position: relative;
}
</style>
