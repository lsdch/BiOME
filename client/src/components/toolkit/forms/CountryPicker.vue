<template>
  <v-autocomplete
    v-model="model"
    v-bind="$attrs"
    label="Country"
    :item-value
    :return-object
    :items="items"
    item-title="name"
    filter-mode="some"
    :loading="isPending || loading"
    :error-messages="error?.detail"
    :custom-filter="
      (_: any, q: string, item: any) => {
        const { code, name }: Country = item?.raw
        if (q == '') return true
        return (
          code.toLowerCase().includes(q.toLowerCase()) ||
          name.toLowerCase().includes(q.toLowerCase())
        )
      }
    "
  >
    <template #item="{ item, props }">
      <v-list-item v-bind="props" :title="item.raw.name">
        <template #append>
          <span class="text-overline">
            {{ item.raw.code }}
          </span>
        </template>
      </v-list-item>
    </template>
    <template #append-inner>
      <v-chip v-if="model" :text="typeof model == 'string' ? model : model.code"></v-chip>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { Country } from '@/api'
import { useCountries } from '@/stores/countries'
import { storeToRefs } from 'pinia'

const model = defineModel<Country | string | undefined | null>({ required: true })

const { countries: items, isPending, error } = storeToRefs(useCountries())

withDefaults(
  defineProps<{
    returnObject?: boolean
    itemValue?: 'code' | 'name'
    loading?: boolean
  }>(),
  { itemValue: 'code' }
)
</script>

<style scoped></style>
