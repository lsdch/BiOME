<template>
  <v-select
    hide-details
    density="compact"
    label="Account status"
    placeholder="Any"
    persistent-placeholder
    clearable
    v-model="model.status"
    :items="statuses"
    color="primary"
    class="mt-1"
  >
    <template v-slot:item="{ props, item }">
      <v-list-item v-bind="props" density="compact">
        <template v-slot:prepend>
          <v-icon v-bind="icon(item)" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script lang="ts">
const statuses = ['Registered user', 'Unregistered', ...$UserRole.enum] as const
export type AccountStatus = (typeof statuses)[number]

export type PersonFilters = {
  term?: string
  status?: AccountStatus
}
</script>

<script setup lang="ts">
import { $UserRole } from '@/api'
import { getRoleIcon } from '@/components/icons/UserRoleIcon'

const model = defineModel<PersonFilters>({ required: true })

function icon(s: AccountStatus) {
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
      return getRoleIcon(s)
  }
}
</script>

<style scoped></style>
