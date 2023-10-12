<template>
  <div>
    <v-form>
      <v-row>
        <v-col md="4">
          <v-text-field label="Name" density="compact"></v-text-field>
        </v-col>
        <v-col md="4">
          <v-select label="Rank" density="compact"></v-select>
        </v-col>
        <v-col md="4">
          <v-select label="Status" density="compact"></v-select>
        </v-col>
      </v-row>
    </v-form>
    <v-data-table :items="taxa" :headers="headers" density="compact"> </v-data-table>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios'
import type { Taxon } from '@/types/taxonomy'
import type { Ref } from 'vue'
import { ref } from 'vue'
import { onMounted } from 'vue'
import { VDataTable } from 'vuetify/labs/components'

import type { TaxonomyTaxonRank } from '@/api/data-contracts'

type UnwrapReadonlyArrayType<A> = A extends Readonly<Array<infer I>>
  ? UnwrapReadonlyArrayType<I>
  : A
type DT = InstanceType<typeof VDataTable>
type ReadonlyDataTableHeader = UnwrapReadonlyArrayType<DT['headers']>

const taxa: Ref<Taxon[]> = ref([])

const headers: ReadonlyDataTableHeader[] = [
  { title: 'Name', key: 'name' },
  { title: 'Code', key: 'code' },
  { title: 'Rank', key: 'rank' },
  { title: 'Status', key: 'status' }
]

type Filters = {
  name: string
  rank: string
}

async function fetch(): Promise<Taxon[]> {
  const response = await axios.get('/api/v1/taxa/', {
    params: {}
  })
  return response.data
}

onMounted(async () => {
  taxa.value = await fetch()
})
</script>

<style scoped></style>
