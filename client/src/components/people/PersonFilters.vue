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
          v-model="model.status"
          :items="statuses"
          color="primary"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
const statuses = ['Registered user', 'Unregistered', ...orderedUserRoles] as const
export type PersonStatus = (typeof statuses)[number]

export type PersonFilters = {
  term: string
  status?: PersonStatus
}
</script>

<script setup lang="ts">
import { orderedUserRoles } from './userRole'

const model = defineModel<PersonFilters>({ required: true })
</script>

<style scoped></style>
