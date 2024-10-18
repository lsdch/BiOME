<template>
  <v-select v-model="model" label="Role" :items="roles" color="primary" v-bind="$attrs">
    <template #prepend-inner>
      <v-icon v-bind="icon(model)" />
    </template>
    <template #item="{ props, item }">
      <v-list-item v-bind="props" density="compact" :subtitle="hints[item.raw]">
        <template #prepend>
          <v-icon v-bind="icon(item.value)" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { $UserRole, UserRole } from '@/api'
import { roleIcon } from '../people/userRole'

const model = defineModel<UserRole>({ required: true })
const roles = $UserRole.enum

function icon(item: UserRole) {
  return roleIcon(item as UserRole)
}

const hints: Record<UserRole, string> = {
  Visitor: 'Visitors have readonly access to the platform content',
  Contributor: 'Contributors may submit content and modify their own submissions',
  Maintainer: 'Maintainers have rights to manage most of the content',
  Admin: 'Admins have full read/write access'
} as const
</script>

<style scoped></style>
