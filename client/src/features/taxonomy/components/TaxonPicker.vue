<template>
  <v-autocomplete
    v-model="model"
    v-model:search="search"
    :label
    :loading
    :chips
    :items
    :multiple
    :return-object
    item-title="name"
    :item-value
    variant="outlined"
    clear-on-select
    placeholder="Enter search term..."
    :error-messages="error?.detail"
    class="taxon-picker"
    :list-props="{ class: 'position-relative' }"
    v-bind="$attrs"
  >
    <!-- no-filter -->
    <template #append-inner="props">
      <slot name="append-inner" v-bind="props" />
    </template>
    <template #prepend-inner="props">
      <slot name="prepend-inner" v-bind="props" />
    </template>
    <template #prepend-item>
      <slot name="prepend-item" />
    </template>
    <template #chip="{ item, props, index, internalItem }" v-if="items">
      <slot name="chip" :props :item>
        <TaxonChip :taxon="internalItem.raw" v-bind="props" />
      </slot>
    </template>
    <template #item="{ props, item }">
      <v-list-item v-bind="props" :title="item.name" class="fuzzy-search-item">
        <template #subtitle>
          {{ item.authorship }}
        </template>
        <template #append>
          <span class="text-muted text-caption">
            {{ item.rank }}
          </span>
          <FTaxonStatusIndicator :status="item.status" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="Multiple extends boolean, ReturnObject extends boolean">
import { Taxon, TaxonRank } from '@/api'
import { searchTaxaOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { keepPreviousData, useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'
import { FTaxonStatusIndicator } from './functionals'
import TaxonChip from './TaxonChip'

export type TaxonPickerProps<
  Multiple extends boolean = false,
  ReturnObject extends boolean = false
> = {
  label?: string
  ranks?: TaxonRank[] | TaxonRank
  threshold?: number
  limit?: number
  chips?: boolean
  multiple?: Multiple
  itemValue?: 'id' | 'name' | 'gbif_id'
  sampledOnly?: boolean
  returnObject?: ReturnObject
}
const model = defineModel<Value<Taxon, ReturnObject, Multiple> | null>()

const {
  label = 'Taxon',
  // threshold = 0.7,
  // limit = 10,
  multiple = false,
  ranks,
  sampledOnly,
  itemValue = 'id',
  returnObject = false
} = defineProps<TaxonPickerProps<Multiple, ReturnObject>>()

// const selectedTaxon = computed({
//   get() {
//     if (returnObject) {
//       return model.value as Value<Taxon, true, Multiple>
//     }
//     if (multiple) {
//       return ((_model.value as Taxon[] | undefined)?.filter((t) =>
//         model.value?.includes(t[itemValue])
//       ) ?? []) as Value<Taxon, true, true>
//     }
//     return _model.value as Value<Taxon, true, false>
//   },
//   set(taxon) {
//     _model.value = taxon as Value<Taxon, ReturnObject, Multiple>
//     console.log('selectedTaxon set', taxon)
//     if (returnObject) {
//       model.value = taxon as Value<Taxon, ReturnObject, Multiple>
//       return
//     }
//     if (multiple) {
//       console.log('multiple', taxon)
//       model.value = (taxon as Taxon[]).map((t) => t[itemValue]) as Value<
//         Taxon,
//         ReturnObject,
//         Multiple
//       >
//       console.log('model.value', model.value)
//       return
//     }
//     model.value = (taxon as Taxon)?.[itemValue] as Value<Taxon, ReturnObject, Multiple>
//   }
// }) as WritableComputedRef<Value<Taxon, true, Multiple>>

// const _model = ref<Value<Taxon, ReturnObject, Multiple>>()

const search = ref<string>()
const {
  data: items,
  isPending: loading,
  error
} = useQuery(
  computed(() => ({
    initialData: Array<Taxon>(),
    // enabled: (search.value?.length ?? 0) > 2,
    placeholderData: keepPreviousData,
    ...searchTaxaOptions({
      query: {
        // search_term: search.value,
        // limit: 20,
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
