<template>
  <v-select v-model="model" label="Role" :items="$UserRole.enum" color="primary" v-bind="$attrs">
    <template #prepend-inner>
      <UserRole.Icon :role="model" />
    </template>
    <template #item="{ props, item }">
      <v-list-item v-bind="props" density="compact" :subtitle="hints[item]">
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
  visitor: 'Visitors have readonly access to the platform content',
  contributor: 'Contributors may submit content and modify their own submissions',
  maintainer: 'Maintainers have rights to manage most of the content',
  admin: 'Admins have full read/write access'
} as const
</script>

<style scoped></style>
