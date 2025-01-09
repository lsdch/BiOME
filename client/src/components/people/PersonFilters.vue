<template>
  <v-row class="ma-0">
    <v-col cols="12" sm="6" class="pa-0">
      <v-list>
        <v-list-item>
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
                  <v-icon v-bind="icon(item.value)" />
                </template>
              </v-list-item>
            </template>
          </v-select>
        </v-list-item>
      </v-list>
    </v-col>
    <v-col cols="12" sm="6" class="pa-0">
      <v-list>
        <v-list-item>
          <InstitutionPicker
            v-model="model.institutions"
            clearable
            label="Institutions"
            placeholder="Any"
            persistent-placeholder
            density="compact"
            class="mt-1"
            multiple
            item-value="code"
          />
        </v-list-item>
      </v-list>
    </v-col>
  </v-row>
</template>

<script lang="ts">
const statuses = ['Registered user', 'Unregistered', ...$UserRole.enum] as const
export type AccountStatus = (typeof statuses)[number]

export type PersonFilters = {
  term?: string
  status?: AccountStatus
  institutions?: string[]
}
</script>

<script setup lang="ts">
import { $UserRole, UserRole } from '@/api'
import { roleIcon } from '../icons/UserRoleIcon'
import InstitutionPicker from './InstitutionPicker.vue'

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
      return roleIcon(s as UserRole)
  }
}
</script>

<style scoped></style>
