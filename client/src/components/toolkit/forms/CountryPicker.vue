<template>
  <v-autocomplete
    v-model="model"
    label="Country"
    :items
    :item-value
    :multiple
    :return-object
    item-title="name"
    filter-mode="some"
    clear-on-select
    :loading="isPending || loading"
    :error-messages="error?.detail"
    :custom-filter="
      (_: any, q: string, item: InternalItem<CountrySummary> | undefined) => {
        if (q == '') return true
        if (!item) return false
        const { code, name } = item.raw
        return (
          code.toLowerCase().includes(q.toLowerCase()) ||
          name.toLowerCase().includes(q.toLowerCase())
        )
      }
    "
    v-bind="$attrs"
  >
    <template #item="{ item, props }">
      <v-list-item
        v-bind="props"
        :title="item.name"
        :subtitle="item.occurrence_count ? `${item.occurrence_count} occurrences` : undefined"
      >
        <template #append>
          <span class="text-overline">
            {{ item.code }}
          </span>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script
  setup
  lang="ts"
  generic="
    ItemValue extends 'code' | 'name',
    Multiple extends boolean,
    ReturnObject extends boolean
  "
>
import { Country, CountrySummary } from '@/api'
import { useCountries } from '@/stores/countries'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { InternalItem } from 'vuetify'
import { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'

type Model = ItemValue extends 'code' | 'name' ? Country[ItemValue] : Country

const model = defineModel<Value<CountrySummary, boolean, Multiple>>()

const { countries, isPending, error } = storeToRefs(useCountries())
const items = computed(() => {
  if (typeof filter === 'function') {
    return countries.value.filter(filter)
  }
  return countries.value
})

const {
  itemValue = 'code',
  loading,
  filter
} = defineProps<{
  returnObject?: ReturnObject
  itemValue?: ItemValue
  loading?: boolean
  multiple?: Multiple
  filter?: (item: CountrySummary) => boolean
}>()
</script>

<style scoped></style>
