<template>
  <v-container>
    <v-row>
      <v-col cols="12" sm="8">
        <v-inline-search-bar label="Search" v-model="model.term" />
      </v-col>
      <v-col cols="12" sm="4">
        <v-select
          hide-details
          density="compact"
          placeholder="All"
          clearable
          v-model="model.kind"
          :items="institutionKindOptions"
          color="primary"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
export type Filters = {
  term: string
  kind?: InstitutionKind
}
</script>

<script setup lang="ts">
import { institutionKindOptions } from './institutionKind'
import { InstitutionKind } from '@/api/models/InstitutionKind.js'

const model = defineModel<Filters>({ required: true })
</script>

<style scoped></style>
