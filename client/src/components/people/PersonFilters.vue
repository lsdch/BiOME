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
        >
          <template v-slot:item="{ props, item }">
            <v-list-item v-bind="props" density="compact">
              <template v-slot:prepend>
                <v-icon v-bind="icon(item.value)" />
              </template>
            </v-list-item>
          </template>
        </v-select>
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
import { UserRole } from '@/api'
import { orderedUserRoles, roleIcon } from './userRole'

const model = defineModel<PersonFilters>({ required: true })

function icon(s: PersonStatus) {
  switch (s) {
    case 'Unregistered':
      return {
        icon: 'mdi-account',
        color: 'primary'
      }
    case 'Registered user':
      return {
        color: 'primary',
        icon: 'mdi-account-key'
      }
    default:
      return roleIcon(s as UserRole)
  }
}
</script>

<style scoped></style>
