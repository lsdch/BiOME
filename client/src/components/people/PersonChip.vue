<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip :text="person.full_name" v-bind="{ ...props, ...$attrs }" />
    </template>
    <v-card
      :title="person.full_name"
      :subtitle="person.alias"
      class="bg-surface-light"
      density="compact"
    >
      <template #prepend>
        <div class="d-flex flex-column align-center">
          <UserRoleIcon :role="person.role" size="small" />
        </div>
      </template>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { PeopleService, PersonInner } from '@/api'
import UserRoleIcon from '../icons/UserRoleIcon.vue'

defineProps<{ person: PersonInner }>()
</script>

<style scoped lang="scss"></style>
