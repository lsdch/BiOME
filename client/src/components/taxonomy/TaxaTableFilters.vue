<template>
  <v-form>
    <v-container>
      <v-row>
        <v-col cols="12" sm="4">
          <v-text-field
            label="Name"
            :model-value="name"
            @update:model-value="debouncedUpdate"
            density="compact"
            clearable
            hide-details
            color="primary"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <v-select
            label="Rank"
            v-model="filters.rank"
            :items="taxonRankOptions"
            density="compact"
            clearable
            single-line
            hide-details
            variant="outlined"
            color="primary"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <v-select
            label="Status"
            v-model="filters.status"
            :items="taxonStatusOptions"
            density="compact"
            clearable
            single-line
            hide-details
            variant="outlined"
            color="primary"
          />
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>

<script setup lang="ts">
import { TaxonRank, TaxonStatus } from '@/api'
import { debounce as debounceFn } from 'vue-debounce'
import { taxonRankOptions, taxonStatusOptions } from './enums'

type TaxaFilters = {
  name?: string
  rank?: TaxonRank
  status?: TaxonStatus
}

const props = withDefaults(
  defineProps<{
    debounce?: number
  }>(),
  { debounce: 300 }
)

const name = defineModel<string>('name')

const debouncedUpdate = debounceFn((val?: string) => {
  name.value = val
}, props.debounce)

const filters = defineModel<TaxaFilters>({ default: {} })
</script>

<style scoped></style>
