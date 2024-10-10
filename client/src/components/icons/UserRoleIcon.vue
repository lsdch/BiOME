<template>
  <v-icon v-bind="{ ...$attrs, ...roleIcon(role) }" />
</template>

<script setup lang="ts">
import { UserRole } from '@/api'
import { roleIcon } from '../people/userRole'

defineProps<{ role: UserRole }>()
</script>

<style scoped></style>
