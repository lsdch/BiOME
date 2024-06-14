<template>
  <v-autocomplete
    v-bind="$attrs"
    label="Country"
    :items="items"
    item-title="name"
    chips
    filter-mode="some"
    :custom-filter="
      (_, q, item) => {
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
  </v-autocomplete>
</template>

<script setup lang="ts">
import { Country, LocationService } from '@/api'
import { ref } from 'vue'

// TODO: use composable with cached list
const items = ref(await LocationService.listCountries())
</script>

<style scoped></style>
