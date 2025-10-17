<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchTerm"
    item-title="canonicalName"
    :items="autocompleteItems"
    :loading="loading"
    label="Taxon name"
    return-object
    auto-select-first
    color="blue"
    :no-data-text="
      searchTerm.trim().length >= 3 ? 'No matching taxa found' : 'At least 3 characters required'
    "
  >
    <template #item="{ props, item }">
      <v-list-item v-bind="props" :title="item?.raw?.canonicalName" :subtitle="item?.raw?.status">
        <template #append>
          <v-chip close>{{ item?.raw?.rank }}</v-chip>
        </template>
      </v-list-item>
    </template>
    <template #append-inner>
      <v-chip v-if="model">
        {{ model.rank }}
      </v-chip>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { endpointGBIF } from '.'
import { TaxonRank } from '@/api'
type TaxonGBIF = {
  key: number
  canonicalName: string
  authorship: string
  rank: string
  status: string
  path?: any
}
const model = ref<TaxonGBIF>()
const searchTerm = ref('')
const props = defineProps<{ rank?: TaxonRank | 'Any' }>()

type Item = {
  canonicalName: string
  rank: string
  status: string
  authorship: string
}
const autocompleteItems = ref<Item[]>([])
const loading = ref(false)

const excludedRanks = ['FORM', 'VARIETY', 'UNRANKED']

watch(searchTerm, async (val: string) => {
  if (val.length >= 3) {
    loading.value = true
    let response = await fetch(
      endpointGBIF('suggest', {
        query: {
          q: val,
          rank: props.rank !== 'Any' ? props.rank : undefined,
          status: 'ACCEPTED'
        }
      }).toString()
    )
    let json: Item[] = await response.json()
    autocompleteItems.value = json.filter(({ rank }) => !excludedRanks.includes(rank))
    loading.value = false
  } else {
    autocompleteItems.value = []
  }
})
</script>

<style lang="scss" scoped></style>
