<template>
  <v-select v-model="model" label="Role" :items="$UserRole.enum" color="primary" v-bind="$attrs">
    <template #prepend-inner>
      <UserRole.Icon :role="model" />
    </template>
    <template #item="{ props, item }">
      <v-list-item v-bind="props" density="compact" :subtitle="hints[item.raw]">
        <template #prepend>
          <UserRole.Icon :role="item" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { UserRole } from '@/api/adapters'
import { $UserRole } from '@/api'

const model = defineModel<UserRole>({ required: true })

const hints: Record<UserRole, string> = {
  Visitor: 'Visitors have readonly access to the platform content',
  Contributor: 'Contributors may submit content and modify their own submissions',
  Maintainer: 'Maintainers have rights to manage most of the content',
  Admin: 'Admins have full read/write access'
} as const
</script>

<style scoped></style>
