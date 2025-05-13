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
      (_: any, q: string, item: InternalItem<CountryWithSitesCount> | undefined) => {
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
        :title="item.raw.name"
        :subtitle="item.raw.sites_count ? `${item.raw.sites_count} sites` : undefined"
      >
        <template #append>
          <span class="text-overline">
            {{ item.raw.code }}
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
import { Country, CountryWithSitesCount } from '@/api'
import { useCountries } from '@/stores/countries'
import { storeToRefs } from 'pinia'
import { InternalItem } from 'vuetify'
import { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'

type Model = ItemValue extends 'code' | 'name' ? Country[ItemValue] : Country

const model = defineModel<Value<CountryWithSitesCount, boolean, Multiple>>()

const { countries: items, isPending, error } = storeToRefs(useCountries())

const { itemValue = 'code', loading } = defineProps<{
  returnObject?: ReturnObject
  itemValue?: ItemValue
  loading?: boolean
  multiple?: Multiple
}>()
</script>

<style scoped></style>
